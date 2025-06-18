package trash

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type TrashInfo struct {
	TrashName    string    `json:"trash_name"`
	OriginalPath string    `json:"original_path"`
	DeletionDate time.Time `json:"deletion_date"`
}

func GetTrashInfoPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot determine user home directory")
	}
	dirPath := filepath.Join(home, ".brm")
	filePath := filepath.Join(dirPath, "trash.json")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.WriteFile(filePath, []byte("[]"), 0644)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

func SaveTrashInfo(path string, entries []TrashInfo) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}

func LoadTrashInfo(path string) ([]TrashInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []TrashInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&entries)
	return entries, err
}

func AddTrashInfoEntry(entry TrashInfo) error {
	var entries []TrashInfo

	filePath, err := GetTrashInfoPath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&entries); err != nil {
			entries = []TrashInfo{}
		}
	}

	entries = append(entries, entry)

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(entries); err != nil {
		return err
	}

	return nil
}
