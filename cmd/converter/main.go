package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals, mnd // will be considered
var charToValue = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3,
	'4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, 'A': 10, 'B': 11,
	'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

func isValidBase(base int) bool {
	return base == 2 || base == 8 || base == 10 || base == 16
}

func convertToDecimal(numberStr, fromBase string) (int, error) {
	numberStr = strings.ToUpper(numberStr)
	decimal := 0
	posWeight := 1

	base, err := strconv.Atoi(fromBase)
	if err != nil || !isValidBase(base) {
		return -1, errors.New("error: поддерживаются только '2, 8, 10, 16' системы")
	}

	for i := len(numberStr) - 1; i >= 0; i-- {
		val, ok := charToValue[rune(numberStr[i])]
		if !ok || val >= base {
			return -1, fmt.Errorf("error: Недопустимый символ '%c' для системы '%s'", numberStr[i], fromBase)
		}

		decimal += val * posWeight
		posWeight *= base
	}

	return decimal, nil
}

func convertFromDecimal(decimal int, toBase string) (string, error) {
	if decimal == 0 {
		return "0", nil
	}

	base, err := strconv.Atoi(toBase)
	if err != nil || !isValidBase(base) {
		return "", errors.New("error: поддерживаются только '2, 8, 10, 16' системы")
	}

	result := ""

	valueToChar := make(map[int]rune, len(charToValue))

	for k, v := range charToValue {
		valueToChar[v] = k
	}

	for decimal > 0 {
		remainder := decimal % base
		char := valueToChar[remainder]
		result = string(char) + result
		decimal /= base
	}

	return result, nil
}

func convert(numberStr, fromBase, toBase string) (string, error) {
	decimal, err := convertToDecimal(numberStr, fromBase)
	if err != nil {
		return "", err
	}

	result, err := convertFromDecimal(decimal, toBase)
	if err != nil {
		return "", err
	}

	return result, nil
}

func convertFromStdin(fromBase, toBase string) error {
	fmt.Println("Ожидаются числа через пробел или Enter. Завершить Ctrl+Z")

	for {
		var numStr string

		_, err := fmt.Fscan(os.Stdin, &numStr)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		result, err := convert(numStr, fromBase, toBase)
		if err != nil {
			return err
		}

		fmt.Printf("%s (%s) --> %s (%s)\n", numStr, fromBase, result, toBase)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "converter",
		Usage: "Конвертирует числа между системами счисления",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "from",
				Aliases:  []string{"f"},
				Usage:    "Исходная система счисления (2, 8, 10, 16)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "to",
				Aliases:  []string{"t"},
				Usage:    "Целевая система счисления (2, 8, 10, 16)",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			fromBase := cCtx.String("from")
			toBase := cCtx.String("to")
			args := cCtx.Args().Slice()

			if len(args) > 0 {
				for _, numStr := range args {
					result, err := convert(numStr, fromBase, toBase)
					if err != nil {
						return cli.Exit(err, 1)
					}

					fmt.Printf("%s (%s) --> %s (%s)\n", numStr, fromBase, result, toBase)
				}

				return nil
			}

			err := convertFromStdin(fromBase, toBase)
			if err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
