package coredb

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/jsii-runtime-go"
)

type DB struct{
	Client *dynamodb.DynamoDB
}

// CreateDynamoDBClient creates a new DynamoDB client and session
func NewDb(s *session.Session) (*DB, error) {

	// Create a DynamoDB client
	db := &DB{
		Client: dynamodb.New(s),
	}
	return db, nil
}

func (c DB) QueryDatabasePartition(table_name string, partition_key string) (*dynamodb.QueryOutput, error) {
	indexName := "date_createdIndex"
	input := &dynamodb.QueryInput{
		TableName: jsii.String(table_name),
		KeyConditions: map[string]*dynamodb.Condition{
			"partition_key": {
				ComparisonOperator: jsii.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: jsii.String(partition_key),
					},
				},
			},
		},
		IndexName: jsii.String(indexName),
		ScanIndexForward: jsii.Bool(false),
	}

	result, err := c.Client.Query(input)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (c DB) QueryDatabase(table_name string, partition_key string, sort_key *string) (*dynamodb.GetItemOutput, error) {
	// Interface to query the database
	input := &dynamodb.GetItemInput{
		TableName: jsii.String(table_name),
		Key: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: jsii.String(partition_key),
			},
		},
	}

	if sort_key != nil {
		input.Key["sort_key"] = &dynamodb.AttributeValue{
			S: jsii.String(*sort_key),
		}
	}

	result, err := c.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}


func (c DB) MutateDatabase(table_name string, partition_key string, sort_key string, attributes interface{}) error {

	// Convert the 'attributes' interface to a map of attribute values
	attributeMap, err := dynamodbattribute.MarshalMap(attributes)

	if err != nil {
		return err
	}

	// Interface to add a new record to the database
	input := &dynamodb.PutItemInput{
		TableName: jsii.String(table_name),
		Item: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: jsii.String(partition_key),
			},
			"sort_key": {
				S: jsii.String(sort_key),
			},
		},
	}

	for attributeName, attributeValue := range attributeMap {
		input.Item[attributeName] = attributeValue
	}

	// Update the database
	_, err = c.Client.PutItem(input)
	
	if err != nil {
		return err
	}
	return nil
}

func (c DB) CacheData(table_name string, partition_key string, sort_key string, attributes interface{}, ttl int64) error {
	// Convert the 'attributes' interface to a map of attribute values
	attributeMap, err := dynamodbattribute.MarshalMap(attributes)
	if err != nil {
		return err
	}

	// Interface to add a new record to the database
	input := &dynamodb.PutItemInput{
		TableName: jsii.String(table_name),
		Item: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: jsii.String(partition_key),
			},
			"sort_key": {
				S: jsii.String(sort_key),
			},
		},
	}

	for attributeName, attributeValue := range attributeMap {
		input.Item[attributeName] = attributeValue
	}

	ttlAttributeValue := &dynamodb.AttributeValue{
		N: jsii.String(strconv.FormatInt(ttl, 10)),
	}
	input.Item["ttl"] = ttlAttributeValue

	// Update the database
	_, err = c.Client.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}


func (c DB) DeleteRecord(table_name string, partition_key string, sort_key *string) error {
	// Interface to delete a record by partition key and optionally a sort key
	input := &dynamodb.DeleteItemInput{
		TableName: jsii.String(table_name),
		Key: map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: jsii.String(partition_key),
			},
		},
	}

	if sort_key != nil {
		input.Key["sort_key"] = &dynamodb.AttributeValue{
			S: jsii.String(*sort_key),
		}
	}

	_, err := c.Client.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (c DB) UpdateRecord(table_name string, partition_key string, sort_key string, attributes interface{}) error {
	
	// Convert the interface to a map[string]*dynamodb.AttributeValue
	attrs, err := dynamodbattribute.MarshalMap(attributes)
	if err != nil {
		return errors.New("failed to convert attributes to DynamoDB format - " + err.Error())
	}

	if len(attrs) == 0 {
		return errors.New("attributes cannot be empty")
	}

	// Create the update expression for SET
	updateExpression := "SET "
	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	for attributeName, attributeValue := range attrs {
		updateExpression += fmt.Sprintf("%s = :%s, ", attributeName, attributeName)
		expressionAttributeValues[fmt.Sprintf(":%s", attributeName)] = attributeValue
	}

	// Remove the trailing comma and space
	updateExpression = updateExpression[:len(updateExpression)-2]

	// Interface to update a record by partition key and optionally a sort key
	input := &dynamodb.UpdateItemInput{
		TableName:                 jsii.String(table_name),
		Key:                       map[string]*dynamodb.AttributeValue{
			"partition_key": {
				S: jsii.String(partition_key),
			},
			"sort_key": {
				S: jsii.String(sort_key),
			},
		},
		UpdateExpression:          jsii.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}


	_, err = c.Client.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}

