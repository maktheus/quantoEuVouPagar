package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/auth"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/middleware"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/repository"
	"go.uber.org/zap"
)

// UpdateSubscriptionRequest representa os dados necessários para atualizar uma assinatura.
type UpdateSubscriptionRequest struct {
	Level string `json:"level"`
}

// UpdateSubscriptionHandler lida com a atualização do nível de assinatura de um usuário.
// @Summary Atualizar Assinatura
// @Description Atualiza o nível de assinatura de um usuário (requer permissão de admin)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID do Usuário"
// @Param subscription body UpdateSubscriptionRequest true "Novo Nível de Assinatura"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /api/admin/users/{id}/subscription [put]
func UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository, log *zap.Logger) {
	if r.Method != http.MethodPut {
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

	// Verifica se o usuário é admin
	if !auth.HasPermission(claims, auth.RoleAdmin) {
		log.Warn("Permissão insuficiente para atualizar assinatura", zap.String("user_role", string(claims.Role)), zap.String("path", r.URL.Path))
		http.Error(w, "Acesso negado. Apenas administradores podem atualizar assinaturas.", http.StatusForbidden)
		return
	}

	// Extrai o ID do usuário da URL
	// Exemplo: /api/admin/users/1/subscription
	path := r.URL.Path
	var userIDStr string
	_, err := http.Sscanf(path, "/api/admin/users/%s/subscription", &userIDStr)
	if err != nil {
		log.Warn("ID do usuário não fornecido ou inválido", zap.String("path", path))
		http.Error(w, "ID do usuário não fornecido ou inválido", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Warn("ID do usuário inválido", zap.String("id", userIDStr))
		http.Error(w, "ID do usuário inválido", http.StatusBadRequest)
		return
	}

	var req UpdateSubscriptionRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Falha ao decodificar JSON", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validação básica
	level := model.SubscriptionLevel(req.Level)
	if level != model.LevelFree && level != model.LevelTier && level != model.LevelGold && level != model.LevelPlatinum {
		log.Warn("Nível de assinatura inválido", zap.String("level", req.Level))
		http.Error(w, "Nível de assinatura inválido", http.StatusBadRequest)
		return
	}

	// Atualiza o nível de assinatura do usuário
	err = userRepo.UpdateSubscriptionLevel(userID, level)
	if err != nil {
		log.Error("Erro ao atualizar nível de assinatura", zap.Error(err))
		if err.Error() == "usuário não encontrado" {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao atualizar nível de assinatura", http.StatusInternalServerError)
		}
		return
	}

	log.Info("Nível de assinatura atualizado com sucesso", zap.Int("user_id", userID), zap.String("level", string(level)), zap.String("admin", claims.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Nível de assinatura atualizado com sucesso"})
}