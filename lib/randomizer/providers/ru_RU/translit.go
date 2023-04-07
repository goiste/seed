package ru

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	RusToASCII = map[rune]string{
		'а': "a",
		'б': "b",
		'в': "v",
		'г': "g",
		'д': "d",
		'е': "e",
		'ё': "yo",
		'ж': "zh",
		'з': "z",
		'и': "i",
		'й': "j",
		'к': "k",
		'л': "l",
		'м': "m",
		'н': "n",
		'о': "o",
		'п': "p",
		'р': "r",
		'с': "s",
		'т': "t",
		'у': "u",
		'ф': "f",
		'х': "h",
		'ц': "c",
		'ч': "ch",
		'ш': "sh",
		'щ': "sch",
		'ъ': "",
		'ы': "y",
		'ь': "",
		'э': "e",
		'ю': "ju",
		'я': "ja",
	}
	hasNoCyrillicReg = regexp.MustCompile(`^[^А-я]*$`)
)

// Transliterate Транслитерация кириллицы
func Transliterate(s string) string {
	if hasNoCyrillicReg.MatchString(s) {
		return s
	}

	runes := []rune(s)

	result := make([]rune, 0, len(runes))

	for i, r := range runes {
		asciiValue, ok := RusToASCII[unicode.ToLower(r)]
		if !ok {
			result = append(result, r)
			continue
		}

		if asciiValue == "" {
			continue
		}

		if !unicode.IsUpper(r) {
			result = append(result, []rune(asciiValue)...)
			continue
		}

		// если это последний символ в слове/строке или следующий символ тоже в верхнем регистре,
		// приводим всю ASCII-последовательность в верхний регистр ('ЩИ' => 'SCHI')
		if i+1 >= len(runes) || unicode.IsUpper(runes[i+1]) || unicode.IsSpace(runes[i+1]) {
			result = append(result, []rune(strings.ToUpper(asciiValue))...)
			continue
		}

		// если следующий символ в нижнем регистре, символы ASCII-последовательности, кроме первого,
		// должны быть в нижнем регистре ('Щи' => 'Schi')
		seq := []rune(asciiValue)
		seq[0] = unicode.ToUpper(seq[0])
		result = append(result, seq...)
	}

	return string(result)
}
