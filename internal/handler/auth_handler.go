package handler

import (
	"encoding/json"
	"net/http"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/auth"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/repository"
	"go.uber.org/zap"
)

// LoginRequest representa os dados necessários para login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse representa os dados retornados após um login bem-sucedido.
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterRequest representa os dados necessários para registro.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler lida com o registro de novos usuários.
// @Summary Registrar Usuário
// @Description Registra um novo usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Dados do Usuário"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository, log *zap.Logger) {
	if r.Method != http.MethodPost {
		log.Warn("Método não permitido", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Falha ao decodificar JSON", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validação básica (pode ser expandida)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Warn("Dados de registro inválidos")
		http.Error(w, "Username, email e password são obrigatórios", http.StatusBadRequest)
		return
	}

	// Cria o usuário (por padrão, todos os novos usuários são "user", não "admin")
	// A criação de usuários admin deve ser feita de outra forma, por exemplo, via script ou endpoint específico
	user, err := userRepo.Create(req.Username, req.Email, req.Password, string(auth.RoleUser))
	if err != nil {
		log.Error("Erro ao criar usuário", zap.Error(err))
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	log.Info("Usuário registrado com sucesso", zap.Int("id", user.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Usuário registrado com sucesso", "user": user})
}

// LoginHandler lida com o login de usuários.
// @Summary Login de Usuário
// @Description Autentica um usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciais do Usuário"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository, log *zap.Logger) {
	if r.Method != http.MethodPost {
		log.Warn("Método não permitido", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Falha ao decodificar JSON", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validação básica
	if req.Username == "" || req.Password == "" {
		log.Warn("Credenciais de login inválidas")
		http.Error(w, "Username e password são obrigatórios", http.StatusBadRequest)
		return
	}

	// Autentica o usuário
	user, err := userRepo.Authenticate(req.Username, req.Password)
	if err != nil {
		log.Warn("Falha na autenticação", zap.Error(err))
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Gera o token JWT
	token, err := auth.GenerateToken(user.ID, user.Username, auth.UserRole(user.Role))
	if err != nil {
		log.Error("Erro ao gerar token", zap.Error(err))
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	log.Info("Usuário autenticado com sucesso", zap.Int("id", user.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}