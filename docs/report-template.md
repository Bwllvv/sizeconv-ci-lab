# Отчёт: настройка GitHub Actions

Статус: ШАБЛОН. Перед экспортом в PDF заполнить ссылки, результаты и приложить
доказательства. Успешный запуск GitHub Actions пока не подтверждён.

Автор: ____________________

Дата: ____________________

## 1. Выбранный проект

Sizeconv — кроссплатформенная консольная утилита на Go для перевода размеров
данных в байты. Например, `sizeconv "1.5 GiB"` выводит `1610612736 bytes`.
Используется библиотека `github.com/dustin/go-humanize v1.0.1`.
Проект содержит собственный код CLI и автоматические тесты корректных и
ошибочных входных данных. Использование кода примера из модуля необходимо
исключить сравнением с материалами курса.

Открытый репозиторий: [ВСТАВИТЬ ФАКТИЧЕСКУЮ ССЫЛКУ]

## 2. Код ci-cd.yml

```yaml
name: Build and test

on:
  push:
  pull_request:
  workflow_dispatch:

permissions:
  contents: read

jobs:
  build:
    name: Build (${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        include:
          - os: ubuntu-latest
            binary: sizeconv
          - os: windows-latest
            binary: sizeconv.exe
          - os: macos-latest
            binary: sizeconv
    steps:
      - name: Checkout sources
        uses: actions/checkout@v7

      - name: Set up Go and restore cache
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache: true
          cache-dependency-path: go.sum

      - name: Install dependencies
        run: go mod download

      - name: Verify dependencies
        run: go mod verify

      - name: Static analysis
        run: go vet ./...

      - name: Run tests
        run: go test -v -count=1 ./...

      - name: Build executable
        run: go build -trimpath -o dist/${{ matrix.binary }} .

      - name: Smoke test
        run: ./dist/${{ matrix.binary }} "1.5 GiB"

      - name: Upload executable
        uses: actions/upload-artifact@v7
        with:
          name: sizeconv-${{ matrix.os }}
          path: dist/${{ matrix.binary }}
          if-no-files-found: error
          retention-days: 14
```

## 3. Описание джобы

Джоба `build` выполняется по матрице ОС: Ubuntu, Windows и macOS на размещённых
GitHub runners. События запуска: push в любую ветку, pull_request
(по умолчанию opened, synchronize, reopened), ручной workflow_dispatch.
Каждый вариант матрицы выполняет одинаковую последовательность:

1. `actions/checkout@v7` получает исходники.
2. `actions/setup-go@v7` устанавливает Go по go.mod и восстанавливает кэш.
3. `go mod download` устанавливает зависимости.
4. `go mod verify` проверяет целостность загруженных модулей.
5. `go vet ./...` выполняет статический анализ.
6. `go test -v -count=1 ./...` запускает тесты конвертации и ошибок ввода.
7. `go build -trimpath` собирает бинарный файл в dist.
8. Пробный запуск проверяет работу собранного файла.
9. `actions/upload-artifact@v7` сохраняет бинарник на 14 дней.

Для Windows используется расширение .exe. Имена артефактов различаются по ОС.
Если бинарник отсутствует, загрузка артефакта завершается ошибкой. При падении
одного варианта остальные продолжают работу: fail-fast отключён.
Права GITHUB_TOKEN ограничены чтением содержимого репозитория.

## 4. Кэширование и артефакты

Кэш модулей Go и результатов компиляции включён через cache: true
в setup-go. Файл go.sum задаёт хеш зависимостей в ключе кэша.
При первом запуске кэш создаётся, при повторном с тем же ключом
восстанавливается. Неизменившиеся модули повторно скачивать не требуется.
Для разных ОС используются отдельные кэши.

Артефакты можно скачать на странице запуска Actions в разделе Artifacts:
sizeconv-ubuntu-latest, sizeconv-windows-latest, sizeconv-macos-latest.
Кэш ускоряет следующие сборки, а артефакт предназначен для получения
готовой программы пользователем.

## 5. Фактическая проверка

Первый запуск: [ССЫЛКА, ДАТА, РЕЗУЛЬТАТ ВСЕХ ТРЁХ ДЖОБ]

Повторный запуск: [ССЫЛКА, ДАТА, РЕЗУЛЬТАТ ВСЕХ ТРЁХ ДЖОБ]

Восстановление кэша: [ТОЧНАЯ СТРОКА ИЗ ЛОГА ПОВТОРНОГО ЗАПУСКА]

Артефакты: [ИМЕНА И РЕЗУЛЬТАТ ПРОВЕРКИ СКАЧИВАНИЯ]

Приложить скриншоты успешного запуска, трёх артефактов и восстановления кэша.

## 6. Результат

[ЗАПОЛНИТЬ ПО ФАКТУ: настроен ли успешный pipeline, подтверждены ли
тесты/сборка на трёх ОС, кэширование, доступность артефактов.]
