package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
)

type Handler struct {
	slack     *Slack
	templates *template.Template
	channels  []string
	username  string
	password  string
}

type Mux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
}

func (h *Handler) Register(mux Mux) {
	mux.Handle("/", authMiddleware(h.username, h.password, "slacker", http.HandlerFunc(h.rootHandler)))
}

type formViewData struct {
	Channels       []string
	Message        SlackMessage
	ErrorMessage   string
	SuccessMessage string
}

func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	data := formViewData{
		Channels: h.channels,
	}
	if r.Method == http.MethodPost {
		h.handlePost(w, r, &data)
	}

	w.Header().Set("Content-Type", "text/html")

	err := h.templates.ExecuteTemplate(w, "form", data)

	if err != nil {
		log.Print("Error rendering template: " + err.Error())
	}
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request, data *formViewData) {
	decoder := schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		data.ErrorMessage = err.Error()
		return
	}

	if err := decoder.Decode(&data.Message, r.Form); err != nil {
		data.ErrorMessage = err.Error()
		return
	}

	if err := h.slack.Send(data.Message); err != nil {
		data.ErrorMessage = "Failed to send Slack message: " + err.Error()
		return
	}
	data.SuccessMessage = "Message sent successfully"
}

func authMiddleware(username string, password string, realm string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get("Authorization")
		if !strings.HasPrefix(v, "Basic ") {
			rejectAuth(realm, w)
			return
		}
		v = v[len("Basic "):]
		bytes, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return
		}
		v = string(bytes)
		creds := strings.SplitN(v, ":", 2)
		if (len(creds) == 2) && creds[0] == username && creds[1] == password {
			next.ServeHTTP(w, r)
		} else {
			rejectAuth(realm, w)
		}
	})
}

func rejectAuth(realm string, w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}
