package store

import (
	"context"
	"testing"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

func TestMemoryPlanStoreDoesNotRetainCallerSlice(t *testing.T) {
	saved := []domain.ReplenishmentPlan{{SKU: "tea", Quantity: 5}}
	store := NewMemoryPlanStore()
	if err := store.Save(context.Background(), saved); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	saved[0].Quantity = 99

	if snapshot := store.Snapshot(); snapshot[0].Quantity != 5 {
		t.Fatalf("snapshot quantity = %d, want 5 after caller mutation", snapshot[0].Quantity)
	}
}
