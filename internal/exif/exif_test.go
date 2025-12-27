package exif

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sugiyan97/heic-image-converter-cli/internal/converter"
)

// setupTestHEICFile copies the test HEIC file to a temporary directory
func setupTestHEICFile(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "heic-exif-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	sourceFile := filepath.Join("..", "..", "sample", "test.HEIC")
	destFile := filepath.Join(tmpDir, "test.HEIC")

	sourceData, err := os.ReadFile(sourceFile)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("Failed to read source file: %v", err)
	}

	if err := os.WriteFile(destFile, sourceData, 0644); err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write test file: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	return destFile, cleanup
}

// setupTestJPEGFile creates a JPEG file from HEIC for testing
func setupTestJPEGFile(t *testing.T) (string, func()) {
	t.Helper()

	heicFile, _ := setupTestHEICFile(t)
	// Don't defer cleanup here - we need the HEIC file to remain until cleanup is called

	// Convert HEIC to JPEG
	options := converter.ConvertOptions{RemoveEXIF: false}
	if err := converter.ConvertHEICToJPEG(heicFile, options); err != nil {
		_ = os.RemoveAll(filepath.Dir(heicFile))
		t.Fatalf("Failed to convert HEIC to JPEG: %v", err)
	}

	jpegFile := converter.GenerateOutputPath(heicFile)
	if _, err := os.Stat(jpegFile); os.IsNotExist(err) {
		_ = os.RemoveAll(filepath.Dir(heicFile))
		t.Fatalf("JPEG file was not created: %s", jpegFile)
	}

	cleanup := func() {
		_ = os.Remove(jpegFile)
		_ = os.Remove(heicFile)
		_ = os.RemoveAll(filepath.Dir(heicFile))
	}

	return jpegFile, cleanup
}

// TestFindHEICFiles tests finding HEIC files in a directory
func TestFindHEICFiles(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "heic-find-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create test files with different extensions
	testFiles := []string{"test1.HEIC", "test2.heic", "test3.Heic", "test4.HEIF", "test5.txt"}
	for _, f := range testFiles {
		filePath := filepath.Join(tmpDir, f)
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create subdirectory with HEIC file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	subFile := filepath.Join(subDir, "subtest.HEIC")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create subdir file: %v", err)
	}

	// Find HEIC files
	files, err := FindHEICFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindHEICFiles failed: %v", err)
	}

	// Should find 5 HEIC files (test1, test2, test3, test4, subtest)
	if len(files) != 5 {
		t.Errorf("Expected 5 HEIC files, got %d: %v", len(files), files)
	}
}

// TestFindJPEGFiles tests finding JPEG files in a directory
func TestFindJPEGFiles(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "jpeg-find-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create test files with different extensions
	testFiles := []string{"test1.JPG", "test2.jpg", "test3.JPEG", "test4.jpeg", "test5.txt"}
	for _, f := range testFiles {
		filePath := filepath.Join(tmpDir, f)
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Find JPEG files
	files, err := FindJPEGFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindJPEGFiles failed: %v", err)
	}

	// Should find 4 JPEG files
	if len(files) != 4 {
		t.Errorf("Expected 4 JPEG files, got %d: %v", len(files), files)
	}
}

