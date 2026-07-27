package v1

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hay-kot/homebox/backend/internal/core/services"
	"github.com/hay-kot/homebox/backend/internal/data/ent/attachment"
	"github.com/hay-kot/homebox/backend/internal/sys/validate"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/hay-kot/httpkit/server"
	"github.com/rs/zerolog/log"
)

// ImportAttachmentBody is the JSON body for importing a server-side file as an attachment.
type ImportAttachmentBody struct {
	SourcePath string `json:"sourcePath"`
	Type       string `json:"type,omitempty"`
}

// HandleImportBrowse godoc
//
//	@Summary  Browse Import Directories
//	@Tags     Import
//	@Produce  json
//	@Param    path     query string false "Subdirectory path to browse"
//	@Param    page     query int    false "Page number (files only, default 1)"
//	@Param    pageSize query int    false "Items per page (files only, default 50)"
//	@Success  200      {object} services.BrowseResult
//	@Router   /v1/import/browse [GET]
//	@Security Bearer
func (ctrl *V1Controller) HandleImportBrowse() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		subPath := r.URL.Query().Get("path")
		page := queryIntOrNegativeOne(r.URL.Query().Get("page"))
		pageSize := queryIntOrNegativeOne(r.URL.Query().Get("pageSize"))

		result, err := ctrl.svc.Import.Browse(subPath, page, pageSize)
		if err != nil {
			log.Err(err).Str("path", subPath).Msg("failed to browse import directory")
			return validate.NewRequestError(err, http.StatusBadRequest)
		}

		return server.JSON(w, http.StatusOK, result)
	}
}

// HandleImportThumb godoc
//
//	@Summary  Get Import Thumbnail
//	@Tags     Import
//	@Produce  image/jpeg
//	@Param    path query string true  "Path to the file in the import directory"
//	@Param    w    query int    false "Thumbnail width"
//	@Success  200  {file} image/jpeg
//	@Router   /v1/import/thumb [GET]
//	@Security Bearer
func (ctrl *V1Controller) HandleImportThumb() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		importPath := r.URL.Query().Get("path")
		if importPath == "" {
			return validate.NewRequestError(
				validate.NewFieldErrors().Append("path", "path is required"),
				http.StatusBadRequest,
			)
		}

		widthStr := r.URL.Query().Get("w")
		width := 0
		if widthStr != "" {
			var err error
			width, err = strconv.Atoi(widthStr)
			if err != nil || width < 0 {
				return validate.NewRequestError(
					validate.NewFieldErrors().Append("w", "invalid width"),
					http.StatusBadRequest,
				)
			}
		}

		target, err := ctrl.svc.Import.ThumbnailPath(importPath, width)
		if err != nil {
			log.Err(err).Str("path", importPath).Msg("failed to get import thumbnail")
			return validate.NewRequestError(err, http.StatusBadRequest)
		}

		http.ServeFile(w, r, target)
		return nil
	}
}

// HandleImportAttachmentCreate godoc
//
//	@Summary  Import Server-Side File as Attachment
//	@Tags     Items Attachments
//	@Accept   json
//	@Produce  json
//	@Param    id      path string               true "Item ID"
//	@Param    payload body ImportAttachmentBody  true "Import data"
//	@Success  201     {object} repo.ItemOut
//	@Router   /v1/items/{id}/attachments/import [POST]
//	@Security Bearer
func (ctrl *V1Controller) HandleImportAttachmentCreate() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var body ImportAttachmentBody
		if err := server.Decode(r, &body); err != nil {
			return validate.NewRequestError(err, http.StatusBadRequest)
		}

		errs := validate.NewFieldErrors()

		if body.SourcePath == "" {
			errs = errs.Append("sourcePath", "sourcePath is required")
		}

		if !errs.Nil() {
			return server.JSON(w, http.StatusUnprocessableEntity, errs)
		}

		// Auto-detect attachment type from extension if not provided
		attachmentType := body.Type
		if attachmentType == "" {
			ext := strings.ToLower(filepath.Ext(body.SourcePath))
			switch ext {
			case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp", ".tiff", ".heic", ".heif":
				attachmentType = attachment.TypePhoto.String()
			default:
				attachmentType = attachment.TypeAttachment.String()
			}
		}

		id, err := ctrl.routeID(r)
		if err != nil {
			return err
		}

		ctx := services.NewContext(r.Context())

		item, err := ctrl.svc.Import.ImportAttachment(
			ctx,
			id,
			body.SourcePath,
			attachment.Type(attachmentType),
		)
		if err != nil {
			log.Err(err).
				Str("sourcePath", body.SourcePath).
				Str("itemID", id.String()).
				Msg("failed to import attachment")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusCreated, item)
	}
}
