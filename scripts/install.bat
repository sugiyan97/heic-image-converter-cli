@echo off
REM HEIC Image Converter インストールスクリプト (Windows Batch)
REM UTF-8文字化け対策
chcp 65001 >nul 2>&1

setlocal enabledelayedexpansion

REM 固定インストール先
set "INSTALL_DIR=%USERPROFILE%\bin\HeicConverter"
set "BINARY_NAME=convert-windows-amd64.exe"
set "UNINSTALL_SCRIPT_PS1=uninstall.ps1"
set "UNINSTALL_SCRIPT_BAT=uninstall.bat"

REM カレントディレクトリからバイナリを検出
set "SCRIPT_DIR=%~dp0"
set "BINARY_PATH=%SCRIPT_DIR%%BINARY_NAME%"

REM バイナリの存在確認
if not exist "%BINARY_PATH%" (
    echo [ERROR] バイナリファイルが見つかりません: %BINARY_PATH%
    echo [ERROR] このスクリプトは、バイナリファイルと同じディレクトリで実行してください。
    exit /b 1
)

echo [INFO] HEIC Image Converter のインストールを開始します...
echo [INFO] インストール先: %INSTALL_DIR%

REM インストール先ディレクトリの作成
if not exist "%INSTALL_DIR%" (
    echo [INFO] インストール先ディレクトリを作成します...
    mkdir "%INSTALL_DIR%"
)

REM 既存バイナリの確認
set "BINARY_DEST=%INSTALL_DIR%\convert.exe"
if exist "%BINARY_DEST%" (
    echo [WARN] 既存のバイナリが見つかりました。上書きして更新します。
)

REM バイナリのコピー
echo [INFO] バイナリをコピーしています...
copy /Y "%BINARY_PATH%" "%BINARY_DEST%"

REM MinGWランタイムDLLのコピー
echo [INFO] 必要なDLLファイルをコピーしています...
set "DLL_COUNT=0"
if exist "%SCRIPT_DIR%libgcc_s_seh-1.dll" (
    copy /Y "%SCRIPT_DIR%libgcc_s_seh-1.dll" "%INSTALL_DIR%\" >nul 2>&1
    set /a DLL_COUNT+=1
)
if exist "%SCRIPT_DIR%libwinpthread-1.dll" (
    copy /Y "%SCRIPT_DIR%libwinpthread-1.dll" "%INSTALL_DIR%\" >nul 2>&1
    set /a DLL_COUNT+=1
)
if exist "%SCRIPT_DIR%libstdc++-6.dll" (
    copy /Y "%SCRIPT_DIR%libstdc++-6.dll" "%INSTALL_DIR%\" >nul 2>&1
    set /a DLL_COUNT+=1
)
if !DLL_COUNT! gtr 0 (
    echo [INFO] !DLL_COUNT! 個のDLLファイルをコピーしました。
) else (
    echo [WARN] DLLファイルが見つかりませんでした。バイナリが正常に動作しない可能性があります。
)

REM アンインストールスクリプトのコピー
set "UNINSTALL_PS1_SOURCE=%SCRIPT_DIR%%UNINSTALL_SCRIPT_PS1%"
set "UNINSTALL_BAT_SOURCE=%SCRIPT_DIR%%UNINSTALL_SCRIPT_BAT%"

if exist "%UNINSTALL_PS1_SOURCE%" (
    echo [INFO] アンインストールスクリプトをコピーしています...
    copy /Y "%UNINSTALL_PS1_SOURCE%" "%INSTALL_DIR%\%UNINSTALL_SCRIPT_PS1%"
)

if exist "%UNINSTALL_BAT_SOURCE%" (
    copy /Y "%UNINSTALL_BAT_SOURCE%" "%INSTALL_DIR%\%UNINSTALL_SCRIPT_BAT%"
)

REM PATH設定の確認
echo.
echo [INFO] PATH設定について
echo [INFO] インストール先 (%INSTALL_DIR%) をPATHに追加すると、どこからでも 'convert' コマンドを実行できます。

REM 現在のユーザー環境変数のPATHを取得
for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "CURRENT_PATH=%%b"

REM 既にPATHに追加されているか確認
set "PATH_ADDED=0"
if defined CURRENT_PATH (
    echo !CURRENT_PATH! | findstr /C:"%INSTALL_DIR%" >nul
    if !errorlevel! equ 0 (
        echo [INFO] 既にPATHに追加されています。
        set "PATH_ADDED=1"
    )
)

if "!PATH_ADDED!"=="0" (
    set /p RESPONSE="PATHに追加しますか？ (Y/n): "
    if /i "!RESPONSE!"=="" set RESPONSE=Y
    if /i "!RESPONSE!"=="Y" (
        echo [INFO] PATH設定を追加しています...
        if defined CURRENT_PATH (
            set "NEW_PATH=!CURRENT_PATH!;%INSTALL_DIR%"
        ) else (
            set "NEW_PATH=%INSTALL_DIR%"
        )
        reg add "HKCU\Environment" /v Path /t REG_EXPAND_SZ /d "!NEW_PATH!" /f >nul 2>&1
        if !errorlevel! equ 0 (
            echo [INFO] PATH設定を追加しました。
            echo [INFO] 新しいコマンドプロンプトまたはPowerShellを開いてください。
            set "PATH_ADDED=1"
        ) else (
            echo [ERROR] PATH設定の追加に失敗しました。
            echo [INFO] 手動でPATHに追加する場合は、環境変数の設定から以下を追加してください:
            echo [INFO]   %INSTALL_DIR%
            set "PATH_ADDED=0"
        )
    ) else (
        echo [INFO] PATH設定をスキップしました。
        echo [INFO] 手動でPATHに追加する場合は、環境変数の設定から以下を追加してください:
        echo [INFO]   %INSTALL_DIR%
        set "PATH_ADDED=0"
    )
)

echo.
echo [INFO] インストールが完了しました！
echo.
echo [INFO] 使用方法:
if "!PATH_ADDED!"=="1" (
    echo   convert --help
) else (
    echo   %BINARY_DEST% --help
)
echo.
echo [INFO] アンインストール方法:
echo   convert --uninstall
echo   または
echo   %INSTALL_DIR%\%UNINSTALL_SCRIPT_BAT%

endlocal

