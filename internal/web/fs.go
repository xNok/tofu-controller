package web

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// FileSystem wraps an embedded filesystem to handle SPA routing.
// It serves index.html for all non-existent file requests (404s),
// enabling client-side routing within the single-page application.
type FileSystem struct {
	fs fs.FS
}

// NewFileSystem creates a new filesystem handler for SPA routing.
func NewFileSystem(efs fs.FS) *FileSystem {
	return &FileSystem{fs: efs}
}

// Open implements fs.FS interface for SPA routing.
func (w *FileSystem) Open(name string) (fs.File, error) {
	// Normalize path
	name = strings.TrimPrefix(name, "/")

	// Try to open the requested file
	f, err := w.fs.Open(filepath.Join("dist", name))
	if err == nil {
		return f, nil
	}

	// If the file doesn't exist, serve index.html for SPA routing
	if !strings.HasSuffix(name, "/") {
		indexPath := filepath.Join("dist", "index.html")
		f, err := w.fs.Open(indexPath)
		if err == nil {
			return f, nil
		}
	}

	return nil, err
}
