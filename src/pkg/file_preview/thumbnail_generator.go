package filepreview

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	maxFileSize = "104857600" // 100MB limit
	outputExt   = ".jpg"
)

type ThumbnailGenerator struct {
	tempFiles     map[string]string
	tempDirectory string
}

func NewThumbnailGenerator() (*ThumbnailGenerator, error) {
	tmp, err := os.MkdirTemp("", "superfiles-*")
	if err != nil {
		return nil, err
	}

	thumbnailGenerator := &ThumbnailGenerator{
		tempFiles:     make(map[string]string),
		tempDirectory: tmp,
	}

	return thumbnailGenerator, nil
}

func (g *ThumbnailGenerator) GetThumbnailOrGenerate(path string) (string, error) {
	file, ok := g.tempFiles[path]

	if ok {
		_, err := os.Stat(file)
		if err == nil {
			return file, nil
		}

		delete(g.tempFiles, path)
	}

	generatedThumbnailPath, err := g.generateThumbnail(path)
	if err != nil {
		return "", err
	}

	g.tempFiles[path] = generatedThumbnailPath

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

	// ffmpeg -v warning -t 60 -hwaccel auto -an -sn -dn -skip_frame nokey -i input.mkv -vf scale='min(1024,iw)':'min(720,ih)':force_original_aspect_ratio=decrease:flags=fast_bilinear -vf "thumbnail" -frames:v 1 -y thumb.jpg
	ffmpeg := exec.Command("ffmpeg",
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
