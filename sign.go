package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Sign создаёт ЭЦП для хеша сообщения e по ГОСТ Р 34.10-2012.
// Возвращает пару (r, s) — компоненты подписи.
//
// Параметры:
//   - d: закрытый ключ (секретное число)
//   - e: хеш сообщения по модулю Q
//
// Возвращает:
//   - r: координата X точки k*G по модулю Q
//   - s: (r*d + k*e) mod Q
//   - err: ошибка генерации случайных данных
//
// Алгоритм (ГОСТ Р 34.10-2012):
//  1. Генерируем случайное k в диапазоне [1, Q-1]
//  2. Вычисляем точку C = k*G
//  3. r = C.X mod Q
//  4. Если r = 0, повторяем с шага 1
//  5. s = (r*d + k*e) mod Q
//  6. Если s = 0, повторяем с шага 1
//  7. Возвращаем (r, s)
func (curve *CurveParams) Sign(d, e *big.Int) (r, s *big.Int, err error) {
	Q := curve.Q
	if Q == nil || Q.Sign() <= 0 {
		return nil, nil, fmt.Errorf("некорректный порядок Q")
	}

	// Базовая точка G
	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	for {
		// Генерируем случайное k в диапазоне [1, Q-1]
		k, err := rand.Int(rand.Reader, Q)
		if err != nil {
			return nil, nil, err
		}
		if k.Sign() == 0 {
			continue
		}

		// Вычисляем точку C = k*G
		C := curve.ScalarMult(k, G)

		// r = C.X mod Q
		r = new(big.Int).Mod(C.X, Q)
		if r.Sign() == 0 {
			continue
		}

		// s = (r*d + k*e) mod Q
		s = new(big.Int).Add(new(big.Int).Mul(r, d), new(big.Int).Mul(k, e))
		s.Mod(s, Q)
		if s.Sign() != 0 {
			return r, s, nil
		}
	}
}
