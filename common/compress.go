package common

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/eolinker/eosc/log"
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

func UnzipFromBytes(packageContent []byte, dest string) error {
	// 通过字节流创建zip的Reader对象
	zr, err := zip.NewReader(bytes.NewReader(packageContent), int64(len(packageContent)))
	if err != nil {
		return err
	}

	// 解压
	return Unzip(zr, dest)
}

func Unzip(src *zip.Reader, dst string) error {
	// 强制转换一遍目录
	dst = filepath.Clean(dst)

	// 遍历压缩文件
	for _, file := range src.File {
		// 在闭包中完成以下操作可以及时释放文件句柄
		err := func() error {
			// 跳过文件夹
			if file.Mode().IsDir() {
				return nil
			}
			decodeName := file.Name
			//if file.Flags == 0 {
			//	//如果标致位是0  则是默认的本地编码   默认为gbk
			//	i := bytes.NewReader([]byte(file.Name))
			//	decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			//	content, _ := io.ReadAll(decoder)
			//	decodeName = string(content)
			//} else {
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			//decodeName = file.Name
			//}
			// 配置输出目标路径
			filename := filepath.Join(dst, decodeName)
			// 创建目标路径所在文件夹
			e := os.MkdirAll(filepath.Dir(filename), 0755)
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

	// 解压完成
	return nil
}
