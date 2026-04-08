package utils

import (
	"encoding/base64"
	"os"
)

func Savebase64ToFile(base64str string, fileName string) error {
	dec, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return err
	}
	f, err := os.Create("./uploads/" + fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(dec); err != nil {
		return err
	}
	return nil
}

func DeleteFile(fileName string) {
	// Abaikan jika fileName kosong
	if fileName == "" {
		return
	}
	// Hapus file di path uploads
	_ = os.Remove("./uploads/" + fileName)
}
