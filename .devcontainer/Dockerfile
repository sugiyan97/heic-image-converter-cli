FROM mcr.microsoft.com/devcontainers/go:1.25

# libde265 とその他の必要な依存関係をインストール
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    libheif-dev \
    libde265-dev \
    libx265-dev \
    pkg-config \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Set CGO_ENABLED environment variable
ENV CGO_ENABLED=1

# 作業ディレクトリを設定
WORKDIR /workspace

