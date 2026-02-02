package internal

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"filippo.io/age"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

// encryptFileWithPassword encrypts a single file using age passphrase-based encryption
//
//nolint:funlen,gocognit,goconst // Encryption requires comprehensive error handling
func encryptFileWithPassword(srcPath, password string, processBar *processbar.Model) error {
	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	// Generate output filename
	outputPath := srcPath + ".age"
	outputPath, err := renameIfDuplicate(outputPath)
	if err != nil {
		return fmt.Errorf("failed to generate unique filename: %w", err)
	}

	// Create process bar entry
	p, err := processBar.SendAddProcessMsg(filepath.Base(srcPath), processbar.OpEncrypt, 1, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process: %w", err)
	}

	// Read source file
	plaintext, err := os.ReadFile(srcPath)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to read file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Create age recipient with password
	recipient, err := age.NewScryptRecipient(password)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to create encryption key"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create encryption key: %w", err)
	}

	// Encrypt the data
	var encrypted bytes.Buffer
	writer, err := age.Encrypt(&encrypted, recipient)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create encryption writer: %w", err)
	}

	if _, err := writer.Write(plaintext); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if err := writer.Close(); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to finalize encryption: %w", err)
	}

	// Write encrypted data to file
	if err := os.WriteFile(outputPath, encrypted.Bytes(), 0600); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to write encrypted file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	// Update process bar as successful
	p.State = processbar.Successful
	p.Done = 1
	p.DoneTime = time.Now()
	if err := processBar.SendUpdateProcessMsg(p, true); err != nil {
		slog.Error("Error sending process update", "error", err)
	}

	return nil
}

// encryptFolderWithPassword creates a tar archive of the folder and encrypts it
//
//nolint:funlen,gocognit // Encryption requires comprehensive error handling
func encryptFolderWithPassword(srcPath, password string, processBar *processbar.Model) error {
	// Check if source directory exists
	info, err := os.Stat(srcPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", srcPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", srcPath)
	}

	// Count files for progress tracking
	fileCount, err := countFiles(srcPath)
	if err != nil {
		slog.Error("Error counting files", "error", err)
		fileCount = 1
	}

	// Generate output filename
	outputPath := srcPath + ".tar.age"
	outputPath, err = renameIfDuplicate(outputPath)
	if err != nil {
		return fmt.Errorf("failed to generate unique filename: %w", err)
	}

	// Create process bar entry
	p, err := processBar.SendAddProcessMsg(filepath.Base(srcPath), processbar.OpEncrypt, fileCount, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process: %w", err)
	}

	// Create tar archive in memory
	var tarBuffer bytes.Buffer
	tarWriter := tar.NewWriter(&tarBuffer)

	// Walk through directory and add files to tar
	srcParentDir := filepath.Dir(srcPath)
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(srcParentDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If it's a file, write content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		p.CurrentFile = filepath.Base(path)
		p.Done++
		processBar.TrySendingUpdateProcessMsg(p)

		return nil
	})

	if err != nil {
		slog.Error("Error while creating tar archive", "error", err)
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to create archive"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create tar archive: %w", err)
	}

	if closeErr := tarWriter.Close(); closeErr != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to finalize archive"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to close tar writer: %w", closeErr)
	}

	// Create age recipient with password
	recipient, err := age.NewScryptRecipient(password)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to create encryption key"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create encryption key: %w", err)
	}

	// Encrypt the tar archive
	var encrypted bytes.Buffer
	writer, err := age.Encrypt(&encrypted, recipient)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create encryption writer: %w", err)
	}

	if _, err := writer.Write(tarBuffer.Bytes()); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if err := writer.Close(); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Encryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to finalize encryption: %w", err)
	}

	// Write encrypted data to file
	if err := os.WriteFile(outputPath, encrypted.Bytes(), 0600); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to write encrypted file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	// Update process bar as successful
	p.State = processbar.Successful
	p.Done = fileCount
	p.DoneTime = time.Now()
	if err := processBar.SendUpdateProcessMsg(p, true); err != nil {
		slog.Error("Error sending process update", "error", err)
	}

	return nil
}

