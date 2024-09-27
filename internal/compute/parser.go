package compute

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	availSymbolsRegexp = `[a-zA-Zа-яА-Я0-9!?,.;:\"\'\ *#-=_@+№%$^\/\\|\[\]]`
	
	GetCmd string = "get"
	SetCmd string = "set"
	DeleteCmd string = "delete"
)

type Parser interface {
	ParseArgs(s string) (string, []string, error)
}

type RequestParser struct{}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

func (b *RequestParser) ParseArgs(s string) (string, []string, error) {
	rawArgs := strings.Split(s, " ")
	command := strings.Trim(rawArgs[0], "\t\n ")
	args := make([]string, 0, len(rawArgs[1:]))
	for _, arg := range rawArgs[1:] {
		args = append(args, strings.Trim(arg, "\t\n "))
	}

	err := b.validate(command, args)
	if err != nil {
		return "", nil, err
	}

	return command, args, nil
}

func (b *RequestParser) validate(command string, args []string) error {
	ln := len(args)

	switch command {
	case GetCmd:
		if ln != 1 {
			return fmt.Errorf("ожидается 1 аргумент, получено %d", ln)
		}
	case SetCmd:
		if ln != 2 {
			return fmt.Errorf("ожидается 2 аргумента, получено %d", ln)
		}
	case DeleteCmd:
		if ln != 1 {
			return fmt.Errorf("ожидается 1 аргумент, получено %d", ln)
		}
	default:
		return errors.New("Неизвестная команда")
	}
	
	r, err := regexp.Compile(availSymbolsRegexp)
	if err != nil {
		return errors.New("regexp.Compile error")
	}

	for i := 1; i < len(args); i++ {
		match := r.MatchString(args[i])
		if !match {
			return errors.New("Неизвестные символы в аргументах")
		}
	}

	return nil
}


