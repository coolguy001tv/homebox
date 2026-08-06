# CI 待办清单（2026-08-06 更新）

## ✅ 已完成

- **修复前端 CI install 失败**：工作流引用改为本地 `partial-*.yaml`（不再跟随上游 `hay-kot/homebox@main` 的 pnpm 6 配置），`pnpm install` 加 `--frozen-lockfile` 严格锁定 lockfile 版本
- **发布目标改为 Docker Hub** `hellocoolguy/homebox`（buildx 三平台 amd64/arm64/arm/v7，移除上游 ghcr.io/hay-kot）
- **移除 goreleaser**（不需要 GitHub 二进制发布）与 **Fly.io deploy**（无 FLY_API_TOKEN）
- **打 tag `v0.1.19` 并推送**，触发发布
- **`.eslintrc.js` → `.eslintrc.cjs`**：修复 `frontend/package.json` 里 `"type": "module"` 导致 ESLint 8 无法加载 CommonJS 配置的问题

### 2026-08-06 修复批次（本地全部验证通过）

| 项 | 修复内容 |
|---|---|
| 前端 lint（10 error） | 修复 6 个文件：删除未使用变量（FileBrowser/toast、index.vue 4 个死代码 computed、items.vue 2 个）、import 顺序、`<template v-for>` key 移到真实元素 |
| 前端 typecheck | 修 5 个既有 TS 错误：`ImportAPI.thumbURL` 重命名 `fileThumbURL`（避免与基类签名冲突）、index.vue `item.value` 非空断言、测试文件补 `LabelCreate.required`、Table.vue 比较断言；tsconfig 加 `ignoreDeprecations`（nuxi typecheck 用 npx 拉 latest TS，对 Nuxt 3.6 生成的 baseUrl/moduleResolution 报弃用错误） |
| 后端 golangci-lint | 修 errcheck（tx.Rollback/f.Close）、gocritic ifElseChain/wrapperFunc、gosimple S1009；`.golangci.yml` 排除 opinionated 检查（captLocal 全大写缩写参数名、exitAfterDefer TestMain/main 退出模式）、移除已弃用 `exportloopref`；CI 固定 `version: v1.64.8`（`latest`=v2 无法加载 v1 配置，CI 报 exit code 3） |
| Integration Tests（vitest 无法启动） | `postcss.config.js` / `tailwind.config.js` 是 CommonJS 语法但 package.json 是 `"type": "module"` → 改为 ESM 语法（`export default` / `import`） |
| notifier 测试 5s 超时 flake | 并发跑完整套件时 singleUse（注册+登录）超过默认 5s，用例加 `15000` 超时 |
| 发布 build 卡死防护 | `partial-publish.yaml` 的 build nightly/release image 步骤加 `timeout-minutes: 30`（v0.1.19 的 release 构建曾卡近 6h 被 GH Actions 超时取消） |

## 🔍 v0.1.19 发布失败原因（已查明）

- **3 个测试 job 失败**：Frontend Lint（10 error）、Frontend Integration Tests（postcss config 加载失败）、Backend golangci-lint（config exit 3）
- **Publish Homebox 被取消**：`build release image` 从 14:00 跑到 19:59（近 6 小时）被 GitHub Actions 6h 超时强制取消 → Docker Hub 上**没有** v0.1.19，`latest` 仍是 0.1.17（7/27）
- 同日 Publish Nightly 三平台构建 28 分钟成功，说明 release 那次是偶发卡死/限流

## ⏭️ 下一步

1. 提交全部修复并推送 main
2. 删除远端 tag `v0.1.19` 并重新打 tag（镜像重新推送覆盖同名标签）
3. 确认新 run 全绿 + Docker Hub 出现 `hellocoolguy/homebox:v0.1.19` 和更新后的 `latest`
