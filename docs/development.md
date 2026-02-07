# 開発

## DevContainersを使用する場合

このプロジェクトはDevContainersに対応しています。VS Codeで開くと、自動的に開発環境がセットアップされます。

詳細は[.devcontainer/README.md](../.devcontainer/README.md)を参照してください。

## ローカル開発環境のセットアップ

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

## 利用可能なMakeコマンド

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

テストの詳細については、[test-cases.md](test-cases.md)を参照してください。
