package models

type TransferDTO struct {
	PlayerID int     `json:"player_id"`
	ToTeamID int     `json:"to_team_id"`
	Price    float64 `json:"price"`
}
