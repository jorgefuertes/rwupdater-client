package client

import (
	"encoding/json"
	"net/http"
)

// GetCatalog - Download catalog for architecture
func GetArchList() ([]string, error) {
	list := make([]string, 0)

	res, err := http.Get(API + "/arch")
	if err != nil {
		return list, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return list, err
	}

	return list, nil
}

// IsArch - Check if arch exists in remote
func IsArch(arch string) bool {
	list, _ := GetArchList()
	for _, a := range list {
		if a == arch {
			return true
		}
	}

	return false
}
