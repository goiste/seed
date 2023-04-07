package seed

import (
	"math"
	"strconv"
	"strings"
	"time"

	ru "github.com/goiste/seed/lib/randomizer/providers/ru_RU"
)

var (
	baseCommands = map[string]GeneratorFunc{
		"int":           getInt,
		"int64":         getInt64,
		"float":         getFloat,
		"word":          getWord,
		"sentence":      getSentence,
		"paragraph":     getParagraph,
		"bool":          getBool,
		"datetime":      getDatetime,
		"uuid":          getUUID,
		"name":          getName,
		"user":          getUser,
		"email":         getEmail,
		"password":      getPassword,
		"password_hash": getHash,
		"phone":         getPhone,
		"address":       getAddress,
		"region":        getRegion,
		"city":          getCity,
		"region_code":   getRegionCode,
		"lat":           getLat,
		"lon":           getLon,
		"url":           getURL,
		"company":       getCompany,
		"inn":           getINN,
		"job_title":     getJobTitle,
		"one_of":        getOneOf,
	}
)

func getInt(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		min = parseIntArgOrDefault(args, 0, math.MinInt32)
		max = parseIntArgOrDefault(args, 1, math.MaxInt32)
	)

	return rf.IntBetween(min, max), nil
}

func getInt64(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		min = parseIntArgOrDefault(args, 0, math.MinInt64)
		max = parseIntArgOrDefault(args, 1, math.MaxInt64)
	)

	return rf.IntBetween(min, max), nil
}

func getFloat(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		precision = parseIntArgOrDefault(args, 0, 2)
		min       = parseIntArgOrDefault(args, 1, math.MinInt)
		max       = parseIntArgOrDefault(args, 2, math.MaxInt)
	)

	return rf.Float(precision, min, max), nil
}

func getWord(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Lorem().Word(), nil
}

func getSentence(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		min = parseIntArgOrDefault(args, 0, 3)
		max = parseIntArgOrDefault(args, 1, 5)
	)

	wordsCount := rf.IntBetween(min, max)

	return strings.Trim(rf.Lorem().Sentence(wordsCount), "."), nil
}

func getParagraph(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		min = parseIntArgOrDefault(args, 0, 1)
		max = parseIntArgOrDefault(args, 1, 3)
	)

	sentences := rf.IntBetween(min, max)

	return rf.Lorem().Paragraph(sentences), nil
}

func getBool(rf *ru.RusFaker, args ...string) (any, error) {
	chancePercent := parseIntArgOrDefault(args, 0, 50)

	return rf.Boolean().BoolWithChance(chancePercent), nil
}

func getDatetime(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		min = parseTimeArgOrDefault(args, 0, time.Unix(0, 0))
		max = parseTimeArgOrDefault(args, 1, time.Now())
	)

	return rf.Time().TimeBetween(min, max), nil
}

func getUUID(rf *ru.RusFaker, _ ...string) (any, error) {
	return UUID4(rf.Generator), nil
}

func getName(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		gender = parseStringArgOrDefault(args, 0, "", false)
		format = parseStringArgOrDefault(args, 1, "", false)
	)

	return rf.PersonRu().Name(gender, format), nil
}

func getUser(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		name = parseStringArgOrDefault(args, 0, "", false)
	)

	return rf.InternetRu().UserName(name), nil
}

func getEmail(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		username = parseStringArgOrDefault(args, 0, "", false)
	)

	return rf.InternetRu().Email(username), nil
}

func getPassword(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Internet().Password(), nil
}

func getHash(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Hash().SHA256(), nil
}

func getPhone(rf *ru.RusFaker, args ...string) (any, error) {
	var (
		code   = parseIntArgOrDefault(args, 0, 0)
		prefix = parseIntArgOrDefault(args, 1, 0)
	)

	return rf.PhoneRu().Number(code, prefix), nil
}

func getAddress(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.AddressRu().Full(), nil
}

func getLat(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Address().Latitude(), nil
}

func getLon(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Address().Longitude(), nil
}

func getURL(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.Internet().URL(), nil
}

func getRegion(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.AddressRu().Region(true, false), nil
}

func getCity(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.AddressRu().City(false, false), nil
}

func getRegionCode(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.AddressRu().RegionCode(), nil
}

func getCompany(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.CompanyRu().Name(), nil
}

func getINN(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.CompanyRu().INN(0), nil
}

func getJobTitle(rf *ru.RusFaker, _ ...string) (any, error) {
	return rf.CompanyRu().JobTitle(), nil
}

func getOneOf(rf *ru.RusFaker, args ...string) (any, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return rf.RandomStringElement(args), nil
}

func parseIntArgOrDefault(args []string, argIndex, defaultInt int) int {
	if argIndex > len(args)-1 {
		return defaultInt
	}

	parsedInt, err := strconv.Atoi(args[argIndex])
	if err != nil {
		return defaultInt
	}

	return parsedInt
}

func parseTimeArgOrDefault(args []string, argIndex int, defaultTime time.Time) time.Time {
	if argIndex > len(args)-1 {
		return defaultTime
	}

	parsedTime, err := time.Parse("2006-01-02T15:04:05", args[argIndex])
	if err != nil {
		return defaultTime
	}

	return parsedTime
}

func parseStringArgOrDefault(args []string, argIndex int, defaultString string, notEmpty bool) string {
	if argIndex > len(args)-1 {
		return defaultString
	}

	if len(args[argIndex]) == 0 && notEmpty {
		return defaultString
	}

	return args[argIndex]
}
