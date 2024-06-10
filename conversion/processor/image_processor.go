package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
)

type ImageProcessor struct {
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *ImageProcessor) Base64Thumbnail(inputPath string) (*client.ImageBase64, error) {
	inputSize, err := p.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.ResizeImage(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return nil, err
			}
		} else {
			if err := p.ResizeImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return nil, err
			}
		}
		b64, err := helper.ImageToBase64(outputPath)
		if err != nil {
			return nil, err
		}
		size, err := p.MeasureImage(outputPath)
		if err != nil {
			return nil, err
		}
		return &client.ImageBase64{
			Base64: *b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	} else {
		b64, err := helper.ImageToBase64(inputPath)
		if err != nil {
			return nil, err
		}
		size, err := p.MeasureImage(inputPath)
		if err != nil {
			return nil, err
		}
		return &client.ImageBase64{
			Base64: *b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	}
}

func (p *ImageProcessor) MeasureImage(inputPath string) (*client.ImageProps, error) {
	size, err := infra.NewCommand().ReadOutput("identify", "-format", "%w,%h", inputPath)
	if err != nil {
		return nil, err
	}
	values := strings.Split(*size, ",")
	width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
	if err != nil {
		return nil, err
	}
	return &client.ImageProps{
		Width:  width,
		Height: height,
	}, nil
}

func (p *ImageProcessor) ResizeImage(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	if err := infra.NewCommand().Exec("convert", "-resize", widthStr+"x"+heightStr, inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailFromImage(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	if err := infra.NewCommand().Exec("convert", "-thumbnail", widthStr+"x"+heightStr, "-background", "white", "-alpha", "remove", "-flatten", fmt.Sprintf("%s[0]", inputPath), outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ConvertImage(inputPath string, outputPath string) error {
	if err := infra.NewCommand().Exec("convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
	if err := infra.NewCommand().Exec("convert", inputPath, "-alpha", "off", outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) DPIFromImage(inputPath string) (*int, error) {
	output, err := infra.NewCommand().ReadOutput("exiftool", "-S", "-s", "-ImageWidth", "-ImageHeight", "-XResolution", "-YResolution", "-ResolutionUnit", inputPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(*output, "\n")
	if len(lines) < 5 || lines[4] != "inches" {
		return helper.ToPtr(72), nil
	}
	xRes, err := strconv.ParseFloat(lines[2], 64)
	if err != nil {
		return nil, err
	}
	yRes, err := strconv.ParseFloat(lines[3], 64)
	if err != nil {
		return nil, err
	}
	return helper.ToPtr(int((xRes + yRes) / 2)), nil
}
