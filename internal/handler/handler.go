package handler

import (
	"net/http"
	"sycretru/internal/service"
	"sycretru/pkg/logger"

	"github.com/gorilla/mux"
)

type handler struct {
	svc service.Gendocer
	log *logger.Log
}

func NewHandler(svc *service.Service, log *logger.Log) *handler {
	return &handler{
		svc: svc,
		log: log,
	}
}

func (h *handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/gendoc", h.gendoc).Methods("GET", "POST")
	r.HandleFunc("/upload/word/{dir}/{file}", h.upload).Methods("GET")
	// r.HandleFunc("/upload/pdf/{dir}/{file}", h.upload).Methods("GET")

	r.Use(Logging)

	return r
}
