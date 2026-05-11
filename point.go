package main

import (
	"math/big"
)

// Point представляет точку на эллиптической кривой.
type Point struct {
	X   *big.Int // x-координата точки
	Y   *big.Int // y-координата точки
	Inf bool     // true для точки на бесконечности
}

// NewInfinityPoint создаёт точку на бесконечности (нейтральный элемент).
func NewInfinityPoint() *Point {
	return &Point{Inf: true}
}

// Neg возвращает точку -P: (x, -y mod p).
func Neg(P *Point, p *big.Int) *Point {
	if P == nil || P.Inf {
		return NewInfinityPoint()
	}
	negY := new(big.Int).Neg(P.Y)
	negY.Mod(negY, p)
	return &Point{
		X:   new(big.Int).Set(P.X),
		Y:   negY,
		Inf: false,
	}
}

// Equal сравнивает точки (с учётом Inf).
func Equal(P, Q *Point) bool {
	if P == nil && Q == nil {
		return true
	}
	if P == nil || Q == nil {
		return false
	}
	if P.Inf && Q.Inf {
		return true
	}
	if P.Inf != Q.Inf {
		return false
	}
	return P.X.Cmp(Q.X) == 0 && P.Y.Cmp(Q.Y) == 0
}
