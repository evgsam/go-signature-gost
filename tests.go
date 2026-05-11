package main

import (
	"fmt"
	"testing"
)

// TestCurveParams_Add тестирует сложение точек на кривой.
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

	fmt.Println("Все тесты пройдены успешно")
}

// MustRunTests запускает все тесты и выводит результат.
func MustRunTests() {
	test := &testing.T{}
	TestCurveParams_Add(test)
	if !test.Failed() {
		fmt.Println("Все тесты пройдены успешно")
	} else {
		fmt.Println("Некоторые тесты не прошли")
	}
}
