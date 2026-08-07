# HomeBox 项目说明

## 项目概况

- **技术栈**: Go + Ent ORM + SQLite + Chi Router (后端) + Nuxt 3.6.5 + Vue 3.5 + DaisyUI 2 + TypeScript (前端)
- **部署**: Docker, 前端 Nuxt 构建产物嵌入 Go binary
- **当前分支**: `main`（历史 `v0.10.3-release` 分支已冻结，不再使用）
- **Docker Hub**: `hellocoolguy/homebox:0.1.x`
- **镜像发布**: 由 GitHub Actions 自动构建（push 到 main → nightly 三平台镜像；push tag `v*` → 跑完整测试后构建 release 镜像 `latest`+`vX.Y.Z` 并部署 docs）。**推荐发布方式**：`git tag vX.Y.Z && git push origin vX.Y.Z`（详见 README「Docker 镜像发布」）；不要手动构建——`scripts/release.sh` 仅应急/本地验证用。

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

---

## 代码审查结论 (2026-07-14)

以下结论记录了上次全面代码审查的结果，方便后续优化时优先查看。

### Ent QueryBuilder 指针语义

`*ItemQuery.Where()` 使用指针接收者 (`func (iq *ItemQuery) Where(...) *ItemQuery`)。调用 `iq.Where()` 时通过指针直接修改了共享 struct (`iq.predicates = append(...)`)，所以 `qb.Where(...)` 和 `qb = qb.Where(...)` 功能等价 — 不是 Bug。

### 已修复

- **AttachmentRepo.Update 事务缺失** (`backend/internal/data/repo/repo_item_attachments.go:90-150`): 用 `ent.Tx()` 包裹所有写操作，保证「更新主照片 → 清除其他 primary 标记」的原子性。修复后未跑测试（Windows 本机缺 CGO），建议在 Linux/Mac 或 CI 中验证。

### 待处理 (按优先级)

**P1 - 代码质量**
- `ItemService.filepath` 字段 (service_items.go:23) 从未使用 — 死代码
- `HandleSetPrimaryPhotos` 描述字符串复制粘贴错误 (v1_ctrl_actions.go:82): 描述为 "ensure asset IDs" 应为 "set primary photos"
- `extractQuery` 每次请求重新定义 (v1_ctrl_items.go:35-80) — 应提取为包级别函数

**P2 - 性能**
- `ZeroOutTimeFields` 逐条 UPDATE (repo_items.go:843-888) — N+1，应批量更新
- `SetPrimaryPhotos` 逐条 UPDATE (repo_items.go:890-937) — 同上
- DaisyUI 31 个主题全量打包 (tailwind.config.js) — ~150KB CSS
- MarkdownIt 每次组件实例重新创建 (components/global/Markdown.vue)

**P3 - 架构一致性**
- 部分 Handler 绕过 Service 层直接调 Repo (v1_ctrl_items.go)
- 两套映射函数工具并存 (map_helpers.go vs automappers.go) — 功能重复

**P4 - 前端**
- `fmtCurrency()` 硬编码 CNY，忽略传入货币参数
- `data-contracts.ts` 全局 `@ts-nocheck` 禁用类型检查
- `route()` 将 undefined/null 转为字符串 (lib/api/base/urls.ts)
- `titlecase()`/`capitalize()` 空字符串边界问题 (lib/strings/index.ts)
- 注册失败 loading 未重置 (pages/index.vue:68-85)
- ItemCard 嵌套 `<a>` 标签 (components/Item/Card.vue) — HTML 无效

**P5 - 依赖更新**
- Nuxt 3.6.5 → 最新版
- DaisyUI 2.x → 4.x/5.x
- TypeScript 5.0 → 5.x+

### 已知特点 (非问题)

- 项目中存在中英混合的代码注释和 UI 文本（如 "购买价格"、"制造商"），这是有意设计的本地化，不是需要修复的问题。
- 前端多处硬编码中文 UI 文本，未使用 i18n 方案。
