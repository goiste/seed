package seed

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	ru "github.com/goiste/seed/lib/randomizer/providers/ru_RU"
	"github.com/jaswdr/faker"
)

const (
	commandPrefix = "$"
	argsSeparator = ","

	defaultArrayFrom = 1
	defaultArrayTo   = 3
	maxArrayTo       = 1000
)

var (
	commandReg = regexp.MustCompile(`^\$(\[(\d+)?,?(\d+)?])?([\w.]+)(:(.*))?$`)
)

type Command struct {
	Command   string
	Args      []string
	ArrayFrom int
	ArrayTo   int
	IsArray   bool
}

type GeneratorFunc func(fake *ru.RusFaker, args ...string) (any, error)

type Randomizer struct {
	rusFaker *ru.RusFaker
	commands map[string]GeneratorFunc
}

func NewRandomizer(randSeed int64) *Randomizer {
	f := faker.New()

	if randSeed > 0 {
		f = faker.NewWithSeed(rand.NewSource(randSeed))
	}

	return &Randomizer{
		rusFaker: ru.NewFaker(&f),
		commands: baseCommands,
	}
}

func (r *Randomizer) Bool(percentage int) bool {
	return r.rusFaker.BoolWithChance(percentage)
}

func (r *Randomizer) RegisterCommands(commands map[string]GeneratorFunc) {
	for command, generatorFunc := range commands {
		r.commands[command] = generatorFunc
	}
}

func (r *Randomizer) RunCommand(cmd Command) (any, error) {
	cmdFunc, ok := r.commands[cmd.Command]
	if !ok {
		return nil, fmt.Errorf("unknown command: %q", cmd.Command)
	}

	if !cmd.IsArray {
		return cmdFunc(r.rusFaker, cmd.Args...)
	}

	cmd.IsArray = false

	count := r.getArrayLen(cmd.ArrayFrom, cmd.ArrayTo)

	result := make([]any, count)

	var err error
	for i := 0; i < count; i++ {
		result[i], err = r.RunCommand(cmd)
		if err != nil {
			return nil, err
		}
	}

	bld := strings.Builder{}
	bld.WriteString("{")
	for i, val := range result {
		bld.WriteString(fmt.Sprintf("%v", val))
		if i < len(result)-1 {
			bld.WriteString(",")
		}
	}
	bld.WriteString("}")

	return bld.String(), nil
}

// ParseCommand Извлекает из строки команду и аргументы
func ParseCommand(str string) (cmd Command, ok bool) {
	if !strings.HasPrefix(str, commandPrefix) {
		return
	}

	matches := commandReg.FindStringSubmatch(str)
	if len(matches) != 7 {
		return
	}

	ok = true

	cmd.Command = matches[4]

	args := strings.Split(matches[6], argsSeparator)
	if len(args) > 0 && (len(args[0]) > 0 || len(args) > 1) { // первый аргумент может быть пустым: command:,arg2
		cmd.Args = args
	}

	if matches[1] == "" {
		return
	}

	arrayFrom, _ := strconv.Atoi(matches[2])
	arrayTo, _ := strconv.Atoi(matches[3])

	if arrayTo >= arrayFrom || arrayTo == 0 {
		cmd.IsArray = true
		cmd.ArrayFrom = arrayFrom
		cmd.ArrayTo = arrayTo
	}

	return
}

func (r *Randomizer) getArrayLen(from int, to int) int {
	if from < 1 && to < 1 {
		from, to = defaultArrayFrom, defaultArrayTo
	}

	if from < 1 {
		from = 1
	}

	if to < from {
		to = maxArrayTo
	}

	return r.rusFaker.IntBetween(from, to)
}
