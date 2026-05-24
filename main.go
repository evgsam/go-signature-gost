package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/mikhirev/gostribog" // Библиотека для хеширования Стрибогом
)

func (curve *CurveParams) hashToNumber(msg []byte) *big.Int { // Хешируем сообщение и преобразуем в число
	hasher := gostribog.New256()
	hasher.Write(msg)
	hash := hasher.Sum(nil)

	e := new(big.Int).SetBytes(hash)
	e.Mod(e, curve.Q)

	if e.Sign() == 0 { // по ГОСТ - если e равно 0
		e.SetInt64(1) // то e должно быть 1
	}
	return e
}

func (curve *CurveParams) Sign(d, e *big.Int) (r, s *big.Int, err error) { // Подпись сообщения
	Q := curve.Q
	if Q == nil || Q.Sign() <= 0 {
		return nil, nil, fmt.Errorf("некорректный порядок Q")
	}

	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	for {
		k, err := rand.Int(rand.Reader, Q) // Генерируем случайное k в диапазоне [1, Q-1]
		if err != nil {                    // Если произошла ошибка при генерации k, возвращаем ошибку
			return nil, nil, err
		}
		if k.Sign() == 0 { // Если k равно 0, генерируем новое k
			continue
		}

		C := curve.ScalarMult(k, G)  // Вычисляем точку C = k*G
		r = new(big.Int).Mod(C.X, Q) // r = C.X mod Q
		if r.Sign() == 0 {           // Если r равно 0, генерируем новое k
			continue
		}

		s = new(big.Int).Add(new(big.Int).Mul(r, d), new(big.Int).Mul(k, e)) // s = (r*d + k*e)
		s.Mod(s, Q)                                                          // s = (r*d + k*e) mod Q
		if s.Sign() != 0 {                                                   // Если s не равно 0, возвращаем r и s
			return r, s, nil
		}
	}
}

func (curve *CurveParams) Verify(H *Point, msg []byte, r, s *big.Int) bool { // Проверка подписи
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
	// G = (GX, GY) - базовая точка
	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	//C' = z1 * G + z2 * H
	p1 := curve.ScalarMult(z1, G)
	p2 := curve.ScalarMult(z2, H)
	Cp := curve.AddPoints(p1, p2)

	// Если C' – точка на бесконечности, подпись недействительна
	if Cp.Inf {
		return false
	}
	// r' = C'.X mod Q
	rPrime := new(big.Int).Mod(Cp.X, Q)
	return rPrime.Cmp(r) == 0
}

func (curve *CurveParams) GenerateKey() (d *big.Int, H *Point, err error) {
	Q := curve.Q
	if Q == nil || Q.Sign() <= 0 {
		return nil, nil, fmt.Errorf("некорректный порядок Q")
	}

	G := &Point{X: curve.GX, Y: curve.GY, Inf: false}

	for {
		d, err = rand.Int(rand.Reader, Q)
		if err != nil {
			return nil, nil, err
		}
		if d.Sign() != 0 {
			break
		}
	}
	H = curve.ScalarMult(d, G)
	return d, H, nil
}

func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")

	// Параметры кривой (тестовый набор)
	params := &CurveParams{
		OID: "1.2.643.2.2.35.1",
		P:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD97"),
		A:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD94"),
		B:   mustParse("A6"),
		M:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		Q:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		GX:  mustParse("1"),
		GY:  mustParse("8D91E471E0989CDA27DF505A453F2B7635294F2DDF23E3B122ACC99C9E9F1E14"),
	}

	// Генерация ключей
	d, H, err := params.GenerateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Закрытый ключ d: %X\n", d)
	fmt.Printf("Открытый ключ H: (%X, %X)\n", H.X, H.Y)

	// Проверка, что открытый ключ лежит на кривой
	if params.IsOnCurve(H) {
		fmt.Println("Открытый ключ корректен (лежит на кривой)")
	} else {
		fmt.Println("Ошибка: открытый ключ не лежит на кривой")
	}
}
