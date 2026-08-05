# CI 待办清单（2026-08-05 创建，明天继续）

## ✅ 已完成

- **修复前端 CI install 失败**：工作流引用改为本地 `partial-*.yaml`（不再跟随上游 `hay-kot/homebox@main` 的 pnpm 6 配置），`pnpm install` 加 `--frozen-lockfile` 严格锁定 lockfile 版本
- **发布目标改为 Docker Hub** `hellocoolguy/homebox`（buildx 三平台 amd64/arm64/arm/v7，移除上游 ghcr.io/hay-kot）
- **移除 goreleaser**（不需要 GitHub 二进制发布）与 **Fly.io deploy**（无 FLY_API_TOKEN）
- **打 tag `v0.1.19` 并推送**，触发发布
- **`.eslintrc.js` → `.eslintrc.cjs`**：修复 `frontend/package.json` 里 `"type": "module"` 导致 ESLint 8 无法加载 CommonJS 配置的问题

## 🔍 待确认

### 1. 确认 v0.1.19 镜像发布结果
- 查看 GitHub Actions **Publish Release** 运行（run id `31012908020`）中 `Publish Tag / Publish Homebox` job 是否成功（当时已登录 Docker Hub 成功，正在构建推送 `latest` + `v0.1.19`）
- 确认 Docker Hub 上 `hellocoolguy/homebox:latest` 与 `hellocoolguy/homebox:v0.1.19` 已存在（三平台）
- 确认 **Publish Dockers** 运行（run id `31012673571`，push main 触发的 `nightly` 镜像）结果
- 链接：https://github.com/coolguy001tv/homebox/actions

## ⚠️ 待修复（CI 全绿）

### 2. 前端 lint 失败（`pnpm run lint:ci`，10 个 error）
`.eslintrc.cjs` 修复后 config 可加载，但暴露 10 个既有代码 error，分布在 6 个文件：

| 文件 | 位置 | 问题 |
|---|---|---|
| `components/Item/FileBrowser.vue` | :112 | `toast` 未使用 |
| `composables/utils.ts` | :28 | `currency` 未使用 |
| `lib/api/classes/items.ts` | :17 | import 顺序（`./import` 应在 `~~/lib/requests` 前） |
| `pages/item/[id]/index.vue` | :261 / :268 / :359 / :366 | `showWarranty` / `warrantyDetails` / `showSold` / `soldDetails` 未使用 |
| `pages/item/[id]/index/edit.vue` | :516 | `<template v-for>` 不能有 key |
| `pages/items.vue` | :57 / :68 | `orderByLabel` / `orderLabel` 未使用 |

> 注意：lint 输出里大量 `Delete ␍`（CRLF）警告**只是本地 Windows 现象**。git 中文件以 LF 存储（`attr/text eol=lf`），CI（Linux）上不会出现，**不用处理**。

### 3. 后端 golangci-lint 失败
- `Backend Server Tests / Go` 的 `golangci-lint` 步骤失败，需本地跑 golangci-lint 或 `task` 查看具体报错（可能是既有 Go 代码问题）

### 4. 修复后重新发版
- 修复上述 lint 后 commit + 推 main，下次打 tag 时 CI 即可全绿
- 如果希望 `v0.1.19` 的 CI 记录也是绿的，可在修复后删除远端 tag 并重新打（镜像会重新推送，覆盖同名标签）
