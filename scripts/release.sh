#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# HomeBox Docker Release Script
# 一键发布脚本：构建 → 检查登录 → 推送 → (可选) latest 标签
# ============================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info()  { printf "${GREEN}[INFO]${NC}  %s\n" "$*"; }
log_warn()  { printf "${YELLOW}[WARN]${NC}  %s\n" "$*"; }
log_error() { printf "${RED}[ERROR]${NC} %s\n" "$*"; }
log_step()  { printf "\n${CYAN}==> %s${NC}\n\n" "$*"; }

# -----------------------------------------------------------
# 1. 解析/确认版本号
# -----------------------------------------------------------
VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    # 交互式询问版本号
    read -r -p "请输入发布版本号 (例如 0.1.11): " VERSION
fi

if [[ -z "$VERSION" ]]; then
    log_error "版本号不能为空，已取消发布。"
    exit 1
fi

# 二次确认
echo ""
read -r -p "确认发布版本 ${VERSION}？(y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    log_info "已取消发布。"
    exit 0
fi

IMAGE="hellocoolguy/homebox:${VERSION}"
LATEST="hellocoolguy/homebox:latest"

# -----------------------------------------------------------
# 2. 检查 Docker 环境
# -----------------------------------------------------------
log_step "检查 Docker 环境..."

if ! command -v docker &> /dev/null; then
    log_error "未找到 docker 命令。请先安装 Docker Desktop 或 Docker Engine。"
    exit 1
fi

if ! docker info &> /dev/null; then
    log_error "Docker 未运行。请先启动 Docker Desktop。"
    exit 1
fi

log_info "Docker 环境正常。"

# -----------------------------------------------------------
# 3. 检查 Docker Hub 登录状态
# -----------------------------------------------------------
log_step "检查 Docker Hub 登录状态..."

DOCKER_CONFIG="${DOCKER_CONFIG:-$HOME/.docker/config.json}"
DOCKER_REGISTRY="https://index.docker.io/v1/"

if [[ -f "$DOCKER_CONFIG" ]]; then
    if grep -q "\"${DOCKER_REGISTRY}\"" "$DOCKER_CONFIG" 2>/dev/null; then
        log_info "已登录 Docker Hub。"
    else
        log_error "未登录 Docker Hub，请先执行: docker login"
        log_error "（Docker Hub 密码参见 Notion 文档）"
        exit 1
    fi
else
    log_error "未找到 Docker 配置文件 (${DOCKER_CONFIG})，请先执行: docker login"
    exit 1
fi

# -----------------------------------------------------------
# 4. 构建 Docker 镜像
# -----------------------------------------------------------
log_step "构建 Docker 镜像: ${IMAGE}"

docker build . -t "$IMAGE"

log_info "镜像构建成功: ${IMAGE}"

# -----------------------------------------------------------
# 5. 推送镜像到 Docker Hub
# -----------------------------------------------------------
log_step "推送镜像到 Docker Hub: ${IMAGE}"

docker push "$IMAGE"

log_info "镜像推送成功: ${IMAGE}"

# -----------------------------------------------------------
# 6. 可选: 更新 latest 标签
# -----------------------------------------------------------
echo ""
read -r -p "是否同时更新 latest 标签并推送？(y/N): " tag_latest

if [[ "$tag_latest" =~ ^[Yy]$ ]]; then
    log_step "更新 latest 标签..."

    docker tag "$IMAGE" "$LATEST"
    log_info "已打标签: ${LATEST}"

    docker push "$LATEST"
    log_info "latest 推送成功。"

    echo ""
    log_info "要验证 latest 是否正确，可以运行:"
    echo "        docker run --rm ${LATEST} --version"
fi

# -----------------------------------------------------------
# 完成
# -----------------------------------------------------------
echo ""
log_info "============================================"
log_info "  发布完成！"
log_info "  版本标签: ${IMAGE}"
if [[ "$tag_latest" =~ ^[Yy]$ ]]; then
    log_info "  latest 标签: 已同步"
fi
log_info "============================================"
echo ""
log_info "可运行以下命令验证本地镜像:"
echo "        docker images hellocoolguy/homebox"
