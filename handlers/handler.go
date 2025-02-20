package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./web/template/" + tmpl + ".html")

	if err != nil {
		return
	}

	if err := t.Execute(w, nil); err != nil {
		return
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	switch r.Method {
	case "GET":
		RenderTemplate(w, "index")
		log.Printf("HTTP Response Code : %v", (http.StatusOK))
	default:
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reset-password" {
		return
	}

	switch r.Method {
	case "GET":
		RenderTemplate(w, "resetPassword")
		log.Printf("HTTP Response Code : %v", (http.StatusOK))
	default:
	}
}

func ResetConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reset-confirmation" {
		return
	}

	switch r.Method {
	case "GET":
		RenderTemplate(w, "resetConfirmation")
		log.Printf("HTTP Response Code : %v", (http.StatusOK))
	default:
	}
}
