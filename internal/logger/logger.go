package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger inicializa o logger.
func InitLogger() {
	// Configuração do logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err = config.Build()
	if err != nil {
		log.Fatalf("Falha ao inicializar o logger: %v", err)
	}
}

// GetLogger retorna a instância do logger.
func GetLogger() *zap.Logger {
	if logger == nil {
		// Se o logger não foi inicializado, cria um logger de desenvolvimento
		l, _ := zap.NewDevelopment()
		return l
	}
	return logger
}