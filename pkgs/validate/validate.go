package validate

import "errors"

func Email(s string) error {
	if len(s) < 1 {
		return errors.New("please fill this out")
	}
	return nil
}

func Password(s string) error {
	if len(s) < 1 {
		return errors.New("please fill this out")
	}
	return nil
}

func LengthBetween(s string, x int, y int) error {
	l := len(s)
	if l < x {
		return errors.New("not long enough")
	}
	if l > y {
		return errors.New("too long")
	}
	return nil
}