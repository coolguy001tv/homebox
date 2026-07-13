#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# HomeBox Git Tag 发布脚本
# 为指定版本号创建并推送 Git 标签
# 用法: ./scripts/git-tag.sh [版本号]
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
    read -r -p "请输入版本号 (例如 0.1.13): " VERSION
fi

if [[ -z "$VERSION" ]]; then
    log_error "版本号不能为空，已取消。"
    exit 1
fi

TAG="v${VERSION}"

# 二次确认
echo ""
read -r -p "确认推送 Git 标签 ${TAG} 到远程仓库？(y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    log_info "已取消。"
    exit 0
fi

# -----------------------------------------------------------
# 2. 检查 Git 仓库状态
# -----------------------------------------------------------
log_step "检查 Git 仓库状态..."

if ! git rev-parse --git-dir &> /dev/null; then
    log_error "当前目录不是 Git 仓库。"
    exit 1
fi

# 检查标签是否已存在
if git tag -l "$TAG" | grep -q "$TAG"; then
    log_warn "标签 ${TAG} 已存在。"
    read -r -p "是否删除旧标签并重新创建？(y/N): " recreate
    if [[ "$recreate" =~ ^[Yy]$ ]]; then
        git tag -d "$TAG"
        git push origin --delete "$TAG" 2>/dev/null || true
        log_info "旧标签 ${TAG} 已删除。"
    else
        log_info "已取消。"
        exit 0
    fi
fi

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    log_warn "检测到未提交的更改。"
    read -r -p "是否继续打标签？(y/N): " force_tag
    if [[ ! "$force_tag" =~ ^[Yy]$ ]]; then
        log_info "已取消。"
        exit 0
    fi
fi

# -----------------------------------------------------------
# 3. 创建并推送标签
# -----------------------------------------------------------
log_step "创建标签: ${TAG}"

git tag -a "$TAG" -m "Release ${TAG}"

log_info "标签 ${TAG} 已创建。"

log_step "推送标签到远程仓库..."

git push origin "$TAG"

log_info "标签 ${TAG} 已推送到远程。"
