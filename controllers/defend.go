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
	Flag string `json:"flag"`
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
	var flag string

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
	err = database.DB.QueryRow(`
		SELECT flag FROM flags WHERE owner_session_id = ?
	`, sessionID).Scan(&flag)
	if err != nil {
		http.Error(w, "flag not found", http.StatusInternalServerError)
		return
	}
	resp := DefendInfo{
		Host: host,
		Port: port,
		Status: "online",
		Flag: flag,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}