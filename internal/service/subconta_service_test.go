package service

import (
	"testing"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

func TestCalcularParcelas_Price(t *testing.T) {
	subconta := model.Subconta{
		ID:          1,
		Descricao:   "Teste Price",
		Valor:       10000.0,  // Valor financiado
		TaxaJuros:   0.01,     // 1% ao mês
		Parcelas:    12,       // 12 meses
		DataInicio:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		Sistema:     "price",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	parcelas, err := CalcularParcelas(subconta)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if len(parcelas) != subconta.Parcelas {
		t.Errorf("Esperava %d parcelas, obteve %d", subconta.Parcelas, len(parcelas))
	}

	// Verifica se a primeira parcela tem o valor esperado (aproximadamente)
	// Usando uma margem de erro de 0.01 para lidar com arredondamentos
	expectedValorPrimeiraParcela := 888.49 // Valor calculado previamente
	if parcelas[0].Valor < expectedValorPrimeiraParcela-0.01 || parcelas[0].Valor > expectedValorPrimeiraParcela+0.01 {
		t.Errorf("Valor da primeira parcela incorreto. Esperava ~%.2f, obteve %.2f", expectedValorPrimeiraParcela, parcelas[0].Valor)
	}

	// Verifica se a última parcela tem o valor esperado (aproximadamente)
	expectedValorUltimaParcela := 843.48 // Valor calculado previamente
	if parcelas[len(parcelas)-1].Valor < expectedValorUltimaParcela-0.01 || parcelas[len(parcelas)-1].Valor > expectedValorUltimaParcela+0.01 {
		t.Errorf("Valor da última parcela incorreto. Esperava ~%.2f, obteve %.2f", expectedValorUltimaParcela, parcelas[len(parcelas)-1].Valor)
	}

	// Verifica se a soma das amortizações é igual ao valor financiado
	somaAmortizacoes := 0.0
	for _, p := range parcelas {
		somaAmortizacoes += p.Principal
	}
	// Usando uma margem de erro de 0.01 para lidar com arredondamentos
	if somaAmortizacoes < subconta.Valor-0.01 || somaAmortizacoes > subconta.Valor+0.01 {
		t.Errorf("Soma das amortizações incorreta. Esperava %.2f, obteve %.2f", subconta.Valor, somaAmortizacoes)
	}
}

func TestCalcularParcelas_SAC(t *testing.T) {
	subconta := model.Subconta{
		ID:          1,
		Descricao:   "Teste SAC",
		Valor:       10000.0,  // Valor financiado
		TaxaJuros:   0.01,     // 1% ao mês
		Parcelas:    12,       // 12 meses
		DataInicio:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		Sistema:     "sac",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	parcelas, err := CalcularParcelas(subconta)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if len(parcelas) != subconta.Parcelas {
		t.Errorf("Esperava %d parcelas, obteve %d", subconta.Parcelas, len(parcelas))
	}

	// No SAC, a amortização é constante
	amortizacaoEsperada := subconta.Valor / float64(subconta.Parcelas)
	for _, p := range parcelas {
		// Usando uma margem de erro de 0.01 para lidar com arredondamentos
		if p.Principal < amortizacaoEsperada-0.01 || p.Principal > amortizacaoEsperada+0.01 {
			t.Errorf("Amortização incorreta. Esperava ~%.2f, obteve %.2f", amortizacaoEsperada, p.Principal)
		}
	}

	// Verifica se a primeira parcela tem o valor esperado
	// Primeira parcela = Amortização + Juros (sobre o valor total)
	expectedValorPrimeiraParcela := amortizacaoEsperada + (subconta.Valor * subconta.TaxaJuros)
	if parcelas[0].Valor < expectedValorPrimeiraParcela-0.01 || parcelas[0].Valor > expectedValorPrimeiraParcela+0.01 {
		t.Errorf("Valor da primeira parcela incorreto. Esperava ~%.2f, obteve %.2f", expectedValorPrimeiraParcela, parcelas[0].Valor)
	}

	// Verifica se a última parcela tem o valor esperado
	// Última parcela = Amortização + Juros (sobre o saldo devedor restante)
	saldoRestante := subconta.Valor - (amortizacaoEsperada * float64(subconta.Parcelas-1))
	expectedValorUltimaParcela := amortizacaoEsperada + (saldoRestante * subconta.TaxaJuros)
	if parcelas[len(parcelas)-1].Valor < expectedValorUltimaParcela-0.01 || parcelas[len(parcelas)-1].Valor > expectedValorUltimaParcela+0.01 {
		t.Errorf("Valor da última parcela incorreto. Esperava ~%.2f, obteve %.2f", expectedValorUltimaParcela, parcelas[len(parcelas)-1].Valor)
	}

	// Verifica se a soma das amortizações é igual ao valor financiado
	somaAmortizacoes := 0.0
	for _, p := range parcelas {
		somaAmortizacoes += p.Principal
	}
	// Usando uma margem de erro de 0.01 para lidar com arredondamentos
	if somaAmortizacoes < subconta.Valor-0.01 || somaAmortizacoes > subconta.Valor+0.01 {
		t.Errorf("Soma das amortizações incorreta. Esperava %.2f, obteve %.2f", subconta.Valor, somaAmortizacoes)
	}
}

func TestCalcularParcelas_SistemaInvalido(t *testing.T) {
	subconta := model.Subconta{
		ID:          1,
		Descricao:   "Teste Sistema Inválido",
		Valor:       10000.0,
		TaxaJuros:   0.01,
		Parcelas:    12,
		DataInicio:  time.Now(),
		Sistema:     "sistema_invalido",
		Tipo:        "teste",
		CriadoEm:    time.Now(),
		AtualizadoEm: time.Now(),
	}

	_, err := CalcularParcelas(subconta)
	if err == nil {
		t.Error("Esperava um erro para sistema inválido, mas não obteve nenhum")
	}
}