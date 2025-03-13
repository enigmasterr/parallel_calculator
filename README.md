# Описание работы программы!

Программа для распределенного вычисления арифметических выражений, написана на языке Golang. Работает следующим образом: 
запускается оркестратор (сервер) (далее есть описание по запуску), который принимает запросы по адресу ```http://localhost:8080/api/v1/calculate.```  
Для распределенной работы запускается агент, который запускает несколько воркеров и опрашивает сервер на наличие задач. Как только появляется задача, агент обрабатывает ее и возвращает ответ обратно в виде json-файла.  


Запрос на сервер отправляется в виде json-файла: ```{"expression: "...."}```, вместо "...." нужно подставить выражение, например "(2+2)*3-7", ответом на запрос приходит json-файл ```{"id": ...}```, с номером выражения. Номер назначается автоматически.  
Запущенный агент при этом делает запрос на вычисление. Как только выражение будет разложено, то в словарь добавляется задание и агент его получает, затем запускает вычисление и отпраляет json-файл с результатом вычисления. Ответ добавляется в map с результатами. После попадает
в дальнейшую работу по разбору выражения.  


#### Запросы к серверу!

Для получения результатов всех выражений надо сделать запрос оркестратору 
```
http://localhost:8080/api/v1/expressions.
```  
Для получения результата одного выражения по ```id```, нужно сделать запрос через командную строку (cmd) 
```
http://localhost:8080/api/v1/expressions/:id
```
вместо ```id``` ставится значение выражения, которое получаете при вводе выражения.  


В программе есть свои тесты для проверки работы по вычислению выражений и http-тесты, для проверки http-запросов.

#### Клонирование репозитория!

```git clone https://github.com/enigmasterr/parallel_calculator.git``` - склонировать репозиторий

#### Запуск оркестратора!

зайти в каталог parallel_calculator - ```cd parallel_calculator```

```go run cmd/main.go``` - запуск оркестратора, по умолчанию порт 8080, либо можно получить через переменную окружения

#### Запуск агента! Запустить новую командную строку (cmd)

зайти в каталог parallel_calculator - ```cd parallel_calculator```

```go run agent/main.go``` - запуск агента, по умолчанию порт 8080, либо можно получить через переменную окружения


## Схема работы приложения!

Все результаты выражения хранятся в срезе Allexpressions - ```[{id1, status1, result1}, {id2, status2, result2}, ...]```  
Все задания агенту находятся в map allTasks ```[id1: {id1, arg1, arg2, op1, op_time1}, id2: ...., ]```  
Все результаты находятся в map allresults ```[{id1: res1}, {id2: res2}, {id2: res2}, ... ]```  


!Нужно быть в папке с клонированным проектом!

cmd1 -- ```go run cmd/main.go```   - поднимаем сервер(оркестратор)  
cmd2 -- ```go run agent/main.go``` - запускаем агента  
cmd3 -- пишем запросы на вычисление выражения ```curl ...```  



```                                    |---------------|    "internal/task", POST            |---------|```  
```           "api/v1/calculate"       |               |   <---------------------------      |         |```  
```cmd3 ---- отправка выражения -----> |  Оркестратор  |        {id, result}                 |  Агент  |```  
```           {"expression":"..."}     |               |                                     |         |```  
```                                    |               |      "internal/task", GET           |         |```  
```cmd3  <---------------------------  |---------------|   ------------------------------>   |---------|```  
```                {id: ...}                             Task{id, arg1, arg2, oper, op_time}            ```  



```         "api/v1/expressions"         |-------------|```  
```cmd3 <------------------------------> | Оркестратор |```  
```      [{id1, status1, result1}, ...]  |             |```  
```          "api/v1/expressions/:id"    |             |```  
```cmd3 <------------------------------> |-------------|```  
```           {id, status, result}                      ```  


```                    Оркестратор                     ```  
```|--------------------------------------------------|```  
```|              AllExpressions                      |```  
```|                   ^                              |```  
```|                   |                              |```  
```|                   |                              |```  
```|                   |       "getresult/id"         |```  
```|              Allresults --------------> Calc()   |```  
```|                           {id: result}           |```  
```|                                                  |```  
```|                                                  |```  
```|              AllTasks   <------------- Calc()    |```  
```|------------------------------------------------- |```  
 

## Проверка программы на тестах!

запускаем тесты (тесты запускались в windows 11). Необходимо открыть еще одну командную строку (cmd) и отправить тесты ниже:

1. Корректные запросы 
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*9\"}" http://localhost:8080/api/v1/calculate
```
Ответ: 200, ```{"result": 18}```

```  
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"1/2\"}" http://localhost:8080/api/v1/calculate  
```  
Ответ: 200, ```{"result": 0.5}```  

```  
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+3)*9\"}" http://localhost:8080/api/v1/calculate  
```  
Ответ: 200, ```{"result": 36}```  


2. Некорректные запросы
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+2*9\"}" http://localhost:8080/api/v1/calculate
```

Ответ: 400, ```{"result": 0}```

```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"yy2+2*9\"}" http://localhost:8080/api/v1/calculate
```

Ответ: 422, ```{"result": 0}```

```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"/\"}" http://localhost:8080/api/v1/calculate
```
Ответ: 400, ```{"result": 0}```
```
curl -X GET http://localhost:8080/api/v1/expressions/:999
```
Ответ: 404, ```{"result": 0}```

** Программа логирует запросы и ответы в консоль!

** в папке pkg/calculation/ добавлены http тесты их всего 3

## Запуск тестов

Запустить все тесты в том числе и http, нужно зайти в папку pkg/calculation/ и ввести go test -v

чтобы зайти в папку нужно набрать cd pkg/calculation/
