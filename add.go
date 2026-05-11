package main

import (
	"math/big"
)

// AddPoints складывает две точки P и Q на кривой и возвращает результат.
func (curve *CurveParams) AddPoints(P, Q *Point) *Point {
	// Обработка нейтральных элементов
	// O + Q = Q
	if P == nil || P.Inf {
		if Q == nil || Q.Inf {
			return NewInfinityPoint()
		}
		return &Point{
			X:   new(big.Int).Set(Q.X),
			Y:   new(big.Int).Set(Q.Y),
			Inf: false,
		}
	}
	// P + O = P
	if Q == nil || Q.Inf {
		return &Point{
			X:   new(big.Int).Set(P.X),
			Y:   new(big.Int).Set(P.Y),
			Inf: false,
		}
	}

	// создаем копии координат, чтобы не изменять исходные точки
	x1 := new(big.Int).Set(P.X)
	y1 := new(big.Int).Set(P.Y)
	x2 := new(big.Int).Set(Q.X)
	y2 := new(big.Int).Set(Q.Y)

	// Проверка на противоположные точки: P = -Q
	if x1.Cmp(x2) == 0 {
		ySum := new(big.Int).Add(y1, y2)
		ySum.Mod(ySum, curve.P)
		if ySum.Sign() == 0 {
			return NewInfinityPoint()
		}
	}

	var lambda *big.Int

	// Случай удвоения: P == Q
	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		if y1.Sign() == 0 {
			return NewInfinityPoint()
		}
		// lambda = (3*x1² + a) / (2*y1) mod p
		numerator := new(big.Int).Exp(x1, big.NewInt(2), nil) // x1^2
		numerator.Mul(numerator, big.NewInt(3))               // 3*x1^2
		numerator.Add(numerator, curve.A)                     // 3*x1^2 + a
		numerator.Mod(numerator, curve.P)

		denominator := new(big.Int).Mul(y1, big.NewInt(2))              // 2*y1
		denominator.Mod(denominator, curve.P)                           // 2*y1 mod p
		denominatorInv := new(big.Int).ModInverse(denominator, curve.P) // (2*y1)^(-1) mod p
		if denominatorInv == nil {
			return NewInfinityPoint()
		}
		lambda = new(big.Int).Mul(numerator, denominatorInv)
		lambda.Mod(lambda, curve.P)
	} else {
		// Случай общего сложения: P != Q
		// lambda = (y2 - y1) / (x2 - x1) mod p
		numerator := new(big.Int).Sub(y2, y1) // (y2 - y1)
		numerator.Mod(numerator, curve.P)

		denominator := new(big.Int).Sub(x2, x1) // (x2 - x1)
		denominator.Mod(denominator, curve.P)   // (x2 - x1) mod p
		denominatorInv := new(big.Int).ModInverse(denominator, curve.P)
		if denominatorInv == nil {
			return NewInfinityPoint()
		}
		lambda = new(big.Int).Mul(numerator, denominatorInv)
		lambda.Mod(lambda, curve.P)
	}

	// Вычисление координат результата
	// x3 = lambda^2 - x1 - x2 mod p
	x3 := new(big.Int).Mul(lambda, lambda) // lambda^2
	x3.Sub(x3, x1)                        // lambda^2 - x1
	x3.Sub(x3, x2)                        // lambda^2 - x1 - x2
	x3.Mod(x3, curve.P)

	// y3 = lambda*(x1 - x3) - y1 mod p
	y3 := new(big.Int).Sub(x1, x3) // (x1 - x3)
	y3.Mul(lambda, y3)             // lambda*(x1 - x3)
	y3.Sub(y3, y1)                 // lambda*(x1 - x3) - y1
	y3.Mod(y3, curve.P)

	return &Point{
		X:   x3,
		Y:   y3,
		Inf: false,
	}
}
