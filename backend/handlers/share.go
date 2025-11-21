package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fileupload/database"
	"fileupload/models"
	"fileupload/utils"
)

type CreateShareLinkRequest struct {
	FileID   string  `json:"file_id"`
	Password *string `json:"password"`
	ExpiresIn *int   `json:"expires_in"` // hours
}

type CreateShareLinkResponse struct {
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
}

type AccessShareLinkRequest struct {
	Password *string `json:"password"`
}

func CreateShareLink(w http.ResponseWriter, r *http.Request) {
	var req CreateShareLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify file exists
	row, _ := database.GetFileByID(req.FileID)
	var file models.File
	var tagsJSON string
	var parentID *string
	if err := row.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
		&file.Description, &file.Version, &file.VersionNumber, &parentID, &file.Path,
		&tagsJSON, &file.CreatedAt, &file.UpdatedAt); err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Generate short code
	shortCode, err := utils.GenerateShortCode(8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hash password if provided
	var hashedPassword *string
	if req.Password != nil && *req.Password != "" {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hashedPassword = &hashed
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		expiresAt = &expiry
	}

	linkID, _ := utils.GenerateID()
	if err := database.SaveShareLink(linkID, req.FileID, shortCode, hashedPassword, expiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateShareLinkResponse{
		ShortCode: shortCode,
		URL:       "/s/" + shortCode,
	})
}

func GetSharedFile(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("code")
	if shortCode == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	// Get share link
	row, _ := database.GetShareLinkByCode(shortCode)
	var link models.ShareLink
	var password *string
	var expiresAt *time.Time
	if err := row.Scan(&link.ID, &link.FileID, &link.ShortCode, &password, &expiresAt, &link.CreatedAt); err != nil {
		http.Error(w, "Share link not found", http.StatusNotFound)
		return
	}

	// Check expiration
	if expiresAt != nil && time.Now().After(*expiresAt) {
		http.Error(w, "Share link has expired", http.StatusGone)
		return
	}

	// Check password
	if password != nil && *password != "" {
		var req AccessShareLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Password == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"password_required": true,
			})
			return
		}

		if !utils.CheckPasswordHash(*req.Password, *password) {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	}

	// Get file
	fileRow, _ := database.GetFileByID(link.FileID)
	var file models.File
	var tagsJSON string
	var parentID *string
	if err := fileRow.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
		&file.Description, &file.Version, &file.VersionNumber, &parentID, &file.Path,
		&tagsJSON, &file.CreatedAt, &file.UpdatedAt); err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file.ParentID = parentID
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &file.Tags)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}
