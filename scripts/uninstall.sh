#!/bin/bash

# HEIC Image Converter アンインストールスクリプト (macOS)

set -e

# 色の定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 固定インストール先
INSTALL_DIR="$HOME/bin/HeicConverter"

# メッセージ表示関数
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# インストール先の確認
if [ ! -d "$INSTALL_DIR" ]; then
    error "インストール先が見つかりません: $INSTALL_DIR"
    error "既にアンインストールされている可能性があります。"
    exit 1
fi

info "HEIC Image Converter のアンインストールを開始します..."
info "削除対象: $INSTALL_DIR"

# 削除前の確認
read -p "本当に削除しますか？ (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "アンインストールをキャンセルしました。"
    exit 0
fi

# 最終確認
warn "警告: $INSTALL_DIR フォルダ全体が削除されます。"
read -p "続行しますか？ (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "アンインストールをキャンセルしました。"
    exit 0
fi

# フォルダごと削除
info "インストールフォルダを削除しています..."
rm -rf "$INSTALL_DIR"
info "インストールフォルダを削除しました。"

# PATH削除の確認
info ""
read -p "PATH設定も削除しますか？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # シェル設定ファイルの検出
    SHELL_CONFIG=""
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        SHELL_CONFIG="$HOME/.bash_profile"
    fi

    if [ -n "$SHELL_CONFIG" ]; then
        # PATH設定の削除
        if grep -q "$INSTALL_DIR" "$SHELL_CONFIG" 2>/dev/null; then
            info "PATH設定を削除しています: $SHELL_CONFIG"
            # HEIC Image Converter関連の行を削除
            sed -i.bak "/# HEIC Image Converter/,+1d" "$SHELL_CONFIG" 2>/dev/null || \
            sed -i '' "/# HEIC Image Converter/,+1d" "$SHELL_CONFIG" 2>/dev/null || \
            sed -i "/$INSTALL_DIR/d" "$SHELL_CONFIG" 2>/dev/null
            # バックアップファイルを削除（存在する場合）
            [ -f "$SHELL_CONFIG.bak" ] && rm "$SHELL_CONFIG.bak"
            info "PATH設定を削除しました。"
            info "新しいターミナルを開くか、以下のコマンドを実行してください:"
            info "  source $SHELL_CONFIG"
        else
            info "PATH設定が見つかりませんでした。"
        fi
    else
        warn "シェル設定ファイルが見つかりませんでした。"
    fi
else
    info "PATH設定の削除をスキップしました。"
fi

info ""
info "アンインストールが完了しました！"

