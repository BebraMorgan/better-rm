package actions

import (
	"brm/localization"
	"brm/trash"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrRemoveRoot      = errors.New(localization.GetMessage("err_remove_root"))
	ErrRemoveTrashSelf = errors.New(localization.GetMessage("err_remove_trash_self"))
)

func GetTrashPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	trashPath := filepath.Join(home, ".trash")

	info, err := os.Stat(trashPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(trashPath, 0750)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else if !info.IsDir() {
		return "", os.ErrInvalid
	}

	return trashPath, nil
}

func getUniquePath(dir, baseName string) (string, error) {
	ext := filepath.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)

	uniqueName := baseName
	counter := 1

	for {
		fullPath := filepath.Join(dir, uniqueName)
		_, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			return fullPath, nil
		} else if err != nil {
			return "", err
		}
		uniqueName = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}
}

func MoveDir(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return MoveFile(path, targetPath)
	})
	if err != nil {
		return err
	}

	return os.RemoveAll(src)
}

func ClearDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}

		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func MoveFile(srcPath, dstPath string) error {
	inputFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(srcPath)
	if err == nil {
		_ = os.Chmod(dstPath, info.Mode())
	}

	return os.Remove(srcPath)
}

func SaveDelete(srcPath string) error {
	absSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	trashPath, err := GetTrashPath()
	if err != nil {
		return err
	}

	absTrashPath, err := filepath.Abs(trashPath)
	if err != nil {
		return err
	}

	if absSrcPath == "/" {
		return ErrRemoveRoot
	}

	if absSrcPath == absTrashPath {
		return os.RemoveAll(absTrashPath)
	}

	fileName := filepath.Base(absSrcPath)

	info, err := os.Stat(absSrcPath)
	if err != nil {
		return err
	}

	dstPath, err := getUniquePath(absTrashPath, fileName)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = MoveDir(absSrcPath, dstPath)
	} else {
		err = MoveFile(absSrcPath, dstPath)
	}
	if err != nil {
		return err
	}

	err = trash.AddTrashInfoEntry(trash.TrashInfo{
		TrashName:    filepath.Base(dstPath),
		OriginalPath: absSrcPath,
		DeletionDate: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func EmptyTrash() error {
	trashPath, err := GetTrashPath()
	if err != nil {
		return err
	}
	return ClearDir(trashPath)
}
func Restore() error {
	filePath, err := trash.GetTrashInfoPath()
	if err != nil {
		return err
	}

	entries, err := trash.LoadTrashInfo(filePath)
	if err != nil {
		return err
	}

	trashFilesDir, err := getTrashFilesDir()
	if err != nil {
		return err
	}

	var remainingEntries []trash.TrashInfo

	for _, entry := range entries {
		trashFilePath := filepath.Join(trashFilesDir, entry.TrashName)

		if _, err := os.Stat(trashFilePath); os.IsNotExist(err) {
			continue
		}

		restoreDir := filepath.Dir(entry.OriginalPath)
		if err := os.MkdirAll(restoreDir, 0755); err != nil {
			remainingEntries = append(remainingEntries, entry)
			continue
		}

		err = os.Rename(trashFilePath, entry.OriginalPath)
		if err != nil {
			remainingEntries = append(remainingEntries, entry)
			continue
		}
	}

	err = trash.SaveTrashInfo(filePath, remainingEntries)
	if err != nil {
		return err
	}

	return nil
}

func getTrashFilesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot determine home directory")
	}
	trashPath := filepath.Join(home, ".trash")
	return trashPath, nil
}
