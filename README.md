# Seed

Библиотека для заполнения БД PostgreSQL начальными данными.

### Использование в качестве библиотеки
```go
import seed "github.com/goiste/seed/lib"
```
### Использование в качестве консольной утилиты
```shell
go install github.com/goiste/seed@latest
```

### Описание входных данных
На вход принимает данные (`[]byte`) в одном из форматов:
- готовый SQL скрипт;
- файл `json`/`yaml` с описанием схем, таблиц и данных.

Возможно добавление парсера для собственного формата входных данных. Собственный парсер должен реализовать метод `Parse(data []byte) ([]Schema, error)`

### Описание структуры файла на примере `yaml`:
```yaml
- schema: schema_test # схема БД
  env_vars: # сессионные переменные
    - key: mod.user
      value: 'seed'
  tables:
    - name: main # имя таблицы
      # первичный ключ (чтобы можно было получить из БД первичные ключи при создании 
      # и подставить их во внешние ключи других таблиц)
      primary_key: id
      columns: # поля таблицы
        - full_name
        - username
        - email
        - ?hired_at # необязательное поле
        - is_active
      values: # массив данных в формате "поле: значение"
        - username: admin
          email: $email:admin # команды для генерации случайных значений, описание ниже
          is_active: true
        - $repeat: 27 # вставить значение со сгенерированными полями 27 раз
          full_name: $name:,fio
          username: $one_of:user1,user2,user3 # случайное значение из указанных
          email: $email
          hired_at: $datetime:2012-12-12T12:00:00,2022-12-12T12:00:00 # timestamp в указанном диапазоне
          is_active: $bool:75
    - name: second
      primary_key: id
      foreign_keys: # внешние ключи
        - column: user_id # поле текущей таблицы, содержащее внешний ключ
          ref_table: main # внешняя таблица
      columns:
        - user_id
        - uuid
        - title
        - ?40description # необязательное поле с вероятностью заполнения 40%
        - address
        - ?75phone # необязательное поле с вероятностью заполнения 75%
        - amount
        - values
      values:
        - $repeat: 42
          # Индекс элемента в списке первичных ключей, полученных при создании внешней таблицы. 
          # В данном примере будет сгенерирован случайный индекс из диапазона 1-27 (включительно): 
          # по числу сгенерированных во внешней таблице данных, кроме первого (admin)
          user_id: $int:1,27
          uuid: $uuid
          title: $sentence:2,5
          description: $paragraph:2,4
          address: $address
          phone: $phone
          amount: $int:111,222
          values: $[,5]float:3,1,3 # массив float длиной от 1 до 5
```

Внешние ключи указываются в виде индекса элемента в списке первичных ключей, полученных при создании внешней таблицы, поэтому таблицы в файле должны идти в правильном порядке: сначала внешняя, потом ссылающиеся на неё.

Примеры файлов можно посмотреть в директории [example](./example)

### Генерация псевдослучайных значений

Данные можно генерировать псевдослучайным образом. Для этого используются специальные команды, начинающиеся с `$`, за которыми через двоеточие идут аргументы, например: `$int:1,100` — целое число в диапазоне от 1 до 100 включительно.

Массив значений обозначается квадратными скобками перед командой: `$[от,до]команда:аргументы`, например: `$[,3]int:1,7` — массив целых чисел от 1 до 7; длина массива от 1 до 3 элементов

### Перечень предопределённых команд:

