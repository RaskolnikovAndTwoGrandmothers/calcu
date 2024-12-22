package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "У тебя своя голова на плечах есть, но если очень надо, то тебе сюда")
	})

	http.HandleFunc("/api/v1/calculate", calculateHandler)

	http.ListenAndServe(":8642", nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Неправильно введено выражение"})
		return
	}

	result, calcErr := performCalculation(req.Expression)
	if calcErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: calcErr.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"Результат": result})
}

func performCalculation(expression string) (float64, error) {
	// Убираем пробелы из выражения (если пользователь ввел выражение без пробелов)
	expression = strings.ReplaceAll(expression, " ", "")

	// Токенизация: разделение на числа и операторы
	tokens := tokenizeExpression(expression)

	if len(tokens) < 3 {
		return 0, fmt.Errorf("недопустимое выражение")
	}

	// 1-й проход: обработка "*" и "/"
	stack := []string{}
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "*" || tokens[i] == "/" {
			if len(stack) == 0 {
				return 0, fmt.Errorf("ошибка в синтаксисе: %s", tokens[i])
			}

			// Берем предыдущий элемент в стеке и текущий токен
			left, err := strconv.ParseFloat(stack[len(stack)-1], 64)
			if err != nil {
				return 0, fmt.Errorf("неверное число: %s", stack[len(stack)-1])
			}
			right, err := strconv.ParseFloat(tokens[i+1], 64)
			if err != nil {
				return 0, fmt.Errorf("неверное число: %s", tokens[i+1])
			}

			// Вычисляем результат (для * или /)
			var result float64
			if tokens[i] == "*" {
				result = left * right
			} else if tokens[i] == "/" {
				if right == 0 {
					return 0, fmt.Errorf("деление на 0")
				}
				result = left / right
			}

			// Заменяем последний элемент в стеке на результат
			stack[len(stack)-1] = fmt.Sprintf("%f", result)
			i++ // Пропускаем следующий токен, так как он уже обработан
		} else {
			// Оператор или число, добавляем в стек
			stack = append(stack, tokens[i])
		}
	}

	// 2-й проход: обработка "+" и "-"
	result, err := strconv.ParseFloat(stack[0], 64)
	if err != nil {
		return 0, fmt.Errorf("неверное число в результате: %s", stack[0])
	}
	for i := 1; i < len(stack); i += 2 {
		operator := stack[i]
		operand, err := strconv.ParseFloat(stack[i+1], 64)
		if err != nil {
			return 0, fmt.Errorf("неверное число: %s", stack[i+1])
		}

		switch operator {
		case "+":
			result += operand
		case "-":
			result -= operand
		default:
			return 0, fmt.Errorf("неподдерживаемая операция: %s", operator)
		}
	}

	return result, nil
}

// Вспомогательная функция для токенизации выражения
func tokenizeExpression(expression string) []string {
	var tokens []string
	num := ""

	for _, ch := range expression {
		if (ch >= '0' && ch <= '9') || ch == '.' {
			num += string(ch)
		} else {
			if num != "" {
				tokens = append(tokens, num)
				num = ""
			}
			tokens = append(tokens, string(ch))
		}
	}

	if num != "" {
		tokens = append(tokens, num)
	}

	return tokens
}
