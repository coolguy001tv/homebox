package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/homebox/backend/internal/data/ent/attachment"
	"github.com/hay-kot/homebox/backend/internal/data/repo"
	"github.com/rs/zerolog/log"
)

// ImportService handles browsing and importing files from server-side directories.
type ImportService struct {
	repo           *repo.AllRepos
	importDirs     []string
	importThumbDir string // NAS 缩略图子目录名
	thumbCacheDir  string // Homebox 本地缩略图缓存目录
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

// BrowseResult separates directories (always all) from paginated files.
type BrowseResult struct {
	Dirs     []FileEntry `json:"dirs"`
	Files    []FileEntry `json:"files"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"` // total number of files (for pagination)
}

var imageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	".gif": true, ".bmp": true, ".tiff": true, ".heic": true, ".heif": true,
}

// Browse lists files and directories under the given sub-path within the
// allowed import directories. An empty subPath returns the root listing.
// Directories are always returned in full; files are paginated by page and pageSize.
func (svc *ImportService) Browse(subPath string, page, pageSize int) (BrowseResult, error) {
	realPath, err := svc.resolvePath(subPath)
	if err != nil {
		return BrowseResult{}, err
	}

	entries, err := os.ReadDir(realPath)
	if err != nil {
		return BrowseResult{}, fmt.Errorf("read directory: %w", err)
	}

	dirs := make([]FileEntry, 0)
	files := make([]FileEntry, 0)

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}

		// Skip thumbnail cache files to avoid listing them as entries.
		if !e.IsDir() && strings.Contains(e.Name(), ".thumb_w") {
			continue
		}

		entryPath := filepath.Join(realPath, e.Name())

		isDir := e.IsDir()
		ext := ""
		if !isDir {
			ext = strings.ToLower(filepath.Ext(e.Name()))
		}

		entry := FileEntry{
			Name:    e.Name(),
			Path:    entryPath,
			Size:    info.Size(),
			IsDir:   isDir,
			Ext:     ext,
			ModTime: info.ModTime(),
			IsImage: imageExtensions[ext],
		}

		if isDir {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// Sort directories and files alphabetically by name (descending)
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name > dirs[j].Name
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})

	// Apply defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	// Paginate files
	totalFiles := len(files)
	start := (page - 1) * pageSize
	if start > totalFiles {
		start = totalFiles
	}
	end := start + pageSize
	if end > totalFiles {
		end = totalFiles
	}

	return BrowseResult{
		Dirs:     dirs,
		Files:    files[start:end],
		Page:     page,
		PageSize: pageSize,
		Total:    totalFiles,
	}, nil
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

// ThumbnailPath returns a path to a thumbnail for a file in the import directory.
// It first checks for a NAS-generated thumbnail, then falls back to a locally
// cached thumbnail in the Homebox data directory, generating one if needed.
func (svc *ImportService) ThumbnailPath(importPath string, width int) (string, error) {
	realPath, err := svc.resolvePath(importPath)
	if err != nil {
		log.Err(err).Str("importPath", importPath).Msg("import thumb: resolvePath failed")
		return "", err
	}

	if width <= 0 {
		log.Info().Str("src", realPath).Msg("import thumb: width <= 0, serving original")
		return realPath, nil
	}

	// 1. Prefer NAS-generated thumbnail (e.g. .@__thumb/<imageName>)
	if nas := svc.findNasThumb(realPath); nas != "" {
		// Clean up any previously cached local thumbnail — the NAS thumb now
		// takes priority, so the local cache would never be used again.
		if cachePath, err := svc.importThumbCachePath(realPath, width); err == nil {
			if err := os.Remove(cachePath); err == nil {
				log.Debug().
					Str("src", realPath).
					Str("cache", cachePath).
					Msg("import thumb: removed stale local cache (NAS thumb now available)")
			}
		}
		log.Info().
			Str("src", realPath).
			Str("nasThumb", nas).
			Int("width", width).
			Msg("import thumb: hit NAS thumbnail")
		return nas, nil
	}

	// 2. Check Homebox local cache
	cachePath, err := svc.importThumbCachePath(realPath, width)
	if err != nil {
		log.Err(err).Str("src", realPath).Msg("import thumb: cache path error")
		return "", err
	}
	if _, err := os.Stat(cachePath); err == nil {
		log.Info().
			Str("src", realPath).
			Str("cache", cachePath).
			Int("width", width).
			Msg("import thumb: hit local cache")
		return cachePath, nil
	}

	// 3. Ensure cache directory exists and generate thumbnail
	if err := os.MkdirAll(svc.thumbCacheDir, 0o755); err != nil {
		log.Err(err).Str("dir", svc.thumbCacheDir).Msg("failed to create thumb cache directory")
		return realPath, nil
	}
	log.Info().
		Str("src", realPath).
		Str("cache", cachePath).
		Int("width", width).
		Msg("import thumb: cache miss, generating thumbnail")
	if err := generateThumb(realPath, cachePath, width); err != nil {
		log.Err(err).
			Str("src", realPath).
			Str("cache", cachePath).
			Int("width", width).
			Msg("failed to generate import thumbnail, falling back to original")
		return realPath, nil
	}

	log.Info().
		Str("src", realPath).
		Str("cache", cachePath).
		Int("width", width).
		Msg("import thumb: generated new thumbnail")
	return cachePath, nil
}

// findNasThumb looks for a NAS-generated thumbnail file for the given source
// file inside the configured import thumb directory (default .@__thumb).
// QNAP NAS stores thumbnails as "default" + original filename (e.g. "defaultIMG_1234.HEIC"),
// so we check both the exact filename and the "default"-prefixed variant.
// Returns empty string when not found or not configured.
func (svc *ImportService) findNasThumb(realPath string) string {
	if svc.importThumbDir == "" {
		return ""
	}

	dir := filepath.Dir(realPath)
	name := filepath.Base(realPath)
	thumbDir := filepath.Join(dir, svc.importThumbDir)

	// Try exact filename match first (e.g. .@__thumb/IMG_1234.HEIC)
	candidate := filepath.Join(thumbDir, name)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}

	// Try "default" prefix (QNAP NAS convention: .@__thumb/defaultIMG_1234.HEIC)
	candidate = filepath.Join(thumbDir, "default"+name)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}

	return ""
}

// importThumbCachePath derives a local cache path under thumbCacheDir for the
// given source file path and width. Uses an MD5 hash of the source path to
// produce a stable, filesystem-safe filename.
func (svc *ImportService) importThumbCachePath(srcPath string, width int) (string, error) {
	hash := md5.Sum([]byte(srcPath))
	hashStr := hex.EncodeToString(hash[:])
	ext := strings.ToLower(filepath.Ext(srcPath))

	if svc.thumbCacheDir == "" {
		return "", fmt.Errorf("thumb cache directory is not set")
	}

	return filepath.Join(svc.thumbCacheDir, fmt.Sprintf("%s_w%d%s", hashStr, width, ext)), nil
}

// SetThumbCacheDir sets the Homebox local thumbnail cache directory.
func (svc *ImportService) SetThumbCacheDir(dir string) {
	svc.thumbCacheDir = dir
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
