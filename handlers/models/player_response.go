package models

type PlayerResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Value    float64 `json:"value"`
	Team     string  `json:"team"`
}
