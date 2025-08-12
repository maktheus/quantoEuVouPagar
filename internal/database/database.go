package database

import (
	"database/sql"
	"fmt"
	"log"

	// Driver para PostgreSQL
	_ "github.com/lib/pq"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/config"
)

// Connect cria uma conexão com o banco de dados PostgreSQL.
func Connect(cfg *config.Config) (*sql.DB, error) {
	// Monta a string de conexão
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Abre a conexão com o banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com o banco de dados: %w", err)
	}

	// Verifica se a conexão é válida
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao pingar o banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return db, nil
}