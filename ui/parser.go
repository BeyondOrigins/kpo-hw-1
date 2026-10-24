package ui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrEmptyLine     = errors.New("пустая команда")
	ErrUnclosedQuote = errors.New("не закрыта кавычка")
	ErrArgsCount     = errors.New("неверное число аргументов")
	ErrBadArg        = errors.New("некорректный аргумент")
)

// Command — разобранная строка ввода: имя команды и её аргументы.
type Command struct {
	Name string
	Args []string
}

// ParseLine разбивает строку на слова по пробелам. Слова с пробелами
// берутся в двойные кавычки: add scooter 1 "Ninebot Max" 7 1.
// Внутри кавычек \" и \\ экранируют кавычку и обратный слэш.
func ParseLine(line string) (Command, error) {
	tokens, err := tokenize(line)
	if err != nil {
		return Command{}, err
	}
	if len(tokens) == 0 {
		return Command{}, ErrEmptyLine
	}
	return Command{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func tokenize(line string) ([]string, error) {
	var (
		tokens   []string
		current  strings.Builder
		inQuote  bool
		escaped  bool
		hasToken bool
	)

	flush := func() {
		if hasToken {
			tokens = append(tokens, current.String())
			current.Reset()
			hasToken = false
		}
	}

	for _, r := range line {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\' && inQuote:
			escaped = true
		case r == '"':
			inQuote = !inQuote
			hasToken = true
		case unicode.IsSpace(r) && !inQuote:
			flush()
		default:
			current.WriteRune(r)
			hasToken = true
		}
	}
	if inQuote {
		return nil, ErrUnclosedQuote
	}
	flush()
	return tokens, nil
}

func (c Command) RequireArgs(n int) error {
	if len(c.Args) != n {
		return fmt.Errorf("%w: команда %q ожидает %d, получено %d",
			ErrArgsCount, c.Name, n, len(c.Args))
	}
	return nil
}

func (c Command) Arg(i int) (string, error) {
	if i < 0 || i >= len(c.Args) {
		return "", fmt.Errorf("%w: нет аргумента №%d", ErrArgsCount, i+1)
	}
	return c.Args[i], nil
}

func (c Command) IntArg(i int) (int, error) {
	s, err := c.Arg(i)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: аргумент №%d (%q) должен быть целым числом",
			ErrBadArg, i+1, s)
	}
	return n, nil
}
