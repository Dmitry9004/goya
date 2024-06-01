# go
REST, разбивающий арифметическое выражение на постфиксное (Обратная польская запись), которое вычисляется посредством агента, принимающего элементарную задачу - выражение с 2 операндами и 1 операцией.

Местом хранения выражений и задач была выбрана база данных KEY-VALUE Pudge из-за простой реализации и поддержкой многопоточности
![изображение](https://github.com/Dmitry9004/go/assets/117633827/2012a1c8-1107-4b3a-9cf4-4b76c32c78e1)

Примеры:

Отправка выражения на обработку:
curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"2+2*2\"}" localhost:8080/api/v1/calculate
Ответ:
{"id": 253234}

Получениe всех выражений:
curl -X POST -H "Content-Type: application/json"  localhost:8080/api/v1/expressions
Ответ:
[{"Id":6813596,"Status":"Done","Result":"6"},
{"Id":2182735,"Status":"Done","Result":"100"}]

Получение выражения по его id:
curl localhost:8080/api/v1/expressions/5433955
Ответ:
{"Id":5433955,"Status":"Done","Result":"0"}


