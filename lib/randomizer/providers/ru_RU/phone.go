package ru

import (
	"fmt"

	"github.com/jaswdr/faker"
)

type phone struct {
	f *faker.Faker
}

// PhoneRu Возвращает структуру для генерации случайного номера телефона
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) PhoneRu() *phone {
	return &phone{f: rf.Faker}
}

// Number Возвращает случайный номер телефона с указанными (или случайными) префиксом страны и кодом
func (p *phone) Number(code, prefix int) string {
	if code == 0 {
		code = p.f.IntBetween(900, 999)
	}

	if prefix == 0 {
		prefix = 7
	}

	firstPart := "8800"
	if prefix != 8 {
		firstPart = fmt.Sprintf("+%d%d", prefix, code)
	}
	return fmt.Sprintf("%s%07d", firstPart, p.f.IntBetween(1, 9999999))
}

func (p *phone) String() string {
	return p.Number(0, 0)
}
