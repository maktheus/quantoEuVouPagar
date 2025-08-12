package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/auth"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/config"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/database"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/handler"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/logger"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/middleware"
	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/repository"
	"go.uber.org/zap"

	// Importa o pacote docs gerado pelo swag
	_ "github.com/matheus-uchoa/quanto-eu-vou-pagar/cmd/server/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Quanto Eu Vou Pagar? API
// @version         1.0
// @description     Esta API permite gerenciar e prever os custos relacionados a grandes compras, como a aquisição de uma casa.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Matheus Uchôa
// @contact.url    https://github.com/matheus-uchoa
// @contact.email  matheus.uchoa@example.com

// @license.name  MIT
// @license.url   https://github.com/matheus-uchoa/quanto-eu-vou-pagar/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 "Type 'Bearer TOKEN' to correctly set the API Key"

func main() {
	// Inicializa o logger
	logger.InitLogger()
	log := logger.GetLogger()
	defer log.Sync() // flushes buffer, if any

	// Carrega as configurações
	cfg := config.LoadConfig()

	// Define a chave secreta para o JWT
	auth.SetSecretKey(cfg.JWTSecret)

	// Conecta ao banco de dados
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados", zap.Error(err))
	}
	defer db.Close()

	// Cria os repositórios
	subcontaRepo := repository.NewSubcontaRepository(db)
	userRepo := repository.NewUserRepository(db, log)
	paymentRepo := repository.NewPaymentRepository(db, log)

	// Cria um mux para as rotas da aplicação
	mux := http.NewServeMux()

	// Configuração das rotas da aplicação
	// Rotas públicas (não requerem autenticação)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.HealthCheckHandler(w, r, db)
	})
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterHandler(w, r, userRepo, log)
	})
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		handler.LoginHandler(w, r, userRepo, log)
	})

	// Cria um sub-mux para rotas protegidas
	protectedMux := http.NewServeMux()
	// Rotas protegidas (requerem autenticação)
	protectedMux.HandleFunc("/subcontas", func(w http.ResponseWriter, r *http.Request) {
		// Passando os repositórios e o logger para o handler
		handler.CreateSubcontaHandler(w, r, subcontaRepo, userRepo, log)
	})
	protectedMux.HandleFunc("/subcontas/", func(w http.ResponseWriter, r *http.Request) {
		// Passando os repositórios e o logger para o handler
		handler.GetParcelasHandler(w, r, subcontaRepo, userRepo, log)
	})
	// Rotas de pagamento
	protectedMux.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		handler.CreatePaymentHandler(w, r, paymentRepo, log)
	})
	protectedMux.HandleFunc("/payments/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetPaymentHandler(w, r, paymentRepo, log)
	})
	// Rotas de administração (requerem autenticação e permissão de admin)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateSubscriptionHandler(w, r, userRepo, log)
	})
	// Aplica o middleware de autenticação às rotas de administração
	adminHandler := middleware.AuthMiddleware(log)(adminMux)
	protectedMux.Handle("/admin/", http.StripPrefix("/admin", adminHandler))

	// Aplica o middleware de autenticação às rotas protegidas
	authMiddleware := middleware.AuthMiddleware(log)
	protectedHandler := authMiddleware(protectedMux)

	// Registra o handler protegido no mux principal
	// Isso significa que todas as rotas começando com /api/ requerem autenticação
	mux.Handle("/api/", http.StripPrefix("/api", protectedHandler))

	// Rota para a documentação do Swagger
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL para o JSON da documentação
	))

	// Rota para as métricas do Prometheus
	mux.Handle("/metrics", middleware.MetricsHandler())

	// Cria um handler principal com o middleware de métricas
	mainHandler := middleware.MetricsMiddleware(mux)

	// Configura o servidor
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mainHandler,
	}

	// Canal para receber sinais do sistema operacional
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Goroutine para iniciar o servidor
	go func() {
		log.Info("Servidor iniciado", zap.String("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Falha ao iniciar o servidor", zap.Error(err))
		}
	}()

	// Goroutine para esperar o sinal de parada
	go func() {
		<-stop
		log.Info("Recebido sinal de interrupção, desligando o servidor...")

		// Cria um contexto com timeout para o shutdown
		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()

		// Tenta desligar o servidor graciosamente
		// Por enquanto, vamos apenas parar o programa
		// if err := server.Shutdown(ctx); err != nil {
		//	 log.Fatalf("Falha ao desligar o servidor graciosamente: %v", err)
		// }

		// Fechar o banco de dados
		db.Close()

		os.Exit(0)
	}()

	// Mantém o programa rodando
	select {}
}