package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/disintegration/imaging"
)

var errorsFileNotSupported = errors.New("This file type currently not supported")
var allowedMimes = map[string]bool{
	// image
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,

	// video
	"video/mp4":       true,
	"video/webm":      true,
	"video/quicktime": true,

	// audio
	"audio/mpeg": true,
	"audio/wav":  true,
	"audio/ogg":  true,

	// document
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

const (
	TypeAudio    = "AUDIO"
	TypeImage    = "IMAGE"
	TypeVideo    = "VIDEO"
	TypeDocument = "DOCUMENT"
)

type FileStorage struct {
	Root        string
	PublicPath  string
	PrivatePath string
}

func mustGetProjectRoot() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exe)
}

func mustCreateDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
}

func NewFileStorage() *FileStorage {
	root := mustGetProjectRoot()
	uploadsDir := filepath.Join(root, "uploads")
	privateDir := filepath.Join(uploadsDir, "private")
	publicDir := filepath.Join(uploadsDir, "public")
	mustCreateDir(uploadsDir)
	mustCreateDir(privateDir)
	mustCreateDir(publicDir)

	fileStoreage := FileStorage{
		Root:        root,
		PublicPath:  publicDir,
		PrivatePath: privateDir,
	}

	return &fileStoreage
}

func (storage *FileStorage) validateFile(f multipart.File) (string, string, error) {
	buf := make([]byte, 512)
	_, err := f.Read(buf)
	if err != nil {
		return "", "", err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return "", "", err
	}

	mimeType := http.DetectContentType(buf)

	if !allowedMimes[mimeType] {
		return "", "", errorsFileNotSupported
	}

	_, after, _ := strings.Cut(mimeType, "/")

	fileExt := "." + after

	mediaType := storage.GetMediaType(mimeType)

	return fileExt, mediaType, nil
}

func (storage *FileStorage) GetMediaType(mime string) string {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return TypeImage
	case strings.HasPrefix(mime, "video/"):
		return TypeVideo
	case strings.HasPrefix(mime, "audio/"):
		return TypeAudio
	case strings.HasPrefix(mime, "application/"):
		return TypeDocument
	default:
		return ""
	}
}

func (storage *FileStorage) compressImage(src multipart.File, dst io.Writer) error {
	img, err := imaging.Decode(src)
	if err != nil {
		return err
	}

	// resize
	img = imaging.Resize(img, 1024, 0, imaging.Lanczos)

	// encode ke JPEG (compressed)
	err = imaging.Encode(dst, img, imaging.JPEG, imaging.JPEGQuality(70))
	if err != nil {
		return err
	}

	return nil
}

func (storage *FileStorage) SavePublicFile(fileHeader *multipart.FileHeader, place ...string) (string, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	ext, mediaType, err := storage.validateFile(f)
	if err != nil {
		return "", err
	}

	filename, err := pkg.GenerateSecureString(40)
	if err != nil {
		return "", err
	}
	filename = filename + ext

	parts := []string{
		storage.PublicPath,
	}

	parts = append(parts, place...)
	parts = append(parts, filename)

	fullPath := filepath.Join(parts...)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer dst.Close()
	if mediaType == TypeImage {
		err := storage.compressImage(f, dst)
		if err != nil {
			return "", err
		}
	}
	if _, err := io.Copy(dst, f); err != nil {
		return "", err
	}
	return filename, nil
}

func (storage *FileStorage) SavePrivateFile(fileheader multipart.FileHeader, place ...string) (string, error) {

	// todo : do some shit
	// ot do me instead
	return "", nil
}
