package transport

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/service"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/store"
)

func TestLoadInputAndWritePlans(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	input := `{"orders":[{"sku":"tea","on_hand":1,"open_purchase":0,"daily_sales":2,"days_of_cover":2}],"policies":[{"sku":"tea","minimum_stock":3,"reorder_multiple":5,"lead_time_days":1,"safety_stock":2}]}`
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadInput(path)
	if err != nil {
		t.Fatalf("LoadInput() error = %v", err)
	}
	if len(loaded.Orders) != 1 || loaded.Policies[0].SKU != "tea" {
		t.Fatalf("LoadInput() = %#v", loaded)
	}

	var output bytes.Buffer
	if err := WritePlans(&output, []domain.ReplenishmentPlan{{SKU: "tea", Quantity: 10}}); err != nil {
		t.Fatalf("WritePlans() error = %v", err)
	}
	if !bytes.Contains(output.Bytes(), []byte(`"quantity": 10`)) {
		t.Fatalf("WritePlans() output = %s", output.String())
	}
}

func TestLoadInputAppliesDefaultSafetyStock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	input := `{"orders":[{"sku":"tea","on_hand":0,"daily_sales":2,"days_of_cover":1}],"policies":[{"sku":"tea","minimum_stock":6,"reorder_multiple":5,"lead_time_days":1}]}`
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadInput(path)
	if err != nil {
		t.Fatalf("LoadInput() error = %v", err)
	}
	if loaded.Policies[0].SafetyStock == nil || *loaded.Policies[0].SafetyStock != 6 {
		t.Fatalf("loaded policy safety stock = %#v, want default 6", loaded.Policies[0].SafetyStock)
	}

	planner := service.NewPlanner(store.NewCatalog(loaded.Policies), store.NewMemoryPlanStore())
	plans, err := planner.Plan(context.Background(), loaded.Orders)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plans) != 1 || plans[0].Quantity != 10 {
		t.Fatalf("Plan() = %#v, want one quantity-10 plan", plans)
	}
}
