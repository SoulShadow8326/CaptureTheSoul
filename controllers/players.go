package controllers

import (
	"net/http"
	"time"
	"database/sql"
	"github.com/google/uuid"
	"CaptureTheSoul/database"
	"fmt"
)

func Index(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	var sessionID string

	if err != nil || cookie.Value == "" {
		sessionID = uuid.New().String()
		flag := fmt.Sprintf("CTS{%s}", uuid.New().String())

		port, containerID, err := StartChallengeContainer(flag)
		if err != nil {
			http.Error(w, "failed to start challenge container", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "session_id",
			Value: sessionID,
			Path: "/",
			HttpOnly: true,
			Expires: time.Now().Add(24 * time.Hour),
		})

		_, err = database.DB.Exec(`INSERT INTO players (session_id, score) VALUES (?, 0)`, sessionID)
		if err != nil {
			http.Error(w, "failed to create player", http.StatusInternalServerError)
			return
		}

		_, err = database.DB.Exec(`INSERT INTO flags (session_id, flag) VALUES (?, ?)`, sessionID, flag)
		if err != nil {
			http.Error(w, "failed to store flag", http.StatusInternalServerError)
			return
		}

		_, err = database.DB.Exec(`INSERT INTO services (session_id, port, container_id) VALUES (?, ?, ?)`, sessionID, port, containerID)
		if err != nil {
			http.Error(w, "failed to store service info", http.StatusInternalServerError)
			return
		}
	} else {
		sessionID = cookie.Value
		var exists string
		err := database.DB.QueryRow(`SELECT session_id FROM players WHERE session_id = ?`, sessionID).Scan(&exists)

		if err == sql.ErrNoRows {
			_, err := database.DB.Exec(`INSERT INTO players (session_id, score) VALUES (?, 0)`, sessionID)
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
