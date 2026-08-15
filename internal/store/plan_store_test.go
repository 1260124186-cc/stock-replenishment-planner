package store

import (
	"context"
	"testing"
	"time"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

func TestMemoryPlanStoreUnlocksAfterEmptyBatch(t *testing.T) {
	store := NewMemoryPlanStore()
	if err := store.Save(context.Background(), nil); err != nil {
		t.Fatalf("empty Save() error = %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- store.Save(context.Background(), []domain.ReplenishmentPlan{{SKU: "tea", Quantity: 5}})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("follow-up Save() error = %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("follow-up Save() blocked after an empty batch")
	}
}
