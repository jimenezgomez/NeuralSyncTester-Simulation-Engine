package dbmanager

import "database/sql"

type QueryManager struct {
	db *sql.DB
}

func NewQueryManager(db *sql.DB) *QueryManager {
	return &QueryManager{db: db}
}

func (m *QueryManager) DB() *sql.DB {
	return m.db
}
