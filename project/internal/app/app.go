package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/powerdigital/project/internal/config"
	"github.com/rs/zerolog"
)

const fileExtJpg = "jpg"

type uriPathDto struct {
	Width  uint
	Height uint
	Path   string
}

type App struct {
	logger zerolog.Logger
	config config.Config
}

func NewApp(logger zerolog.Logger, config config.Config) App {
	return App{
		logger: logger,
		config: config,
	}
}

func (app App) ResizeImage(w http.ResponseWriter, r *http.Request) {
	fileDest, err := uploadRemoteFile(w, r, app.config)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	fileBytes, err := os.ReadFile(fileDest)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(fileBytes)
}

func uploadRemoteFile(w http.ResponseWriter, r *http.Request, config config.Config) (fileDest string, err error) {
	dto, err := getRequestDto(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fileHash := getRequestFileHash(*dto)
	fileDest = fmt.Sprintf("%s/%s.%s", config.CacheFolder, fileHash, fileExtJpg)
	fileBytes, err := os.ReadFile(fileDest)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
		return
	}

	urlData, err := url.Parse(dto.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("wrong image path provided"))
		return
	}

	if len(urlData.Scheme) == 0 {
		urlData.Scheme = "https"
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlData.String(), nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	req.Header = r.Header

	cli := &http.Client{}
	file, err := cli.Do(req)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer file.Body.Close()

	img, err := jpeg.Decode(file.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	thumbnail := resize.Thumbnail(dto.Width, dto.Height, img, resize.Lanczos3)

	out, err := os.Create(fileDest)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer out.Close()

	jpeg.Encode(out, thumbnail, nil)

	return fileDest, err
}

func getRequestDto(r *http.Request) (*uriPathDto, error) {
	params := strings.Split(r.RequestURI, "/")

	width, err := strconv.Atoi(params[2])
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(params[3])
	if err != nil {
		return nil, err
	}

	path := strings.Join(params[4:], "/")
	if len(path) == 0 {
		return nil, errors.New("empty file path")
	}

	return &uriPathDto{
		Width:  uint(width),
		Height: uint(height),
		Path:   path,
	}, nil
}

func getRequestFileHash(dto uriPathDto) string {
	filepath := fmt.Sprintf("%d-%d-%s", dto.Width, dto.Height, dto.Path)
	hash := sha256.Sum256([]byte(filepath))
	return hex.EncodeToString(hash[:])
}
