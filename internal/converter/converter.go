// Package converter provides functionality for converting HEIC image files to JPEG format.
package converter

import (
	"bytes"
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

	// maxEXIFSegmentPayload is the largest EXIF payload that can fit in a single
	// JPEG APP1 segment. A segment's length field is 2 bytes and covers itself,
	// so the payload (which starts with the "Exif\0\0" marker) may be at most
	// 0xFFFF - 2 bytes.
	maxEXIFSegmentPayload = 0xFFFF - 2

	// jpegAPP1Marker is the JPEG marker used for EXIF (and XMP) segments.
	jpegAPP1Marker = 0xE1
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
	// RemoveEXIF controls whether EXIF metadata from the source HEIC file is
	// carried over to the converted JPEG. When true, the JPEG is written
	// without any EXIF data. When false, EXIF metadata found in the HEIC
	// source is embedded into the output JPEG.
	RemoveEXIF bool
}

// ConvertHEICToJPEG converts a HEIC file to JPEG format
func ConvertHEICToJPEG(inputPath string, options ConvertOptions) error {
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

	// Extract EXIF metadata from the source HEIC file, unless the caller
	// asked for it to be stripped. Extraction failures (e.g. no EXIF present)
	// are non-fatal: the conversion simply proceeds without EXIF data.
	var exifSegment []byte
	if !options.RemoveEXIF {
		if exifData, exifErr := goheif.ExtractExif(file); exifErr == nil {
			exifSegment = buildEXIFAPP1Segment(exifData)
		}
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

	// Encode as JPEG into a buffer so an EXIF segment can be spliced in
	// right after the SOI marker.
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: JPEGQuality}
	if err := jpeg.Encode(&buf, encodeImg, opts); err != nil {
		return fmt.Errorf("JPEGファイルのエンコードに失敗しました: %w", err)
	}

	if err := writeJPEGWithEXIF(outFile, buf.Bytes(), exifSegment); err != nil {
		return fmt.Errorf("JPEGファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// writeJPEGWithEXIF writes JPEG data to w, inserting exifSegment (a complete
// APP1 marker segment, or nil) immediately after the leading SOI marker.
func writeJPEGWithEXIF(w *os.File, jpegData []byte, exifSegment []byte) error {
	if len(exifSegment) == 0 {
		_, err := w.Write(jpegData)
		return err
	}

	// jpegData always starts with the 2-byte SOI marker (0xFFD8).
	if _, err := w.Write(jpegData[:2]); err != nil {
		return err
	}
	if _, err := w.Write(exifSegment); err != nil {
		return err
	}
	_, err := w.Write(jpegData[2:])
	return err
}

// buildEXIFAPP1Segment builds a complete JPEG APP1 marker segment (marker +
// length + payload) embedding the given EXIF payload. The payload is
// expected to already carry the "Exif\0\0" marker, as returned by
// goheif.ExtractExif. Returns nil if there is no usable EXIF data or the
// payload is too large to fit in a single APP1 segment.
func buildEXIFAPP1Segment(exifData []byte) []byte {
	if len(exifData) == 0 || len(exifData) > maxEXIFSegmentPayload {
		return nil
	}

	length := len(exifData) + 2 // length field covers itself
	segment := make([]byte, 0, length+2)
	segment = append(segment, 0xFF, jpegAPP1Marker)
	segment = append(segment, byte(length>>8), byte(length&0xFF))
	segment = append(segment, exifData...)
	return segment
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
