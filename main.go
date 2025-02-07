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
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Inicializar conexión a SQL Server
	if err := db.InitSQLServer(); err != nil {
		log.Fatalf("Error al inicializar SQL Server: %v", err)
	}
	defer db.SQLServerDB.Close()

	// Obtener las variables de entorno
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	// Si no se proporciona el host, se usa por defecto "localhost:3306"
	if mysqlHost == "" {
		mysqlHost = "localhost:3306"
	}

	// Construir el DSN para MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlDatabase)

	// Inicializar la conexión a MySQL
	if err := db.InitMySQL(dsn); err != nil {
		log.Fatal("Error al inicializar conexión a MySQL:", err)
	}
	defer db.MySQLDB.Close()

	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configurar rutas centralizadas
	routes.SetupRoutes()

	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
