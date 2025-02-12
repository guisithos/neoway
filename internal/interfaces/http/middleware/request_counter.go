package middleware

import (
	"log"
	"net/http"

	"github.com/guisithos/neoway/internal/application/services"
)

func RequestCounter(metricsService services.MetricsService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Contador assincrono
			go func() {
				err := metricsService.IncrementRequestCount(r.Method, r.URL.Path)
				if err != nil {
					// Log o erro mas não afeta a req
					log.Printf("erro ao incrementar o contador de requisições: %v", err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
