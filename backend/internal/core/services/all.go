// Package services provides the core business logic for the application.
package services

import (
	"path/filepath"
	"strings"

	"github.com/hay-kot/homebox/backend/internal/data/repo"
)

type AllServices struct {
	User              *UserService
	Group             *GroupService
	Items             *ItemService
	BackgroundService *BackgroundService
	Import            *ImportService
}

type OptionsFunc func(*options)

type options struct {
	autoIncrementAssetID bool
	importDirs           string
}

func WithAutoIncrementAssetID(v bool) func(*options) {
	return func(o *options) {
		o.autoIncrementAssetID = v
	}
}

func WithImportDirs(dirs string) func(*options) {
	return func(o *options) {
		o.importDirs = dirs
	}
}

func New(repos *repo.AllRepos, opts ...OptionsFunc) *AllServices {
	if repos == nil {
		panic("repos cannot be nil")
	}

	options := &options{
		autoIncrementAssetID: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &AllServices{
		User:  &UserService{repos},
		Group: &GroupService{repos},
		Items: &ItemService{
			repo:                 repos,
			autoIncrementAssetID: options.autoIncrementAssetID,
		},
		BackgroundService: &BackgroundService{repos},
		Import: &ImportService{
			repo:       repos,
			importDirs: parseImportDirs(options.importDirs),
		},
	}
}

// parseImportDirs splits a comma-separated list of import directory paths
// and returns a cleaned, non-empty slice.
func parseImportDirs(dirs string) []string {
	if dirs == "" {
		return nil
	}
	parts := strings.Split(dirs, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = filepath.Clean(strings.TrimSpace(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
