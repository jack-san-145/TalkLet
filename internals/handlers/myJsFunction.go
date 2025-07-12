package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
)

var MyJsonConvFunc = template.FuncMap{
	"tojson": func(v any) template.JS {
		data, err := json.Marshal(v)
		if err != nil {
			fmt.Println("error while mapping tojson function - ", err)
			return ""
		}
		return template.JS(data)

	},
}
