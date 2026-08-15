#!/usr/bin/env sh
set -eu

for platform in linux/amd64 linux/arm64; do
  tag="stock-replenishment-planner:${platform#linux/}"
  docker build --platform "$platform" -f benzhi.Dockerfile -t "$tag" .
  docker run --rm --platform "$platform" "$tag" go build ./...
  docker run --rm --platform "$platform" "$tag"
done
