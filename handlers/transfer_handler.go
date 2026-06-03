package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Manchita09/transfer-api/handlers/models"
)

func TransferPlayer(w http.ResponseWriter, r *http.Request) {

	// SOLO POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var dto models.TransferDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// VALIDACIONES BÁSICAS
	if dto.PlayerID <= 0 || dto.ToTeamID <= 0 || dto.Price <= 0 {
		http.Error(w, "Datos incompletos o inválidos", http.StatusBadRequest)
		return
	}

	db := GetDB()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error iniciando transacción", http.StatusInternalServerError)
		return
	}

	// rollback seguro
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var fromTeamID int
	var playerValue float64

	// Obtener info del jugador
	err = tx.QueryRow(`
		SELECT team_id, value 
		FROM players 
		WHERE id = $1
	`, dto.PlayerID).Scan(&fromTeamID, &playerValue)

	if err == sql.ErrNoRows {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error consultando jugador", http.StatusInternalServerError)
		return
	}

	// evitar misma transferencia
	if fromTeamID == dto.ToTeamID {
		http.Error(w, "El jugador ya pertenece a este equipo", http.StatusBadRequest)
		return
	}

	// Verificar presupuesto
	var toTeamBudget float64

	err = tx.QueryRow(`
		SELECT budget 
		FROM teams 
		WHERE id = $1
	`, dto.ToTeamID).Scan(&toTeamBudget)

	if err == sql.ErrNoRows {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error consultando equipo", http.StatusInternalServerError)
		return
	}

	if toTeamBudget < dto.Price {
		http.Error(w, "Fondos insuficientes", http.StatusBadRequest)
		return
	}

	// Restar dinero comprador
	_, err = tx.Exec(`
		UPDATE teams 
		SET budget = budget - $1 
		WHERE id = $2
	`, dto.Price, dto.ToTeamID)

	if err != nil {
		http.Error(w, "Error actualizando comprador", http.StatusInternalServerError)
		return
	}

	// Sumar dinero vendedor
	_, err = tx.Exec(`
		UPDATE teams 
		SET budget = budget + $1 
		WHERE id = $2
	`, dto.Price, fromTeamID)

	if err != nil {
		http.Error(w, "Error actualizando vendedor", http.StatusInternalServerError)
		return
	}

	// Actualizar jugador
	_, err = tx.Exec(`
		UPDATE players 
		SET team_id = $1 
		WHERE id = $2
	`, dto.ToTeamID, dto.PlayerID)

	if err != nil {
		http.Error(w, "Error actualizando jugador", http.StatusInternalServerError)
		return
	}

	// Guardar transferencia
	_, err = tx.Exec(`
		INSERT INTO transfers (player_id, from_team, to_team, price)
		VALUES ($1, $2, $3, $4)
	`, dto.PlayerID, fromTeamID, dto.ToTeamID, dto.Price)

	if err != nil {
		http.Error(w, "Error guardando transferencia", http.StatusInternalServerError)
		return
	}

	// CONFIRMAR
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error confirmando transacción", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transferencia realizada correctamente",
	})
}
func GetTransfers(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
		SELECT 
			p.name,
			t1.name as from_team,
			t2.name as to_team,
			tr.price
		FROM transfers tr
		JOIN players p ON tr.player_id = p.id
		JOIN teams t1 ON tr.from_team = t1.id
		JOIN teams t2 ON tr.to_team = t2.id
		ORDER BY tr.id DESC
	`)
	if err != nil {
		http.Error(w, "Error en DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Transfer struct {
		Player   string  `json:"player"`
		FromTeam string  `json:"from_team"`
		ToTeam   string  `json:"to_team"`
		Price    float64 `json:"price"`
	}

	var transfers []Transfer

	for rows.Next() {
		var t Transfer
		rows.Scan(&t.Player, &t.FromTeam, &t.ToTeam, &t.Price)
		transfers = append(transfers, t)
	}

	json.NewEncoder(w).Encode(transfers)
}
func ResetTransfers(w http.ResponseWriter, r *http.Request) {

	_, err := db.Exec("DELETE FROM transfers")
	if err != nil {
		http.Error(w, "Error limpiando historial", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Historial reiniciado",
	})
}
