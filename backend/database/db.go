package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	return createTables()
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		original_name TEXT NOT NULL,
		size INTEGER NOT NULL,
		mime_type TEXT NOT NULL,
		description TEXT,
		version TEXT,
		version_number INTEGER DEFAULT 1,
		parent_id TEXT,
		path TEXT NOT NULL,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_id) REFERENCES files(id)
	);

	CREATE TABLE IF NOT EXISTS file_versions (
		id TEXT PRIMARY KEY,
		file_id TEXT NOT NULL,
		version TEXT NOT NULL,
		version_number INTEGER NOT NULL,
		size INTEGER NOT NULL,
		path TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (file_id) REFERENCES files(id)
	);

	CREATE TABLE IF NOT EXISTS share_links (
		id TEXT PRIMARY KEY,
		file_id TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		password TEXT,
		expires_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (file_id) REFERENCES files(id)
	);

	CREATE TABLE IF NOT EXISTS upload_sessions (
		id TEXT PRIMARY KEY,
		file_name TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		total_chunks INTEGER NOT NULL,
		chunks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_files_parent ON files(parent_id);
	CREATE INDEX IF NOT EXISTS idx_files_created ON files(created_at);
	CREATE INDEX IF NOT EXISTS idx_share_links_code ON share_links(short_code);
	`

	_, err := DB.Exec(schema)
	return err
}

func SaveFile(id, name, originalName, mimeType, description, version, path string, size int64, versionNumber int, parentID *string, tags []string) error {
	tagsJSON, _ := json.Marshal(tags)

	query := `INSERT INTO files (id, name, original_name, size, mime_type, description, version, version_number, parent_id, path, tags)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query, id, name, originalName, size, mimeType, description, version, versionNumber, parentID, path, string(tagsJSON))
	return err
}

func GetFileByID(id string) (*sql.Row, error) {
	query := `SELECT id, name, original_name, size, mime_type, description, version, version_number, parent_id, path, tags, created_at, updated_at
			  FROM files WHERE id = ?`
	return DB.QueryRow(query, id), nil
}

func GetFilesByName(name string) (*sql.Rows, error) {
	query := `SELECT id, name, original_name, size, mime_type, description, version, version_number, parent_id, path, tags, created_at, updated_at
			  FROM files WHERE original_name = ? ORDER BY version_number DESC`
	return DB.Query(query, name)
}

func GetAllFiles(limit, offset int, tag string) (*sql.Rows, error) {
	var query string
	var args []interface{}

	if tag != "" {
		query = `SELECT id, name, original_name, size, mime_type, description, version, version_number, parent_id, path, tags, created_at, updated_at
				 FROM files WHERE tags LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
		args = []interface{}{"%" + tag + "%", limit, offset}
	} else {
		query = `SELECT id, name, original_name, size, mime_type, description, version, version_number, parent_id, path, tags, created_at, updated_at
				 FROM files ORDER BY created_at DESC LIMIT ? OFFSET ?`
		args = []interface{}{limit, offset}
	}

	return DB.Query(query, args...)
}

func SaveShareLink(id, fileID, shortCode string, password *string, expiresAt *time.Time) error {
	query := `INSERT INTO share_links (id, file_id, short_code, password, expires_at) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, id, fileID, shortCode, password, expiresAt)
	return err
}

func GetShareLinkByCode(shortCode string) (*sql.Row, error) {
	query := `SELECT id, file_id, short_code, password, expires_at, created_at FROM share_links WHERE short_code = ?`
	return DB.QueryRow(query, shortCode), nil
}

func SaveUploadSession(id, fileName string, fileSize int64, totalChunks int) error {
	chunks := []int{}
	chunksJSON, _ := json.Marshal(chunks)

	query := `INSERT INTO upload_sessions (id, file_name, file_size, total_chunks, chunks) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, id, fileName, fileSize, totalChunks, string(chunksJSON))
	return err
}

func GetUploadSession(id string) (*sql.Row, error) {
	query := `SELECT id, file_name, file_size, total_chunks, chunks, created_at FROM upload_sessions WHERE id = ?`
	return DB.QueryRow(query, id), nil
}

func UpdateUploadSessionChunks(id string, chunks []int) error {
	chunksJSON, _ := json.Marshal(chunks)
	query := `UPDATE upload_sessions SET chunks = ? WHERE id = ?`
	_, err := DB.Exec(query, string(chunksJSON), id)
	return err
}

func DeleteUploadSession(id string) error {
	query := `DELETE FROM upload_sessions WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}
