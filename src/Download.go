package src

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
)

var File bytes.Buffer

// Download 下载
func Download(url string, name string) {
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"正在下载模板压缩包:",
	)

	io.Copy(io.MultiWriter(&File, bar), resp.Body)

	FileZip, err := io.ReadAll(&File)
	if err != nil {
		fmt.Println(err)
		return
	}

	ZipReader, err := zip.NewReader(bytes.NewReader(FileZip), int64(len(FileZip)))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解 .zip 压缩
	if err = Unzip("./", ZipReader, name); err != nil {
		fmt.Println(err)
		return
	}

}
