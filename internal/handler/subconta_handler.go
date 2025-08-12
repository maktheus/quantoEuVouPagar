package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/auth"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/repository"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/service"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/middleware"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/subscription"
	"go.uber.org/zap"
)

// HealthCheckHandler retorna um status 200 OK para indicar que o serviço está funcionando.
// @Summary Health Check
// @Description Verifica se o serviço está funcionando
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// TODO: Verificar a conexão com o banco de dados
	// if err := db.Ping(); err != nil {
	//     http.Error(w, "Database connection failed", http.StatusInternalServerError)
	//     return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Quanto Eu Vou Pagar? service is up and running"})
}

// CreateSubcontaHandler lida com a criação de uma nova subconta.
// @Summary Criar Subconta
// @Description Cria uma nova subconta e calcula suas parcelas
// @Tags subcontas
// @Accept json
// @Produce json
// @Param subconta body model.Subconta true "Dados da Subconta"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /api/subcontas [post]
func CreateSubcontaHandler(w http.ResponseWriter, r *http.Request, repo *repository.SubcontaRepository, userRepo *repository.UserRepository, log *zap.Logger) {
	if r.Method != http.MethodPost {
		log.Warn("Método não permitido", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verifica as permissões do usuário
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Usuário não autenticado", zap.String("path", r.URL.Path))
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	// Busca o usuário completo para verificar o nível de assinatura
	user, err := userRepo.GetByUsername(claims.Username)
	if err != nil {
		log.Error("Erro ao buscar usuário", zap.Error(err))
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	// Por exemplo, apenas usuários "admin" podem criar subcontas
	// Esta é uma regra de negócio arbitrária, pode ser modificada
	// Vamos mudar para verificar a feature
	if !auth.HasPermission(claims, auth.RoleAdmin) {
		// Verifica se o usuário tem a feature para criar subcontas
		if !subscription.HasFeature(user, subscription.FeatureCreateSubcontas) {
			log.Warn("Permissão insuficiente para criar subcontas", zap.String("user_role", string(claims.Role)), zap.String("subscription_level", string(user.SubscriptionLevel)), zap.String("path", r.URL.Path))
			http.Error(w, "Acesso negado. Você não tem permissão para criar subcontas.", http.StatusForbidden)
			return
		}

		// Se o usuário for Free, verifica o limite de subcontas
		if user.SubscriptionLevel == model.LevelFree {
			// Aqui você precisaria contar o número de subcontas do usuário
			// Por simplicidade, vamos assumir que não há limite na criação, mas o acesso à listagem pode ser limitado
			// Ou você pode implementar uma lógica para contar as subcontas
		}
	}

	var subconta model.Subconta
	err = json.NewDecoder(r.Body).Decode(&subconta)
	if err != nil {
		log.Error("Falha ao decodificar JSON", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	log.Info("Recebida requisição para criar subconta", zap.String("descricao", subconta.Descricao), zap.String("user", claims.Username))

	// Valida os dados da subconta
	if err := service.ValidarSubconta(subconta); err != nil {
		log.Warn("Dados da subconta inválidos", zap.Error(err))
		http.Error(w, fmt.Sprintf("Dados da subconta inválidos: %v", err), http.StatusBadRequest)
		return
	}

	// Se o usuário for Free, verifica o valor da subconta (limite de R$ 10.000,00)
	if user.SubscriptionLevel == model.LevelFree && subconta.Valor > 10000.0 {
		log.Warn("Valor da subconta excede o limite para o nível Free", zap.Float64("valor", subconta.Valor), zap.String("user", claims.Username))
		http.Error(w, "Valor da subconta excede o limite para o nível Free (R$ 10.000,00)", http.StatusForbidden)
		return
	}

	// Calcula as parcelas usando o serviço
	parcelas, err := service.CalcularParcelas(subconta)
	if err != nil {
		log.Error("Erro ao calcular parcelas", zap.Error(err))
		http.Error(w, fmt.Sprintf("Erro ao calcular parcelas: %v", err), http.StatusBadRequest)
		return
	}

	// Salva a subconta e as parcelas no banco de dados
	subconta, err = repo.Create(subconta, parcelas)
	if err != nil {
		log.Error("Erro ao salvar subconta", zap.Error(err))
		http.Error(w, fmt.Sprintf("Erro ao salvar subconta: %v", err), http.StatusInternalServerError)
		return
	}

	log.Info("Subconta criada com sucesso", zap.Int("id", subconta.ID), zap.String("user", claims.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Subconta criada com sucesso", "subconta": subconta, "parcelas": parcelas})
}

// GetParcelasHandler lida com a obtenção das parcelas de uma subconta.
// @Summary Obter Parcelas de uma Subconta
// @Description Obtém todas as parcelas associadas a uma subconta específica
// @Tags subcontas
// @Produce json
// @Param id path int true "ID da Subconta"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /api/subcontas/{id}/parcelas [get]
func GetParcelasHandler(w http.ResponseWriter, r *http.Request, repo *repository.SubcontaRepository, userRepo *repository.UserRepository, log *zap.Logger) {
	if r.Method != http.MethodGet {
		log.Warn("Método não permitido", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verifica as permissões do usuário
	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Usuário não autenticado", zap.String("path", r.URL.Path))
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	// Busca o usuário completo para verificar o nível de assinatura
	user, err := userRepo.GetByUsername(claims.Username)
	if err != nil {
		log.Error("Erro ao buscar usuário", zap.Error(err))
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	// Verifica se o usuário tem a feature para visualizar parcelas
	if !subscription.HasFeature(user, subscription.FeatureViewParcelas) {
		log.Warn("Permissão insuficiente para visualizar parcelas", zap.String("user_role", string(claims.Role)), zap.String("subscription_level", string(user.SubscriptionLevel)), zap.String("path", r.URL.Path))
		http.Error(w, "Acesso negado. Você não tem permissão para visualizar parcelas.", http.StatusForbidden)
		return
	}

	// Extrai o ID da subconta da URL
	// Exemplo: /api/subcontas/1/parcelas
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		log.Warn("ID da subconta não fornecido")
		http.Error(w, "ID da subconta não fornecido", http.StatusBadRequest)
		return
	}

	subcontaIDStr := parts[1] // parts[0] é "subcontas"
	subcontaID, err := strconv.Atoi(subcontaIDStr)
	if err != nil {
		log.Warn("ID da subconta inválido", zap.String("id", subcontaIDStr))
		http.Error(w, "ID da subconta inválido", http.StatusBadRequest)
		return
	}

	log.Info("Recebida requisição para obter parcelas", zap.Int("subconta_id", subcontaID), zap.String("user", claims.Username))

	// Busca as parcelas no banco de dados
	parcelas, err := repo.GetAllParcelasBySubcontaID(subcontaID)
	if err != nil {
		log.Error("Erro ao buscar parcelas", zap.Error(err))
		http.Error(w, fmt.Sprintf("Erro ao buscar parcelas: %v", err), http.StatusInternalServerError)
		return
	}

	log.Info("Parcelas obtidas com sucesso", zap.Int("subconta_id", subcontaID), zap.Int("count", len(parcelas)), zap.String("user", claims.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"parcelas": parcelas})
}