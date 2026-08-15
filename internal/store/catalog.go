package store

import (
	"context"
	"fmt"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

type Catalog interface {
	Lookup(context.Context, string) (domain.ReorderPolicy, error)
}

type MemoryCatalog struct {
	policies map[string]domain.ReorderPolicy
}

func NewCatalog(policies []domain.ReorderPolicy) *MemoryCatalog {
	catalog := &MemoryCatalog{policies: make(map[string]domain.ReorderPolicy, len(policies))}
	for _, policy := range policies {
		catalog.policies[policy.SKU] = policy
	}
	return catalog
}

func (c *MemoryCatalog) Lookup(ctx context.Context, sku string) (domain.ReorderPolicy, error) {
	if err := ctx.Err(); err != nil {
		return domain.ReorderPolicy{}, err
	}
	policy, ok := c.policies[sku]
	if !ok {
		return domain.ReorderPolicy{}, fmt.Errorf("%w: %s", domain.ErrUnknownSKU, sku)
	}
	return policy, nil
}
