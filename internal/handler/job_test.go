package handler

import (
	"bytes"
	"encoding/json"
	"github.com/blazee5/imageChecker/internal/domain"
	"github.com/blazee5/imageChecker/internal/service"
	"github.com/blazee5/imageChecker/internal/service/mocks"
	"github.com/blazee5/imageChecker/lib/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateJob(t *testing.T) {
	type fields struct {
		service *mocks.Job
	}

	tests := []struct {
		name         string
		input        domain.CreateJobRequest
		mockFunc     func(f *fields)
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name: "success",
			input: domain.CreateJobRequest{
				Name:      "test",
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.service.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"message": "success",
			},
		},
		{
			name: "server error",
			input: domain.CreateJobRequest{
				Name:      "test",
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc: func(f *fields) {
				f.service.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"message": "server error",
			},
		},
		{
			name: "missing job name",
			input: domain.CreateJobRequest{
				Image:     "nginx",
				IsPrivate: false,
				Username:  "",
				Password:  "",
			},
			mockFunc:     func(f *fields) {},
			wantStatus:   http.StatusBadRequest,
			wantResponse: map[string]interface{}{"message": "Key: 'CreateJobRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		{
			name: "private registry",
			input: domain.CreateJobRequest{
				Name:      "private-job",
				Image:     "my.private.registry/image",
				IsPrivate: true,
				Username:  "user",
				Password:  "password",
			},
			mockFunc: func(f *fields) {
				f.service.On("CreateJob", mock.Anything, mock.AnythingOfType("domain.CreateJobRequest")).
					Return(nil)
			},
			wantStatus:   http.StatusOK,
			wantResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				service: mocks.NewJob(t),
			}
			tt.mockFunc(&f)

			log := logger.NewLogger()
			svc := &service.Service{
				Job: f.service,
			}

			h := NewHandler(log, svc)
			r := gin.Default()
			r.POST("/check", h.CreateJob)

			jsonData, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/check", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response map[string]interface{}
			if tt.wantResponse != nil {
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantResponse != nil {
				assert.Equal(t, tt.wantResponse, response)
			}

			f.service.AssertExpectations(t)
		})
	}
}
