package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api/user"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers(ctx context.Context) []domain.User {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User)
}

func (m *MockUserService) GetUser(ctx context.Context, id int) *domain.User {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User)
}

func (m *MockUserService) AddUser(ctx context.Context, user domain.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Remove(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) SearchUser(ctx context.Context, firstName, lastName string) ([]domain.User, error) {
	args := m.Called(ctx, firstName, lastName)
	return args.Get(0).([]domain.User), args.Error(1)
}

func setupTest(t *testing.T) (*user.User, *MockUserService) {
	mockService := new(MockUserService)
	api := user.NewUser(mockService, nil)
	return api, mockService
}

func TestGetUsers(t *testing.T) {
	api, mockService := setupTest(t)

	expectedUsers := []domain.User{
		{
			ID:        1,
			Login:     "testuser1",
			FirstName: "Test",
			LastName:  "User1",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
		{
			ID:        2,
			Login:     "testuser2",
			FirstName: "Test",
			LastName:  "User2",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
	}

	mockService.On("GetUsers", mock.Anything).Return(expectedUsers)

	req := httptest.NewRequest("GET", "/user/get_list", nil)
	w := httptest.NewRecorder()

	api.GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, response)
}

func TestGetUser(t *testing.T) {
	api, mockService := setupTest(t)

	expectedUser := &domain.User{
		ID:        1,
		Login:     "testuser",
		FirstName: "Test",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
	}

	mockService.On("GetUser", mock.Anything, 1).Return(expectedUser)

	req := httptest.NewRequest("GET", "/user/get?id=1", nil)
	w := httptest.NewRecorder()

	api.GetUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, &response)
}

func TestAddUser(t *testing.T) {
	api, mockService := setupTest(t)

	newUser := domain.User{
		Login:     "newuser",
		FirstName: "New",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
	}

	mockService.On("AddUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(1, nil)

	userJSON, _ := json.Marshal(newUser)
	req := httptest.NewRequest("POST", "/user/add", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.AddUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		ID int `json:"id"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
}

func TestUpdateUser(t *testing.T) {
	api, mockService := setupTest(t)

	userToUpdate := domain.User{
		ID:        1,
		Login:     "updateduser",
		FirstName: "Updated",
		LastName:  "User",
		Birthday:  time.Now(),
		GenderID:  1,
		City:      "Moscow",
		Enabled:   true,
	}

	mockService.On("UpdateUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)

	userJSON, _ := json.Marshal(userToUpdate)
	req := httptest.NewRequest("PUT", "/user/update", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	api.UpdateUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveUser(t *testing.T) {
	api, mockService := setupTest(t)

	mockService.On("Remove", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/user/remove?id=1", nil)
	w := httptest.NewRecorder()

	api.Remove(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUser(t *testing.T) {
	api, mockService := setupTest(t)

	expectedUsers := []domain.User{
		{
			ID:        1,
			Login:     "testuser1",
			FirstName: "Test",
			LastName:  "User1",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
		{
			ID:        2,
			Login:     "testuser2",
			FirstName: "Test",
			LastName:  "User2",
			Birthday:  time.Now(),
			GenderID:  1,
			City:      "Moscow",
			Enabled:   true,
		},
	}

	mockService.On("SearchUser", mock.Anything, "Test", "User").Return(expectedUsers, nil)

	req := httptest.NewRequest("GET", "/user/search?first_name=Test&last_name=User", nil)
	w := httptest.NewRecorder()

	api.SearchUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.User
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, response)
}
