package main

import (
	"fmt"
	"runtime"
	"tooling"
)

func Construct(destFolder string) (string, string) {
	var dest string
	var file string

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		dest = "C:\\Tooling\\" + destFolder + "\\"
		file = dest + "quickView.exe"
	case "linux":
		dest = "/Tooling/" + destFolder + "/"
		file = dest + "quickview"
	case "darwin":
		dest = "/Tooling/" + destFolder + "/"
		file = dest + "QuickView.app"
	default:
		dest = ""
		file = ""
	}

	return dest, file
}

func GetDownloadURL() (url string) {
	var apiUrl string = "https://api.github.com/repos/Implycitt/quickView/releases/latest"

	release, err := tooling.GetRecentTag(apiUrl)
	tooling.Check(err)
	tagName := release.TagName[1:]

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s-win.zip", tagName)
	case "linux":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s.zip", tagName)
	case "darwin":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s-arm64-mac.zip", tagName)
	default:
		url = ""
	}
	return url
}
