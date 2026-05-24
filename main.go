package main

import (
	"fmt"
	"math/big"
)

// main — точка входа в программу.
// Демонстрирует генерацию ключей, подпись сообщения и проверку подписи по ГОСТ Р 34.10-2012.
func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")

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

	// Генерация ключей
	d, H, err := params.GenerateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Закрытый ключ d: %X\n", d)
	fmt.Printf("Открытый ключ H: (%X, %X)\n", H.X, H.Y)

	// Проверка, что открытый ключ лежит на кривой
	if params.IsOnCurve(H) {
		fmt.Println("Открытый ключ корректен (лежит на кривой)")
	} else {
		fmt.Println("Ошибка: открытый ключ не лежит на кривой")
	}

	// ---- Демонстрация подписи и проверки ----
	msg := []byte("Тестовое сообщение для подписи по ГОСТ Р 34.10-2012")
	fmt.Printf("\nСообщение: %s\n", msg)

	// Хеш сообщения
	e := params.hashToNumber(msg)
	fmt.Printf("Хеш (e): %X\n", e)

	// Подпись
	r, s, err := params.Sign(d, e)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Подпись: r = %X\n", r)
	fmt.Printf("        s = %X\n", s)

	// Проверка подписи
	valid := params.Verify(H, msg, r, s)
	fmt.Printf("Проверка подписи: %v\n", valid)

	// Проверка с изменённой подписью (неверное s)
	sBad := new(big.Int).Set(s)
	sBad.Add(sBad, big.NewInt(1))
	validBad := params.Verify(H, msg, r, sBad)
	fmt.Printf("Проверка подписи с изменённым s: %v\n", validBad)
}
