package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"os"
	"strconv"
	"todo-app/dto"
	"todo-app/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	Client       *dynamodb.Client
	Table        string
	CounterTable string
}

var (
	db DynamoDBClient
)

func init() {
	utils.LoadEnv()
	db = *NewDynamoDBClient()
}

func GetDB() *DynamoDBClient {
	return &db
}

func NewDynamoDBClient() *DynamoDBClient {
	region := os.Getenv("AWS_REGION")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	counterTable := "counter"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load AWS configuration: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBClient{
		Client:       client,
		Table:        tableName,
		CounterTable: counterTable,
	}
}

func (db *DynamoDBClient) GetNextID(ctx context.Context) (string, error) {
	counterName := "TaskID"

	result, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &db.CounterTable,
		Key: map[string]types.AttributeValue{
			"CounterName": &types.AttributeValueMemberS{Value: counterName},
		},
		UpdateExpression: aws.String("SET CurrentValue = CurrentValue + :incr"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return "", err
	}

	newID := result.Attributes["CurrentValue"].(*types.AttributeValueMemberN).Value
	return newID, nil
}

func (db *DynamoDBClient) CreateTask(ctx context.Context, task dto.Task) (dto.Task, error) {

	taskID, err := db.GetNextID(ctx)
	if err != nil {
		return dto.Task{}, err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.Table,
		Item: map[string]types.AttributeValue{
			"Id":          &types.AttributeValueMemberN{Value: taskID},
			"Title":       &types.AttributeValueMemberS{Value: task.Title},
			"Description": &types.AttributeValueMemberS{Value: task.Description},
		},
	})

	if err != nil {
		return dto.Task{}, err
	}

	uintValue, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		return dto.Task{}, err
	}

	result := dto.Task{
		Id:          uint(uintValue),
		Title:       task.Title,
		Description: task.Description,
	}

	return result, nil
}

func (db *DynamoDBClient) GetTask(ctx context.Context, id string) (*dto.Task, error) {
	res, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &db.Table,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberN{Value: id},
		},
	})
	if err != nil || res.Item == nil {
		return &dto.Task{
			Id:          0,
			Description: "Not found",
			Title:       "Not found",
		}, err
	}

	title := res.Item["Title"]
	description := res.Item["Description"]
	newID := res.Item["Id"].(*types.AttributeValueMemberN).Value
	uintValue, err := strconv.ParseUint(newID, 10, 32)
	if err != nil {
		return &dto.Task{
			Id:          0,
			Description: "Not found",
			Title:       "Not found",
		}, err
	}
	task := dto.Task{
		Title:       title.(*types.AttributeValueMemberS).Value,
		Description: description.(*types.AttributeValueMemberS).Value,
		Id:          uint(uintValue),
	}
	return &task, err
}

func (db *DynamoDBClient) GetAllTasks(ctx context.Context) ([]dto.Task, error) {
	result, err := db.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &db.Table,
	})
	if err != nil {
		return nil, err
	}

	var tasks []dto.Task
	for _, item := range result.Items {

		title := item["Title"]
		description := item["Description"]
		id := item["Id"].(*types.AttributeValueMemberN).Value
		uintValue, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return tasks, err
		}
		task := dto.Task{
			Title:       title.(*types.AttributeValueMemberS).Value,
			Description: description.(*types.AttributeValueMemberS).Value,
			Id:          uint(uintValue),
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
