Программа для вычисления арифметических выражений, написана на языке Golang. Работает следующим образом: 
запускается сервер(далее есть описание по запуску), и принимает запросы по адресу http://localhost:8080/api/v1/calculate. 
Запрос отправляется в виде json-файла: {"expression: "...."}, вместо "...." нужно подставить выражение, например "(2+2)*3-7", 
и соответственно получает ответ в виде json-файла: {"result": "....."}, "....." результат вычисления выражения, для выражения выше получим "5.0".

В программе есть свои тесты для проверки работы по вычислению выражений и http-тесты, для проверки http-запросов.

** Клонирование репозитория!

git clone https://github.com/enigmasterr/calchttp.git - склонировать репозиторий

** Запуск сервера!

зайти в каталог calchttp - cd calchttp

go run cmd/main.go - запуск проекта, по умолчанию порт 8080, либо можно получить через переменную окружения

** Проверка программы на тестах!

запускаем тесты (тесты запускались в windows 11). Необходимо зайти в командную строку и отправить тесты ниже:

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*9\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+2*9\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"yy2+2*9\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*9\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"/\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"1/2\"}" http://localhost:8080/api/v1/calculate

curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+3)*9\"}" http://localhost:8080/api/v1/calculate

** Программа логирует запросы и ответы в консоль!

** в папке pkg/calculation/ добавлены http тесты их всего 3

** Запуск тестов

Запустить все тесты в том числе и http, нужно зайти в папку pkg/calculation/ и ввести go test -v

чтобы зайти в папку нужно набрать cd pkg/calculation/
