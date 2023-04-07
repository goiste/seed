package ru

import (
	"fmt"
	"strconv"

	"github.com/jaswdr/faker"
)

var (
	companyPrefixes = []string{
		"ООО", "ЗАО", "ООО Компания", "ОАО", "ПАО", "МКК", "МФО",
	}
	companyNameElements = []string{
		"ЖелДор", "Гараж", "Цемент", "Асбоцемент", "Строй", "Лифт", "Креп", "Авто", "Теле", "Транс", "Алмаз", "Метиз",
		"Мотор", "Рос", "Тяж", "Тех", "Сантех", "Урал", "Башкир", "Тверь", "Казань", "Обл", "Бух", "Хоз", "Электро",
		"Текстиль", "Восток", "Орион", "Юпитер", "Финанс", "Микро", "Радио", "Мобайл", "Дизайн", "Метал", "Нефть",
		"Телеком", "Инфо", "Сервис", "Софт", "IT", "Рыб", "Глав", "Вектор", "Рем", "Гор", "Газ", "Монтаж", "Мор",
		"Реч", "Флот", "Cиб", "Каз", "Инж", "Вод", "Пив", "Хмель", "Мяс", "Томск", "Омск", "Север", "Лен",
	}
	companyNameSuffixes = []string{
		"Маш", "Наладка", "Экспедиция", "Пром", "Комплекс", "Машина", "Снос", "-М", "Лизинг", "Траст", "Снаб", "-H",
		"Трест", "Банк", "Опт", "Проф", "Сбыт", "Центр",
	}
	jobTitles = []string{
		"Администратор", "Арт-директор", "Архивист", "Бариста", "Бармен", "Бизнес-аналитик", "Бухгалтер", "Ветеринар",
		"Водитель", "Водолаз", "Геймдизайнер", "Детектив", "Диджей", "Диктор", "Зубной техник", "Лесоруб", "Лингвист",
		"Машинист", "Менеджер", "Музыкант", "Научный сотрудник", "Офис-менеджер", "Печатник", "Пианист", "Писатель",
		"Программист", "Продюсер", "Промоутер", "Психолог", "Редактор", "Системный аналитик", "Стилист", "Столяр",
		"Сторож", "Технический писатель", "Учёный", "Физик", "Финансовый советник", "Фотограф", "Фрезеровщик",
		"Художник", "Экономист", "Электромонтёр",
	}
)

type company struct {
	f *faker.Faker
}

// CompanyRu Возвращает структуру для генерации случайных данных компании
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) CompanyRu() *company {
	return &company{f: rf.Faker}
}

// JobTitle Возвращает случайную должность
func (c *company) JobTitle() string {
	return c.f.RandomStringElement(jobTitles)
}

// INN Возвращает случайный (валидный) ИНН юридического лица
func (c *company) INN(regionCode int) string {
	if regionCode == 0 {
		regionCode = c.f.IntBetween(1, 92)
	}
	inn := fmt.Sprintf("%02d%07d", regionCode, c.f.IntBetween(0, 9999999))
	return inn + innCheckSum(inn)
}

// KPP Возвращает случайный (валидный) КПП на основе ИНН
func (c *company) KPP(inn string) string {
	if len(inn) < 4 {
		inn = c.INN(0)
	}
	return fmt.Sprintf("%s01001", inn[:4])
}

// Name Возвращает случайное наименование компании
func (c *company) Name() string {
	return c.String()
}

// String Строковое представление компании (сгенерированное название)
func (c *company) String() string {
	names := []string{
		fmt.Sprintf("%s %s", // ООО Мобайл
			c.f.RandomStringElement(companyPrefixes),
			c.f.RandomStringElement(companyNameElements),
		),
		fmt.Sprintf("%s %s%s", // ЗАО СантехФинанс
			c.f.RandomStringElement(companyPrefixes),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
		),
		fmt.Sprintf("%s %s%s%s", // ЗАО РосСофтСбыт
			c.f.RandomStringElement(companyPrefixes),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameSuffixes),
		),
		fmt.Sprintf("%s %s%s%s", // ПАО ВостокТяжГаз
			c.f.RandomStringElement(companyPrefixes),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
		),
		fmt.Sprintf("%s %s%s%s%s", // ОАО МикроАвтоСибОпт
			c.f.RandomStringElement(companyPrefixes),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameElements),
			c.f.RandomStringElement(companyNameSuffixes),
		),
	}
	return c.f.RandomStringElement(names)
}

// Вычисление контрольной цифры ИНН юридического лица
//
// see https://ru.wikipedia.org/wiki/%D0%98%D0%B4%D0%B5%D0%BD%D1%82%D0%B8%D1%84%D0%B8%D0%BA%D0%B0%D1%86%D0%B8%D0%BE%D0%BD%D0%BD%D1%8B%D0%B9_%D0%BD%D0%BE%D0%BC%D0%B5%D1%80_%D0%BD%D0%B0%D0%BB%D0%BE%D0%B3%D0%BE%D0%BF%D0%BB%D0%B0%D1%82%D0%B5%D0%BB%D1%8C%D1%89%D0%B8%D0%BA%D0%B0%23%D0%92%D1%8B%D1%87%D0%B8%D1%81%D0%BB%D0%B5%D0%BD%D0%B8%D0%B5_%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%BE%D0%BB%D1%8C%D0%BD%D1%8B%D1%85_%D1%86%D0%B8%D1%84%D1%80
func innCheckSum(inn string) string {
	if len(inn) != 9 {
		return ""
	}

	multipliers := map[int]int{1: 2, 2: 4, 3: 10, 4: 3, 5: 5, 6: 9, 7: 4, 8: 6, 9: 8}

	sum := 0

	for i := 1; i <= 9; i++ {
		intVal, err := strconv.Atoi(string(inn[i-1]))
		if err != nil {
			return ""
		}

		sum += intVal * multipliers[i]
	}

	return fmt.Sprintf("%d", (sum%11)%10)
}
