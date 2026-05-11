package main

import (
	"math/big"
)

// CurveParams содержит параметры эллиптической кривой для ЭЦП по ГОСТ.
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

// IsValidCurveParams проверяет, что базовая точка G лежит на кривой.
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
