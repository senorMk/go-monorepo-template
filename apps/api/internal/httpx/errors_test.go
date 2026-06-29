package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, Conflict("Email already registered"))

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}

	var body struct {
		StatusCode int    `json:"statusCode"`
		Error      string `json:"error"`
		Message    string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.StatusCode != 409 || body.Error != "CONFLICT" || body.Message != "Email already registered" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestWriteErrorMasksUnknownErrors(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, errors.New("some internal failure"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}

	var body struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Error != "INTERNAL_ERROR" {
		t.Fatalf("error = %q, want INTERNAL_ERROR", body.Error)
	}
	if body.Message == "some internal failure" {
		t.Fatal("internal error detail should not leak to the client")
	}
}
