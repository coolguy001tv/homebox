
# Build Nuxt
# 前端产物 (.output/public) 是纯静态文件，与目标架构无关，
# 钉在 BUILDPLATFORM 原生构建，避免 QEMU 模拟下 pnpm install/nuxt build 超时
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder
WORKDIR  /app
RUN npm install -g pnpm@9
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --shamefully-hoist
COPY frontend .
RUN pnpm build

# Build API
# CGO 已禁用，Go 可原生交叉编译：钉在 BUILDPLATFORM，按 TARGETPLATFORM 注入 GOOS/GOARCH/GOARM，
# 避免 QEMU 模拟下 go build 超时
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG BUILD_TIME
ARG COMMIT
ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
RUN apk update && \
    apk upgrade && \
    apk add --update git build-base gcc g++

WORKDIR /go/src/app
COPY ./backend .
RUN go get -d -v ./...
RUN rm -rf ./app/api/public
COPY --from=frontend-builder /app/.output/public ./app/api/static/public
# arm/v7 平台的 TARGETVARIANT 为 "v7"，GOARM 需要 "7"，去掉前缀 v
RUN set -eux; \
    export GOOS="$TARGETOS" GOARCH="$TARGETARCH"; \
    if [ -n "$TARGETVARIANT" ]; then export GOARM="${TARGETVARIANT#v}"; fi; \
    CGO_ENABLED=0 go build \
    -ldflags "-s -w -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME -X main.version=$VERSION"  \
    -o /go/bin/api \
    -v ./app/api/*.go

# Production Stage
FROM alpine:3.21

ENV HBOX_MODE=production
ENV HBOX_STORAGE_DATA=/data/
ENV HBOX_STORAGE_SQLITE_URL=/data/homebox.db?_pragma=busy_timeout=2000&_pragma=journal_mode=WAL&_fk=1

RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /go/bin/api /app

RUN chmod +x /app/api

LABEL Name=homebox Version=0.0.1
LABEL org.opencontainers.image.source="https://github.com/hay-kot/homebox"
EXPOSE 7745
WORKDIR /app
VOLUME [ "/data" ]

ENTRYPOINT [ "/app/api" ]
CMD [ "/data/config.yml" ]
