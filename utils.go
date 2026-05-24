package main

import (
	"fmt"
	"math/big"
)

// mustParse преобразует шестнадцатеричную строку в *big.Int.
// Используется для загрузки параметров кривой из hex-представления.
// При ошибке парсинга паникует.
func mustParse(hexStr string) *big.Int {
	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		panic("неверная шестнадцатеричная строка: " + hexStr)
	}
	return n
}

// PrintPoint выводит точку в удобочитаемом формате.
// Для точки на бесконечности выводит "O".
func PrintPoint(label string, p *Point) {
	if p.Inf {
		fmt.Printf("%s: O (точка на бесконечности)\n", label)
		return
	}
	fmt.Printf("%s: (%X, %X)\n", label, p.X, p.Y)
}
