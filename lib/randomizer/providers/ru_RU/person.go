package ru

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/jaswdr/faker"
)

const (
	GenderMale   = "m"
	GenderFemale = "f"

	FormatLastFirst = "lf"
	FormatFirstLast = "fl"
	FormatFIO       = "fio"
	FormatIOF       = "iof"

	femaleLastNameSuffix = "а"
)

var (
	firstNamesMale = []string{
		"Абрам", "Август", "Адам", "Адриан", "Аким", "Александр", "Алексей", "Альберт", "Ананий", "Анатолий", "Андрей",
		"Антон", "Антонин", "Аполлон", "Аркадий", "Арсений", "Артемий", "Артур", "Артём", "Афанасий", "Богдан",
		"Болеслав", "Борис", "Бронислав", "Вадим", "Валентин", "Валериан", "Валерий", "Василий", "Вениамин", "Викентий",
		"Виктор", "Виль", "Виталий", "Витольд", "Влад", "Владимир", "Владислав", "Владлен", "Всеволод", "Вячеслав",
		"Гавриил", "Гарри", "Геннадий", "Георгий", "Герасим", "Герман", "Глеб", "Гордей", "Григорий", "Давид", "Дан",
		"Даниил", "Данила", "Денис", "Дмитрий", "Добрыня", "Донат", "Евгений", "Егор", "Ефим", "Захар", "Иван", "Игнат",
		"Игнатий", "Игорь", "Илларион", "Илья", "Иммануил", "Иннокентий", "Иосиф", "Ираклий", "Кирилл", "Клим",
		"Константин", "Кузьма", "Лаврентий", "Лев", "Леонид", "Макар", "Максим", "Марат", "Марк", "Матвей", "Милан",
		"Мирослав", "Михаил", "Назар", "Нестор", "Никита", "Никодим", "Николай", "Олег", "Павел", "Платон", "Прохор",
		"Пётр", "Радислав", "Рафаил", "Роберт", "Родион", "Роман", "Ростислав", "Руслан", "Сава", "Савва", "Святослав",
		"Семён", "Сергей", "Спартак", "Станислав", "Степан", "Стефан", "Тарас", "Тимофей", "Тимур", "Тит", "Трофим",
		"Феликс", "Филипп", "Фёдор", "Эдуард", "Эрик", "Юлиан", "Юлий", "Юрий", "Яков", "Ян", "Ярослав", "Милан",
	}
	firstNamesFemale = []string{
		"Александра", "Алина", "Алиса", "Алла", "Альбина", "Алёна", "Анастасия", "Анжелика", "Анна", "Антонина",
		"Анфиса", "Валентина", "Валерия", "Варвара", "Василиса", "Вера", "Вероника", "Виктория", "Владлена", "Галина",
		"Дарья", "Диана", "Дина", "Доминика", "Ева", "Евгения", "Екатерина", "Елена", "Елизавета", "Жанна", "Зинаида",
		"Злата", "Зоя", "Изабелла", "Изольда", "Инга", "Инесса", "Инна", "Ирина", "Искра", "Капитолина", "Клавдия",
		"Клара", "Клементина", "Кристина", "Ксения", "Лада", "Лариса", "Лидия", "Лилия", "Любовь", "Людмила", "Люся",
		"Майя", "Мальвина", "Маргарита", "Марина", "Мария", "Марта", "Надежда", "Наталья", "Нелли", "Ника", "Нина",
		"Нонна", "Оксана", "Олеся", "Ольга", "Полина", "Рада", "Раиса", "Регина", "Рената", "Розалина", "Светлана",
		"Софья", "София", "Таисия", "Тамара", "Татьяна", "Ульяна", "Фаина", "Федосья", "Флорентина", "Эльвира",
		"Эмилия", "Эмма", "Юлия", "Яна", "Ярослава",
	}
	middleNamesMale = []string{
		"Александрович", "Алексеевич", "Андреевич", "Дмитриевич", "Евгеньевич", "Сергеевич", "Иванович", "Фёдорович",
		"Львович", "Романович", "Владимирович", "Борисович", "Максимович",
	}
	middleNamesFemale = []string{
		"Александровна", "Алексеевна", "Андреевна", "Дмитриевна", "Евгеньевна", "Сергеевна", "Ивановна", "Фёдоровна",
		"Львовна", "Романовна", "Владимировна", "Борисовна", "Максимовна",
	}
	lastNamesCommon = []string{
		"Смирнов", "Иванов", "Кузнецов", "Соколов", "Попов", "Лебедев", "Козлов", "Новиков", "Морозов", "Петров",
		"Волков", "Соловьёв", "Васильев", "Зайцев", "Павлов", "Семёнов", "Голубев", "Виноградов", "Богданов",
		"Воробьёв", "Фёдоров", "Михайлов", "Беляев", "Тарасов", "Белов", "Комаров", "Орлов", "Киселёв", "Макаров",
		"Андреев", "Ковалёв", "Ильин", "Гусев", "Титов", "Кузьмин", "Кудрявцев", "Баранов", "Куликов", "Алексеев",
		"Степанов", "Яковлев", "Сорокин", "Сергеев", "Романов", "Захаров", "Борисов", "Королёв", "Герасимов",
		"Пономарёв", "Григорьев", "Лазарев", "Медведев", "Ершов", "Никитин", "Соболев", "Рябов", "Поляков", "Цветков",
		"Данилов", "Жуков", "Фролов", "Журавлёв", "Николаев", "Крылов", "Максимов", "Сидоров", "Осипов", "Белоусов",
		"Федотов", "Дорофеев", "Егоров", "Матвеев", "Бобров", "Дмитриев", "Калинин", "Анисимов", "Петухов", "Антонов",
		"Тимофеев", "Никифоров", "Веселов", "Филиппов", "Марков", "Большаков", "Суханов", "Миронов", "Ширяев",
		"Александров", "Коновалов", "Шестаков", "Казаков", "Ефимов", "Денисов", "Громов", "Фомин", "Давыдов",
		"Мельников", "Щербаков", "Блинов", "Колесников", "Карпов", "Афанасьев", "Власов", "Маслов", "Исаков", "Тихонов",
		"Аксёнов", "Гаврилов", "Родионов", "Котов", "Горбунов", "Кудряшов", "Быков", "Зуев", "Третьяков", "Савельев",
		"Панов", "Рыбаков", "Суворов", "Абрамов", "Воронов", "Мухин", "Архипов", "Трофимов", "Мартынов", "Емельянов",
		"Горшков", "Чернов", "Овчинников", "Селезнёв", "Панфилов", "Копылов", "Михеев", "Галкин", "Назаров", "Лобанов",
		"Лукин", "Беляков", "Потапов", "Некрасов", "Хохлов", "Жданов", "Наумов", "Шилов", "Воронцов", "Ермаков",
		"Дроздов", "Игнатьев", "Савин", "Логинов", "Сафонов", "Капустин", "Кириллов", "Моисеев", "Елисеев", "Кошелев",
		"Костин", "Горбачёв", "Орехов", "Ефремов", "Исаев", "Евдокимов", "Калашников", "Кабанов", "Носков", "Юдин",
		"Кулагин", "Лапин", "Прохоров", "Нестеров", "Харитонов", "Агафонов", "Муравьёв", "Ларионов", "Федосеев",
		"Зимин", "Пахомов", "Шубин", "Игнатов", "Филатов", "Крюков", "Рогов", "Кулаков", "Терентьев", "Молчанов",
		"Владимиров", "Артемьев", "Гурьев", "Зиновьев", "Гришин", "Кононов", "Дементьев", "Ситников", "Симонов",
		"Мишин", "Фадеев", "Комиссаров", "Мамонтов", "Носов", "Гуляев", "Шаров", "Устинов", "Вишняков", "Евсеев",
		"Лаврентьев", "Брагин", "Константинов", "Корнилов", "Авдеев", "Зыков", "Бирюков", "Шарапов", "Никонов", "Щукин",
		"Дьячков", "Одинцов", "Сазонов", "Якушев", "Красильников", "Гордеев", "Самойлов", "Князев", "Беспалов",
		"Уваров", "Шашков", "Бобылёв", "Доронин", "Белозёров", "Рожков", "Самсонов", "Мясников", "Лихачёв", "Буров",
		"Сысоев", "Фомичёв", "Русаков", "Стрелков", "Гущин", "Тетерин", "Колобов", "Субботин", "Фокин", "Блохин",
		"Селиверстов", "Пестов", "Кондратьев", "Силин", "Меркушев", "Лыткин", "Туров",
	}
	formats = map[string]string{
		FormatLastFirst: "{{.LastName}} {{.FirstName}}",
		FormatFirstLast: "{{.FirstName}} {{.LastName}}",
		FormatFIO:       "{{.LastName}} {{.FirstName}} {{.MiddleName}}",
		FormatIOF:       "{{.FirstName}} {{.MiddleName}} {{.LastName}}",
	}
)

