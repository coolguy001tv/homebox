# TODO

后续升级计划：

## 前端
- [ ] **Nuxt 3.6.5 → 3.11+** — 修复 Node 22 的 `fs.Stats` deprecation warning（`DEP0180`），同时升级 Nitro。开发时无实质影响，仅警告。

## 后端
- [ ] **依赖批量升级** — `go get -u ./...` 刷新所有 Go 依赖（ent、chi、zerolog、crypto 等），跑完测试确认。

## 已完成

CI/CD、DevContainer、发布、Git 等已完成事项见归档记录 [`history/TODO-ci.md`](history/TODO-ci.md)。
