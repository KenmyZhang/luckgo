package model

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
	"strings"
)

var versions = []string{
	"1.0",
}

var CurrentVersion string = versions[0]
var BuildDate string
var BuildHash string

const (
	LOG_FILENAME = "critic.log"
)

var (
	commonBaseSearchPaths = []string{
		".",
		"..",
		"../..",
		"../../..",
	}
)

const (
	InvalidParam        = 4000
	InvalidFrontImg     = 4001
	InvalidBackImg      = 4002
	InvalidEmail        = 4003
	InvalidMobile       = 4004
	InvalidWechat       = 4005
	InvalidQQ           = 4006
	InvalidDescription  = 4007
	InvalidType         = 4008
	AlreadyApproved     = 4009
	AlreadySubmit       = 4010
	InvalidStatus       = 4011
	AlreadyHandle       = 4012
	RecordNotFound      = 4013
	InvalidScrshot      = 4014
	InvalidLink         = 4015
	InternalServerError = 5000
)

type Err struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func (e *Err) Error() string {
	if e == nil {
		return ""
	}
	return "code:" + strconv.Itoa(e.Code) + ", result:" + e.Result
}

func NewInvalidParamError(code int, where, parameter, details string) *Err {
	var message string
	if details != "" {
		message = ", details:" + details + ", where:" + where
	} else {
		message = ", where:" + where
	}
	return &Err{Code: code, Result: "Invalid " + parameter + " patameter" + message}
}

func NewInternalServerError(where, details string) *Err {
	var message string
	if details != "" {
		message = "details:" + details + ", where:" + where
	} else {
		message = "where:" + where
	}
	return &Err{Code: InternalServerError, Result: "Internal Server Error," + message}
}

// GetMillis is a convience method to get milliseconds since epoch.
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func FindConfigFile(fileName string) (path string) {
	found := FindFile(filepath.Join("config", fileName))
	if found == "" {
		found = FindPath(fileName, []string{"."}, nil)
	}

	return found
}

func FindPath(path string, baseSearchPaths []string, filter func(os.FileInfo) bool) string {
	//判斷是否是絕對路徑
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}

		return ""
	}

	searchPaths := []string{}
	for _, baseSearchPath := range baseSearchPaths {
		searchPaths = append(searchPaths, baseSearchPath)
	}

	var binaryDir string
	//返回启动当前进程的可执行文件的路径名称。
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			if exe, err = filepath.Abs(exe); err == nil {
				binaryDir = filepath.Dir(exe)
			}
		}
	}
	if binaryDir != "" {
		for _, baseSearchPath := range baseSearchPaths {
			searchPaths = append(
				searchPaths,
				filepath.Join(binaryDir, baseSearchPath),
			)
		}
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, path))
		if err != nil {
			continue
		} else if fileInfo, err := os.Stat(found); err == nil {
			if filter != nil {
				if filter(fileInfo) {
					return found
				}
			} else {
				return found
			}
		}
	}

	return ""
}

func GetLogFileLocation(fileLocation string) string {
	if fileLocation == "" {
		fileLocation, _ = FindDir("logs")
	}

	return filepath.Join(fileLocation, LOG_FILENAME)
}

func FindDir(dir string) (string, bool) {
	found := FindPath(dir, commonBaseSearchPaths, func(fileInfo os.FileInfo) bool {
		return fileInfo.IsDir()
	})
	if found == "" {
		return "./", false
	}

	return found, true
}

func FindFile(path string) string {
	return FindPath(path, commonBaseSearchPaths, func(fileInfo os.FileInfo) bool {
		return !fileInfo.IsDir()
	})
}

func checkFileExit(filename string) bool {
	var exit = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exit = false
	}
	return exit
}

func GetFile(filename string) (*os.File, error) {
	if checkFileExit(filename) {
		return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		index := strings.LastIndex(filename, "/")
		// 创建文件夹
		err := os.MkdirAll(filename[:index], os.ModePerm)
		if err != nil {
			return nil, err
		}
		return os.Create(filename)
	}

}