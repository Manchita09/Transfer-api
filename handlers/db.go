package handlers

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB inicializa la conexión a PostgreSQL
func InitDB(connStr string) {
	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión:", err)
	}

	// Verificar conexión real
	if err = db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}

	// Configuración del pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Conectado a PostgreSQL")
}

// GetDB devuelve la instancia de la base de datos
func GetDB() *sql.DB {
	return db
}
