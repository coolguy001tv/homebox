<div align="center">
  <img src="/docs/docs/assets/img/lilbox.svg" height="200"/>
</div>

<h1 align="center" style="margin-top: -10px"> HomeBox </h1>
<p align="center">
   <a href="https://hay-kot.github.io/homebox/">Docs</a>
   |
   <a href="https://homebox.fly.dev">Demo</a>
   |
   <a href="https://discord.gg/tuncmNrE4z">Discord</a>
</p>

## Quick Start

## Docker 镜像发布

镜像由 **GitHub Actions 自动构建**，正常情况下无需手动构建/推送。

### GitHub Actions 何时自动构建

| 事件 | 触发的 workflow | 产物 |
|---|---|---|
| push 到 `main` | `publish.yaml` | 三平台 nightly 镜像 `hellocoolguy/homebox:nightly` |
| push tag `v*`（如 `v0.1.21`） | `tag.yaml` | 跑完整测试（后端 + 前端）→ 通过后构建三平台 release 镜像 `hellocoolguy/homebox:latest` + `v0.1.21` → 部署 docs |
| PR 到 `main` | `pull-requests.yaml` | 只跑测试，不构建镜像 |

三平台 = `linux/amd64`、`linux/arm64`、`linux/arm/v7`（buildx + QEMU 交叉构建）。

### 推荐发布流程（用 GitHub Actions 构建）

```bash
# 1. 代码已合并到 main 后，打版本 tag 并推送：
git tag v0.1.21
git push origin v0.1.21
# 2. GitHub Actions 自动完成：测试 → 三平台构建 → 推送 Docker Hub（latest + v0.1.21）→ 部署 docs
# 3. 在仓库 Actions 页面查看 Publish Release 运行进度
```

版本号规则：`vX.Y.Z`。线上最新版本可参考 [Docker Hub](https://hub.docker.com/repository/docker/hellocoolguy/homebox/general)。

### 手动构建（不推荐，仅应急/本地验证）

仓库提供 `scripts/release.sh` 一键脚本（构建 → 推送 → 可选 latest / git tag）：

```bash
./scripts/release.sh 0.1.21 -y --latest --git-tag
```

或手动执行：

```bash
docker build . -t hellocoolguy/homebox:0.1.21   # 记得先打开本机 Docker Desktop，否则找不到 docker 命令
docker push hellocoolguy/homebox:0.1.21
```

> **为什么推荐 GitHub Actions 构建而不是手动？**
> - 手动 `docker build` 是**单平台**（仅本机 CPU 架构），CI 是**三平台**（amd64/arm64/arm/v7）
> - CI 在干净环境构建（固定 Go 1.25 / Node 22 工具链），不依赖本机环境、可复现；手动构建受本机工具链影响，产物可能与线上不一致
> - 手动发版容易漏推 `latest` 或 git tag；CI 自动带完整测试保障
> - **除非 CI 不可用，否则不要手动发版**。若只是想触发 CI 构建，直接 `git tag` + `git push origin` 即可，无需跑 `release.sh`

如果本地构建遇到网络问题，可先设置代理再构建：

1. `export http_proxy=http://127.0.0.1:7890`
2. 测试网站访问： `curl -I https://google.com`
3. 后续可以在 proxy(Clash) 中查看连接是否正常。

`docker images` 可查看本机已构建的镜像（例如上面的 `hellocoolguy/homebox:0.1.21`）。

## 分支注意事项
目前是基于v0.10.3-release分支进行的开发，不要切换到main上，名字的由来是当前分支是官方的v0.10.3版本切出来的。

## 一些异常情况

### `ERROR: failed to solve: node:22-alpine: failed to resolve source metadata for docker.io/library/node:22-alpine: ...`

直接手动拉一下： `docker pull node:22-alpine`

> 注意：这里以 Dockerfile 中实际使用的基础镜像为准，当前为 `node:22-alpine`。

### ` failed to solve: archive/tar: unknown file mode ?rwxr-xr-x`

如果遇到此问题，请检查 Docker Desktop 是否已更新到最新版本，并确保文件系统权限正常。

## 本地开发注意事项

1. 执行`task go:run`和`task ui:dev`
2. 数据库homebox.db拷贝到`backend/.data`下，对应的资源文件目录（6816a....）也放置在同目录下
3. 由于docker中的路径问题，需要修改数据库，在dataGrip或其他数据库APP中修改数据：

```sh
UPDATE documents
SET path= replace(path, '/data/','.data\')
WHERE  path like '%/data/%';

```
4. 打开前端页面`http://localhost:3000`

[Configuration & Docker Compose](https://hay-kot.github.io/homebox/quick-start)

```bash {"id":"01J05HHH061CS0NPR62CBA4E9K"}
# If using the rootless image, ensure data 
# folder has correct permissions
mkdir -p /path/to/data/folder
chown 65532:65532 -R /path/to/data/folder
docker run -d \
  --name homebox \
  --restart unless-stopped \
  --publish 3100:7745 \
  --env TZ=Europe/Bucharest \
  --volume /path/to/data/folder/:/data \
  ghcr.io/hay-kot/homebox:latest
# ghcr.io/hay-kot/homebox:latest-rootless


```

## Credits

- Logo by [@lakotelman](https://github.com/lakotelman)

## Browse NAS Files 功能说明

### 配置方法

通过环境变量 `HBOX_IMPORT_DIRS` 指定可在编辑页中浏览的服务器端目录，Docker 挂载时建议使用 `:ro` 只读模式。

示例（挂载多个目录）：

```yaml
environment:
  - HBOX_IMPORT_DIRS=/import/Pictures,/import/Docs

volumes:
  - /share/CoolGuy/Pictures:/import/Pictures:ro
  - /share/CoolGuy/Docs:/import/Docs:ro
```

### 工作原理

导入文件时在 `/data/<gid>/documents/` 下创建**符号链接（symlink）**指向原始文件，不复制文件内容。删除附件时仅删除 symlink，原始文件不受影响。

### 更换映射目录后保持已有 symlink 有效

如果将来需要更换 NAS 上的源目录（例如从 `/share/CoolGuy/Pictures` 迁移到 `/share/CoolGuy/Pictures2`），已创建的大量 symlink 仍指向旧的容器内路径 `/import/Pictures/xxx.jpg`，直接去掉旧映射会导致这些链接失效。

**解决方案：docker-compose 中同时挂载新旧目录，旧目录保持只读，新目录用于后续使用：**

```yaml
environment:
  - HBOX_IMPORT_DIRS=/import/Pictures2,/import/Docs2   # Browse NAS Files 只显示新目录

volumes:
  - /share/Public/App/homebox:/data/
  - /share/CoolGuy/Pictures:/import/Pictures:ro        # 旧目录保留，已有 symlink 继续有效
  - /share/CoolGuy/Docs:/import/Docs:ro                # 旧目录保留
  - /share/CoolGuy/Pictures2:/import/Pictures2:ro      # 新目录，供 Browse NAS Files 使用
  - /share/CoolGuy/Docs2:/import/Docs2:ro              # 新目录
```

**关键点**：
- `HBOX_IMPORT_DIRS` 只写新目录路径，Browse NAS Files 弹窗只显示新目录
- 旧目录映射**不写在** `HBOX_IMPORT_DIRS` 中（不在文件浏览器中显示），但**保留 volumes 挂载**以维持已有 symlink 有效
- 旧映射路径保持不变（如 `/import/Pictures`），否则已有 symlink 会断裂
