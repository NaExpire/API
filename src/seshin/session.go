package seshin

import (
	"database/sql"

	"github.com/nu7hatch/gouuid"
)

// GenerateSessionID will create and return a unique identifier which is used to
// validate user sessions.
func GenerateSessionID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// CreateSession will, given a database connection and a session ID, create a
// new
func CreateSession(db *sql.DB, sessionID string) error {
	// store given session id at the row of the session table with
	// assiociated user id. return errors if encountered
	_, err := db.Exec("INSERT INTO sessions ('session-content') VALUES (?)", sessionID)
	return err
}

func ValidateSession(db *sql.DB, sessionID string) (bool, error) {
	// look up session in session id column, throw error if not found.
	rows, err := db.Query("SELECT 'session-content' FROM sessions WHERE 'session-content'=?", sessionID)
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, err
}

func InvalidateSession(db *sql.DB, sessionID string) error {
	// remove session from session id table.
	_, err := db.Query("DELETE FROM sessions WHERE 'session-content'=?", sessionID)
	return err
}
