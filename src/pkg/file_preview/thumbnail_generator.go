package filepreview

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = "104857600" // 100MB limit
	outputExt   = ".jpg"
)

type ThumbnailGenerator struct {
	// This is a cache. Key -> Video file path, Value -> Thumbnail file path
	// TODO: We can potentially make it persisitent, preventing generation
	// of thumbnail on every launch or superfile
	tempFilesCache map[string]string
	tempDirectory  string
	mu             sync.Mutex
}

func NewThumbnailGenerator() (*ThumbnailGenerator, error) {
	if !isFFmpegInstalled() {
		return nil, errors.New("ffmpeg is not installed")
	}

	tmp, err := os.MkdirTemp("", "superfiles-*")
	if err != nil {
		return nil, err
	}

	thumbnailGenerator := &ThumbnailGenerator{
		tempFilesCache: make(map[string]string),
		tempDirectory:  tmp,
	}

	return thumbnailGenerator, nil
}

func (g *ThumbnailGenerator) GetThumbnailOrGenerate(path string) (string, error) {
	g.mu.Lock()
	file, ok := g.tempFilesCache[path]
	g.mu.Unlock()

	if ok {
		_, err := os.Stat(file)
		if err == nil {
			return file, nil
		}

		g.mu.Lock()
		delete(g.tempFilesCache, path)
		g.mu.Unlock()
	}

	generatedThumbnailPath, err := g.generateThumbnail(path)
	if err != nil {
		return "", err
	}

	g.mu.Lock()
	g.tempFilesCache[path] = generatedThumbnailPath
	g.mu.Unlock()

	return generatedThumbnailPath, nil
}

func (g *ThumbnailGenerator) generateThumbnail(inputPath string) (string, error) {
	fileExt := filepath.Ext(inputPath)
	filename := filepath.Base(inputPath)
	baseName := filename[:len(filename)-len(fileExt)]

	outputFile, err := os.CreateTemp(g.tempDirectory, "*-"+baseName+outputExt)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	outputFilePath := outputFile.Name()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ffmpeg -v warning -t 60 -hwaccel auto -an -sn -dn -skip_frame nokey -i input.mkv -vf scale='min(1024,iw)':'min(720,ih)':force_original_aspect_ratio=decrease:flags=fast_bilinear -vf "thumbnail" -frames:v 1 -y thumb.jpg
	ffmpeg := exec.CommandContext(ctx, "ffmpeg",
		"-v", "warning", // set log level to warning
		"-an",       // disable Audio stream
		"-sn",       // disable Subtitle stream
		"-dn",       // disable data stream
		"-t", "180", // process maximum 180s of the video (the first 3 min)
		"-hwaccel", "auto", // Use Hardware Acceleration if available
		"-skip_frame", "nokey", // skip non-key frames
		"-i", inputPath, // set input file
		"-vf", "thumbnail", // use ffmpeg default thumbnail filter
		"-frames:v", "1", // output only one frame (one image)
		"-f", "image2", // set format to image2
		"-fs", maxFileSize, // limit the max file size to match image previewer limit
		"-y", outputFilePath, // set the outputFile and overwrite it without confirmation if already exists
	)

	err = ffmpeg.Run()
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}

func (g *ThumbnailGenerator) CleanUp() error {
	return os.RemoveAll(g.tempDirectory)
}

func isFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}
