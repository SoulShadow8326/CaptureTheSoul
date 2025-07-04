package controllers

import(
	"CaptureTheSoul/database"
	"database/sql"
	"encoding/json"
	"net/http"
)
type DefendInfo struct{
	Host string `json:"host"`
	Port int `json:"port"`
	Status string `json:"status"`
}

func DefendStatus(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == ""{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	var host string
	var port int 

	err = database.DB.QueryRow(`
		SELECT host, port FROM services WHERE session_id = ?
	`, sessionID).Scan(&host, &port)
	if err == sql.ErrNoRows{
		http.Error(w, "service not found", http.StatusNotFound)
		return
	} else if err != nil{
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	resp := DefendInfo{
		Host: host,
		Port: port,
		Status: "online",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}