package main

import (
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

	// Construir el DSN para MySQL usando variables de entorno
	mysqlDSN := os.Getenv("MYSQL_USER") + ":" +
		os.Getenv("MYSQL_PASSWORD") + "@tcp(" +
		os.Getenv("MYSQL_HOST") + ")/" +
		os.Getenv("MYSQL_DATABASE")

	// Inicializar la conexión a MySQL
	if err := db.InitMySQL(mysqlDSN); err != nil {
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
