# AuthCheck
Часть сервиса аутентификации пользователей, обеспечивающую безопасную и масштабируемую авторизацию через REST API.

## Содержание
- [Технологии](#технологии)
- [Начало работы](#установка)
- [Использование](#использование)
- [Тестирование](#тестирование)
- [Contributing](#contributing)
- [FAQ](#FAQ)
- [Команда проекта](#команда-проекта)

## Технологии
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT](https://pkg.go.dev/github.com/dgrijalva/jwt-go)

## Установка
Для работы с проектом требуется только **Docker**. Его наличине **обязательно**!

Клонирование репозитория:
```sh
$ git clone https://github.com/MuhaFAH/AuthCheck.git
```

Перейдите в директорию проекта (**AuthCheck/**) и установите зависимости:
```sh
$ go mod tidy
```
Создайте env-файл из шаблона:
```sh
$ cp .env.example .env
```
Из директории проекта перейдите в папку configs
```sh
cd configs
```
Запустите сборку приложения следующей командой:
```sh
sudo docker compose --env-file ../.env up --build
```
В конечном итоге в консоли должен быть примерно такой вывод:
```sh
 ✔ Container postgres    Created                                                                                                                                                 0.0s 
 ✔ Container migrations  Created                                                                                                                                                 0.0s 
 ✔ Container app         Recreated   
postgres    | 2024-08-25 11:08:16.980 UTC [1] LOG:  database system is ready to accept connections
migrations exited with code 0
```
Это означает, что все контейнеры успешно запущены, сервер поднят.
## Использование
### Добавление пользователя
Добавление пользователя в базу данных и получение токенов происходит путем отправки запроса по адресу:
```
http://localhost:8080/auth/{GUID}
```
Где вместо **{GUID}** вводится GUID пользователя. При верном запросе пользователя, в ответном запросе сервера будут токены.<br><br>
Пример отправки запроса и получения ответа сервера:
```sh
curl http://localhost:8080/auth/3F2504E0-4F89-11D3-9A0C-0305E82C3474

{"answer":"OK","reason":"","status":200,"tokens":{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1ODg1MjIsImlwIjoiMTcyLjE4LjAuMSJ9.yLw5H5R4n7DAmGydFKSd56tnr3C3Kzf5Nrw5-YUJGFgLN9WEsCJqKsEUoeAbFWFze7AY9xITIODU_R8xjziFBA","refresh_token":"unujY5xAdM7LbNBRmNV9/ebp3nH/09ZFoKUFqRMEIV8="}}
```
Пример неправильного запроса (неверный формат GUID):
```sh
curl http://localhost:8080/auth/3F2504E0-4F89-11D3-9A0C

{"answer":"access denied","reason":"invalid user guid","status":403,"tokens":{"access_token":"","refresh_token":""}}
```
### Замена токенов
Замена токенов существующего пользователя происходит путем отправки запроса по адресу:
```
http://localhost:8080/refresh/{GUID}
```
Где вместо **{GUID}** вводится GUID пользователя. В теле запроса должен быть указан актуальный Refresh-токен пользователя. При верном запросе пользователя, в ответном запросе сервера будут токены.<br><br>
Пример отправки запроса и получения ответа сервера:
```sh
curl -d '{"refresh_token":"1DAC+2m1DsfYw1PRq0a508hNTXzUXN0o5eLZF1x7JiI="}' http://localhost:8080/refresh/3F2504E0-4F89-11D3-9A0C-0305E82C3490

{"answer":"OK","reason":"","status":200,"tokens":{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ1ODkzNjIsImlwIjoiMTcyLjE4LjAuMSJ9.JWrnPf3uJ6P-10I9GN5rCi6huSD2_knuNjBW6biTInKGGUZuOPqKhAduOCtFRAAm-Po1PO06r5XFdOwQ6SQSDQ","refresh_token":"+UnkyGYwHcCN85G5HbKnvHcD5lMLxFDLhgoC4nAckKs="}}
```
>Если с момента прошлого запроса поменялся IP-адрес клиента, будет отправлено предупреждающее письмо на его почту. В данном случае используются моки, все письма появляются в папке **email/**.

Пример неправильного запроса (неверный Refresh-токен или GUID):
```sh
curl -d '{"refresh_token":"1DAC+2m1DsfYw1PRq0a508h"}' http://localhost:8080/refresh/3F2504E0-4F89-11D3-9A0C-0305E82C3490

{"answer":"access denied","reason":"invalid refresh token","status":403,"tokens":{"access_token":"","refresh_token":""}}
```

## Тестирование
На данный момент тестирование отсутствует.

## Contributing
По всем вопросам, отсутствующим в FAQ, обращаться в telegram: @ChanChinCho

## FAQ
### Как изменить порт, данные БД и так далее?
Все изменения вносятся в файл .env, расположенный в корневой директории проекта.

### Находятся ли мои токены в базе данных?
Все токены шифруются, а в базе данных находится только дополнительно хэшированный Refresh-токен. Он нужен для валидации токена пользователя.

### Почему нет тестов?
Они будут добавлены, но не скоро.

## Команда проекта

- [Муха](tg://resolve?domain=chanchincho) — делал всё.
