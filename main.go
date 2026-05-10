package main

import (
	"fmt"
	"math/big"
)

type CurveParams struct {
	OID *big.Int // Идентификатор набора параметров
	P   *big.Int // модуль поля
	A   *big.Int // коэффициент a
	B   *big.Int // коэффициент b
	M   *big.Int // порядок группы точек кривой
	Q   *big.Int // порядок подгруппы, в которой строится ЭЦП
	GX  *big.Int // x‑координата базовой точки P
	GY  *big.Int // y‑координата базовой точки P
}

func main() {
	fmt.Println("ЭЦП ГОСТ 256 бит")
}
