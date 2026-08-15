# 修复前故障复现（Docker）

## 项目与标准命令

库存补货计划器根据商品库存、在途采购、日销量和补货策略生成建议采购量。仓库根目录中的标准命令如下：

```sh
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

修复前，完整测试会在空批次路径上失败。

## 环境构建与编译

修复前状态已实际执行以下命令：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner-005-base:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner-005-base:amd64 go build ./...

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner-005-base:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner-005-base:arm64 go build ./...
```

两个平台的镜像构建和容器内编译均成功，故障在下面的测试命令中触发。

## 故障触发步骤

在仓库根目录使用修复前状态构建镜像后，执行：

```sh
docker run --rm --platform linux/arm64 stock-replenishment-planner-005-base:arm64 go test ./internal/service ./internal/store
```

该命令会先处理一个没有补货项的批次，再验证后续保存操作和空批次保存行为。

## 实际错误输出

```text
--- FAIL: TestPlannerDoesNotPersistEmptyPlanBatch (0.00s)
    planner_test.go:53: Save() calls = 1, want 0 for an empty plan batch
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/service	0.003s
--- FAIL: TestMemoryPlanStoreUnlocksAfterEmptyBatch (0.10s)
    plan_store_test.go:28: follow-up Save() blocked after an empty batch
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/store	0.104s
FAIL
exit_status=1
```

## 期望行为

没有任何补货项的批次不应被记录为有效计划；之后提交的正常补货请求应及时完成并正常保存结果。
