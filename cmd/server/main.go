package main

import (
	"com.github/EnesHarman/url-shortener/config"
	"com.github/EnesHarman/url-shortener/internal/controllers"
	"com.github/EnesHarman/url-shortener/internal/db"
	"com.github/EnesHarman/url-shortener/internal/kafka"
	"com.github/EnesHarman/url-shortener/internal/repository"
	"com.github/EnesHarman/url-shortener/internal/services"
	gin "github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configs, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	g := gin.Default()
	generateUrlController(g, configs)
	g.Run(":" + configs.Server.Port)

}

func generateUrlController(e *gin.Engine, configs *config.Config) { // I miss Spring Boot
	clickEventProducer := kafka.NewClickEventProducer()
	pool := db.GetConnectionPool(context.Background(), configs.Postgresql)
	repository := repository.NewUrlRepository(pool)
	service := services.NewUrlService(repository, configs.UrlShortener, *clickEventProducer)
	controller := controllers.NewUrlController(service)
	controller.RegisterRoutes(e)

	shutDownProducers(clickEventProducer)
}

func shutDownProducers(clickEventProducer *kafka.ClickEventProducer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutdown signal received, closing producer...")
		if err := clickEventProducer.ShutDown(); err != nil {
			log.Fatalf("Failed to shutdown producer: %v", err)
		}
		os.Exit(0)
	}()
}
