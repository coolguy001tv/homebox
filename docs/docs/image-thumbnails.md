# Image Thumbnails

Homebox 内置了图片自动缩略图功能，可以显著减少图片附件的传输流量。本页介绍该功能的使用方法和注意事项。

## 概述

默认情况下，Homebox 通过 `/api/v1/items/{id}/attachments/{attachment_id}` 直接返回用户上传的原始图片文件。对于手机拍摄的高分辨率照片（通常 4000×3000 像素、4-8MB），在列表卡片等小尺寸展示场景中会产生大量不必要的带宽消耗。

缩略图功能通过以下方式解决这个问题：

- **按需生成**：首次请求时动态生成缩略图并写入磁盘缓存
- **磁盘缓存**：后续相同尺寸的请求直接返回缓存文件，零 CPU 开销
- **原图保留**：原始文件不受影响，下载和弹窗查看仍使用原图

## API 路由

### 缩略图接口

```
GET /api/v1/items/{id}/attachments/{attachment_id}/thumb?w={width}
```

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `id` | path | 是 | 物品 UUID |
| `attachment_id` | path | 是 | 附件 UUID |
| `w` | query | 否 | 目标宽度（像素）。省略时返回原图。 |

### 鉴权

缩略图接口使用与原附件接口相同的鉴权中间件（`assetMW`），支持：

- **Bearer Token**：标准 `Authorization: Bearer <token>` 请求头
- **access_token**：URL 查询参数 `?access_token=<token>`，适用于 `<img>` 标签直接渲染

### 行为说明

| 场景 | 行为 |
|------|------|
| `?w=400` | 缩放至宽度 400px，高度等比缩放，输出 JPEG (quality 85) |
| `?w=0` 或不传 `w` | 返回原始文件 |
| 原图宽度 ≤ 目标宽度 | 不缩放，仅重新编码为 JPEG |
| 缓存文件已存在 | 直接返回缓存，跳过编解码 |
| 非图片附件 | 返回原始文件（不处理） |

### 示例

```bash
# 获取 400px 宽度的缩略图
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:7745/api/v1/items/872f387e/attachments/751e69f1/thumb?w=400"

# 配合 access_token 在 img 标签中使用
<img src="/api/v1/items/{itemId}/attachments/{attachmentId}/thumb?w=400&access_token={token}" />
```

## 缓存机制

### 缓存位置

缩略图缓存文件与原图存放在同一目录下，文件名格式为：

```
<原文件名>.thumb_w<宽度>.jpg
```

例如：

```
# 原图
./.data/<gid>/documents/a1b2c3d4.jpg

# 400px 缩略图缓存
./.data/<gid>/documents/a1b2c3d4.thumb_w400.jpg

# 800px 缩略图缓存
./.data/<gid>/documents/a1b2c3d4.thumb_w800.jpg
```

### 缓存管理

- 缓存文件**不会自动清理**。如需清理，手动删除 `*.thumb_w*.jpg` 文件即可，下次请求时会自动重新生成。
- 删除原图附件时，关联的缩略图缓存**不会自动删除**（可定期清理）。
- 缓存文件的总大小远小于原图，通常无需关注。

### 生成失败

如果缩略图生成失败（例如源文件损坏），服务会记录错误日志并**回退到返回原图**，不会中断请求。

## 前端集成

前端已在以下场景使用缩略图：

| 页面 | 缩略图宽度 | 说明 |
|------|------------|------|
| 物品列表卡片 (`Card.vue`) | 400px | 卡片首图，`aspect-[4/3]` 容器 |
| 表格视图 (`Table.vue`) | 400px | 表格首图列 |
| 物品详情页 (`index.vue`) | 800px | 图片列表区域，`max-h-[200px]` 内展示 |
| 编辑页 (`edit.vue`) | 100px | 附件列表缩略图，高度 48px |
| 详情页弹窗 | 原图 | 点击图片放大查看时使用原始分辨率 |

### 开发者用法

在代码中使用 `BaseAPI.thumbURL()` 生成缩略图 URL：

```typescript
const api = useUserApi();

// 获取 400px 宽度的缩略图 URL（自动附加 access_token）
const thumbUrl = api.thumbURL(itemId, attachmentId, 400);

// 获取原始图片 URL（下载等场景）
const originalUrl = api.authURL(`/items/${itemId}/attachments/${attachmentId}`);
```

## 支持的图片格式

以下格式支持缩略图生成：

- JPEG (`.jpg`, `.jpeg`)
- PNG (`.png`)
- WebP (`.webp`)
- GIF (`.gif`，仅第一帧)
- BMP (`.bmp`)
- TIFF (`.tiff`)

所有格式的缩略图统一输出为 **JPEG** 格式，质量 **85%**。GIF 动画仅保留第一帧。

## 技术参数

| 参数 | 值 |
|------|-----|
| 输出格式 | JPEG |
| 输出质量 | 85（1-100） |
| 缩放算法 | Lanczos |
| 自动旋转 | 是（读取 EXIF Orientation） |
| 依赖库 | `github.com/disintegration/imaging` |

## 流量对比示例

以一张 4000×3000 像素、4.5MB 的典型手机照片为例：

| 展示场景 | 无缩略图 | 有缩略图 | 节省 |
|----------|----------|----------|------|
| 列表卡片 (400px) | 4.5 MB | ~30 KB | 99.3% |
| 详情页 (800px) | 4.5 MB | ~100 KB | 97.8% |
| 编辑页预览 (100px) | 4.5 MB | ~5 KB | 99.9% |
| 弹窗查看原图 | 4.5 MB | 4.5 MB | 0%（有意保留） |

## 注意事项

1. **首次请求延迟**：首次请求某个尺寸的缩略图时需要解码、缩放、编码，会对该请求增加数百毫秒延迟。后续请求不受影响。
2. **磁盘空间**：每个缩略图缓存文件通常为 10-100KB，对于百张图片的规模影响可忽略。如空间敏感，可定期清理 `*.thumb_w*.jpg` 文件。
3. **升级兼容**：缩略图路由为新增接口，不影响已有的原图路由和 API。现有前端代码无缝升级。
4. **GIF 动画**：缩略图仅保留 GIF 的第一帧，如需查看动画请使用原图路由。
5. **透明通道**：PNG/WebP 的透明背景在缩略图中会变为白色（JPEG 不支持透明）。

## 配置

缩略图功能暂无可配置项，使用以下硬编码值：

- JPEG 质量：`85`
- 缩放算法：`Lanczos`

如需调整，可修改 `backend/internal/core/services/service_items_thumbnails.go` 中的 `ThumbnailQuality` 常量。
