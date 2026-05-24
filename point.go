package main

import (
	"math/big"
)

// Point представляет точку на эллиптической кривой над конечным полем.
// Точка на бесконечности (нейтральный элемент группы) обозначается Inf = true.
type Point struct {
	X   *big.Int // x-координата точки
	Y   *big.Int // y-координата точки
	Inf bool     // true для точки на бесконечности
}

// NewInfinityPoint создаёт точку на бесконечности (нейтральный элемент группы).
// Точка на бесконечности - аналог нуля в аддитивной группе точек кривой.
func NewInfinityPoint() *Point {
	return &Point{Inf: true}
}

// Neg возвращает точку, противоположную P: -P = (x, -y mod p).
// Для точки на бесконечности возвращает точку на бесконечности.
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

// Equal сравнивает две точки на кривой.
// Возвращает true, если точки совпадают (включая случай точки на бесконечности).
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
