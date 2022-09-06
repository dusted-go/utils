package webfile

import (
	"crypto/sha256"
	"encoding/hex"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/dusted-go/fault/fault"
)

// MimeType returns the media type of a multipart.File object.
func MimeType(f multipart.File, h *multipart.FileHeader) (string, error) {
	mimeType := mime.TypeByExtension(filepath.Ext(h.Filename))
	if len(mimeType) > 0 {
		return mimeType, nil
	}

	buff := make([]byte, 512)
	_, err := f.Read(buff)
	if err != nil {
		return "", fault.SystemWrap(err, "webfile", "MimeType", "error reading the first 512 bytes")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return "", fault.SystemWrap(err, "webfile", "MimeType", "error seeking file")
	}

	return http.DetectContentType(buff), nil
}

// Hash computes a hash of a given multipart.File.
// Apply an optional salt to create a unique hash if needed.
func Hash(f multipart.File, h *multipart.FileHeader, salt []byte) (string, error) {
	buff := make([]byte, h.Size)
	_, err := f.Read(buff)
	if err != nil {
		return "", fault.SystemWrap(err, "webfile", "Hash", "error reading file")
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return "", fault.SystemWrap(err, "webfile", "Hash", "error seeking file")
	}

	buff = append(buff, salt...)
	hasher := sha256.New()
	_, err = hasher.Write(buff)
	if err != nil {
		return "", fault.SystemWrap(err, "webfile", "Hash", "error hashing file")
	}
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}
