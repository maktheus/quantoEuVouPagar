package model

import (
	"time"
)

// Subconta representa uma parcela ou parte de uma compra maior.
// Todos os juros são tratados internamente como juros mensais.
type Subconta struct {
	ID          int       `json:"id"`
	Descricao   string    `json:"descricao"`             // Descrição da subconta
	Valor       float64   `json:"valor"`                 // Valor principal da subconta
	TaxaJuros   float64   `json:"taxa_juros"`            // Taxa de juros mensal (já convertida se for anual)
	Parcelas    int       `json:"parcelas"`              // Número total de parcelas
	DataInicio  time.Time `json:"data_inicio"`           // Data de início do pagamento
	Sistema     string    `json:"sistema"`               // Sistema de amortização (ex: "price", "sac")
	Tipo        string    `json:"tipo"`                  // Tipo (ex: "financiamento", "boleto")
	CriadoEm    time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Parcela representa uma parcela individual de uma Subconta.
type Parcela struct {
	ID         int       `json:"id"`
	SubcontaID int       `json:"subconta_id"` // Foreign Key para Subconta
	Numero     int       `json:"numero"`      // Número da parcela (1, 2, 3, ...)
	Valor      float64   `json:"valor"`       // Valor total da parcela (principal + juros)
	Principal  float64   `json:"principal"`   // Parte do valor que amortiza a dívida
	Juros      float64   `json:"juros"`       // Parte do valor referente aos juros
	DataVenc   time.Time `json:"data_venc"`   // Data de vencimento da parcela
	Pago       bool      `json:"pago"`        // Indica se a parcela foi paga
}

// Investimento representa um investimento com aportes.
type Investimento struct {
	ID             int       `json:"id"`
	Descricao      string    `json:"descricao"`       // Descrição do investimento
	ValorInicial   float64   `json:"valor_inicial"`   // Valor inicial do investimento
	TaxaRendimento float64   `json:"taxa_rendimento"` // Taxa de rendimento mensal (já convertida se for anual)
	Aportes        []Aporte  `json:"aportes"`         // Lista de aportes realizados
	CriadoEm       time.Time `json:"criado_em"`
	AtualizadoEm   time.Time `json:"atualizado_em"`
}

// Aporte representa um aporte realizado em um Investimento.
type Aporte struct {
	ID             int       `json:"id"`
	InvestimentoID int       `json:"investimento_id"` // Foreign Key para Investimento
	Valor          float64   `json:"valor"`           // Valor do aporte
	Data           time.Time `json:"data"`            // Data do aporte
}

// Divida representa uma dívida que pode ser amortizada.
// Pode ser uma subconta ou uma dívida avulsa.
type Divida struct {
	ID           int       `json:"id"`
	Descricao    string    `json:"descricao"`     // Descrição da dívida
	ValorDevido  float64   `json:"valor_devido"`  // Valor total devido
	TaxaJuros    float64   `json:"taxa_juros"`    // Taxa de juros mensal
	DataInicio   time.Time `json:"data_inicio"`   // Data de início da dívida
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Amortizacao representa uma amortização aplicada a uma Dívida.
type Amortizacao struct {
	ID        int       `json:"id"`
	DividaID  int       `json:"divida_id"`   // Foreign Key para Divida
	Valor     float64   `json:"valor"`       // Valor amortizado
	Data      time.Time `json:"data"`        // Data da amortização
	Tipo      string    `json:"tipo"`        // Tipo de amortização (ex: "reduzir_valor", "reduzir_prazo")
}