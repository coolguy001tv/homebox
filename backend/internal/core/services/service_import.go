package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/homebox/backend/internal/data/ent/attachment"
	"github.com/hay-kot/homebox/backend/internal/data/repo"
	"github.com/rs/zerolog/log"
)

// ImportService handles browsing and importing files from server-side directories.
type ImportService struct {
	repo       *repo.AllRepos
	importDirs []string
}

// FileEntry represents a file or directory in the import browser.
type FileEntry struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isDir"`
	Ext     string    `json:"extension"`
	ModTime time.Time `json:"modified"`
	IsImage bool      `json:"isImage"`
}

var imageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	".gif": true, ".bmp": true, ".tiff": true, ".heic": true, ".heif": true,
}

// Browse lists files and directories under the given sub-path within the
// allowed import directories. An empty subPath returns the root listing.
func (svc *ImportService) Browse(subPath string) ([]FileEntry, error) {
	realPath, err := svc.resolvePath(subPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(realPath)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	result := make([]FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}

		// Skip thumbnail cache files to avoid listing them as entries.
		// They match the pattern: <original>.thumb_w<width>.jpg
		if !e.IsDir() && strings.Contains(e.Name(), ".thumb_w") {
			continue
		}

		entryPath := filepath.Join(realPath, e.Name())

		isDir := e.IsDir()
		ext := ""
		if !isDir {
			ext = strings.ToLower(filepath.Ext(e.Name()))
		}

		result = append(result, FileEntry{
			Name:    e.Name(),
			Path:    entryPath,
			Size:    info.Size(),
			IsDir:   isDir,
			Ext:     ext,
			ModTime: info.ModTime(),
			IsImage: imageExtensions[ext],
		})
	}

	return result, nil
}

// ImportAttachment imports a file from a server-side directory as an item
// attachment by creating a symlink in the documents directory.
func (svc *ImportService) ImportAttachment(ctx Context, itemID uuid.UUID, sourcePath string, attachmentType attachment.Type) (repo.ItemOut, error) {
	// Validate the source path is within allowed directories
	_, err := svc.resolvePath(sourcePath)
	if err != nil {
		return repo.ItemOut{}, fmt.Errorf("invalid source path: %w", err)
	}

	// Verify the source file exists and is not a directory
	info, err := os.Stat(sourcePath)
	if err != nil {
		return repo.ItemOut{}, fmt.Errorf("stat source file: %w", err)
	}
	if info.IsDir() {
		return repo.ItemOut{}, fmt.Errorf("source path is a directory, not a file")
	}

	filename := filepath.Base(sourcePath)

	// Get the Item to verify ownership
	_, err = svc.repo.Items.GetOneByGroup(ctx, ctx.GID, itemID)
	if err != nil {
		return repo.ItemOut{}, err
	}

	// Create the document via symlink
	doc, err := svc.repo.Docs.CreateFromPath(ctx, ctx.GID, filename, sourcePath)
	if err != nil {
		log.Err(err).Msg("failed to create document from path")
		return repo.ItemOut{}, err
	}

	// Create the attachment
	_, err = svc.repo.Attachments.Create(ctx, itemID, doc.ID, attachmentType)
	if err != nil {
		log.Err(err).Msg("failed to create attachment")
		return repo.ItemOut{}, err
	}

	return svc.repo.Items.GetOneByGroup(ctx, ctx.GID, itemID)
}

// ThumbnailPath generates a thumbnail for a file in the import directory
// and returns the cache path. The thumbnail is cached alongside the source file.
func (svc *ImportService) ThumbnailPath(importPath string, width int) (string, error) {
	realPath, err := svc.resolvePath(importPath)
	if err != nil {
		return "", err
	}

	if width <= 0 {
		return realPath, nil
	}

	cachePath := thumbCachePath(realPath, width)

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	if err := generateThumb(realPath, cachePath, width); err != nil {
		log.Err(err).
			Str("src", realPath).
			Str("cache", cachePath).
			Int("width", width).
			Msg("failed to generate import thumbnail, falling back to original")
		return realPath, nil
	}

	return cachePath, nil
}

// resolvePath validates that the given public path is within one of the
// allowed import directories and returns the resolved absolute path.
// If publicPath is already an absolute path under an import dir, it is
// validated and returned as-is. Otherwise it is treated as a relative
// sub-path and joined with each import dir.
func (svc *ImportService) resolvePath(publicPath string) (string, error) {
	if len(svc.importDirs) == 0 {
		return "", fmt.Errorf("no import directories configured")
	}

	clean := filepath.Clean(publicPath)

	// If clean is ".", treat as empty (root listing)
	if clean == "." {
		clean = ""
	}

	// If the path is already absolute, validate it directly instead of
	// joining it with the base dir (which would duplicate the prefix).
	if filepath.IsAbs(clean) && clean != "" {
		return svc.validateAbsPath(clean)
	}

	// Relative path: join with each base directory and validate.
	for _, base := range svc.importDirs {
		full := filepath.Join(base, clean)
		resolved, err := svc.validateAbsPath(full)
		if err == nil {
			return resolved, nil
		}
	}

	// Fallback: try with Abs
	for _, base := range svc.importDirs {
		full := filepath.Join(base, clean)
		absFull, _ := filepath.Abs(full)
		if resolved, err := svc.validateAbsPath(absFull); err == nil {
			return resolved, nil
		}
	}

	return "", fmt.Errorf("path is not within any allowed import directory")
}

// validateAbsPath checks that the given absolute path is within or equal to
// one of the allowed import directories and returns the cleaned resolved path.
func (svc *ImportService) validateAbsPath(absPath string) (string, error) {
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		realPath = absPath
	}
	realPath = filepath.Clean(realPath)

	for _, base := range svc.importDirs {
		if isWithin(realPath, base) {
			return realPath, nil
		}
	}

	// Fallback: try without symlink resolution (case-insensitive on Windows)
	for _, base := range svc.importDirs {
		absBase, _ := filepath.Abs(base)
		if isWithin(absPath, absBase) {
			return filepath.Clean(absPath), nil
		}
	}

	return "", fmt.Errorf("path is not within any allowed import directory")
}

// isWithin returns true if target is equal to base or is a descendant of base.
func isWithin(target, base string) bool {
	target = filepath.Clean(target)
	base = filepath.Clean(base)

	// Case-insensitive comparison on Windows
	if strings.EqualFold(target, base) {
		return true
	}

	basePrefix := base + string(os.PathSeparator)
	return len(target) > len(basePrefix) && strings.EqualFold(target[:len(basePrefix)], basePrefix)
}
