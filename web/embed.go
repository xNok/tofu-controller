package web

import "embed"

//go:embed dist/*
var webFS embed.FS

// GetFS returns the embedded web filesystem containing the built frontend assets.
func GetFS() embed.FS {
	return webFS
}
