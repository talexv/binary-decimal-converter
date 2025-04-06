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
		fmt.Fprintln(os.Stderr, "Формат использования: <имя программы> <режим> <число>")
		fmt.Fprintln(os.Stderr, "Режимы: 'b' - конвертировать бинарное число, 'd' - конвертировать десятичное число")

		os.Exit(1)
	}

	mode := os.Args[1]
	input := os.Args[2]

	switch mode {
	case "b":
		decimal, err := binaryToDecimal(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Бинарное %s в десятичной системе: %d\n", input, decimal)

	case "d":
		binary, err := decimalToBinary(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Число %s в двоичной системе: %s\n", input, binary)

	default:
		fmt.Fprintln(os.Stderr, "error: ожидается режим 'b' или 'd'")
		os.Exit(1)
	}
}

func binaryToDecimal(binary string) (int, error) {
	decimal := 0
	binaryWeight := 1

	for i := len(binary) - 1; i >= 0; i-- {
		if binary[i] != '0' && binary[i] != '1' {
			// return -1, fmt.Errorf("%w: ожидается бинарное число", os.ErrInvalid)
			return -1, errors.New("error: ожидается бинарное число")
		}

		if binary[i] == '1' {
			decimal += binaryWeight
		}

		binaryWeight *= 2
	}

	return decimal, nil
}

func decimalToBinary(decimal string) (string, error) {
	num, err := strconv.Atoi(decimal)
	if err != nil || num < 0 {
		// return "", fmt.Errorf("%w: ожидается целое число >= 0", os.ErrInvalid)
		return "", errors.New("error: ожидается целое число >= 0")
	}

	if num == 0 {
		return "0", nil
	}

	const base = 2

	binary := ""

	for num > 0 {
		bit := num % base
		binary = string('0'+rune(bit)) + binary
		num /= base
	}

	return binary, nil
}
