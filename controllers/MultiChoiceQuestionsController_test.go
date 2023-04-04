package controllers

import (
	"AAT_Api/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMultiChoiceQuestionController(t *testing.T) {
	t.Run("CreateMultiChoiceQuestion", TestCreateMultiChoiceQuestion)
	t.Run("GetMultiChoiceQuestions", TestGetMultiChoiceQuestions)
	t.Run("GetMultiChoiceQuestion", TestGetMultiChoiceQuestion)
}

func TestCreateMultiChoiceQuestion(t *testing.T) {
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
			name: "Create a new multi choice question",
			args: args{
				req: httptest.NewRequest("POST", "/multi-choice", strings.NewReader(
					`{
						  "id": null,
						  "title": "What is the capital of France?",
						  "description": "This is a sample question",
						  "credit": 10,
						  "answers": [
							{
							  "content": "Paris",
							  "is_correct": true
							},
							{
							  "content": "London",
							  "is_correct": false
							},
							{
							  "content": "New York",
							  "is_correct": false
							}
						  ],
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
				"answers": [
					{
						"content": "Paris",
						"is_correct": true
					},
					{
						"content": "London",
						"is_correct": false
					},
					{
						"content": "New York",
						"is_correct": false
					}
					],
				"tags": [
					"France",
					"Capital"
					]
				}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run test here
			// Create a new HTTP response recorder
			rr := httptest.NewRecorder()

			// Create a new router instance
			r := gin.Default()

			// Connect database
			models.ConnectDatabase()

			// Define the route
			r.POST("/multi-choice", CreateMultiChoiceQuestion)

			// Dispatch the HTTP request
			r.ServeHTTP(rr, tt.args.req)

			// Check the status code
			if status := rr.Code; status != tt.args.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.args.expectedStatus)
			}

			// Check the response body
			if compareJSON(rr.Body.String(), tt.args.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.args.expectedBody)
			}
		})
	}
}

func TestGetMultiChoiceQuestions(t *testing.T) {
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
			name: "Get all multi choice questions",
			args: args{
				req:            httptest.NewRequest("GET", "/multi-choice", nil),
				expectedStatus: http.StatusOK,
				expectedBody: `[
					{
						"id": 1,
						"title": "What is the capital of France?",
						"description": "This is a sample question",
						"credit": 10,
						"answers": [
							{
								"content": "Paris",
								"is_correct": true
							},
							{
								"content": "London",
								"is_correct": false
							},
							{
								"content": "New York",
								"is_correct": false
							}
						],
						"tags": [
							"France",
							"Capital"
						]
					}
				]`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run test here
			// Create a new HTTP response recorder
			rr := httptest.NewRecorder()

			// Create a new router instance
			r := gin.Default()

			// Connect database
			models.ConnectDatabase()

			// Define the route
			r.GET("/multi-choice", GetMultiChoiceQuestions)

			// Dispatch the HTTP request
			r.ServeHTTP(rr, tt.args.req)

			// Check the status code
			if status := rr.Code; status != tt.args.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.args.expectedStatus)
			}

			// Check the response body
			if compareJSON(rr.Body.String(), tt.args.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.args.expectedBody)
			}
		})
	}
}

func TestGetMultiChoiceQuestion(t *testing.T) {
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
			name: "Get a multi choice question",
			args: args{
				req:            httptest.NewRequest("GET", "/multi-choice/1", nil),
				expectedStatus: http.StatusOK,
				expectedBody: `{
					"id": 1,
					"title": "What is the capital of France?",
					"description": "This is a sample question",
					"credit": 10,
					"answers": [
						{
							"content": "Paris",
							"is_correct": true
						},
						{
							"content": "London",
							"is_correct": false
						},
						{
							"content": "New York",
							"is_correct": false
						}
					],
					"tags": [
						"France",
						"Capital"
					]
				}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run test here
			// Create a new HTTP response recorder
			rr := httptest.NewRecorder()

			// Create a new router instance
			r := gin.Default()

			// Define the route
			r.GET("/multi-choice/:id", GetMultiChoiceQuestion)

			// Connect database
			models.ConnectDatabase()

			// Dispatch the HTTP request
			r.ServeHTTP(rr, tt.args.req)

			// Check the status code
			if status := rr.Code; status != tt.args.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.args.expectedStatus)
			}

			// Check the response body
			if compareJSON(rr.Body.String(), tt.args.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.args.expectedBody)
			}
		})
	}
}

func compareJSON(a, b string) bool {
	var o1 interface{}
	var o2 interface{}
	err := json.Unmarshal([]byte(a), &o1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(b), &o2)
	if err != nil {
		panic(err)
	}
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return string(aBytes) == string(bBytes)
}
