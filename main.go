package main

import (
	"fmt"
	"math/big"
	"testing"
)

type CurveParams struct {
	OID string   // Идентификатор набора параметров
	P   *big.Int // модуль поля
	A   *big.Int // коэффициент a
	B   *big.Int // коэффициент b
	M   *big.Int // порядок группы точек кривой
	Q   *big.Int // порядок подгруппы, в которой строится ЭЦП
	GX  *big.Int // x‑координата базовой точки P
	GY  *big.Int // y‑координата базовой точки P
}

// Структура для точки кривой
type Point struct {
	X   *big.Int // x-координата точки
	Y   *big.Int // y-координата точки
	Inf bool     // true для точки на бесконечности
}

// Функция для создания точки на бесконечности
func NewInfinityPoint() *Point {
	return &Point{Inf: true}
}

// Функция для сложения двух точек P и Q на кривой
func (curve *CurveParams) AddPoints(P, Q *Point) *Point {
	// Обработка нейтральных элементов
	// O + Q = Q
	if P == nil || P.Inf {	// если P - точка на бесконечности, то результатом будет Q
		if Q == nil || Q.Inf { // если Q тоже точка на бесконечности, то результат - точка на бесконечности
			return NewInfinityPoint()
		}
		return &Point{	// если P - точка на бесконечности, а Q - нет, то результат - Q
			X:   new(big.Int).Set(Q.X),
			Y:   new(big.Int).Set(Q.Y),
			Inf: false,
		}
	}
	// P + O = P
	if Q == nil || Q.Inf { // если Q - точка на бесконечности, то результатом будет P
		return &Point{ 
			X:   new(big.Int).Set(P.X),
			Y:   new(big.Int).Set(P.Y),
			Inf: false,
		}
	}

	//создаем копии координат, чтобы не изменять исходные точки
	x1 := new(big.Int).Set(P.X)
	y1 := new(big.Int).Set(P.Y)
	x2 := new(big.Int).Set(Q.X)
	y2 := new(big.Int).Set(Q.Y)

	// 2. Проверка на противоположные точки: P = -Q
	if x1.Cmp(x2) == 0 { // сравниваем что x1 == x2
		ySum := new(big.Int).Add(y1, y2) //скалярное сложение y1 + y2
		ySum.Mod(ySum, curve.P)          // берем по модулю p, так как мы работаем в поле GF(p)
		if ySum.Sign() == 0 {            // если y1 + y2 == 0, то P и Q - противоположные точки, и их сумма - точка на бесконечности
			return NewInfinityPoint() // Если y1 + y2 != 0, то P и Q - это одна и та же точка, и мы должны выполнить удвоение
		}
	}

	var lambda *big.Int // переменная для хранения коэффициента наклона (lambda)

	// 3. Случай удвоения: P == Q
	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 { // если P и Q - это одна и та же точка, то выполняем удвоение
		if y1.Sign() == 0 { // если y1 == 0, то результат - точка на бесконечности
			return NewInfinityPoint()
		}
		// lambda = (3*x1² + a) / (2*y1) mod p
		numerator := new(big.Int).Exp(x1, big.NewInt(2), nil) //x1^2
		numerator.Mul(numerator, big.NewInt(3))               // 3*x1^2
		numerator.Add(numerator, curve.A)                     // 3*x1^2 + a
		numerator.Mod(numerator, curve.P)					 // берем по модулю p

		denominator := new(big.Int).Mul(y1, big.NewInt(2))              // 2*y1
		denominator.Mod(denominator, curve.P)                           // 2*y1 mod p
		denominatorInv := new(big.Int).ModInverse(denominator, curve.P) // переворачиваем знаменатель, чтобы из деления получить умножение: (2*y1)^(-1) mod p
		if denominatorInv == nil {                                      // если обратного элемента не существует, то результат -
			// точка на бесконечности
			return NewInfinityPoint()
		}
		lambda = new(big.Int).Mul(numerator, denominatorInv) // умножаем числитель на обратный элемент знаменателя: (3*x1^2 + a) * (2*y1)^(-1)
		lambda.Mod(lambda, curve.P)                          // берем по модулю p
	} else {
		// 4. Случай общего сложения: P != Q
		// Вычисляем lambda для общего сложения: lambda = (y2 - y1) / (x2 - x1) mod p
		numerator:= new(big.Int).Sub(y2,y1) // (y2 - y1)
		numerator.Mod(numerator,curve.P)	// (y2 - y1) mod p

		denominator:=new(big.Int).Sub(x2,x1) // (x2 - x1)
		denominator.Mod(denominator,curve.P) // (x2 - x1) mod p
		denominatorInv := new(big.Int).ModInverse(denominator, curve.P) // переворачиваем знаменатель чтобы из деления получить умножение: (x2 - x1)^(-1) mod p
		if denominatorInv == nil {                                      // если обратного элемента не существует, то результат -
			// точка на бесконечности
			return NewInfinityPoint()
		}
		lambda=new(big.Int).Mul(numerator,denominatorInv) // умножаем числитель на обратный элемент знаменателя: (y2 - y1) * (x2 - x1)^(-1) 
		lambda.Mod(lambda,curve.P)	// берем по модулю p
	}
	
	 // 5. Вычисление координат результата
    // x3 = lambda^2 - x1 - x2 mod p
	x3:=new(big.Int).Mul(lambda,lambda) // lambda^2
	x3.Sub(x3,x1) // lambda^2 - x1
	x3.Sub(x3,x2) // lambda^2 - x1 - x2
	x3.Mod(x3,curve.P) // y3 = lambda*(x1 - x3) - y1 mod p
	// y3 = lambda*(x1 - x3) - y1 mod p
	y3:=new(big.Int).Sub(x1,x3) // (x1 - x3)
	y3.Mul(lambda,y3) //lambda*(x1 - x3)
	y3.Sub(y3,y1) // lambda*(x1 - x3) - y1
	y3.Mod(y3,curve.P) // y3= lambda*(x1 - x3) - y1 mod p

	 return &Point{	// возвращаем результат сложения точек P и Q
        X:   x3,
        Y:   y3,
        Inf: false,
    }
}

