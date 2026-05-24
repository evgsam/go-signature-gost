package main

import "math/big"

// Verify проверяет ЭЦП (r, s) для сообщения msg по ГОСТ Р 34.10-2012.
//
// Параметры:
//   - H: открытая ключевая точка (H = d*G, где d — закрытый ключ)
//   - msg: исходное сообщение
//   - r, s: компоненты подписи
//
// Возвращает true, если подпись корректна.
//
// Алгоритм проверки (ГОСТ Р 34.10-2012):
//  1. Проверяем корректность r, s: 0 < r < Q, 0 < s < Q
//  2. e = hash(msg) mod Q
//  3. v = e^(-1) mod Q
//  4. z1 = s*v mod Q, z2 = -r*v mod Q
//  5. Cp = z1*G + z2*H
//  6. rPrime = Cp.X mod Q
//  7. rPrime == r -> подпись верна
func (curve *CurveParams) Verify(H *Point, msg []byte, r, s *big.Int) bool {
	Q := curve.Q

	// Проверяем, что r и s в правильном диапазоне (0 < r < Q и 0 < s < Q)
	if r.Sign() <= 0 || r.Cmp(Q) >= 0 {
		return false
	}
	if s.Sign() <= 0 || s.Cmp(Q) >= 0 {
		return false
	}

	// Вычисляем e = hash(msg) mod Q
	e := curve.hashToNumber(msg)

	// Вычисляем v = e^(-1) mod Q
	v := new(big.Int).ModInverse(e, Q)
	if v == nil {
		return false
	}

	// Вычисляем z1 = s * v mod Q и z2 = -r * v mod Q
	z1 := new(big.Int).Mul(s, v)
	z1.Mod(z1, Q)
	z2 := new(big.Int).Mul(r, v)
	z2.Neg(z2)
	z2.Mod(z2, Q)

	// G = (GX, GY) — базовая точка
	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	// C' = z1 * G + z2 * H
	p1 := curve.ScalarMult(z1, G)
	p2 := curve.ScalarMult(z2, H)
	Cp := curve.AddPoints(p1, p2)

	// Если C' — точка на бесконечности, подпись недействительна
	if Cp.Inf {
		return false
	}

	// r' = C'.X mod Q
	rPrime := new(big.Int).Mod(Cp.X, Q)
	return rPrime.Cmp(r) == 0
}
