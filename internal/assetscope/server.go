package assetscope

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var dashboardFiles embed.FS

// NewHandler returns an HTTP handler for the embedded AssetScope dashboard.
func NewHandler() (http.Handler, error) {
	webFiles, err := fs.Sub(dashboardFiles, "web")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(webFiles)), nil
}
