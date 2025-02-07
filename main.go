package main

import (
	"fmt"
	"go_api/db"
	"go_api/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Intentar cargar .env pero no fallar si no existe
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Inicializar conexión a SQL Server
	if err := db.InitSQLServer(); err != nil {
		log.Printf("Error al inicializar SQL Server: %v", err)
		return
	}
	defer db.SQLServerDB.Close()

	// Inicializar la conexión a MySQL
	if err := db.InitMySQL(); err != nil {
		log.Printf("Error al inicializar MySQL: %v", err)
		return
	}
	defer db.MySQLDB.Close()

	// Configurar rutas centralizadas
	routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
