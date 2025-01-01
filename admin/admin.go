package admin

import (
	"net/http"
	"html/template"
)

func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/admin.html")
		tmpl.Execute(w, map[string]interface{}{})
	}
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	
}