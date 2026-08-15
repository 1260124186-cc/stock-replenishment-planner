package domain

import "testing"

func TestRecommendedQuantityUsesLeadTimeAndSafetyStock(t *testing.T) {
	safetyStock := 4
	quantity := RecommendedQuantity(
		OrderSignal{SKU: "tea", OnHand: 4, OpenPurchase: 1, DailySales: 3, DaysOfCover: 2},
		ReorderPolicy{SKU: "tea", MinimumStock: 5, ReorderMultiple: 5, LeadTimeDays: 2, SafetyStock: &safetyStock},
	)

	if quantity != 15 {
		t.Fatalf("quantity = %d, want 15", quantity)
	}
}

func TestNormalizePolicyDefaultsSafetyStock(t *testing.T) {
	policy, err := NormalizePolicy(ReorderPolicy{SKU: "tea", MinimumStock: 6, ReorderMultiple: 5})
	if err != nil {
		t.Fatalf("NormalizePolicy() error = %v", err)
	}
	if *policy.SafetyStock != 6 {
		t.Fatalf("safety stock = %d, want 6", *policy.SafetyStock)
	}
}
