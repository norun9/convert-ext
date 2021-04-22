package cvt

import (
	"fmt"
	"github.com/first-task/pkg/errof"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ImageCvtGlue struct {
	InputDir  string
	OutputDir string
	BeforeExt string
	AfterExt  string
}

func NewImageCvtGlue(
	inputDir,
	outputDir,
	beforeExt,
	afterExt string,
) *ImageCvtGlue {
	if strings.Index(beforeExt, ".") == -1 {
		beforeExt = fmt.Sprintf(".%s", beforeExt)
	}
	if strings.Index(afterExt, ".") == -1 {
		afterExt = fmt.Sprintf(".%s", afterExt)
	}
	return &ImageCvtGlue{
		InputDir:  inputDir,
		OutputDir: outputDir,
		BeforeExt: beforeExt,
		AfterExt:  afterExt,
	}
}

func (c *ImageCvtGlue) Exec() (err error) {
	var srcPaths []string
	if srcPaths, err = c.pathWalk(); err != nil {
		return err
	}
	log.Println("SrcPath:", srcPaths)
	if err = c.convert(srcPaths); err != nil {
		return err
	}
	return nil
}

func (c *ImageCvtGlue) convert(srcPaths []string) (err error) {
	for _, srcPath := range srcPaths {
		var file *os.File
		if file, err = os.Open(srcPath); err != nil {
			return errors.Wrap(errof.ErrOpenSrcFile, err.Error())
		}
		img, _, err := image.Decode(file)

		dstPath := c.getDstPath(file.Name())

		var dst *os.File
		if dst, err = os.Create(dstPath); err != nil {
			return errors.Wrap(errof.ErrCreateDstFile, err.Error())
		}

		switch filepath.Ext(dstPath) {
		case ".png":
			if err = png.Encode(dst, img); err != nil {
				return errors.Wrap(errof.ErrEncodePngImg, err.Error())
			}
		case ".jpg", ".jpeg":
			if err = jpeg.Encode(dst, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
				return errors.Wrap(errof.ErrEncodeJpgImg, err.Error())
			}
		case ".gif":
			if err = gif.Encode(dst, img, nil); err != nil {
				return errors.Wrap(errof.ErrEncodeGifImg, err.Error())
			}
		}

		// image.Imageへとデコード
		if err = file.Close(); err != nil {
			return errors.Wrap(errof.ErrCloseSrcFile, err.Error())
		}
		if err = dst.Close(); err != nil {
			return errors.Wrap(errof.ErrCloseSrcFile, err.Error())
		}
	}
	return nil
}

func (c *ImageCvtGlue) pathWalk() (srcPaths []string, err error) {
	rootDir := getRootDir()
	if err = filepath.Walk(filepath.Join(rootDir, c.InputDir), func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(errof.ErrWalkingSrcPath, err.Error())
		}
		// ファイルが存在しているパスかどうかを確認
		var fileExists bool
		if fileExists, err = isFileExist(srcPath); err != nil {
			return err
		}
		// ファイルかつ指定した拡張子であれば配列に格納
		if fileExists && filepath.Ext(srcPath) == c.BeforeExt {
			srcPaths = append(srcPaths, srcPath)
		}
		return nil
	}); err != nil {
		return srcPaths, err
	}
	return srcPaths, nil
}

func (c *ImageCvtGlue) getDstPath(path string) string {
	fileDir := filepath.Dir(path)
	fileNameWithoutExt := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
	return fmt.Sprintf("%s%s", filepath.Join(fileDir, fileNameWithoutExt), c.AfterExt)
}

func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	cvtDir := filepath.Dir(b)
	internalDir := filepath.Dir(cvtDir)
	return filepath.Dir(internalDir)
}

func isFileExist(filepath string) (isFile bool, err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(errof.ErrGetSrcFileInfo, err.Error())
	}
	if !fileInfo.IsDir() {
		return true, nil
	}
	return false, nil
}
