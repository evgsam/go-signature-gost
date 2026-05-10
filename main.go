package main

import (
	"fmt"
	"math/big"
)

type CurveParams struct {
	OID string   // Идентификатор набора параметров
	P   *big.Int // модуль поля
	A   *big.Int // коэффициент a
	B   *big.Int // коэффициент b
	M   *big.Int // порядок группы точек кривой
	Q   *big.Int // порядок подгруппы, в которой строится ЭЦП
	GX  *big.Int // x‑координата базовой точки P
	GY  *big.Int // y‑координата базовой точки P
}

func mustParse(hexStr string) *big.Int {
	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		panic("неверная шестнадцатеричная строка: " + hexStr)
	}
	return n
}

func IsValidCurveParams(params *CurveParams) bool {
	x, y, p, a, b := params.GX, params.GY, params.P, params.A, params.B
	//lhs = y^2 mod p
	lhs := new(big.Int).Exp(y, big.NewInt(2), p)
	// rhs = (x^3 + a*x + b) mod p
	rhs := new(big.Int).Mul(a, x)
	rhs.Add(rhs, b)
	rhs.Add(rhs, new(big.Int).Exp(x, big.NewInt(3), nil))
	rhs.Mod(rhs, p)
	return lhs.Cmp(rhs) == 0
}

func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")
	params := &CurveParams{
		OID: "1.2.643.2.2.35.1",
		P:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD97"),
		A:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD94"),
		B:   mustParse("A6"),
		M:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		Q:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		GX:  mustParse("1"),
		GY:  mustParse("8D91E471E0989CDA27DF505A453F2B7635294F2DDF23E3B122ACC99C9E9F1E14"), //в big-endian виде
	}
	if !IsValidCurveParams(params) {
		panic("параметры кривой не проходят проверку")
	}
	fmt.Println("Параметры кривой корректны")
}
