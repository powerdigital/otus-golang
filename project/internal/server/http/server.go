package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

const folder = "/tmp/"

type Server struct {
	server *http.Server
}

type requestDto struct {
	Width  uint
	Height uint
	Path   string
}

type RequestHandler struct{}

type Application interface {
	ResizeImage(imagePath string) error
}

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

	fileBytes, err := ioutil.ReadFile(fileDest)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(fileBytes)
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
		w.Write([]byte("incorrect parameters provided"))
		return
	}

	urlData, err := url.Parse(dto.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("wrong image path provided"))
		return
	}

	scheme := urlData.Scheme
	if len(scheme) == 0 {
		scheme = "https"
	}

	fileSource := fmt.Sprintf("%s://%s", scheme, strings.TrimLeft(urlData.Path, "/"))
	file, err := http.Get(fileSource)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	img, err := jpeg.Decode(file.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	file.Body.Close()

	thumbnail := resize.Thumbnail(dto.Width, dto.Height, img, resize.Lanczos3)

	filePath := strings.Split(dto.Path, "/")
	filename := filePath[len(filePath)-1]

	fileDest = folder + filename
	out, err := os.Create(fileDest)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer out.Close()

	jpeg.Encode(out, thumbnail, nil)

	return fileDest, err
}
