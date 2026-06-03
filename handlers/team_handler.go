package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
		SELECT id, name, budget FROM teams
	`)
	if err != nil {
		http.Error(w, "Error en DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Team struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Budget float64 `json:"budget"`
	}

	var teams []Team

	for rows.Next() {
		var t Team
		rows.Scan(&t.ID, &t.Name, &t.Budget)
		teams = append(teams, t)
	}

	json.NewEncoder(w).Encode(teams)
}

func GetTeamByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	var team struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Budget float64 `json:"budget"`
	}

	err := db.QueryRow(`
		SELECT id, name, budget FROM teams WHERE id=$1
	`, id).Scan(&team.ID, &team.Name, &team.Budget)

	if err != nil {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(team)
}
func GetTeamRanking(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
		SELECT name, budget 
		FROM teams 
		ORDER BY budget DESC
	`)
	if err != nil {
		http.Error(w, "Error en DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Team struct {
		Name   string  `json:"name"`
		Budget float64 `json:"budget"`
	}

	var teams []Team

	for rows.Next() {
		var t Team
		rows.Scan(&t.Name, &t.Budget)
		teams = append(teams, t)
	}

	json.NewEncoder(w).Encode(teams)
}
