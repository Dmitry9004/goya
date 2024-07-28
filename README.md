# goya
REST, разбивающий арифметическое выражение на постфиксное (Обратная польская запись), которое вычисляется посредством агента, принимающего элементарную задачу - выражение с 2 операндами и 1 операцией.<br />

Местом хранения выражений и задач была выбрана база данных KEY-VALUE Pudge из-за простой реализации и поддержкой многопоточности<br />
Приложение работает только с целыми числами
![изображение](https://github.com/Dmitry9004/go/assets/117633827/2012a1c8-1107-4b3a-9cf4-4b76c32c78e1)<br />

Примеры:<br />

Отправка выражения на обработку:<br />
```curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"2+2*2\"}" localhost:8080/api/v1/calculate```<br />
Ответ:<br />
```{"id": 253234}```<br />

Получениe всех выражений:<br />
```curl localhost:8080/api/v1/expressions```<br />
Ответ:<br />
```[{"Id":6813596,"Status":"Done","Result":"6"}```<br />
```{"Id":2182735,"Status":"Done","Result":"100"}]```<br />

Получение выражения по его id:<br />
```curl localhost:8080/api/v1/expressions/5433955```<br />
Ответ:<br />
```{"Id":5433955,"Status":"Done","Result":"0"}```<br />

Для запуска проекта:<br />
```cd %GOPATH% (example - "C:\Program Files\Go\src")```<br />
```git clone https://github.com/Dmitry9004/goya.git```<br />
```go run goya\project\internal\app\main.go (As admin)```<br />

Регистрация пользователя:<br />
````curl  -X POST -H "Content-Type: application/json" -d "{\"login\":\"user-test\",\"password\":\"pass-test\"}" localhost:8080/auth/register````

Аутентификация пользователя:<br />
````curl  -X POST -H "Content-Type: application/json" -d "{\"login\":\"user-test\",\"password\":\"pass-test\"}" localhost:8080/auth/login````

````curl  -X POST -H "Content-Type: application/json" -d "{\"login\":\"user-test\",\"password\":\"pass-test\"}" localhost:8080/auth/register````
