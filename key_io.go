package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

// savePrivateKey сохраняет закрытый ключ d как hex-строку в файл.
func savePrivateKey(d *big.Int, filename string) error {
	hexStr := d.Text(16)
	return os.WriteFile(filename, []byte(hexStr), 0600)
}

// loadPrivateKey загружает закрытый ключ d из файла (hex-строка).
func loadPrivateKey(filename string) (*big.Int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	d := new(big.Int)
	d.SetString(strings.TrimSpace(string(data)), 16)
	return d, nil
}

// savePublicKey сохраняет открытую ключевую точку H в файл.
// Формат: X\nY (две строки, каждая в hex-формате).
func savePublicKey(H *Point, filename string) error {
	content := fmt.Sprintf("%X\n%X", H.X, H.Y)
	return os.WriteFile(filename, []byte(content), 0644)
}

// loadPublicKey загружает открытую ключевую точку H из файла.
// Ожидается формат: X\nY (две строки, каждая в hex-формате).
func loadPublicKey(filename string) (*Point, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("некорректный формат файла открытого ключа (ожидаются X и Y)")
	}
	x := new(big.Int)
	x.SetString(strings.TrimSpace(lines[0]), 16)
	y := new(big.Int)
	y.SetString(strings.TrimSpace(lines[1]), 16)
	return &Point{X: x, Y: y, Inf: false}, nil
}
