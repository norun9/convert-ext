package cvt

import (
	"fmt"
	"github.com/first-task/pkg/errof"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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
	RemoveSrc bool
}

func NewImageCvtGlue(
	inputDir,
	outputDir,
	beforeExt,
	afterExt string,
	removeSrc bool,
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
		RemoveSrc: removeSrc,
	}
}

// Exec:
func (c *ImageCvtGlue) Exec() (err error) {
	var srcPaths []string
	if srcPaths, err = c.pathWalk(); err != nil {
		return err
	}
	if err = c.convert(srcPaths); err != nil {
		return err
	}
	return nil
}

// convert:
func (c *ImageCvtGlue) convert(srcPaths []string) (err error) {
	for _, srcPath := range srcPaths {
		var file *os.File
		if file, err = os.Open(srcPath); err != nil {
			return errors.Wrap(errof.ErrOpenSrcFile, err.Error())
		}
		// イメージのデコード
		var img image.Image
		if img, _, err = image.Decode(file); err != nil {
			return errors.Wrap(errof.ErrDecodeImage, err.Error())
		}
		var dstPath string
		if dstPath, err = c.getDstPath(file.Name()); err != nil {
			return err
		}
		// ファイルの作成
		var dst *os.File
		if dst, err = createFile(dstPath); err != nil {
			return err
		}
		//　作成したファイルにデコードしたイメージをエンコード
		if err = encodeImage(dstPath, dst, img); err != nil {
			return err
		}
		// ファイルの削除
		if c.RemoveSrc {
			if err = removeSrc(srcPath); err != nil {
				return err
			}
		}
		if err = file.Close(); err != nil {
			return errors.Wrap(errof.ErrCloseSrcFile, err.Error())
		}
		if err = dst.Close(); err != nil {
			return errors.Wrap(errof.ErrCloseSrcFile, err.Error())
		}
	}
	return nil
}

// pathWalk:
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

// getDstPath:
func (c *ImageCvtGlue) getDstPath(path string) (dstPath string, err error) {
	fileNameWithoutExt := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
	if c.BeforeExt == "" {
		fileDir := filepath.Dir(path)
		return fmt.Sprintf("%s%s", filepath.Join(fileDir, fileNameWithoutExt), c.AfterExt), nil
	}
	// 画像の出力先が指定されている場合は指定されたディレクトリを作成
	fileDir := filepath.Join(getRootDir(), c.OutputDir)
	// 指定したディレクトリが既に存在していれば新規でディレクトリを作成しない
	var isDir bool
	if isDir, err = isDirExist(fileDir); err != nil {
		return "", err
	}
	if !isDir {
		if err = createDir(fileDir); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%s", filepath.Join(fileDir, fileNameWithoutExt), c.AfterExt), nil
}

// encodeImage:
func encodeImage(dstPath string, dst *os.File, img image.Image) (err error) {
	switch filepath.Ext(dstPath) {
	case ".png":
		if err = png.Encode(dst, img); err != nil {
			return errors.Wrap(errof.ErrEncodePNGImg, err.Error())
		}
	case ".jpg", ".jpeg":
		if err = jpeg.Encode(dst, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return errors.Wrap(errof.ErrEncodeJPGImg, err.Error())
		}
	case ".gif":
		if err = gif.Encode(dst, img, nil); err != nil {
			return errors.Wrap(errof.ErrEncodeGIFImg, err.Error())
		}
	}
	return nil
}

// removeSrc:
func removeSrc(srcPath string) (err error) {
	if err = os.Remove(srcPath); err != nil {
		return errors.Wrap(errof.ErrEncodeGIFImg, err.Error())
	}
	return nil
}

func createFile(dstPath string) (dst *os.File, err error) {
	if dst, err = os.Create(dstPath); err != nil {
		return dst, errors.Wrap(errof.ErrCreateDstFile, err.Error())
	}
	return dst, nil
}

// createDir:
func createDir(dirPath string) (err error) {
	// 既にディレクトリが存在しているかを確認する
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return errors.Wrap(errof.ErrCreateDirectory, err.Error())
	}
	return nil
}

// getRootDir:
func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	cvtDir := filepath.Dir(b)
	internalDir := filepath.Dir(cvtDir)
	return filepath.Dir(internalDir)
}

// isFileExist:
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

// isDirExist:
func isDirExist(dir string) (isDir bool, err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(errof.ErrGetDirInfo, err.Error())
	}
	if !fileInfo.IsDir() {
		return false, nil
	}
	return true, nil
}
