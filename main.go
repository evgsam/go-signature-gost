package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// main — точка входа в программу.
// Демонстрирует генерацию ключей, подпись сообщения и проверку подписи по ГОСТ Р 34.10-2012.
func main() {
	// Инициализируем сканер для чтения ввода пользователя
	scanner := bufio.NewScanner(os.Stdin)

	// Параметры кривой (набор 1.2.643.2.2.35.1 — ГОСТ Р 34.10-2015, 256 бит)
	params := &CurveParams{
		OID: "1.2.643.2.2.35.1",
		P:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD97"),
		A:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD94"),
		B:   mustParse("A6"),
		M:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		Q:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		GX:  mustParse("1"),
		GY:  mustParse("8D91E471E0989CDA27DF505A453F2B7635294F2DDF23E3B122ACC99C9E9F1E14"),
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
