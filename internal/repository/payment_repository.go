package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"go.uber.org/zap"
)

// PaymentRepository define os métodos para interagir com a tabela de pagamentos no banco de dados.
type PaymentRepository struct {
	db  *sql.DB
	log *zap.Logger
}

// NewPaymentRepository cria uma nova instância de PaymentRepository.
func NewPaymentRepository(db *sql.DB, log *zap.Logger) *PaymentRepository {
	return &PaymentRepository{db: db, log: log}
}

// CreatePayment insere um novo pagamento no banco de dados.
func (r *PaymentRepository) CreatePayment(payment model.Payment) (model.Payment, error) {
	query := `
	INSERT INTO payments (user_id, amount, currency, status, method, transaction_id, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	var id int
	createdAt := time.Now()
	updatedAt := time.Now()

	err := r.db.QueryRow(query,
		payment.UserID, payment.Amount, payment.Currency, payment.Status, payment.Method,
		payment.TransactionID, payment.Description, createdAt, updatedAt).
		Scan(&id)
	if err != nil {
		r.log.Error("Falha ao inserir pagamento", zap.Error(err))
		return payment, fmt.Errorf("falha ao inserir pagamento: %w", err)
	}

	payment.ID = id
	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	r.log.Info("Pagamento criado com sucesso", zap.Int("id", id))
	return payment, nil
}

// GetPaymentByID busca um pagamento pelo ID.
func (r *PaymentRepository) GetPaymentByID(paymentID int) (model.Payment, error) {
	var payment model.Payment
	var createdAt, updatedAt time.Time

	query := `
	SELECT id, user_id, amount, currency, status, method, transaction_id, description, created_at, updated_at
	FROM payments
	WHERE id = $1`

	err := r.db.QueryRow(query, paymentID).Scan(
		&payment.ID, &payment.UserID, &payment.Amount, &payment.Currency, &payment.Status,
		&payment.Method, &payment.TransactionID, &payment.Description, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return payment, fmt.Errorf("pagamento não encontrado")
		}
		r.log.Error("Falha ao buscar pagamento", zap.Error(err))
		return payment, fmt.Errorf("falha ao buscar pagamento: %w", err)
	}

	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	return payment, nil
}

// UpdatePaymentStatus atualiza o status de um pagamento.
func (r *PaymentRepository) UpdatePaymentStatus(paymentID int, status model.PaymentStatus) error {
	query := `
	UPDATE payments
	SET status = $1, updated_at = $2
	WHERE id = $3`

	updatedAt := time.Now()
	result, err := r.db.Exec(query, status, updatedAt, paymentID)
	if err != nil {
		r.log.Error("Falha ao atualizar status do pagamento", zap.Error(err))
		return fmt.Errorf("falha ao atualizar status do pagamento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("Falha ao obter número de linhas afetadas", zap.Error(err))
		return fmt.Errorf("falha ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pagamento não encontrado")
	}

	r.log.Info("Status do pagamento atualizado com sucesso", zap.Int("id", paymentID), zap.String("status", string(status)))
	return nil
}