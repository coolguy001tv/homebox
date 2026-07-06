# [Choice] Node.js version (use -bookworm variants on local arm64/Apple Silicon): 22, 20, 18, 22-bookworm, 20-bookworm, 18-bookworm
ARG VARIANT=22-bookworm
FROM mcr.microsoft.com/devcontainers/typescript-node:1-${VARIANT}

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin