package assetscope

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashboardServesAssetScopeHome(t *testing.T) {
	handler, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "AssetScope") {
		t.Fatal("dashboard response does not contain AssetScope")
	}
	if !strings.Contains(response.Body.String(), "/api/assetscope/compare") {
		t.Fatal("dashboard does not post to the AssetScope API route")
	}
}

func TestDashboardReturnsNotFoundForUnknownAssets(t *testing.T) {
	handler, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/missing.js", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
