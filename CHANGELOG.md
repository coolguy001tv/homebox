# 版本更新日志

## v0.1.13 (2026-07-13)

### 🐛 问题修复

- **Browse NAS Files 子目录导航路径错误**：修复从父目录点击子文件夹时面包屑导航路径丢失中间层级的问题（如 `Home/Q1/Q3` 变成 `Home/Q3`）。改为基于当前路径累积拼接，不再依赖不稳定的绝对路径前缀剥离方式。

### 🛠️ 开发工具

- **Git 标签发布脚本**：`scripts/release.sh` 新增 Git tag 推送步骤（Docker 发布后询问是否推送版本标签）。新增独立 `scripts/git-tag.sh`，支持直接执行 `./scripts/git-tag.sh [版本号]` 单独推送标签，含重复标签检测、工作区状态检查。

---

## v0.1.12 (2026-07-10)

### ⚡ 性能优化：Browse NAS Files 分页与无限滚动

- **后端分页**：`/api/v1/import/browse` 新增 `page`/`pageSize` 查询参数，文件按页返回（默认每页 50 个），目录始终全部返回（导航元素）
- **前端无限滚动**：FileBrowser 改为 IntersectionObserver 哨兵机制，滚动到底部自动加载下一页，不再一次性渲染全部文件
- **缩略图加载优化**：每页仅渲染 50 个文件条目，配合 `loading="lazy"` + `decoding="async"`，大幅减少并发的缩略图生成请求，避免 NAS I/O 被占满

### 🔧 API 变化

`GET /api/v1/import/browse` 响应格式从 `FileEntry[]` 变更为：

```json
{
  "dirs": [...],     // 所有目录（始终全部，已排序）
  "files": [...],    // 当前页文件（已排序）
  "page": 1,
  "pageSize": 50,
  "total": 200       // 文件总数
}
```

### 🖼️ 缩略图缓存优化：避免污染 NAS 目录

- **优先使用 NAS 缩略图**：`ThumbnailPath` 优先读取 NAS 自带缩略图（`.@__thumb/<imageName>`），画质更好、无需额外 I/O
- **缓存移至本地**：无法使用 NAS 缩略图时，不再将缓存写入 NAS 源目录，改为缓存到 Homebox 数据目录 `<dataDir>/import-thumbs/`，消除 NAS 自身对缓存文件重复生成缩略图的浪费
- **可配置缩略图目录**：新增 `HBOX_IMPORT_THUMB_DIR` 环境变量，可自定义 NAS 缩略图子目录名（默认 `.@__thumb`）

---

## v0.1.11 (2026-07-10)

### 🖼️ 新增功能：Browse NAS Files（NAS 文件直接导入）

- **服务端文件浏览**：编辑物品页新增 "Browse NAS Files" 按钮，可直接浏览 NAS 上的文件目录，无需从 PC 上传
- **Symlink 导入**：导入文件时使用符号链接（symlink），不复制文件内容，**零额外磁盘空间占用**
- **缩略图预览**：文件浏览器中图片自动显示缩略图，支持横向图片友好展示
- **目录导航**：支持面包屑导航浏览多层子目录
- **安全控制**：
  - 通过环境变量 `HBOX_IMPORT_DIRS` 白名单控制可浏览的目录
  - 路径穿越防护（`../../../etc/passwd` 等攻击会被拦截）
  - 建议 Docker 挂载时使用 `:ro` 只读模式

### 🐛 问题修复

- **README 过时内容更新**：更新版本号至 0.1.11，修复 Docker build 命令和故障排查部分的镜像引用
- **前端组件导入**：显式 import FileBrowser 组件，确保 Nuxt dev server 正确识别

### 🐳 Docker 部署配置

部署时需在 `docker-compose.yml` 中新增以下配置：

```yaml
environment:
  - HBOX_IMPORT_DIRS=/import/Pictures   # 允许浏览的目录（逗号分隔多个）

volumes:
  - /your/nas/pictures:/import/Pictures:ro   # :ro = 只读，防止误修改原始文件
```

如需挂载多个目录，参考以下配置：

```yaml
environment:
  - HBOX_IMPORT_DIRS=/import/Pictures,/import/Docs

volumes:
  - /share/CoolGuy/Pictures:/import/Pictures:ro
  - /share/CoolGuy/Docs:/import/Docs:ro
```

> **📌 更换映射目录后如何保持已有 symlink 有效？**
>
> 如果将来需要更换源目录（例如从 `/share/CoolGuy/Pictures` 迁移到 `/share/CoolGuy/Pictures2`），直接在 docker-compose 中**同时挂载新旧目录**即可：
>
> ```yaml
> environment:
>   - HBOX_IMPORT_DIRS=/import/Pictures2          # 只写新目录，Browse NAS Files 不显示旧目录
>
> volumes:
>   - /share/CoolGuy/Pictures:/import/Pictures:ro   # 旧目录保留，已有 symlink 继续有效
>   - /share/CoolGuy/Pictures2:/import/Pictures2:ro # 新目录供 Browse NAS Files 使用
> ```
>
> 关键点：`HBOX_IMPORT_DIRS` 只写新路径，旧路径仅保留 volumes 挂载（不在文件浏览器中出现）。旧路径如变动会导致已有 symlink 断裂。

---

## v0.1.10 (2026-07-07)

### 🖼️ 图片功能

- **新增 HEIC 格式支持**：支持上传和查看 HEIC 格式的图片，包括大图查看和下载原图功能
- **缩略图方案**：新增缩略图生成方案，优化图片列表的加载性能
- **优化缩略图显示**：优化缩略图在不同场景（列表、详情等）下的显示效果

### 🔧 基础设施

- **更新 Node、Go 等依赖版本**：升级项目依赖的 Node.js 和 Go 版本
- **支持新版 Go Taskfile**：适配新版 Taskfile 格式，确保构建任务正常运行
- **AssetID 类型变更**：统一 AssetID 的数据类型

### 🛠️ 开发工具

- **新增 SQLite MCP**：集成 SQLite MCP 工具，方便开发时直接查询本地数据库
- **新增一键发布脚本**：`scripts/release.sh`，自动化 Docker 构建、登录检查、推送及 latest 标签同步，版本号需交互确认

---

## v0.1.9

初始版本（基于官方 v0.10.3 分支）
