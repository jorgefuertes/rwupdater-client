package client

import (
	"encoding/json"
	"net/http"
)

// GetCatalog - Download catalog for architecture
func GetCatalog(arch string) (*[]File, error) {
	var files []File

	res, err := http.Get(API + "/files/catalog/" + arch)
	if err != nil {
		return &files, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&files); err != nil {
		return &files, err
	}

	return &files, nil
}
