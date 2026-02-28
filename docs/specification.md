# HEIC Image Converter 仕様書

## 1. 概要

### 1.1 プロダクト名

HEIC Image Converter

### 1.2 目的

HEIC（High Efficiency Image Container）形式の画像ファイルをJPEG形式に変換するためのコマンドラインツール。将来的に他の画像形式への変換もサポート予定。

### 1.3 対応プラットフォーム

- Windows (amd64)
- macOS (arm64) - Apple Siliconのみサポート（Intel版は非サポート）
- Linux (amd64/arm64)

### 1.4 技術スタック

- **言語**: Go 1.23以上
- **主要ライブラリ**:
  - `github.com/adrium/goheif`: HEIC形式のデコード
  - `github.com/dsoprea/go-exif/v3`: EXIF情報の処理
  - `github.com/dsoprea/go-jpeg-image-structure/v2`: JPEG構造の操作

## 2. 機能仕様

### 2.1 コア機能

#### 2.1.1 HEICからJPEGへの変換

- **入力形式**: HEIC (.heic)
- **出力形式**: JPEG (.jpg)
- **品質設定**: JPEG品質95（固定）
- **出力先**: 入力ファイルと同じディレクトリ
- **ファイル名**: 入力ファイル名の拡張子を`.jpg`に変更

**処理フロー**:
1. HEICファイルを開く
2. HEIC画像をデコード
3. RGB形式に変換（RGBA/NRGBA/YCbCrからRGBAへ）
4. EXIF情報の処理（保持/削除）
5. JPEG形式でエンコード
6. 出力ファイルに保存

#### 2.1.2 元HEICファイルの削除

- 変換完了後、元のHEICファイルを削除する機能
- **注意**: 現時点では実装されていないが、将来的な機能として想定

### 2.2 EXIF情報管理機能

#### 2.2.1 EXIF情報の表示

- `--show-exif`オプションでEXIF情報を表示
- 表示される主要なEXIFタグ:
  - DateTime, DateTimeOriginal
  - Make, Model
  - Orientation
  - XResolution, YResolution, ResolutionUnit
  - Software, Artist, Copyright
  - ExifVersion
  - Flash, FocalLength, FNumber, ExposureTime, ISOSpeedRatings
  - GPSInfo
  - ImageWidth, ImageLength

#### 2.2.2 EXIF情報の削除

- `--remove-exif`オプションでEXIF情報を削除して変換
- プライバシー保護の用途に使用可能

#### 2.2.3 EXIF情報の保持

- デフォルトでは可能な限りEXIF情報を保持
- `--remove-exif`が指定されていない場合、元のHEICファイルからEXIF情報を抽出し、JPEGファイルに埋め込む

#### 2.2.4 EXIF情報のチェック

- `--check-exif`オプションでJPEGファイルのEXIF情報の有無をチェック
- チェック結果:
  - ✓ EXIF情報は削除されています
  - ✗ EXIF情報が残っています（検出された主要なEXIFタグも表示）
- サマリー表示:
  - 総ファイル数
  - EXIF削除済み数
  - EXIF残存数
  - エラー数

### 2.3 ファイル処理モード

#### 2.3.1 カレントディレクトリ処理

- 引数なしで実行した場合、カレントディレクトリ内の全HEICファイルを再帰的に検索して変換

#### 2.3.2 単一ファイル処理

- ファイルパスを1つ指定して実行した場合、そのファイルのみを変換

#### 2.3.3 ディレクトリ処理

- ディレクトリパスを1つ指定して実行した場合、そのディレクトリ内の全HEICファイルを再帰的に検索して変換

## 3. コマンドライン仕様

### 3.1 基本コマンド

```bash
heic-convert [オプション] [ファイル/ディレクトリ]
```

### 3.2 オプション

| オプション | 説明 |
|-----------|------|
| `--show-exif` | EXIF情報を表示 |
| `--remove-exif` | EXIF情報を削除して変換 |
| `--check-exif` | JPGファイルのEXIF削除をチェック |

### 3.3 使用例

#### 3.3.1 変換コマンド

