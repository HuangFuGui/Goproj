package main

import (
	"os"
	"sort"
	"fmt"
	"log"
	"crypto/md5"
	"sync"
	"path/filepath"
	"io/ioutil"
	"errors"
)

type result struct {
	path 	string
	sum 	[md5.Size]byte
	err 	error
}

func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	ch := make(chan result)
	errch := make(chan error, 1)//记录walk过程中的错误：1、例如os.Lstat，2、被done取消、3、...
	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			//如果walkFn返回的err!=nil，walk函数就会返回，不再递归下降遍历
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case ch <- result{path, md5.Sum(data), err}:
				case <- done:
				}
				wg.Done()
			}()
			select {
			case <- done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		//此时walk函数已经返回，所有文件的wg.Add调用都已经结束
		//启动一个goroutine，待所有结果都计算完毕后，关闭ch
		go func() {
			wg.Wait()
			close(ch)
		}()
		errch <- err
	}()
	return ch, errch
}

//并发遍历文件夹并计算MD5值，注意程序设计思路
func MD5AllParallel(root string) (map[string][md5.Size]byte, error) {
	//done channel用来标记结束遍历文件夹
	done := make(chan struct{})
	defer close(done)//除此之外可能还有别的原因使done channel关闭

	ch, errch := sumFiles(done, root)

	m := make(map[string][md5.Size]byte)
	for r := range ch {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <- errch; err != nil {
		return nil, err
	}
	return m, nil
}

func main(){
	m, err := MD5AllParallel(os.Args[1])
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