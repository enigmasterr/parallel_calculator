package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if err != nil {
		return nil, err
	}
	var task TaskF
	err = json.NewDecoder(resp.Body).Decode(&task)
	fmt.Println(task)
	if err != nil {
		return nil, err
	}
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
	data, err := json.Marshal(&ans)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}
	_, err = http.Post("http://localhost:8080/submitResult", "application/json", bytes.NewBuffer(data))
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
				if err != nil {
					fmt.Errorf("Some trouble getting task")
				}
				if task != nil {
					res := computeTask(task)
					err := sendResult(res, task.ID)
					if err != nil {
						fmt.Printf("Some problem occured in sending %v", err)
					}
				} else {
					log.Printf("Worker dont get any task(((")
				}
				time.Sleep(2 * time.Second)
			}
		}()
	}
	wg.Wait()
}
