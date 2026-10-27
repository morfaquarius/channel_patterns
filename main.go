package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	ID  int
	URL string
}

type Result struct {
	Job      Job
	Status   string
	Duration time.Duration
}

func Worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	status := map[int]string{
		1: "Обработан",
		2: "Ошибка",
	}
	for job := range jobs {
		randDuration := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(randDuration) // Симуляция HTTP-запроса
		results <- Result{Job: job, Status: status[rand.Intn(2)+1], Duration: randDuration}
	}
}

func main() {
	jobs := make(chan Job, 10)
	results := make(chan Result, 10)
	var wg sync.WaitGroup
	numWorkers := 5
	// Создание воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Worker(jobs, results, &wg)
	}

	numJobs := 10
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, URL: "http://example.com/" + strconv.Itoa(j)}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	statistics := []Result{}
	for res := range results {
		statistics = append(statistics, res)
	}

	fmt.Println("Статистика обработки:")
	for _, stat := range statistics {
		fmt.Printf("ID задания: %d, URL: %s, Статус: %s, Время выполнения: %v\n", stat.Job.ID, stat.Job.URL, stat.Status, stat.Duration)
	}

	countSuccess := 0
	for _, stat := range statistics {
		if stat.Status == "Обработан" {
			countSuccess++
		}
	}

	statDuration := time.Duration(0)
	for i := 0; i < len(statistics); i++ {
		statDuration += statistics[i].Duration
	}
	statDuration = statDuration / time.Duration(len(statistics))

	fmt.Printf("Количество успешно обработанных URL: %d\n", countSuccess)
	fmt.Printf("Среднее время выполнения: %v\n", statDuration)
}
