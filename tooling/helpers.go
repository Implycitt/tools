package tooling

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Check(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func DownloadFile(filePath string, url string) (err error) {
	var file = "quickView.zip"
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		panic(err)
	}

	out, err := os.Create(filePath + file)
	Check(err)
	defer out.Close()

	resp, err := http.Get(url)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	Check(err)

	return nil
}

func Unzip(file string, destFolder string) {
	archive, err := zip.OpenReader(destFolder + file)
	Panic(err)
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destFolder, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destFolder)+string(os.PathSeparator)) {
			fmt.Println("Invalid Path")
			return
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		Panic(err)
		defer destinationFile.Close()

		fileInArchive, err := f.Open()
		Panic(err)
		defer fileInArchive.Close()

		if _, err := io.Copy(destinationFile, fileInArchive); err != nil {
			panic(err)
		}
	}

}

func ConstructPath(destFolder string) (destination string) {
	var dest string

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		dest = "C:\\Tooling\\" + destFolder + "\\"
	case "linux":
		dest = "/Tooling/" + destFolder + "/"
	case "darwin":
		dest = "/Tooling/" + destFolder + "/"
	default:
		dest = ""
	}

	return dest
}

func ClearPath(destination string) (err error) {
	err = os.RemoveAll(destination)
	Check(err)

	err = os.MkdirAll(destination, 0755)
	return err
}

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
