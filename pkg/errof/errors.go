package errof

type UserErr string

func (e UserErr) Error() (msg string) {
	var ok bool
	if msg, ok = ErrCodeNames[e]; !ok {
		return string(e)
	}
	return msg
}

var ErrCodeNames = map[UserErr]string{
	ErrWalkingPath: "指定されたディレクトリのトラバースに失敗しました",
	ErrGetFileInfo: "指定されたファイル情報の取得に失敗しました",
	ErrOpenFile:    "指定されたファイルの展開に失敗しました",
}

var (
	ErrWalkingPath UserErr = "ErrWalkingPath"
	ErrGetFileInfo UserErr = "ErrGetFileInfo"
	ErrOpenFile    UserErr = "ErrOpenFile"
)
