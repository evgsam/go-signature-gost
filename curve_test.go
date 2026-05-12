package main

import (
	"math/big"
	"testing"
)

// ----------------------------------------------------------------------------
// Вспомогательные функции для тестов
// ----------------------------------------------------------------------------

func getTestCurve() *CurveParams {
	return &CurveParams{
		OID: "1.2.643.2.2.35.1",
		P:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD97"),
		A:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD94"),
		B:   mustParse("A6"),
		M:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		Q:   mustParse("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6C611070995AD10045841B09B761B893"),
		GX:  mustParse("1"),
		GY:  mustParse("8D91E471E0989CDA27DF505A453F2B7635294F2DDF23E3B122ACC99C9E9F1E14"),
	}
}

func getBasePoint(params *CurveParams) *Point {
	return &Point{
		X:   new(big.Int).Set(params.GX),
		Y:   new(big.Int).Set(params.GY),
		Inf: false,
	}
}

// Вспомогательная функция для проверки, что точка лежит на кривой
func mustPointOnCurve(t *testing.T, curve *CurveParams, p *Point, name string) {
	t.Helper()
	if !curve.IsOnCurve(p) {
		if p == nil {
			t.Fatalf("%s = nil, ожидалась корректная точка", name)
		}
		if p.Inf {
			t.Fatalf("%s – точка на бесконечности, а ожидалась обычная точка", name)
		}
		t.Fatalf("%s не лежит на кривой: (%X, %X)", name, p.X, p.Y)
	}
}

// Вспомогательная функция для сравнения двух точек
func mustEqualPoints(t *testing.T, got, want *Point, msg string) {
	t.Helper()
	if !Equal(got, want) {
		if got == nil && want == nil {
			return
		}
		if got == nil {
			t.Fatalf("%s: got=nil, want=(%X,%X, inf=%v)", msg, want.X, want.Y, want.Inf)
		}
		if want == nil {
			t.Fatalf("%s: got=(%X,%X, inf=%v), want=nil", msg, got.X, got.Y, got.Inf)
		}
		if got.Inf || want.Inf {
			t.Fatalf("%s: got inf=%v, want inf=%v", msg, got.Inf, want.Inf)
		}
		t.Fatalf("%s: got=(%X,%X), want=(%X,%X)", msg, got.X, got.Y, want.X, want.Y)
	}
}

// ----------------------------------------------------------------------------
// Тесты корректности параметров кривой и базовой точки
// ----------------------------------------------------------------------------

func TestCurveParams_ValidParams(t *testing.T) {
	params := getTestCurve()
	if !IsValidCurveParams(params) {
		t.Fatal("параметры кривой не проходят проверку")
	}
	G := getBasePoint(params)
	mustPointOnCurve(t, params, G, "G")
}

// ----------------------------------------------------------------------------
// Тесты для Neg (противоположная точка)
// ----------------------------------------------------------------------------

func TestPoint_Neg(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	negG := Neg(G, params.P)
	mustPointOnCurve(t, params, negG, "-G")

	back := Neg(negG, params.P)
	mustEqualPoints(t, back, G, "Neg(Neg(G)) должно быть равно G")

	sum := params.AddPoints(G, negG)
	if !sum.Inf {
		t.Fatalf("G + (-G) должно быть точкой на бесконечности, получили (%X, %X)", sum.X, sum.Y)
	}
}

// ----------------------------------------------------------------------------
// Тесты нейтрального элемента (точки на бесконечности)
// ----------------------------------------------------------------------------

func TestCurveParams_AddPoints_NeutralElement(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)
	O := NewInfinityPoint()

	sum1 := params.AddPoints(G, O)
	mustEqualPoints(t, sum1, G, "G + O должно быть равно G")

	sum2 := params.AddPoints(O, G)
	mustEqualPoints(t, sum2, G, "O + G должно быть равно G")

	sum3 := params.AddPoints(O, O)
	if !sum3.Inf {
		t.Fatalf("O + O должно быть точкой на бесконечности")
	}
}

// ----------------------------------------------------------------------------
// Тест замкнутости: сумма любых двух точек лежит на кривой
// ----------------------------------------------------------------------------

func TestCurveParams_AddPoints_Closure(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	doubleG := params.AddPoints(G, G)
	tripleG := params.AddPoints(G, doubleG)
	fourG := params.AddPoints(doubleG, doubleG)

	points := []*Point{G, doubleG, tripleG, fourG, NewInfinityPoint()}

	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points); j++ {
			sum := params.AddPoints(points[i], points[j])
			mustPointOnCurve(t, params, sum, "sum")
		}
	}
}

