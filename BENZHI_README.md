# 库存补货计划器

库存补货计划器根据商品的现货、在途采购、日销量和补货策略生成建议采购量，帮助仓库在覆盖周期内避免缺货。项目可在本机或容器中直接从源码构建和运行，不依赖外部网络服务或数据库。

在仓库根目录执行：

```sh
go build ./...
go run ./cmd/replenishment --input ./examples/daily.json
go test ./...
```

## 固定环境

项目的语言版本固定在 `go.mod` 的 `go 1.26.0`。实际 Docker 文件为 `benzhi.Dockerfile`，使用 `golang:1.26.2-bookworm` 并设置 `GOTOOLCHAIN=local`，因此容器不会自动切换到其他 Go 工具链。镜像保留完整 Go 工具链，所有编译都在容器内从源码执行。

以下命令分别验证 `linux/amd64` 和 `linux/arm64`：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t stock-replenishment-planner:amd64 .
docker run --rm --platform linux/amd64 stock-replenishment-planner:amd64 go build ./...
docker run --rm --platform linux/amd64 stock-replenishment-planner:amd64

docker build --platform linux/arm64 -f benzhi.Dockerfile -t stock-replenishment-planner:arm64 .
docker run --rm --platform linux/arm64 stock-replenishment-planner:arm64 go build ./...
docker run --rm --platform linux/arm64 stock-replenishment-planner:arm64
```

也可以执行：

```sh
./build_benzhi_docker.sh
```

脚本会依次为两个平台构建镜像，在每个容器内运行 `go build ./...`，然后启动项目入口并输出计划 JSON。手工进行测试时，可在已构建镜像中执行：

```sh
docker run --rm --platform linux/amd64 stock-replenishment-planner:amd64 go test ./...
docker run --rm --platform linux/arm64 stock-replenishment-planner:arm64 go test ./...
```

当本机或容器内的 `go build ./...` 与 `go test ./...` 均以退出码 `0` 结束，且运行入口能输出包含 `plans` 的 JSON 时，表示环境与项目行为通过。