// decryptFileWithPassword decrypts an age-encrypted file
//
//nolint:funlen // Decryption requires comprehensive error handling
func decryptFileWithPassword(srcPath, password string, processBar *processbar.Model) error {
	// Check if source file exists and has .age extension
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	if !strings.HasSuffix(strings.ToLower(srcPath), ".age") {
		return errors.New("file does not have .age extension")
	}

	// Generate output filename by removing .age extension
	outputPath := strings.TrimSuffix(srcPath, ".age")
	outputPath, err := renameIfDuplicate(outputPath)
	if err != nil {
		return fmt.Errorf("failed to generate unique filename: %w", err)
	}

	// Create process bar entry
	p, err := processBar.SendAddProcessMsg(filepath.Base(srcPath), processbar.OpDecrypt, 1, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process: %w", err)
	}

	// Read encrypted file
	encrypted, err := os.ReadFile(srcPath)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to read file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Create age identity with password
	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to create decryption key"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create decryption key: %w", err)
	}

	// Decrypt the data
	reader, err := age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Decryption failed - wrong password?"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to decrypt (wrong password?): %w", err)
	}

	plaintext, err := io.ReadAll(reader)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Decryption failed"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to read decrypted data: %w", err)
	}

	// Write decrypted data to file
	if err := os.WriteFile(outputPath, plaintext, 0600); err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to write decrypted file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	// Update process bar as successful
	p.State = processbar.Successful
	p.Done = 1
	p.DoneTime = time.Now()
	if err := processBar.SendUpdateProcessMsg(p, true); err != nil {
		slog.Error("Error sending process update", "error", err)
	}

	return nil
}

// decryptAndExtractTarWithPassword decrypts a .tar.age file and extracts its contents
//
//nolint:funlen,gocognit,gosec // Tar extraction requires comprehensive error handling and security checks
func decryptAndExtractTarWithPassword(srcPath, password string, processBar *processbar.Model) error {
	// Check if source file exists and has .tar.age extension
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	if !strings.HasSuffix(strings.ToLower(srcPath), ".tar.age") {
		return errors.New("file does not have .tar.age extension")
	}

	// Create process bar entry
	p, err := processBar.SendAddProcessMsg(filepath.Base(srcPath), processbar.OpDecrypt, 1, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process: %w", err)
	}

	// Read encrypted file
	encrypted, err := os.ReadFile(srcPath)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to read file"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Create age identity with password
	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Failed to create decryption key"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to create decryption key: %w", err)
	}

	// Decrypt the data
	reader, err := age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = "Decryption failed - wrong password?"
		p.DoneTime = time.Now()
		if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
			slog.Error("Failed to send process update", "error", sendErr)
		}
		return fmt.Errorf("failed to decrypt (wrong password?): %w", err)
	}

	// Extract tar archive
	tarReader := tar.NewReader(reader)
	extractDir := filepath.Dir(srcPath)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.State = processbar.Failed
			p.ErrorMsg = "Failed to extract archive"
			p.DoneTime = time.Now()
			if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
				slog.Error("Failed to send process update", "error", sendErr)
			}
			return fmt.Errorf("failed to read tar archive: %w", err)
		}

		targetPath := filepath.Join(extractDir, header.Name)

		// Security check: prevent directory traversal
		if !strings.HasPrefix(targetPath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
			continue
		}

		p.CurrentFile = filepath.Base(header.Name)
		processBar.TrySendingUpdateProcessMsg(p)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0700); err != nil {
				p.State = processbar.Failed
				p.ErrorMsg = "Failed to create directory"
				p.DoneTime = time.Now()
				if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
					slog.Error("Failed to send process update", "error", sendErr)
				}
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directories if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0700); err != nil {
				p.State = processbar.Failed
				p.ErrorMsg = "Failed to create directory"
				p.DoneTime = time.Now()
				if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
					slog.Error("Failed to send process update", "error", sendErr)
				}
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Use 0600 for files (instead of tar header mode for security)
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				p.State = processbar.Failed
				p.ErrorMsg = "Failed to create file"
				p.DoneTime = time.Now()
				if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
					slog.Error("Failed to send process update", "error", sendErr)
				}
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				if closeErr := outFile.Close(); closeErr != nil {
					slog.Error("Failed to close file", "error", closeErr)
				}
				p.State = processbar.Failed
				p.ErrorMsg = "Failed to extract file"
				p.DoneTime = time.Now()
				if sendErr := processBar.SendUpdateProcessMsg(p, true); sendErr != nil {
					slog.Error("Failed to send process update", "error", sendErr)
				}
				return fmt.Errorf("failed to extract file: %w", err)
			}
			if err := outFile.Close(); err != nil {
				slog.Error("Failed to close file", "error", err)
			}
		}
	}

	// Update process bar as successful
	p.State = processbar.Successful
	p.Done = 1
	p.DoneTime = time.Now()
	if err := processBar.SendUpdateProcessMsg(p, true); err != nil {
		slog.Error("Error sending process update", "error", err)
	}

	return nil
}
