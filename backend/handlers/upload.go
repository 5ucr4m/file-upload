package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fileupload/database"
	"fileupload/models"
	"fileupload/utils"
)

type InitUploadRequest struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	TotalChunks int    `json:"total_chunks"`
}

type InitUploadResponse struct {
	UploadID string `json:"upload_id"`
}

type CompleteUploadRequest struct {
	UploadID    string   `json:"upload_id"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

func InitUpload(w http.ResponseWriter, r *http.Request) {
	var req InitUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadID, err := utils.GenerateID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := database.SaveUploadSession(uploadID, req.FileName, req.FileSize, req.TotalChunks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create temp directory for chunks
	tempDir := filepath.Join("uploads", "temp", uploadID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InitUploadResponse{UploadID: uploadID})
}

func UploadChunk(w http.ResponseWriter, r *http.Request) {
	uploadID := r.URL.Query().Get("upload_id")
	chunkIndex, _ := strconv.Atoi(r.URL.Query().Get("chunk_index"))

	if uploadID == "" {
		http.Error(w, "upload_id is required", http.StatusBadRequest)
		return
	}

	// Get upload session
	row, _ := database.GetUploadSession(uploadID)
	var session models.UploadSession
	var chunksJSON string
	if err := row.Scan(&session.ID, &session.FileName, &session.FileSize, &session.TotalChunks, &chunksJSON, &session.CreatedAt); err != nil {
		http.Error(w, "Upload session not found", http.StatusNotFound)
		return
	}

	json.Unmarshal([]byte(chunksJSON), &session.Chunks)

	// Save chunk
	tempDir := filepath.Join("uploads", "temp", uploadID)
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndex))

	file, err := os.Create(chunkPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update chunks list
	session.Chunks = append(session.Chunks, chunkIndex)
	if err := database.UpdateUploadSessionChunks(uploadID, session.Chunks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"chunks_uploaded": len(session.Chunks),
		"total_chunks":    session.TotalChunks,
	})
}

func CompleteUpload(w http.ResponseWriter, r *http.Request) {
	var req CompleteUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get upload session
	row, _ := database.GetUploadSession(req.UploadID)
	var session models.UploadSession
	var chunksJSON string
	if err := row.Scan(&session.ID, &session.FileName, &session.FileSize, &session.TotalChunks, &chunksJSON, &session.CreatedAt); err != nil {
		http.Error(w, "Upload session not found", http.StatusNotFound)
		return
	}

	json.Unmarshal([]byte(chunksJSON), &session.Chunks)

	// Verify all chunks are uploaded
	if len(session.Chunks) != session.TotalChunks {
		http.Error(w, "Not all chunks uploaded", http.StatusBadRequest)
		return
	}

	// Merge chunks
	fileID, _ := utils.GenerateID()
	finalDir := filepath.Join("uploads", "files")
	os.MkdirAll(finalDir, 0755)

	ext := filepath.Ext(session.FileName)
	finalPath := filepath.Join(finalDir, fileID+ext)

	finalFile, err := os.Create(finalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer finalFile.Close()

	tempDir := filepath.Join("uploads", "temp", req.UploadID)
	for i := 0; i < session.TotalChunks; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.Copy(finalFile, chunkFile)
		chunkFile.Close()
	}

	// Clean up temp directory
	os.RemoveAll(tempDir)

	// Check if file with same name exists (versioning)
	var parentID *string
	versionNumber := 1
	version := req.Version

	rows, err := database.GetFilesByName(session.FileName)
	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var existingFile models.File
			var tagsJSON string
			var parentIDNull *string
			rows.Scan(&existingFile.ID, &existingFile.Name, &existingFile.OriginalName, &existingFile.Size,
				&existingFile.MimeType, &existingFile.Description, &existingFile.Version, &existingFile.VersionNumber,
				&parentIDNull, &existingFile.Path, &tagsJSON, &existingFile.CreatedAt, &existingFile.UpdatedAt)

			versionNumber = existingFile.VersionNumber + 1
			if version == "" {
				version = fmt.Sprintf("v%d", versionNumber)
			}
		}
	}

	if version == "" {
		version = "v1"
	}

	// Detect MIME type
	mimeType := "application/octet-stream"
	if strings.HasSuffix(session.FileName, ".apk") {
		mimeType = "application/vnd.android.package-archive"
	} else if strings.HasSuffix(session.FileName, ".pdf") {
		mimeType = "application/pdf"
	} else if strings.HasSuffix(session.FileName, ".jpg") || strings.HasSuffix(session.FileName, ".jpeg") {
		mimeType = "image/jpeg"
	} else if strings.HasSuffix(session.FileName, ".png") {
		mimeType = "image/png"
	}

	// Save to database
	if err := database.SaveFile(fileID, session.FileName, session.FileName, mimeType, req.Description, version, finalPath, session.FileSize, versionNumber, parentID, req.Tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete upload session
	database.DeleteUploadSession(req.UploadID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"file_id":  fileID,
		"version":  version,
		"version_number": versionNumber,
	})
}

func CancelUpload(w http.ResponseWriter, r *http.Request) {
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "upload_id is required", http.StatusBadRequest)
		return
	}

	// Clean up temp directory
	tempDir := filepath.Join("uploads", "temp", uploadID)
	os.RemoveAll(tempDir)

	// Delete upload session
	database.DeleteUploadSession(uploadID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetUploadStatus(w http.ResponseWriter, r *http.Request) {
	uploadID := r.URL.Query().Get("upload_id")
	if uploadID == "" {
		http.Error(w, "upload_id is required", http.StatusBadRequest)
		return
	}

	row, _ := database.GetUploadSession(uploadID)
	var session models.UploadSession
	var chunksJSON string
	if err := row.Scan(&session.ID, &session.FileName, &session.FileSize, &session.TotalChunks, &chunksJSON, &session.CreatedAt); err != nil {
		http.Error(w, "Upload session not found", http.StatusNotFound)
		return
	}

	json.Unmarshal([]byte(chunksJSON), &session.Chunks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
