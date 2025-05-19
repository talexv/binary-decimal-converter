package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals, mnd // will be considered
var charToValue = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3,
	'4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, 'A': 10, 'B': 11,
	'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

//nolint:gochecknoglobals // will be considered
var valueToChar = func() map[int]rune {
	m := make(map[int]rune, len(charToValue))
	for k, v := range charToValue {
		m[v] = k
	}

	return m
}()

var ErrBadInput = errors.New("ошибка во входных данных")

func isValidBase(base int) bool {
	return base == 2 || base == 8 || base == 10 || base == 16
}

func convertToDecimal(numberStr, fromBase string) (int, error) {
	numberStr = strings.ToUpper(numberStr)
	decimal := 0
	posWeight := 1

	base, err := strconv.Atoi(fromBase)
	if err != nil || !isValidBase(base) {
		return -1, fmt.Errorf("%w: поддерживаются только '2, 8, 10, 16' системы", ErrBadInput)
	}

	for i := len(numberStr) - 1; i >= 0; i-- {
		val, ok := charToValue[rune(numberStr[i])]
		if !ok || val >= base {
			return -1, fmt.Errorf("%w: недопустимый символ '%c' для системы '%s'", ErrBadInput, numberStr[i], fromBase)
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
		return "", fmt.Errorf("%w: поддерживаются только '2, 8, 10, 16' системы", ErrBadInput)
	}

	result := ""

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

func printConverted(w io.Writer, numberStr, fromBase, result, toBase string) error {
	_, err := fmt.Fprintf(w, "%s (%s) --> %s (%s)\n", numberStr, fromBase, result, toBase)
	return err
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

		err = printConverted(os.Stdout, numStr, fromBase, result, toBase)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:gochecknoglobals // will be considered
var convertForm = `
	<form action="/convert" method="POST">
		<input name="number" placeholder="Число">
		<input name="from" placeholder="Из системы">
		<input name="to" placeholder="В систему">
		<button type="submit">Конвертировать</button>
	</form>
`

func convertHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, convertForm)

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
			return
		}

		number := r.FormValue("number")
		fromBase := r.FormValue("from")
		toBase := r.FormValue("to")

		result, err := convert(number, fromBase, toBase)
		if err != nil {
			if errors.Is(err, ErrBadInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		err = printConverted(w, number, fromBase, result, toBase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, convertForm)
	}
}

//nolint:mnd // will be considered
func main() {
	if len(os.Args) > 1 && os.Args[1] == "http" {
		mux := http.NewServeMux()
		mux.HandleFunc("/convert", convertHandler)

		server := http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		fmt.Println("HTTP сервер запущен на http://localhost:8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}

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

			if len(args) == 0 {
				err := convertFromStdin(fromBase, toBase)
				if err != nil {
					return cli.Exit(err, 1)
				}

				return nil
			}

			for _, numStr := range args {
				result, err := convert(numStr, fromBase, toBase)
				if err != nil {
					return cli.Exit(err, 1)
				}

				err = printConverted(os.Stdout, numStr, fromBase, result, toBase)
				if err != nil {
					return cli.Exit(err, 1)
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
