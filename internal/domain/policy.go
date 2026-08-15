package domain

import (
	"fmt"
	"strings"
)

func NormalizePolicy(policy ReorderPolicy) (ReorderPolicy, error) {
	policy.SKU = strings.TrimSpace(policy.SKU)
	if policy.SKU == "" {
		return ReorderPolicy{}, fmt.Errorf("policy SKU is required")
	}
	if policy.MinimumStock < 0 || policy.ReorderMultiple <= 0 || policy.LeadTimeDays < 0 {
		return ReorderPolicy{}, fmt.Errorf("policy for %s has invalid limits", policy.SKU)
	}
	if policy.SafetyStock == nil {
		defaultSafetyStock := policy.MinimumStock
		policy.SafetyStock = &defaultSafetyStock
	}
	if *policy.SafetyStock < 0 {
		return ReorderPolicy{}, fmt.Errorf("policy for %s has negative safety stock", policy.SKU)
	}
	return policy, nil
}

func RecommendedQuantity(signal OrderSignal, policy ReorderPolicy) int {
	coverageTarget := signal.DailySales*(signal.DaysOfCover+policy.LeadTimeDays) + *policy.SafetyStock
	if coverageTarget < policy.MinimumStock {
		coverageTarget = policy.MinimumStock
	}

	available := signal.OnHand + signal.OpenPurchase
	if available >= coverageTarget {
		return 0
	}

	shortfall := coverageTarget - available
	return roundUp(shortfall, policy.ReorderMultiple)
}

func roundUp(value, multiple int) int {
	return ((value + multiple - 1) / multiple) * multiple
}
