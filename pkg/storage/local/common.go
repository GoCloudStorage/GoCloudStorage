package local

import (
	"fmt"
	"io"
	"os"
)

func copyFileToFile(srcPath, dstPath string) (n int, err error) {
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0755)
	if err != nil {
		return 0, fmt.Errorf("failed to open src file: %s, err: %v", srcPath, err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return 0, fmt.Errorf("failed to open src file: %s, err: %v", srcPath, err)
	}
	defer dstFile.Close()

	tmpData := make([]byte, 1024)

	var total int

	for {
		n, err := srcFile.Read(tmpData)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read file data, file: %s, err: %v", srcPath, err)
		}

		writeN, err := dstFile.Write(tmpData[:n])
		if err != nil {
			return 0, fmt.Errorf("failed to write final file data, file: %s, chunk file: %s, err: %v", dstPath, srcPath, err)
		}

		if writeN != n {
			return 0, fmt.Errorf("write file not complete")
		}

		total += n
	}
	return total, nil
}

func saveFile(filepath string, data io.Reader) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("failed to open file, [file path] : %s, err: %v", filepath, err)
	}

	for {
		tmpData := make([]byte, 1024)
		n, err := data.Read(tmpData)
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed write file to %s, err: %v", filepath, err)
		}

		writeN, err := file.Write(tmpData[:n])
		if err != nil {
			return fmt.Errorf("failed write file to %s, err: %v", filepath, err)
		}
		if writeN != n {
			return fmt.Errorf("write file not complete")
		}
	}

	return nil
}
