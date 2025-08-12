package model

import (
	"testing"
	"time"
)

func TestSubcontaCreation(t *testing.T) {
	// Teste básico para verificar se a struct Subconta pode ser instanciada
	subconta := Subconta{
		ID:          1,
		Descricao:   "Financiamento de casa",
		Valor:       300000.0,
		TaxaJuros:   0.01, // 1% ao mês
		Parcelas:    120,  // 10 anos
		DataInicio:  time.Now(),
		Sistema:     "price",
		Tipo:        "financiamento",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	if subconta.Descricao != "Financiamento de casa" {
		t.Errorf("Esperava 'Financiamento de casa', obteve '%s'", subconta.Descricao)
	}

	if subconta.Valor != 300000.0 {
		t.Errorf("Esperava 300000.0, obteve %f", subconta.Valor)
	}

	// Adicione mais verificações conforme necessário
}

func TestInvestimentoCreation(t *testing.T) {
	// Teste básico para verificar se a struct Investimento pode ser instanciada
	investimento := Investimento{
		ID:             1,
		Descricao:      "Investimento inicial",
		ValorInicial:   10000.0,
		TaxaRendimento: 0.005, // 0.5% ao mês
		Aportes:        []Aporte{},
		CriadoEm:       time.Now(),
		AtualizadoEm:   time.Now(),
	}

	if investimento.Descricao != "Investimento inicial" {
		t.Errorf("Esperava 'Investimento inicial', obteve '%s'", investimento.Descricao)
	}

	if investimento.ValorInicial != 10000.0 {
		t.Errorf("Esperava 10000.0, obteve %f", investimento.ValorInicial)
	}
}

func TestDividaCreation(t *testing.T) {
	// Teste básico para verificar se a struct Divida pode ser instanciada
	divida := Divida{
		ID:           1,
		Descricao:    "Empréstimo pessoal",
		ValorDevido:  5000.0,
		TaxaJuros:    0.02, // 2% ao mês
		DataInicio:   time.Now(),
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}

	if divida.Descricao != "Empréstimo pessoal" {
		t.Errorf("Esperava 'Empréstimo pessoal', obteve '%s'", divida.Descricao)
	}

	if divida.ValorDevido != 5000.0 {
		t.Errorf("Esperava 5000.0, obteve %f", divida.ValorDevido)
	}
}