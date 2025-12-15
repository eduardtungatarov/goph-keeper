[![Go](https://img.shields.io/badge/Go-1.25%2B-blue?logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Required-blue?logo=docker)](https://docker.com)

###Запуск сервера:

`make up` - поднимаем зависимые контейнеры (нужен docker)

`make server` - запускаем сервер (нужен go 1.25+)

###Запуск клиента:

Выбери свою платформу

`./bin/keeper-darwin-amd64`

`./bin/keeper-darwin-arm64`

`./bin/keeper-linux-amd64`

`./bin/keeper-windows-amd64`

###Команды клиента:

#### Аутентификация:

Регистрация

`./keeper register user2 -p secret123`

Вход

`./keeper login user1 -p secret123`

#### Работа с данными:

Список всех данных

`./keeper list`

Сохранить пароль

`./keeper create pwd`  

Сохранить карту

`./keeper create card`

Сохранить бинарник

`./keeper create binary -d "data..." -T "config.txt"`

Прочитать данные с ID=123

`./keeper read 123`

Удалить данные с ID=123

`./keeper delete 123`
