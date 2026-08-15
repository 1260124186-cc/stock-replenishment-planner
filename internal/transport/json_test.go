package transport

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
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

func TestLoadInputRejectsInvalidPolicy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	input := `{"orders":[],"policies":[{"sku":"tea","minimum_stock":1,"reorder_multiple":0,"lead_time_days":1,"safety_stock":1}]}`
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadInput(path)
	if err == nil {
		t.Fatal("LoadInput() error = nil, want invalid policy error")
	}
	if !strings.Contains(err.Error(), "validate policy 0") || !strings.Contains(err.Error(), "invalid limits") {
		t.Fatalf("LoadInput() error = %v, want clear policy validation error", err)
	}
}
