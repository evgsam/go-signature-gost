package main

import (
	"math/big"
)

// AddPoints складывает две точки P и Q на эллиптической кривой и возвращает результат R = P + Q.
// Использует стандартные формулы группового закона для эллиптических кривых над GF(p).
//
// Формулы:
//   - O + Q = Q, P + O = P (нейтральный элемент)
//   - P + (-P) = O (противоположные точки)
//   - P = Q: удвоение, lambda = (3*x1^2 + a) / (2*y1) mod p
//   - P != Q: сложение, lambda = (y2 - y1) / (x2 - x1) mod p
//   - x3 = lambda^2 - x1 - x2 mod p, y3 = lambda*(x1 - x3) - y1 mod p
func (curve *CurveParams) AddPoints(P, Q *Point) *Point {
	// Обработка нейтральных элементов (точка на бесконечности O)

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

	// Создаём копии координат, чтобы не изменять исходные точки
	x1 := new(big.Int).Set(P.X)
	y1 := new(big.Int).Set(P.Y)
	x2 := new(big.Int).Set(Q.X)
	y2 := new(big.Int).Set(Q.Y)

	// Проверка на противоположные точки: P = -Q -> P + Q = O
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
		// lambda = (3*x1^2 + a) / (2*y1) mod p
		numerator := new(big.Int).Exp(x1, big.NewInt(2), nil) // x1^2
		numerator.Mul(numerator, big.NewInt(3))               // 3*x1^2
		numerator.Add(numerator, curve.A)                     // 3*x1^2 + a
		numerator.Mod(numerator, curve.P)

		denominator := new(big.Int).Mul(y1, big.NewInt(2))              // 2*y1
		denominator.Mod(denominator, curve.P)                           // 2*y1 mod p
		denominatorInv := new(big.Int).ModInverse(denominator, curve.P) // (2*y1)^(-1) mod p
		if denominatorInv == nil {
			return NewInfinityPoint() // обратного элемента не существует
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
	x3.Sub(x3, x1)                         // lambda^2 - x1
	x3.Sub(x3, x2)                         // lambda^2 - x1 - x2
	x3.Mod(x3, curve.P)

	// y3 = lambda*(x1 - x3) - y1 mod p
	y3 := new(big.Int).Sub(x1, x3) // (x1 - x3)
	y3.Mul(lambda, y3)             // lambda*(x1 - x3)
	y3.Sub(y3, y1)                 // λ(x₁ - x₃) - y₁
	y3.Mod(y3, curve.P)

	return &Point{
		X:   x3,
		Y:   y3,
		Inf: false,
	}
}

// Double удваивает точку P: R = P + P.
// Является обёрткой над AddPoints(P, P).
func (curve *CurveParams) Double(P *Point) *Point {
	return curve.AddPoints(P, P)
}

// ScalarMult умножает точку P на скаляр k: R = k*P.
// Использует двоичный метод (double-and-add).
//
// Алгоритм:
//  1. Если k < 0, заменяем (k, P) на (-k, -P)
//  2. Если k = 0 или P = O, возвращаем O
//  3. Сканируем биты k: удваиваем base, добавляем result при единичном бите
func (curve *CurveParams) ScalarMult(k *big.Int, P *Point) *Point {
	if k.Sign() < 0 {
		// k отрицательное: k*P = (-k)*(-P)
		return curve.ScalarMult(new(big.Int).Neg(k), Neg(P, curve.P))
	}

	if k.Sign() == 0 {
		// 0*P = O
		return NewInfinityPoint()
	}
	if P == nil || P.Inf {
		return NewInfinityPoint()
	}

	// Двоичный метод (double-and-add)
	result := NewInfinityPoint()
	base := &Point{
		X:   new(big.Int).Set(P.X),
		Y:   new(big.Int).Set(P.Y),
		Inf: P.Inf,
	}
	kCopy := new(big.Int).Set(k)

	for kCopy.Sign() > 0 {
		if kCopy.Bit(0) == 1 {
			result = curve.AddPoints(result, base)
		}
		base = curve.Double(base)
		kCopy.Rsh(kCopy, 1)
	}
	return result
}
