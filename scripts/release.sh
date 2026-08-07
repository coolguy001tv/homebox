#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# HomeBox Docker Release Script
# 一键发布脚本：构建 → 检查登录 → 推送 → (可选) latest 标签
#
# ⚠️ 手动发版脚本，仅供应急/本地验证。正常情况下镜像由 GitHub Actions
# 自动构建，无需手动执行：
#   - push 到 main  → 自动构建 nightly 镜像（三平台）
#   - push tag v*   → 自动跑测试并构建 release 镜像（latest + vX.Y.Z，三平台）
# 推荐发布方式：git tag vX.Y.Z && git push origin vX.Y.Z
# 注意：手动构建是单平台（仅本机架构），且推 git tag 会再次触发 CI 构建。
#
# 用法:
#   交互模式:  ./scripts/release.sh [版本号]
#   AI/CI 模式: ./scripts/release.sh <版本号> -y [--latest] [--git-tag]
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

usage() {
    cat <<EOF
用法: $0 [选项] <版本号>

选项:
  -y, --yes       非交互模式，跳过所有确认提示
  --latest        同时推送 latest 标签（仅 -y 模式生效）
  --git-tag       同时推送 Git 标签到远程仓库（仅 -y 模式生效）
  -h, --help      显示此帮助信息

示例:
  $0 0.1.15                    交互模式
  $0 0.1.15 -y                 AI 模式，仅构建推送
  $0 0.1.15 -y --latest        AI 模式，同时推送 latest
  $0 0.1.15 -y --latest --git-tag  AI 模式，全量发布
EOF
    exit 0
}

# -----------------------------------------------------------
# 0. 解析命令行参数
# -----------------------------------------------------------
YES_MODE=false
PUSH_LATEST=false
PUSH_GIT_TAG=false
VERSION=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -y|--yes)
            YES_MODE=true
            shift
            ;;
        --latest)
            PUSH_LATEST=true
            shift
            ;;
        --git-tag)
            PUSH_GIT_TAG=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        -*)
            log_error "未知选项: $1"
            echo ""
            usage
            ;;
        *)
            if [[ -z "$VERSION" ]]; then
                VERSION="$1"
            else
                log_error "多余的参数: $1"
                usage
            fi
            shift
            ;;
    esac
done

# -----------------------------------------------------------
# 1. 解析/确认版本号
# -----------------------------------------------------------
if [[ -z "$VERSION" ]]; then
    if $YES_MODE; then
        log_error "非交互模式下必须提供版本号。"
        echo ""
        usage
    fi
    read -r -p "请输入发布版本号 (例如 0.1.15): " VERSION
fi

if [[ -z "$VERSION" ]]; then
    log_error "版本号不能为空，已取消发布。"
    exit 1
fi

# 二次确认（交互模式）
if ! $YES_MODE; then
    echo ""
    read -r -p "确认发布版本 ${VERSION}？(y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        log_info "已取消发布。"
        exit 0
    fi
else
    log_info "非交互模式，发布版本: ${VERSION}"
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
if $YES_MODE; then
    if $PUSH_LATEST; then
        DO_LATEST=true
    else
        DO_LATEST=false
        log_info "跳过 latest 标签（使用 --latest 启用）。"
    fi
else
    echo ""
    read -r -p "是否同时更新 latest 标签并推送？(y/N): " tag_latest
    DO_LATEST=false
    [[ "$tag_latest" =~ ^[Yy]$ ]] && DO_LATEST=true
fi

if $DO_LATEST; then
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
# 7. 可选: 推送 Git 标签
# -----------------------------------------------------------
if $YES_MODE; then
    if $PUSH_GIT_TAG; then
        DO_GIT_TAG=true
    else
        DO_GIT_TAG=false
        log_info "跳过 Git 标签（使用 --git-tag 启用）。"
    fi
else
    echo ""
    read -r -p "是否推送 Git 标签 v${VERSION} 到远程仓库？(y/N): " tag_git
    DO_GIT_TAG=false
    [[ "$tag_git" =~ ^[Yy]$ ]] && DO_GIT_TAG=true
fi

if $DO_GIT_TAG; then
    log_step "推送 Git 标签..."

    if ! git rev-parse --git-dir &> /dev/null; then
        log_error "当前目录不是 Git 仓库，无法推送标签。"
    else
        TAG="v${VERSION}"

        # 检查标签是否已存在
        if git tag -l "$TAG" | grep -q "$TAG"; then
            log_warn "标签 ${TAG} 已存在。"
            if $YES_MODE; then
                log_info "非交互模式：删除旧标签并重新创建。"
                git tag -d "$TAG"
                git push origin --delete "$TAG" 2>/dev/null || true
            else
                read -r -p "是否删除旧标签并重新创建？(y/N): " recreate
                if [[ "$recreate" =~ ^[Yy]$ ]]; then
                    git tag -d "$TAG"
                    git push origin --delete "$TAG" 2>/dev/null || true
                    log_info "旧标签 ${TAG} 已删除。"
                else
                    log_info "已取消 Git 标签。"
                    DO_GIT_TAG=false
                fi
            fi
        fi

        if $DO_GIT_TAG; then
            # 检查是否有未提交的更改
            if ! git diff-index --quiet HEAD -- 2>/dev/null; then
                log_warn "检测到未提交的更改。"
                if $YES_MODE; then
                    log_info "非交互模式：继续打标签。"
                else
                    read -r -p "是否继续打标签？(y/N): " force_tag
                    if [[ ! "$force_tag" =~ ^[Yy]$ ]]; then
                        log_info "已跳过 Git 标签。"
                        DO_GIT_TAG=false
                    fi
                fi
            fi
        fi

        if $DO_GIT_TAG; then
            git tag -a "$TAG" -m "Release ${TAG}"
            git push origin "$TAG"
            log_info "Git 标签 ${TAG} 已推送。"
        fi
    fi
fi

# -----------------------------------------------------------
# 完成
# -----------------------------------------------------------
echo ""
log_info "============================================"
log_info "  发布完成！"
log_info "  版本标签: ${IMAGE}"
if $DO_LATEST; then
    log_info "  latest 标签: 已同步"
fi
if $DO_GIT_TAG; then
    log_info "  Git 标签: v${VERSION} 已推送"
fi
log_info "============================================"
echo ""
log_info "可运行以下命令验证本地镜像:"
echo "        docker images hellocoolguy/homebox"
