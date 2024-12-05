package service

import (
	"context"
	"sort"
	"strconv"
	"todo-app/database"
	"todo-app/dto"
	"todo-app/utils/log"
)

var (
	taskCache     = make(map[string]dto.Task)
	fetchRequest  = make(chan string, 10)
	fetchResponse = make(chan dto.Task, 10)
)

func WorkerRoutine() {
	for taskID := range fetchRequest {
		log.Info("Worker fetching task:", taskID)
		ctx := context.Background()
		task, err := database.GetDB().GetTask(ctx, taskID)
		if err != nil {
			log.Error("Error fetching task:", err)
			fetchResponse <- dto.Task{
				Id:          0,
				Title:       "Not found",
				Description: "Not found",
			}
			continue
		}
		taskCache[taskID] = *task
		log.Info("Task fetched by worker:", task)
		fetchResponse <- *task
	}
}

func GetAllTasks() (*[]dto.Task, error) {

	dbClient := database.GetDB()
	ctx := context.Background()
	tasks, err := dbClient.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		taskCache[strconv.Itoa(int(task.Id))] = task
		log.Info("Task added to cache:", task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})

	return &tasks, nil

}

func GetTask(id string) (*dto.Task, error) {

	task, found := taskCache[id]

	if found {
		log.Info("Task found in cache:", id)
		return &task, nil
	}

	log.Info("Task not in cache. Requesting worker:", id)
	fetchRequest <- id

	fetchedTask := <-fetchResponse
	log.Info("Task fetched from db:", fetchedTask)

	return &fetchedTask, nil

}

func AddTask(taskReq *dto.Task) (*dto.Task, error) {

	dbClient := database.GetDB()
	ctx := context.Background()
	task, err := dbClient.CreateTask(ctx, *taskReq)

	if err != nil {
		return nil, err
	}

	taskCache[strconv.Itoa(int(task.Id))] = task
	log.Info("Task added to cache:", task)

	return &task, nil

}
