# Sizeconv: GitHub Actions lab

Учебная кроссплатформенная утилита на Go. Переводит размер данных в байты:

```sh
go run . "1.5 GiB"
# 1610612736 bytes
go run . "1 GB"
# 1000000000 bytes
```

Зависимость: `github.com/dustin/go-humanize v1.0.1`. Различает десятичные
единицы (GB = 10^9) и двоичные (GiB = 2^30).

## Локальная проверка

Нужен Go >= 1.24.

```sh
go mod download
go mod verify
go vet ./...
go test -v -count=1 ./...
go build -trimpath -o dist/sizeconv .
./dist/sizeconv "1.5 GiB"
```

На Windows имя исполняемого файла: `sizeconv.exe`.

## Workflow

Файл: `.github/workflows/ci-cd.yml`.

События: каждый `push`, создание/обновление/повторное открытие
`pull_request`, ручной `workflow_dispatch`.

Джоба `build` запускается в трёх вариантах матрицы: `ubuntu-latest`,
`windows-latest`, `macos-latest`. Каждый вариант получает исходники,
устанавливает Go, восстанавливает кэш, устанавливает и проверяет зависимости,
выполняет статический анализ, тесты, сборку и пробный запуск программы.

Использованные Actions:

| Action | Назначение |
| --- | --- |
| `actions/checkout@v7` | Получение исходников репозитория |
| `actions/setup-go@v7` | Установка Go из go.mod; кэш модулей и результатов компиляции |
| `actions/upload-artifact@v7` | Загрузка бинарника как скачиваемого артефакта на 14 дней |

Кэш включён через `cache: true`, файл зависимостей для ключа — `go.sum`.
setup-go разделяет кэш по ОС, архитектуре, версии Go и хешу go.sum.
При подходящем кэше модули восстанавливаются, поэтому `go mod download`
не скачивает их повторно. При изменении зависимостей или отсутствии/удалении
кэша скачивание потребуется снова. Тесты принудительно выполняются повторно
благодаря `-count=1`.

Артефакты: `sizeconv-ubuntu-latest`, `sizeconv-windows-latest`,
`sizeconv-macos-latest`. Откройте успешный запуск во вкладке Actions,
прокрутите до Artifacts и скачайте нужный архив. Для скачивания может
потребоваться вход в GitHub. После распаковки на Linux/macOS при необходимости
выполните `chmod +x sizeconv`.

Отдельного развёртывания на сервер нет: детальные требования задания
предусматривают сборку, тесты, кэш и сохранение артефакта.

## Публикация с macOS

```sh
brew install gh
gh auth login --hostname github.com --git-protocol https --web
gh auth setup-git
gh api user --jq .login
```

Убедитесь, что показан нужный аккаунт. Затем из каталога проекта:

```sh
git init -b main
git add .
git commit -m "Add size converter and cross-platform CI"
gh repo create sizeconv-ci-lab --public --source=. --remote=origin --push
gh browse --actions
```

Если Git запросит имя или email автора, задайте свои значения через
`git config user.name` и `git config user.email`, затем повторите commit.
Если репозиторий с таким именем существует, выберите другое имя в команде
`gh repo create`.

## Проверка GitHub Actions

1. Дождитесь зелёного первого запуска и всех трёх вариантов `build`.
2. На странице запуска проверьте три скачиваемых артефакта.
3. Запустите workflow ещё раз: Actions → Build and test → Run workflow → main.
4. Дождитесь зелёного второго запуска. В `Set up Go and restore cache`
   найдите сообщение об успешном восстановлении кэша. Не ориентируйтесь
   только на скорость: сеть и загрузка runner могут влиять на время.
5. Сохраните ссылки на оба запуска и скриншоты: зелёные джобы,
   артефакты, восстановление кэша.

При ошибке откройте красный шаг и сохраните текст ошибки.

## Отчёт

Шаблон: `docs/report-template.md`. Его нужно дополнить ссылками и результатами
реальных запусков, затем экспортировать в PDF. Локальные результаты проверки
сохранены отдельно в `docs/local-verification.txt`; они не заменяют запуск CI.

Документация Actions:

- https://github.com/actions/checkout
- https://github.com/actions/setup-go#caching-dependency-files-and-build-outputs
- https://github.com/actions/upload-artifact
