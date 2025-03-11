package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
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

func getTask() (*TaskF, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	fmt.Println("Trying to get task")
	fmt.Println("Agent, getTask ", resp)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Println("Задача не найдена.")
		} else {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}
	}

	var task TaskF
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if err = json.Unmarshal(data, &task); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	fmt.Println("Agent getting task from orchestrator and got that --- ", task)
	return &task, nil
}

func computeTask(task *TaskF) float64 {
	var ans float64
	if task.Operation == "+" {
		ans = task.Arg1 + task.Arg2
	} else if task.Operation == "-" {
		ans = task.Arg1 - task.Arg2
	} else if task.Operation == "*" {
		ans = task.Arg1 * task.Arg2
	} else {
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
	_, err = http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send result: %v", err)
	}
	return nil
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				task, err := getTask()
				fmt.Println("Answer taken --- ", task)
				if err != nil {
					fmt.Errorf("Some trouble getting task")
				}
				if task != nil {
					res := computeTask(task)
					fmt.Println("Result of operation ---- ", res)
					err := sendResult(res, task.ID)
					if err != nil {
						fmt.Printf("Some problem occured in sending %v", err)
					}
				} else {
					log.Printf("Worker dont get any task(((")
				}
				time.Sleep(3 * time.Second)
			}
		}()
	}
	wg.Wait()
}
