FROM --platform=$BUILDPLATFORM golang:1.26.2-bookworm

ENV GOTOOLCHAIN=local
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["go", "run", "./cmd/replenishment", "--input", "/app/examples/daily.json"]
