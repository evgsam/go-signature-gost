package main

import (
	"fmt"
	"math/big"
	"os"
)

// SignFile подписывает файл и сохраняет подпись в файл.
// Алгоритм:
//  1. Читает содержимое файла
//  2. Вычисляет хеш сообщения (hashToNumber)
//  3. Формирует подпись (r, s) по закрытому ключу d
//  4. Сохраняет подпись в signatureFile
func SignFile(params *CurveParams, msgFile, signatureFile string, d *big.Int) error {
	// Читаем содержимое файла
	data, err := os.ReadFile(msgFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Вычисляем хеш сообщения
	e := params.hashToNumber(data)

	// Формируем подпись
	r, s, err := params.Sign(d, e)
	if err != nil {
		return fmt.Errorf("ошибка создания подписи: %w", err)
	}

	// Сохраняем подпись
	if err := saveSignature(r, s, signatureFile); err != nil {
		return fmt.Errorf("ошибка сохранения подписи: %w", err)
	}

	fmt.Println("Подпись успешно создана: ")
	fmt.Printf(" r = %X\n", r)
	fmt.Printf(" s = %X\n", s)
	return nil
}

// VerifyFile проверяет подпись файла.
// Алгоритм:
//  1. Читает содержимое файла
//  2. Загружает подпись (r, s) из файла
//  3. Загружает открытый ключ H
//  4. Проверяет подпись (Verify)
//  5. Выводит результат
func VerifyFile(params *CurveParams, msgFile, signatureFile, publicKeyFile string) error {
	// Читаем содержимое файла
	data, err := os.ReadFile(msgFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Загружаем подпись
	r, s, err := loadSignature(signatureFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки подписи: %w", err)
	}

	// Загружаем открытый ключ
	H, err := loadPublicKey(publicKeyFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки открытого ключа: %w", err)
	}

	// Проверяем подпись
	valid := params.Verify(H, data, r, s)
	if valid {
		fmt.Println("Подпись корректна")
		return nil
	}
	fmt.Println("Подпись НЕ корректна")
	return fmt.Errorf("подпись недействительна")
}
