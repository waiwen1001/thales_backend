package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Paginate struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

func GetIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	aid := vars["id"]
	if aid == "" {
		return 0, http.ErrMissingFile
	}

	id, err := strconv.Atoi(aid)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetPaginateFromRequest(r *http.Request) (Paginate, error) {
	defaultPage := 1
	defaultPageSize := 20

	query := r.URL.Query()
	pageStr := query.Get("page")
	if pageStr == "" {
		pageStr = strconv.Itoa(defaultPage)
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = defaultPage
	}

	// Parse pageSize
	pageSizeStr := query.Get("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = strconv.Itoa(defaultPageSize)
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}

	// Return Paginate struct
	return Paginate{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func Pagination(page int, pageSize int, totalCount int) Paginate {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if totalCount < 0 {
		totalCount = 0
	}

	return Paginate{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}
}

func StoreFile(r *http.Request, formFieldName string) (string, error) {
	file, header, err := r.FormFile(formFieldName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file: %w", err)
	}
	defer file.Close()

	// Generate a unique filename and save the file
	filename := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), header.Filename)
	filepath := "./static/" + filename // Example: save files under "./static/uploads/"

	out, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer out.Close()

	// Copy the uploaded file's contents to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// Return the relative file path
	return filename, nil
}
