// Package cli provides the command-line interface for the HEIC image converter.
// It handles command parsing, flag management, and orchestrates the conversion process.
package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/sugiyan97/heic-image-converter-cli/internal/converter"
	"github.com/sugiyan97/heic-image-converter-cli/internal/exif"
)

var (
	// Version is overwritten at build time via ldflags (default: "v0.0.0")
	Version     = "v0.0.0"
	showEXIF    bool
	removeEXIF  bool
	checkEXIF   bool
	uninstall   bool
	showVersion bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "heic-convert [ファイル/ディレクトリ]",
	Short: "HEIC画像をJPEG形式に変換する（現時点ではJPEG形式をサポート）",
	Long: `HEIC Image Converterは、HEIC形式の画像ファイルを他の画像形式に変換するコマンドラインツールです。
現時点ではJPEG形式への変換をサポートしています。

引数なしで実行した場合、カレントディレクトリ内の全HEICファイルを再帰的に検索して変換します。
ファイルパスまたはディレクトリパスを指定することで、特定のファイルやディレクトリを処理できます。`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConvert,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&showEXIF, "show-exif", false, "EXIF情報を表示します")
	rootCmd.Flags().BoolVar(&removeEXIF, "remove-exif", false, "EXIF情報を削除して変換します")
	rootCmd.Flags().BoolVar(&checkEXIF, "check-exif", false, "JPEGファイルのEXIF情報の有無をチェックします")
	rootCmd.Flags().BoolVar(&uninstall, "uninstall", false, "アンインストールを実行します")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "バージョンを表示します")
}

func runConvert(_ *cobra.Command, args []string) error {
	// バージョン表示モード
	if showVersion {
		fmt.Println(Version)
		return nil
	}

	// アンインストールモード
	if uninstall {
		return runUninstall()
	}

	// EXIFチェックモード
	if checkEXIF {
		return runCheckEXIF(args)
	}

	// EXIF表示モード（変換なし）
	if showEXIF && !removeEXIF {
		return runShowEXIF(args)
	}

	// 変換モード
	return runConvertMode(args)
}

func runCheckEXIF(args []string) error {
	var targetPath string
	if len(args) > 0 {
		targetPath = args[0]
	} else {
		targetPath = "."
	}

	// パスの存在確認
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("パスが見つかりません: %w", err)
	}

	var jpegFiles []string
	if info.IsDir() {
		// ディレクトリの場合、JPEGファイルを再帰的に検索
		files, err := exif.FindJPEGFiles(targetPath)
		if err != nil {
			return fmt.Errorf("JPEGファイルの検索に失敗しました: %w", err)
		}
		jpegFiles = files
	} else {
		// ファイルの場合
		if !exif.IsJPEGFile(targetPath) {
			return fmt.Errorf("指定されたファイルはJPEGファイルではありません: %s", targetPath)
		}
		jpegFiles = []string{targetPath}
	}

	// EXIFチェック
	var hasEXIFCount, noEXIFCount, errorCount int
	for _, jpegPath := range jpegFiles {
		hasEXIF, tags, err := exif.CheckEXIFInJPEG(jpegPath)
		if err != nil {
			fmt.Printf("✗ エラー: %s - %v\n", jpegPath, err)
			errorCount++
			continue
		}

		if hasEXIF {
			fmt.Printf("✗ EXIF情報が残っています: %s\n", jpegPath)
			if len(tags) > 0 {
				fmt.Printf("  検出された主要なEXIFタグ: %s\n", tags[0])
				if len(tags) > 1 {
					fmt.Printf("  (他 %d 個のタグ)\n", len(tags)-1)
				}
			}
			hasEXIFCount++
		} else {
			fmt.Printf("✓ EXIF情報は削除されています: %s\n", jpegPath)
			noEXIFCount++
		}
	}

	// サマリー表示
	fmt.Printf("\n=== チェック結果 ===\n")
	fmt.Printf("総ファイル数: %d\n", len(jpegFiles))
	fmt.Printf("EXIF削除済み: %d\n", noEXIFCount)
	fmt.Printf("EXIF残存: %d\n", hasEXIFCount)
	fmt.Printf("エラー: %d\n", errorCount)

	return nil
}

func runShowEXIF(args []string) error {
	var targetPath string
	if len(args) > 0 {
		targetPath = args[0]
	} else {
		targetPath = "."
	}

	// パスの存在確認
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("パスが見つかりません: %w", err)
	}

	var heicFiles []string
	if info.IsDir() {
		// ディレクトリの場合、HEICファイルを再帰的に検索
		files, err := exif.FindHEICFiles(targetPath)
		if err != nil {
			return fmt.Errorf("HEICファイルの検索に失敗しました: %w", err)
		}
		heicFiles = files
	} else {
		// ファイルの場合
		if !exif.IsHEICFile(targetPath) {
			return fmt.Errorf("指定されたファイルはHEICファイルではありません: %s", targetPath)
		}
		heicFiles = []string{targetPath}
	}

	if len(heicFiles) == 0 {
		fmt.Println("EXIF情報を表示するHEICファイルが見つかりませんでした。")
		return nil
	}

	// EXIF情報の表示
	var errorCount int
	for _, heicPath := range heicFiles {
		if err := exif.ShowEXIFFromHEIC(heicPath); err != nil {
			fmt.Printf("警告: %s のEXIF情報の表示に失敗しました: %v\n", heicPath, err)
			errorCount++
		}
	}

	// サマリー表示
	if len(heicFiles) > 1 {
		fmt.Printf("\n=== 表示結果 ===\n")
		fmt.Printf("表示成功: %d\n", len(heicFiles)-errorCount)
		fmt.Printf("表示失敗: %d\n", errorCount)
	}

	return nil
}

