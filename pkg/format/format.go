package format

import (
	"fmt"
	"strconv"
)

func FormatLargeNumber(num string) string {
	// Парсим строку в число с плавающей точкой
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "N/A" // Возвращаем "N/A", если число невалидно
	}

	// Форматируем число в зависимости от его величины
	switch {
	case f >= 1e12: // Триллионы
		return fmt.Sprintf("%.2fT", f/1e12)
	case f >= 1e9: // Миллиарды
		return fmt.Sprintf("%.2fB", f/1e9)
	case f >= 1e6: // Миллионы
		return fmt.Sprintf("%.2fM", f/1e6)
	default: // Меньше миллиона
		return fmt.Sprintf("%.2f", f)
	}
}

func FormatLargeNumberForPercent(num string) string {
	// Парсим строку в число с плавающей точкой
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "N/A" // Возвращаем "N/A", если число невалидно
	}

	// Форматируем число в зависимости от его величины
	switch {
	case f > 0:
		return fmt.Sprintf("+%.2f", f)
	default:
		return fmt.Sprintf("%.2f", f)
	}

}

func Float(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
