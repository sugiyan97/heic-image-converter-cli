# HEIC Image Converter アンインストールスクリプト (Windows PowerShell)
#
# 実行ポリシーエラーが発生する場合:
# PowerShellの実行ポリシーが制限されている場合、以下のいずれかの方法で実行してください:
#
# 方法1: 実行ポリシーをバイパスして実行
#   powershell -ExecutionPolicy Bypass -File .\uninstall.ps1
#
# 方法2: 実行ポリシーを変更（推奨）
#   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
#   その後、通常通り .\uninstall.ps1 を実行
#
# 方法3: uninstall.bat を使用（推奨）
#   コマンドプロンプトまたはPowerShellで uninstall.bat を実行

$ErrorActionPreference = "Stop"

# 固定インストール先
$InstallDir = Join-Path $env:USERPROFILE "bin\HeicConverter"

# メッセージ表示関数
function Write-Info {
    Write-Host "[INFO] $args" -ForegroundColor Green
}

function Write-Warn {
    Write-Host "[WARN] $args" -ForegroundColor Yellow
}

function Write-Error {
    Write-Host "[ERROR] $args" -ForegroundColor Red
}

# インストール先の確認
if (-not (Test-Path $InstallDir)) {
    Write-Error "インストール先が見つかりません: $InstallDir"
    Write-Error "既にアンインストールされている可能性があります。"
    exit 1
}

Write-Info "HEIC Image Converter のアンインストールを開始します..."
Write-Info "削除対象: $InstallDir"

# 削除前の確認
$Response = Read-Host "本当に削除しますか？ (y/N)"
if ($Response -ne "y" -and $Response -ne "Y") {
    Write-Info "アンインストールをキャンセルしました。"
    exit 0
}

# 削除対象ファイルの確認
Write-Info "削除対象ファイルを確認しています..."
$DllFiles = @("libgcc_s_seh-1.dll", "libwinpthread-1.dll", "libstdc++-6.dll")
$DllCount = 0
foreach ($DllFile in $DllFiles) {
    $DllPath = Join-Path $InstallDir $DllFile
    if (Test-Path $DllPath) {
        $DllCount++
    }
}
$BinaryPath = Join-Path $InstallDir "convert.exe"
if (Test-Path $BinaryPath) {
    Write-Info "バイナリファイル: convert.exe"
}
if ($DllCount -gt 0) {
    Write-Info "DLLファイル: $DllCount 個"
}

# 最終確認
Write-Warn "警告: $InstallDir フォルダ全体が削除されます。"
$Response2 = Read-Host "続行しますか？ (y/N)"
if ($Response2 -ne "y" -and $Response2 -ne "Y") {
    Write-Info "アンインストールをキャンセルしました。"
    exit 0
}

# フォルダごと削除（バイナリ、DLL、スクリプトを含む）
Write-Info "インストールフォルダを削除しています..."
Remove-Item -Path $InstallDir -Recurse -Force
Write-Info "インストールフォルダを削除しました。"

# PATH削除の確認
Write-Info ""
$PathResponse = Read-Host "PATH設定も削除しますか？ (y/N)"
if ($PathResponse -eq "y" -or $PathResponse -eq "Y") {
    # 現在のユーザー環境変数のPATHを取得
    $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($CurrentPath -like "*$InstallDir*") {
        Write-Info "PATH設定を削除しています..."
        # PATHから該当パスを削除
        $PathArray = $CurrentPath -split ';' | Where-Object { $_ -ne $InstallDir -and $_ -ne "" }
        $NewPath = $PathArray -join ';'
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Info "PATH設定を削除しました。"
        Write-Info "新しいコマンドプロンプトまたはPowerShellを開いてください。"
    } else {
        Write-Info "PATH設定が見つかりませんでした。"
    }
} else {
    Write-Info "PATH設定の削除をスキップしました。"
}

Write-Info ""
Write-Info "アンインストールが完了しました！"

