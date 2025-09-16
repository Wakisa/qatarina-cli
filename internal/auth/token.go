package auth

import (
	"os"
	"path/filepath"
)

const tokenFileName = "token.json"

func tokenFilePath() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, ".qatarina", tokenFileName)
}

func SaveToken(token string) error {
	os.MkdirAll(filepath.Dir(tokenFilePath()), 0700)
	return os.WriteFile(tokenFilePath(), []byte(token), 0600)
}
func LoadToken() (string, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DeleteToken() error {
	return os.Remove(tokenFilePath())
}
