package entity

import "errors"

var (
	LimitLengthIncorrect  = errors.New("лимит должен быть в диапазоне от 1 до 1000")
	OffsetLengthIncorrect = errors.New("сдвиг должен быть в диапазоне от 0 до 500")
)

type Sorting struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (s *Sorting) Validate() error {
	if s.Limit <= 0 || s.Limit > 1000 {
		return LimitLengthIncorrect
	}
	if s.Offset < 0 || s.Offset > 500 {
		return OffsetLengthIncorrect
	}

	return nil
}
