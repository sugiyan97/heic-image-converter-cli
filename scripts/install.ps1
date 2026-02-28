# HEIC Image Converter インストールスクリプト (Windows PowerShell)

$ErrorActionPreference = "Stop"

# 固定インストール先
$InstallDir = Join-Path $env:USERPROFILE "bin\HeicConverter"
$BinaryName = "convert-windows-amd64.exe"
$UninstallScriptPS1 = "uninstall.ps1"
$UninstallScriptBAT = "uninstall.bat"

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

# カレントディレクトリからバイナリを検出
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinaryPath = Join-Path $ScriptDir $BinaryName

# バイナリの存在確認
if (-not (Test-Path $BinaryPath)) {
    Write-Error "バイナリファイルが見つかりません: $BinaryPath"
    Write-Error "このスクリプトは、バイナリファイルと同じディレクトリで実行してください。"
    exit 1
}

Write-Info "HEIC Image Converter のインストールを開始します..."
Write-Info "インストール先: $InstallDir"

# インストール先ディレクトリの作成
if (-not (Test-Path $InstallDir)) {
    Write-Info "インストール先ディレクトリを作成します..."
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# 既存バイナリの確認と旧 convert.exe の削除
$BinaryDest = Join-Path $InstallDir "heic-convert.exe"
$OldBinaryDest = Join-Path $InstallDir "convert.exe"
if (Test-Path $OldBinaryDest) {
    Write-Warn "旧バイナリ convert.exe を削除します。"
    Remove-Item -Path $OldBinaryDest -Force
}
if (Test-Path $BinaryDest) {
    Write-Warn "既存のバイナリが見つかりました。上書きして更新します。"
}

# バイナリのコピー
Write-Info "バイナリをコピーしています..."
Copy-Item -Path $BinaryPath -Destination $BinaryDest -Force

# MinGWランタイムDLLのコピー
Write-Info "必要なDLLファイルをコピーしています..."
$DllFiles = @("libgcc_s_seh-1.dll", "libwinpthread-1.dll", "libstdc++-6.dll")
$DllCount = 0
foreach ($DllFile in $DllFiles) {
    $DllSource = Join-Path $ScriptDir $DllFile
    if (Test-Path $DllSource) {
        $DllDest = Join-Path $InstallDir $DllFile
        Copy-Item -Path $DllSource -Destination $DllDest -Force
        $DllCount++
    }
}
if ($DllCount -gt 0) {
    Write-Info "$DllCount 個のDLLファイルをコピーしました。"
} else {
    Write-Warn "DLLファイルが見つかりませんでした。バイナリが正常に動作しない可能性があります。"
}

# アンインストールスクリプトのコピー
$UninstallPS1Source = Join-Path $ScriptDir $UninstallScriptPS1
$UninstallBATSource = Join-Path $ScriptDir $UninstallScriptBAT

if (Test-Path $UninstallPS1Source) {
    Write-Info "アンインストールスクリプトをコピーしています..."
    Copy-Item -Path $UninstallPS1Source -Destination (Join-Path $InstallDir $UninstallScriptPS1) -Force
}

if (Test-Path $UninstallBATSource) {
    Copy-Item -Path $UninstallBATSource -Destination (Join-Path $InstallDir $UninstallScriptBAT) -Force
}

# PATH設定の確認
Write-Info ""
Write-Info "PATH設定について"
Write-Info "インストール先 $InstallDir をPATHに追加すると、どこからでも 'heic-convert' コマンドを実行できます。"

# 現在のユーザー環境変数のPATHを取得
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")

# 既にPATHに追加されているか確認
if ($CurrentPath -like "*$InstallDir*") {
    Write-Info "既にPATHに追加されています。"
} else {
    $Response = Read-Host "PATHに追加しますか？ (Y/n)"
    if ($Response -eq "" -or $Response -eq "Y" -or $Response -eq "y") {
        Write-Info "PATH設定を追加しています..."
        $NewPath = $CurrentPath
        if ($NewPath -and -not $NewPath.EndsWith(";")) {
            $NewPath += ";"
        }
        $NewPath += $InstallDir
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Info "PATH設定を追加しました。"
        Write-Info "新しいコマンドプロンプトまたはPowerShellを開いてください。"
    } else {
        Write-Info "PATH設定をスキップしました。"
        Write-Info "手動でPATHに追加する場合は、環境変数の設定から以下を追加してください:"
        Write-Info "  $InstallDir"
    }
}

Write-Info ""
Write-Info "インストールが完了しました！"
Write-Info ""
Write-Info "使用方法:"
if ($CurrentPath -like "*$InstallDir*" -or ($Response -eq "" -or $Response -eq "Y" -or $Response -eq "y")) {
    Write-Info "  heic-convert --help"
} else {
    Write-Info "  $BinaryDest --help"
}
Write-Info ""
Write-Info "アンインストール方法:"
Write-Info "  heic-convert --uninstall"
Write-Info "  または"
Write-Info "  $InstallDir\$UninstallScriptPS1"

