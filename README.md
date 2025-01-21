# wb-current-weather
[![Go Report Card](https://goreportcard.com/badge/github.com/devalv/wb-current-weather)](https://goreportcard.com/report/github.com/devalv/wb-current-weather)
[![CodeQL](https://github.com/devalv/wb-current-weather/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/devalv/wb-current-weather/actions/workflows/github-code-scanning/codeql)
[![CodeQL](https://github.com/devalv/wb-inbox-mail-count/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/devalv/wb-inbox-mail-count/actions/workflows/github-code-scanning/codeql)

## Отображение текущей погоды
![пример](example.jpg)

## Установка и конфигурация
Результатом работы будет отображение текущего прогноза погоды для выбранного города.
Данные будут получены от [OpenWeather](https://openweathermap.org/current#cityid), с помощью API v2.5
(https://api.openweathermap.org/data/2.5/weather?id={city id}&appid={API key}&units={units}&lang={lang})
Перед началом работы Вам необходимо найти ID города [тут](https://bulk.openweathermap.org/sample/) и получить собственный API-ключ для обращения к API.

### Установка собранного bin-файла
1. Загрузите соответствующую версию из раздела [релизы](https://github.com/devalv/wb-current-weather/releases)
2. Скопируйте исполняемый файл в /usr/local/bin (или иной каталог доступный waybar на запуск)
3. Создайте файл-конфигурации по инструкции описанной ниже
4. Проверьте запуск командой `wbcw -config /home/user/.config/wb-current-weather/config.yml`
5. Если на 4м шаге произошли ошибки - активируйте ключ debug в config.yml и повторите запуск
6. Добавьте отображение статуса в waybar (инструкция ниже)

### Содержимое конфигурационного файла приложения (config.yml)
```
debug: false
city_id: 498817
weather_api_token: "you-api-key"
units: "metric"
lang: "ru"
```

#### Параметры по-умолчанию:
    units - metric
    lang  - ru

### Добавление запуска в waybar (~/.config/waybar/config.jsonc)
1. Добавьте отображение вывода в раздел **modules-right** (или иной)
```json
"modules-right": [
    ...
    "battery",
    "custom/wbcw",
    ...
],
```
2. Добавьте обработчик вывода
```json
    ...
   "custom/wbcw": {
     "exec" : "wbcw -config /home/user/.config/wb-current-weather/config.yml",
    "return-type": "json",
    "interval": 300,
     "format": "{}"
    },
    "battery": {
        "format": "{icon} {capacity}%",
        "format-icons": ["", "", "", "", ""]
    },
    ...
```

### Настройка отступов для waybar (~/.config/waybar/style.css)
```css
#custom-wbcw {
    color: @text;
    padding-right: 13px;
 }
```

## Установка для разработки
1. Убедитесь, что установлена подходящая версия [Go](https://go.dev/dl/) - **1.23**.

2. Запустите **make** команду для установки утилит разработки.

```bash
make setup
```

### Make команды
- **setup**   - установка утилит для разработки/проверки
- **fmt**     - запуск gofmt и goimports
- **test**    - запуск тестов
- **cover**   - вывод % покрытия тестов
- **build**   - сборка исполняемого файла


## Структура проекта
```
wb-current-weather/
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
|   |   ├── forecast.go
│   │   └── waybar.go

```

<!-- ## Сборка deb-пакета -->
<!-- TODO: актуализировать для v0.2 -->

## TODO v0.2
- TODO: автоматизировать сборку deb-пакета в github
- TODO: автоматизировать сборку bin-артефактов в github
- TODO: тесты
- TODO: сборка debian-пакета
