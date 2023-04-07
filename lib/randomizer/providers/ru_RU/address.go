package ru

import (
	"fmt"

	"github.com/jaswdr/faker"
)

const (
	cityPrefix        = "город"
	cityPrefixShort   = "г."
	regionSuffix      = "область"
	regionSuffixShort = "обл."
)

var (
	streetPrefix = []string{"пер.", "ул.", "пр.", "шоссе", "пл.", "бульвар", "въезд", "спуск", "проезд", "наб."}
	countries    = []string{
		"Абхазия", "Австралия", "Австрия", "Азербайджан", "Албания", "Алжир", "Американское Самоа", "Ангилья", "Ангола", "Андорра", "Антарктида", "Антигуа и Барбуда", "Аргентина", "Армения", "Аруба", "Афганистан",
		"Багамы", "Бангладеш", "Барбадос", "Бахрейн", "Беларусь", "Белиз", "Бельгия", "Бенин", "Бермуды", "Болгария", "Боливия", "Бонэйр, Синт-Эстатиус и Саба", "Босния и Герцеговина", "Ботсвана", "Бразилия", "Британская Территория в Индийском Океане", "Британские Виргинские Острова", "Бруней", "Буркина-Фасо", "Бурунди", "Бутан",
		"Вануату", "Ватикан", "Венгрия", "Венесуэла", "Великобритания", "Виргинские Острова Соединённых Штатов", "Вьетнам",
		"Габон", "Гаити", "Гайана", "Гамбия", "Гана", "Гваделупа", "Гватемала", "Гвинея", "Гвинея-Бисау", "Германия", "Гернси", "Гибралтар", "Гондурас", "Гонконг", "Гренада", "Гренландия", "Греция", "Грузия", "Гуам",
		"Дания", "Демократическая Республика Конго", "Джерси", "Джибути", "Доминика", "Доминиканская Республика",
		"Египет",
		"Замбия", "Западная Сахара", "Зимбабве",
		"Израиль", "Индия", "Индонезия", "Иордания", "Ирак", "Иран", "Ирландия", "Исландия", "Испания", "Италия",
		"Йемен",
		"Кабо-Верде", "Казахстан", "Камбоджа", "Камерун", "Канада", "Катар", "Кения", "Кипр", "Киргизия", "Кирибати", "Китай", "Кокосовые острова", "Колумбия", "Коморы", "Конго", "Корейская Народно-Демократическая Республика", "Корея", "Коста-Рика", "Кот-д\"Ивуар", "Куба", "Кувейт", "Кюрасао",
		"Лаос", "Латвия", "Лесото", "Либерия", "Ливан", "Ливия", "Литва", "Лихтенштейн", "Люксембург",
		"Маврикий", "Мавритания", "Мадагаскар", "Майотта", "Макао", "Малави", "Малайзия", "Мали", "Малые Тихоокеанские Отдаленные Острова Соединенных Штатов", "Мальдивы", "Мальта", "Марокко", "Мартиника", "Маршалловы Острова", "Мексика", "Микронезия", "Мозамбик", "Молдова", "Монако", "Монголия", "Монтсеррат", "Мьянма",
		"Намибия", "Науру", "Непал", "Нигер", "Нигерия", "Нидерланды", "Никарагуа", "Ниуэ", "Новая Зеландия", "Новая Каледония", "Норвегия",
		"Объединенные Арабские Эмираты", "Оман", "Острова Кайман", "Острова Кука", "Острова Теркс и Кайкос", "Остров Буве", "Остров Мэн", "Остров Норфолк", "Остров Рождества", "Остров Херд и Острова Макдональд",
		"Пакистан", "Палау", "Палестина", "Панама", "Папуа-Новая Гвинея", "Парагвай", "Перу", "Питкерн", "Польша", "Португалия", "Пуэрто-Рико",
		"Республика Македония", "Реюньон", "Россия", "Руанда", "Румыния",
		"Самоа", "Сан-Марино", "Сан-Томе и Принсипи", "Саудовская Аравия", "Свазиленд", "Святая Елена, Остров Вознесения, Тристан-да-кунья", "Северные Марианские Острова", "Сейшелы", "Сен-Бартелеми", "Сен-Мартен", "Сенегал", "Сент-Винсент и Гренадины", "Сент-Китс и Невис", "Сент-Люсия", "Сент-Пьер и Микелон", "Сербия", "Сингапур", "Сирийская Арабская Республика", "Словакия", "Словения", "Соединенные Штаты Америки", "Соломоновы Острова", "Сомали", "Судан", "Суринам", "Сьерра-Леоне",
		"Таджикистан", "Таиланд", "Тайвань", "Танзания", "Тимор-лесте", "Того", "Токелау", "Тонга", "Тринидад и Тобаго", "Тувалу", "Тунис", "Туркмения", "Турция",
		"Уганда", "Узбекистан", "Украина", "Уоллис и Футуна", "Уругвай",
		"Фарерские острова", "Фиджи", "Филиппины", "Финляндия", "Фолклендские острова", "Франция", "Французская Гвиана", "Французская Полинезия", "Французские Южные Территории",
		"Хорватия",
		"Центрально-Африканская Республика",
		"Чад", "Черногория", "Чехия", "Чили",
		"Швейцария", "Швеция", "Шпицберген и Ян-Майен", "Шри-Ланка",
		"Эквадор", "Экваториальная Гвинея", "Эландские Острова", "Эль-Сальвадор", "Эритрея", "Эстония", "Эфиопия",
		"Южная Африка", "Южная Джорджия и Южные Сандвичевы Острова", "Южная Осетия", "Южный Судан",
		"Ямайка", "Япония",
	}
	regions = []string{
		"Амурская", "Архангельская", "Астраханская",
		"Белгородская", "Брянская",
		"Владимирская", "Волгоградская", "Вологодская", "Воронежская",
		"Ивановская", "Иркутская",
		"Калининградская", "Калужская", "Кемеровская", "Кировская", "Костромская", "Курганская", "Курская",
		"Ленинградская", "Липецкая",
		"Магаданская", "Московская", "Мурманская",
		"Нижегородская", "Новгородская", "Новосибирская",
		"Омская", "Оренбургская", "Орловская",
		"Пензенская", "Псковская",
		"Ростовская", "Рязанская",
		"Самарская", "Саратовская", "Сахалинская", "Свердловская", "Смоленская",
		"Тамбовская", "Тверская", "Томская", "Тульская", "Тюменская",
		"Ульяновская",
		"Челябинская", "Читинская",
		"Ярославская",
	}
	cities = []string{
		"Балашиха", "Видное", "Волоколамск", "Воскресенск", "Дмитров", "Домодедово", "Дорохово", "Егорьевск", "Зарайск",
		"Истра", "Кашира", "Клин", "Коломна", "Красногорск", "Лотошино", "Луховицы", "Люберцы", "Можайск", "Москва",
		"Мытищи", "Наро-Фоминск", "Ногинск", "Одинцово", "Озёры", "Орехово-Зуево", "Павловский Посад", "Подольск",
		"Пушкино", "Раменское", "Сергиев Посад", "Серебряные Пруды", "Серпухов", "Солнечногорск", "Ступино", "Талдом",
		"Чехов", "Шатура", "Шаховская", "Щёлково",
	}
	streets = []string{
		"Косиора", "Ладыгина", "Ленина", "Ломоносова", "Домодедовская", "Гоголя", "1905 года", "Чехова", "Сталина",
		"Космонавтов", "Гагарина", "Славы", "Бухарестская", "Будапештсткая", "Балканская",
	}
)

