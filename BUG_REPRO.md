# 修复前故障复现（Docker）

## 项目与标准命令

库存补货计划器根据商品库存、销量和补货策略生成建议采购量。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

其中 `go test ./...` 会触发下述缺失安全库存字段的故障。

## 环境构建与编译

已实际执行以下命令；linux/amd64 和 linux/arm64 的镜像构建及容器内 `go build ./...` 均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner-delivery-002-base:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner-delivery-002-base:amd64 go build ./...

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner-delivery-002-base:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner-delivery-002-base:arm64 go build ./...
```

## 故障触发步骤

在仓库根目录先构建 linux/arm64 镜像，再执行：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner-delivery-002-base:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner-delivery-002-base:arm64 go test ./internal/service ./internal/transport
```

## 实际错误输出

```text
--- FAIL: TestPlannerUsesDefaultSafetyStockForProgrammaticPolicy (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x126108]

goroutine 9 [running]:
testing.tRunner.func1.2({0x15bc60, 0x2c8dd0})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x15bc60?, 0x2c8dd0?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
github.com/zhangchengcheng/stock-replenishment-planner/internal/domain.RecommendedQuantity(...)
	/app/internal/domain/policy.go:27
github.com/zhangchengcheng/stock-replenishment-planner/internal/service.(*Planner).Plan(0x36bb45eafec8, {0x197998, 0x2fa2a0}, {0x36bb45eaff28, 0x1, 0x0?})
	/app/internal/service/planner.go:28 +0x128
github.com/zhangchengcheng/stock-replenishment-planner/internal/service.TestPlannerUsesDefaultSafetyStockForProgrammaticPolicy(0x36bb45ee4908)
	/app/internal/service/planner_test.go:46 +0x1c4
testing.tRunner(0x36bb45ee4908, 0x1949d0)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/service	0.005s
--- FAIL: TestLoadInputAppliesDefaultSafetyStock (0.00s)
    json_test.go:48: loaded policy safety stock = (*int)(nil), want default 6
FAIL
FAIL	github.com/zhangchengcheng/stock-replenishment-planner/internal/transport	0.001s
FAIL

exit status: 1
```

## 期望行为

当商品策略未填写安全库存时，文件导入和程序化提交都应生成采用默认补货规则的计划，并以正常退出结果返回；不应中断补货计算。
