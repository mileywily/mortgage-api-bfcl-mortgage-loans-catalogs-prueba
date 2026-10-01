package create_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	createApp "api-hello-world/internal/app/example/create"
	appMocks "api-hello-world/internal/app/example/create/mocks"
	"api-hello-world/internal/core/domain"
	createHandler "api-hello-world/internal/infra/http/gin/example/create"
)

func TestNewHandler_Handle(t *testing.T) {
	validBody := `{"body":{"field2":"value"}}`

	tests := []struct {
		name                string
		xFieldHeader        string
		body                string
		controllerResponse  createApp.ResponseDto
		controllerError     error
		expectedStatus      int
		expectedMessage     string
		expectError         bool
		expectControllerHit bool
	}{
		{
			name:                "success",
			xFieldHeader:        "headerValue",
			body:                validBody,
			controllerResponse:  createApp.ResponseDto{Code: "222", Message: "result"},
			expectedStatus:      http.StatusOK,
			expectedMessage:     "result",
			expectControllerHit: true,
		},
		{
			name:        "missing mandatory header produces bind error",
			body:        validBody,
			expectError: true,
		},
		{
			name:         "malformed body produces bind error",
			xFieldHeader: "headerValue",
			body:         `{"body":`,
			expectError:  true,
		},
		{
			name:                "controller error is delegated",
			xFieldHeader:        "headerValue",
			body:                validBody,
			controllerError:     &domain.ServiceError{Message: "request failed"},
			expectError:         true,
			expectControllerHit: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			controllerMock := new(appMocks.Controller)
			controllerMock.On("Create", mock.Anything, mock.Anything).Return(tt.controllerResponse, tt.controllerError).Maybe()

			var capturedErrors []*gin.Error
			engine := gin.New()
			engine.Use(func(c *gin.Context) {
				c.Next()
				capturedErrors = c.Errors
			})
			engine.POST("/anything/:something/example", createHandler.NewHandler(controllerMock, createHandler.NewDecoder()))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/anything/foo/example", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.xFieldHeader != "" {
				req.Header.Set("X-Field", tt.xFieldHeader)
			}
			engine.ServeHTTP(w, req)

			if tt.expectError {
				assert.NotEmpty(t, capturedErrors)
			} else {
				require.Empty(t, capturedErrors)
				assert.Equal(t, tt.expectedStatus, w.Code)
				var res createApp.ResponseDto
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
				assert.Equal(t, tt.expectedMessage, res.Message)
			}
			if tt.expectControllerHit {
				controllerMock.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			} else {
				controllerMock.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestNewDecoder_UriBindError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "missing uri param produces bind error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/anything//example", bytes.NewBufferString(`{"body":{"field2":"value"}}`))
			c.Request.Header.Set("X-Field", "headerValue")

			_, err := createHandler.NewDecoder()(c)

			assert.Error(t, err)
		})
	}
}
