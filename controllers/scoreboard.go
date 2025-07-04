package controllers

import (
	"CaptureTheSoul/database"
	"encoding/json"
	"net/http"
)

type ScoreEntry struct{
	Rank int `json:"rank"`
	SessionID string `json:"session_id"`
	Score int `json:"score"`
}

func Scoreboard(w http.ResponseWriter, r *http.Request){
	rows, err := database.DB.Query(`
		SELECT session_id, score
		FROM players
		ORDER BY score DESC, created_at ASC
	`)
	if err != nil{
		http.Error(w, "Failed to fetch scoreboard", http.StatusInternalServerError)
		return
	}
	var results []ScoreEntry
	rank := 1
	for rows.Next(){
		var session string
		var score int
		err := rows.Scan(&session, &score)
		if err !=  nil{
			continue
		}
		results = append(results, ScoreEntry{
			Rank: rank,
			SessionID: session,
			Score: score,
		})
		rank++
	}
	w.Header().Set("Content-Type", "applications/json")
	json.NewEncoder(w).Encode(results)
}