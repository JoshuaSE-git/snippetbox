package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JoshuaSE-git/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal()
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(
		t,
		res.Header.Get("Content-Security-Policy"),
		"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
	)

	assert.Equal(t, res.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, res.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, res.Header.Get("X-XXS-Protection"), "0")
	assert.Equal(t, res.Header.Get("Server"), "Go")

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
