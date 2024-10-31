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
)

type CopyInfo struct {
	fromPath    string
	toPath      string
	offset      int64
	limit       int64
	bytesToCopy int64
	fromFile    *os.File
	toFile      *os.File
}

func NewCopyInfo(fromPath, toPath string, offset, limit int64) *CopyInfo {
	return &CopyInfo{
		fromPath: fromPath,
		toPath:   toPath,
		offset:   offset,
		limit:    limit,
	}
}

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

func (ci *CopyInfo) openSrcFile() error {
	if ci.offset < 0 {
		return ErrInvalidOffset
	}
	if ci.limit < 0 {
		return ErrInvalidLimit
	}
	var err error
	ci.fromFile, err = os.Open(ci.fromPath)
	if err != nil {
		return err
	}
	// проверки для входного файла
	fromFileInfo, err := ci.fromFile.Stat()
	if err != nil {
		return err
	}
	if fromFileInfo.IsDir() {
		return ErrUnsupportedFile
	}
	if fromFileInfo.Mode().IsRegular() {
		// offset больше, чем размер файла - невалидная ситуация;
		if ci.offset > fromFileInfo.Size() {
			return ErrOffsetExceedsFileSize
		}
		ci.bytesToCopy = fromFileInfo.Size() - ci.offset
		if ci.limit > 0 {
			ci.bytesToCopy = min(ci.limit, ci.bytesToCopy)
		}
	} else {
		// размер файла неизвестен, тогда нужен лимит
		if ci.limit == 0 {
			return ErrInvalidLimitPositive
		}
		ci.bytesToCopy = ci.limit
	}
	if ci.offset > 0 {
		// проматываем от начала входного файла
		_, err = ci.fromFile.Seek(ci.offset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ci *CopyInfo) doCopy() error {
	// start new bar
	bar := pb.Full.Start64(ci.bytesToCopy)
	defer bar.Finish()
	// create proxy reader
	barReader := bar.NewProxyReader(ci.fromFile)

	for bytesCopied := int64(0); bytesCopied < ci.bytesToCopy; {
		stepToCopy := min(CopyStep, ci.bytesToCopy-bytesCopied)
		written, err := io.CopyN(ci.toFile, barReader, stepToCopy)
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

func (ci *CopyInfo) Copy() error {
	// открываем входной файл
	err := ci.openSrcFile()
	if err != nil {
		if ci.fromFile != nil {
			ci.fromFile.Close()
		}
		return err
	}
	defer ci.fromFile.Close()

	// создаем выходной файл
	useTempFile := sameRealFile(ci.fromPath, ci.toPath)
	if useTempFile {
		// файл источник и файл приемник совпадают - копирование через временный файл
		ci.toFile, err = os.CreateTemp(filepath.Dir(ci.toPath), "copy_temp")
	} else {
		ci.toFile, err = os.Create(ci.toPath)
	}
	if err != nil {
		return err
	}
	defer func() {
		ci.toFile.Close()
		if useTempFile {
			err = os.Rename(ci.toFile.Name(), ci.toPath)
			if err != nil {
				os.Remove(ci.toFile.Name())
			}
		}
	}()
	return ci.doCopy()
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	ci := NewCopyInfo(fromPath, toPath, offset, limit)
	return ci.Copy()
}
