package services

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// Register additional image decoders (blank imports trigger init())
	_ "github.com/gen2brain/heic"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ThumbnailQuality is the JPEG quality for generated thumbnails (1-100).
const ThumbnailQuality = 85

// Thumbnail returns the filesystem path to a thumbnail of the given width.
// If the thumbnail already exists on disk it is returned immediately.
// Otherwise the source image is decoded, resized, encoded as JPEG, and cached.
func (svc *ItemService) Thumbnail(ctx context.Context, attachmentID uuid.UUID, width int) (string, error) {
	doc, err := svc.AttachmentPath(ctx, attachmentID)
	if err != nil {
		return "", fmt.Errorf("thumbnail: get attachment: %w", err)
	}

	srcPath := doc.Path

	if width <= 0 {
		return srcPath, nil
	}

	cachePath := thumbCachePath(srcPath, width)

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	if err := generateThumb(srcPath, cachePath, width); err != nil {
		log.Err(err).
			Str("src", srcPath).
			Str("cache", cachePath).
			Int("width", width).
			Msg("failed to generate thumbnail, falling back to original")
		return srcPath, nil
	}

	return cachePath, nil
}

// thumbCachePath derives the cache filename for a given source path and width.
func thumbCachePath(srcPath string, width int) string {
	ext := strings.ToLower(filepath.Ext(srcPath))
	base := strings.TrimSuffix(srcPath, ext)
	return fmt.Sprintf("%s.thumb_w%d.jpg", base, width)
}

// generateThumb reads srcPath, resizes to the given width (height proportional),
// and writes a JPEG to cachePath.
func generateThumb(srcPath, cachePath string, width int) error {
	src, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}

	bounds := src.Bounds()
	if bounds.Dx() <= width {
		// Source is already small enough; just re-encode as JPEG
		return saveJPEG(src, cachePath)
	}

	resized := imaging.Resize(src, width, 0, imaging.Lanczos)
	return saveJPEG(resized, cachePath)
}

// saveJPEG encodes img as a JPEG file at the given path with the default quality.
func saveJPEG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create cache file: %w", err)
	}
	defer func() { _ = f.Close() }()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: ThumbnailQuality})
}
