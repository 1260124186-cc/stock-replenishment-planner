package service

import (
	"context"
	"fmt"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/store"
)

type Planner struct {
	catalog store.Catalog
	plans   store.PlanStore
}

func NewPlanner(catalog store.Catalog, plans store.PlanStore) *Planner {
	return &Planner{catalog: catalog, plans: plans}
}

func (p *Planner) Plan(ctx context.Context, signals []domain.OrderSignal) ([]domain.ReplenishmentPlan, error) {
	plans := make([]domain.ReplenishmentPlan, 0, len(signals))
	for _, signal := range signals {
		policy, err := p.catalog.Lookup(ctx, signal.SKU)
		if err != nil {
			return nil, fmt.Errorf("load reorder policy for %s: %w", signal.SKU, err)
		}

		quantity := domain.RecommendedQuantity(signal, policy)
		if quantity == 0 {
			continue
		}
		plans = append(plans, domain.ReplenishmentPlan{
			SKU:      signal.SKU,
			Quantity: quantity,
			Reason:   "coverage below reorder target",
		})
	}

	if len(plans) == 0 {
		return plans, nil
	}

	if err := p.plans.Save(ctx, plans); err != nil {
		return nil, fmt.Errorf("save replenishment plans: %w", err)
	}
	return plans, nil
}
