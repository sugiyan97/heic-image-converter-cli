# HEIC Image Converter

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
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

### ビルド済みバイナリを使用する場合

GitHub Releasesから各プラットフォーム用のバイナリをダウンロードしてください。

- [最新リリース](https://github.com/sugiyan97/heic-image-converter-cli/releases/latest)

### ソースからビルドする場合

#### 前提条件

- **Go**: 1.23以上
- **CGO**: 有効（`goheif`ライブラリの要件）
- **システム依存ライブラリ**:
  - **Linux**: `libheif-dev`, `libde265-dev`, `x265-dev`
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
  - Linux: `sudo apt-get install libheif-dev libde265-dev x265-dev`
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
