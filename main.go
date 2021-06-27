package main

import (
	"context"
	"encoding/json"
	"github.com/ervitis/exchangerateapp/exchangeclient"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type (
	handler struct {
		client *exchangeclient.ExchangeApi
	}
)

func (h *handler) healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}
}

func (h *handler) ratingConverterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/rate/from/")
		if p == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fromCurrency := strings.TrimRight(p, "?")
		to := r.URL.Query().Get("to")
		if fromCurrency == "" || to == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		toCurrencies := strings.Split(to, ",")

		cc, err := h.client.ConvertCurrency(fromCurrency, toCurrencies)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}
		b, _ := json.Marshal(cc)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func main() {
	s := make(chan os.Signal, 2)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	ec, err := exchangeclient.NewClient()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{client: ec}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.healthCheckHandler())
	mux.HandleFunc("/rate/from/{fromCurrency}?to={toCurrencies}", h.ratingConverterHandler())

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func(server *http.Server) {
		log.Fatal(server.ListenAndServe())
	}(server)

	<- s
	log.Println("Terminating...")
}
