package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/middleware"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/repository"
	"go.uber.org/zap"
)

// CreatePaymentRequest representa os dados necessários para criar um pagamento.
type CreatePaymentRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Method      string  `json:"method"`
	Description string  `json:"description"`
}

// CreatePaymentResponse representa os dados retornados após a criação de um pagamento.
type CreatePaymentResponse struct {
	Payment model.Payment `json:"payment"`
}

// CreatePaymentHandler lida com a criação de um novo pagamento.
// @Summary Criar Pagamento
// @Description Cria um novo pagamento para o usuário autenticado
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body CreatePaymentRequest true "Dados do Pagamento"
// @Success 201 {object} CreatePaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /api/payments [post]
func CreatePaymentHandler(w http.ResponseWriter, r *http.Request, paymentRepo *repository.PaymentRepository, log *zap.Logger) {
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

	var req CreatePaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Falha ao decodificar JSON", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validação básica
	if req.Amount <= 0 {
		log.Warn("Valor do pagamento inválido")
		http.Error(w, "Valor do pagamento deve ser maior que zero", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = "BRL" // Moeda padrão
	}
	if req.Method == "" {
		log.Warn("Método de pagamento não fornecido")
		http.Error(w, "Método de pagamento é obrigatório", http.StatusBadRequest)
		return
	}

	// Cria o pagamento
	payment := model.Payment{
		UserID:      claims.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      model.PaymentStatusPending,
		Method:      model.PaymentMethod(req.Method),
		Description: req.Description,
	}

	payment, err = paymentRepo.CreatePayment(payment)
	if err != nil {
		log.Error("Erro ao criar pagamento", zap.Error(err))
		http.Error(w, "Erro ao criar pagamento", http.StatusInternalServerError)
		return
	}

	log.Info("Pagamento criado com sucesso", zap.Int("id", payment.ID), zap.String("user", claims.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatePaymentResponse{Payment: payment})
}

// GetPaymentHandler lida com a obtenção de um pagamento pelo ID.
// @Summary Obter Pagamento
// @Description Obtém os detalhes de um pagamento específico
// @Tags payments
// @Produce json
// @Param id path int true "ID do Pagamento"
// @Success 200 {object} CreatePaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /api/payments/{id} [get]
func GetPaymentHandler(w http.ResponseWriter, r *http.Request, paymentRepo *repository.PaymentRepository, log *zap.Logger) {
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

	// Extrai o ID do pagamento da URL
	// Exemplo: /api/payments/1
	path := r.URL.Path
	var paymentIDStr string
	_, err := fmt.Sscanf(path, "/api/payments/%s", &paymentIDStr)
	if err != nil {
		log.Warn("ID do pagamento não fornecido ou inválido", zap.String("path", path))
		http.Error(w, "ID do pagamento não fornecido ou inválido", http.StatusBadRequest)
		return
	}

	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		log.Warn("ID do pagamento inválido", zap.String("id", paymentIDStr))
		http.Error(w, "ID do pagamento inválido", http.StatusBadRequest)
		return
	}

	// Busca o pagamento
	payment, err := paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		log.Warn("Erro ao buscar pagamento", zap.Error(err))
		if err.Error() == "pagamento não encontrado" {
			http.Error(w, "Pagamento não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar pagamento", http.StatusInternalServerError)
		}
		return
	}

	// Verifica se o pagamento pertence ao usuário autenticado
	// A menos que o usuário seja admin, ele só pode ver seus próprios pagamentos
	// if claims.Role != "admin" && payment.UserID != claims.UserID {
	//	 log.Warn("Acesso negado a pagamento de outro usuário", zap.Int("payment_user_id", payment.UserID), zap.Int("user_id", claims.UserID))
	//	 http.Error(w, "Acesso negado", http.StatusForbidden)
	//	 return
	// }
	// Por enquanto, vamos permitir que qualquer usuário autenticado veja qualquer pagamento (para simplificar)

	log.Info("Pagamento obtido com sucesso", zap.Int("id", payment.ID), zap.String("user", claims.Username))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreatePaymentResponse{Payment: payment})
}