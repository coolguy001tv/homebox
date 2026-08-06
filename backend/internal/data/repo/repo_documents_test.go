package repo

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hay-kot/homebox/backend/internal/data/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func useDocs(t *testing.T, num int) []DocumentOut {
	t.Helper()

	results := make([]DocumentOut, 0, num)
	ids := make([]uuid.UUID, 0, num)

	for i := 0; i < num; i++ {
		doc, err := tRepos.Docs.Create(context.Background(), tGroup.ID, DocumentCreate{
			Title:   fk.Str(10) + ".md",
			Content: bytes.NewReader([]byte(fk.Str(10))),
		})

		require.NoError(t, err)
		assert.NotNil(t, doc)
		results = append(results, doc)
		ids = append(ids, doc.ID)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			err := tRepos.Docs.Delete(context.Background(), id)
			if err != nil {
				assert.True(t, ent.IsNotFound(err))
			}
		}
	})

	return results
}

func TestDocumentRepository_CreateUpdateDelete(t *testing.T) {
	temp := t.TempDir()
	r := DocumentRepository{
		db:  tClient,
		dir: temp,
	}

	type args struct {
		ctx context.Context
		gid uuid.UUID
		doc DocumentCreate
	}
	tests := []struct {
		name    string
		content string
		args    args
		title   string
		wantErr bool
	}{
		{
			name:    "basic create",
			title:   "test.md",
			content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			args: args{
				ctx: context.Background(),
				gid: tGroup.ID,
				doc: DocumentCreate{
					Title:   "test.md",
					Content: bytes.NewReader([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Document
			got, err := r.Create(tt.args.ctx, tt.args.gid, tt.args.doc)
			require.NoError(t, err)
			assert.Equal(t, tt.title, got.Title)
			assert.Equal(t, fmt.Sprintf("%s/%s/documents", temp, tt.args.gid), filepath.Dir(got.Path))

			ensureRead := func() {
				// Read Document
				bts, err := os.ReadFile(got.Path)
				require.NoError(t, err)
				assert.Equal(t, tt.content, string(bts))
			}
			ensureRead()

			// Update Document
			got, err = r.Rename(tt.args.ctx, got.ID, "__"+tt.title+"__")
			require.NoError(t, err)
			assert.Equal(t, "__"+tt.title+"__", got.Title)

			ensureRead()

			// Delete Document
			err = r.Delete(tt.args.ctx, got.ID)
			require.NoError(t, err)

			_, err = os.Stat(got.Path)
			require.Error(t, err)
		})
	}
}

// openFDCount returns the number of file descriptors in this process that
// point at path. It is only meaningful on Linux, which exposes the process
// fd table via /proc/self/fd; on other platforms it returns -1 and callers
// should skip the assertion.
func openFDCount(path string) int {
	fds, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return -1
	}

	count := 0
	for _, fd := range fds {
		link, err := os.Readlink(filepath.Join("/proc/self/fd", fd.Name()))
		if err != nil {
			continue
		}
		if link == path {
			count++
		}
	}
	return count
}

// TestDocumentRepository_CreateClosesFile is a regression test for a file
// handle leak: Create() opened the target file with os.Create but never
// closed it, leaving the descriptor open after the call returned. On Windows
// that made Delete() fail with "file in use" (HTTP 500 when removing an
// attachment); on Linux the descriptor leaked silently. A successful Create()
// must not leave any descriptor to the file open.
func TestDocumentRepository_CreateClosesFile(t *testing.T) {
	temp := t.TempDir()
	r := DocumentRepository{
		db:  tClient,
		dir: temp,
	}

	doc, err := r.Create(context.Background(), tGroup.ID, DocumentCreate{
		Title:   "leak-test.md",
		Content: bytes.NewReader([]byte("hello")),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = r.Delete(context.Background(), doc.ID)
	})

	if fdCount := openFDCount(doc.Path); fdCount >= 0 {
		require.Zero(t, fdCount, "Create() left %d open file handle(s) to %s", fdCount, doc.Path)
	}
}
