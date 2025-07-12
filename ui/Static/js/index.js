const input_msg=document.getElementById("input-msg")
  const socket=new WebSocket("ws://localhost:8989/ws")
  socket.onmessage=function(event)
  {
    const chat_messages=document.getElementById("chat-messages")
    console.log(chat_messages)
    const received_msg=event.data
    console.log("received msg = ",received_msg)
    console.log(chat_messages)
    const received_msg_div=document.createElement("div")
    received_msg_div.textContent=received_msg
    received_msg_div.classList="message friend"
    chat_messages.appendChild(received_msg_div)
    scrollBottom()
  }


  function sendThisMSg()
  {
    const chat_messages=document.getElementById("chat-messages")
    console.log(chat_messages)
    const input_msg=document.getElementById("input-msg")
    const send_msg=input_msg.value
    input_msg.value=""
    if(send_msg=="")
    {
      console.log("enter nothing")
    }else{
      socket.send(send_msg)
      console.log(send_msg)
      const send_msg_div=document.createElement("div")
      send_msg_div.textContent=send_msg
      send_msg_div.classList="message you"
      chat_messages.appendChild(send_msg_div)
      scrollBottom()
    }
  }

  function scrollBottom()
  {
    const chat=document.getElementById("chat-messages")
    chat.scrollTop=chat.scrollHeight
  }

  input_msg.addEventListener("keypress",(event)=>{
    if(event.key=="Enter"){
      sendThisMSg()
    }
  })
  input_msg.addEventListener("blur",()=>{
    setTimeout(()=>input_msg.focus(),10)
  })
