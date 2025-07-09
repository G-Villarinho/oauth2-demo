package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

func NewConfig() (*Environment, error) {
	var environment Environment

	if err := LoadEnviroment(&environment); err != nil {
		return nil, err
	}

	return &environment, nil
}

func LoadEnviroment(environment *Environment) error {
	if environment == nil {
		return errors.New("environment is nil")
	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	if _, err := env.UnmarshalFromEnviron(environment); err != nil {
		return fmt.Errorf("load environment variables: %w", err)
	}

	if environment.Key.PrivateKey == "" || environment.Key.PublicKey == "" {
		privateKey, err := LoadKeyFromFile("../ecdsa_private.pem")
		if err != nil {
			return fmt.Errorf("load private key: %w", err)
		}

		publicKey, err := LoadKeyFromFile("../ecdsa_public.pem")
		if err != nil {
			return fmt.Errorf("load public key: %w", err)
		}

		environment.Key.PrivateKey = privateKey
		environment.Key.PublicKey = publicKey
	}

	return nil
}

func LoadKeyFromFile(filename string) (string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(currentFile)
	fullPath := filepath.Join(baseDir, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	return strings.TrimSpace(string(data)), nil
}
