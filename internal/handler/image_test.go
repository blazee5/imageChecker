package handler

import (
	"bytes"
	"encoding/json"
	"github.com/blazee5/imageChecker/internal/service"
	"github.com/blazee5/imageChecker/internal/service/mocks"
	"github.com/blazee5/imageChecker/lib/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Check(t *testing.T) {
	type fields struct {
		service *mocks.Image
	}

	tests := []struct {
		name         string
		input        domain.CheckImageRequest
		mockFunc     func(f *fields)
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name: "success",
			input: domain.CheckImageRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.service.On("CheckImage", mock.Anything, mock.AnythingOfType("domain.CheckImageRequest")).
					Return(true, nil)
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"message": true,
			},
		},
		{
			name: "bad request",
			input: domain.CheckImageRequest{
				Image:     "",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc:   func(f *fields) {},
			wantStatus: http.StatusBadRequest,
			wantResponse: map[string]interface{}{
				"message": "Key: 'CheckImageRequest.Image' Error:Field validation for 'Image' failed on the 'required' tag",
			},
		},
		{
			name: "server error",
			input: domain.CheckImageRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.service.On("CheckImage", mock.Anything, mock.AnythingOfType("domain.CheckImageRequest")).
					Return(false, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"message": "server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				service: mocks.NewImage(t),
			}
			tt.mockFunc(&f)

			log := logger.NewLogger()
			svc := &service.Service{
				Image: f.service,
			}
			h := NewHandler(log, svc)

			r := gin.Default()
			RegisterHandlers(r, h)

			jsonData, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			req, _ := http.NewRequest(http.MethodGet, "/check-image", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantResponse, response)

			f.service.AssertExpectations(t)
		})
	}
}
