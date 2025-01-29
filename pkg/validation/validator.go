package validation

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrContainsInvalidChars = errors.New("song text contains invalid characters")
	ErrInvalidStructure     = errors.New("song text must have at least one non-empty line")
)

// ValidateSongText проверяет корректность текста песни.
func ValidateSongText(songText string) error {
	// Удалить лишние пробелы
	songText = strings.TrimSpace(songText)

	// Разбить текст на строки
	lines := strings.Split(songText, "\n")

	// Проверка, что есть хотя бы одна строка с содержимым
	hasValidLine := false
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			hasValidLine = true
			break
		}
	}

	if !hasValidLine {
		return ErrInvalidStructure // Текст не содержит ни одной строки с содержимым
	}

	// Проверка на запрещенные символы
	for _, r := range songText {
		if unicode.IsControl(r) && r != '\n' { // Разрешены только переносы строки
			return ErrContainsInvalidChars
		}
	}
	return nil
}
