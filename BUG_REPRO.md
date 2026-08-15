# 修复前故障复现（Docker）

## 项目与标准命令

库存补货计划器根据商品库存、在途采购、日销量和补货策略生成建议采购量。以下命令在仓库根目录执行：

```sh
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

## 环境构建与编译

在修复前的代码状态下，以下命令已实际执行。两个平台的镜像构建和容器内编译均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner:amd64 go build ./...

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner:arm64 go build ./...
```

## 故障触发步骤

在修复前的代码状态下，于仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/amd64 stock-replenishment-planner:amd64 go test ./internal/service ./internal/store
```

## 实际错误输出

```text
--- FAIL: TestPlannerKeepsEarlierBatchStableAfterLaterPlan (0.00s)
    planner_test.go:61: first plan quantity = 6, want 3 after later plan
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/service	0.007s
--- FAIL: TestMemoryPlanStoreDoesNotRetainCallerSlice (0.00s)
    plan_store_test.go:19: snapshot quantity = 99, want 5 after caller mutation
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/store	0.013s
FAIL
```

该命令的退出状态为 `1`。

## 期望行为

同一进程生成后续补货批次或调用方修改原始输入后，已经返回的首批计划和历史记录都应保持生成当时的商品、数量与原因，不应被后续操作覆盖。
