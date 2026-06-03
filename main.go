package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Manchita09/transfer-api/handlers"
)

func main() {

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=127.0.0.1 port=5432 user=postgres password=1234 dbname=TransferMarketDB sslmode=disable"
	}

	handlers.InitDB(connStr)

	router := mux.NewRouter().StrictSlash(true)

	seedDefaultUser()

	// CORS
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// rutas
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	router.HandleFunc("/players", handlers.GetPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", handlers.GetPlayerByID).Methods("GET")
	router.HandleFunc("/players", handlers.CreatePlayer).Methods("POST")
	router.HandleFunc("/players/{id}", handlers.UpdatePlayer).Methods("PUT")
	router.HandleFunc("/players/{id}", handlers.DeletePlayer).Methods("DELETE")
	router.HandleFunc("/transfer", handlers.TransferPlayer).Methods("POST")
	router.HandleFunc("/teams", handlers.GetTeams).Methods("GET")
	router.HandleFunc("/teams/{id}", handlers.GetTeamByID).Methods("GET")
	router.HandleFunc("/ranking", handlers.GetTeamRanking).Methods("GET")
	router.HandleFunc("/transfers", handlers.GetTransfers).Methods("GET")
	router.HandleFunc("/reset-transfers", handlers.ResetTransfers).Methods("DELETE")
	router.HandleFunc("/top-players", handlers.GetTopPlayers).Methods("GET")

	// static
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))),
	)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	router.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Servidor corriendo en http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func seedDefaultUser() {
	db := handlers.GetDB()

	_, err := db.Exec(`
	INSERT INTO users (username, password)
	SELECT 'admin', '1234'
	WHERE NOT EXISTS (
		SELECT 1 FROM users WHERE username = 'admin'
	)
	`)
	if err != nil {
		log.Println("error creando usuario por defecto:", err)
	}
}
