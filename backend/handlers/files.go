package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"fileupload/database"
	"fileupload/models"
)

func GetFiles(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	tag := r.URL.Query().Get("tag")

	if limit == 0 {
		limit = 50
	}

	rows, err := database.GetAllFiles(limit, offset, tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var tagsJSON string
		var parentID *string

		if err := rows.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
			&file.Description, &file.Version, &file.VersionNumber, &parentID, &file.Path,
			&tagsJSON, &file.CreatedAt, &file.UpdatedAt); err != nil {
			continue
		}

		file.ParentID = parentID
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &file.Tags)
		}

		files = append(files, file)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	row, _ := database.GetFileByID(fileID)
	var file models.File
	var tagsJSON string
	var parentID *string

	if err := row.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
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

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	row, _ := database.GetFileByID(fileID)
	var file models.File
	var tagsJSON string
	var parentID *string

	if err := row.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
		&file.Description, &file.Version, &file.VersionNumber, &parentID, &file.Path,
		&tagsJSON, &file.CreatedAt, &file.UpdatedAt); err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileData, err := os.ReadFile(file.Path)
	if err != nil {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+file.OriginalName)
	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	w.Write(fileData)
}

func GetFileVersions(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	rows, err := database.GetFilesByName(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var tagsJSON string
		var parentID *string

		if err := rows.Scan(&file.ID, &file.Name, &file.OriginalName, &file.Size, &file.MimeType,
			&file.Description, &file.Version, &file.VersionNumber, &parentID, &file.Path,
			&tagsJSON, &file.CreatedAt, &file.UpdatedAt); err != nil {
			continue
		}

		file.ParentID = parentID
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &file.Tags)
		}

		files = append(files, file)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
