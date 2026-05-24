package main

import (
	"math/big"

	"github.com/mikhirev/gostribog" // Библиотека для хеширования Стрибог
)

// hashToNumber хеширует сообщение и преобразует хеш в число по модулю Q.
// Реализует шаг хеширования по ГОСТ Р 34.10-2012.
//
// Алгоритм:
//  1. Вычисляется хеш msg алгоритмом Стрибог-256
//  2. Байтовый хеш преобразуется в большое целое число
//  3. Результат берется по модулю Q (порядок подгруппы)
//  4. Если результат = 0, заменяется на 1 (по требованиям ГОСТ)
func (curve *CurveParams) hashToNumber(msg []byte) *big.Int {
	// Создаём хешер Стрибог-256
	hasher := gostribog.New256()
	hasher.Write(msg)
	hash := hasher.Sum(nil)

	// Преобразуем байтовый хеш в большое целое число
	e := new(big.Int).SetBytes(hash)
	// Берём по модулю Q (порядок подгруппы)
	e.Mod(e, curve.Q)

	// По ГОСТ: если e = 0 mod Q, то e = 1
	if e.Sign() == 0 {
		e.SetInt64(1)
	}
	return e
}
