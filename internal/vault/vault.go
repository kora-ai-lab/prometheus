package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Vault struct {
	path string
}

func New(path string) *Vault {
	return &Vault{path: path}
}

func (v *Vault) Set(key, value string) error {
	data, err := v.load()
	if err != nil {
		return err
	}
	data[key] = value
	return v.save(data)
}

func (v *Vault) Get(key string) (string, error) {
	data, err := v.load()
	if err != nil {
		return "", err
	}
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("vault key %q not found", key)
	}
	return val, nil
}

func (v *Vault) List() ([]string, error) {
	data, err := v.load()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys, nil
}

func (v *Vault) load() (map[string]string, error) {
	if _, err := os.Stat(v.path); os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	raw, err := os.ReadFile(v.path)
	if err != nil {
		return nil, err
	}
	plain, err := decrypt(raw)
	if err != nil {
		return nil, err
	}
	var data map[string]string
	if err := json.Unmarshal(plain, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (v *Vault) save(data map[string]string) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	enc, err := encrypt(raw)
	if err != nil {
		return err
	}
	return os.WriteFile(v.path, enc, 0o600)
}

func encrypt(plain []byte) ([]byte, error) {
	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return append(nonce, gcm.Seal(nil, nonce, plain, nil)...), nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	body := ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, body, nil)
}

func deriveKey() []byte {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	hostname, _ := os.Hostname()
	sum := sha256.Sum256([]byte(strings.TrimSpace(hostname + ":" + user + ":prometheus-v1")))
	return sum[:]
}