```bash
# カレントディレクトリの全HEICファイルを変換
heic-convert

# 単一ファイルを変換
heic-convert input.HEIC

# 指定ディレクトリ内の全HEICファイルを変換
heic-convert /path/to/directory

# EXIF情報を表示して変換
heic-convert --show-exif input.HEIC

# EXIF情報を削除して変換
heic-convert --remove-exif input.HEIC

# EXIF情報を表示してから削除して変換
heic-convert --show-exif --remove-exif input.HEIC
```

#### 3.3.2 EXIFチェックコマンド

```bash
# カレントディレクトリの全JPGファイルのEXIFをチェック
heic-convert --check-exif

# 単一JPGファイルのEXIFをチェック
heic-convert --check-exif input.jpg

# 指定ディレクトリ内の全JPGファイルのEXIFをチェック
heic-convert --check-exif /path/to/directory
```

## 4. 画像変換仕様

### 4.1 色空間変換

- **入力**: HEIC形式（様々な色空間に対応）
- **中間処理**: RGBA形式に統一
- **出力**: JPEG形式（RGB）

**対応する入力色空間**:

- RGBA: 直接変換（最適化済み）
- NRGBA: 直接変換（最適化済み）
- YCbCr: 直接変換（最適化済み）
- その他: 汎用変換処理

### 4.2 アルファチャンネル処理

- アルファチャンネルが存在する場合、白背景に合成してからJPEGに変換
- JPEGはアルファチャンネルをサポートしないため、事前に合成が必要

### 4.3 品質設定

- JPEG品質: 95（固定）
- 将来的に品質設定をオプション化する可能性あり

## 5. エラーハンドリング

### 5.1 エラーケース

| エラーケース | 処理 |
|------------|------|
| ファイルが見つからない | エラーメッセージを表示してスキップ |
| HEICファイルのデコード失敗 | エラーメッセージを表示してスキップ |
| JPEGファイルのエンコード失敗 | エラーメッセージを表示してスキップ |
| EXIF情報の抽出失敗 | 警告を表示して続行（EXIFなしで変換） |
| EXIF情報の埋め込み失敗 | エラーメッセージを表示してスキップ |

### 5.2 エラーメッセージ形式

- 変換失敗: `✗ 変換失敗: [ファイル名]`
- 変換成功: `✓ 変換完了: [ファイル名]`
- 警告: `警告: [警告内容]`

## 6. パフォーマンス仕様

### 6.1 処理速度

- 画像サイズと複雑さに依存
- 最適化された色空間変換処理により、RGBA/NRGBA/YCbCr形式の場合は高速処理

### 6.2 メモリ使用量

- 画像全体をメモリに読み込む必要があるため、大きな画像ファイルの場合はメモリ使用量が増加
- ストリーミング処理は未対応

## 7. 将来の拡張予定

### 7.1 対応形式の拡張

- **入力形式の拡張**:
  - HEIF (.heif)
  - その他の画像形式（検討中）

- **出力形式の拡張**:
  - PNG
  - WebP
  - AVIF
  - その他の形式（検討中）

### 7.2 機能の拡張

- 品質設定のオプション化
- リサイズ機能
- バッチ処理の最適化
- 並列処理による高速化
- 元HEICファイルの自動削除オプション
- 出力ディレクトリの指定
- 出力ファイル名のカスタマイズ

### 7.3 ユーザビリティの向上

- プログレスバーの表示
- 詳細ログモード
- 設定ファイルのサポート

## 8. 制限事項

### 8.1 現在の制限

- JPEG品質が固定（95）
- 出力ディレクトリが入力ファイルと同じ（変更不可）
- 出力ファイル名が自動生成（カスタマイズ不可）
- 並列処理未対応（順次処理）
- ストリーミング処理未対応（メモリに全画像を読み込む）

### 8.2 技術的制限

- `goheif`ライブラリがCGOを必要とするため、クロスコンパイルが複雑
- macOS向けバイナリはmacOS環境でビルドする必要がある（osxcrossが必要）

## 9. ビルド仕様

### 9.1 ビルド要件

