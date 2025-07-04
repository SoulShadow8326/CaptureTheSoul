package controllers

import (
	"CaptureTheSoul/database"
	"encoding/json"
	"net/http"
)

type ChallengeEntry struct{
	Host string `json:"host"`
	Port int `json:"port"`
	Status string `json:"status"`
}

func ListChallenges(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == ""{
		http.Error(w, "your session doesnt exist how did you mess this up", http.StatusUnauthorized)
		return
	}
	rows, err := database.DB.Query(`
		SELECT host, port FROM services
		WHERE session_id != ?
	`, cookie.Value)
	if err != nil{
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []ChallengeEntry
	for rows.Next(){
		var host string
		var port int
		err := rows.Scan(&host, &port)
		if err != nil{
			continue
		}

		results = append(results, ChallengeEntry{
			Host: host,
			Port: port,
			Status: "online", //fix
		})
	}
	data, err := json.Marshal(results)
	if err != nil{
		http.Error(w, "Encoding err", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}