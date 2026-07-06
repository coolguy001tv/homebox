# TODO

后续升级计划：

## 前端
- [ ] **Nuxt 3.6.5 → 3.11+** — 修复 Node 22 的 `fs.Stats` deprecation warning（`DEP0180`），同时升级 Nitro。开发时无实质影响，仅警告。

## 后端
- [ ] **依赖批量升级** — `go get -u ./...` 刷新所有 Go 依赖（ent、chi、zerolog、crypto 等），跑完测试确认。

## CI/CD
- [ ] **Node 版本统一** — CI 中 `partial-frontend.yaml` 的 `node-version: 18` 同步到 22。

## DevContainer
- [ ] **`node_modules` volume 挂载** — 避免每次 rebuild 容器都要重装前端依赖。
