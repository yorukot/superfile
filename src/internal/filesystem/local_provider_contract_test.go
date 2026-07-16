package filesystem

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestLocalProviderContract(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		root := t.TempDir()
		utils.SetupDirectories(t, filepath.Join(root, "alpha-dir"))
		utils.SetupFilesWithData(t, []byte("alpha"), filepath.Join(root, "beta.txt"))

		session := newLocalTestSession(t, root)
		entries, err := session.List(context.Background(), NewLocalPath(root))
		require.NoError(t, err)
		require.Len(t, entries, 2)

		names := []string{entries[0].Name, entries[1].Name}
		sort.Strings(names)
		assert.Equal(t, []string{"alpha-dir", "beta.txt"}, names)
	})

	t.Run("stat", func(t *testing.T) {
		root := t.TempDir()
		filePath := filepath.Join(root, "stat.txt")
		utils.SetupFilesWithData(t, []byte("hello"), filePath)

		session := newLocalTestSession(t, root)
		stat, err := session.Stat(context.Background(), NewLocalPath(filePath))
		require.NoError(t, err)

		assert.Equal(t, "stat.txt", stat.Name)
		assert.EqualValues(t, 5, stat.Size)
		assert.False(t, stat.IsDir)
		assert.Equal(t, string(ProviderLocal), stat.ProviderID)
	})

	t.Run("read", func(t *testing.T) {
		root := t.TempDir()
		filePath := filepath.Join(root, "read.txt")
		utils.SetupFilesWithData(t, []byte("reader"), filePath)

		session := newLocalTestSession(t, root)
		reader, err := session.Read(context.Background(), NewLocalPath(filePath))
		require.NoError(t, err)
		defer reader.Close()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "reader", string(data))
	})

	t.Run("write", func(t *testing.T) {
		root := t.TempDir()
		filePath := filepath.Join(root, "write.txt")

		session := newLocalTestSession(t, root)
		err := session.Create(
			context.Background(),
			NewLocalPath(filePath),
			bytes.NewReader([]byte("writer")),
			CreateOptions{
				Mode:      utils.UserFilePerm,
				Overwrite: true,
			},
		)
		require.NoError(t, err)

		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, "writer", string(data))
	})

	t.Run("mkdir", func(t *testing.T) {
		root := t.TempDir()
		dirPath := filepath.Join(root, "nested", "child")

		session := newLocalTestSession(t, root)
		err := session.Mkdir(context.Background(), NewLocalPath(dirPath), MkdirOptions{
			Mode:    utils.UserDirPerm,
			Parents: true,
		})
		require.NoError(t, err)
		assert.DirExists(t, dirPath)
	})

	t.Run("rename", func(t *testing.T) {
		root := t.TempDir()
		sourcePath := filepath.Join(root, "before.txt")
		destinationPath := filepath.Join(root, "after.txt")
		utils.SetupFilesWithData(t, []byte("rename"), sourcePath)

		session := newLocalTestSession(t, root)
		err := session.Rename(
			context.Background(),
			NewLocalPath(sourcePath),
			NewLocalPath(destinationPath),
			RenameOptions{
				Overwrite: true,
			},
		)
		require.NoError(t, err)

		assert.NoFileExists(t, sourcePath)
		assert.FileExists(t, destinationPath)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("file", func(t *testing.T) {
			root := t.TempDir()
			filePath := filepath.Join(root, "delete.txt")
			utils.SetupFilesWithData(t, []byte("delete"), filePath)

			session := newLocalTestSession(t, root)
			err := session.Delete(context.Background(), NewLocalPath(filePath), DeleteOptions{Recursive: true})
			require.NoError(t, err)
			assert.NoFileExists(t, filePath)
		})

		t.Run("recursive directory", func(t *testing.T) {
			root := t.TempDir()
			dirPath := filepath.Join(root, "delete-dir")
			utils.SetupDirectories(t, filepath.Join(dirPath, "child"))
			utils.SetupFilesWithData(t, []byte("delete"), filepath.Join(dirPath, "child", "file.txt"))

			session := newLocalTestSession(t, root)
			err := session.Delete(context.Background(), NewLocalPath(dirPath), DeleteOptions{Recursive: true})
			require.NoError(t, err)
			assert.NoDirExists(t, dirPath)
		})

		t.Run("missing recursive path keeps current local semantics", func(t *testing.T) {
			root := t.TempDir()
			missingPath := filepath.Join(root, "missing.txt")
			neighborPath := filepath.Join(root, "neighbor.txt")
			utils.SetupFilesWithData(t, []byte("still here"), neighborPath)

			session := newLocalTestSession(t, root)
			err := session.Delete(context.Background(), NewLocalPath(missingPath), DeleteOptions{Recursive: true})
			require.NoError(t, err)
			assert.FileExists(t, neighborPath)
		})
	})

	t.Run("copy", func(t *testing.T) {
		root := t.TempDir()
		sourceDir := filepath.Join(root, "copy-src")
		destinationDir := filepath.Join(root, "copy-dst")
		nestedDir := filepath.Join(sourceDir, "nested")
		nestedFile := filepath.Join(nestedDir, "file.txt")
		utils.SetupDirectories(t, nestedDir)
		utils.SetupFilesWithData(t, []byte("copy"), nestedFile)

		session := newLocalTestSession(t, root)
		err := session.Copy(context.Background(), NewLocalPath(sourceDir), NewLocalPath(destinationDir), CopyOptions{
			Overwrite: true,
			Recursive: true,
		})
		require.NoError(t, err)

		assert.DirExists(t, destinationDir)
		assert.FileExists(t, filepath.Join(destinationDir, "nested", "file.txt"))
		copiedData, err := os.ReadFile(filepath.Join(destinationDir, "nested", "file.txt"))
		require.NoError(t, err)
		assert.Equal(t, "copy", string(copiedData))
	})

	t.Run("move", func(t *testing.T) {
		root := t.TempDir()
		sourceDir := filepath.Join(root, "move-src")
		destinationDir := filepath.Join(root, "move-dst")
		utils.SetupDirectories(t, filepath.Join(sourceDir, "nested"))
		utils.SetupFilesWithData(t, []byte("move"), filepath.Join(sourceDir, "nested", "file.txt"))

		session := newLocalTestSession(t, root)
		err := session.Move(context.Background(), NewLocalPath(sourceDir), NewLocalPath(destinationDir), MoveOptions{
			Overwrite: true,
			Recursive: true,
		})
		require.NoError(t, err)

		assert.NoDirExists(t, sourceDir)
		assert.FileExists(t, filepath.Join(destinationDir, "nested", "file.txt"))
		movedData, err := os.ReadFile(filepath.Join(destinationDir, "nested", "file.txt"))
		require.NoError(t, err)
		assert.Equal(t, "move", string(movedData))
	})
}

func newLocalTestSession(t *testing.T, root string) Session {
	t.Helper()

	provider := NewLocalProvider()
	session, err := provider.Open(context.Background(), Location{
		Provider: ProviderLocal,
		Path:     NewLocalPath(root),
		Label:    localSessionLabel,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, session.Close())
	})

	return session
}
