# TODO

后续升级计划：

## 前端
- [ ] **Nuxt 3.6.5 → 3.11+** — 修复 Node 22 的 `fs.Stats` deprecation warning（`DEP0180`），同时升级 Nitro。开发时无实质影响，仅警告。

## 后端
- [ ] **依赖批量升级** — `go get -u ./...` 刷新所有 Go 依赖（ent、chi、zerolog、crypto 等），跑完测试确认。

## CI/CD
- [x] **Node 版本统一** — CI 中 `partial-frontend.yaml` 的 `node-version` 已同步到 22（commit `26516b5`，2026-08-07；node 18 下 pnpm/action-setup@v6 崩溃）。

## DevContainer
- [ ] **`node_modules` volume 挂载** — 避免每次 rebuild 容器都要重装前端依赖。

## 发布
- [x] **验证 Docker 发版流程** — Docker Hub 发布已加入 CI/CD：`publish.yaml`（push main → nightly）与 `tag.yaml`（push tag `v*` → release 镜像 `latest` + `vX.Y.Z`）buildx 三平台自动构建推送，已实测 v0.1.21 全流程通过。手动 `scripts/release.sh` 仅应急使用（详见 README「Docker 镜像发布」）。rootless 镜像未做（当前只发标准镜像）。
