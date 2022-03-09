package untar

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

func UntarFile(srcFile string, dstDir string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()
	return Untar(src, dstDir)
}

func Untar(reader io.Reader, dstDir string) error {
	mimetype.SetLimit(4096)
	buf := make([]byte, 0, 4096)
	buffer := bytes.NewBuffer(buf)
	tee := io.TeeReader(reader, buffer)
	mime, err := mimetype.DetectReader(tee)
	if err != nil && err != io.EOF {
		return err
	}
	mr := io.MultiReader(buffer, reader)
	if mime.Is("application/x-bzip2") {
		return untarBzip(mr, dstDir)
	} else if mime.Is("application/gzip") {
		return untarGzip(mr, dstDir)
	} else {
		return untar(mr, dstDir)
	}
}

func untarBzip(reader io.Reader, dstDir string) error {
	bzipReader := bzip2.NewReader(reader)
	return untar(bzipReader, dstDir)
}

func untarGzip(reader io.Reader, dstDir string) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	return untar(gzipReader, dstDir)
}

func untar(r io.Reader, dst string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			prevDir, _ := filepath.Split(target)
			if _, err := os.Stat(prevDir); os.IsNotExist(err) {
				if err := os.MkdirAll(prevDir, 0755); err != nil {
					return err
				}
			}
			if err = func() error {
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := io.Copy(f, tr); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
	}
}
