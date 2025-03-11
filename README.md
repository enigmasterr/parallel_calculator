** Описание работы программы!

Программа для распределенного вычисления арифметических выражений, написана на языке Golang. Работает следующим образом: 
запускается оркестратор (сервер) (далее есть описание по запуску), а также запускается агент(далее также будет описание запуска) и принимает запросы по адресу http://localhost:8080/api/v1/calculate. 
Запрос отправляется в виде json-файла: {"expression: "...."}, вместо "...." нужно подставить выражение, например "(2+2)*3-7", ответом на запрос приходит json-файл {"id": ...}, с номером выражения. Номер назначается автоматически.
Запущенный агент при этом делает запрос на вычисление. Как только выражение будет разложено, то в словарь добавляется задание и агент его получает, затем запускает вычисление и отпраляет json-файл с результатом вычисления. Ответ добавляется в map с результатами. После попадает
в дальнейшую работу по разбору выражения.

Для получения результатов всех выражений надо сделать запрос оркестратору http://localhost:8080/api/v1/expressions.

Для получения результата одного выражения по id, нужно сделать запрос через командную строку (cmd) http://localhost:8080/api/v1/expressions/:id, вместо id ставится значение выражения, которое получаете при вводе выражения.

В программе есть свои тесты для проверки работы по вычислению выражений и http-тесты, для проверки http-запросов.

** Клонирование репозитория!

git clone https://github.com/enigmasterr/parallel_calculator.git - склонировать репозиторий

** Запуск оркестратора!

зайти в каталог parallel_calculator - cd parallel_calculator

go run cmd/main.go - запуск оркестратора, по умолчанию порт 8080, либо можно получить через переменную окружения

** Запуск агента! Запустить новую командную строку (cmd)

зайти в каталог parallel_calculator - cd parallel_calculator

go run agent/main.go - запуск агента, по умолчанию порт 8080, либо можно получить через переменную окружения

** Проверка программы на тестах!

запускаем тесты (тесты запускались в windows 11). Необходимо открыть еще одну командную строку (cmd) и отправить тесты ниже:

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
