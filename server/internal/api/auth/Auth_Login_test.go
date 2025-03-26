package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Successfully decode valid JSON request body into AuthUser struct
func TestLoginSuccessfullyDecodesValidJSON(t *testing.T) {
	// Arrange
	mockService := &mockAuthService{}
	mockService.On("Login", mock.Anything, "testuser", "testpass").Return(nil)

	auth := NewAuth(mockService, nil, "test-secret")

	reqBody := `{"user": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	auth.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// Handle malformed JSON in request body
func TestLoginHandlesMalformedJSON(t *testing.T) {
	// Arrange
	mockService := &mockAuthService{}
	auth := NewAuth(mockService, nil, "test-secret")

	reqBody := `{"user": "testuser", "password": ` // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Act
	auth.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Login")
}
