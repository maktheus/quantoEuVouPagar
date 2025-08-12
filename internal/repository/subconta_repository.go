package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"go.uber.org/zap"
)

// SubcontaRepository define os métodos para interagir com a tabela de subcontas no banco de dados.
type SubcontaRepository struct {
	db  *sql.DB
	log *zap.Logger
}

// NewSubcontaRepository cria uma nova instância de SubcontaRepository.
func NewSubcontaRepository(db *sql.DB) *SubcontaRepository {
	return &SubcontaRepository{db: db, log: zap.L()}
}

// Create insere uma nova subconta e suas parcelas no banco de dados.
// TODO: Implementar transação para garantir consistência entre subconta e parcelas.
func (r *SubcontaRepository) Create(subconta model.Subconta, parcelas []model.Parcela) (model.Subconta, error) {
	// Inserir a subconta
	querySubconta := `
	INSERT INTO subcontas (descricao, valor, taxa_juros, parcelas, data_inicio, sistema, tipo, criado_em, atualizado_em)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	var id int
	createdAt := time.Now()
	updatedAt := time.Now()

	err := r.db.QueryRow(querySubconta,
		subconta.Descricao, subconta.Valor, subconta.TaxaJuros, subconta.Parcelas,
		subconta.DataInicio, subconta.Sistema, subconta.Tipo, createdAt, updatedAt).
		Scan(&id)
	if err != nil {
		r.log.Error("Falha ao inserir subconta", zap.Error(err))
		return subconta, fmt.Errorf("falha ao inserir subconta: %w", err)
	}

	subconta.ID = id
	subconta.CriadoEm = createdAt
	subconta.AtualizadoEm = updatedAt

	// Inserir as parcelas
	// TODO: Usar transação
	queryParcela := `
	INSERT INTO parcelas (subconta_id, numero, valor, principal, juros, data_venc, pago, criado_em, atualizado_em)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, parcela := range parcelas {
		createdAt := time.Now()
		updatedAt := time.Now()

		_, err := r.db.Exec(queryParcela,
			id, parcela.Numero, parcela.Valor, parcela.Principal, parcela.Juros,
			parcela.DataVenc, parcela.Pago, createdAt, updatedAt)
		if err != nil {
			// TODO: Em caso de erro, deveríamos fazer rollback da inserção da subconta
			r.log.Error("Falha ao inserir parcela", zap.Error(err))
			return subconta, fmt.Errorf("falha ao inserir parcela: %w", err)
		}
	}

	r.log.Info("Subconta e suas parcelas inseridas com sucesso", zap.Int("id", id))
	return subconta, nil
}

// GetAllParcelasBySubcontaID busca todas as parcelas associadas a um ID de subconta.
// TODO: Implementar
func (r *SubcontaRepository) GetAllParcelasBySubcontaID(subcontaID int) ([]model.Parcela, error) {
	// Query para buscar as parcelas
	query := `
	SELECT id, subconta_id, numero, valor, principal, juros, data_venc, pago, criado_em, atualizado_em
	FROM parcelas
	WHERE subconta_id = $1
	ORDER BY numero`

	rows, err := r.db.Query(query, subcontaID)
	if err != nil {
		r.log.Error("Falha ao buscar parcelas", zap.Error(err))
		return nil, fmt.Errorf("falha ao buscar parcelas: %w", err)
	}
	defer rows.Close()

	var parcelas []model.Parcela
	for rows.Next() {
		var p model.Parcela
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&p.ID, &p.SubcontaID, &p.Numero, &p.Valor, &p.Principal, &p.Juros,
			&p.DataVenc, &p.Pago, &createdAt, &updatedAt,
		)
		if err != nil {
			r.log.Error("Falha ao escanear parcela", zap.Error(err))
			return nil, fmt.Errorf("falha ao escanear parcela: %w", err)
		}
		// p.CriadoEm e p.AtualizadoEm não estão no model.Parcela, então não vamos preencher
		parcelas = append(parcelas, p)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("Falha ao iterar sobre as parcelas", zap.Error(err))
		return nil, fmt.Errorf("falha ao iterar sobre as parcelas: %w", err)
	}

	r.log.Info("Parcelas obtidas com sucesso", zap.Int("subconta_id", subcontaID), zap.Int("count", len(parcelas)))
	return parcelas, nil
}