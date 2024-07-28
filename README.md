# goya
REST, разбивающий арифметическое выражение на постфиксное (Обратная польская запись), которое вычисляется посредством агента, принимающего элементарную задачу - выражение с 2 операндами и 1 операцией.<br />

База данныз - sqlite3.<br />
Поддержка регистрации и аутентификации пользователя.(JWT)<br />
API работает только с целыми числами.<br />
Взаимодействия между агентом и сервером осуществляется посредством gRPC.<br />
![изображение](https://github.com/Dmitry9004/go/assets/117633827/2012a1c8-1107-4b3a-9cf4-4b76c32c78e1)<br />

Примеры:<br />

Отправка выражения на обработку:<br />
```curl -X POST -H "Content-Type: application/json" -H "Authorization: this example token ..." -d "{\"expression\":\"2+2*412\"}" localhost:8080/api/v1/calculate```<br />
Ответ:<br />
```{"id": 253234}```<br />

Получениe всех выражений:<br />
```curl -H "Authorization: this example token ..."  localhost:8080/api/v1/expressions```<br />
Ответ:<br />
```[{"Id":6813596,"Status":"Done","Result":"6"}```<br />
```{"Id":2182735,"Status":"Done","Result":"100"}]```<br />

Получение выражения по его id:<br />
```curl -H "Authorization: this example token ..." localhost:8080/api/v1/expressions/5433955```<br />
Ответ:<br />
```{"Id":5433955,"Status":"Done","Result":"0"}```<br />

Для запуска проекта: (для запуска требуется установленный компилятор gcc (CGO)) <br />
```cd %GOPATH% (example - "C:\Program Files\Go\src")```<br />
```git clone https://github.com/Dmitry9004/goya.git```<br />
```go run goya\project\internal\app\main.go (As admin)```<br />

Регистрация пользователя:<br />
````curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"user-test\",\"password\":\"pass-test\"}" localhost:8080/auth/register````

Аутентификация пользователя:<br />
````curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"user-test\",\"password\":\"pass-test\"}" localhost:8080/auth/login````
Ответ: <br />
````{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjIxNzU4MDUsInVzZXJfaWQiOjZ9.UDrrQMVghpzFD-VpO1mFOrumWetmOmiEj_zLjub1NjI"}````

Запуск тестов:<br />
````cd goya\project\tests````
````go run tests -v````
