package imgproc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestImage(t *testing.T) []byte {
	data, err := os.ReadFile("testdata/test.jpg")
	require.NoError(t, err)
	return data
}

func TestProcessImage_ValidJPEG_Thumb(t *testing.T) {
	dir := t.TempDir()
	resizer, err := NewResizer(dir)
	require.NoError(t, err)

	imageBytes := loadTestImage(t)
	filename, err := resizer.ProcessImage(imageBytes, "thumb", "abc123")

	require.NoError(t, err)
	assert.Equal(t, "abc123_thumb.webp", filename)
	assert.FileExists(t, filepath.Join(dir, filename))
}

func TestProcessImage_ValidJPEG_Medium(t *testing.T) {
	dir := t.TempDir()
	resizer, err := NewResizer(dir)
	require.NoError(t, err)

	imageBytes := loadTestImage(t)
	filename, err := resizer.ProcessImage(imageBytes, "medium", "abc123")

	require.NoError(t, err)
	assert.Equal(t, "abc123_medium.webp", filename)
	assert.FileExists(t, filepath.Join(dir, filename))
}

func TestProcessImage_ValidJPEG_Large(t *testing.T) {
	dir := t.TempDir()
	resizer, err := NewResizer(dir)
	require.NoError(t, err)

	imageBytes := loadTestImage(t)
	filename, err := resizer.ProcessImage(imageBytes, "large", "abc123")

	require.NoError(t, err)
	assert.Equal(t, "abc123_large.webp", filename)
	assert.FileExists(t, filepath.Join(dir, filename))
}

func TestProcessImage_CorruptBytes(t *testing.T) {
	dir := t.TempDir()
	resizer, err := NewResizer(dir)
	require.NoError(t, err)

	_, err = resizer.ProcessImage([]byte("not an image"), "thumb", "bad")

	require.Error(t, err)
}

func TestProcessImage_EmptyBytes(t *testing.T) {
	dir := t.TempDir()
	resizer, err := NewResizer(dir)
	require.NoError(t, err)

	_, err = resizer.ProcessImage([]byte{}, "thumb", "empty")

	require.Error(t, err)
}
