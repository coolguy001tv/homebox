# HomeBox 项目说明

## 项目概况

- **技术栈**: Go + Ent ORM + SQLite + Chi Router (后端) + Nuxt 3.6.5 + Vue 3.5 + DaisyUI 2 + TypeScript (前端)
- **部署**: Docker, 前端 Nuxt 构建产物嵌入 Go binary
- **当前分支**: `main`（历史 `v0.10.3-release` 分支已冻结，不再使用）
- **Docker Hub**: `hellocoolguy/homebox:0.1.x`
- **镜像发布**: 发布方式见 README「Docker 镜像发布」。不要手动构建——`scripts/release.sh` 仅应急/本地验证用。

## 项目结构

```
backend/app/api/          — HTTP handlers, 路由, 中间件
backend/internal/core/services/ — 业务逻辑层
backend/internal/data/repo/     — 数据访问层 (Ent ORM 封装)
backend/internal/data/ent/      — Ent 生成代码 (大部分不可手动修改)
backend/internal/data/ent/schema/ — Ent schema 定义 (可手动修改)
backend/pkgs/              — 独立工具包 (mailer, hasher, pathlib, set)
frontend/pages/            — Nuxt 页面组件
frontend/components/       — Vue 组件
frontend/lib/api/          — API 客户端 (BaseAPI + 各资源子类)
```

## 架构分层

```
HTTP Handler → Service → Repository → Ent Client → SQLite
```

**注意**: 部分 Handler 直接调用 Repository，跳过了 Service 层，分层不够一致。

## 开发命令

| 命令 | 说明 |
|------|------|
| `task go:run` | 启动后端 |
| `task ui:dev` | 启动前端 dev server |
| `task go:test` | 运行 Go 测试 (需要 CGO 环境) |
| `task ui:check` | 前端类型检查 |
| `task pr` | PR 前完整检查 |
| `task go:build` | 构建后端二进制 |

## 测试基础设施

- 后端测试使用 SQLite 内存数据库 (`file:ent?mode=memory&cache=shared&_fk=1`)
- 使用 `github.com/mattn/go-sqlite3`，需要 CGO 环境
- CGO 不可用时测试无法运行，但编译可通过 `CGO_ENABLED=0 go build ./...`
- 测试入口: `backend/internal/data/repo/main_test.go`

## 关键特性

- **Browse NAS Files**: 通过 symlink 导入服务器端文件
- **缩略图**: 优先 NAS 自带缩略图 (`.@__thumb/`)，否则本地缓存
- **WebSocket 事件推送**: 数据变更时通知前端
- **通知系统**: Shoutrrr 多平台 (Discord, Telegram, Slack 等)

## 已知约束

- Ent `*ItemQuery.Where()` 使用指针接收者直接修改共享 struct，`qb.Where(...)` 与 `qb = qb.Where(...)` 等价 — 不是 Bug。
- 代码注释和 UI 文本中英混合（如 "购买价格"、"制造商"）是有意设计，不是需要修复的问题。
