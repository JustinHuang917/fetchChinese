package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regExp = "([\u4e00-\u9fa5]+)+?"
	filter = flag.String("filter", "*.*", " filter")
	r      = regexp.MustCompile(regExp)
	action = flag.String("action", "fetch", "fetch")
	dir    = flag.String("dir", "", "dir to fecth/reverse")
)

type Config struct {
	Comments map[string]CommentsConfig `json:Comments`
}

var AppConfig *Config

func init() {
	err := load("./config.json")
	if err != nil {
		fmt.Println("Init Config Error:", err)
	}
}
func load(cfgPath string) error {
	file, err := os.Open(string(cfgPath))
	if err != nil {
		return err
	}
	AppConfig = &Config{}
	dec := json.NewDecoder(file)
	if err = dec.Decode(AppConfig); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if *action == "fetch" {
		if *dir != "" {
			rs := fecthDir(*dir)
			outputResult(rs)
		}
	} else if *action == "reverse" {
		panic("Not Implemented")
	}
}
func outputResult(rs []fectchResult) {
	if rs != nil {
		for index, result := range rs {
			if result.err == nil {
				fmt.Printf("%d, filePath:%s,row:%d,col:%d length:%d,word:%s\n", index, result.filePath, result.rowNum, result.colNum, result.length, result.word)
			}
		}
	}
}

type CommentsConfig struct {
	SingleLine_Begin string
	Multiline_Begin  string
	Multiline_End    string
}

type fectchResult struct {
	filePath string
	fileName string
	rowNum   int
	colNum   int
	word     string
	length   int
	err      error
}
type line struct {
	rawContent        string
	beginCommentIndex int
}

func constructLines(rows []string, fileExtName string) []*line {
	inMultiLinesScope := false
	lines := make([]*line, 0, 10)
	for _, row := range rows {
		l := &line{
			rawContent:        row,
			beginCommentIndex: -1,
		}
		if comment, ok := AppConfig.Comments[fileExtName]; ok {
			if inMultiLinesScope {
				l.beginCommentIndex = 0
			} else {
				if comment.Multiline_Begin != "" {
					singleLineCommentIndex := strings.Index(row, comment.SingleLine_Begin)
					if singleLineCommentIndex >= 0 {
						l.beginCommentIndex = singleLineCommentIndex
					}
				}
			}
			if comment.Multiline_Begin != "" {
				multilineBeginIndex := strings.Index(row, comment.Multiline_Begin)
				if multilineBeginIndex >= 0 {
					l.beginCommentIndex = multilineBeginIndex
					inMultiLinesScope = true
				}
				multilineEndIndex := strings.Index(row, comment.Multiline_End)
				if multilineEndIndex >= 0 {
					inMultiLinesScope = false
				}
			}
		}
		lines = append(lines, l)
	}
	return lines

}
func fecthOneFile(path string, f os.FileInfo) []fectchResult {
	if f == nil {
		return nil
	}
	if f.IsDir() {
		return nil
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		return nil
	} else if ok, _ := filepath.Match(*filter, f.Name()); ok {
		rs := make([]fectchResult, 0, 10)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			rs = append(rs, *&fectchResult{err: err})
			return rs
		}
		rows := strings.Split(string(content), "\n")
		lines := constructLines(rows, filepath.Ext(path))
		for index, l := range lines {
			result := r.FindAllStringIndex(l.rawContent, -1)
			if len(result) > 0 {
				for _, cell := range result {
					if len(cell) == 2 {
						if (l.beginCommentIndex == -1) || (cell[0] < l.beginCommentIndex) {
							w := l.rawContent[cell[0]:cell[1]]
							fc := *&fectchResult{
								filePath: path,
								fileName: f.Name(),
								rowNum:   index,
								colNum:   cell[0],
								word:     w,
								length:   len(w),
							}
							rs = append(rs, fc)
						}
					}
				}
			}
		}
		return rs
	} else {
		return nil
	}
}

func fecthDir(dirPath string) []fectchResult {
	frs := make([]fectchResult, 0, 10)
	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		rs := fecthOneFile(path, f)
		if rs != nil {
			frs = append(frs, rs...)
		}
		return nil
	})
	return frs
}
