package ru

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jaswdr/faker"
)

var (
	freeMailDomains = []string{
		"yandex.ru", "ya.ru", "narod.ru", "gmail.com", "mail.ru", "list.ru", "bk.ru", "inbox.ru", "rambler.ru",
		"hotmail.com", "vk.com",
	}
	invalidEmailCharactersReg = regexp.MustCompile(`[^a-z0-9._%+\-]+`)
)

type internet struct {
	f      *faker.Faker
	person *person
}

// InternetRu Возвращает структуру для генерации случайных адресов электронной почты на бесплатных почтовых доменах
//
//goland:noinspection GoExportedFuncWithUnexportedType
func (rf *RusFaker) InternetRu() *internet {
	return &internet{f: rf.Faker, person: rf.PersonRu()}
}

// FreeMailDomain Возвращает случайный бесплатный почтовый домен
func (i *internet) FreeMailDomain() string {
	return i.f.RandomStringElement(freeMailDomains)
}

// Email Возвращает случайный адрес электронной почты по имени пользователя или со случайным именем
func (i *internet) Email(name string) string {
	if len(name) < 1 {
		name = i.person.Name("", "fl")
	}

	return fmt.Sprintf("%s@%s", i.UserName(name), i.FreeMailDomain())
}

// UserName Возвращает имя пользователя
func (i *internet) UserName(name string) string {
	if len(name) < 1 {
		name = i.person.Name("", "fl")
	}

	name = Transliterate(strings.ToLower(name))

	nameParts := strings.Split(name, " ")

	username := nameParts[len(nameParts)-1]
	for _, n := range nameParts[:len(nameParts)-1] {
		if len(n) < 1 {
			continue
		}
		username = fmt.Sprintf("%s.%s", string(n[0]), username)
	}

	return invalidEmailCharactersReg.ReplaceAllString(username, "_")
}