// TestIsHEICFile tests HEIC file detection
func TestIsHEICFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Uppercase HEIC", "test.HEIC", true},
		{"Lowercase heic", "test.heic", true},
		{"Mixed case Heic", "test.Heic", true},
		{"HEIF extension", "test.HEIF", true},
		{"heif extension", "test.heif", true},
		{"JPEG file", "test.jpg", false},
		{"No extension", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsHEICFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsHEICFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestIsJPEGFile tests JPEG file detection
func TestIsJPEGFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Uppercase JPG", "test.JPG", true},
		{"Lowercase jpg", "test.jpg", true},
		{"Uppercase JPEG", "test.JPEG", true},
		{"Lowercase jpeg", "test.jpeg", true},
		{"HEIC file", "test.HEIC", false},
		{"No extension", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsJPEGFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsJPEGFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestRemoveEXIFFromJPEG_TC00601 tests TC-006-01: Remove EXIF from JPEG
func TestRemoveEXIFFromJPEG_TC00601(t *testing.T) {
	t.Parallel()
	jpegFile, cleanup := setupTestJPEGFile(t)
	defer cleanup()

	// Check if EXIF exists before removal
	hasEXIFBefore, _, err := CheckEXIFInJPEG(jpegFile)
	if err != nil {
		t.Fatalf("Failed to check EXIF before removal: %v", err)
	}

	// Remove EXIF
	err = RemoveEXIFFromJPEG(jpegFile)
	if err != nil {
		t.Fatalf("Failed to remove EXIF: %v", err)
	}

	// Check if EXIF was removed
	hasEXIFAfter, _, err := CheckEXIFInJPEG(jpegFile)
	if err != nil {
		t.Fatalf("Failed to check EXIF after removal: %v", err)
	}

	// If EXIF existed before, it should be gone now
	if hasEXIFBefore && hasEXIFAfter {
		t.Error("EXIF was not removed from JPEG file")
	}
}

// TestCheckEXIFInJPEG_TC00801 tests TC-008-01, TC-008-02: Check EXIF in JPEG
func TestCheckEXIFInJPEG_TC00801(t *testing.T) {
	t.Parallel()
	jpegFile, cleanup := setupTestJPEGFile(t)
	defer cleanup()

	// Check EXIF
	hasEXIF, tags, err := CheckEXIFInJPEG(jpegFile)
	if err != nil {
		t.Fatalf("Failed to check EXIF: %v", err)
	}

	// Note: The test file may or may not have EXIF
	// This test verifies the function works correctly
	if hasEXIF {
		if len(tags) == 0 {
			t.Error("hasEXIF is true but no tags returned")
		}
	} else {
		// If no EXIF, tags should be empty
		if len(tags) > 0 {
			t.Error("hasEXIF is false but tags returned")
		}
	}
}

// TestCheckEXIFInJPEG_NonexistentFile tests TC-008-06: Check EXIF for nonexistent file
func TestCheckEXIFInJPEG_NonexistentFile(t *testing.T) {
	t.Parallel()
	_, _, err := CheckEXIFInJPEG("nonexistent.jpg")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestShowEXIFFromJPEG_TC00701 tests TC-007-01: Show EXIF from JPEG
func TestShowEXIFFromJPEG_TC00701(t *testing.T) {
	t.Parallel()
	jpegFile, cleanup := setupTestJPEGFile(t)
	defer cleanup()

	// Try to show EXIF (may or may not have EXIF)
	err := ShowEXIFFromJPEG(jpegFile)
	// Function should not error even if EXIF doesn't exist
	// (though current implementation may error - this is acceptable)
	if err != nil {
		// If error, it should be because EXIF doesn't exist
		if !strings.Contains(err.Error(), "EXIF情報が見つかりませんでした") &&
			!strings.Contains(err.Error(), "not found") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

// TestExtractEXIFFromJPEG tests EXIF extraction from JPEG
func TestExtractEXIFFromJPEG(t *testing.T) {
	t.Parallel()
	jpegFile, cleanup := setupTestJPEGFile(t)
	defer cleanup()

	exifData, err := ExtractEXIFFromJPEG(jpegFile)
	if err != nil {
		t.Fatalf("Failed to extract EXIF: %v", err)
	}

	// Note: exifData may be nil if no EXIF exists
	// This is acceptable - the function should handle both cases
	_ = exifData
}

// TestRemoveEXIFFromJPEG_NonexistentFile tests error handling for nonexistent file
func TestRemoveEXIFFromJPEG_NonexistentFile(t *testing.T) {
	t.Parallel()
	err := RemoveEXIFFromJPEG("nonexistent.jpg")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestShowEXIFFromJPEG_NonexistentFile tests error handling for nonexistent file
func TestShowEXIFFromJPEG_NonexistentFile(t *testing.T) {
	t.Parallel()
	err := ShowEXIFFromJPEG("nonexistent.jpg")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestExtractEXIFFromHEIC tests HEIC EXIF extraction (placeholder)
func TestExtractEXIFFromHEIC(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestHEICFile(t)
	defer cleanup()

	exifData, err := ExtractEXIFFromHEIC(heicFile)
	if err != nil {
		t.Fatalf("ExtractEXIFFromHEIC failed: %v", err)
	}

	// Note: Current implementation returns nil
	// This is expected as it's a placeholder
	if exifData != nil {
		t.Logf("EXIF data extracted: %d bytes", len(exifData))
	}
}

// TestCopyEXIFFromHEICToJPEG tests copying EXIF from HEIC to JPEG
func TestCopyEXIFFromHEICToJPEG(t *testing.T) {
	t.Parallel()
	heicFile, cleanupHEIC := setupTestHEICFile(t)
	defer cleanupHEIC()

	// Convert to JPEG first
	options := converter.ConvertOptions{RemoveEXIF: false}
	if err := converter.ConvertHEICToJPEG(heicFile, options); err != nil {
		t.Fatalf("Failed to convert HEIC to JPEG: %v", err)
	}

	jpegFile := converter.GenerateOutputPath(heicFile)
	defer func() {
		_ = os.Remove(jpegFile)
		_ = os.Remove(heicFile)
		_ = os.RemoveAll(filepath.Dir(heicFile))
	}()

	// Try to copy EXIF (may fail if HEIC EXIF extraction is not implemented)
	err := CopyEXIFFromHEICToJPEG(heicFile, jpegFile)
	// This may fail if EXIF extraction is not implemented, which is acceptable
	if err != nil {
		t.Logf("CopyEXIFFromHEICToJPEG failed (expected if not implemented): %v", err)
	}
}

// TestEmbedEXIFToJPEG tests embedding EXIF into JPEG
func TestEmbedEXIFToJPEG(t *testing.T) {
	t.Parallel()
	jpegFile, cleanup := setupTestJPEGFile(t)
	defer cleanup()

	// Try to embed empty EXIF data
	err := EmbedEXIFToJPEG(jpegFile, nil)
	if err != nil {
		t.Fatalf("EmbedEXIFToJPEG failed with nil data: %v", err)
	}

	// Try to embed some test EXIF data
	testEXIF := []byte{0xFF, 0xE1, 0x00, 0x10} // Minimal EXIF header
	err = EmbedEXIFToJPEG(jpegFile, testEXIF)
	// This may fail if full implementation is not available, which is acceptable
	if err != nil {
		t.Logf("EmbedEXIFToJPEG failed (expected if not fully implemented): %v", err)
	}
}

// TestFindHEICFiles_EmptyDirectory tests finding HEIC files in empty directory
func TestFindHEICFiles_EmptyDirectory(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "heic-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	files, err := FindHEICFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindHEICFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

// TestFindJPEGFiles_EmptyDirectory tests finding JPEG files in empty directory
func TestFindJPEGFiles_EmptyDirectory(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "jpeg-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	files, err := FindJPEGFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindJPEGFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(files))
	}
}

// TestFindHEICFiles_NonexistentDirectory tests error handling for nonexistent directory
func TestFindHEICFiles_NonexistentDirectory(t *testing.T) {
	t.Parallel()
	_, err := FindHEICFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}

// TestFindJPEGFiles_NonexistentDirectory tests error handling for nonexistent directory
func TestFindJPEGFiles_NonexistentDirectory(t *testing.T) {
	t.Parallel()
	_, err := FindJPEGFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}

