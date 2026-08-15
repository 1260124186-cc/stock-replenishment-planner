package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/zhangchengcheng/stock-replenishment-planner/internal/service"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/store"
	"github.com/zhangchengcheng/stock-replenishment-planner/internal/transport"
)

func main() {
	inputPath := flag.String("input", "", "path to a replenishment input JSON file")
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "missing required -input path")
		os.Exit(2)
	}

	input, err := transport.LoadInput(*inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	planner := service.NewPlanner(store.NewCatalog(input.Policies), store.NewMemoryPlanStore())
	plans, err := planner.Plan(context.Background(), input.Orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := transport.WritePlans(os.Stdout, plans); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
