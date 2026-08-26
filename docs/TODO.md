# TODO

按优先级整理的待办（源自 2026-07-14 代码审查，2026-08-26 从 AGENTS.md 迁入；历史细节见 git log）。

## P1 性能

- [ ] `ZeroOutTimeFields` 逐条 UPDATE（backend/internal/data/repo/repo_items.go）— N+1，应批量更新
- [ ] `SetPrimaryPhotos` 逐条 UPDATE（backend/internal/data/repo/repo_items.go）— N+1，应批量更新
- [ ] DaisyUI 31 个主题全量打包（frontend/tailwind.config.js）— 约 150KB CSS
- [ ] MarkdownIt 每次组件实例重新创建（frontend/components/global/Markdown.vue）

## P2 架构一致性

- [ ] 部分 Handler 绕过 Service 层直接调 Repo（backend/app/api/handlers/v1/v1_ctrl_items.go）
- [ ] 两套映射函数工具并存（backend/internal/data/repo/map_helpers.go vs automappers.go）— 功能重复

## P3 前端

- [ ] `data-contracts.ts` 全局 `@ts-nocheck` 禁用类型检查（frontend/lib/api/types/data-contracts.ts）
- [ ] `route()` 将 undefined/null 转为字符串（frontend/lib/api/base/urls.ts）
- [ ] `titlecase()`/`capitalize()` 空字符串边界问题（frontend/lib/strings/index.ts）
- [ ] 注册失败 loading 未重置（frontend/pages/index.vue:68-85）
- [ ] ItemCard 嵌套 `<a>` 标签（frontend/components/Item/Card.vue）— HTML 无效

## P4 依赖升级

- [ ] Nuxt 3.6.5 → 3.11+ — 修复 Node 22 的 `fs.Stats` deprecation warning（`DEP0180`），同时升级 Nitro。开发时无实质影响，仅警告。
- [ ] DaisyUI 2.x → 4.x/5.x
- [ ] TypeScript 5.0 → 5.x+
- [ ] 后端依赖批量升级 — `go get -u ./...` 刷新所有 Go 依赖（ent、chi、zerolog、crypto 等），跑完测试确认。

## 已完成

CI/CD、DevContainer、发布、Git 等已完成事项见归档记录 [`history/TODO-ci.md`](history/TODO-ci.md)。
