package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/mux"
)

const (
	DEFAULT_PORT = "8080"
)

type Page struct {
	Title string
	Body  string
}

func HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/home.html")
	p := Page{
		Title: "Welcome",
		Body:  "Hello, World!",
	}
	t.Execute(responseWriter, p)
}

func main() {
	var port string
	logger := lager.NewLogger("my-app")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.DEBUG)
	logger.RegisterSink(sink)

	appEnv, _ := cfenv.Current()
	logger.Info("Just logging out some application environment info", lager.Data{"appEnv": fmt.Sprintf("%+v", appEnv)})

	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	srv := http.Server{
		Handler:        r,
		Addr:           fmt.Sprintf(":%s", port),
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info("Starting server")
	srv.ListenAndServe()

}