type address struct {
	f *faker.Faker
}

// AddressRu Возвращает структуру для создания случайных элементов адреса
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) AddressRu() *address {
	return &address{f: rf.Faker}
}

// Country Возвращает случайную страну из списка
func (a *address) Country() string {
	return a.f.RandomStringElement(countries)
}

// Region Возвращает случайную область из списка
//
// withSuffix - использовать суффикс "область", short - использовать короткий суффикс "обл."
func (a *address) Region(withSuffix, short bool) string {
	reg := a.f.RandomStringElement(regions)

	if !withSuffix {
		return reg
	}

	suf := regionSuffix
	if short {
		suf = regionSuffixShort
	}

	return fmt.Sprintf("%s %s", reg, suf)
}

// RegionCode Возвращает случайный код региона
func (a *address) RegionCode() int {
	return a.f.IntBetween(1, 92)
}

// City Возвращает случайный город из списка
//
// withPrefix - использовать префикс "город", short - использовать короткий префикс "г."
func (a *address) City(withPrefix, short bool) string {
	city := a.f.RandomStringElement(cities)

	if !withPrefix {
		return city
	}

	pref := cityPrefix
	if short {
		pref = cityPrefixShort
	}

	return fmt.Sprintf("%s %s", pref, city)
}

// Street Возвращает случайную улицу из списка
//
// withPrefix - использовать случайный префикс из списка streetPrefix
func (a *address) Street(withPrefix bool) string {
	str := a.f.RandomStringElement(streets)

	if !withPrefix {
		return str
	}

	pref := a.f.RandomStringElement(streetPrefix)

	return fmt.Sprintf("%s %s", pref, str)
}

// BuildingNumber Возвращает случайный номер дома в диапазоне от 1 до 199 включительно
func (a *address) BuildingNumber() int {
	return a.f.IntBetween(1, 199)
}

// PostCode Возвращает случайный почтовый индекс из диапазона допустимых индексов Почты России
func (a *address) PostCode() int {
	return a.f.IntBetween(100000, 694923)
}

// Full Возвращает полный адрес, сформированный случайным образом
func (a *address) Full() string {
	return a.String()
}

// String Возвращает строковое представление случайным образом сформированного адреса
func (a *address) String() string {
	return fmt.Sprintf(
		"%d, %s, %s, %s, %d", // 477813, Оренбургская область, город Чехов, проезд Гагарина, 134
		a.PostCode(),
		a.Region(true, false),
		a.City(true, false),
		a.Street(true),
		a.BuildingNumber(),
	)
}