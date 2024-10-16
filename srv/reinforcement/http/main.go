package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/ollama/ollama/api"
	"github.com/rs/cors"
)

type FraudRequest struct {
	Username        string `json:"username"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	CompanyName     string `json:"company_name"`
	CompanyIndustry string `json:"company_industry"`
	CompanyAddress  string `json:"company_address"`
	CompanyCity     string `json:"company_city"`
	CompanyState    string `json:"company_state"`
	CompanyZipcode  string `json:"company_zipcode"`
}

type FraudResponse struct {
	IsSpam          bool   `json:"is_spam"`
	Reason          string `json:"reason"`
	ConfidenceLevel string `json:"confidence_level"`
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func HttpHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			HandleErrorResponse(w, err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func HandleErrorResponse(w http.ResponseWriter, e error) error {
	return WriteJSON(w, http.StatusBadRequest, e)
}

func CheckForFraud(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	var req FraudRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return HandleErrorResponse(w, err)
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	return WriteJSON(w, http.StatusOK, "OK")
}

func registerRoutes() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/fraud", HttpHandlerFunc(CheckForFraud)).Methods("POST")
	return r
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "graceful timeout")
	flag.Parse()
	port := "3001"
	routes := registerRoutes()
	crossOrigin := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})
	handler := crossOrigin.Handler(routes)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		slog.Info("server running on port :3001")
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	server.Shutdown(ctx)
	slog.Info("shutting down application")
	os.Exit(0)
}
