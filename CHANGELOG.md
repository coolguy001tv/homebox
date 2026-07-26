# Changelog

## [Unreleased]

### Fixed

- **NAS 缩略图本地缓存自动清理**：当 NAS 的 `.@__thumb/` 缩略图可用时，自动删除之前生成的本地缓存文件（`./.data/import-thumbs/`）。解决了以下场景的磁盘空间浪费问题：用户先于 NAS 定时任务访问图片 → Homebox 生成本地缓存 → NAS 后续生成缩略图 → 本地缓存永不再用却驻留磁盘。清理采用"访问时顺带清理"策略，仅在 Tier 1（NAS 缩略图）命中时执行一次 `os.Remove`，零额外开销。
