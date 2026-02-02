# HEIC Image Converter

[![Go Version](https://img.shields.io/badge/Go-1.25.6+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![GitHub](https://img.shields.io/badge/GitHub-sugiyan97%2Fheic--image--converter--cli-black?style=flat-square&logo=github)](https://github.com/sugiyan97/heic-image-converter-cli)

HEIC（High Efficiency Image Container）形式の画像ファイルを他の画像形式に変換するコマンドラインツールです。現時点ではJPEG形式への変換をサポートしています。

## 目次

- [機能](#機能)
- [クイックスタート](#クイックスタート)
- [インストール](#インストール)
- [使用方法](#使用方法)
- [トラブルシューティング](#トラブルシューティング)
- [開発](#開発)
- [テスト](#テスト)
- [ライセンス](#ライセンス)
- [貢献](#貢献)

## 機能

- ✅ **HEICからJPEGへの変換** - HEIC形式の画像をJPEG形式に変換
- ✅ **単一ファイル変換** - 指定したHEICファイルを個別に変換
- ✅ **ディレクトリ一括変換** - ディレクトリ内の全HEICファイルを再帰的に検索して一括変換
- ✅ **EXIF情報の管理** - EXIF情報の保持・削除・表示・チェック機能
- ✅ **高品質変換** - JPEG品質95で高品質な変換を実現
- ✅ **クロスプラットフォーム** - Windows、macOS、Linuxに対応

## クイックスタート

### 最も簡単な使い方

```bash
# カレントディレクトリの全HEICファイルを変換
convert
```

### 単一ファイルの変換

```bash
# ファイルを指定して変換
convert photo.HEIC
```

### EXIF情報を削除して変換

```bash
# プライバシー保護のためEXIF情報を削除
convert --remove-exif photo.HEIC
```

## インストール

### ビルド済みバイナリを使用する場合（推奨）

GitHub Releasesから各プラットフォーム用のZIPファイルをダウンロードしてインストールします。

- [最新リリース](https://github.com/sugiyan97/heic-image-converter-cli/releases/latest)

#### macOS の場合

1. `convert-darwin-arm64.zip`をダウンロード
2. ZIPファイルを展開
3. ターミナルで展開したディレクトリに移動し、以下を実行：

   ```bash
   ./install.sh
   ```

   必要であれば、実行権限を付与：

   ```bash
   chmod +x install.sh
   ```
4. インストール先は自動的に`~/bin/HeicConverter`に設定されます
5. PATH設定の確認（Y/n）で、Yを選択するとPATHに自動追加されます

#### Windows の場合

1. `convert-windows-amd64.zip`をダウンロード
2. ZIPファイルを展開
3. PowerShellまたはコマンドプロンプトで展開したディレクトリに移動し、以下を実行：
   ```powershell
   # PowerShellの場合
   .\install.ps1
   ```
   または
   ```cmd
   # コマンドプロンプトの場合
   install.bat
   ```
4. インストール先は自動的に`%USERPROFILE%\bin\HeicConverter`に設定されます
5. PATH設定の確認（Y/n）で、Yを選択するとPATHに自動追加されます

#### 固定インストール先

- **macOS**: `~/bin/HeicConverter`
- **Windows**: `%USERPROFILE%\bin\HeicConverter`

インストール先は固定されており、変更できません。これにより、シンプルで一貫したインストール体験を提供します。

#### アップデート

新しいバージョンをインストールする場合は、同じ手順でインストールスクリプトを実行してください。既存のバイナリが自動的に上書きされます。

#### アンインストール

以下のいずれかの方法でアンインストールできます：

1. **バイナリから直接実行**（推奨）:
   ```bash
   convert --uninstall
   ```

2. **アンインストールスクリプトを直接実行**:
   ```bash
   # macOS
   ~/bin/HeicConverter/uninstall.sh
   
   # Windows
   %USERPROFILE%\bin\HeicConverter\uninstall.ps1
   # または
   %USERPROFILE%\bin\HeicConverter\uninstall.bat
   ```

アンインストール時は、`HeicConverter`フォルダ全体が削除されます。PATH設定も削除するかどうかを選択できます。

### ソースからビルドする場合

#### 前提条件

- **Go**: 1.25.6以上
- **CGO**: 有効（`goheif`ライブラリの要件）
- **システム依存ライブラリ**:
  - **Linux**: `libheif-dev`, `libde265-dev`, `libx265-dev`
  - **macOS**: Homebrewで`libheif`をインストール
  - **Windows**: 適切なCライブラリが必要

#### ビルド手順

```bash
# リポジトリのクローン
git clone https://github.com/sugiyan97/heic-image-converter-cli.git
cd heic-image-converter-cli

# 依存関係のダウンロード
make deps

# 現在のプラットフォーム向けにビルド
make build

# または、すべてのプラットフォーム向けにビルド
make build-all
```

ビルドされたバイナリは`bin/`ディレクトリに生成されます。

## 使用方法

### 基本的な使用方法

```bash
# カレントディレクトリの全HEICファイルを変換
convert

# 単一ファイルを変換
convert input.HEIC

# 指定ディレクトリ内の全HEICファイルを変換
convert /path/to/directory
```

### オプション

#### EXIF情報の表示

```bash
# EXIF情報を表示してから変換
convert --show-exif input.HEIC

# ディレクトリ内の全ファイルのEXIF情報を表示
convert --show-exif /path/to/directory
```

#### EXIF情報の削除

```bash
# EXIF情報を削除して変換（プライバシー保護）
convert --remove-exif input.HEIC

# ディレクトリ内の全ファイルからEXIF情報を削除
convert --remove-exif /path/to/directory
```

#### EXIF情報の表示と削除を同時に実行

```bash
# EXIF情報を表示してから削除して変換
convert --show-exif --remove-exif input.HEIC
```

#### EXIF情報のチェック

```bash
# カレントディレクトリの全JPEGファイルのEXIF情報をチェック
convert --check-exif

# 単一JPEGファイルのEXIF情報をチェック
convert --check-exif input.jpg

# 指定ディレクトリ内の全JPEGファイルのEXIF情報をチェック
convert --check-exif /path/to/directory
```

### 使用例

#### 例1: iPhoneで撮影した写真を一括変換

```bash
# 写真フォルダ内の全HEICファイルをJPEGに変換
convert ~/Pictures/iPhone
```

#### 例2: プライバシー保護のためEXIF情報を削除

```bash
# SNSに投稿する前にEXIF情報を削除
convert --remove-exif --show-exif ~/Pictures/iPhone
```

#### 例3: 変換後のJPEGファイルのEXIF情報を確認

```bash
# 変換後のJPEGファイルにEXIF情報が残っていないか確認
convert --check-exif ~/Pictures/iPhone
```

## トラブルシューティング

### よくある問題

#### ビルドエラー: CGO関連のエラー

**問題**: `go build`時にCGO関連のエラーが発生する

**解決方法**:
- CGOが有効になっているか確認: `CGO_ENABLED=1`
- 必要なシステムライブラリがインストールされているか確認
  - Linux: `sudo apt-get install libheif-dev libde265-dev libx265-dev`
  - macOS: `brew install libheif`

#### 変換エラー: ファイルが見つからない

**問題**: `エラー: ファイルまたはディレクトリが見つかりません`が表示される

**解決方法**:
- ファイルパスが正しいか確認
- ファイルが存在するか確認: `ls -la <ファイルパス>`
- 相対パスと絶対パスの違いを確認

#### 変換エラー: HEICファイルのデコード失敗

**問題**: `✗ 変換失敗`が表示される

**解決方法**:
- HEICファイルが破損していないか確認
- ファイルが正しいHEIC形式か確認
- 他のHEICファイルで試してみる

#### EXIF情報が表示されない

**問題**: `--show-exif`オプションを使用してもEXIF情報が表示されない

**解決方法**:
- 元のHEICファイルにEXIF情報が含まれているか確認
- 他のツール（例: `exiftool`）でEXIF情報を確認

#### インストールエラー: スクリプトの実行権限がない（macOS）

**問題**: `./install.sh`を実行すると「Permission denied」エラーが表示される

**解決方法**:
```bash
chmod +x install.sh
./install.sh
```

#### インストールスクリプトの文字化け（Windows）

**問題**: `install.bat`や`uninstall.bat`を実行すると日本語が文字化けする

**解決方法**:
- 最新のリリースでは文字化け対策が含まれています。最新版をダウンロードしてください
- それでも文字化けする場合は、コマンドプロンプトのコードページをUTF-8に変更:
  ```cmd
  chcp 65001
  ```

#### インストールエラー: PowerShell実行ポリシー（Windows）

**問題**: PowerShellスクリプトが実行できない

**解決方法**:
- **推奨**: バッチファイル（`install.bat`、`uninstall.bat`）を使用してください
- または、実行ポリシーを変更:
  ```powershell
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```
- または、実行ポリシーをバイパスして実行:
  ```powershell
  powershell -ExecutionPolicy Bypass -File .\install.ps1
  powershell -ExecutionPolicy Bypass -File .\uninstall.ps1
  ```

#### インストール後、`convert`コマンドが見つからない

**問題**: インストール後も`convert`コマンドが実行できない

**解決方法**:
- PATH設定を確認してください
- 新しいターミナル/コマンドプロンプトを開いてください
- 手動でPATHに追加する場合は、以下を実行：
  - macOS: `export PATH="$HOME/bin/HeicConverter:$PATH"`をシェル設定ファイルに追加
  - Windows: 環境変数の設定から`%USERPROFILE%\bin\HeicConverter`をPATHに追加

#### Windows: DLLが見つからないエラー

**問題**: `libstdc++-6.dll`、`libwinpthread-1.dll`、`libgcc_s_seh-1.dll`が見つからない

**解決方法**:
- **最新のリリース（v1.0.1以降）**: 必要なDLLはZIPファイルに同梱されています。ZIPファイルを展開した際に、バイナリと同じディレクトリにDLLが含まれていることを確認してください
- **古いバージョンを使用している場合**:
  - MSYS2/MinGWがインストールされている場合、`C:\tools\msys64\mingw64\bin`をPATHに追加してください
  - または、必要なDLLを`%USERPROFILE%\bin\HeicConverter`にコピーしてください
  - 必要なDLL: `libgcc_s_seh-1.dll`、`libwinpthread-1.dll`、`libstdc++-6.dll`

### サポート

問題が解決しない場合は、[GitHub Issues](https://github.com/sugiyan97/heic-image-converter-cli/issues)で報告してください。

## 開発

### DevContainersを使用する場合

このプロジェクトはDevContainersに対応しています。VS Codeで開くと、自動的に開発環境がセットアップされます。

詳細は[.devcontainer/README.md](.devcontainer/README.md)を参照してください。

### ローカル開発環境のセットアップ

```bash
# リポジトリのクローン
git clone https://github.com/sugiyan97/heic-image-converter-cli.git
cd heic-image-converter-cli

# 依存関係のインストール
make deps

# ビルド
make build

# テストの実行
make test

# リンターの実行
make lint
```

### 利用可能なMakeコマンド

```bash
make help          # 利用可能なコマンドを表示
make build         # 現在のプラットフォーム向けにビルド
make build-all     # すべてのプラットフォーム向けにビルド
make test          # テストを実行
make test-coverage # カバレッジ付きテストを実行
make lint          # リンターを実行
make clean         # ビルド成果物を削除
```

## テスト

```bash
# テストの実行
make test

# カバレッジ付きテストの実行
make test-coverage
```

テストの詳細については、[docs/test-cases.md](docs/test-cases.md)を参照してください。

## ライセンス

このプロジェクトは[MIT License](LICENSE)の下で公開されています。

## 貢献

プルリクエストを歓迎します！

### 貢献の手順

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

### 貢献ガイドライン

- 大きな変更の場合は、まず[イシュー](https://github.com/sugiyan97/heic-image-converter-cli/issues)を開いて変更内容を議論してください
- コードは既存のスタイルに合わせてください
- テストを追加・更新してください
- ドキュメントを更新してください

---

**プロジェクト**: [HEIC Image Converter](https://github.com/sugiyan97/heic-image-converter-cli)  
**作成者**: [sugiyan97](https://github.com/sugiyan97)
