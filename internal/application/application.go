package application

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/enigmasterr/calchttp/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	log.Println("get request - ", request)
	if err != nil {
		type ErrStr struct {
			Error string `json:"error"`
		}
		ansJson := ErrStr{Error: "Request is not valid"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ansJson)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	//fmt.Println(result)
	if err != nil {
		// "error": "Expression is not valid"
		type ErrStr struct {
			Error string `json:"error"`
		}
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, calculation.ErrInvalidExpression) {
			ansJson := ErrStr{Error: "Expression is not valid"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ansJson)
			log.Printf("err: %s", err.Error())
		} else if errors.Is(err, calculation.ErrStrangeSymbols) {
			ansJson := ErrStr{Error: "Expression is not valid"}
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ansJson)
			log.Printf("err: %s", err.Error())
		} else {
			ansJson := ErrStr{Error: "Internal server error"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ansJson)
			log.Printf("err: %s", err.Error())
		}

	} else {
		type ResStr struct {
			Result string `json:"result"`
		}
		convRes := fmt.Sprintf("%f", result)
		ansJson := ResStr{Result: string(convRes)}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ansJson)
		log.Printf("send json {\"result\": \"%s\"}", string(convRes))
		// fmt.Fprintf(w, "result: %f", result)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
