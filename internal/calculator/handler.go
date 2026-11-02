package calculator

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
)

type multiplicationRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type multiplicationResponse struct {
	Result float64 `json:"resultado"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"codigo"`
		Message string `json:"mensagem"`
	} `json:"erro"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := errorResponse{}
	response.Error.Code = code
	response.Error.Message = message

	_ = json.NewEncoder(w).Encode(response)
}

func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED", "Método não permitido.")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType,
			"UNSUPPORTED_MEDIA_TYPE", "O conteúdo deve ser JSON.")
		return
	}

	defer r.Body.Close()

	var raw json.RawMessage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest,
			"INVALID_JSON", "JSON inválido.")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest,
			"INVALID_JSON", "Apenas um objeto JSON é permitido.")
		return
	}

	raw = bytes.TrimSpace(raw)

	if len(raw) == 0 || raw[0] != '{' {
		writeError(w, http.StatusBadRequest,
			"INVALID_JSON", "O corpo deve ser um objeto JSON.")
		return
	}

	var request multiplicationRequest
	requestDecoder := json.NewDecoder(bytes.NewReader(raw))
	requestDecoder.DisallowUnknownFields()

	if err := requestDecoder.Decode(&request); err != nil {
		var typeError *json.UnmarshalTypeError

		if errors.As(err, &typeError) {
			writeError(w, http.StatusBadRequest,
				"INVALID_INPUT", "Os operandos devem ser números.")
			return
		}

		writeError(w, http.StatusBadRequest,
			"INVALID_JSON", "JSON inválido.")
		return
	}

	if request.A == nil || request.B == nil {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT", "Os operandos a e b são obrigatórios.")
		return
	}

	result := *request.A * *request.B

	if math.IsInf(result, 0) || math.IsNaN(result) {
		writeError(w, http.StatusUnprocessableEntity,
			"INVALID_RESULT", "Resultado fora do intervalo permitido.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(multiplicationResponse{
		Result: result,
	})
}
