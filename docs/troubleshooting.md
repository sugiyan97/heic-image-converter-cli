# トラブルシューティング

## よくある問題

### ビルドエラー: CGO関連のエラー

**問題**: `go build`時にCGO関連のエラーが発生する

**解決方法**:
- CGOが有効になっているか確認: `CGO_ENABLED=1`
- 必要なシステムライブラリがインストールされているか確認
  - Linux: `sudo apt-get install libheif-dev libde265-dev libx265-dev`
  - macOS: `brew install libheif`

### 変換エラー: ファイルが見つからない

**問題**: `エラー: ファイルまたはディレクトリが見つかりません`が表示される

**解決方法**:
- ファイルパスが正しいか確認
- ファイルが存在するか確認: `ls -la <ファイルパス>`
- 相対パスと絶対パスの違いを確認

### 変換エラー: HEICファイルのデコード失敗

**問題**: `✗ 変換失敗`が表示される

**解決方法**:
- HEICファイルが破損していないか確認
- ファイルが正しいHEIC形式か確認
- 他のHEICファイルで試してみる

### EXIF情報が表示されない

**問題**: `--show-exif`オプションを使用してもEXIF情報が表示されない

**解決方法**:
- 元のHEICファイルにEXIF情報が含まれているか確認
- 他のツール（例: `exiftool`）でEXIF情報を確認

### インストールエラー: スクリプトの実行権限がない（macOS）

**問題**: `./install.sh`を実行すると「Permission denied」エラーが表示される

**解決方法**:
```bash
chmod +x install.sh
./install.sh
```

### インストールスクリプトの文字化け（Windows）

**問題**: `install.bat`や`uninstall.bat`を実行すると日本語が文字化けする

**解決方法**:
- 最新のリリースでは文字化け対策が含まれています。最新版をダウンロードしてください
- それでも文字化けする場合は、コマンドプロンプトのコードページをUTF-8に変更:
  ```cmd
  chcp 65001
  ```

### インストールエラー: PowerShell実行ポリシー（Windows）

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

### インストール後、`heic-convert`コマンドが見つからない

**問題**: インストール後も`heic-convert`コマンドが実行できない

**解決方法**:
- PATH設定を確認してください
- 新しいターミナル/コマンドプロンプトを開いてください
- 手動でPATHに追加する場合は、以下を実行：
  - macOS: `export PATH="$HOME/bin/HeicConverter:$PATH"`をシェル設定ファイルに追加
  - Windows: 環境変数の設定から`%USERPROFILE%\bin\HeicConverter`をPATHに追加

**Windows の場合（`convert` と入力すると「無効なドライブ指定です」と表示される）**:
- Windows にはシステムの `convert.exe`（FAT を NTFS に変換するコマンド）が標準で含まれており、PATH 上で先に検出されることがあります。
- 本ツールではコマンド名を **`heic-convert`** にしています。`heic-convert --help` や `heic-convert input.HEIC` のように `heic-convert` を使ってください。

### Windows: DLLが見つからないエラー

**問題**: `libstdc++-6.dll`、`libwinpthread-1.dll`、`libgcc_s_seh-1.dll`が見つからない

**解決方法**:
- **最新のリリース（v1.0.1以降）**: 必要なDLLはZIPファイルに同梱されています。ZIPファイルを展開した際に、バイナリと同じディレクトリにDLLが含まれていることを確認してください
- **古いバージョンを使用している場合**:
  - MSYS2/MinGWがインストールされている場合、`C:\tools\msys64\mingw64\bin`をPATHに追加してください
  - または、必要なDLLを`%USERPROFILE%\bin\HeicConverter`にコピーしてください
  - 必要なDLL: `libgcc_s_seh-1.dll`、`libwinpthread-1.dll`、`libstdc++-6.dll`

## サポート

問題が解決しない場合は、[GitHub Issues](https://github.com/sugiyan97/heic-image-converter-cli/issues)で報告してください。
