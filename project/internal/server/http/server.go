package internalhttp

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

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

// TODO: move to config.
const (
	folder     = "/tmp"
	fileExtJpg = "jpg"
)

type Server struct {
	server *http.Server
}

type requestDto struct {
	Width  uint
	Height uint
	Path   string
}

type RequestHandler struct{}

func NewServer() *Server {
	return &Server{
		server: &http.Server{
			Addr:    ":8888",
			Handler: createHandler(),
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("http server starting error: %w", err)
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	return nil
}

func createHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/resize/{width}/{height}/{path:.+}", ResizeImage).Methods("GET")

	return router
}

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	fileDest, err := uploadRemoteFile(w, r)
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

func getRequestFileHash(dto requestDto) string {
	filepath := fmt.Sprintf("%d-%d-%s", dto.Width, dto.Height, dto.Path)
	hash := sha256.Sum256([]byte(filepath))
	return hex.EncodeToString(hash[:])
}

func getRequestDto(r *http.Request) (*requestDto, error) {
	vars := mux.Vars(r)

	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return nil, err
	}

	path := vars["path"]
	if len(path) == 0 {
		return nil, errors.New("empty file path")
	}

	return &requestDto{
		Width:  uint(width),
		Height: uint(height),
		Path:   path,
	}, nil
}

func uploadRemoteFile(w http.ResponseWriter, r *http.Request) (fileDest string, err error) {
	dto, err := getRequestDto(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fileHash := getRequestFileHash(*dto)
	fileDest = fmt.Sprintf("%s/%s.%s", folder, fileHash, fileExtJpg)
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

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlData.String(), nil)
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
