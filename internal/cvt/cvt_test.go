package cvt

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"testing"
)

func RemoveAll(t *testing.T, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Errorf("failed remove dir: %v", err)
	}
}

func TestPathWalk(t *testing.T) {
	var currentDir string
	var err error
	if currentDir, err = os.Getwd(); err != nil {
		t.Errorf("failed get current dir: %v", err)
	}
	vectors := map[string]struct {
		inputDir  string
		beforeExt string
		srcPaths  []string
		expected  []string
		wantErr   error
	}{
		"OK": {
			inputDir:  "internal/cvt/walktest",
			beforeExt: ".jpg",
			srcPaths: []string{
				fmt.Sprintf("%s/walktest/test001.jpg", currentDir),
				fmt.Sprintf("%s/walktest/test002.jpg", currentDir),
				fmt.Sprintf("%s/walktest/test003.png", currentDir),
			},
			expected: []string{
				fmt.Sprintf("%s/walktest/test001.jpg", currentDir),
				fmt.Sprintf("%s/walktest/test002.jpg", currentDir),
			},
		},
	}
	for k, v := range vectors {
		if err = createDir("walktest"); err != nil {
			t.Errorf("failed create dir: %v", err)
		}
		for _, srcPath := range v.srcPaths {
			if _, err = createFile(srcPath); err != nil {
				t.Errorf("failed create dir: %v", err)
			}
		}
		imageCvtGlue := NewImageCvtGlue(v.inputDir, "", v.beforeExt, "", false)
		actual, err := imageCvtGlue.pathWalk()
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %s, pathWalk() = %v, want %v", k, errors.Cause(err), v.wantErr)
		}

		// 実際の値を取得後ディレクトリ削除
		RemoveAll(t, fmt.Sprintf("%s/walktest", currentDir))
		// 実際に取得した配列長とテストケースの配列長が異なる場合
		if len(v.expected) != len(actual) {
			t.Errorf("the length of arrays are different test: %s length of expected %d length of actual %d", k, len(v.expected), len(actual))
		}
		for idx, expected := range v.expected {
			if expected != actual[idx] {
				t.Errorf("test: %s expected %s actual %s", k, expected, actual[idx])
			}
		}
	}
}
