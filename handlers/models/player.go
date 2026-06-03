package models

type Player struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Value    float64 `json:"value"`
	TeamID   int     `json:"team_id"`
}
