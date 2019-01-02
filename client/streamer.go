package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// ListRecursive list a dir recursively
func ListRecursive(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	fileList := make([]string, 0, len(files))
	for _, f := range files {
		absPath, err := filepath.Abs(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}
		if f.IsDir() {
			nestedFiles, err := ListRecursive(absPath)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, nestedFiles...)
		} else {
			if f.Mode()&os.ModeSymlink == os.ModeSymlink {
				continue
			}
			fileList = append(fileList, absPath)
		}
	}
	return fileList, nil
}

func readFiles(fileChan chan string, concurrency int) chan string {
	outputChan := make(chan string, 100)
	wg := sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for filePath := range fileChan {
				file, err := os.Open(filePath)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					outputChan <- scanner.Text()
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
		go func() {
			wg.Wait()
			close(outputChan)
		}()
	}
	return outputChan
}
