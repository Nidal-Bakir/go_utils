package funcUtils

import (
    "unicode/utf8"
    "github.com/shopspring/decimal"
)

// language agnostic function to trim a string str to maxLength
// e.g 1: TrimString("abcd",2) // ab
// e.g 2: TrimString("Ἄγγελος",2) // Ἄγ
func TrimString(str string, maxLength int) string {
	byteLenOfStr := len(str)
	accumulateByteSize := 0
	byteSize := 0

	// as log as we are lees the then the maxLength
	// and we did not exceed the number of bytes in str
	for runeCount := 0; runeCount < maxLength && accumulateByteSize < byteLenOfStr; runeCount++ {
		_, byteSize = utf8.DecodeRuneInString(str[accumulateByteSize:])
		accumulateByteSize += byteSize
	}

	return str[:accumulateByteSize]
}


func calcWithFractionRemainder(cost float64, installmentsCount int) (priceWithFR decimal.Decimal, price decimal.Decimal) {
	dCost := decimal.NewFromFloat(cost)
	if installmentsCount <= 1 {
		return dCost, dCost
	}
	dInstallmentsCount := decimal.NewFromInt(int64(installmentsCount))
	dInstallmentsCountSub1 := decimal.NewFromInt(int64(installmentsCount - 1))

	// n = (cost/c) * (c -1) :: e.g: cost=1, c=3;; => n = 0.66666666666666...
	// => z = cost - n :: e.g: (1 - 0.6666666...) => 0.3333333333333334
	price = dCost.Div(dInstallmentsCount)
	priceWithFR = dCost.Sub(price.Mul(dInstallmentsCountSub1))

	return priceWithFR, price
}
