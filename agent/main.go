package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type TaskF struct {
	ID             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}
type Task struct {
	Task TaskF `json:"task"`
}

var PORT string

func getTask() (*TaskF, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%v/internal/task", PORT))
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Println("Задача не найдена.")
		} else {
			log.Printf("Unexpected status code: %d\n", resp.StatusCode)
		}
	}

	var task TaskF
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
	}

	if err = json.Unmarshal(data, &task); err != nil {
		log.Printf("Error unmarshaling JSON: %v\n", err)
	}
	return &task, nil
}

func computeTask(task *TaskF) float64 {
	var ans float64

	if task.Operation == "+" {
		timeAdd, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
		time.Sleep(time.Duration(timeAdd) * time.Millisecond)
		ans = task.Arg1 + task.Arg2
	} else if task.Operation == "-" {
		timeAdd, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
		time.Sleep(time.Duration(timeAdd) * time.Millisecond)
		ans = task.Arg1 - task.Arg2
	} else if task.Operation == "*" {
		timeAdd, _ := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
		time.Sleep(time.Duration(timeAdd) * time.Millisecond)
		ans = task.Arg1 * task.Arg2
	} else {
		timeAdd, _ := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
		time.Sleep(time.Duration(timeAdd) * time.Millisecond)
		ans = task.Arg1 / task.Arg2
	}
	return ans
}

func sendResult(res float64, id int) error {
	type jsonData struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	ans := jsonData{ID: id, Result: res}
	data, err := json.Marshal(ans)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}
	_, err = http.Post(fmt.Sprintf("http://localhost:%v/internal/task", PORT), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}
	return nil
}

func main() {
	var wg sync.WaitGroup

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	compPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Printf("COMPUTING_POWER have to be a number")
		compPower = 2
	}
	PORT = os.Getenv("PORT")
	for i := 0; i < compPower; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task, err := getTask()
				if err != nil {
					fmt.Errorf("Some trouble getting task")
				}
				if task != nil {
					log.Printf("Данные от оркестратора получены: %+v\n", task)
					res := computeTask(task)
					log.Printf("Получен результат работы операции: %+v\n", res)
					err := sendResult(res, task.ID)
					if err != nil {
						fmt.Printf("Some problem occured in sending %v", err)
					}
					log.Println("Результат операции отравлен оркестратору!")
				} else {
					log.Printf("Worker dont get any task(((")
				}
			}
		}()
	}
	wg.Wait()
}
