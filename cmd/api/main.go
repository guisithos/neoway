package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/guisithos/neoway/internal/application/services"
	"github.com/guisithos/neoway/internal/domain/models"
	"github.com/guisithos/neoway/internal/infrastructure/repository"
	"github.com/guisithos/neoway/internal/interfaces/http/handlers"
	"github.com/guisithos/neoway/internal/interfaces/http/routes"
)

func main() {
	// db conn
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	// db retry conn
	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}
		log.Printf("falhou a conexão com o banco de dados (tentativa %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		log.Fatal("falha ao conectar ao banco de dados após várias tentativas:", err)
	}

	// Migrate
	err = db.AutoMigrate(&models.Client{}, &models.RequestMetrics{})
	if err != nil {
		log.Fatal("falha ao migrar o banco de dados:", err)
	}

	// Inicia repos e servicos
	clientRepo := repository.NewClientRepository(db)
	metricsRepo := repository.NewRequestMetricsRepository(db)

	clientService := services.NewClientService(clientRepo)
	metricsService := services.NewMetricsService(metricsRepo)

	// Inicia handlers
	clientHandler := handlers.NewClientHandler(clientService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)

	// Configura router com middleware e rotas
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Configura rotas
	routes.SetupRoutes(r, clientHandler, metricsHandler, metricsService)

	// Inicia serv
	log.Println("Servidor iniciando na porta :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
