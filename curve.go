package main

import (
	"math/big"
)

// CurveParams содержит параметры эллиптической кривой для ЭЦП по ГОСТ Р 34.10.
// Кривая задаётся уравнением: y^2 = x^3 + a*x + b (mod p)
type CurveParams struct {
	OID string   // Идентификатор набора параметров (ASN.1 OID)
	P   *big.Int // модуль поля (простое число p)
	A   *big.Int // коэффициент a уравнения кривой
	B   *big.Int // коэффициент b уравнения кривой
	M   *big.Int // порядок группы точек кривой
	Q   *big.Int // порядок подгруппы, в которой строится ЭЦП
	GX  *big.Int // x-координата базовой точки G
	GY  *big.Int // y-координата базовой точки G
}

// IsValidCurveParams проверяет, что базовая точка G лежит на кривой.
// Проверяется равенство: y^2 = x^3 + a*x + b (mod p)
func IsValidCurveParams(params *CurveParams) bool {
	x, y, p, a, b := params.GX, params.GY, params.P, params.A, params.B

	// lhs = y^2 mod p
	lhs := new(big.Int).Exp(y, big.NewInt(2), p)

	// rhs = (x^3 + a*x + b) mod p
	rhs := new(big.Int).Mul(a, x)
	rhs.Add(rhs, b)
	rhs.Add(rhs, new(big.Int).Exp(x, big.NewInt(3), nil))
	rhs.Mod(rhs, p)

	return lhs.Cmp(rhs) == 0
}
