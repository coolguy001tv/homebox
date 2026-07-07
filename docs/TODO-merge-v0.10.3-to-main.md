# TODO: 合并 v0.10.3-release → main

> **状态**: 待执行  
> **创建日期**: 2026-07-06  
> **源分支**: `v0.10.3-release` (124 commits ahead of main)  
> **目标分支**: `main`

---

## 一、合并前准备

- [ ] **1.1 确认 CI 通过** — 确保 `v0.10.3-release` 分支的 CI/CD pipeline 全部绿色
- [ ] **1.2 确认 main 分支最新** — `git checkout main && git pull origin main`，确保基于最新 main
- [ ] **1.3 阅读变更摘要** — 关键变更：
  - HEIC 格式支持（大图查看、下载原图）
  - 缩略图方案
  - MCP SQLite 支持
  - Node 18 → 22、Go 版本升级
  - AssetID type 变更
  - 数据库 schema 变化（检查是否需要迁移）
  - 前端依赖升级（Nuxt、Vue 等）
- [ ] **1.4 备份数据库** — 确认 staging/本地数据库有备份，以防 schema 变更导致问题

## 二、执行合并

- [ ] **2.1 创建合并分支** — `git checkout -b merge/v0.10.3-to-main main`
- [ ] **2.2 执行合并** — `git merge v0.10.3-release`
- [ ] **2.3 手动解决冲突** — 重点关注：
  - `.devcontainer/Dockerfile` — Node 版本差异（18 vs 22）
  - `.devcontainer/devcontainer.json` — devcontainer 配置差异
  - `.devcontainer/devcontainer-lock.json` — 新增文件（go feature 配置）
  - 前端依赖文件（`pnpm-lock.yaml` 等）— 可能有大量依赖树冲突
  - Go 依赖文件（`go.mod`、`go.sum`）
- [ ] **2.4 验证冲突解决** — 逐个文件检查，确保两边改动都得到正确合并

## 三、合并后验证

- [ ] **3.1 后端编译** — `task build:backend` 或 `go build ./...` 必须通过
- [ ] **3.2 前端编译** — `task build:frontend` 或 `cd frontend && pnpm build` 必须通过
- [ ] **3.3 运行测试** — `task test`（Go test）+ `cd frontend && pnpm test`（Vitetest）
- [ ] **3.4 本地启动验证** — `task dev` 启动完整应用，手动检查：
  - [ ] 首页加载正常
  - [ ] 列表/详情页正常
  - [ ] HEIC 图片上传和查看
  - [ ] 缩略图生成正常
  - [ ] 图片比例 4/3 显示正确
  - [ ] 排序功能正常（默认创建时间倒序）
  - [ ] 标签筛选 AND 逻辑正确
  - [ ] SQLite MCP 功能正常
- [ ] **3.5 数据库迁移** — 确认 ent schema 变更自动迁移成功，无数据丢失

## 四、发布流程

- [ ] **4.1 推送合并分支** — `git push origin merge/v0.10.3-to-main`
- [ ] **4.2 创建 PR** — 在 GitHub 上创建 `merge/v0.10.3-to-main` → `main` 的 Pull Request
- [ ] **4.3 Code Review** — 至少一人 review 合并结果，重点关注冲突解决部分
- [ ] **4.4 合并 PR** — 使用 "Merge commit"（不用 squash，保留完整提交历史）
- [ ] **4.5 创建 Tag** — 在 main 上打 `v0.10.3` tag，推送到 GitHub
- [ ] **4.6 创建 GitHub Release** — 基于 tag 创建 release，编写 Release Notes，重点列出：
  - HEIC 支持
  - 缩略图方案
  - MCP SQLite
  - 版本升级（Node 22、Go）
  - Bug 修复汇总
- [ ] **4.7 Docker 镜像发布** — 构建并推送标准镜像 + rootless 镜像到 Docker Hub

## 五、发布后清理

- [ ] **5.1 验证 Docker 镜像** — `docker pull` 最新镜像，确认能正常启动
- [ ] **5.2 更新 v0.10.3-release 分支** — 将 main 合并回 `v0.10.3-release`，保持同步：
  ```bash
  git checkout v0.10.3-release
  git merge main
  git push origin v0.10.3-release
  ```
- [ ] **5.3 删除合并分支** — `git branch -d merge/v0.10.3-to-main` 本地和远程
- [ ] **5.4 通知相关方** — 在相关频道/issue 中通知发布完成

---

## 已知风险

| 风险 | 影响 | 应对 |
|------|------|------|
| 数据库 schema 变更 | 旧数据不兼容 | 合并前备份，确认 ent 自动迁移覆盖所有变更 |
| pnpm-lock.yaml 冲突 | 依赖解析失败 | 优先保留 `v0.10.3-release` 版本（更新），手动验证 |
| Node 22 兼容性 | 运行时错误 | CI 中已是 Node 22，本地验证充分即可 |
| HEIC 解码依赖 | 部分平台可能不支持 | 确认 CGO/libheif 依赖在 Docker 镜像中已安装 |

---

## 参考

- [Homebox 主仓库](https://github.com/coolguy001tv/homebox)
- 分支: `v0.10.3-release` → `main`
- 提交数: 124 commits
- 变更文件: 174 files, +6996 / -12131 lines
