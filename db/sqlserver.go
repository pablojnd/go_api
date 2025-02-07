package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

// SQLServerDB es la variable global que contendrá la conexión a SQL Server.
var SQLServerDB *sql.DB

// InitSQLServer establece la conexión a SQL Server utilizando las variables de entorno.
func InitSQLServer() error {
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s", server, port, user, password, database)
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	SQLServerDB = db
	log.Println("Conexión a SQL Server establecida correctamente")
	return nil
}
