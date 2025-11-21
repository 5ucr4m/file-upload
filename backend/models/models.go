package models

import "time"

type File struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	OriginalName string   `json:"original_name"`
	Size        int64     `json:"size"`
	MimeType    string    `json:"mime_type"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	VersionNumber int     `json:"version_number"`
	ParentID    *string   `json:"parent_id"`
	Path        string    `json:"path"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FileVersion struct {
	ID          string    `json:"id"`
	FileID      string    `json:"file_id"`
	Version     string    `json:"version"`
	VersionNumber int     `json:"version_number"`
	Size        int64     `json:"size"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShareLink struct {
	ID          string     `json:"id"`
	FileID      string     `json:"file_id"`
	ShortCode   string     `json:"short_code"`
	Password    *string    `json:"password,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UploadChunk struct {
	UploadID    string `json:"upload_id"`
	ChunkIndex  int    `json:"chunk_index"`
	TotalChunks int    `json:"total_chunks"`
	Data        []byte `json:"-"`
}

type UploadSession struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	TotalChunks int       `json:"total_chunks"`
	Chunks      []int     `json:"chunks"`
	CreatedAt   time.Time `json:"created_at"`
}
