package main

import (
	"math/big"
)

// IsOnCurve проверяет, лежит ли точка на эллиптической кривой.
func (curve *CurveParams) IsOnCurve(p *Point) bool {
	if p.Inf {
		return true
	}
	// y^2 mod p
	lhs := new(big.Int).Exp(p.Y, big.NewInt(2), curve.P)
	// x^3 + a*x + b mod p
	x3 := new(big.Int).Exp(p.X, big.NewInt(3), curve.P)
	ax := new(big.Int).Mul(curve.A, p.X)
	ax.Mod(ax, curve.P)
	rhs := new(big.Int).Add(x3, ax)
	rhs.Add(rhs, curve.B)
	rhs.Mod(rhs, curve.P)
	return lhs.Cmp(rhs) == 0
}
