package emails

import (
	"database/sql"
	"net/http"
)

type DBHandler struct {
	DB *sql.DB
}

func (db *DBHandler) SendHandler(w http.ResponseWriter, r *http.Request) {}
