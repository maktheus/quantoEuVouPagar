package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

func TestHealthCheckHandler(t *testing.T) {
	// Cria uma requisição HTTP GET para o endpoint /health
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Cria um ResponseRecorder para gravar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// Chama o handler com a requisição e o gravador
	handler.ServeHTTP(rr, req)

	// Verifica o código de status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verifica o Content-Type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestCreateSubcontaHandler_ValidJSON(t *testing.T) {
	// Cria um payload JSON válido para a subconta
	subconta := model.Subconta{
		Descricao:    "Teste de Subconta",
		Valor:        1000.0,
		TaxaJuros:    0.01,
		Parcelas:     12,
		DataInicio:   time.Now(),
		Sistema:      "price",
		Tipo:         "teste",
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}

	// Converte a struct para JSON
	jsonSubconta, err := json.Marshal(subconta)
	if err != nil {
		t.Fatal(err)
	}

	// Cria uma requisição HTTP POST com o payload JSON
	req, err := http.NewRequest("POST", "/subcontas", bytes.NewBuffer(jsonSubconta))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Cria um ResponseRecorder para gravar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateSubcontaHandler)

	// Chama o handler com a requisição e o gravador
	handler.ServeHTTP(rr, req)

	// Verifica o código de status
	// Como a lógica ainda é um placeholder, esperamos 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Verifica o Content-Type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestCreateSubcontaHandler_InvalidJSON(t *testing.T) {
	// Cria um payload JSON inválido
	invalidJSON := []byte(`{descricao: "Teste de Subconta", valor:}`)

	// Cria uma requisição HTTP POST com o payload JSON inválido
	req, err := http.NewRequest("POST", "/subcontas", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Cria um ResponseRecorder para gravar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateSubcontaHandler)

	// Chama o handler com a requisição e o gravador
	handler.ServeHTTP(rr, req)

	// Verifica o código de status
	// Esperamos 400 Bad Request para JSON inválido
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestCreateSubcontaHandler_WrongMethod(t *testing.T) {
	// Cria uma requisição HTTP GET para o endpoint /subcontas (que espera POST)
	req, err := http.NewRequest("GET", "/subcontas", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Cria um ResponseRecorder para gravar a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateSubcontaHandler)

	// Chama o handler com a requisição e o gravador
	handler.ServeHTTP(rr, req)

	// Verifica o código de status
	// Esperamos 405 Method Not Allowed
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}