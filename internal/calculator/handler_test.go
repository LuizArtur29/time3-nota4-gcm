package calculator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMultiplyHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		contentType    string
		expectedStatus int
		expectedResult float64
		expectedError  string
	}{
		{
			name:           "multiplicacao de inteiros",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedResult: 50,
		},
		{
			name:           "multiplicacao de decimais",
			method:         http.MethodPost,
			body:           `{"a":2.5,"b":4}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedResult: 10,
		},
		{
			name:           "multiplicacao por zero",
			method:         http.MethodPost,
			body:           `{"a":10,"b":0}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedResult: 0,
		},
		{
			name:           "multiplicacao com negativos",
			method:         http.MethodPost,
			body:           `{"a":-3,"b":5}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedResult: -15,
		},
		{
			name:           "campo obrigatorio ausente",
			method:         http.MethodPost,
			body:           `{"a":10}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
		},
		{
			name:           "operando nao numerico",
			method:         http.MethodPost,
			body:           `{"a":"dez","b":5}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
		},
		{
			name:           "json malformado",
			method:         http.MethodPost,
			body:           `{"a":10,`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
		{
			name:           "metodo nao permitido",
			method:         http.MethodGet,
			body:           "",
			contentType:    "application/json",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "METHOD_NOT_ALLOWED",
		},
		{
			name:           "content type invalido",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5}`,
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedError:  "UNSUPPORTED_MEDIA_TYPE",
		},
		{
			name:           "multiplos objetos JSON",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5} {"a":2,"b":3}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
		{
			name:           "corpo null",
			method:         http.MethodPost,
			body:           `null`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
		{
			name:           "corpo array",
			method:         http.MethodPost,
			body:           `[10, 5]`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
		{
			name:           "corpo string",
			method:         http.MethodPost,
			body:           `"teste"`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
		{
			name:           "campo desconhecido",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5,"c":2}`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				tt.method,
				"/multiplicar",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", tt.contentType)

			recorder := httptest.NewRecorder()
			MultiplyHandler(recorder, request)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"status esperado %d, obtido %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if tt.expectedStatus == http.StatusOK {
				var response multiplicationResponse

				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("resposta JSON inválida: %v", err)
				}

				if response.Result != tt.expectedResult {
					t.Errorf(
						"resultado esperado %v, obtido %v",
						tt.expectedResult,
						response.Result,
					)
				}

				return
			}

			var response errorResponse

			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("erro JSON inválido: %v", err)
			}

			if response.Error.Code != tt.expectedError {
				t.Errorf(
					"erro esperado %s, obtido %s",
					tt.expectedError,
					response.Error.Code,
				)
			}
		})
	}
}
