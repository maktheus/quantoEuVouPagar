package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository define os métodos para interagir com a tabela de usuários no banco de dados.
type UserRepository struct {
	db  *sql.DB
	log *zap.Logger
}

// NewUserRepository cria uma nova instância de UserRepository.
func NewUserRepository(db *sql.DB, log *zap.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

// Create insere um novo usuário no banco de dados.
func (r *UserRepository) Create(username, email, password string, role model.UserRole) (model.User, error) {
	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.log.Error("Falha ao gerar hash da senha", zap.Error(err))
		return model.User{}, fmt.Errorf("falha ao gerar hash da senha: %w", err)
	}

	// Inserir o usuário
	// Por padrão, todos os novos usuários começam com o nível "free"
	query := `
	INSERT INTO users (username, email, password, role, subscription_level, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id int
	createdAt := time.Now()
	updatedAt := time.Now()

	err = r.db.QueryRow(query, username, email, string(hashedPassword), string(role), string(model.LevelFree), createdAt, updatedAt).Scan(&id)
	if err != nil {
		r.log.Error("Falha ao inserir usuário", zap.Error(err))
		return model.User{}, fmt.Errorf("falha ao inserir usuário: %w", err)
	}

	user := model.User{
		ID:                id,
		Username:          username,
		Email:             email,
		Password:          "", // Não retornamos a senha
		Role:              role,
		SubscriptionLevel: model.LevelFree, // Nível inicial
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	r.log.Info("Usuário criado com sucesso", zap.Int("id", id))
	return user, nil
}

// GetByUsername busca um usuário pelo nome de usuário.
func (r *UserRepository) GetByUsername(username string) (model.User, error) {
	var user model.User
	var hashedPassword string
	var role string
	var subscriptionLevel string
	var createdAt, updatedAt time.Time

	query := `
	SELECT id, username, email, password, role, subscription_level, created_at, updated_at
	FROM users
	WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &hashedPassword, &role, &subscriptionLevel, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("usuário não encontrado")
		}
		r.log.Error("Falha ao buscar usuário", zap.Error(err))
		return user, fmt.Errorf("falha ao buscar usuário: %w", err)
	}

	user.Password = hashedPassword
	user.Role = model.UserRole(role)
	user.SubscriptionLevel = model.SubscriptionLevel(subscriptionLevel)
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, nil
}

// Authenticate verifica se o nome de usuário e senha são válidos.
func (r *UserRepository) Authenticate(username, password string) (model.User, error) {
	user, err := r.GetByUsername(username)
	if err != nil {
		return user, err
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("senha inválida")
	}

	// Limpa a senha antes de retornar
	user.Password = ""
	return user, nil
}

// UpdateSubscriptionLevel atualiza o nível de assinatura de um usuário.
func (r *UserRepository) UpdateSubscriptionLevel(userID int, level model.SubscriptionLevel) error {
	query := `
	UPDATE users
	SET subscription_level = $1, updated_at = $2
	WHERE id = $3`

	updatedAt := time.Now()
	result, err := r.db.Exec(query, string(level), updatedAt, userID)
	if err != nil {
		r.log.Error("Falha ao atualizar nível de assinatura do usuário", zap.Error(err))
		return fmt.Errorf("falha ao atualizar nível de assinatura do usuário: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("Falha ao obter número de linhas afetadas", zap.Error(err))
		return fmt.Errorf("falha ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	r.log.Info("Nível de assinatura do usuário atualizado com sucesso", zap.Int("id", userID), zap.String("level", string(level)))
	return nil
}