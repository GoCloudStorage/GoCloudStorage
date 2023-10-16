package local

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
)

var (
	Client *chunkStorage
)

type chunkStorage struct {
	root string
}

func Init(root string) {
	Client = &chunkStorage{
		root: root,
	}
}

func (c *chunkStorage) SaveChunk(key string, chunkNumber int, data io.Reader, offset int64) error {

	fileDir := path.Join(c.root, key)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir, err: %v", err)
	}
	filepath := path.Join(fileDir, strconv.Itoa(chunkNumber))
	return saveFile(filepath, data, offset)
}

func (c *chunkStorage) MergeChunk(key string, totalSize int) (filePath string, err error) {
	fileDir := path.Join(c.root, key)
	dirs, err := os.ReadDir(fileDir)

	if err != nil {
		if err2 := removeDir(fileDir); err2 != nil {
			return "", fmt.Errorf("%v, %v", err, err2)
		}
		return "", fmt.Errorf("failed to read dir [%s], err: %v", fileDir, err)
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	// create final local file
	finalFilePath := path.Join(fileDir, "data")
	for i, part := range dirs {
		srcPath := path.Join(fileDir, part.Name())
		size, err := copyFileToFile(srcPath, finalFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to copy [%d] chunk file to dst file, err: %v", i, err)
		}
		totalSize -= size
	}

	if totalSize != 0 {
		if err2 := removeDir(fileDir); err2 != nil {
			return "", fmt.Errorf("%v, %v", err, err2)
		}
		return "", fmt.Errorf("merge chunk not complete, %d", totalSize)
	}

	// remove chunk file
	for i, part := range dirs {
		if part.Name() != "data" {
			filePath := path.Join(fileDir, part.Name())
			if err := os.Remove(filePath); err != nil {
				if err2 := removeDir(fileDir); err2 != nil {
					return "", fmt.Errorf("%v, %v", err, err2)
				}
				return "", fmt.Errorf("failed to remove [%d] chunk file, err: %v", i, err)
			}
		}
	}
	return finalFilePath, nil
}

func (c *chunkStorage) Remove(key string) error {
	fileDir := path.Join(c.root, key)
	return removeDir(fileDir)
}