func runConvertMode(args []string) error {
	var targetPath string
	if len(args) > 0 {
		targetPath = args[0]
	} else {
		targetPath = "."
	}

	// パスの存在確認
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("パスが見つかりません: %w", err)
	}

	var heicFiles []string
	if info.IsDir() {
		// ディレクトリの場合、HEICファイルを再帰的に検索
		files, err := exif.FindHEICFiles(targetPath)
		if err != nil {
			return fmt.Errorf("HEICファイルの検索に失敗しました: %w", err)
		}
		heicFiles = files
	} else {
		// ファイルの場合
		if !exif.IsHEICFile(targetPath) {
			return fmt.Errorf("指定されたファイルはHEICファイルではありません: %s", targetPath)
		}
		heicFiles = []string{targetPath}
	}

	if len(heicFiles) == 0 {
		fmt.Println("変換対象のHEICファイルが見つかりませんでした。")
		return nil
	}

	// 変換オプション
	options := converter.ConvertOptions{
		RemoveEXIF: removeEXIF,
	}

	// 変換処理
	var successCount, errorCount int
	for _, heicPath := range heicFiles {
		// EXIF情報の表示（変換前にHEICファイルから表示）
		if showEXIF {
			if err := exif.ShowEXIFFromHEIC(heicPath); err != nil {
				fmt.Printf("警告: %s のEXIF情報の表示に失敗しました: %v\n", heicPath, err)
			}
		}

		// HEIC→JPEG変換
		if err := converter.ConvertHEICToJPEG(heicPath, options); err != nil {
			fmt.Printf("✗ 変換失敗: %s - %v\n", heicPath, err)
			errorCount++
			continue
		}

		// 出力ファイルパス
		outputPath := converter.GenerateOutputPath(heicPath)

		// EXIF情報の処理
		if removeEXIF {
			// EXIF情報を削除
			if err := exif.RemoveEXIFFromJPEG(outputPath); err != nil {
				fmt.Printf("警告: %s のEXIF情報の削除に失敗しました: %v\n", outputPath, err)
			}
		}
		// EXIF情報を保持（HEICからJPEGへコピー）
		// 現時点では、goheifがEXIF抽出を直接サポートしていないため、
		// この機能は将来的に実装予定
		// if err := exif.CopyEXIFFromHEICToJPEG(heicPath, outputPath); err != nil {
		// 	fmt.Printf("警告: %s のEXIF情報の保持に失敗しました: %v\n", outputPath, err)
		// }

		fmt.Printf("✓ 変換完了: %s -> %s\n", heicPath, outputPath)
		successCount++
	}

	// サマリー表示
	if len(heicFiles) > 1 {
		fmt.Printf("\n=== 変換結果 ===\n")
		fmt.Printf("変換成功: %d\n", successCount)
		fmt.Printf("変換失敗: %d\n", errorCount)
	}

	return nil
}

func runUninstall() error {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("ホームディレクトリを取得できませんでした: %w", err)
	}

	// 固定インストール先
	var installDir string
	var uninstallScript string
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		installDir = filepath.Join(homeDir, "bin", "HeicConverter")
		uninstallScript = filepath.Join(installDir, "uninstall.ps1")
		
		// PowerShellスクリプトが存在するか確認
		if _, err := os.Stat(uninstallScript); os.IsNotExist(err) {
			// バッチファイルを試す
			uninstallScript = filepath.Join(installDir, "uninstall.bat")
			if _, err := os.Stat(uninstallScript); os.IsNotExist(err) {
				return fmt.Errorf("アンインストールスクリプトが見つかりません: %s", installDir)
			}
			cmd = exec.Command("cmd", "/c", uninstallScript)
		} else {
			cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", uninstallScript)
		}
	} else {
		installDir = filepath.Join(homeDir, "bin", "HeicConverter")
		uninstallScript = filepath.Join(installDir, "uninstall.sh")
		
		if _, err := os.Stat(uninstallScript); os.IsNotExist(err) {
			return fmt.Errorf("アンインストールスクリプトが見つかりません: %s", installDir)
		}
		cmd = exec.Command("bash", uninstallScript)
	}

	// アンインストールスクリプトを実行
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("アンインストールスクリプトを実行します: %s\n", uninstallScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("アンインストールスクリプトの実行に失敗しました: %w", err)
	}

	return nil
}
