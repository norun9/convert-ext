package cvt

import (
	"bytes"
	"fmt"
	"github.com/first-task/pkg/errof"
	"github.com/pkg/errors"
	"image"
	"io"
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
	var imageFiles []string
	if imageFiles, err = c.pathWalk(); err != nil {
		return err
	}
	log.Println("images:", imageFiles)
	c.convert(imageFiles)
	return nil
}

func (c *ImageCvtGlue) convert(files []string) (err error) {
	for _, file := range files {
		var sf *os.File
		if sf, err = os.Open(file); err != nil {
			return errors.Wrap(errof.ErrOpenFile, err.Error())
		}
		buf := new(bytes.Buffer)
		io.Copy(buf, sf)
		img, _, err := image.Decode(sf)

		// image.Imageへとデコード
		log.Println(img)
		if err = sf.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *ImageCvtGlue) pathWalk() (output []string, err error) {
	var imageFiles []string
	rootDir := getRootDir()
	if err = filepath.Walk(filepath.Join(rootDir, c.InputDir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(errof.ErrWalkingPath, err.Error())
		}
		// ファイルが存在しているパスかどうかを確認
		var fileExists bool
		if fileExists, err = isFileExist(path); err != nil {
			return err
		}
		// ファイルかつ指定した拡張子であれば配列に格納
		if fileExists && filepath.Ext(path) == c.BeforeExt {
			imageFiles = append(imageFiles, path)
		}
		return nil
	}); err != nil {
		return []string{}, err
	}
	return imageFiles, nil
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
		if !os.IsExist(err) {
			return false, nil //Directory does not exist
		} else {
			return false, errors.Wrap(errof.ErrGetFileInfo, err.Error()) //Something wrong
		}
	}
	if !fileInfo.IsDir() {
		return true, nil //It's file
	}
	return false, nil
}
