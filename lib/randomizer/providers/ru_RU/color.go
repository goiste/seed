package ru

import (
	"github.com/jaswdr/faker"
)

var (
	colors = []string{
		"красный", "зелёный", "синий", "голубой", "малиновый", "жёлтый", "черный", "серый", "белый", "тёмно-бордовый",
		"тёмно-синий", "оливковый", "фиолетовый", "бирюзовый", "зеленовато-голубой", "серебро", "фуксия", "янтарный",
	}
)

type color struct {
	f *faker.Faker
}

// ColorRu Возвращает структуру для генерации случайных названий цветов
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) ColorRu() *color {
	return &color{f: rf.Faker}
}

// Name Возвращает случайное название цвета
func (c *color) Name() string {
	return c.String()
}

// String Возвращает случайное название цвета
func (c *color) String() string {
	return c.f.RandomStringElement(colors)
}
