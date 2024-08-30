package main

import (
	"com.github/EnesHarman/url-shortener/config"
	"com.github/EnesHarman/url-shortener/internal/controllers"
	"com.github/EnesHarman/url-shortener/internal/db"
	"com.github/EnesHarman/url-shortener/internal/repository"
	"com.github/EnesHarman/url-shortener/internal/services"
	gin "github.com/gin-gonic/gin"
	"golang.org/x/net/context"
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

	pool := db.GetConnectionPool(context.Background(), configs.Postgresql)
	repository := repository.NewUrlRepository(pool)
	service := services.NewUrlService(repository, configs.UrlShortener)
	controller := controllers.NewUrlController(service)
	controller.RegisterRoutes(e)
}
