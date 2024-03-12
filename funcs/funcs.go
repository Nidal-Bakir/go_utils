package funcUtils

import "unicode/utf8"

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
