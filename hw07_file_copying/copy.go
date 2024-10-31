package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const CopyStep = 8192

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidOffset         = errors.New("offset cannot be negative")
	ErrInvalidLimit          = errors.New("limit cannot be negative")
	ErrInvalidLimitPositive  = errors.New("limit must be positive")
	ErrEmptyFileName         = errors.New("empty file name")
)

func sameRealFile(srcFileName, destFileName string) bool {
	destFI, err := os.Stat(destFileName)
	// файл результата существует ?
	if err == nil {
		srcFI, _ := os.Stat(srcFileName)
		// файл результата и исходный файл - один и тот же обычный файл ?
		return os.SameFile(srcFI, destFI) && destFI.Mode().IsRegular()
	}
	return false
}

func checkArgs(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}
	if fromPath == "" || toPath == "" {
		return ErrEmptyFileName
	}
	return nil
}

func prepareSrcFile(fromPath string, offset, limit int64, bytesToCopy *int64) (*os.File, error) {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return nil, err
	}
	// проверки для входного файла
	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return nil, err
	}
	if fromFileInfo.IsDir() {
		return nil, ErrUnsupportedFile
	}
	if fromFileInfo.Mode().IsRegular() {
		// offset больше, чем размер файла - невалидная ситуация;
		if offset > fromFileInfo.Size() {
			return nil, ErrOffsetExceedsFileSize
		}
		*bytesToCopy = fromFileInfo.Size() - offset
		if limit > 0 {
			*bytesToCopy = min(limit, *bytesToCopy)
		}
	} else {
		// размер файла неизвестен, тогда нужен лимит
		if limit == 0 {
			return nil, ErrInvalidLimitPositive
		}
		*bytesToCopy = limit
	}
	if offset > 0 {
		// проматываем от начала входного файла
		_, err = fromFile.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}
	return fromFile, nil
}

func doCopy(fromFile *os.File, toFile *os.File, bytesToCopy int64) error {
	// start new bar
	bar := pb.Full.Start64(bytesToCopy)
	defer bar.Finish()
	// create proxy reader
	barReader := bar.NewProxyReader(fromFile)

	for bytesCopied := int64(0); bytesCopied < bytesToCopy; {
		stepToCopy := min(CopyStep, bytesToCopy-bytesCopied)
		written, err := io.CopyN(toFile, barReader, stepToCopy)
		bytesCopied += written
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// проверки аргументов
	if err := checkArgs(fromPath, toPath, offset, limit); err != nil {
		return err
	}
	// подготовка источника данных
	var bytesToCopy int64
	fromFile, err := prepareSrcFile(fromPath, offset, limit, &bytesToCopy)
	if err != nil {
		if fromFile != nil {
			fromFile.Close()
		}
		return err
	}
	defer fromFile.Close()

	// создаем выходной файл
	var toFile *os.File
	useTempFile := sameRealFile(fromPath, toPath)
	if useTempFile {
		// файл источник и файл приемник совпадают - копирование через временный файл
		toFile, err = os.CreateTemp(filepath.Dir(toPath), "copy_temp")
	} else {
		toFile, err = os.Create(toPath)
	}
	if err != nil {
		return err
	}
	defer func() {
		toFile.Close()
		if useTempFile {
			err = os.Rename(toFile.Name(), toPath)
			if err != nil {
				os.Remove(toFile.Name())
			}
		}
	}()
	return doCopy(fromFile, toFile, bytesToCopy)
}
