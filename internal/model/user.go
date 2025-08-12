package model

import (
	"time"
)

// UserRole representa o papel do usuário.
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// SubscriptionLevel representa o nível de assinatura do usuário.
type SubscriptionLevel string

const (
	LevelFree     SubscriptionLevel = "free"
	LevelTier     SubscriptionLevel = "tier"
	LevelGold     SubscriptionLevel = "gold"
	LevelPlatinum SubscriptionLevel = "platinum"
)

// User representa um usuário do sistema.
type User struct {
	ID               int               `json:"id"`
	Username         string            `json:"username"`
	Email            string            `json:"email"`
	Password         string            `json:"-"` // Não serializar a senha
	Role             UserRole          `json:"role"`              // "user" ou "admin"
	SubscriptionLevel SubscriptionLevel `json:"subscription_level"` // "free", "tier", "gold", "platinum"
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}