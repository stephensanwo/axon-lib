package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	axon_coredb "github.com/stephensanwo/axon-lib/coredb"
	axon_types "github.com/stephensanwo/axon-lib/types"
)

type Node struct {
	Session axon_types.Session
	AwsSession *aws_session.Session
}

func (no *Node) GetNodes(a *axon_types.AxonContext, folder_id string, note_id string) (*[]axon_types.Node, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(no.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch nodes - " + err.Error())
	}

	var nodes []axon_types.Node

	nodeResult, err := db.QueryDatabase(axon_types.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", no.Session.SessionData.User.Email, folder_id, note_id), nil)

	if err != nil {
		return nil, errors.New("could not fetch nodes - " + err.Error())
	}

	// Unmarshal the DynamoDB item into a Note | Node | Edges structs
	if err := dynamodbattribute.UnmarshalMap(nodeResult.Item, &nodes); err != nil {
		return nil, err
	}
	
	return &nodes, err

}

func (no *Node) CreateNode(a *axon_types.AxonContext, userNodeData axon_types.NodeData, clientRefPosition axon_types.Position, folder_id string, note_id string) (*axon_types.Node, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(no.AwsSession)
	if err != nil {
		return nil, errors.New("could not create node - " + err.Error())
	}

	// Confirm that note exists
	var note axon_types.Note 

	noteResult, err := db.QueryDatabase(axon_types.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", no.Session.SessionData.User.Email, folder_id), &note_id)

	if noteResult.Item == nil || err != nil {
		return nil, errors.New("could not fetch note data - " + err.Error())
	}

	// Unmarshal the DynamoDB item into a Note struct
	if err := dynamodbattribute.UnmarshalMap(noteResult.Item, &note); err != nil {
		return nil, err
	}

	//  Create node object
	node := axon_types.Node{
		UserId:   no.Session.SessionData.User.UserId,
		FolderID: note.FolderID,
		NoteID:   note.NoteID,
		NodeID:   uuid.New().String(),
		// Provided by user
		Data: userNodeData,
		// Defined by client state
		Position: clientRefPosition,
		// Not provided on creation
		Content: axon_types.NodeContent{
			MarkDown: "",
		},
		// Not provided on creation
		Styles:     axon_types.NodeStyles{},
		LastEdited: time.Now(),
	}

	// Add node to Database
	err = db.MutateDatabase(axon_types.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", no.Session.SessionData.User.Email, folder_id, note.NoteID), node.NodeID, node)

	if err != nil {
		return nil, errors.New("could not create node - " + err.Error())
	}
	
	return &node, err

}

func (no *Node) FindNode(a *axon_types.AxonContext, folder_id string, note_id string, node_id string) (*axon_types.Node, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(no.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch node - " + err.Error())
	}
	
	// Fetch the Node
	nodeResult, err := db.QueryDatabase(axon_types.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", no.Session.SessionData.User.Email, folder_id, note_id), &node_id)

	var node axon_types.Node

	// Unmarshal the DynamoDB item into a Node struct
	if err := dynamodbattribute.UnmarshalMap(nodeResult.Item, &node); err != nil {
		return nil, err
	}
		
	return &node, err
}

func (no *Node) DeleteNode(a *axon_types.AxonContext, folder_id string, note_id string, node_id string) (*string, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(no.AwsSession)
	if err != nil {
		return nil, errors.New("could not delete node - " + err.Error())
	}

	err = db.DeleteRecord(axon_types.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", no.Session.SessionData.User.Email, folder_id, note_id), &node_id)

	if err != nil {
		return nil, errors.New("could not delete node or node does not exist - " + err.Error())
	}

	return &node_id, err

}

func (no *Node) UpdateNode(a *axon_types.AxonContext, userNodeData axon_types.NodeData, clientRefPosition axon_types.Position, userContent axon_types.NodeContent, userStyles axon_types.NodeStyles, folder_id string, note_id string, node_id string) (*string, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(no.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note detail - " + err.Error())
	}

	// Create a map to store the updated attributes
	updatedAttributes := make(map[string]*dynamodb.AttributeValue)

	// Update the userNodeData fields if provided
	if userNodeData.Label != "" {
		updatedAttributes["userNodeData.label"] = &dynamodb.AttributeValue{S: &userNodeData.Label}
	}

	if userNodeData.Label != "" {
		updatedAttributes["userNodeData.label"] = &dynamodb.AttributeValue{S: &userNodeData.Label}
	}

	if userNodeData.Description != "" {
		updatedAttributes["userNodeData.description"] = &dynamodb.AttributeValue{S: &userNodeData.Description}
	}

	if userNodeData.NodeCategory != "" {
		updatedAttributes["userNodeData.nodeCategory"] = &dynamodb.AttributeValue{S: &userNodeData.NodeCategory}
	}

	// Update the clientRefPosition fields if provided
	if clientRefPosition.X != 0 {
		updatedAttributes["clientRefPosition.x"] = &dynamodb.AttributeValue{N: jsii.String(strconv.Itoa(clientRefPosition.X))}
	}
	if clientRefPosition.Y != 0 {
		updatedAttributes["clientRefPosition.y"] = &dynamodb.AttributeValue{N: jsii.String(strconv.Itoa(clientRefPosition.Y))}
	}

	// Update the userContent fields if provided
	if userContent.MarkDown != "" {
		updatedAttributes["userContent.markdown"] = &dynamodb.AttributeValue{S: &userContent.MarkDown}
	}

	// Update the userStyles fields if provided
	if len(userStyles.BackgroundStyles) > 0 {
		// Convert the map[string]interface{} to map[string]*dynamodb.AttributeValue
		backgroundStylesAV, err := dynamodbattribute.MarshalMap(userStyles.BackgroundStyles)
		if err != nil {
			return nil, err
		}
		updatedAttributes["userStyles.background_styles"] = &dynamodb.AttributeValue{
			M: backgroundStylesAV,
		}
	}
	if len(userStyles.LabelStyles) > 0 {
		// Convert the map[string]interface{} to map[string]*dynamodb.AttributeValue
		labelStylesAV, err := dynamodbattribute.MarshalMap(userStyles.LabelStyles)
		if err != nil {
			return nil, err
		}
		updatedAttributes["userStyles.label_styles"] = &dynamodb.AttributeValue{
			M: labelStylesAV,
		}
	}
	if len(userStyles.DescriptionStyles) > 0 {
		// Convert the map[string]interface{} to map[string]*dynamodb.AttributeValue
		descriptionStylesAV, err := dynamodbattribute.MarshalMap(userStyles.DescriptionStyles)
		if err != nil {
			return nil, err
		}
		updatedAttributes["userStyles.description_styles"] = &dynamodb.AttributeValue{
			M: descriptionStylesAV,
		}
	}
	// Update the LastEdited field with the current timestamp
	updatedAttributes["last_edited"] = &dynamodb.AttributeValue{
		S: jsii.String(time.Now().Format(time.RFC3339)),
	}

	err = db.UpdateRecord(axon_types.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", no.Session.SessionData.User.Email, folder_id, note_id), node_id, updatedAttributes)

	return &node_id, err

}
