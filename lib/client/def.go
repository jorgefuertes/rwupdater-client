package client

const API = "https://core.abadiaretro.com"

// File
type File struct {
	ID        string `json:"id,omitempty"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Core      string `json:"core,omitempty"`
	Version   string `json:"version,omitempty"`
	Timestamp int64  `json:"ts"`
}
