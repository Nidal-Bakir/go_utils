package funcUtils

import (
	"io"
	"slices"
	"sync"
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


func CalcWithFractionRemainder(cost float64, installmentsCount int) (priceWithFR decimal.Decimal, price decimal.Decimal) {
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

func SplitAndGather[I any, O any](in <-chan I, count int, proc func(I) O) []O {
	out := make(chan O)

	var wg sync.WaitGroup
	wg.Add(count)

	for range count {
		go func() {
			defer wg.Done()
			for v := range in {
				out <- proc(v)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	result := make([]O, 0)
	for v := range out {
		result = append(result, v)
	}
	return result
}


func CountLettersFromIReader(r io.Reader, out map[string]int) error {
	isPossibleRune := func(p []byte) bool {
		l := len(p)
		return l > 0 && l <= 3 && utf8.RuneStart(p[0])
	}

	buf := make([]byte, 2*1024)

	var possibleRune []byte = nil

	for {
		n, err := r.Read(buf)

		if possibleRune != nil {
			buf = slices.Concat(possibleRune, buf[:n])
			n = len(buf)
			possibleRune = nil
		}

		totalReadBytes := 0
		for totalReadBytes < n {
			bufScope := buf[totalReadBytes:n]

			r, readBytes := utf8.DecodeRune(bufScope)

			if r != utf8.RuneError {
				totalReadBytes += readBytes
				out[string(r)]++
				continue
			}

			// we have rune error. can this bytes assemble a rune in the next Read?
			// it is possible that the rune constructed of 4/3/2 bytes and we only
			// have the first bytes of the rune. We need to get the remaining bytes
			// in the next Read
			if isPossibleRune(bufScope) {
				possibleRune = make([]byte, len(bufScope))
				copy(possibleRune, bufScope)
				break
			}

			return fmt.Errorf("error while decoding a rune: %v", bufScope)
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
