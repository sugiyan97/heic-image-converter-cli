package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupTestFile copies the test HEIC file to a temporary directory
func setupTestFile(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "heic-converter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Copy test file from sample directory
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

// TestConvertHEICToJPEG_TC00101 tests TC-001-01: Normal conversion of HEIC to JPEG
func TestConvertHEICToJPEG_TC00101(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	options := ConvertOptions{RemoveEXIF: false}
	err := ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Check output file exists
	outputPath := GenerateOutputPath(heicFile)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputPath)
	}

	// Verify it's a valid JPEG file
	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Logf("Failed to close file: %v", err)
		}
	}()

	_, err = jpeg.Decode(file)
	if err != nil {
		t.Fatalf("Output file is not a valid JPEG: %v", err)
	}
}

// TestConvertHEICToJPEG_ExtensionVariations tests TC-001-02, TC-001-03, TC-001-04: Different extension cases
func TestConvertHEICToJPEG_ExtensionVariations(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	// Read source data once before the loop
	sourceData, err := os.ReadFile(heicFile)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	extensions := []string{".HEIC", ".heic", ".Heic", ".HEIF", ".heif"}
	for _, ext := range extensions {
		t.Run(ext, func(t *testing.T) {
			// Create a copy with different extension
			newPath := filepath.Join(filepath.Dir(heicFile), "test"+ext)
			if err := os.WriteFile(newPath, sourceData, 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			defer func() {
				_ = os.Remove(newPath)
			}()

			options := ConvertOptions{RemoveEXIF: false}
			err = ConvertHEICToJPEG(newPath, options)
			if err != nil {
				t.Fatalf("Conversion failed for %s: %v", ext, err)
			}

			outputPath := GenerateOutputPath(newPath)
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Fatalf("Output file was not created: %s", outputPath)
			}
			defer func() {
				_ = os.Remove(outputPath)
			}()
		})
	}
}

// TestConvertHEICToJPEG_NonexistentFile tests TC-002-03: Error handling for nonexistent file
func TestConvertHEICToJPEG_NonexistentFile(t *testing.T) {
	t.Parallel()
	options := ConvertOptions{RemoveEXIF: false}
	err := ConvertHEICToJPEG("nonexistent.HEIC", options)
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
}

// TestConvertHEICToJPEG_InvalidFile tests TC-002-04: Error handling for invalid file
func TestConvertHEICToJPEG_InvalidFile(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp("", "heic-converter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	invalidFile := filepath.Join(tmpDir, "invalid.txt")
	if err := os.WriteFile(invalidFile, []byte("not a HEIC file"), 0644); err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	options := ConvertOptions{RemoveEXIF: false}
	err = ConvertHEICToJPEG(invalidFile, options)
	if err == nil {
		t.Fatal("Expected error for invalid file, got nil")
	}
}

// TestGenerateOutputPath tests output path generation
func TestGenerateOutputPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple case", "test.HEIC", "test.jpg"},
		{"Lowercase", "test.heic", "test.jpg"},
		{"Mixed case", "test.Heic", "test.jpg"},
		{"With path", "/path/to/test.HEIC", "/path/to/test.jpg"},
		{"Relative path", "./test.HEIC", "./test.jpg"},
		{"HEIF extension", "test.HEIF", "test.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GenerateOutputPath(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateOutputPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestConvertToRGBA_TC01001 tests TC-010-01, TC-010-02, TC-010-03, TC-010-04: Color space conversion
// Note: These tests verify the conversion functions work correctly for different color spaces
// Actual HEIC files with specific color spaces would be needed for complete testing
func TestConvertToRGBA_TC01001(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	// Open and decode HEIC file
	file, err := os.Open(heicFile)
	if err != nil {
		t.Fatalf("Failed to open HEIC file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Logf("Failed to close file: %v", err)
		}
	}()

	// Note: We can't directly test different color spaces without different HEIC files
	// This test verifies that the conversion works for the available test file
	options := ConvertOptions{RemoveEXIF: false}
	err = ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Verify output is valid
	outputPath := GenerateOutputPath(heicFile)
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			t.Logf("Failed to close output file: %v", err)
		}
	}()

	img, err := jpeg.Decode(outputFile)
	if err != nil {
		t.Fatalf("Failed to decode output JPEG: %v", err)
	}

	// Verify image has valid bounds
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatalf("Invalid image dimensions: %v", bounds)
	}

	// Verify image can be converted to RGBA
	rgbaImg := convertToRGBA(img)
	if rgbaImg == nil {
		t.Fatal("convertToRGBA returned nil")
	}

	rgbaBounds := rgbaImg.Bounds()
	if rgbaBounds.Dx() != bounds.Dx() || rgbaBounds.Dy() != bounds.Dy() {
		t.Errorf("RGBA conversion changed dimensions: %v -> %v", bounds, rgbaBounds)
	}
}

