package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func (curve *CurveParams) GenerateKey() (d *big.Int, H *Point, err error) {
	Q := curve.Q
	if Q == nil || Q.Sign() <= 0 {
		return nil, nil, fmt.Errorf("некорректный порядок Q")
	}

	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	for {
		d, err = rand.Int(rand.Reader, Q)
		if err != nil {
			return nil, nil, err
		}
		if d.Sign() != 0 {
			break
		}
	}
	H = curve.ScalarMult(d, G)
	return d, H, nil
}

func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")

	// Параметры кривой (тестовый набор)
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
}
