package ru

import (
	"github.com/jaswdr/faker"
)

// RusFaker Обёртка для faker.Faker
//
// Добавляет русскую локализацию для некоторых типов данных
type RusFaker struct {
	*faker.Faker
}

// NewFaker Создаёт экземпляр RusFaker
func NewFaker(f *faker.Faker) *RusFaker {
	return &RusFaker{f}
}
