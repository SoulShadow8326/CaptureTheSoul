package controllers

import(
	"net/http"
	"time"
	"database/sql"
	"github.com/google/uuid"
	"CaptureTheSoul/database"
)

func Index(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("session_id")
	var sessionID string
	if err != nil || cookie.Value == ""{
		sessionID = uuid.New().String()
		http.SetCookie(w, &http.Cookie{
			Name: "session_id",
			Value: sessionID,
			Path: "/",
			HttpOnly: true,
			Expires: time.Now().Add(24 * time.Hour),
		})
		_, err := database.DB.Exec(`
			INSERT INTO players (session_id, score)
			VALUES (?, 0)
		`, sessionID)
		if err != nil {
			http.Error(w, "failed to create player session", http.StatusInternalServerError)
			return
		}
		shortID := sessionID[:5]
		host := "player-"+shortID+".ctf.local"
		port := 9000+len(sessionID)%1000
		_, err = database.DB.Exec(`
			INSERT INTO services (session_id, host, port)
			VALUES (?,?,?)
		`, sessionID, host, port)
		if err != nil{
			http.Error(w, "failed to assign a service", http.StatusInternalServerError)
			return
		}
		flag:="CTS{"+uuid.New().String()[:8]+"}"
		value := 100
		_, err = database.DB.Exec(`
			INSERT INTO flags (flag, owner_session_id, value)
			VALUES (?, ?, ?)
		`, flag, sessionID, value)

		if err != nil{
			http.Error(w, "failed to assign a flag", http.StatusInternalServerError)
			return
		}
	} else {
		sessionID = cookie.Value
		var exists string
		err := database.DB.QueryRow(`
			SELECT session_id FROM players WHERE session_id = ?
		`, sessionID).Scan(&exists)

		if err == sql.ErrNoRows {
			_, err := database.DB.Exec(`
				INSERT INTO players (session_id, score)
				VALUES (?, 0)
			`, sessionID)
			if err != nil {
				http.Error(w, "failed to recreate session", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Welcome to CaptureTheSoul"))
}