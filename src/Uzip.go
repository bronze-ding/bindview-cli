package src

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(Path string, Zip *zip.Reader, newName string) error {
	// 强制转换一遍目录
	dst := filepath.Clean(Path)

	// 遍历压缩文件
	for _, file := range Zip.File {
		// 在闭包中完成以下操作可以及时释放文件句柄
		err := func() error {
			// 跳过文件夹
			if file.Mode().IsDir() {
				return nil
			}

			// 配置输出目标路径
			filename := filepath.Join(dst, file.Name)
			// 创建目标路径所在文件夹
			e := os.MkdirAll(filepath.Dir(filename), os.ModeDir)
			if e != nil {
				return e
			}

			// 打开这个压缩文件
			zfr, e := file.Open()
			if e != nil {
				return e
			}
			defer zfr.Close()

			// 创建目标文件
			fw, e := os.Create(filename)
			if e != nil {
				return e
			}
			defer fw.Close()

			// 执行拷贝
			_, e = io.Copy(fw, zfr)
			if e != nil {
				return e
			}

			// 拷贝成功
			return nil
		}()

		// 是否发生异常
		if err != nil {
			return err
		}
	}

	// 重命名
	NewFileName(Zip.File[0].Name, newName)
	// 解压完成
	return nil
}
