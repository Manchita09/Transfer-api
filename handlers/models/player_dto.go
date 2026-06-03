package models

type CreatePlayerDTO struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Value    float64 `json:"value"`
	TeamID   int     `json:"team_id"`
}
