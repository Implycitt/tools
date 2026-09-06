package template

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"tooling"
)

func GrabTemplate(url string, path string) error {
	tmpDir, err := os.MkdirTemp("", "tmpClone")
	tooling.Panic(err)
	defer os.RemoveAll(tmpDir)

	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse", url, tmpDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %v: %s", err, string(output))
	}

	sparseCmd := exec.Command("git", "-C", tmpDir, "sparse-checkout", "set", path)
	if output, err := sparseCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git sparse-checkout failed: %v: %s", err, string(output))
	}

	currDir, err := os.Getwd()
	tooling.Panic(err)

	srcPath := filepath.Join(tmpDir, path)
	destPath := filepath.Join(currDir, filepath.Base(path))

	return copyDir(srcPath, destPath)
}

func copyDir(src string, destination string) error {
	info, err := os.Stat(src)
	tooling.Panic(err)

	if err := os.MkdirAll(destination, info.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	tooling.Panic(err)

	for _, entry := range entries {
		sourceFilePath := filepath.Join(src, entry.Name())
		destinationFilePath := filepath.Join(destination, entry.Name())

		if entry.IsDir() {
			if err := copyDir(sourceFilePath, destinationFilePath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourceFilePath, destinationFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src string, destination string) error {
	in, err := os.Open(src)
	tooling.Panic(err)
	defer in.Close()

	out, err := os.Create(destination)
	tooling.Panic(err)
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