- Go 1.23以上
- CGO有効（goheifライブラリの要件）
- プラットフォーム固有のビルドツール（macOSの場合はosxcrossまたはmacOS環境）

### 9.2 ビルドコマンド

```bash
# 現在のプラットフォーム向けにビルド
make build-go

# Windows向けにビルド
make build-windows

# macOS向けにビルド
make build-macos

# すべてのプラットフォーム向けにビルド
make build-all
```

### 9.3 出力バイナリ

- `bin/convert`: 現在のプラットフォーム向け
- `bin/convert-windows-amd64.exe`: Windows向け
- `bin/convert-darwin-arm64`: macOS (Apple Silicon)向け

## 10. テスト仕様

### 10.1 テスト実装

このプロジェクトには包括的なテストスイートが実装されています。テストは`main_test.go`に実装されており、以下の内容が含まれています：

- **機能テスト**: HEICからJPEGへの変換、EXIF情報の処理、エラーハンドリングなど
- **統合テスト**: 複数の操作を組み合わせたエンドツーエンドテスト
- **ユニットテスト**: 個別の関数の動作確認

### 10.2 テスト実行方法

```bash
# ビルド（テスト実行に必要）
make build-go

# テストを実行
make test
```

または直接実行：

```bash
go test -v ./...
```

### 10.3 テストデータ

テストは`sample/test.HEIC`をテストデータとして使用します。テスト実行時には、このファイルが一時ディレクトリに複製されて使用されます。

### 10.4 CI/CD

GitHub Actionsを使用してCI/CDパイプラインが設定されています：

- **Lint**: `.github/workflows/lint.yml` - コードの静的解析
- **Test**: `.github/workflows/test.yml` - テストの自動実行

Pull Requestを作成すると、自動的にlintとtestが実行されます。

### 10.5 テストケース仕様

詳細なテストケース仕様については[テストケース仕様書](test-cases.md)を参照してください。

### 10.6 テスト対象

- HEICファイルのデコード
- JPEGファイルのエンコード
- EXIF情報の抽出・埋め込み・削除
- 色空間変換
- エラーハンドリング

### 10.7 テスト環境

- 各プラットフォーム（Windows/macOS/Linux）での動作確認
- 様々なサイズのHEICファイルでのテスト
- EXIF情報の有無による動作確認

## 11. セキュリティ考慮事項

### 11.1 ファイルアクセス

- 指定されたファイル/ディレクトリのみにアクセス
- シンボリックリンクの処理はOSの標準動作に依存

### 11.2 EXIF情報のプライバシー

- EXIF情報には位置情報（GPS）や撮影日時などの個人情報が含まれる可能性がある
- `--remove-exif`オプションを使用することで、プライバシー保護が可能

## 12. 依存関係

### 12.1 直接依存

- `github.com/adrium/goheif v0.0.0-20230113233934-ca402e77a786`
- `github.com/dsoprea/go-exif/v3 v3.0.1`
- `github.com/dsoprea/go-jpeg-image-structure/v2 v2.0.0-20221012074422-4f3f7e934102`

### 12.2 間接依存

- `github.com/dsoprea/go-iptc`
- `github.com/dsoprea/go-logging`
- `github.com/dsoprea/go-photoshop-info-format`
- `github.com/dsoprea/go-utility/v2`
- `github.com/go-errors/errors`
- `github.com/go-xmlfmt/xmlfmt`
- `github.com/golang/geo`
- `github.com/rwcarlsen/goexif`
- `golang.org/x/net`
- `gopkg.in/yaml.v2`

## 13. バージョン管理

### 13.1 バージョン番号

- セマンティックバージョニング（SemVer）を採用予定
- 現時点では開発中（バージョン未定義）

### 13.2 変更履歴

- 変更履歴は別途管理（CHANGELOG.md等）

## 14. ライセンス

### 14.1 プロジェクトライセンス

- 未定義（要確認）

### 14.2 依存ライブラリのライセンス

- 各依存ライブラリのライセンスに準拠

---

**最終更新日**: 2025年（作成日）
**ドキュメントバージョン**: 1.0
