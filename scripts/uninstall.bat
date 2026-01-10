@echo off
REM HEIC Image Converter アンインストールスクリプト (Windows Batch)
REM UTF-8文字化け対策
chcp 65001 >nul 2>&1

setlocal enabledelayedexpansion

REM 固定インストール先
set "INSTALL_DIR=%USERPROFILE%\bin\HeicConverter"

REM インストール先の確認
if not exist "%INSTALL_DIR%" (
    echo [ERROR] インストール先が見つかりません: %INSTALL_DIR%
    echo [ERROR] 既にアンインストールされている可能性があります。
    exit /b 1
)

echo [INFO] HEIC Image Converter のアンインストールを開始します...
echo [INFO] 削除対象: %INSTALL_DIR%

REM 削除前の確認
set /p RESPONSE="本当に削除しますか？ (y/N): "
if /i not "!RESPONSE!"=="y" (
    echo [INFO] アンインストールをキャンセルしました。
    exit /b 0
)

REM 削除対象ファイルの確認
echo [INFO] 削除対象ファイルを確認しています...
set "DLL_COUNT=0"
if exist "%INSTALL_DIR%\libgcc_s_seh-1.dll" set /a DLL_COUNT+=1
if exist "%INSTALL_DIR%\libwinpthread-1.dll" set /a DLL_COUNT+=1
if exist "%INSTALL_DIR%\libstdc++-6.dll" set /a DLL_COUNT+=1
if exist "%INSTALL_DIR%\convert.exe" (
    echo [INFO] バイナリファイル: convert.exe
)
if !DLL_COUNT! gtr 0 (
    echo [INFO] DLLファイル: !DLL_COUNT! 個
)

REM 最終確認
echo [WARN] 警告: %INSTALL_DIR% フォルダ全体が削除されます。
set /p RESPONSE2="続行しますか？ (y/N): "
if /i not "!RESPONSE2!"=="y" (
    echo [INFO] アンインストールをキャンセルしました。
    exit /b 0
)

REM フォルダごと削除（バイナリ、DLL、スクリプトを含む）
echo [INFO] インストールフォルダを削除しています...
rd /s /q "%INSTALL_DIR%"
echo [INFO] インストールフォルダを削除しました。

REM PATH削除の確認
echo.
set /p PATH_RESPONSE="PATH設定も削除しますか？ (y/N): "
if /i "!PATH_RESPONSE!"=="y" (
    REM 現在のユーザー環境変数のPATHを取得
    for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "CURRENT_PATH=%%b"
    
    REM PATHに含まれているか確認
    echo !CURRENT_PATH! | findstr /C:"%INSTALL_DIR%" >nul
    if !errorlevel! equ 0 (
        echo [INFO] PATH設定を削除しています...
        REM PATHから該当パスを削除
        set "NEW_PATH=!CURRENT_PATH!"
        set "NEW_PATH=!NEW_PATH:%INSTALL_DIR%;=!"
        set "NEW_PATH=!NEW_PATH:;%INSTALL_DIR%=!"
        set "NEW_PATH=!NEW_PATH:%INSTALL_DIR%=!"
        reg add "HKCU\Environment" /v Path /t REG_EXPAND_SZ /d "!NEW_PATH!" /f >nul 2>&1
        echo [INFO] PATH設定を削除しました。
        echo [INFO] 新しいコマンドプロンプトまたはPowerShellを開いてください。
    ) else (
        echo [INFO] PATH設定が見つかりませんでした。
    )
) else (
    echo [INFO] PATH設定の削除をスキップしました。
)

echo.
echo [INFO] アンインストールが完了しました！

endlocal

