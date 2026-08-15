package store

import (
	"context"
	"sync"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

type PlanStore interface {
	Save(context.Context, []domain.ReplenishmentPlan) error
	Snapshot() []domain.ReplenishmentPlan
}

type MemoryPlanStore struct {
	mu      sync.Mutex
	batches [][]domain.ReplenishmentPlan
}

func NewMemoryPlanStore() *MemoryPlanStore {
	return &MemoryPlanStore{}
}

func (s *MemoryPlanStore) Save(ctx context.Context, plans []domain.ReplenishmentPlan) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// 拷贝入参切片，避免调用方或后续批次的突变回写覆盖已保存的历史记录
	saved := make([]domain.ReplenishmentPlan, len(plans))
	copy(saved, plans)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches = append(s.batches, saved)
	return nil
}

func (s *MemoryPlanStore) Snapshot() []domain.ReplenishmentPlan {
	s.mu.Lock()
	defer s.mu.Unlock()
	var plans []domain.ReplenishmentPlan
	for _, batch := range s.batches {
		plans = append(plans, batch...)
	}
	return plans
}
