package filepreview

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yorukot/superfile/src/internal/common"
)

type thumbnailGeneratorInterface interface {
	supportsExt(ext string) bool
	generateThumbnail(inputPath string, outputPathWithoutExt string) (string, error)
}

type VideoGenerator struct {
	ffprobeAvailable bool
}

func newVideoGenerator() (*VideoGenerator, error) {
	if !isFFmpegInstalled() {
		return nil, errors.New("ffmpeg is not installed")
	}

	return &VideoGenerator{ffprobeAvailable: isFFprobeInstalled()}, nil
}

func (g *VideoGenerator) supportsExt(ext string) bool {
	return common.VideoExtensions[strings.ToLower(ext)]
}

type videoMetadata struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func getVideoDuration(ctx context.Context, inputPath string) (float64, error) {
	probeCtx, cancel := context.WithTimeout(ctx, videoProbeTimeout)
	defer cancel()

	ffprobe := exec.CommandContext(probeCtx, "ffprobe",
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "json=c=1",
		inputPath,
	)

	output, err := ffprobe.Output()
	if err != nil {
		return 0, fmt.Errorf("error probing video duration: %w", err)
	}

	var metadata videoMetadata
	err = json.Unmarshal(output, &metadata)
	if err != nil {
		return 0, fmt.Errorf("error parsing ffprobe output: %w", err)
	}

	if metadata.Format.Duration == "" {
		return 0, errors.New("ffprobe output missing duration")
	}

	duration, err := strconv.ParseFloat(metadata.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing video duration %q: %w", metadata.Format.Duration, err)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid video duration: %f", duration)
	}

	return duration, nil
}

func formatVideoThumbnailSeekSeconds(seekSeconds float64) string {
	return strconv.FormatFloat(seekSeconds, 'f', videoThumbSeekTimestampPrecision, 64)
}

func videoThumbnailScaleFilter() string {
	return fmt.Sprintf(
		"scale='min(%d,iw)':'min(%d,ih)':force_original_aspect_ratio=decrease:flags=fast_bilinear",
		videoThumbMaxWidth,
		videoThumbMaxHeight,
	)
}

func runVideoThumbnailFFmpeg(ctx context.Context, inputPath string, outputPath string, seekSeconds float64) error {
	ffmpegArgs := []string{
		"-v", "warning", // set log level to warning
		"-hwaccel", "auto", // Use Hardware Acceleration if available
		"-threads", "1", // limit CPU spikes from video decoding
		"-an", // disable Audio stream
		"-sn", // disable Subtitle stream
		"-dn", // disable data stream
	}

	if seekSeconds > 0 {
		ffmpegArgs = append(ffmpegArgs, "-ss", formatVideoThumbnailSeekSeconds(seekSeconds))
	}

	ffmpegArgs = append(ffmpegArgs,
		"-skip_frame", "nokey", // skip non-key frames
		"-i", inputPath, // set input file
		"-vframes", "1", // output only one frame (one image)
		"-q:v", strconv.Itoa(videoThumbQuality), // set JPEG quality
		"-vf", videoThumbnailScaleFilter(), // scale down for terminal preview
		"-f", "image2", // set format to image2
		"-fs", maxVideoFileSizeForThumb, // limit the max file size to match image previewer limit
		"-y", outputPath, // set the outputFile and overwrite it without confirmation if already exists
	)

	ffmpeg := exec.CommandContext(ctx, "ffmpeg", ffmpegArgs...)
	if err := ffmpeg.Run(); err != nil {
		return err
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("generated thumbnail is missing: %w", err)
	}

	if info.Size() == 0 {
		return errors.New("generated thumbnail is empty")
	}

	return nil
}

func (g *VideoGenerator) generateThumbnail(inputPath string, outputPathWithoutExt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), thumbGenerationTimeout)
	defer cancel()
	outputPath := outputPathWithoutExt + thumbOutputExt

	seekSeconds := 0.0

	if g.ffprobeAvailable {
		videoDuration, err := getVideoDuration(ctx, inputPath)
		if err != nil {
			slog.Debug("Error probing video duration, thumbnail generation starts at beginning of video", "error", err)
		}
		seekSeconds = videoDuration * videoThumbSeekRatio
	}

	err := runVideoThumbnailFFmpeg(ctx, inputPath, outputPath, seekSeconds)
	if err != nil {
		return "", fmt.Errorf("error generating video thumbnail, outputPath: %s : %w", outputPath, err)
	}

	return outputPath, nil
}

