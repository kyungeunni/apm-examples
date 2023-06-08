package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
)

func write(writer http.ResponseWriter, message string) {
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func englishHandler(writer http.ResponseWriter, request *http.Request) {
	// make outgoing http requests to add more spans
	resp, err := ctxhttp.Get(request.Context(), tracingClient, "https://kyungeun.kim")
	if err != nil {
		apm.CaptureError(request.Context(), err).Send()
		http.Error(writer, "failed to query", 500)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		apm.CaptureError(request.Context(), err).Send()
		http.Error(writer, "failed to parse body", 500)
		return
	}
	write(writer, string(body[:]))
}

func frenchHandler(writer http.ResponseWriter, request *http.Request) {
	span, _ := apm.StartSpan(request.Context(), "Delay", "request")
	defer span.End()
	time.Sleep(3 * time.Second)
	span.Context.SetLabel("reason", "time.Sleep")
	write(writer, "Salut, web!")
}

func hindiHandler(writer http.ResponseWriter, request *http.Request) {
	write(writer, "Namaste, web!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", englishHandler)
	mux.HandleFunc("/salut", frenchHandler)
	mux.HandleFunc("/namaste", hindiHandler)
	err := http.ListenAndServe("localhost:8080", apmhttp.Wrap(mux))
	if err != nil {
		log.Fatal(err)
	}
}
