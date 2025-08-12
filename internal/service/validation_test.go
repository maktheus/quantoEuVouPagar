package service

import (
	"testing"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

func TestValidarSubconta_Valida(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta válida",
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err != nil {
		t.Errorf("Esperava que a subconta fosse válida, mas obteve erro: %v", err)
	}
}

func TestValidarSubconta_DescricaoVazia(t *testing.T) {
	subconta := model.Subconta{
		// Descrição vazia
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para descrição vazia, mas não obteve nenhum")
	}
}

func TestValidarSubconta_ValorZero(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com valor zero",
		Valor:       0.0, // Valor zero
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para valor zero, mas não obteve nenhum")
	}
}

func TestValidarSubconta_ValorNegativo(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com valor negativo",
		Valor:       -1000.0, // Valor negativo
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para valor negativo, mas não obteve nenhum")
	}
}

func TestValidarSubconta_TaxaJurosNegativa(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com taxa de juros negativa",
		Valor:       1000.0,
		TaxaJuros:   -0.01, // Taxa de juros negativa
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para taxa de juros negativa, mas não obteve nenhum")
	}
}

func TestValidarSubconta_ParcelasZero(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com zero parcelas",
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    0, // Zero parcelas
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para zero parcelas, mas não obteve nenhum")
	}
}

func TestValidarSubconta_DataInicioZero(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com data de início zero",
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Time{}, // Data de início zero
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para data de início zero, mas não obteve nenhum")
	}
}

func TestValidarSubconta_SistemaInvalido(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com sistema inválido",
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "sistema_invalido", // Sistema inválido
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para sistema inválido, mas não obteve nenhum")
	}
}

func TestValidarSubconta_TipoVazio(t *testing.T) {
	subconta := model.Subconta{
		Descricao:   "Subconta com tipo vazio",
		Valor:       1000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "price",
		// Tipo vazio
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	err := ValidarSubconta(subconta)
	if err == nil {
		t.Error("Esperava erro para tipo vazio, mas não obteve nenhum")
	}
}