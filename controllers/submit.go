package controllers

import (
	"CaptureTheSoul/database"
	"database/sql"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var flagPattern = regexp.MustCompile(`^CTS\{[a-zA-Z0-9_\-]+\}$`)

func SubmitFlag(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	flag := strings.TrimSpace(string(body))

	if !flagPattern.MatchString(flag) {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	var owner string
	var value int
	err = database.DB.QueryRow(`
		SELECT session_id, value FROM flags WHERE flag = ?
	`, flag).Scan(&owner, &value)

	if err == sql.ErrNoRows {
		http.Error(w, "Flag not found", http.StatusForbidden)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if owner == sessionID {
		http.Error(w, "owners cannot submit their own flags", http.StatusForbidden)
		return
	}

	var exists int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM submissions
		WHERE session_id = ? AND flag = ?
	`, sessionID, flag).Scan(&exists)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "Flag already submitted", http.StatusConflict)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO submissions (session_id, flag)
		VALUES (?, ?)
	`, sessionID, flag)
	if err != nil {
		http.Error(w, "Failed to submit flag", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec(`
		UPDATE players SET score = score + ?
		WHERE session_id = ?
	`, value, sessionID)
	if err != nil {
		http.Error(w, "Failed to update score", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Flag accepted"))
}
