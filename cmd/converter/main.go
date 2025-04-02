package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	const lenArgs = 3

	if len(os.Args) != lenArgs {
		fmt.Println("Формат использования: <имя программы> <режим> <число>")
		fmt.Println("Режимы: 'b' - конвертировать бинарное число, 'd' - конвертировать десятичное число")

		return
	}

	mode := os.Args[1]
	input := os.Args[2]

	switch mode {
	case "b":
		decimal, err := binaryToDecimal(input)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Бинарное %s в десятичной системе: %d\n", input, decimal)

	case "d":
		decimal, err := strconv.Atoi(input)
		if err != nil || decimal < 0 {
			fmt.Println("error: ожидается целое число >= 0")
			return
		}

		binary := decimalToBinary(decimal)
		fmt.Printf("Число %d в двоичной системе: %s\n", decimal, binary)

	default:
		fmt.Println("error: ожидается режим 'b' или 'd'")
	}
}

func binaryToDecimal(binary string) (int, error) {
	decimal := 0
	binaryWeight := 1

	for i := len(binary) - 1; i >= 0; i-- {
		if binary[i] != '0' && binary[i] != '1' {
			return -1, errors.New("error: ожидается бинарное число")
		}

		if binary[i] == '1' {
			decimal += binaryWeight
		}

		binaryWeight *= 2
	}

	return decimal, nil
}

func decimalToBinary(decimal int) string {
	if decimal == 0 {
		return "0"
	}

	const base = 2

	binary := ""

	for decimal > 0 {
		bit := decimal % base
		binary = string('0'+rune(bit)) + binary
		decimal /= base
	}

	return binary
}
