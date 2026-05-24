package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

// saveSignature сохраняет подпись (r, s) в файл.
// Формат: r\ns (две строки, каждая в hex-формате).
func saveSignature(r, s *big.Int, filename string) error {
	content := fmt.Sprintf("%X\n%X", r, s)
	return os.WriteFile(filename, []byte(content), 0644)
}

// loadSignature загружает подпись (r, s) из файла.
// Ожидается формат: r\ns (две строки, каждая в hex-формате).
func loadSignature(filename string) (*big.Int, *big.Int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return nil, nil, fmt.Errorf("некорректный формат файла подписи (ожидаются r и s)")
	}
	r := new(big.Int)
	r.SetString(strings.TrimSpace(lines[0]), 16)
	s := new(big.Int)
	s.SetString(strings.TrimSpace(lines[1]), 16)
	return r, s, nil
}
