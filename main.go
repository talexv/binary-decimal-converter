package main

func BinaryToDecimal(binary string) int {
	decimal := 0
	binaryWeight := 1

	for i := len(binary) - 1; i >= 0; i-- {
		if binary[i] == '1' {
			decimal += binaryWeight
		}

		binaryWeight *= 2
	}

	return decimal
}

func DecimalToBinary(decimal int) string {
	const base = 2

	binary := ""

	for decimal > 0 {
		bit := decimal % base
		binary = string('0'+rune(bit)) + binary
		decimal /= base
	}

	return binary
}
