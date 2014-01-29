package main

import (
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

type fectchResult struct {
	filePath string
	fileName string
	rowNum   int
	colNum   int
	word     string
	length   int
	err      error
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
		for index, row := range rows {
			result := r.FindAllStringIndex(row, -1)
			if len(result) > 0 {
				for _, cell := range result {
					if len(cell) == 2 {
						w := row[cell[0]:cell[1]]
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
