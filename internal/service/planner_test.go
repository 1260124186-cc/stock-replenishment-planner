package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/store"
)

func TestPlannerCreatesAndPersistsPlans(t *testing.T) {
	safetyStock := 2
	plans := store.NewMemoryPlanStore()
	planner := NewPlanner(store.NewCatalog([]domain.ReorderPolicy{{
		SKU: "tea", MinimumStock: 4, ReorderMultiple: 5, LeadTimeDays: 1, SafetyStock: &safetyStock,
	}}), plans)

	got, err := planner.Plan(context.Background(), []domain.OrderSignal{{
		SKU: "tea", OnHand: 2, DailySales: 2, DaysOfCover: 2,
	}})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(got) != 1 || got[0].Quantity != 10 {
		t.Fatalf("Plan() = %#v, want one quantity-10 plan", got)
	}
	if snapshot := plans.Snapshot(); len(snapshot) != 1 || snapshot[0] != got[0] {
		t.Fatalf("persisted plans = %#v, want %#v", snapshot, got)
	}
}

func TestPlannerPreservesUnknownSKUError(t *testing.T) {
	planner := NewPlanner(store.NewCatalog(nil), store.NewMemoryPlanStore())
	_, err := planner.Plan(context.Background(), []domain.OrderSignal{{SKU: "missing"}})
	if !errors.Is(err, domain.ErrUnknownSKU) {
		t.Fatalf("Plan() error = %v, want unknown SKU", err)
	}
}

func TestPlannerReturnsErrorForInvalidProgrammaticPolicy(t *testing.T) {
	safetyStock := 1
	plans := store.NewMemoryPlanStore()
	planner := NewPlanner(store.NewCatalog([]domain.ReorderPolicy{{
		SKU: "tea", MinimumStock: 1, ReorderMultiple: 0, SafetyStock: &safetyStock,
	}}), plans)

	_, err := planner.Plan(context.Background(), []domain.OrderSignal{{
		SKU: "tea", OnHand: 0, DailySales: 2, DaysOfCover: 1,
	}})
	if err == nil {
		t.Fatal("Plan() error = nil, want invalid policy error")
	}
	if !strings.Contains(err.Error(), "validate reorder policy for tea") || !strings.Contains(err.Error(), "invalid limits") {
		t.Fatalf("Plan() error = %v, want clear policy validation error", err)
	}
	if snapshot := plans.Snapshot(); len(snapshot) != 0 {
		t.Fatalf("persisted plans = %#v, want no plans after validation failure", snapshot)
	}
}