type fullName struct {
	FirstName  string
	LastName   string
	MiddleName string
}

func (n fullName) String() string {
	return fmt.Sprintf("%s %s %s", n.LastName, n.FirstName, n.MiddleName)
}

type person struct {
	f *faker.Faker
}

// PersonRu Возвращает структуру для генерации случайных данных пользователя
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) PersonRu() *person {
	return &person{f: rf.Faker}
}

// FirstNameMale Возвращает женское имя
func (p *person) FirstNameMale() string {
	return p.f.RandomStringElement(firstNamesMale)
}

// FirstNameFemale Возвращает мужское имя
func (p *person) FirstNameFemale() string {
	return p.f.RandomStringElement(firstNamesFemale)
}

// MiddleNameMale Возвращает мужское отчество
func (p *person) MiddleNameMale() string {
	return p.f.RandomStringElement(middleNamesMale)
}

// MiddleNameFemale Возвращает женское отчество
func (p *person) MiddleNameFemale() string {
	return p.f.RandomStringElement(middleNamesFemale)
}

// LastNameMale Возвращает мужскую фамилию
func (p *person) LastNameMale() string {
	return p.f.RandomStringElement(lastNamesCommon)
}

// LastNameFemale Возвращает женскую фамилию
func (p *person) LastNameFemale() string {
	return p.f.RandomStringElement(lastNamesCommon) + femaleLastNameSuffix
}

