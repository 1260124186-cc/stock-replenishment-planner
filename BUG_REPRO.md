# 修复前故障复现（Docker）

## 项目与标准命令

库存补货计划器根据商品库存、销量和补货策略生成补货建议。进入仓库根目录后，可执行以下标准命令：

```sh
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

## 环境构建与编译

分别执行以下命令构建两个平台的镜像，并在容器内从源码编译项目：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner-delivery-004-base:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner-delivery-004-base:amd64 go build ./...

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner-delivery-004-base:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner-delivery-004-base:arm64 go build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功。目标故障在下节的测试命令中触发。

## 故障触发步骤

在仓库根目录先按上述 `linux/amd64` 命令构建镜像，再执行：

```sh
docker run --rm --platform linux/amd64 stock-replenishment-planner-delivery-004-base:amd64 go test ./internal/service ./internal/store
```

## 实际错误输出

```text
--- FAIL: TestPlannerStopsWhenRequestContextIsCanceled (0.00s)
    planner_test.go:52: Plan() error = <nil>, want context.Canceled
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/service	0.001s
--- FAIL: TestMemoryPlanStoreHonorsCanceledContext (0.00s)
    plan_store_test.go:17: Save() error = <nil>, want context.Canceled
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/store	0.002s
FAIL
```

命令以退出码 `1` 结束。

## 期望行为

当补货请求已经取消时，生成流程应返回 `context.Canceled`，且该批次不应继续保存为历史计划。
