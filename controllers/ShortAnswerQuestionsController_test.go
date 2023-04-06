package controllers

import (
	"AAT_Api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortAnswerQuestionController(t *testing.T) {
	t.Run("CreateShortAnswerQuestion", testCreateShortAnswerQuestion)
	t.Run("GetShortAnswerQuestions", testGetShortAnswerQuestions)
	t.Run("GetShortAnswerQuestion", testGetShortAnswerQuestion)
	t.Run("UpdateShortAnswerQuestion", testUpdateShortAnswerQuestion)
	t.Run("DeleteShortAnswerQuestion", testDeleteShortAnswerQuestion)
}

func testCreateShortAnswerQuestion(t *testing.T) {
	type args struct {
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Create a new short answer question",
			args: args{
				req: httptest.NewRequest("POST", "/short-answer", strings.NewReader(
					`{
						  "id": null,
						  "title": "What is the capital of France?",
						  "description": "This is a sample question",
						  "credit": 10,
						  "correct_answer": "Paris",
						  "tags": [
							"France",
							"Capital"
						  ]
						}`)),
				expectedStatus: http.StatusCreated,
				expectedBody: `{
				"id": 1,
				"title": "What is the capital of France?",
				"description": "This is a sample question",
				"credit": 10,
				"correct_answer": "Paris",
				"tags": [
					"France",
					"Capital"
					]
				}`,
			},
		},
	}
	for _, tt := range tests {
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Create a new router instance
		router := gin.Default()
		router.POST("/short-answer", CreateShortAnswerQuestion)
		// Connect database
		models.ConnectDatabase()
		// Dispatch request
		router.ServeHTTP(rr, tt.args.req)
		// Check status code
		if status := rr.Code; status != tt.args.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.args.expectedStatus)
		}
		// Check response body
		if compareJSON(rr.Body.String(), tt.args.expectedBody) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.args.expectedBody)
		}
	}
}

func testDeleteShortAnswerQuestion(t *testing.T) {
	type args struct {
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Delete a short answer question",
			args: args{
				req:            httptest.NewRequest("DELETE", "/short-answer/1", nil),
				expectedStatus: http.StatusNoContent,
				expectedBody: `{
					"message": "Question deleted successfully"
				}`,
			},
		},
	}
	for _, tt := range tests {
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Create a new router instance
		router := gin.Default()
		router.DELETE("/short-answer/:id", DeleteShortAnswerQuestion)
		// Connect database
		models.ConnectDatabase()
		// Dispatch request
		router.ServeHTTP(rr, tt.args.req)
		// Check status code
		if status := rr.Code; status != tt.args.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.args.expectedStatus)
		}
	}
}

func testGetShortAnswerQuestion(t *testing.T) {
	type args struct {
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get a short answer question",
			args: args{
				req:            httptest.NewRequest("GET", "/short-answer/1", nil),
				expectedStatus: http.StatusOK,
				expectedBody: `{
					"id": 1,
					"title": "What is the capital of France?",
					"description": "This is a sample question",
					"credit": 10,
					"correct_answer": "Paris",
					"tags": [
						"France",
						"Capital"
						]
					}`,
			},
		},
	}
	for _, tt := range tests {
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Create a new router instance
		router := gin.Default()
		router.GET("/short-answer/:id", GetShortAnswerQuestion)
		// Connect database
		models.ConnectDatabase()
		// Dispatch request
		router.ServeHTTP(rr, tt.args.req)
		// Check status code
		if status := rr.Code; status != tt.args.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.args.expectedStatus)
		}
		// Check response body
		if compareJSON(rr.Body.String(), tt.args.expectedBody) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.args.expectedBody)
		}
	}
}

func testGetShortAnswerQuestions(t *testing.T) {
	type args struct {
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get all short answer questions",
			args: args{
				req:            httptest.NewRequest("GET", "/short-answer", nil),
				expectedStatus: http.StatusOK,
				expectedBody: `[
					{
						"id": 1,
						"title": "What is the capital of France?",
						"description": "This is a sample question",
						"credit": 10,
						"correct_answer": "Paris",
						"tags": [
							"France",
							"Capital"
							]
					}
				]`,
			},
		},
	}
	for _, tt := range tests {
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Create a new router instance
		router := gin.Default()
		router.GET("/short-answer", GetShortAnswerQuestions)
		// Connect database
		models.ConnectDatabase()
		// Dispatch request
		router.ServeHTTP(rr, tt.args.req)
		// Check status code
		if status := rr.Code; status != tt.args.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.args.expectedStatus)
		}
		// Check response body
		if compareJSON(rr.Body.String(), tt.args.expectedBody) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.args.expectedBody)
		}
	}
}

func testUpdateShortAnswerQuestion(t *testing.T) {
	type args struct {
		req            *http.Request
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Update a short answer question",
			args: args{
				req: httptest.NewRequest("PUT", "/short-answer/1", strings.NewReader(`{
					"id": 1,
					"title": "What is the capital of France?",
					"description": "This is a sample question",
					"credit": 10,
					"correct_answer": "Paris",
					"tags": [
						"France",
						"Capital"
						]
					}`)),
				expectedStatus: http.StatusOK,
				expectedBody: `{
					"id": 1,
					"title": "What is the capital of France?",
					"description": "This is a sample question",
					"credit": 10,
					"correct_answer": "Paris",
					"tags": [
						"France",
						"Capital"
						]
					}`,
			},
		},
	}
	for _, tt := range tests {
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Create a new router instance
		router := gin.Default()
		router.PUT("/short-answer/:id", UpdateShortAnswerQuestion)
		// Connect database
		models.ConnectDatabase()
		// Dispatch request
		router.ServeHTTP(rr, tt.args.req)
		// Check status code
		if status := rr.Code; status != tt.args.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.args.expectedStatus)
		}
		// Check response body
		if compareJSON(rr.Body.String(), tt.args.expectedBody) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.args.expectedBody)
		}
	}
}