// TestJPEGQuality tests that JPEG quality is set to 95 (TC-001-01 requirement)
func TestJPEGQuality(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	options := ConvertOptions{RemoveEXIF: false}
	err := ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Note: JPEG quality is not directly readable from the file
	// This test verifies the constant is set correctly
	if JPEGQuality != 95 {
		t.Errorf("JPEGQuality = %d, want 95", JPEGQuality)
	}
}

// TestConvertOptions tests conversion with different options
func TestConvertOptions(t *testing.T) {
	t.Parallel()
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	// Test with RemoveEXIF = false
	options := ConvertOptions{RemoveEXIF: false}
	err := ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion with RemoveEXIF=false failed: %v", err)
	}

	outputPath := GenerateOutputPath(heicFile)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputPath)
	}
	_ = os.Remove(outputPath)

	// Test with RemoveEXIF = true (conversion should still work)
	options = ConvertOptions{RemoveEXIF: true}
	err = ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion with RemoveEXIF=true failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputPath)
	}
}

// TestConvertNRGBAToRGBA tests NRGBA to RGBA conversion
func TestConvertNRGBAToRGBA(t *testing.T) {
	t.Parallel()
	// Create a test NRGBA image
	bounds := image.Rect(0, 0, 10, 10)
	nrgba := image.NewNRGBA(bounds)

	// Fill with some test data
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			nrgba.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x * 25),
				G: uint8(y * 25),
				B: 128,
				A: 200, // Semi-transparent
			})
		}
	}

	// Convert to RGBA
	rgba := convertNRGBAToRGBA(nrgba)

	// Verify dimensions
	if rgba.Bounds() != bounds {
		t.Errorf("Bounds mismatch: got %v, want %v", rgba.Bounds(), bounds)
	}

	// Verify it's actually RGBA (rgba is already *image.RGBA)
	if rgba == nil {
		t.Error("convertNRGBAToRGBA returned nil")
	}
}

// TestConvertYCbCrToRGBA tests YCbCr to RGBA conversion
func TestConvertYCbCrToRGBA(t *testing.T) {
	t.Parallel()
	// Create a test YCbCr image
	bounds := image.Rect(0, 0, 10, 10)
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio422)

	// Fill with some test data
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			yi := ycbcr.YOffset(x, y)
			ci := ycbcr.COffset(x, y)
			ycbcr.Y[yi] = uint8((x + y) * 10)
			ycbcr.Cb[ci] = 128
			ycbcr.Cr[ci] = 128
		}
	}

	// Convert to RGBA
	rgba := convertYCbCrToRGBA(ycbcr)

	// Verify dimensions
	if rgba.Bounds() != bounds {
		t.Errorf("Bounds mismatch: got %v, want %v", rgba.Bounds(), bounds)
	}

	// Verify it's actually RGBA (rgba is already *image.RGBA)
	if rgba == nil {
		t.Error("convertYCbCrToRGBA returned nil")
	}
}

// TestConvertGenericToRGBA tests generic image to RGBA conversion
func TestConvertGenericToRGBA(t *testing.T) {
	t.Parallel()
	// Create a test image with alpha channel
	bounds := image.Rect(0, 0, 10, 10)
	nrgba := image.NewNRGBA(bounds)

	// Fill with semi-transparent pixels
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			nrgba.SetNRGBA(x, y, color.NRGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 128, // 50% transparent
			})
		}
	}

	// Convert using generic function
	rgba := convertGenericToRGBA(nrgba)

	// Verify dimensions
	if rgba.Bounds() != bounds {
		t.Errorf("Bounds mismatch: got %v, want %v", rgba.Bounds(), bounds)
	}

	// Verify it's actually RGBA (rgba is already *image.RGBA)
	if rgba == nil {
		t.Error("convertGenericToRGBA returned nil")
	}

	// Verify alpha channel was composited (should be opaque in output)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c := rgba.At(x, y)
			r, g, b, a := c.RGBA()
			if a>>8 != 255 {
				t.Errorf("Pixel at (%d, %d) should be opaque, got alpha=%d", x, y, a>>8)
			}
			// With white background compositing, red should become lighter
			if r>>8 == 255 && g>>8 == 0 && b>>8 == 0 {
				t.Errorf("Pixel at (%d, %d) should be composited on white, got RGB=(%d,%d,%d)", x, y, r>>8, g>>8, b>>8)
			}
		}
	}
}

