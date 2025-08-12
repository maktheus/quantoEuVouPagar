package subscription

import (
	"fmt"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

// Feature representa uma funcionalidade do sistema.
type Feature string

const (
	// Funcionalidades básicas (Free)
	FeatureViewParcelas       Feature = "view_parcelas"
	FeatureSimulateJuros      Feature = "simulate_juros"
	FeatureCreateSubcontas    Feature = "create_subcontas"
	FeatureBasicDocumentation Feature = "basic_documentation"

	// Funcionalidades do Tier
	FeatureUnlimitedSubcontas     Feature = "unlimited_subcontas"
	FeatureUnlimitedSimulateJuros Feature = "unlimited_simulate_juros"
	FeatureTrackInvestments       Feature = "track_investments"
	FeatureManageDebtsBasic       Feature = "manage_debts_basic"
	FeatureExportPDFLimited       Feature = "export_pdf_limited"

	// Funcionalidades do Gold
	FeatureManageDebtsAdvanced Feature = "manage_debts_advanced"
	FeatureTrackInvestmentsAdvanced Feature = "track_investments_advanced"
	FeatureVisualizationAdvanced    Feature = "visualization_advanced"
	FeatureExportPDFUnlimited       Feature = "export_pdf_unlimited"
	FeaturePrioritySupport          Feature = "priority_support"

	// Funcionalidades do Platinum
	FeatureAIAnalysis          Feature = "ai_analysis"
	FeatureCalendarIntegration Feature = "calendar_integration"
	FeatureCustomReports       Feature = "custom_reports"
	FeatureSupport247          Feature = "support_247"
)

// featureAccessMap define quais funcionalidades estão disponíveis para cada nível de assinatura.
var featureAccessMap = map[model.SubscriptionLevel][]Feature{
	model.LevelFree: {
		FeatureViewParcelas,
		FeatureSimulateJuros, // Com limite
		FeatureCreateSubcontas, // Com limite (3)
		FeatureBasicDocumentation,
	},
	model.LevelTier: {
		// Todas as features do Free
		FeatureViewParcelas,
		FeatureSimulateJuros,
		FeatureCreateSubcontas,
		FeatureBasicDocumentation,
		// Features adicionais do Tier
		FeatureUnlimitedSubcontas,
		FeatureUnlimitedSimulateJuros,
		FeatureTrackInvestments,
		FeatureManageDebtsBasic,
		FeatureExportPDFLimited, // 5 relatórios por mês
	},
	model.LevelGold: {
		// Todas as features do Tier
		FeatureViewParcelas,
		FeatureSimulateJuros,
		FeatureCreateSubcontas,
		FeatureBasicDocumentation,
		FeatureUnlimitedSubcontas,
		FeatureUnlimitedSimulateJuros,
		FeatureTrackInvestments,
		FeatureManageDebtsBasic,
		FeatureExportPDFLimited,
		// Features adicionais do Gold
		FeatureManageDebtsAdvanced,
		FeatureTrackInvestmentsAdvanced,
		FeatureVisualizationAdvanced,
		FeatureExportPDFUnlimited,
		FeaturePrioritySupport,
	},
	model.LevelPlatinum: {
		// Todas as features do Gold
		FeatureViewParcelas,
		FeatureSimulateJuros,
		FeatureCreateSubcontas,
		FeatureBasicDocumentation,
		FeatureUnlimitedSubcontas,
		FeatureUnlimitedSimulateJuros,
		FeatureTrackInvestments,
		FeatureManageDebtsBasic,
		FeatureExportPDFLimited,
		FeatureManageDebtsAdvanced,
		FeatureTrackInvestmentsAdvanced,
		FeatureVisualizationAdvanced,
		FeatureExportPDFUnlimited,
		FeaturePrioritySupport,
		// Features adicionais do Platinum
		FeatureAIAnalysis,
		FeatureCalendarIntegration,
		FeatureCustomReports,
		FeatureSupport247,
	},
}

// HasFeature verifica se um usuário tem acesso a uma determinada funcionalidade.
func HasFeature(user model.User, feature Feature) bool {
	// Usuários admin têm acesso a todas as funcionalidades
	if user.Role == model.RoleAdmin {
		return true
	}

	// Obtém a lista de features disponíveis para o nível de assinatura do usuário
	availableFeatures, ok := featureAccessMap[user.SubscriptionLevel]
	if !ok {
		// Se o nível de assinatura não estiver mapeado, nega o acesso
		return false
	}

	// Verifica se a feature está na lista de features disponíveis
	for _, f := range availableFeatures {
		if f == feature {
			return true
		}
	}

	return false
}

// CheckFeatureAndLimit verifica se o usuário tem acesso a uma funcionalidade e, se for o caso, verifica um limite numérico.
// Por exemplo, verificar se o usuário pode criar mais subcontas.
func CheckFeatureAndLimit(user model.User, feature Feature, currentValue, limit int) error {
	// Verifica se o usuário tem acesso à funcionalidade
	if !HasFeature(user, feature) {
		return fmt.Errorf("acesso negado: funcionalidade '%s' não disponível para o nível '%s'", feature, user.SubscriptionLevel)
	}

	// Se o usuário tem acesso ilimitado, não há limite
	// (A lógica de "ilimitado" é implementada nas regras de negócio, não aqui)
	// Por exemplo, "unlimited_subcontas" significa que não há limite no número de subcontas.

	// Se houver um limite, verifica se o valor atual excede o limite
	// Esta parte é um exemplo e pode ser adaptada conforme a necessidade
	switch feature {
	case FeatureCreateSubcontas:
		// Para o nível Free, o limite é 3 subcontas
		if user.SubscriptionLevel == model.LevelFree && currentValue >= limit {
			return fmt.Errorf("limite de subcontas atingido para o nível '%s'. Limite: %d", user.SubscriptionLevel, limit)
		}
	case FeatureSimulateJuros:
		// Para o nível Free, o limite é R$ 10.000,00
		// Este limite é verificado na lógica de simulação de juros, não aqui
		// Esta função é mais para verificar limites baseados em contagem
	}

	// Se não houver erro, o acesso é permitido
	return nil
}