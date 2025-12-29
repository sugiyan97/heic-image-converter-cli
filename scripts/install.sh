#!/bin/bash

# HEIC Image Converter インストールスクリプト (macOS)

set -e

# 色の定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 固定インストール先
INSTALL_DIR="$HOME/bin/HeicConverter"
BINARY_NAME="convert-darwin-arm64"
UNINSTALL_SCRIPT="uninstall.sh"

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

# カレントディレクトリからバイナリを検出
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$SCRIPT_DIR/$BINARY_NAME"

# バイナリの存在確認
if [ ! -f "$BINARY_PATH" ]; then
    error "バイナリファイルが見つかりません: $BINARY_PATH"
    error "このスクリプトは、バイナリファイルと同じディレクトリで実行してください。"
    exit 1
fi

info "HEIC Image Converter のインストールを開始します..."
info "インストール先: $INSTALL_DIR"

# インストール先ディレクトリの作成
if [ ! -d "$INSTALL_DIR" ]; then
    info "インストール先ディレクトリを作成します..."
    mkdir -p "$INSTALL_DIR"
fi

# 既存バイナリの確認
if [ -f "$INSTALL_DIR/convert" ]; then
    warn "既存のバイナリが見つかりました。上書きして更新します。"
fi

# バイナリのコピー
info "バイナリをコピーしています..."
cp "$BINARY_PATH" "$INSTALL_DIR/convert"

# 実行権限の付与
info "実行権限を設定しています..."
chmod +x "$INSTALL_DIR/convert"

# macOSのquarantine属性を削除
info "quarantine属性を削除しています..."
xattr -d com.apple.quarantine "$INSTALL_DIR/convert" 2>/dev/null || true

# アンインストールスクリプトのコピー
if [ -f "$SCRIPT_DIR/$UNINSTALL_SCRIPT" ]; then
    info "アンインストールスクリプトをコピーしています..."
    cp "$SCRIPT_DIR/$UNINSTALL_SCRIPT" "$INSTALL_DIR/$UNINSTALL_SCRIPT"
    chmod +x "$INSTALL_DIR/$UNINSTALL_SCRIPT"
else
    warn "アンインストールスクリプトが見つかりません: $SCRIPT_DIR/$UNINSTALL_SCRIPT"
fi

# PATH設定の確認
info ""
info "PATH設定について"
info "インストール先 ($INSTALL_DIR) をPATHに追加すると、どこからでも 'convert' コマンドを実行できます。"

# シェルの検出
SHELL_CONFIG=""
if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ] || [ -f "$HOME/.bash_profile" ]; then
    SHELL_CONFIG="$HOME/.bash_profile"
else
    # デフォルトでzshを試す
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        SHELL_CONFIG="$HOME/.bash_profile"
    fi
fi

if [ -n "$SHELL_CONFIG" ]; then
    # 既にPATHに追加されているか確認
    PATH_LINE="export PATH=\"$INSTALL_DIR:\$PATH\""
    if grep -q "$INSTALL_DIR" "$SHELL_CONFIG" 2>/dev/null; then
        info "既にPATHに追加されています: $SHELL_CONFIG"
    else
        read -p "PATHに追加しますか？ (Y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
            info "PATH設定を追加しています: $SHELL_CONFIG"
            echo "" >> "$SHELL_CONFIG"
            echo "# HEIC Image Converter" >> "$SHELL_CONFIG"
            echo "$PATH_LINE" >> "$SHELL_CONFIG"
            info "PATH設定を追加しました。"
            info "新しいターミナルを開くか、以下のコマンドを実行してください:"
            info "  source $SHELL_CONFIG"
        else
            info "PATH設定をスキップしました。"
            info "手動でPATHに追加する場合は、以下を $SHELL_CONFIG に追加してください:"
            info "  $PATH_LINE"
        fi
    fi
else
    warn "シェル設定ファイルが見つかりませんでした。"
    warn "手動でPATHに追加する場合は、以下をシェル設定ファイルに追加してください:"
    warn "  export PATH=\"$INSTALL_DIR:\$PATH\""
fi

info ""
info "インストールが完了しました！"
info ""
info "使用方法:"
if grep -q "$INSTALL_DIR" "${SHELL_CONFIG:-$HOME/.zshrc}" 2>/dev/null || grep -q "$INSTALL_DIR" "${SHELL_CONFIG:-$HOME/.bash_profile}" 2>/dev/null; then
    info "  convert --help"
else
    info "  $INSTALL_DIR/convert --help"
fi
info ""
info "アンインストール方法:"
info "  convert --uninstall"
info "  または"
info "  $INSTALL_DIR/$UNINSTALL_SCRIPT"