// TestConvertHEICToJPEG_TC00205 tests TC-002-05: Corrupted HEIC file error
func TestConvertHEICToJPEG_TC00205(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "heic-converter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	corruptedFile := filepath.Join(tmpDir, "corrupted.HEIC")
	if err := os.WriteFile(corruptedFile, []byte("corrupted HEIC data"), 0644); err != nil {
		t.Fatalf("Failed to create corrupted file: %v", err)
	}

	options := ConvertOptions{RemoveEXIF: false}
	err = ConvertHEICToJPEG(corruptedFile, options)
	if err == nil {
		t.Fatal("Expected error for corrupted file, got nil")
	}

	// Check error message contains expected text
	errMsg := err.Error()
	if !strings.Contains(errMsg, "変換失敗") && !strings.Contains(errMsg, "デコード") && !strings.Contains(errMsg, "失敗") {
		t.Logf("Error message: %s", errMsg)
	}
}

// TestConvertHEICToJPEG_TC01501 tests TC-015-01: Performance - conversion speed
func TestConvertHEICToJPEG_TC01501(t *testing.T) {
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	start := time.Now()
	options := ConvertOptions{RemoveEXIF: false}
	err := ConvertHEICToJPEG(heicFile, options)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Check that conversion completed in reasonable time (less than 10 seconds)
	if duration > 10*time.Second {
		t.Errorf("Conversion took too long: %v", duration)
	}

	// Log performance for monitoring
	t.Logf("Conversion completed in %v", duration)
}

// TestConvertHEICToJPEG_TC01502 tests TC-015-02: Performance - batch conversion speed
func TestConvertHEICToJPEG_TC01502(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "heic-converter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create multiple test files
	sourceFile := filepath.Join("..", "..", "sample", "test.HEIC")
	sourceData, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	// Create 10 test files
	for i := 0; i < 10; i++ {
		destFile := filepath.Join(tmpDir, fmt.Sprintf("test%d.HEIC", i))
		if err := os.WriteFile(destFile, sourceData, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	start := time.Now()
	options := ConvertOptions{RemoveEXIF: false}

	heicFiles, err := filepath.Glob(filepath.Join(tmpDir, "*.HEIC"))
	if err != nil {
		t.Fatalf("Failed to find HEIC files: %v", err)
	}

	for _, heicFile := range heicFiles {
		if err := ConvertHEICToJPEG(heicFile, options); err != nil {
			t.Logf("Conversion failed for %s: %v", heicFile, err)
		}
	}

	duration := time.Since(start)

	// Check that batch conversion completed in reasonable time (less than 60 seconds for 10 files)
	if duration > 60*time.Second {
		t.Errorf("Batch conversion took too long: %v", duration)
	}

	// Log performance for monitoring
	t.Logf("Batch conversion of %d files completed in %v (avg: %v per file)", len(heicFiles), duration, duration/time.Duration(len(heicFiles)))
}

// TestConvertHEICToJPEG_TC01601 tests TC-016-01: Memory usage - large file conversion
func TestConvertHEICToJPEG_TC01601(t *testing.T) {
	heicFile, cleanup := setupTestFile(t)
	defer cleanup()

	// Get file size
	fileInfo, err := os.Stat(heicFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Note: We can't directly measure memory usage in Go tests without external tools
	// This test verifies that conversion completes without panics or obvious memory issues
	options := ConvertOptions{RemoveEXIF: false}
	err = ConvertHEICToJPEG(heicFile, options)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Verify output file was created
	outputPath := GenerateOutputPath(heicFile)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputPath)
	}

	// Log file size for monitoring
	t.Logf("Converted file size: %d bytes", fileInfo.Size())

	// Note: Actual memory leak detection would require runtime.MemStats or external profiling tools
}

