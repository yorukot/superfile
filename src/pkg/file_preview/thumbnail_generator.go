package filepreview

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yorukot/superfile/src/internal/common"
)

type thumbnailGeneratorInterface interface {
	supportsExt(ext string) bool
	generateThumbnail(inputPath string, outputDir string) (string, error)
}

type VideoGenerator struct{}

func newVideoGenerator() (*VideoGenerator, error) {
	if !isFFmpegInstalled() {
		return nil, errors.New("ffmpeg is not installed")
	}

	return &VideoGenerator{}, nil
}

func (g *VideoGenerator) supportsExt(ext string) bool {
	return common.VideoExtensions[strings.ToLower(ext)]
}

func (g *VideoGenerator) generateThumbnail(inputPath string, outputDir string) (string, error) {
	fileExt := filepath.Ext(inputPath)
	filename := filepath.Base(inputPath)
	baseName := filename[:len(filename)-len(fileExt)]

	outputFile, err := os.CreateTemp(outputDir, "*-"+baseName+thumbOutputExt)
	if err != nil {
		return "", err
	}
	outputFilePath := outputFile.Name()
	_ = outputFile.Close()

	ctx, cancel := context.WithTimeout(context.Background(), thumbGenerationTimeout)
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
		"-fs", maxVideoFileSizeForThumb, // limit the max file size to match image previewer limit
		"-y", outputFilePath, // set the outputFile and overwrite it without confirmation if already exists
	)

	err = ffmpeg.Run()
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}

type pdfGenerator struct{}

func newPdfGenerator() (*pdfGenerator, error) {
	if !isPoppolerInstalled() {
		return nil, errors.New("poppler is not installed")
	}

	return &pdfGenerator{}, nil
}

func (g *pdfGenerator) supportsExt(ext string) bool {
	return strings.ToLower(ext) == ".pdf"
}

func (g *pdfGenerator) generateThumbnail(inputPath string, outputDir string) (string, error) {
	fileExt := filepath.Ext(inputPath)
	filename := filepath.Base(inputPath)
	baseName := filename[:len(filename)-len(fileExt)]
	outputPath := filepath.Join(outputDir, baseName)

	ctx, cancel := context.WithTimeout(context.Background(), thumbGenerationTimeout)
	defer cancel()

	// pdftoppm -singlefile -png prefixFilename
	pdfttoppm := exec.CommandContext(ctx, "pdftoppm",
		"-singlefile", // output only the first page as image
		"-jpeg",       // Image extension
		inputPath,     // Set input file
		outputPath,    // The output prefix. (pdftoppm will add the .jpg ext)

	)

	err := pdfttoppm.Run()
	if err != nil {
		fmt.Printf("Thumbnail: %s %s", outputPath, err)
		return "", err
	}

	return outputPath + thumbOutputExt, nil
}

type ThumbnailGenerator struct {
	// This is a cache. Key -> Video file path, Value -> Thumbnail file path
	// TODO: We can potentially make it persistent, preventing generation
	// of thumbnail on every launch or superfile
	tempFilesCache map[string]string
	tempDirectory  string
	mu             sync.Mutex
	generators     []thumbnailGeneratorInterface
}

func NewThumbnailGenerator() (*ThumbnailGenerator, error) {
	tmp, err := os.MkdirTemp("", "superfiles-*")
	if err != nil {
		return nil, err
	}

	generators := []thumbnailGeneratorInterface{}

	pdf, err := newPdfGenerator()
	if err != nil {
		slog.Error("Error while trying to create pdfGenerator", "error", err)
	} else {
		generators = append(generators, pdf)
	}

	video, err := newVideoGenerator()
	if err != nil {
		slog.Error("Error while trying to create videoGenerator", "error", err)
	} else {
		generators = append(generators, video)
	}

	thumbnailGenerator := &ThumbnailGenerator{
		tempFilesCache: make(map[string]string),
		tempDirectory:  tmp,
		generators:     generators,
	}

	return thumbnailGenerator, nil
}

func (g *ThumbnailGenerator) SupportsExt(ext string) bool {
	for i := range g.generators {
		if g.generators[i].supportsExt(ext) {
			return true
		}
	}

	return false
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

func (g *ThumbnailGenerator) generateThumbnail(path string) (string, error) {
	for index := range g.generators {
		generator := g.generators[index]

		if !generator.supportsExt(filepath.Ext(path)) {
			continue
		}

		generatedThumbnailPath, err := generator.generateThumbnail(path, g.tempDirectory)
		if err != nil {
			return "", err
		}

		return generatedThumbnailPath, nil
	}

	return "", errors.New("unsupported file format")
}

func (g *ThumbnailGenerator) CleanUp() error {
	return os.RemoveAll(g.tempDirectory)
}

func isPoppolerInstalled() bool {
	_, err := exec.LookPath("pdftoppm")
	return err == nil
}

func isFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}
