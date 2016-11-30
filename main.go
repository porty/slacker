package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	t, err := getTemplates()
	if err != nil {
		log.Fatal(err)
	}

	s := Slack{config.SlackMessageURL}

	h := Handler{
		slack:     &s,
		templates: t,
		channels:  config.SlackChannels,
		username:  config.Username,
		password:  config.Password,
	}
	h.Register(http.DefaultServeMux)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	bail := make(chan struct{})

	go func() {
		log.Printf("Listening on http://0.0.0.0:%d/", config.Port)
		addr := fmt.Sprintf(":%d", config.Port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Print(err)
		}
		close(bail)
	}()

	select {
	case <-c:
		log.Print("Received signal, quitting")
	case <-bail:
		log.Print("Must have had problems with listening for HTTP")
	}

}

func getTemplates() (*template.Template, error) {
	b, err := Asset("templates/form.html")
	if err != nil {
		return nil, err
	}
	return template.New("form").Parse(string(b))
	//return template.ParseFiles("form.html")
}
