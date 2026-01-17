package main

import (
	"bytes"
	"io"
	"os"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"strings"
)

// Build the binary once for all tests
func TestMain(m *testing.M) {
	// Build shred binary
	if err := exec.Command("go", "build", "-o", "./shred").Run(); err != nil {
		fmt.Printf("Failed to build shred: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove("./shred") // Cleanup after all tests
	
	// Run tests
	os.Exit(m.Run())
}



func TestShredBasic(t *testing.T) {
	dir := "/tmp"
	fpath := filepath.Join(dir, "test")
	data := []byte("secret")
	
	if err := os.WriteFile(fpath, data, 0644); err != nil {
		t.Fatal(err)
	}
	
	// Shred with 1 pass, remove
	cmd := exec.Command("./shred", "-n", "1", "-u", fpath)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	
	// Verify removed
	if _, err := os.Stat(fpath); !os.IsNotExist(err) {
		t.Error("file not removed")
	}
}

func TestShredContentOverwritten(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "test")
	orig := []byte("hello123")
	
	os.WriteFile(fpath, orig, 0644)
	
	// Shred 1 pass, no remove
	exec.Command("./shred", "-n", "1", fpath).Run()
	
	f, err := os.Open(fpath)
	if err != nil {
		t.Fatal(err)
	}
	newdata, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	
	// Verify no original data remains
	if bytes.Contains(newdata, orig) {
		t.Errorf("original data still present: %s", newdata[:10])
	}
}

func TestShredZeroPass(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "test")
	
	os.WriteFile(fpath, []byte("data"), 0644)
	
	exec.Command("./shred", "-n", "1", "-z", fpath).Run()
	
	f, _ := os.Open(fpath)
	data, _ := io.ReadAll(f)
	f.Close()
	
	// Should be all zeros
	for _, b := range data {
		if b != 0 {
			t.Errorf("not zero: %v", data[:10])
			break
		}
	}
}

func TestShredDir(t *testing.T) {
	dir := t.TempDir()
	
	cmd := exec.Command("./shred", dir)
	output, err := cmd.CombinedOutput()

	if (err == nil){
		t.Fatal("'shred' should return error on directory.")
	}

	errorOutput := string(output)
	if !strings.Contains(errorOutput, "is a directory") {
		t.Errorf("Expected 'is a directory' error, got: %s", errorOutput)
	}

	exitErr := err.(*exec.ExitError)
	
	if exitErr == nil{
		t.Fatalf("Expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Errorf("Expected non-zero exit code, got %d", exitErr.ExitCode())
	}
}

func TestShredSpecialFile(t *testing.T) {
	
	cmd := exec.Command("./shred", "/dev/null")
	output, err := cmd.CombinedOutput()

	if (err == nil){
		t.Fatal("'shred' should return error on special file.")
	}

	errorOutput := string(output)
	if !strings.Contains(errorOutput, "is not a regular file") {
		t.Errorf("Expected 'is not a regular file' error, got: %s", errorOutput)
	}
	exitErr := err.(*exec.ExitError)
	
	if exitErr == nil{
		t.Fatalf("Expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() == 0 {
		t.Errorf("Expected non-zero exit code, got %d", exitErr.ExitCode())
	}
}