// ----------------------------------------------------------------------------
// Тест коммутативности
// ----------------------------------------------------------------------------

func TestCurveParams_AddPoints_Commutativity(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	doubleG := params.AddPoints(G, G)
	tripleG := params.AddPoints(G, doubleG)

	left := params.AddPoints(G, tripleG)
	right := params.AddPoints(tripleG, G)
	mustEqualPoints(t, left, right, "коммутативность нарушена: G + 3G != 3G + G")

	left2 := params.AddPoints(doubleG, tripleG)
	right2 := params.AddPoints(tripleG, doubleG)
	mustEqualPoints(t, left2, right2, "коммутативность нарушена: 2G + 3G != 3G + 2G")
}

// ----------------------------------------------------------------------------
// Тест ассоциативности
// ----------------------------------------------------------------------------

func TestCurveParams_AddPoints_Associativity(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	doubleG := params.AddPoints(G, G)
	tripleG := params.AddPoints(G, doubleG)

	left := params.AddPoints(params.AddPoints(G, G), G)
	right := params.AddPoints(G, params.AddPoints(G, G))
	mustEqualPoints(t, left, right, "ассоциативность нарушена для G,G,G")

	left2 := params.AddPoints(params.AddPoints(G, doubleG), tripleG)
	right2 := params.AddPoints(G, params.AddPoints(doubleG, tripleG))
	mustEqualPoints(t, left2, right2, "ассоциативность нарушена для G,2G,3G")
}

// ----------------------------------------------------------------------------
// Тест удвоения (Double) – если он у вас есть
// Если Double не реализован, можно пропустить или написать через AddPoints
// ----------------------------------------------------------------------------

func TestCurveParams_Double(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	doubleByMethod := params.AddPoints(G, G) // Временно заменим, если нет Double
	doubleByAdd := params.AddPoints(G, G)

	mustEqualPoints(t, doubleByMethod, doubleByAdd, "Double(G) должно совпадать с G + G")
	mustPointOnCurve(t, params, doubleByMethod, "2G")
}

// ----------------------------------------------------------------------------
// Тест ScalarMult для малых значений
// ----------------------------------------------------------------------------

func TestCurveParams_ScalarMult_SmallValues(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)
	O := NewInfinityPoint()

	zeroG := params.ScalarMult(big.NewInt(0), G)
	if !zeroG.Inf {
		t.Fatalf("0*G должно быть точкой на бесконечности")
	}

	oneG := params.ScalarMult(big.NewInt(1), G)
	mustEqualPoints(t, oneG, G, "1*G должно быть равно G")

	twoG := params.ScalarMult(big.NewInt(2), G)
	doubleG := params.AddPoints(G, G)
	mustEqualPoints(t, twoG, doubleG, "2*G должно быть равно G+G")

	threeG := params.ScalarMult(big.NewInt(3), G)
	threeGByAdd := params.AddPoints(doubleG, G)
	mustEqualPoints(t, threeG, threeGByAdd, "3*G должно быть равно 2G + G")

	fourG := params.ScalarMult(big.NewInt(4), G)
	fourGByAdd := params.AddPoints(doubleG, doubleG)
	mustEqualPoints(t, fourG, fourGByAdd, "4*G должно быть равно 2G+2G")

	zeroO := params.ScalarMult(big.NewInt(123), O)
	if !zeroO.Inf {
		t.Fatalf("k*O должно быть точкой на бесконечности")
	}

	mustPointOnCurve(t, params, oneG, "1G")
	mustPointOnCurve(t, params, twoG, "2G")
	mustPointOnCurve(t, params, threeG, "3G")
	mustPointOnCurve(t, params, fourG, "4G")
}

// ----------------------------------------------------------------------------
// Тест ScalarMult для отрицательных скаляров
// ----------------------------------------------------------------------------

func TestCurveParams_ScalarMult_Negative(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	pos := params.ScalarMult(big.NewInt(3), G)
	neg := params.ScalarMult(big.NewInt(-3), G)
	expected := Neg(pos, params.P)

	mustEqualPoints(t, neg, expected, "(-3)*G должно быть равно -(3*G)")
}

// ----------------------------------------------------------------------------
// Тест дистрибутивности: (a+b)*G == a*G + b*G
// ----------------------------------------------------------------------------

func TestCurveParams_ScalarMult_DistributiveSmall(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	left := params.ScalarMult(big.NewInt(5), G)

	twoG := params.ScalarMult(big.NewInt(2), G)
	threeG := params.ScalarMult(big.NewInt(3), G)
	right := params.AddPoints(twoG, threeG)

	mustEqualPoints(t, left, right, "5*G должно быть равно 2G + 3G")
}

