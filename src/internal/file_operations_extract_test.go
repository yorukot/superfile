package internal

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"golift.io/xtractr"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestExtractCompressFile(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			files := map[string]string{
				"file.txt":      "first file",
				"nested/中文.txt": "nested file",
			}
			writeExtractionArchive(t, archive, format, files)

			require.NoError(t, extractCompressFile(archive, dest, extractionProcessBar(t)))
			for name, want := range files {
				content, err := os.ReadFile(filepath.Join(dest, filepath.FromSlash(name)))
				require.NoError(t, err)
				require.Equal(t, want, string(content))
			}
		})
	}
}

func TestExtractCompressFilePreservesTarPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix execute permissions are not supported on Windows")
	}
	for _, test := range []struct {
		name string
		mode os.FileMode
	}{
		{name: "executable", mode: 0o755},
		{name: "regular", mode: 0o644},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			content := []byte("archive member\n")
			// A reference file applies the current umask without changing process-wide state.
			reference := filepath.Join(root, "reference")
			require.NoError(t, os.WriteFile(reference, content, test.mode))
			want, err := os.Stat(reference)
			require.NoError(t, err)

			var data bytes.Buffer
			writer := tar.NewWriter(&data)
			require.NoError(t, writer.WriteHeader(&tar.Header{
				Name: "member",
				Mode: int64(test.mode),
				Size: int64(len(content)),
			}))
			_, err = writer.Write(content)
			require.NoError(t, err)
			require.NoError(t, writer.Close())
			archive := filepath.Join(root, "archive.tar")
			require.NoError(t, os.WriteFile(archive, data.Bytes(), utils.ExtractedFileMode))

			dest := filepath.Join(root, "output")
			require.NoError(t, extractCompressFile(archive, dest, extractionProcessBar(t)))
			got, err := os.Stat(filepath.Join(dest, "member"))
			require.NoError(t, err)
			require.Equal(t, want.Mode().Perm(), got.Mode().Perm(),
				"archive=%04o expected=%04o got=%04o", test.mode, want.Mode().Perm(), got.Mode().Perm())
		})
	}
}

func TestExtractCompressFileRejectsSiblingPathEscape(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			outside := filepath.Join(root, "output-sibling", "keep.txt")
			utils.SetupDirectories(t, filepath.Dir(outside))
			utils.SetupFilesWithData(t, []byte("keep me"), outside)
			writeExtractionArchive(t, archive, format, map[string]string{
				"../output-sibling/keep.txt": "overwritten",
			})

			err := extractCompressFile(archive, dest, extractionProcessBar(t))
			require.ErrorIs(t, err, xtractr.ErrInvalidPath)
			content, err := os.ReadFile(outside)
			require.NoError(t, err)
			require.Equal(t, "keep me", string(content))
		})
	}
}

func TestExtractCompressFileRejectsOutputSymlink(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			outside := filepath.Join(root, "outside", "keep.txt")
			utils.SetupDirectories(t, dest, filepath.Dir(outside))
			utils.SetupFilesWithData(t, []byte("keep me"), outside)
			if err := os.Symlink(filepath.Dir(outside), filepath.Join(dest, "link")); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			writeExtractionArchive(t, archive, format, map[string]string{
				"link/keep.txt": "overwritten",
			})

			err := extractCompressFile(archive, dest, extractionProcessBar(t))
			require.ErrorIs(t, err, xtractr.ErrInvalidPath)
			content, err := os.ReadFile(outside)
			require.NoError(t, err)
			require.Equal(t, "keep me", string(content))
		})
	}
}

func extractionProcessBar(t *testing.T) *processbar.Model {
	t.Helper()
	model := processbar.New()
	model.ListenForChannelUpdates()
	t.Cleanup(model.SendStopListeningMsgBlocking)
	return &model
}

func writeExtractionArchive(t *testing.T, path, format string, files map[string]string) {
	t.Helper()
	var data bytes.Buffer
	switch format {
	case "zip":
		writer := zip.NewWriter(&data)
		for name, content := range files {
			file, err := writer.Create(name)
			require.NoError(t, err)
			_, err = io.WriteString(file, content)
			require.NoError(t, err)
		}
		require.NoError(t, writer.Close())
	case "tar":
		writer := tar.NewWriter(&data)
		for name, content := range files {
			require.NoError(t, writer.WriteHeader(&tar.Header{
				Name: name,
				Mode: utils.ExtractedFileMode,
				Size: int64(len(content)),
			}))
			_, err := io.WriteString(writer, content)
			require.NoError(t, err)
		}
		require.NoError(t, writer.Close())
	default:
		t.Fatalf("unsupported archive format: %s", format)
	}
	require.NoError(t, os.WriteFile(path, data.Bytes(), utils.ExtractedFileMode))
}
