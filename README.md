тг юз автора: @dhysxxd
# Calculator Web Service

Calculator Web Service — это веб-сервис, который вычисляет арифметические выражения, отправленные через HTTP методом `POST`. Сервис принимает арифметическое выражение, вычисляет его результат и отправляет соответствующий HTTP-ответ.

---

## ▎Как работает

API предоставляет один эндпоинт:

### **POST /api/v1/calculate**

#### Запрос:  
Отправьте JSON с одним параметром `expression`, который содержит строку арифметического выражения.

Пример запроса:

```json
{
    "expression": "2+2*2"
}
```

#### Ответы:

1. **Успешное выполнение**  
Если вычисление прошло успешно, сервис возвращает результат с кодом **200**:

```json
{
    "result": "6"
}
```

2. **Некорректное выражение**  
Если введённое выражение содержит недопустимые символы или является ошибочным, сервис вернёт код **422 (Unprocessable Entity)** с сообщением об ошибке:

```json
{
    "error": "expression is not valid"
}
```

Примеры некорректных данных:
- Неподдерживаемые символы (буквы `a-z`, знаки $, %, и т.д.)
- Пустая строка в поле `expression`
- Неполное арифметическое выражение (например, `"2+"`)

3. **Внутренняя ошибка**  
Если произошла ошибка в процессе выполнения, сервис вернёт код **500 (Internal Server Error)**:

```json
{
    "error": "Internal server error"
}
```

---

## ▎Как запустить проект

Для запуска проекта выполните следующие шаги:

1. Убедитесь, что у вас установлен Go версии 1.22 или выше.
2. Склонируйте репозиторий:
    ```bash
    https://github.com/RaskolnikovAndTwoGrandmothers/calcu.git
    ```
    Сделайте чтобы установить все зависимости
    ```bash
    go mod tidy
    ```
3. Запустите сервис командой:
    ```bash
    go run cmd/main.go
    ```
4. По умолчанию сервис будет доступен локально на порту **8642**. URL сервиса:  
    `http://localhost:8642`

---

## ▎Примеры использования

### 1. Успешный запрос:

Отправляем математическое выражение `2+2*2`:

```bash
curl --location 'http://localhost:8642/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

Ответ:
```json
{
    "result": "6"
}
```

### 2. Ошибка 422 (некорректное выражение):

Отправляем некорректное выражение `2+abc`:

```bash
curl --location 'http://localhost:8642/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+abc"
}'
```


Ответ:
```json
{
    "error":"expression is not valid"
}
```

### 3. Ошибка 500 (внутренняя ошибка):

Отправляем некорректное выражение `5/0`(деление на ноль я добавил только чтобы продемонстрировать 500 т.к. сказали это сделать):

```bash
curl --location 'http://localhost:8642/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "5/0"
}'
```


Ответ:
```json
{
    "error":"Internal server error"
}
```
