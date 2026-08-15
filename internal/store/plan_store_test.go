package store

import (
	"context"
	"errors"
	"testing"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

func TestMemoryPlanStoreHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := NewMemoryPlanStore().Save(ctx, []domain.ReplenishmentPlan{{SKU: "tea", Quantity: 5}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Save() error = %v, want context.Canceled", err)
	}
}
