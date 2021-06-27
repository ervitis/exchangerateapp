package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ervitis/exchangerateapp/exchangeclient"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

type (
	handler struct {
		client        *exchangeclient.ExchangeApi
		ratingPattern *regexp.Regexp
	}

	responseConvertedCurrency map[string]string
)

const (
	quantityIndex = "quantity"
	fromCurrencyIndex = "fromCurrency"
)

var (
	urlPattern = fmt.Sprintf(`\/rate\/quantity\/(?P<%s>[0-9]*[.]?[0-9]+)\/from\/(?P<%s>\w+)`, quantityIndex, fromCurrencyIndex)
)

func (h *handler) healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`ok`))
	}
}

func (h *handler) ratingConverterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := h.ratingPattern.FindStringSubmatch(r.URL.Path)
		params := make(map[string]string, 0)

		for i, v := range h.ratingPattern.SubexpNames() {
			if i > 0 && i <= len(p) {
				params[v] = p[i]
			}
		}

		toCurrencies := strings.Split(r.URL.Query().Get("to"), ",")

		cc, err := h.client.ConvertCurrency(params[fromCurrencyIndex], toCurrencies)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}

		quantity, err := strconv.ParseFloat(params[quantityIndex], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}

		rcc := make(responseConvertedCurrency, 0)
		for k, v := range cc {
			rcc[k] = fmt.Sprintf(`%.2f %s`, v * quantity, params[fromCurrencyIndex])
		}

		b, _ := json.Marshal(rcc)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func main() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ec, err := exchangeclient.NewClient()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{client: ec, ratingPattern: regexp.MustCompile(urlPattern)}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.healthCheckHandler())
	mux.HandleFunc("/rate/quantity/", h.ratingConverterHandler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func(server *http.Server) {
		log.Println("Listening at http://localhost:8080")
		log.Fatal(server.ListenAndServe())
	}(server)

	<-s
	log.Println("Terminating...")
}
