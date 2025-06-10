package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nouvadev/veritas/internal/api/handlers"
	"github.com/nouvadev/veritas/internal/config"
	database "github.com/nouvadev/veritas/internal/database/sqlc"
	"github.com/nouvadev/veritas/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of the Querier interface.
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateURL(ctx context.Context, originalUrl string) (int64, error) {
	args := m.Called(ctx, originalUrl)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) UpdateShortCode(ctx context.Context, params database.UpdateShortCodeParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockQuerier) GetURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func TestCreateShortURL(t *testing.T) {
	testCases := []struct {
		name               string
		body               io.Reader
		setupMock          func(*MockQuerier)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Success",
			body: strings.NewReader(`{"original_url": "https://example.com"}`),
			setupMock: func(mq *MockQuerier) {
				expectedID := int64(123)
				shortCode := utils.ToBase62(uint64(expectedID))

				mq.On("CreateURL", mock.Anything, "https://example.com").Return(expectedID, nil).Once()
				mq.On("UpdateShortCode", mock.Anything, database.UpdateShortCodeParams{
					ShortCode: shortCode,
					ID:        expectedID,
				}).Return(nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `{"short_url":"` + utils.ToBase62(123) + `"}`,
		},
		{
			name:               "Invalid JSON Body",
			body:               strings.NewReader(`{"original_url":}`),
			setupMock:          func(mq *MockQuerier) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "Invalid request body"}`,
		},
		{
			name:               "Invalid URL",
			body:               strings.NewReader(`{"original_url": "not-a-valid-url"}`),
			setupMock:          func(mq *MockQuerier) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "Invalid URL"}`,
		},
		{
			name: "CreateURL Fails",
			body: strings.NewReader(`{"original_url": "https://example.com"}`),
			setupMock: func(mq *MockQuerier) {
				mq.On("CreateURL", mock.Anything, "https://example.com").Return(int64(0), errors.New("db error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error": "Failed to create URL"}`,
		},
		{
			name: "UpdateShortCode Fails",
			body: strings.NewReader(`{"original_url": "https://example.com"}`),
			setupMock: func(mq *MockQuerier) {
				expectedID := int64(123)
				shortCode := utils.ToBase62(uint64(expectedID))

				mq.On("CreateURL", mock.Anything, "https://example.com").Return(expectedID, nil).Once()
				mq.On("UpdateShortCode", mock.Anything, database.UpdateShortCodeParams{
					ShortCode: shortCode,
					ID:        expectedID,
				}).Return(errors.New("db error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error": "Failed to update short code"}`,
		},
		{
			name:               "Empty Request Body",
			body:               strings.NewReader(""),
			setupMock:          func(mq *MockQuerier) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "Invalid request body"}`,
		},
		{
			name:               "Empty URL String",
			body:               strings.NewReader(`{"original_url": ""}`),
			setupMock:          func(mq *MockQuerier) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "Invalid URL"}`,
		},
		{
			name:               "Missing URL Field",
			body:               strings.NewReader(`{}`),
			setupMock:          func(mq *MockQuerier) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "Invalid URL"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockQuerier := new(MockQuerier)
			tc.setupMock(mockQuerier)

			app := &config.AppConfig{
				Querier: mockQuerier,
				Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
			}
			handler := handlers.NewURLHandler(app)

			req := httptest.NewRequest("POST", "/api/v1/urls", tc.body)
			rr := httptest.NewRecorder()

			// Execute
			handler.CreateShortURL(rr, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, rr.Code)

			if tc.expectedStatusCode == http.StatusCreated {
				var resp handlers.URLResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, utils.ToBase62(123), resp.ShortURL)
			} else {
				assert.JSONEq(t, tc.expectedBody, rr.Body.String())
			}

			mockQuerier.AssertExpectations(t)
		})
	}
}