type pdfGenerator struct{}

func newPdfGenerator() (*pdfGenerator, error) {
	if !isPopplerInstalled() {
		return nil, errors.New("poppler is not installed")
	}

	return &pdfGenerator{}, nil
}

func (g *pdfGenerator) supportsExt(ext string) bool {
	return strings.ToLower(ext) == ".pdf"
}

func (g *pdfGenerator) generateThumbnail(inputPath string, outputPathWithoutExt string) (string, error) {
	outputPath := outputPathWithoutExt + thumbOutputExt
	ctx, cancel := context.WithTimeout(context.Background(), thumbGenerationTimeout)
	defer cancel()

	// pdftoppm -singlefile -png prefixFilename
	pdftoppm := exec.CommandContext(ctx, "pdftoppm",
		"-singlefile",        // output only the first page as image
		"-jpeg",              // Image extension
		inputPath,            // Set input file
		outputPathWithoutExt, // The output prefix. (pdftoppm will add the .jpg ext)
	)

	err := pdftoppm.Run()
	if err != nil {
		return "", fmt.Errorf("error generating pdf thumbnail, outputPath: %s : %w",
			outputPath, err)
	}

	return outputPath, nil
}

type psGenerator struct{}

func newPsGenerator() (*psGenerator, error) {
	if !isGhostscriptInstalled() {
		return nil, errors.New("ghostscript is not installed")
	}

	return &psGenerator{}, nil
}

func (g *psGenerator) supportsExt(ext string) bool {
	extension := strings.ToLower(ext)
	return extension == ".ps" || extension == ".eps"
}

func (g *psGenerator) generateThumbnail(inputPath string, outputPathWithoutExt string) (string, error) {
	outputPath := outputPathWithoutExt + thumbOutputExt
	ctx, cancel := context.WithTimeout(context.Background(), thumbGenerationTimeout)
	defer cancel()

	// gs -dSAFER -dBATCH -dNOPAUSE -sPageList=1 -sDEVICE=jpeg -r150 -sOutputFile=output.jpg input.ps
	outputParam := "-sOutputFile=" + outputPath
	gs := exec.CommandContext(ctx, "gs",
		"-dSAFER", "-dBATCH", "-dNOPAUSE", // Standard GS operators
		"-sPageList=1",  // Output only the first page
		"-sDEVICE=jpeg", // Output format
		"-r150",         // Resolution (the same as for pdf)
		outputParam,     // Result (variable because of golangci-lint)
		inputPath,       // Input file
	)

	err := gs.Run()
	if err != nil {
		return "", fmt.Errorf("error generating ps thumbnail, outputPath: %s : %w",
			outputPath, err)
	}

	return outputPath, nil
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
		slog.Debug("Error while trying to create pdfGenerator", "error", err)
	} else {
		generators = append(generators, pdf)
	}

	ps, err := newPsGenerator()
	if err != nil {
		slog.Debug("Error while trying to create psGenerator", "error", err)
	} else {
		generators = append(generators, ps)
	}

	video, err := newVideoGenerator()
	if err != nil {
		slog.Debug("Error while trying to create videoGenerator", "error", err)
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
	fileExt := filepath.Ext(path)
	for index := range g.generators {
		generator := g.generators[index]

		if !generator.supportsExt(fileExt) {
			continue
		}
		filename := filepath.Base(path)
		baseName := filename[:len(filename)-len(fileExt)]

		outputPathWithoutExt := filepath.Join(g.tempDirectory,
			fmt.Sprintf("%s-%d", baseName, time.Now().UnixNano()))

		outputPath, err := generator.generateThumbnail(path, outputPathWithoutExt)
		if err != nil {
			return "", err
		}

		return outputPath, nil
	}

	return "", errors.New("unsupported file format")
}

func (g *ThumbnailGenerator) CleanUp() error {
	return os.RemoveAll(g.tempDirectory)
}

func isPopplerInstalled() bool {
	_, err := exec.LookPath("pdftoppm")
	return err == nil
}

func isGhostscriptInstalled() bool {
	_, err := exec.LookPath("gs")
	return err == nil
}

func isFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func isFFprobeInstalled() bool {
	_, err := exec.LookPath("ffprobe")
	return err == nil
}
