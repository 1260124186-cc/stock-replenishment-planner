# 修复前故障复现（Docker）

## 项目与标准命令

库存补货计划器根据商品现货、在途采购、日销量和补货策略生成建议采购量。以下步骤使用未修复的提交 `b35884da3b6e666a0665d85458e9011b49136903`，并在仓库根目录执行：

```sh
git checkout b35884da3b6e666a0665d85458e9011b49136903
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

其中 `go build ./...` 可以成功，正常示例能够输出计划 JSON；`go test ./...` 会因下述无效策略测试失败。

## 环境构建与编译

已实际执行以下命令。`linux/amd64` 和 `linux/arm64` 的镜像均构建成功，且各自容器内的 `go build ./...` 均成功。

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner-bug003-base:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner-bug003-base:amd64 go build ./...

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner-bug003-base:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner-bug003-base:arm64 go build ./...
```

## 故障触发步骤

在上述未修复提交的仓库根目录执行：

```sh
docker run --rm --platform linux/amd64 stock-replenishment-planner-bug003-base:amd64 go test ./internal/service ./internal/transport
```

## 实际错误输出

```text
--- FAIL: TestPlannerReturnsErrorForInvalidProgrammaticPolicy (0.00s)
panic: runtime error: integer divide by zero [recovered, repanicked]

goroutine 9 [running]:
testing.tRunner.func1.2({0x15bc60, 0x2c8da0})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x15bc60?, 0x2c8da0?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
github.com/zhangchengcheng/stock-replenishment-planner/internal/domain.roundUp(...)
	/app/internal/domain/policy.go:42
github.com/zhangchengcheng/stock-replenishment-planner/internal/domain.RecommendedQuantity(...)
	/app/internal/domain/policy.go:38
github.com/zhangchengcheng/stock-replenishment-planner/internal/service.(*Planner).Plan(0x3efea86e5ec8, {0x1979d8, 0x2fa2a0}, {0x3efea86e5f28, 0x1, 0xebbef?})
	/app/internal/service/planner.go:28 +0x3cc
github.com/zhangchengcheng/stock-replenishment-planner/internal/service.TestPlannerReturnsErrorForInvalidProgrammaticPolicy(0x3efea879a908)
	/app/internal/service/planner_test.go:47 +0x1d0
testing.tRunner(0x3efea879a908, 0x194a00)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/service	0.005s
--- FAIL: TestLoadInputRejectsInvalidPolicy (0.00s)
    json_test.go:44: LoadInput() error = nil, want invalid policy error
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/transport	0.001s
FAIL
```

命令以退出码 `1` 结束。

## 期望行为

补货倍数为零的策略应在文件导入和程序化调用时返回清晰错误，不应进入补货计算或导致批次中断。
