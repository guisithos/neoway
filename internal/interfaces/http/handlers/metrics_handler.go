package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/guisithos/neoway/internal/application/services"
)

type MetricsHandler struct {
	metricsService services.MetricsService
}

type UptimeResponse struct {
	Uptime       string     `json:"uptime"`
	StartTime    time.Time  `json:"start_time"`
	LastRestart  *time.Time `json:"last_restart,omitempty"`
	UptimeString string     `json:"uptime_string"`
}

func NewMetricsHandler(metricsService services.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: metricsService}
}

func (h *MetricsHandler) GetUptime(w http.ResponseWriter, r *http.Request) {
	uptime := h.metricsService.GetUptime()
	startTime := h.metricsService.GetStartTime()
	lastRestart := h.metricsService.GetLastRestart()

	response := UptimeResponse{
		Uptime:       uptime.String(),
		StartTime:    startTime,
		LastRestart:  lastRestart,
		UptimeString: formatUptime(uptime),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func (h *MetricsHandler) GetRequestMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.metricsService.GetRequestMetrics()
	if err != nil {
		http.Error(w, "erro ao buscar métricas de requisições", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
