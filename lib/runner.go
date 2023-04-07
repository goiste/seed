package lib

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	pg "github.com/goiste/seed/lib/db"
	"github.com/goiste/seed/lib/query"
	seed "github.com/goiste/seed/lib/randomizer"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate go install github.com/valyala/quicktemplate/qtc@v1.7.0
//go:generate qtc -dir=query

const (
	nullableFieldPrefix = "?"
)

var percentageReg = regexp.MustCompile(`^\d{0,2}`)

// Command Команда SQL
type Command struct {
	Cmd  string
	Args []any
}

type Runner struct {
	db        *pgxpool.Pool
	rand      *seed.Randomizer
	batchSize int
}

func NewRunner(db *pgxpool.Pool, randSeed int64, batchSize int) *Runner {
	return &Runner{
		db:        db,
		rand:      seed.NewRandomizer(randSeed),
		batchSize: batchSize,
	}
}

type randomizeResponse struct {
	valuesMap map[string]any
	err       error
	last      bool
}

// RegisterRandomizerCommands Добавляет новые команды для генератора псевдослучайных значений
func (r *Runner) RegisterRandomizerCommands(commands map[string]seed.GeneratorFunc) {
	r.rand.RegisterCommands(commands)
}

// RunSQL Выполняет готовый SQL-скрипт
func (r *Runner) RunSQL(ctx context.Context, sql string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql)
	if err != nil {
		//goland:noinspection GoUnhandledErrorResult
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// RunSchema Выполняет команды схемы и всех её таблиц
func (r *Runner) RunSchema(ctx context.Context, schema Schema) error {
	commands := make([]Command, 0, 10)

	// устанавливаем сессионные переменные
	for _, envVar := range schema.EnvVars {
		commands = append(commands, Command{
			Cmd:  "select set_config($1, $2, true)",
			Args: []interface{}{envVar.Key, envVar.Value},
		})
	}

	// первичные ключи, возвращённые при вставке данных в таблицы
	returnedPrimaryKeys := make(map[string][]any, len(schema.Tables))

	for _, table := range schema.Tables {
		if len(table.Values) == 0 {
			continue
		}

		returnedKeys, err := r.runTable(ctx, table, schema.Name, commands, returnedPrimaryKeys)
		if err != nil {
			return err
		}

		returnedPrimaryKeys[table.Name] = returnedKeys
	}

	return nil
}

// Выполняет команду вставки значений в таблицу
func (r *Runner) runTable(ctx context.Context, table Table, schemaName string, schemaCommands []Command, tablesPrimaryKeys map[string][]any) ([]any, error) {
	values := r.randomizeValues(table.Columns, table.Values)

	clearColumns(table.Columns)

	valuesToWrite := make([]map[string]any, 0, r.batchSize)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback(ctx)

	returnedValues := make([]any, 0, 1000)

	total := 0

	for value := range values {
		if value.err != nil {
			return nil, fmt.Errorf("randomize value error: %w", value.err)
		}

		if !value.last {
			valuesToWrite = append(valuesToWrite, value.valuesMap)

			if len(valuesToWrite) < r.batchSize {
				continue
			}
		}

		{
			tableValues, err := r.fillForeignKeys(valuesToWrite, table.ForeignKeys, tablesPrimaryKeys)
			if err != nil {
				return nil, err
			}

			cmd := Command{
				Cmd:  query.InsertData(schemaName, table.Name, table.PrimaryKey, table.Columns, len(tableValues)),
				Args: flattenValues(tableValues, table.Columns),
			}

			var returned []any

			if table.PrimaryKey != "" {
				returned = make([]any, 0, len(tableValues))
			}

			// установка сессионных переменных
			for _, schemaCommand := range schemaCommands {
				if _, err = tx.Exec(ctx, schemaCommand.Cmd, schemaCommand.Args...); err != nil {
					return nil, err
				}
			}

			if returned != nil {
				var rows pgx.Rows

				rows, err = tx.Query(ctx, cmd.Cmd, cmd.Args...)
				if err == nil {
					err = pg.ScanReturnedIds(&returned, rows)
				}
			} else {
				_, err = tx.Exec(ctx, cmd.Cmd, cmd.Args...)
			}

			if err != nil {
				return nil, err
			}

			returnedValues = append(returnedValues, returned...)

			total += len(valuesToWrite)

			go fmt.Print("\rprocessed: ", total)
		}

		valuesToWrite = valuesToWrite[:0]
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	return returnedValues, tx.Commit(ctx)
}

// Подставляет внешние ключи в значения
func (r *Runner) fillForeignKeys(tableValues []map[string]any, tableForeignKeys []ForeignKey, tablesPrimaryKeys map[string][]any) ([]map[string]any, error) {
	// внешние ключи не используются, возвращаем неизменные данные
	if len(tableForeignKeys) == 0 {
		return tableValues, nil
	}

	// мапа вида {колонка_с_внешним_ключом: слайс_ключей_внешней_таблицы}
	fkMap := make(map[string][]any, len(tableForeignKeys))
	for _, foreignKey := range tableForeignKeys {
		fkMap[foreignKey.Column] = tablesPrimaryKeys[foreignKey.RefTable]
	}

	for i, value := range tableValues {
		for k, v := range value {
			// проверяем, надо ли подставлять внешний ключ в это поле
			ids, ok := fkMap[k]
			if !ok {
				continue
			}

			// если нужен внешний ключ, значение должно быть индексом массива
			fkIdx, err := strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				return nil, fmt.Errorf("foreign key must be an index of reference table values, got %#v", v)
			}

			if fkIdx > len(ids)-1 {
				fkIdx = rand.Intn(len(ids))
			}

			// подставляем ключ внешней таблицы по индексу
			value[k] = ids[fkIdx]
		}
		tableValues[i] = value
	}

	return tableValues, nil
}

// Подставляет случайные значения в поля таблицы
func (r *Runner) randomizeValues(columns []string, values []map[string]any) <-chan randomizeResponse {
	ch := make(chan randomizeResponse, r.batchSize)

	cols := make([]string, len(columns))
	copy(cols, columns)

	go func() {
		defer close(ch)

		for _, valueMap := range values {
			var repeat bool
			var repeatCount int

			// если задан ключ для повтора, считываем количество повторений
			if rep, ok := valueMap["$repeat"]; ok {
				repeat = true

				cnt, err := strconv.Atoi(fmt.Sprintf("%v", rep))
				if err != nil {
					ch <- randomizeResponse{
						err: fmt.Errorf("cannot parse repeat count: %w", err),
					}
					return
				}

				repeatCount = cnt
			}

			// если повтор не требуется, заполняем случайными данными одно значение
			if !repeat {
				for field := range valueMap {
					err := r.randomizeValue(field, valueMap)
					if err != nil {
						ch <- randomizeResponse{
							err: err,
						}
						return
					}
				}
				ch <- randomizeResponse{
					valuesMap: valueMap,
				}
				continue
			}

			if repeatCount == 0 {
				continue
			}

			// создаём repeatCount клонов переданного значения, заполняем случайными данными, добавляем в результат
			for i := 0; i < repeatCount; i++ {
				tmp := make(map[string]any, len(valueMap))
				for k, v := range valueMap {
					tmp[k] = v
				}

				// сохраняем порядок полей для повторяемости псевдослучайных значений
				for _, field := range cols {
					err := r.randomizeValue(field, tmp)
					if err != nil {
						ch <- randomizeResponse{
							err: err,
						}
						return
					}
				}

				ch <- randomizeResponse{
					valuesMap: tmp,
				}
			}
		}

		ch <- randomizeResponse{
			last: true,
		}
	}()

	return ch
}

func (r *Runner) randomizeValue(field string, valueMap map[string]any) error {
	var (
		value      = valueMap[field]
		validField = field
		needToFill bool
	)

	if strings.HasPrefix(field, nullableFieldPrefix) {
		validField, needToFill = r.checkNullable(field[1:])

		value = valueMap[validField]
		delete(valueMap, field)

		if !needToFill {
			valueMap[validField] = nil
			return nil
		}
	}

	strVal, ok := value.(string)
	if !ok {
		return nil
	}

	cmd, ok := seed.ParseCommand(strVal)
	if !ok {
		return nil
	}

	val, err := r.rand.RunCommand(cmd)
	if err != nil {
		return err
	}

	valueMap[validField] = val

	return nil
}

func (r *Runner) checkNullable(field string) (string, bool) {
	number, clearField := percentageReg.FindString(field), percentageReg.ReplaceAllString(field, "")
	percentage, _ := strconv.Atoi(number)
	if percentage < 1 || percentage > 99 {
		percentage = 50
	}
	return clearField, r.rand.Bool(percentage)
}

// Перекладывает данные в плоскую структуру для использования в качестве аргументов prepared запроса
func flattenValues(mapValues []map[string]any, columns []string) (flatValues []any) {
	// список полей с их порядковыми номерами
	colMap := make(map[string]int, len(columns))
	for i, column := range columns {
		colMap[column] = i
	}

	flatValues = make([]any, 0, len(mapValues)*len(columns))

	for _, mapValue := range mapValues {
		row := make([]any, len(columns))

		for col, val := range mapValue {
			if idx, ok := colMap[col]; ok {
				row[idx] = val
			}
		}

		flatValues = append(flatValues, row...)
	}

	return
}

func clearColumns(columns []string) {
	for i := range columns {
		if strings.HasPrefix(columns[i], nullableFieldPrefix) {
			columns[i] = percentageReg.ReplaceAllString(columns[i][1:], "")
		}
	}
}
