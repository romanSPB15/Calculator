﻿## Веб-сервис для вычисления арифметических выражений на Go

### Описание

Программа представляет собой веб-сервис, который вычисляет арифметические выражения. Пользователь может использовать как целые, так и вещественные числа (разделитель - точка или запятая). Калькулятор поддерживает базовые операции (/, *, -, +`), а также скобки для обозначения приоритета операции.

Пользователь отправляет арифметическое выражение по HTTP и получает в ответ его результат. У сервиса URL [/api/v1/calculate](http://127.0.0.1:8080/api/v1/calculate). Пользователь отправляет на этот URL POST-запрос с телом:
```bash
{
    "expression": "выражение, которое ввёл пользователь"
}
```
В ответ пользователь получает HTTP-ответ с телом:
```bash
{
    "result": "результат выражения"
}
```
и кодом `200`, если выражение вычислено успешно, либо HTTP-ответ с телом:
```bash
{
    "error": "Expression is not valid"
}
```
и кодом `422`, если входные данные не соответствуют требованиям приложения — например, кроме цифр и разрешённых операций пользователь ввёл символ английского алфавита.

Ещё один вариант HTTP-ответа:
```bash
{
    "error": "Internal server error"
}
```
и код `500` в случае какой-либо иной ошибки («Что-то пошло не так»).

### Установка
Для установки программы сначала необходимо выполнить следущюю команду:
```bash
git clone https://github.com/romanSPB15/Calculator
```

### Запуск программы
Для запуска программы необходимо перейти в директорию с проектом и ввести команду:
```bash
go run ./cmd/main.go
```

### Работа с сервисом

Существует множество вариантов для работы с сервисом.

Если речь идет о работе с сервисом через **Postman**, то необходимо зайти на страницу с отправкой запросов и выбрать тип запроса **POST**. В поле с вводом **URL** необходимо прописать:
```bash
http://127.0.0.1:80/api/v1/calculate
```
и для создания тела запроса следует перейти в **Body** -> **raw**, и в поле для написания тела запроса написать данные в *JSON формате*. Например:
```bash
{
    "expression": "2+2"
}
```
После этого можно отправлять запрос на сервер.

Можно воспользоваться и другим вариантом - командной строкой. Если у вас *Windows*, то рекомендую использовать *Git Bash*. Через терминал мы можем реализовывать запросы через `curl`. Например:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

### Примеры работы с сервисом

Если пользователь введет корректный запрос с корректно написанным арифметическим выражением, то он получит верный ответ.
При вводе данного запроса:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "2+2*2"
}'
```
вы получите ответ:
```bash
{"result":"6"}
```
с кодом `200`

Далее рассмотрим пример с использованием скобок и отрицательных чисел. При вводе запроса:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "-2*(3-1)+8"
}'
```
вы получите ответ:
```bash
{"result":"4"}
```
с кодом 200

При делении на 0 пользователю выводится ошибка. При вводе запроса:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "4/0"
}'
```
пользователь получает код `500`.

При вводе запроса:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "as"
}'
```
пользователь получит ответ с ошибкой:
```bash
{"error":"Expression is not valid"}
```
и код 422

При нарушении самого запроса также будет выведена ошибка. При вводе запроса:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "as
}'
```
пользователю выводит:
```bash
{"error": "Internal server error"}
```
и код 500

Если же запрос не POST, а какой-либо другой, то опять же выводится ошибка. Например, при вводе запроса:
```bash
curl --location --request GET 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "as"
}'
```
пользователю выводит:
```bash
{"error": "Internal server error"}
```
и код `500`

Следует так же уточнить, что на пустое арифметическое выражение программа так же будет отвечать ошибкой. Например, на запрос:
```bash
curl --location 'localhost:80/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": ""
}'
```
пользователь увидит:
```bash
{"error":"Expression is not valid"}
```
и код 422.
