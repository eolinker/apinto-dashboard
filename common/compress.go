package common

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/eolinker/eosc/log"
	"io"
	"os"
	"path"
)

// DeCompress 解压 tar.gz
func DeCompress(srcFile io.Reader, dest string) error {
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filePath := path.Join(dest, hdr.Name)
		// 根据文件类型进行不同的处理
		switch hdr.Typeflag {
		case tar.TypeDir:
			// 如果是目录，创建目录
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				log.Error("安装插件失败, 无法创建目录:", err)
				return err
			}
		case tar.TypeReg:
			// 如果是普通文件，创建文件并写入内容
			file, err := os.Create(filePath)
			if err != nil {
				log.Error("安装插件失败, 无法创建文件:", err)
				return err
			}
			defer file.Close()

			_, err = io.Copy(file, tr)
			if err != nil {
				log.Error("安装插件失败, 无法写入文件内容:", err)
				return err
			}
		default:
			log.Errorf("未知文件类型: %s in %s\n", hdr.Typeflag, hdr.Name)
			return fmt.Errorf("未知文件类型: %s in %s\n", hdr.Typeflag, hdr.Name)
		}
	}
	return nil
}
