package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	imageMagicNumbers = map[string][][]byte{
		".jpg":  {{0xFF, 0xD8, 0xFF}},
		".jpeg": {{0xFF, 0xD8, 0xFF}},
		".png":  {{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
		".gif":  {{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, {0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
		".webp": {{0x52, 0x49, 0x46, 0x46}},
	}
	
	zipMagicNumbers = [][]byte{
		{0x50, 0x4B, 0x03, 0x04},
		{0x50, 0x4B, 0x05, 0x06},
		{0x50, 0x4B, 0x07, 0x08},
	}
)

func ValidateImageMagicNumber(file io.Reader, ext string) error {
	magicNumbers, ok := imageMagicNumbers[ext]
	if !ok {
		return nil
	}

	buf := make([]byte, 16)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]

	if len(buf) == 0 {
		return errors.New("empty file")
	}

	for _, magic := range magicNumbers {
		if bytes.HasPrefix(buf, magic) {
			return nil
		}
	}

	return fmt.Errorf("file content does not match the declared image type (ext: %s, header: %x)", ext, buf[:min(8, len(buf))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ValidateZipMagicNumber(file io.Reader) error {
	buf := make([]byte, 16)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	buf = buf[:n]

	for _, magic := range zipMagicNumbers {
		if bytes.HasPrefix(buf, magic) {
			return nil
		}
	}

	return errors.New("file content does not match zip format")
}

type resetReader struct {
	io.Reader
	buffer *bytes.Buffer
}

func NewResetReader(r io.Reader) *resetReader {
	return &resetReader{
		Reader: r,
		buffer: &bytes.Buffer{},
	}
}

func (r *resetReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		r.buffer.Write(p[:n])
	}
	return n, err
}

func (r *resetReader) Reset() io.Reader {
	return io.MultiReader(r.buffer, r.Reader)
}
