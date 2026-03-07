package utils

import (
	"bytes"
	"testing"
)

func TestValidateImageMagicNumber_JPEG(t *testing.T) {
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46}
	reader := bytes.NewReader(jpegHeader)
	
	err := ValidateImageMagicNumber(reader, ".jpg")
	if err != nil {
		t.Errorf("valid JPEG should pass validation: %v", err)
	}
}

func TestValidateImageMagicNumber_PNG(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	reader := bytes.NewReader(pngHeader)
	
	err := ValidateImageMagicNumber(reader, ".png")
	if err != nil {
		t.Errorf("valid PNG should pass validation: %v", err)
	}
}

func TestValidateImageMagicNumber_GIF(t *testing.T) {
	gifHeader := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}
	reader := bytes.NewReader(gifHeader)
	
	err := ValidateImageMagicNumber(reader, ".gif")
	if err != nil {
		t.Errorf("valid GIF should pass validation: %v", err)
	}
}

func TestValidateImageMagicNumber_Invalid(t *testing.T) {
	invalidHeader := []byte{0x00, 0x00, 0x00, 0x00}
	reader := bytes.NewReader(invalidHeader)
	
	err := ValidateImageMagicNumber(reader, ".jpg")
	if err == nil {
		t.Error("invalid image should fail validation")
	}
}

func TestValidateImageMagicNumber_FakeJPEG(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	reader := bytes.NewReader(pngHeader)
	
	err := ValidateImageMagicNumber(reader, ".jpg")
	if err == nil {
		t.Error("PNG file with .jpg extension should fail validation")
	}
}

func TestValidateZipMagicNumber_Valid(t *testing.T) {
	zipHeader := []byte{0x50, 0x4B, 0x03, 0x04}
	reader := bytes.NewReader(zipHeader)
	
	err := ValidateZipMagicNumber(reader)
	if err != nil {
		t.Errorf("valid ZIP should pass validation: %v", err)
	}
}

func TestValidateZipMagicNumber_ValidEmpty(t *testing.T) {
	emptyZipHeader := []byte{0x50, 0x4B, 0x05, 0x06}
	reader := bytes.NewReader(emptyZipHeader)
	
	err := ValidateZipMagicNumber(reader)
	if err != nil {
		t.Errorf("valid empty ZIP should pass validation: %v", err)
	}
}

func TestValidateZipMagicNumber_Invalid(t *testing.T) {
	invalidHeader := []byte{0x00, 0x00, 0x00, 0x00}
	reader := bytes.NewReader(invalidHeader)
	
	err := ValidateZipMagicNumber(reader)
	if err == nil {
		t.Error("invalid ZIP should fail validation")
	}
}

func TestValidateZipMagicNumber_FakeZip(t *testing.T) {
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	reader := bytes.NewReader(jpegHeader)
	
	err := ValidateZipMagicNumber(reader)
	if err == nil {
		t.Error("JPEG file should fail ZIP validation")
	}
}

func TestValidateImageMagicNumber_UnknownExtension(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x00}
	reader := bytes.NewReader(data)
	
	err := ValidateImageMagicNumber(reader, ".unknown")
	if err != nil {
		t.Errorf("unknown extension should skip validation: %v", err)
	}
}
