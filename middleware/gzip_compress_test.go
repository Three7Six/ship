package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgfone/ship"
)

func TestGzip(t *testing.T) {
	s := ship.New()
	assert := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := s.NewContext(req, rec)

	// Skip if no Accept-Encoding header
	handler := Gzip()(func(ctx ship.Context) error {
		ctx.Response().Write([]byte("test"))
		return nil
	})

	handler(ctx)
	assert.Equal("test", rec.Body.String())

	// Gzip
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(ship.HeaderAcceptEncoding, "gzip")
	rec = httptest.NewRecorder()
	ctx = s.NewContext(req, rec)
	handler(ctx)
	assert.Equal("gzip", rec.Header().Get(ship.HeaderContentEncoding))
	assert.Contains(rec.Header().Get(ship.HeaderContentType), ship.MIMETextPlain)
	reader, err := gzip.NewReader(rec.Body)
	if assert.NoError(err) {
		buf := new(bytes.Buffer)
		defer reader.Close()
		buf.ReadFrom(reader)
		assert.Equal("test", buf.String())
	}
}

func TestGzipNoContent(t *testing.T) {
	s := ship.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(ship.HeaderAcceptEncoding, "gzip")
	rec := httptest.NewRecorder()
	ctx := s.NewContext(req, rec)
	handler := Gzip()(func(ctx ship.Context) error {
		return ctx.NoContent(http.StatusNoContent)
	})

	if assert.NoError(t, handler(ctx)) {
		assert.Empty(t, rec.Header().Get(ship.HeaderContentEncoding))
		assert.Empty(t, rec.Header().Get(ship.HeaderContentType))
		assert.Equal(t, 0, len(rec.Body.Bytes()))
	}
}

func TestGzipErrorReturned(t *testing.T) {
	s := ship.New()
	s.Use(Gzip())
	s.R("/", func(ctx ship.Context) error { return ship.ErrNotFound })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(ship.HeaderAcceptEncoding, "gzip")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Empty(t, rec.Header().Get(ship.HeaderContentEncoding))
}
