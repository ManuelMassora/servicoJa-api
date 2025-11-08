package main

import (
	// "context"
	// "fmt"
	"log"
	// "net/http"

	// "github.com/gin-contrib/cors"
	"github.com/ManuelMassora/servicoJa-api/internal/config"
	"github.com/ManuelMassora/servicoJa-api/internal/db"
	"github.com/ManuelMassora/servicoJa-api/internal/di"
	"github.com/ManuelMassora/servicoJa-api/internal/routes"
	"github.com/gin-gonic/gin"
	// "github.com/markbates/goth/gothic"
)

type appConfig struct {
	*config.Config
}

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	server := gin.Default()

	// server.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           86400,
	// }))

	db, err := db.InitDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	container := di.NewContainer(db)
	routes.SetRoutes(server, container)

	log.Printf("Servidor rodando na porta :%s", cfg.ServerPort)
	if err := server.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}