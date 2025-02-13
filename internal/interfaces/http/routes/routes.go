package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guisithos/neoway/internal/application/services"
	"github.com/guisithos/neoway/internal/interfaces/http/handlers"
	"github.com/guisithos/neoway/internal/interfaces/http/middleware"
)

func SetupRoutes(r *chi.Mux, clientHandler *handlers.ClientHandler, metricsHandler *handlers.MetricsHandler, metricsService services.MetricsService) {
	// add o middleware de contador de reqs a todas as rotas
	r.Use(middleware.RequestCounter(metricsService))

	r.Post("/clients", clientHandler.CreateClient)
	r.Get("/clients", clientHandler.ListClients)
	r.Get("/clients/document/{document}", clientHandler.GetClientByDocument)
	r.Get("/clients/name/{name}", clientHandler.GetClientsByName)

	// add o middleware de debug para as rotas de block/unblock
	r.Put("/clients/{document}/block", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Route params: %+v", chi.URLParam(r, "document"))
		clientHandler.BlockClient(w, r)
	})
	r.Put("/clients/{document}/unblock", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Route params: %+v", chi.URLParam(r, "document"))
		clientHandler.UnblockClient(w, r)
	})

	// rotas de metrics
	r.Route("/metrics", func(r chi.Router) {
		r.Get("/uptime", metricsHandler.GetUptime)
		r.Get("/requests", metricsHandler.GetRequestMetrics)
	})
}
