# DevContainer

このディレクトリには、VS Code の [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers) 拡張機能を使って開発環境を構築するための設定が含まれています。

## 構成

- `devcontainer.json`: DevContainer の設定ファイル
- `Dockerfile`: 開発コンテナのビルド定義

`Dockerfile` は `mcr.microsoft.com/devcontainers/go:1.25` イメージをベースに、HEIC変換に必要な以下の依存ライブラリをインストールします。

- `libheif-dev`
- `libde265-dev`
- `libx265-dev`
- `pkg-config`

また、CGoを利用するため `CGO_ENABLED=1` が設定されています。

## 使い方

1. VS Code で本リポジトリを開きます。
2. コマンドパレットから `Dev Containers: Reopen in Container` を実行します。
3. コンテナのビルドが完了すると、`postCreateCommand` により `go mod download` が自動実行され、依存関係が取得された状態で開発を開始できます。

推奨拡張機能として `golang.go` が自動的にインストールされます。
