package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

// func main() {
// 	const lenArgs = 3

// 	if len(os.Args) != lenArgs {
// 		fmt.Fprintln(os.Stderr, "Формат использования: <имя программы> <режим> <число>")
// 		fmt.Fprintln(os.Stderr, "Режимы: 'b' - конвертировать бинарное число, 'd' - конвертировать десятичное число")

// 		os.Exit(1)
// 	}

// 	mode := os.Args[1]
// 	input := os.Args[2]

// 	switch mode {
// 	case "b":
// 		decimal, err := binaryToDecimal(input)
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			os.Exit(1)
// 		}

// 		fmt.Fprintf(os.Stdout, "Бинарное %s в десятичной системе: %s\n", input, decimal)

// 	case "d":
// 		binary, err := decimalToBinary(input)
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			os.Exit(1)
// 		}

// 		fmt.Fprintf(os.Stdout, "Число %s в двоичной системе: %s\n", input, binary)

// 	default:
// 		fmt.Fprintln(os.Stderr, "error: ожидается режим 'b' или 'd'")
// 		os.Exit(1)
// 	}
// }

func main() {
	app := &cli.App{
		Name:      "converter",
		Usage:     "Конвертирует числа между бинарной и десятичной систем счисления",
		ArgsUsage: "arguments...",
		Commands: []*cli.Command{
			{
				Name:      "binary",
				Aliases:   []string{"b"},
				Usage:     "Конвертирует бинарное число в десятичное",
				ArgsUsage: "<бинарное число>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 1 {
						return cli.Exit("Ожидается 1 аргумент - <бинарное число>", 1)
					}
					input := cCtx.Args().Get(0)
					decimal, err := binaryToDecimal(input)
					if err != nil {
						return cli.Exit(err, 1)
					}

					fmt.Printf("Бинарное %s в десятичной системе: %s\n", input, decimal)
					return nil
				},
			},
			{
				Name:      "decimal",
				Aliases:   []string{"d"},
				Usage:     "Конвертирует десятичное число в бинарное",
				ArgsUsage: "<десятичное число>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 1 {
						return cli.Exit("Ожидается 1 аргумент - <десятичное число>", 1)
					}
					input := cCtx.Args().Get(0)
					binary, err := decimalToBinary(input)
					if err != nil {
						return cli.Exit(err, 1)
					}

					fmt.Printf("Число %s в двоичной системе: %s\n", input, binary)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func binaryToDecimal(binary string) (string, error) {
	decimal := 0
	binaryWeight := 1

	for i := len(binary) - 1; i >= 0; i-- {
		if binary[i] != '0' && binary[i] != '1' {
			return "", errors.New("error: ожидается бинарное число")
		}

		if binary[i] == '1' {
			decimal += binaryWeight
		}

		binaryWeight *= 2
	}

	return strconv.Itoa(decimal), nil
}

func decimalToBinary(decimal string) (string, error) {
	num, err := strconv.Atoi(decimal)
	if err != nil || num < 0 {
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
