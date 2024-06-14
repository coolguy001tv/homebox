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
`docker build . -t hellocoolguy/homebox:0.0.3`
记得在wsl(ubuntu)下运行，否则有权限问题

目前已经设置了wsl的代理，如果后续有问题，可以重新设置一下，设置方式如下：
1. `vi ~/.bashrc` 添加行`export http_proxy=http://127.0.0.1:7890`
2. 执行`source ~/.bashrc`（或相应的配置文件）使更改生效
3. 测试网站访问： `curl -I https://google.com`
4. 后续可以在proxy(Clash)中查看连接看看是否正常。

## 发布流程
1. 确保上面的docker build命令正常成功执行
2. 登录到docker hub: `docker login` (密码参见notion文档)
3. `docker push hellocoolguy/homebox:0.0.3`
4. 可选的 latest tag操作：
    1. docker tag hellocoolguy/homebox:0.0.3 hellocoolguy/homebox:latest
    2. docker push hellocoolguy/homebox:latest

## 一些有用的脚本
1. `docker images`可以查看当前已经build好的image，比如前面通过`docker build`出来的·hellocoolguy/homebox:0.0.3`

## 一些异常情况
### `ERROR: failed to solve: node:18-alpine: failed to resolve source metadata for docker.io/library/node:18-alpine: ...: 8134 != 7826: failed precondition`
直接手动拉一下： `docker pull node:18-alpine`
### ` failed to solve: archive/tar: unknown file mode ?rwxr-xr-x`
win上没有这个权限，请在wsl(ubuntu)下运行

## 本地开发注意事项
1. 执行`task go:run`和`task ui:dev`
2. 数据库homebox.db拷贝到`backend/.data`下，对应的资源文件目录（6816a....）也放置在同目录下
3. 由于docker中的路径问题，需要修改数据库，在dataGrip或其他数据库APP中修改数据：

```
UPDATE documents
SET path= replace(path, '/data/','.data\')
WHERE  path like '%/data/%';
```



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
