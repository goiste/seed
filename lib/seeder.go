package lib

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pg "github.com/goiste/seed/lib/db"
	seed "github.com/goiste/seed/lib/randomizer"
	"github.com/rs/zerolog/log"
)

// Вариант заполнения БД
//
// @see TypeFile, TypeSQL
type typeOfSeeding string

const (
	// TypeSQL Готовый SQL-скрипт
	TypeSQL typeOfSeeding = "sql"
	// TypeFile Файл с данными для вставки
	TypeFile typeOfSeeding = "file"

	lockFileName   = "seed.lock"
	lockTimeFormat = time.RFC3339

	minBatchSize = 100
	maxBatchSize = 25000
)

type DataParser interface {
	Parse(data []byte) ([]Schema, error)
}

type Config struct {
	File       string
	DSN        string
	SeedType   typeOfSeeding
	DataParser DataParser
	RandSeed   int64
	BatchSize  int
}

type Seeder struct {
	seedType   typeOfSeeding
	data       []byte
	workingDir string
	parser     DataParser
	runner     *Runner
	batchSize  int
}

// NewSeeder Создаёт новый экземпляр seeder'а
func NewSeeder(ctx context.Context, config Config) (*Seeder, error) {
	if config.File == "" {
		return nil, fmt.Errorf("input file is required")
	}

	if config.DSN == "" {
		return nil, fmt.Errorf("DB connection string is required")
	}

	data, err := os.ReadFile(config.File)
	if err != nil {
		return nil, fmt.Errorf("cannot read seed file: %w", err)
	}

	db, err := pg.Connect(ctx, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	fileExt := strings.ToLower(
		strings.Trim(filepath.Ext(config.File), "."),
	)

	seedType := getSeedType(config.SeedType, fileExt)
	dataParser := getDataParser(config.DataParser, fileExt)

	if dataParser == nil && seedType != TypeSQL {
		return nil, fmt.Errorf("there is no predefined data parser for filetype %q", fileExt)
	}

	batchSize := config.BatchSize
	if batchSize < minBatchSize || batchSize > maxBatchSize {
		batchSize = maxBatchSize
	}

	return &Seeder{
		seedType:   seedType,
		data:       data,
		workingDir: filepath.Dir(config.File),
		parser:     dataParser,
		runner:     NewRunner(db, config.RandSeed, batchSize),
	}, nil
}

// RegisterRandomizerCommands Добавляет новые команды для генератора псевдослучайных значений
func (s *Seeder) RegisterRandomizerCommands(commands map[string]seed.GeneratorFunc) {
	s.runner.RegisterRandomizerCommands(commands)
}

// Seed Запускает заполнение БД
func (s *Seeder) Seed(ctx context.Context) error {
	lockTime, err := s.LockTime()
	if err != nil {
		return fmt.Errorf("cannot read lock time: %w", err)
	}

	// запускаем seed только один раз
	if lockTime != nil {
		log.Debug().Msgf("seed was already been launched at %s, skip", lockTime.Format(lockTimeFormat))
		return nil
	}

	if s.seedType == TypeSQL {
		err = s.runner.RunSQL(ctx, string(s.data))
		if err != nil {
			return err
		}

		return s.Lock()
	}

	schemas, err := s.parser.Parse(s.data)
	if err != nil {
		return fmt.Errorf("parse data error: %w", err)
	}

	for _, schema := range schemas {
		err = s.runner.RunSchema(ctx, schema)
		if err != nil {
			return fmt.Errorf("run schema %q error: %w", schema.Name, err)
		}
	}

	return s.Lock()
}

// LockTime Возвращает время последнего запуска seed
func (s *Seeder) LockTime() (*time.Time, error) {
	fileName := filepath.Join(s.workingDir, lockFileName)

	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("error reading file %q: %w", fileName, err)
	}

	tm, err := time.Parse(lockTimeFormat, string(data))
	if err != nil {
		return nil, fmt.Errorf("error parsing time from string %q: %w", string(data), err)
	}

	return &tm, nil
}

// Lock Записывает текущее время в качестве времени последнего запуска seed во избежание повторного запуска
func (s *Seeder) Lock() error {
	fileName := filepath.Join(s.workingDir, lockFileName)

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %q: %w", fileName, err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	_, err = f.WriteString(time.Now().Format(lockTimeFormat))
	if err != nil {
		return fmt.Errorf("error writing file %q: %w", fileName, err)
	}

	return nil
}

// Unlock Удаляет lock-файл (например, если надо запустить seed с новыми данными)
func (s *Seeder) Unlock() error {
	fileName := filepath.Join(s.workingDir, lockFileName)

	if err := os.Remove(fileName); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing file %q: %w", fileName, err)
	}

	return nil
}

//goland:noinspection GoExportedFuncWithUnexportedType
func TypeFromString(s string) (typeOfSeeding, bool) {
	switch s {
	case string(TypeFile):
		return TypeFile, true
	case string(TypeSQL):
		return TypeSQL, true
	}

	return "", false
}

// Если seedType установлен, используем его, если нет, определяем по расширению
func getSeedType(seedType typeOfSeeding, fileExt string) typeOfSeeding {
	if seedType != "" {
		return seedType
	}

	if fileExt == "sql" {
		return TypeSQL
	}

	return TypeFile
}

// Если dataParser установлен, используем его, если нет, определяем по расширению
func getDataParser(dataParser DataParser, fileExt string) DataParser {
	if dataParser != nil {
		return dataParser
	}

	switch fileExt {
	case "json":
		return new(JSONParser)
	case "yml", "yaml":
		return new(YAMLParser)
	default:
		return nil
	}
}
