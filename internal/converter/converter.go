// Package converter provides functionality for converting HEIC image files to JPEG format.
package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
)

const (
	// JPEGQuality is the quality setting for JPEG encoding (0-100)
	JPEGQuality = 95
)

func init() {
	// goheif's default decode path hands back Y/Cb/Cr slices that alias
	// the underlying C decoder's buffer, which is freed as soon as
	// goheif.Decode returns (via its internal defer dec.Free()). Any
	// pixel access afterwards -- including jpeg.Encode's direct
	// image.YCbCr fast path -- is a use-after-free that segfaults
	// intermittently depending on whether the freed memory has been
	// reused yet. SafeEncoding makes goheif copy the buffer into
	// Go-managed memory (via C.GoBytes) before freeing it.
	goheif.SafeEncoding = true
}

// ConvertOptions holds options for HEIC to JPEG conversion
type ConvertOptions struct {
	RemoveEXIF bool
}

// ConvertHEICToJPEG converts a HEIC file to JPEG format
func ConvertHEICToJPEG(inputPath string, _ ConvertOptions) error {
	// Open HEIC file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't fail conversion
			_ = closeErr
		}
	}()

	// Decode HEIC image
	img, err := goheif.Decode(file)
	if err != nil {
		return fmt.Errorf("HEICファイルのデコードに失敗しました: %w", err)
	}

	// Generate output file path
	outputPath := GenerateOutputPath(inputPath)

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("出力ファイルを作成できませんでした: %w", err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			// Log error but don't fail conversion
			_ = closeErr
		}
	}()

	// jpeg.Encode has a fast path for *image.YCbCr and *image.Gray that writes
	// the image directly without per-pixel color conversion. goheif.Decode
	// always returns *image.YCbCr, so pass it straight through in that case
	// and only fall back to an RGBA conversion for other color models (e.g.
	// ones with an alpha channel that needs to be composited away).
	encodeImg := img
	switch img.(type) {
	case *image.YCbCr, *image.Gray:
		// Already directly encodable by jpeg.Encode; no conversion needed.
	default:
		encodeImg = convertToRGBA(img)
	}

	// Encode as JPEG
	opts := &jpeg.Options{Quality: JPEGQuality}
	if err := jpeg.Encode(outFile, encodeImg, opts); err != nil {
		return fmt.Errorf("JPEGファイルのエンコードに失敗しました: %w", err)
	}

	return nil
}

// convertToRGBA converts an image to RGBA format.
// Handles color spaces that jpeg.Encode cannot write directly (RGBA, NRGBA,
// and other generic image.Image implementations), notably ones with an
// alpha channel that needs to be composited away.
func convertToRGBA(img image.Image) image.Image {
	switch src := img.(type) {
	case *image.RGBA:
		// Already RGBA, return as is
		return src
	case *image.NRGBA:
		// Convert NRGBA to RGBA
		return convertNRGBAToRGBA(src)
	default:
		// Generic conversion for other types
		return convertGenericToRGBA(img)
	}
}

// convertNRGBAToRGBA converts NRGBA to RGBA
func convertNRGBAToRGBA(src *image.NRGBA) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcIdx := src.PixOffset(x, y)
			dstIdx := dst.PixOffset(x, y)

			r := uint32(src.Pix[srcIdx+0])
			g := uint32(src.Pix[srcIdx+1])
			b := uint32(src.Pix[srcIdx+2])
			a := uint32(src.Pix[srcIdx+3])

			// Premultiply alpha
			if a < 255 {
				r = r * a / 255
				g = g * a / 255
				b = b * a / 255
			}

			// Composite on white background if alpha < 255
			if a < 255 {
				alpha := 255 - a
				r = r + alpha
				g = g + alpha
				b = b + alpha
			}

			dst.Pix[dstIdx+0] = uint8(r)
			dst.Pix[dstIdx+1] = uint8(g)
			dst.Pix[dstIdx+2] = uint8(b)
			dst.Pix[dstIdx+3] = 255
		}
	}

	return dst
}

// convertGenericToRGBA converts any image type to RGBA
// Handles alpha channel by compositing on white background
func convertGenericToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Scale from 16-bit to 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Composite on white background if alpha < 255
			if a8 < 255 {
				alpha := float64(a8) / 255.0
				white := 255.0
				r8 = uint8(float64(r8)*alpha + white*(1.0-alpha))
				g8 = uint8(float64(g8)*alpha + white*(1.0-alpha))
				b8 = uint8(float64(b8)*alpha + white*(1.0-alpha))
			}

			dst.SetRGBA(x, y, color.RGBA{
				R: r8,
				G: g8,
				B: b8,
				A: 255,
			})
		}
	}

	return dst
}

// GenerateOutputPath generates the output JPEG file path from input HEIC path
func GenerateOutputPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	basePath := strings.TrimSuffix(inputPath, ext)
	return basePath + ".jpg"
}
