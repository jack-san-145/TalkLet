let alpineVars
  document.addEventListener("alpine:initialized",()=>{
    const el=document.querySelector("#ui-variables")
    alpineVars = Alpine.$data(el)
    console.log(alpineVars)
    console.log(Alpine.store('globalVar').chatlist)
    // console.log("users - ",alphineVars.users[0].name)
    console.log("ischat - ",alpineVars.IsChat)
    console.log("active_chat from js  - ",Alpine.store('globalVar').active_chat)
    Alpine.store('globalVar').chatlist.map((contact)=>{
      contact.created_at=convertTimeForChatList(contact.created_at)
    })
  })  

  function activeChatfunc(contact){
    LoadOlderMessages(contact.contact_id); 
    alpineVars.IsChat=true; 
    Alpine.store('globalVar').active_chat = contact.contact_id;
    alpineVars.active_user = contact;
    console.log('active chat - ',Alpine.store('globalVar').active_chat)
  }

  let input_msg
  const socket=new WebSocket(`ws://${window.location.hostname}:8989/ws`)
  const search_input=document.getElementsByClassName("search-input")
  socket.onopen=()=>{
    console.log("websocket connected")
  }
  socket.onclose=()=>{
    console.log("websocket closed")
  }

  function addLastMsg(receiver_id,last_msg,time){
    const contact = Alpine.store('globalVar').chatlist.find((contact => contact.contact_id == receiver_id ))
    console.log("contact for last msg - ",contact,receiver_id)
    contact.last_msg=last_msg
    contact.created_at=time
    
  }

  function msgYou(received_msg){
      const chat_messages=document.getElementById("chat-messages")
      const send_you_div=document.createElement("div")
      send_you_div.classList="message you"
      const span1=document.createElement("span")
      span1.classList="text"
      span1.textContent=received_msg.content
      const span2=document.createElement("span")
      span2.classList="msg-you-time"
      console.log("before calculated - ",received_msg.created_at)
      span2.textContent=calculateTimeforMsg(received_msg.created_at)
      send_you_div.appendChild(span1)
      send_you_div.appendChild(span2)
      chat_messages.appendChild(send_you_div)
      addLastMsg(received_msg.receiver_id,received_msg.content,convertTimeForChatList(received_msg.created_at))
      scrollBottom()
      addToFirst(received_msg.receiver_id)
      console.log("received_msg.id - msg you",received_msg.receiver_id)
      console.log('chat_messages - ',chat_messages)
      console.log('msd div - ',send_you_div)
  }

  function msgFrd(received_msg){
    const chat_messages=document.getElementById("chat-messages")
    const send_frd_div=document.createElement("div")
    send_frd_div.classList="message friend"
    const span1=document.createElement("span")
    span1.classList="text"
    span1.textContent=received_msg.content
    const span2=document.createElement("span")
    span2.classList="msg-frd-time"
    span2.textContent=calculateTimeforMsg(received_msg.created_at)
    send_frd_div.appendChild(span1)
    send_frd_div.appendChild(span2)
    chat_messages.appendChild(send_frd_div)
    addLastMsg(received_msg.sender_id,received_msg.content,convertTimeForChatList(received_msg.created_at))
    scrollBottom()
    addToFirst(received_msg.sender_id)
    console.log("received_msg.id - msg frd",received_msg.sender_id)
  }

  socket.onmessage=function(event)
  {
    const received_msg=JSON.parse(event.data)
    console.log("received msg from server - ",received_msg)
    if(received_msg.is_ack=="ack")
    {
      msgYou(received_msg)
    }else{
      msgFrd(received_msg)
    }
    
  }
  function addToFirst(id){
    console.log("to first id - ",id)
    const tofirst=document.getElementById(String(id))
    console.log("tofirst - ",tofirst)
    const chat_list=document.getElementsByClassName("chat-list")[0]
    console.log(chat_list)
    chat_list.prepend(tofirst)
    
    console.log("findUserById - ",findUserById(id))
  }

  function findUserById(id){
   const index = id-1
   console.log(Alpine.store('globalVar').chatlist[index])
  }

  function sendThisMSg()
  {
    console.log("sending msg working .. ")
    const chat_messages=document.getElementById("chat-messages")
    const input_msg=document.getElementById("input-msg")
    if (input_msg.type="text"){
      console.log("true it is text - ",input_msg.type)
    }
    console.log("active_chat js  - ",Alpine.store('globalVar').active_chat)
    console.log(input_msg.value)
    const send_msg={
      "receiver_id" :  Alpine.store('globalVar').active_chat,
      "content" : input_msg.value,
      "type" : input_msg.type
    }
    if(input_msg.value!=""){
      socket.send(JSON.stringify(send_msg))
      console.log("sent before ack - ",send_msg)
      input_msg.value=""
     
    }
  }

 

  
  function calculateTimeforMsg(goTimeString) {
  // Remove microseconds and timezone name
  const cleaned = goTimeString
    .replace(/\.\d+/, '')        // remove `.123456`
    .replace(/ [A-Z]+$/, '');    // remove `IST`

  const date = new Date(cleaned);

  if (isNaN(date)) return 'Invalid Date';

  // Format to "hh:mm AM/PM"
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  });
}


  function scrollBottom()
  {
    const chat=document.getElementById("chat-messages")
    chat.scrollTop=chat.scrollHeight
  }

  function listenInput(){

    input_msg=document.getElementById("input-msg")
    input_msg.addEventListener("keypress",(event)=>{
      console.log("listening input . .")
      if(event.key=="Enter"){
        sendThisMSg()
      }
    })
  }

  function convertTimeForChatList(goTimeString) {
  // Step 1: Clean the Go timestamp
      const cleaned = goTimeString
        .replace(/\.\d+/, '')        // remove microseconds
        .replace(/ [A-Z]+$/, '');    // remove timezone label like IST

      const date = new Date(cleaned);
      if (isNaN(date)) return 'Invalid Date';

      // Step 2: Reference today and yesterday
      const now = new Date();
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
      const yesterday = new Date(today);
      yesterday.setDate(today.getDate() - 1);

      const inputDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());

      // Step 3: Compare
      if (inputDate.getTime() === today.getTime()) {
        return 'Today';
      } else if (inputDate.getTime() === yesterday.getTime()) {
        return 'Yesterday';
      } else {
        const mm = String(date.getMonth() + 1).padStart(2, '0');
        const dd = String(date.getDate()).padStart(2, '0');
        const yy = String(date.getFullYear()).slice(-2);
        return `${mm}/${dd}/${yy}`;
      }
}

function search(){
  Alpine.store('globalVar').chatlist.filter()
}

function handleScroll(event){
  const el=event.target
  console.log('scroll chat - ',Alpine.store('globalVar').active_chat)
  if(el.scrollTop==0){
    console.log('reached top - ')
    h()
  } 
}
