package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"transfer-api/models"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var dto models.LoginDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	log.Println("Username recibido:", dto.Username)
	log.Println("Password recibido:", dto.Password)

	var storedPassword string

	err = db.QueryRow(
		"SELECT password FROM users WHERE username = $1",
		dto.Username,
	).Scan(&storedPassword)

	// SI NO EXISTE USUARIO
	if err != nil {
		log.Println("Usuario no encontrado:", err)
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// COMPARACIÓN
	if storedPassword != dto.Password {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// LOGIN CORRECTO
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login exitoso",
	})
}
