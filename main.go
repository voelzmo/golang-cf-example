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
	DEFAULT_PORT     = "8080"
	DEFAULT_APP_NAME = "my-app"
)

type Page struct {
	Title    string
	AppName  string
	AppIndex int
}

func HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/home.html")

	appEnv, _ := cfenv.Current()
	p := Page{
		Title:    "Welcome",
		AppName:  appEnv.Name,
		AppIndex: appEnv.Index,
	}
	t.Execute(responseWriter, p)
}

func main() {
	var (
		port    string
		appName string
	)

	appEnv, _ := cfenv.Current()
	if appName = appEnv.Name; len(appName) == 0 {
		appName = DEFAULT_APP_NAME
	}

	logger := lager.NewLogger(appName)
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.DEBUG)
	logger.RegisterSink(sink)

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
