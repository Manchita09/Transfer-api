package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"transfer-api/models"

	"github.com/gorilla/mux"
)

func GetPlayers(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
		SELECT p.id, p.name, p.position, p.value, t.name
		FROM players p
		JOIN teams t ON p.team_id = t.id
	`)
	if err != nil {
		log.Println("DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var players []models.PlayerResponse

	for rows.Next() {
		var player models.PlayerResponse

		err := rows.Scan(
			&player.ID,
			&player.Name,
			&player.Position,
			&player.Value,
			&player.Team, //ahora es string
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		log.Println("ROWS ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func GetPlayerByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var player models.PlayerResponse

	err = db.QueryRow(`
		SELECT p.id, p.name, p.position, p.value, t.name
		FROM players p
		JOIN teams t ON p.team_id = t.id
		WHERE p.id = $1
	`, id).Scan(
		&player.ID,
		&player.Name,
		&player.Position,
		&player.Value,
		&player.Team,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("QUERY ERROR:", err)
		http.Error(w, "Error en servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func CreatePlayer(w http.ResponseWriter, r *http.Request) {

	var dto models.CreatePlayerDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Validación básica
	if dto.Name == "" || dto.Position == "" || dto.Value <= 0 || dto.TeamID <= 0 {
		http.Error(w, "Campos inválidos", http.StatusBadRequest)
		return
	}

	// VALIDACIÓN DE DUPLICADO
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM players 
		WHERE name=$1 AND team_id=$2
	`, dto.Name, dto.TeamID).Scan(&count)

	if err != nil {
		log.Println("CHECK ERROR:", err)
		http.Error(w, "Error al validar jugador", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "El jugador ya existe en este equipo", http.StatusBadRequest)
		return
	}

	// INSERT
	_, err = db.Exec(`
		INSERT INTO players (name, position, value, team_id)
		VALUES ($1, $2, $3, $4)
	`, dto.Name, dto.Position, dto.Value, dto.TeamID)

	if err != nil {
		log.Println("CREATE ERROR:", err)
		http.Error(w, "Error al crear jugador", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Jugador creado correctamente",
	})
}

func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var dto models.CreatePlayerDTO

	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if dto.Name == "" || dto.Position == "" || dto.Value <= 0 || dto.TeamID <= 0 {
		http.Error(w, "Campos inválidos", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE players
		SET name = $1,
		    position = $2,
		    value = $3,
		    team_id = $4
		WHERE id = $5
	`, dto.Name, dto.Position, dto.Value, dto.TeamID, id)

	if err != nil {
		log.Println("UPDATE ERROR:", err)
		http.Error(w, "Error al actualizar jugador", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Jugador actualizado correctamente",
	})
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM players WHERE id=$1", id)

	if err != nil {
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Jugador eliminado",
	})

}
func GetTopPlayers(w http.ResponseWriter, r *http.Request) {

	rows, _ := db.Query(`
		SELECT name, value 
		FROM players 
		ORDER BY value DESC 
		LIMIT 5
	`)

	type Player struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}

	var list []Player

	for rows.Next() {
		var p Player
		rows.Scan(&p.Name, &p.Value)
		list = append(list, p)
	}

	json.NewEncoder(w).Encode(list)
}