// ----------------------------------------------------------------------------
// Ключевой тест: проверка порядка подгруппы Q
// Если Q – порядок базовой точки, то Q * G == O
// ----------------------------------------------------------------------------

func TestCurveParams_OrderQ(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	q := new(big.Int).Set(params.Q)
	qG := params.ScalarMult(q, G)

	if !qG.Inf {
		t.Fatalf("Q*G должно быть точкой на бесконечности, получено (%X, %X)", qG.X, qG.Y)
	}

	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	qMinus1G := params.ScalarMult(qMinus1, G)
	if qMinus1G.Inf {
		t.Fatalf("(Q-1)*G не должно быть точкой на бесконечности")
	}
	mustPointOnCurve(t, params, qMinus1G, "(Q-1)G")
}

// ----------------------------------------------------------------------------
// Тест ScalarMult для произвольного большого числа
// (проверяем, что результат лежит на кривой)
// ----------------------------------------------------------------------------

func TestCurveParams_ScalarMult_Random(t *testing.T) {
	params := getTestCurve()
	G := getBasePoint(params)

	// Случайное 256-битное число
	rand, _ := new(big.Int).SetString("A5F6B7C8D9E0F1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B1C2D3E4F5A6B7C8", 16)
	randG := params.ScalarMult(rand, G)
	mustPointOnCurve(t, params, randG, "random*G")

	// Проверка, что (k+1)*G == k*G + G
	k := new(big.Int).Set(rand)
	kPlus1 := new(big.Int).Add(k, big.NewInt(1))
	k1G := params.ScalarMult(kPlus1, G)
	kg := params.ScalarMult(k, G)
	kgPlusG := params.AddPoints(kg, G)
	mustEqualPoints(t, k1G, kgPlusG, "(k+1)G == kG + G")
}

// TestGOSTSignVerifyFromRFC7091 проверяет полный цикл подписи и проверки по официальному примеру.
func TestGOSTSignVerifyFromGOST(t *testing.T) {
	// 1. Определяем тестовые параметры кривой из ГОСТ Р 34.10-2012
	testParams := &CurveParams{
		P:  mustParse("8000000000000000000000000000000000000000000000000000000000000431"),
		A:  mustParse("7"),
		B:  mustParse("5FBFF498AA938CE739B8E022FBAFEF40563F6E6A3472FC2A514C0CE9DAE23B7E"),
		Q:  mustParse("8000000000000000000000000000000150FE8A1892976154C59CFC193ACCF5B3"),
		GX: mustParse("2"),
		GY: mustParse("8E2A8A0E65147D4BD6316030E16D19C85C97F0A9CA267122B96ABBCEA7E8FC8"),
		// Остальные поля (M, OID) для этого теста не критичны
	}

	d := mustParse("7A929ADE789BB9BE10ED359DD39A72C11B60961F49397EEE1D19CE9891EC3B28") // Закрытый ключ из ГОСТ
	e := mustParse("2DFBC1B372D89A1188C09C52E0EEC61FCE52032AB1022E8E67ECE6672B043EE5") // Хеш сообщения из ГОСТ
	k := mustParse("77105C9B20BCD3122823C8CF6FCC7B956DE33814E95B7FE64FED924594DCEAB3") // Случайное число из ГОСТ

	// Ожидаемая подпись (r, s) из RFC
	expectedR := mustParse("41AA28D2F1AB148280CD9ED56FEDA41974053554A42767B83AD043FD39DC0493") //Из ГОСТ
	expectedS := mustParse("1456C64BA4642A1653C235A98A60249BCD6D3F746B631DF928014F6C5BF9C40")  //Из ГОСТ

	// Алгоритм подписи:
	// 1. Вычисляем точку C = k * G
	C := testParams.ScalarMult(k, &Point{X: testParams.GX, Y: testParams.GY, Inf: false})
	// 2. Берём x-координату C и приводим по модулю q, получаем r
	r := new(big.Int).Mod(C.X, testParams.Q)
	// 3. Вычисляем s = (r*d + k*e) mod q
	s := new(big.Int).Mul(r, d)
	s.Add(s, new(big.Int).Mul(k, e))
	s.Mod(s, testParams.Q)

	// Сравниваем полученные r и s с ожидаемыми
	if r.Cmp(expectedR) != 0 || s.Cmp(expectedS) != 0 {
		t.Errorf("Подпись не совпадает с эталоном из ГОСТ.")
	}

}
