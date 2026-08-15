package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/domain"
)

type Input struct {
	Orders   []domain.OrderSignal   `json:"orders"`
	Policies []domain.ReorderPolicy `json:"policies"`
}

func LoadInput(path string) (Input, error) {
	file, err := os.Open(path)
	if err != nil {
		return Input{}, fmt.Errorf("open input: %w", err)
	}
	defer file.Close()

	var input Input
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return Input{}, fmt.Errorf("decode input: %w", err)
	}
	for index, policy := range input.Policies {
		normalizedPolicy, err := domain.NormalizePolicy(policy)
		if err != nil {
			return Input{}, fmt.Errorf("validate policy %d: %w", index, err)
		}
		input.Policies[index] = normalizedPolicy
	}
	return input, nil
}

func WritePlans(writer io.Writer, plans []domain.ReplenishmentPlan) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		Plans []domain.ReplenishmentPlan `json:"plans"`
	}{Plans: plans})
}