func mustParse(hexStr string) *big.Int {
	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		panic("неверная шестнадцатеричная строка: " + hexStr)
	}
	return n
}

func IsValidCurveParams(params *CurveParams) bool {
	x, y, p, a, b := params.GX, params.GY, params.P, params.A, params.B
	//lhs = y^2 mod p
	lhs := new(big.Int).Exp(y, big.NewInt(2), p)
	// rhs = (x^3 + a*x + b) mod p
	rhs := new(big.Int).Mul(a, x)
	rhs.Add(rhs, b)
	rhs.Add(rhs, new(big.Int).Exp(x, big.NewInt(3), nil))
	rhs.Mod(rhs, p)
	return lhs.Cmp(rhs) == 0
}


// Neg возвращает точку -P: (x, -y mod p)
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

// Equal сравнивает точки (с учётом Inf)
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

// IsOnCurve проверяет, лежит ли точка на эллиптической кривой
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

// Этот тестовый случай использует твои же CurveParams — его можно подключить к main.go.
func TestCurveParams_Add(t *testing.T) {
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

	if !IsValidCurveParams(params) {
		t.Fatal("параметры кривой не проходят проверку в тесте")
	}

	G := &Point{X: params.GX, Y: params.GY, Inf: false}
	O := NewInfinityPoint()

	// 1. G + O = G
	sum1 := params.AddPoints(G, O)
	if !Equal(sum1, G) {
		t.Errorf("G + O хотели получить G, получили (%X, %X)", sum1.X, sum1.Y)
	}

	// 2. O + G = G
	sum2 := params.AddPoints(O, G)
	if !Equal(sum2, G) {
		t.Errorf("O + G хотели получить G, получили (%X, %X)", sum2.X, sum2.Y)
	}

	// 3. G + (-G) = O
	negG := Neg(G, params.P)
	sum3 := params.AddPoints(G, negG)
	if !sum3.Inf {
		t.Errorf("G + (-G) хотели получить O, получили (%X, %X)", sum3.X, sum3.Y)
	}

	// 4. O + O = O
	sum4 := params.AddPoints(O, O)
	if !sum4.Inf {
		t.Errorf("O + O хотели получить O, получили (%X, %X)", sum4.X, sum4.Y)
	}

	// 5. 2G лежит на кривой
	doubleG := params.AddPoints(G, G)
	if !params.IsOnCurve(doubleG) {
		t.Errorf("2G не лежит на кривой: (%X, %X)", doubleG.X, doubleG.Y)
	}

	// 6. G + 2G = 3G, и 3G тоже лежит на кривой
	tripleG := params.AddPoints(G, doubleG)
	if !params.IsOnCurve(tripleG) {
		t.Errorf("3G не лежит на кривой: (%X, %X)", tripleG.X, tripleG.Y)
	}

	// 7. (G + 2G) + (-2G) = G
	neg2G := Neg(doubleG, params.P)
	sum5 := params.AddPoints(tripleG, neg2G)
	if !Equal(sum5, G) {
		t.Errorf("(G + 2G) + (-2G) хотели получить G, а получили (%X, %X)", sum5.X, sum5.Y)
	}
}


func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")
	test := &testing.T{}
	TestCurveParams_Add(test)
	if !test.Failed() {
		fmt.Println("Все тесты пройдены успешно")
	} else {
		fmt.Println("Некоторые тесты не прошли")
	}
}
