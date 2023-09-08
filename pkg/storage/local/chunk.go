package local

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
)

type chunkUploader struct {
}

func (receiver chunkUploader) saveChunk(fileDir string, partNum int, data io.Reader) error {
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir, err: %v", err)
	}
	filepath := path.Join(fileDir, strconv.Itoa(partNum))
	return saveFile(filepath, data)
}

func (receiver chunkUploader) mergeChunk(fileDir string, partSize int, totalSize int) error {
	dirs, err := os.ReadDir(fileDir)
	if err != nil {
		return fmt.Errorf("failed to read dir [%s], err: %v", fileDir, err)
	}
	if len(dirs) != partSize {
		return fmt.Errorf("file chunk not complete, need %d have %d", partSize, len(dirs))
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	// create final storage file
	finalFilePath := path.Join(fileDir, "data")
	for i, part := range dirs {
		srcPath := path.Join(fileDir, part.Name())
		size, err := copyFileToFile(srcPath, finalFilePath)
		if err != nil {
			return fmt.Errorf("failed to copy [%d] chunk file to dst file, err: %v", i, err)
		}
		totalSize -= size
	}

	if totalSize != 0 {
		os.Remove(finalFilePath)
		return fmt.Errorf("merge chunk not complete, %d", totalSize)
	}

	// remove chunk file
	for i, part := range dirs {
		if part.Name() != "data" {
			filePath := path.Join(fileDir, part.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove [%d] chunk file, err: %v", i, err)
			}
		}
	}
	return nil
}
