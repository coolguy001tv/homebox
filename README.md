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

## Docker build脚本
可以参考[这里](https://hub.docker.com/repository/docker/hellocoolguy/homebox/general)看看线上最新版本的信息

注意：每次记得修改以下的版本号信息，当前使用的是0.1.11 （可以当前页面全局替换0.1.11到0.1.X）

`docker build . -t hellocoolguy/homebox:0.1.11`

记得先打开本地的Docker Desktop，否则找不到docker命令

如果遇到网络问题需要设置代理，可以在终端中设置：

1. `export http_proxy=http://127.0.0.1:7890`
2. 测试网站访问： `curl -I https://google.com`
3. 后续可以在proxy(Clash)中查看连接看看是否正常。

## 发布流程

1. 确保上面的docker build命令正常成功执行
2. 登录到docker hub: `docker login` (密码参见notion文档)
3. `docker push hellocoolguy/homebox:0.1.11`
4. 可选的 latest tag操作：
   1. docker tag hellocoolguy/homebox:0.1.11 hellocoolguy/homebox:latest
   2. docker push hellocoolguy/homebox:latest

## 一些有用的脚本

1. `docker images`可以查看当前已经build好的image，比如前面通过`docker build`出来的`hellocoolguy/homebox:0.1.11`

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
