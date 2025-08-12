package model

import (
	"time"
)

// PaymentStatus representa o status de um pagamento.
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod representa o método de pagamento.
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBoleto     PaymentMethod = "boleto"
	PaymentMethodPix        PaymentMethod = "pix"
)

// Payment representa um pagamento.
type Payment struct {
	ID            int           `json:"id"`
	UserID        int           `json:"user_id"`         // ID do usuário que fez o pagamento
	Amount        float64       `json:"amount"`          // Valor do pagamento
	Currency      string        `json:"currency"`        // Moeda (ex: "BRL")
	Status        PaymentStatus `json:"status"`          // Status do pagamento
	Method        PaymentMethod `json:"method"`          // Método de pagamento
	TransactionID string        `json:"transaction_id"`  // ID da transação no gateway de pagamento
	Description   string        `json:"description"`     // Descrição do pagamento
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Subscription representa uma assinatura.
type Subscription struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`      // ID do usuário
	PlanID     int       `json:"plan_id"`      // ID do plano
	Status     string    `json:"status"`       // "active", "inactive", "cancelled"
	StartDate  time.Time `json:"start_date"`   // Data de início
	EndDate    time.Time `json:"end_date"`     // Data de término (pode ser zero para assinaturas contínuas)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Plan representa um plano de assinatura.
type Plan struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`         // Nome do plano
	Description string  `json:"description"`  // Descrição do plano
	Price       float64 `json:"price"`        // Preço do plano
	Currency    string  `json:"currency"`     // Moeda (ex: "BRL")
	Interval    string  `json:"interval"`     // Intervalo de cobrança (ex: "monthly", "yearly")
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}