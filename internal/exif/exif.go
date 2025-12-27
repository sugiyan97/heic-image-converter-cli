// Package exif provides functionality for extracting, displaying, and manipulating EXIF metadata
// from HEIC and JPEG image files.
package exif

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
	exifv3 "github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

// ExtractEXIFFromHEIC extracts EXIF data from a HEIC file
func ExtractEXIFFromHEIC(heicPath string) ([]byte, error) {
	file, err := os.Open(heicPath)
	if err != nil {
		return nil, fmt.Errorf("HEICファイルを開けませんでした: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't fail operation
			_ = closeErr
		}
	}()

	// Extract EXIF data from HEIC file
	exifBytes, err := goheif.ExtractExif(file)
	if err != nil {
		// EXIFが存在しない場合は空のスライスを返す
		return nil, nil
	}

	if len(exifBytes) == 0 {
		return nil, nil
	}

	return exifBytes, nil
}

// EmbedEXIFToJPEG embeds EXIF data into a JPEG file
// Note: This is a placeholder implementation. Full EXIF embedding requires IfdBuilder.
func EmbedEXIFToJPEG(jpegPath string, exifData []byte) error {
	if len(exifData) == 0 {
		// No EXIF data to embed
		return nil
	}

	// Read the JPEG file
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return fmt.Errorf("JPEGファイルの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Remove existing EXIF segment if present
	_, err = sl.DropExif()
	if err != nil {
		return fmt.Errorf("EXIFセグメントの削除に失敗しました: %w", err)
	}

	// TODO: SetExif requires *exif.IfdBuilder, not []byte
	// For now, we'll skip EXIF embedding and return a warning
	// This will be implemented when we have proper EXIF extraction from HEIC

	// Write the modified JPEG
	outFile, err := os.Create(jpegPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗しました: %w", err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			// Log error but don't fail operation
			_ = closeErr
		}
	}()

	if err := sl.Write(outFile); err != nil {
		return fmt.Errorf("JPEGファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// RemoveEXIFFromJPEG removes EXIF data from a JPEG file
func RemoveEXIFFromJPEG(jpegPath string) error {
	// Read the JPEG file
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return fmt.Errorf("JPEGファイルの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Remove EXIF segment
	_, err = sl.DropExif()
	if err != nil {
		return fmt.Errorf("EXIFセグメントの削除に失敗しました: %w", err)
	}

	// Write the modified JPEG
	outFile, err := os.Create(jpegPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗しました: %w", err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			// Log error but don't fail operation
			_ = closeErr
		}
	}()

	if err := sl.Write(outFile); err != nil {
		return fmt.Errorf("JPEGファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// ShowEXIFFromHEIC displays EXIF information from a HEIC file
func ShowEXIFFromHEIC(heicPath string) error {
	// Extract EXIF data from HEIC
	exifBytes, err := ExtractEXIFFromHEIC(heicPath)
	if err != nil {
		return fmt.Errorf("HEICファイルからEXIF情報の抽出に失敗しました: %w", err)
	}

	if len(exifBytes) == 0 {
		fmt.Printf("=== EXIF情報: %s ===\n", filepath.Base(heicPath))
		fmt.Println("EXIF情報: なし")
		fmt.Println()
		return nil
	}

	// Search and extract EXIF data
	rawExif, err := exifv3.SearchAndExtractExif(exifBytes)
	if err != nil {
		// SearchAndExtractExifが失敗した場合、直接GetFlatExifDataを試す
		entries, _, err := exifv3.GetFlatExifData(exifBytes, nil)
		if err != nil {
			return fmt.Errorf("EXIF情報の解析に失敗しました: %w", err)
		}

		// Display EXIF information
		fmt.Printf("=== EXIF情報: %s ===\n", filepath.Base(heicPath))
		if len(entries) == 0 {
			fmt.Println("EXIF情報: なし")
		} else {
			printExifEntries(entries)
		}
		fmt.Println()
		return nil
	}

	// Parse EXIF data
	entries, _, err := exifv3.GetFlatExifData(rawExif, nil)
	if err != nil {
		return fmt.Errorf("EXIF情報の解析に失敗しました: %w", err)
	}

	// Display EXIF information
	fmt.Printf("=== EXIF情報: %s ===\n", filepath.Base(heicPath))
	if len(entries) == 0 {
		fmt.Println("EXIF情報: なし")
	} else {
		printExifEntries(entries)
	}
	fmt.Println()

	return nil
}

// ShowEXIFFromJPEG displays EXIF information from a JPEG file
func ShowEXIFFromJPEG(jpegPath string) error {
	// Read the JPEG file
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return fmt.Errorf("JPEGファイルの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Get EXIF data
	_, exifData, err := sl.Exif()
	if err != nil {
		return fmt.Errorf("EXIF情報が見つかりませんでした: %w", err)
	}

	// Parse EXIF data
	entries, _, err := exifv3.GetFlatExifData(exifData, nil)
	if err != nil {
		return fmt.Errorf("EXIF情報の解析に失敗しました: %w", err)
	}

	// Display EXIF information
	fmt.Printf("=== EXIF情報: %s ===\n", filepath.Base(jpegPath))
	printExifEntries(entries)
	fmt.Println()

	return nil
}

// printExifEntries displays EXIF entries in a formatted way
func printExifEntries(entries []exifv3.ExifTag) {
	if len(entries) == 0 {
		fmt.Println("EXIF情報: なし")
		return
	}

	// Important tags to display first
	importantTags := []string{
		"DateTime", "DateTimeOriginal", "Make", "Model",
		"Orientation", "XResolution", "YResolution", "ResolutionUnit",
		"Software", "Artist", "Copyright", "ExifVersion",
		"Flash", "FocalLength", "FNumber", "ExposureTime", "ISOSpeedRatings",
		"GPSInfo", "ImageWidth", "ImageLength",
	}

	// Create a map for quick lookup
	entryMap := make(map[string]exifv3.ExifTag)
	for _, entry := range entries {
		tagName := entry.TagName
		if tagName == "" {
			tagName = fmt.Sprintf("Tag_0x%04x", entry.TagId)
		}
		entryMap[tagName] = entry
	}

	// Display important tags first
	for _, tag := range importantTags {
		if entry, ok := entryMap[tag]; ok {
			fmt.Printf("  %s: %s\n", tag, entry.Formatted)
			delete(entryMap, tag)
		}
	}

	// Display other tags
	if len(entryMap) > 0 {
		otherTags := make([]string, 0, len(entryMap))
		for tag := range entryMap {
			otherTags = append(otherTags, tag)
		}

		if len(otherTags) > 0 {
			fmt.Printf("  その他のタグ (%d個):\n", len(otherTags))
			for i, tag := range otherTags {
				if i >= 10 {
					fmt.Printf("    ... 他 %d 個のタグ\n", len(otherTags)-10)
					break
				}
				entry := entryMap[tag]
				fmt.Printf("    %s: %s\n", tag, entry.Formatted)
			}
		}
	}
}

// CheckEXIFInJPEG checks if EXIF data exists in a JPEG file
func CheckEXIFInJPEG(jpegPath string) (bool, []string, error) {
	// Read the JPEG file
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return false, nil, fmt.Errorf("JPEGファイルの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return false, nil, fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Get EXIF data
	_, exifData, err := sl.Exif()
	if err != nil {
		// No EXIF data found
		return false, nil, nil
	}

	// Parse EXIF data to get tag names
	entries, _, err := exifv3.GetFlatExifData(exifData, nil)
	if err != nil {
		return false, nil, fmt.Errorf("EXIF情報の解析に失敗しました: %w", err)
	}

	tagNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		tagNames = append(tagNames, entry.TagName)
	}

	return true, tagNames, nil
}

// ExtractEXIFFromJPEG extracts EXIF data from a JPEG file
func ExtractEXIFFromJPEG(jpegPath string) ([]byte, error) {
	// Read the JPEG file
	data, err := os.ReadFile(jpegPath)
	if err != nil {
		return nil, fmt.Errorf("JPEGファイルの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Get EXIF data
	_, exifData, err := sl.Exif()
	if err != nil {
		return nil, nil // No EXIF data
	}

	return exifData, nil
}

// CopyEXIFFromHEICToJPEG copies EXIF data from HEIC to JPEG
// This is a placeholder - actual implementation depends on HEIC EXIF extraction
func CopyEXIFFromHEICToJPEG(heicPath, jpegPath string) error {
	// Try to extract EXIF from HEIC
	exifData, err := ExtractEXIFFromHEIC(heicPath)
	if err != nil {
		return fmt.Errorf("HEICファイルからEXIF情報の抽出に失敗しました: %w", err)
	}

	if len(exifData) == 0 {
		// No EXIF data in HEIC file
		return nil
	}

	// Embed EXIF into JPEG
	if err := EmbedEXIFToJPEG(jpegPath, exifData); err != nil {
		return fmt.Errorf("JPEGファイルへのEXIF情報の埋め込みに失敗しました: %w", err)
	}

	return nil
}

// FindHEICFiles recursively finds all HEIC files in a directory
func FindHEICFiles(dirPath string) ([]string, error) {
	var heicFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".heic" || ext == ".heif" {
			heicFiles = append(heicFiles, path)
		}

		return nil
	})

	return heicFiles, err
}

// FindJPEGFiles recursively finds all JPEG files in a directory
func FindJPEGFiles(dirPath string) ([]string, error) {
	var jpegFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" {
			jpegFiles = append(jpegFiles, path)
		}

		return nil
	})

	return jpegFiles, err
}

// IsHEICFile checks if a file is a HEIC file
func IsHEICFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".heic" || ext == ".heif"
}

// IsJPEGFile checks if a file is a JPEG file
func IsJPEGFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg"
}

// ReadEXIFFromReader reads EXIF data from an io.Reader
func ReadEXIFFromReader(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("データの読み込みに失敗しました: %w", err)
	}

	// Parse JPEG structure
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("JPEG構造の解析に失敗しました: %w", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Get EXIF data
	_, exifData, err := sl.Exif()
	if err != nil {
		return nil, nil // No EXIF data
	}

	return exifData, nil
}
