package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB es la variable global que contendrá la conexión a MySQL.
var MySQLDB *sql.DB

// InitMySQL establece la conexión a MySQL utilizando el DSN proporcionado.
// Ejemplo de DSN: "usuario:contraseña@tcp(localhost:3306)/nombre_basedatos"
func InitMySQL(dsn string) error {
	var err error
	MySQLDB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Verifica la conexión.
	if err = MySQLDB.Ping(); err != nil {
		return err
	}

	log.Println("Conexión a MySQL establecida correctamente")
	return nil
}
