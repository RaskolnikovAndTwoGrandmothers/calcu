package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "У тебя своя голова на плечах есть, но если очень надо, то тебе сюда"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedCode   int
	}{
		{"Simple Addition", "1 + 2", 3, http.StatusOK},
		{"Simple Subtraction", "5 - 2", 3, http.StatusOK},
		{"Simple Multiplication", "3 * 4", 12, http.StatusOK},
		{"Simple Division", "8 / 4", 2, http.StatusOK},
		{"Invalid Expression", "2 +", 0, http.StatusUnprocessableEntity},
		{"Division by Zero", "8 / 0", 0, http.StatusInternalServerError},
		{"Invalid Token", "3 & 5", 0, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(CalculationRequest{Expression: tt.expression})
			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(calculateHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var response map[string]float64
				json.NewDecoder(rr.Body).Decode(&response)
				if response["Результат"] != tt.expectedResult {
					t.Errorf("handler returned unexpected result: got %v want %v",
						response["Результат"], tt.expectedResult)
				}
			}
		})
	}
}