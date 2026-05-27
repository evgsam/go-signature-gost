package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// main — точка входа в программу.
// Демонстрирует генерацию ключей, подпись сообщения и проверку подписи по ГОСТ Р 34.10-2012 (512 бит).
func main() {
	// Инициализируем сканер для чтения ввода пользователя
	scanner := bufio.NewScanner(os.Stdin)

	// Параметры кривой id-GostR3410-2012-512A (OID 1.2.643.7.1.2.1.2.1)
	params := &CurveParams{
		OID: "1.2.643.7.1.2.1.2.1",
		P:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFDC7"),
		A:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFDC4"),
		B:   mustParse("E8C2505DEDFC86DDC1BD0B2B6667F1DA34B82574761CB0E879BD081CFD0B6265EE3CB090F30D27614CB4574010DA90DD862EF9D4EBEE4761503190785A71C760"),
		Q:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF27E69532F48D89116FF22B8D4E0560609B4B38ABFAD2B85DCACDB1411F10B275"),
		GX:  mustParse("3"),
		GY:  mustParse("7503CFE87A836AE3A61B8816E25450E6CE5E1C93ACF1ABC1778064FDCBEFA921DF1626BE4FD036E93D75E6A50E3A41E98028FE5FC235F5B889A589CB5215F2A4"),
	}

	// Основной цикл программы — отображение меню до выбора выхода
	for {
		fmt.Println("1. Сгенерировать ключевую пару")
		fmt.Println("2. Подписать файл")
		fmt.Println("3. Проверить подпись")
		fmt.Println("4. Выход")
		fmt.Print(">> ")

		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		// Генерация ключевой пары
		case "1":
			fmt.Println("Генерация ключей...")
			d, H, err := params.GenerateKey()
			if err != nil {
				fmt.Printf("Ошибка генерации: %v\n", err)
				continue
			}

			// Сохраняем закрытый ключ
			if err := savePrivateKey(d, "private.key"); err != nil {
				fmt.Printf("Ошибка сохранения закрытого ключа: %v\n", err)
				continue
			}

			// Сохраняем открытый ключ
			if err := savePublicKey(H, "public.key"); err != nil {
				fmt.Printf("Ошибка сохранения открытого ключа: %v\n", err)
				continue
			}

			// Проверка, что открытый ключ лежит на кривой
			if params.IsOnCurve(H) {
				fmt.Println("Ключевая пара создана")
				fmt.Println("Открытый ключ корректен (лежит на кривой)")
				fmt.Println("Закрытый ключ: private.key")
				fmt.Println("Открытый ключ: public.key")
			} else {
				fmt.Println("Ошибка: открытый ключ не лежит на кривой")
			}

		// Подпись файла
		case "2":
			fmt.Print("Введите путь к файлу для подписи: ")
			if !scanner.Scan() {
				return
			}
			msgFile := strings.TrimSpace(scanner.Text())

			fmt.Print("Введите путь к закрытому ключу (private.key): ")
			if !scanner.Scan() {
				return
			}
			privKeyFile := strings.TrimSpace(scanner.Text())

			fmt.Print("Введите путь для сохранения подписи (signature.sig): ")
			if !scanner.Scan() {
				return
			}
			sigFile := strings.TrimSpace(scanner.Text())

			// Загружаем закрытый ключ
			d, err := loadPrivateKey(privKeyFile)
			if err != nil {
				fmt.Printf("Ошибка загрузки закрытого ключа: %v\n", err)
				continue
			}

			// Подписываем файл
			if err := SignFile(params, msgFile, sigFile, d); err != nil {
				fmt.Printf("Ошибка подписи: %v\n", err)
				continue
			}

		// Проверка подписи
		case "3":
			fmt.Print("Введите путь к подписываемому файлу: ")
			if !scanner.Scan() {
				return
			}
			msgFile := strings.TrimSpace(scanner.Text())

			fmt.Print("Введите путь к файлу подписи (signature.sig): ")
			if !scanner.Scan() {
				return
			}
			sigFile := strings.TrimSpace(scanner.Text())

			fmt.Print("Введите путь к открытому ключу (public.key): ")
			if !scanner.Scan() {
				return
			}
			pubKeyFile := strings.TrimSpace(scanner.Text())

			// Проверяем подпись
			if err := VerifyFile(params, msgFile, sigFile, pubKeyFile); err != nil {
				fmt.Printf("Ошибка проверки подписи: %v\n", err)
				continue
			}

		// Выход из программы
		case "4":
			fmt.Println("Выход.")
			return

		default:
			fmt.Println("Неверный пункт меню")
		}
	}
}
