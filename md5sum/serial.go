package main

import (
	"crypto/md5"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"sort"
	"fmt"
)

//没有使用并发，而是简单的通过filepath.Walk()逐个读取根目录下的文件并计算MD5
func MD5AllSerial(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

//计算特定目录下所有文件的md5值
//然后按照路径名顺序打印结果
func main(){
	m, err := MD5AllSerial(os.Args[1])
	if err != nil {
		log.Fatalf("main main() m, err := MD5All(os.Args[1]) error => %v\n", err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}