| Команда        | Описание                                              | Перечень аргументов                                                                                                                                                                                    | Примеры                                           |
|----------------|-------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------|
| $int           | Целое 32-разрядное число                              | Минимальное значение, максимальное значение                                                                                                                                                            | $int<br/>$int:1,100<br/>$int:,42<br/>$int:3       |
| $int64         | Целое 64-разрядное число                              | См. $int                                                                                                                                                                                               | См. $int                                          |
| $float         | Число с плавающей точкой                              | Знаков после запятой (по умолч. 2), минимальное значение, максимальное значение                                                                                                                        | $float<br/>$float:3<br/>$float:,1,10<br/>$float,3 |
| $word          | 1 слово из набора Lorem ipsum                         | Не имеет аргументов                                                                                                                                                                                    | $word                                             |
| $sentence      | Несколько слов из набора Lorem ipsum                  | Минимальное кол-во слов (по умолч. 3), максимальное кол-во (по умолч. 5)                                                                                                                               | $sentence<br/>$sentence:,10<br/>$sentence:2,2     |
| $paragraph     | Несколько предложений из слов из набора Lorem ipsum   | Минимальное кол-во предложений (по умолч. 1), максимальное кол-во (по умолч. 3)                                                                                                                        | $paragraph<br/>$paragraph:,5<br/>$paragraph:2,2   |
| $bool          | Булево значение                                       | Вероятность выпадения true в процентах (по умолч. 50)                                                                                                                                                  | $bool<br/>$bool:75                                |
| $datetime      | Дата и время                                          | Минимальная дата (по умолч. 1970-01-01), максимальная дата (по умолч. "сейчас")<br/>Формат: "YYYY-MM-DDThh:mm:ss"                                                                                      | $datetime<br/>$datetime:,2012-12-12T00:00:00      |
| $uuid          | UUID                                                  | Не имеет аргументов                                                                                                                                                                                    | $uuid                                             |
| $name          | Полное имя в заданном формате [RU]                    | Пол (по умолч. случайный), формат (по умолч. fio)<br/>Варианты пола: m (муж), f (жен)<br/>Варианты формата: fl (Имя Фамилия), lf (Фамилия Имя), fio (Фамилия Имя Отчество), iof (Имя Отчество Фамилия) | $name<br/>$name:f<br/>$name:,iof                  |
| $user          | Имя пользователя                                      | Имя, на основании которого будет сформировано имя пользователя (по умолч. сгенерированные Имя Фамилия, будет преобразовано в i.familija)                                                               | $user                                             |
| $email         | Адрес электронной почты на бесплатном почтовом домене | Имя пользователя для почтового ящика (по умолч. см. $user)                                                                                                                                             | $email<br/>$email:admin                           |
| $password      | Пароль                                                | Не имеет аргументов                                                                                                                                                                                    | $password                                         |
| $password_hash | Хеш SHA-256                                           | Не имеет аргументов                                                                                                                                                                                    | $password_hash                                    |
| $phone         | Номер телефона в международном формате                | Код (по умолч. от 900 до 999), префикс (по умолч. 7)<br/>Если префикс 8, код всегда будет 800                                                                                                          | $phone<br/>$phone:812<br/>$phone:,375             |
| $address       | Почтовый адрес [RU]                                   | Не имеет аргументов                                                                                                                                                                                    | $address                                          |
| $region        | Регион [RU]                                           | Не имеет аргументов                                                                                                                                                                                    | $region                                           |
| $city          | Город [RU]                                            | Не имеет аргументов                                                                                                                                                                                    | $city                                             |
| $region_code   | Код региона [RU]                                      | Не имеет аргументов                                                                                                                                                                                    | $region_code                                      |
| $lat           | Широта                                                | Не имеет аргументов                                                                                                                                                                                    | $lat                                              |
| $lon           | Долгота                                               | Не имеет аргументов                                                                                                                                                                                    | $lon                                              |
| $url           | URL-адрес                                             | Не имеет аргументов                                                                                                                                                                                    | $url                                              |
| $company       | Название организации [RU]                             | Не имеет аргументов                                                                                                                                                                                    | $company                                          |
| $inn           | ИНН юридического лица [RU]                            | Не имеет аргументов                                                                                                                                                                                    | $inn                                              |
| $job_title     | Название должности [RU]                               | Не имеет аргументов                                                                                                                                                                                    | $job_title                                        |
| $one_of        | Случайный элемент из списка                           | Список элементов через запятую                                                                                                                                                                         | $one_of:elem1,elem2                               |

Также предусмотрен механизм добавления собственных команд. Для этого необходимо вызвать метод `RegisterRandomizerCommands()`, передав команды в виде словаря, ключами которого являются названия команд (без префикса `$`), а значениями — функции со следующей сигнатурой: `func(fake *ru.RusFaker, args ...string) (any, error)` Например:
```go
import (
    ...
    seed "github.com/goiste/seed/lib"
    rnd "github.com/goiste/seed/lib/randomizer"
    ru "github.com/goiste/seed/lib/randomizer/providers/ru_RU"
)

...

seeder, _ := seed.NewSeeder(seed.Config{
    File: file,
    DSN:  dsn,
})

seeder.RegisterRandomizerCommands(map[string]rnd.GeneratorFunc{
	"big_lorem": func(fake *ru.RusFaker, args ...string) (any, error) {
		return fake.Lorem().Text(1800), nil
    }   
})

seeder.Seed(ctx)
...
```
Название команды должно содержать только символы `A-Za-z0-9._`, регистр учитывается.
