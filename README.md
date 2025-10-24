# ldap314ki

[![Go Report Card](https://goreportcard.com/badge/github.com/devalv/ldap314ki)](https://goreportcard.com/report/github.com/devalv/ldap314ki)
[![CodeQL](https://github.com/devalv/ldap314ki/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/devalv/ldap314ki/actions/workflows/codeql-analysis.yml)

## TODO: пример работы
TODO: заполнить

## Установка и конфигурация
TODO: заполнить


### Содержимое конфигурационного файла приложения (config.yml)
TODO: дополнить
```
debug: false
```

## Установка для разработки
1. Убедитесь, что установлена подходящая версия [Go](https://go.dev/dl/) - **1.25**.

2. Запустите **make** команду для установки утилит разработки.

```bash
make setup
```

### Make команды
- **setup**   - установка утилит для разработки/проверки
- **fmt**     - запуск gofmt и goimports
- **test**    - запуск тестов
- **build**   - сборка исполняемого файла


## Структура проекта
TODO: ещё раз подумать над структурой
```
ldap314ki/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
|   ├── app/
│       └── app.go           // Методы работы с приложением
|   ├── config/              // Хранение конфигураций для всех частей проекта
│   │   └── config.go
|   ├── transport/           // Часть на получение внутри
│   │   ├── http/
│   │   ├── grpc/
│   │   └── messaging/       // Консьюмеры
|   ├── domain/              // Обобщенные структуры / константы / ошибки
|   |   ├── models/
│   │   ├── errors/
│   │   └── consts/
|   |       └──consts.go
|   ├── usecase/             // Бизнес логика
│   │   └── ldap.go
```

## TODO v0.2
- TODO: тесты
