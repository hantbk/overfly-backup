package decompressor

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "decompressor_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock compressed file
	mockCompressedFile := filepath.Join(tempDir, "test.tar.gz")
	err = createMockCompressedFile(mockCompressedFile)
	assert.NoError(t, err)

	// Run the decompressor
	err = Run(mockCompressedFile, "test")
	assert.NoError(t, err)

	// Check if the original file was removed
    _, err = os.Stat(mockCompressedFile)
    assert.True(t, os.IsNotExist(err), "Original compressed file should be removed")

    // Check if the model directory exists
    modelDir := filepath.Join(tempDir, "test")
    _, err = os.Stat(modelDir)
    assert.NoError(t, err, "Model directory should exist")

    // Check if the extracted file exists in the model directory
    extractedFile := filepath.Join(modelDir, "sample.txt")
    _, err = os.Stat(extractedFile)
    assert.NoError(t, err, "Extracted file should exist in the model directory")
}

func createMockCompressedFile(filePath string) error {
    // Create a temporary directory for the model
    tempDir := filepath.Join(filepath.Dir(filePath), "decompressor_test")
    if err := os.MkdirAll(tempDir, 0755); err != nil {
        return err
    }
    defer os.RemoveAll(tempDir)

    // Create the model directory (test)
    modelDir := filepath.Join(tempDir, "test")
    if err := os.MkdirAll(modelDir, 0755); err != nil {
        return err
    }

    // Create a sample file
    sampleFile := filepath.Join(modelDir, "sample.txt")
    if err := os.WriteFile(sampleFile, []byte("Sample content"), 0644); err != nil {
        return err
    }

    // Create archive.tar inside the model directory
    archiveTar := filepath.Join(modelDir, "archive.tar")
    cmd := exec.Command("tar", "-cvf", archiveTar, "-C", modelDir, "sample.txt")
    if err := cmd.Run(); err != nil {
        return err
    }

    // Create the final tar.gz file
    cmd = exec.Command("tar", "-czvf", filePath, "-C", tempDir, "test")
    return cmd.Run()
}
