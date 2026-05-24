package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateKey генерирует пару ключей (закрытый d, открытый H) по ГОСТ Р 34.10-2012.
//
// Алгоритм:
//  1. Генерируем случайное число d в диапазоне [1, Q-1] - закрытый ключ
//  2. Вычисляем H = d*G - открытая ключевая точка
//
// Возвращает:
//   - d: закрытый ключ (секретное число)
//   - H: открытая ключевая точка на кривой
//   - err: ошибка генерации случайных данных
func (curve *CurveParams) GenerateKey() (d *big.Int, H *Point, err error) {
	Q := curve.Q
	if Q == nil || Q.Sign() <= 0 {
		return nil, nil, fmt.Errorf("некорректный порядок Q")
	}

	// Базовая точка G
	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	for {
		// Генерируем случайный закрытый ключ d в диапазоне [1, Q-1]
		d, err = rand.Int(rand.Reader, Q)
		if err != nil {
			return nil, nil, err
		}
		if d.Sign() != 0 {
			break
		}
	}

	// Вычисляем открытый ключ H = d*G
	H = curve.ScalarMult(d, G)
	return d, H, nil
}