// Gender Возвращает случайный пол
func (p *person) Gender() string {
	return p.f.RandomStringElement([]string{GenderMale, GenderFemale})
}

// FirstName Возвращает имя в зависимости от пола
func (p *person) FirstName(gender string) string {
	switch gender {
	case GenderMale:
		return p.FirstNameMale()
	case GenderFemale:
		return p.FirstNameFemale()
	default:
		return p.FirstName(p.Gender())
	}
}

// LastName Возвращает фамилию в зависимости от пола
func (p *person) LastName(gender string) string {
	switch gender {
	case GenderMale:
		return p.LastNameMale()
	case GenderFemale:
		return p.LastNameFemale()
	default:
		return p.LastName(p.Gender())
	}
}

// MiddleName Возвращает отчество в зависимости от пола
func (p *person) MiddleName(gender string) string {
	switch gender {
	case GenderMale:
		return p.MiddleNameMale()
	case GenderFemale:
		return p.MiddleNameFemale()
	default:
		return p.MiddleName(p.Gender())
	}
}

// Name Возвращает имя в соответствии с полом и форматом
func (p *person) Name(gender, format string) string {
	if gender != GenderMale && gender != GenderFemale {
		gender = p.Gender()
	}

	name := fullName{
		FirstName:  p.FirstName(gender),
		LastName:   p.LastName(gender),
		MiddleName: p.MiddleName(gender),
	}

	formatString, ok := formats[format]
	if !ok {
		formatString = p.f.RandomStringMapValue(formats)
	}

	tmpl, err := template.New("name").Parse(formatString)
	if err != nil {
		return name.String()
	}

	buff := bytes.Buffer{}

	err = tmpl.Execute(&buff, name)
	if err != nil {
		return name.String()
	}

	return buff.String()
}

func (p *person) String() string {
	return p.Name("", "")
}
