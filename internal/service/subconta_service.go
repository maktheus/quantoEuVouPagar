package service

import (
	"fmt"
	"math"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

// CalcularParcelas calcula as parcelas de uma subconta com base no sistema de amortização (Price ou SAC).
func CalcularParcelas(subconta model.Subconta) ([]model.Parcela, error) {
	switch subconta.Sistema {
	case "price":
		return calcularParcelasPrice(subconta)
	case "sac":
		return calcularParcelasSAC(subconta)
	default:
		return nil, fmt.Errorf("sistema de amortização '%s' não suportado", subconta.Sistema)
	}
}

// calcularParcelasPrice calcula as parcelas usando o Sistema Price (amortização constante).
func calcularParcelasPrice(subconta model.Subconta) ([]model.Parcela, error) {
	var parcelas []model.Parcela
	valorPresente := subconta.Valor
	taxaJuros := subconta.TaxaJuros
	numeroParcelas := subconta.Parcelas

	// Fórmula do Sistema Price para calcular o valor da parcela:
	// PMT = PV * [i * (1 + i)^n] / [(1 + i)^n - 1]
	// Onde:
	// PMT = valor da parcela
	// PV = valor presente (valor financiado)
	// i = taxa de juros
	// n = número de parcelas

	// Calcula (1 + i)^n
	potencia := math.Pow(1+taxaJuros, float64(numeroParcelas))

	// Calcula o valor da parcela
	valorParcela := valorPresente * (taxaJuros * potencia) / (potencia - 1)

	saldoDevedor := valorPresente

	for i := 1; i <= numeroParcelas; i++ {
		// Calcula os juros sobre o saldo devedor
		juros := saldoDevedor * taxaJuros

		// Calcula a parte da parcela que amortiza a dívida
		amortizacao := valorParcela - juros

		// Atualiza o saldo devedor
		saldoDevedor -= amortizacao

		// Cria a parcela
		parcela := model.Parcela{
			ID:         i, // ID temporário, será definido pelo banco de dados
			SubcontaID: subconta.ID,
			Numero:     i,
			Valor:      valorParcela,
			Principal:  amortizacao,
			Juros:      juros,
			DataVenc:   subconta.DataInicio.AddDate(0, i-1, 0), // Adiciona meses
			Pago:       false,
		}

		parcelas = append(parcelas, parcela)
	}

	return parcelas, nil
}

// calcularParcelasSAC calcula as parcelas usando o Sistema SAC (amortização constante).
func calcularParcelasSAC(subconta model.Subconta) ([]model.Parcela, error) {
	var parcelas []model.Parcela
	valorPresente := subconta.Valor
	taxaJuros := subconta.TaxaJuros
	numeroParcelas := subconta.Parcelas

	// No SAC, a amortização é constante
	amortizacao := valorPresente / float64(numeroParcelas)

	saldoDevedor := valorPresente

	for i := 1; i <= numeroParcelas; i++ {
		// Calcula os juros sobre o saldo devedor
		juros := saldoDevedor * taxaJuros

		// Calcula o valor da parcela (amortização + juros)
		valorParcela := amortizacao + juros

		// Atualiza o saldo devedor
		saldoDevedor -= amortizacao

		// Cria a parcela
		parcela := model.Parcela{
			ID:         i, // ID temporário, será definido pelo banco de dados
			SubcontaID: subconta.ID,
			Numero:     i,
			Valor:      valorParcela,
			Principal:  amortizacao,
			Juros:      juros,
			DataVenc:   subconta.DataInicio.AddDate(0, i-1, 0), // Adiciona meses
			Pago:       false,
		}

		parcelas = append(parcelas, parcela)
	}

	return parcelas, nil
}