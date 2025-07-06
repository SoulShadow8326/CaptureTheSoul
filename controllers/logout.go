package controllers

import (
	"CaptureTheSoul/database"
	"net/http"
	"os/exec"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}
	sessionID := cookie.Value

	var containerID string
	err = database.DB.QueryRow(`SELECT container_id FROM services WHERE session_id = ?`, sessionID).Scan(&containerID)
	if err == nil && containerID != "" {
		exec.Command("docker", "rm", "-f", containerID).Run()
	}

	database.DB.Exec(`DELETE FROM services WHERE session_id = ?`, sessionID)
	database.DB.Exec(`DELETE FROM flags WHERE session_id = ?`, sessionID)
	database.DB.Exec(`DELETE FROM players WHERE session_id = ?`, sessionID)
	database.DB.Exec(`DELETE FROM submissions WHERE session_id = ?`, sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Write([]byte("Session ended"))
}
