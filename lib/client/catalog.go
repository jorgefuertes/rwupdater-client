package client

import (
	"encoding/json"
	"net/http"

	"git.martianoids.com/queru/retroupdater-client/lib/file"
)

// GetCatalog - Download catalog for architecture
func GetCatalog(arch string) (*file.Catalog, error) {
	cat := new(file.Catalog)

	res, err := http.Get(API + "/files/catalog/" + arch)
	if err != nil {
		return cat, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&cat); err != nil {
		return cat, err
	}

	return cat, nil
}
