package services

import (
	"sync"
	"time"

	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/guisithos/neoway/internal/infrastructure/repository"
)

type MetricsService interface {
	GetUptime() time.Duration
	GetStartTime() time.Time
	GetLastRestart() *time.Time
	RecordRestart()
	IncrementRequestCount(method, path string) error
	GetRequestMetrics() ([]models.RequestMetrics, error)
}

type metricsService struct {
	startTime   time.Time
	lastRestart *time.Time
	mu          sync.RWMutex
	metricsRepo repository.RequestMetricsRepository
}

func NewMetricsService(metricsRepo repository.RequestMetricsRepository) MetricsService {
	return &metricsService{
		startTime:   time.Now(),
		metricsRepo: metricsRepo,
	}
}

func (s *metricsService) GetUptime() time.Duration {
	return time.Since(s.startTime)
}

func (s *metricsService) GetStartTime() time.Time {
	return s.startTime
}

func (s *metricsService) GetLastRestart() *time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastRestart
}

func (s *metricsService) RecordRestart() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.lastRestart = &now
}

// Registra o horário da última reinicialização do serviço
func (s *metricsService) IncrementRequestCount(method, path string) error {
	return s.metricsRepo.IncrementCount(method, path)
}

// Obtém todas as métricas de requisições
func (s *metricsService) GetRequestMetrics() ([]models.RequestMetrics, error) {
	return s.metricsRepo.GetAllMetrics()
}
