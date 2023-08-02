package core

import (
	"errors"
	"fmt"
	"time"

	aws_session "github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	axon_coredb "github.com/stephensanwo/axon-lib/coredb"
	axon_types "github.com/stephensanwo/axon-lib/types"
)

type Note struct {
	Session axon_types.Session
		AwsSession *aws_session.Session
}

// Gets the note data by ID and all the nodes and edges associated with it
func (n *Note) GetNoteDetail(a *axon_types.AxonContext, folder_id string, note_id string) (*axon_types.NoteDetail, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note detail - " + err.Error())
	}

	// Fetch the Note
	noteResult, err := db.QueryDatabase(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id), &note_id)

	if err != nil {
		return nil, errors.New("could not fetch note - " + err.Error())
	}

	// Fetch Nodes and Edges
	nodeResult, err := db.QueryDatabasePartition(axon_coredb.AXON_TABLE, fmt.Sprintf("NODE#%s#%s#%s", n.Session.SessionData.User.Email, folder_id, note_id))

	if err != nil {
		return nil, errors.New("could not fetch node details - " + err.Error())
	}


	edgeResult, err := db.QueryDatabasePartition(axon_coredb.AXON_TABLE, fmt.Sprintf("EDGE#%s#%s#%s", n.Session.SessionData.User.Email, folder_id, note_id))

	if err != nil {
		return nil, errors.New("could not fetch edge details - " + err.Error())
	}

	var noteData axon_types.NoteDetail
	var note axon_types.Note
	var nodes []axon_types.Node
	var edges []axon_types.Edge

	// Unmarshal the DynamoDB item into a Note | Node | Edges structs
	if err := dynamodbattribute.UnmarshalMap(noteResult.Item, &note); err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(nodeResult.Items, &nodes); err != nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(edgeResult.Items, &edges); err != nil {
		return nil, err
	}

	noteData.UserId = note.UserId
	noteData.FolderID = note.FolderID
	noteData.NoteID = note.NoteID
	noteData.NoteName = note.NoteName
	noteData.Description = note.Description
	noteData.DateCreated = note.DateCreated
	noteData.LastEdited = note.LastEdited
	noteData.Nodes = nodes
	noteData.Edges = edges

	return &noteData, err
}

func (n *Note) GetNotes(a *axon_types.AxonContext, folder_id string) (*[]axon_types.Note, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note detail - " + err.Error())
	}

	// Fetch the Note
	notesResult, err := db.QueryDatabasePartition(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id))

	if err != nil {
		return nil, errors.New("could not fetch notes - " + err.Error())
	}

	var notes []axon_types.Note

	// Unmarshal the DynamoDB item into a Note struct
	if err := dynamodbattribute.UnmarshalListOfMaps(notesResult.Items, &notes); err != nil {
		return nil, err
	}

	return &notes, err

}

func (n *Note) CreateNote(a *axon_types.AxonContext, note_name string, description string, folder_id string) (*string, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note detail - " + err.Error())
	}

	//  Create note object
	note := axon_types.Note{
		UserId:      n.Session.SessionData.User.UserId,
		FolderID:    folder_id,
		NoteID:      uuid.New().String(),
		NoteName:        note_name,
		Description: description,
		DateCreated: time.Now(),
		LastEdited:  time.Now(),
	}

	// Add note to Database
	err = db.MutateDatabase(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id), note.NoteID, note)

	if err != nil {
		return nil, errors.New("could not create note - " + err.Error())
	}

	return &note.NoteID, err
}

func (n *Note) FindNote(a *axon_types.AxonContext, folder_id string, note_id string) (*axon_types.Note, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note - " + err.Error())
	}

	// Fetch the Note
	noteResult, err := db.QueryDatabase(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id), &note_id)

	var note axon_types.Note

	// Unmarshal the DynamoDB item into a Note struct
	if err := dynamodbattribute.UnmarshalMap(noteResult.Item, &note); err != nil {
		return nil, err
	}
	
	return &note, err
}

func (n *Note) DeleteNote(a *axon_types.AxonContext, folder_id string, note_id string) (*string, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not delete note - " + err.Error())
	}
	
	err = db.DeleteRecord(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id), &note_id)

	if err != nil {
		return nil, errors.New("could not delete note or note does not exist - " + err.Error())
	}

	return &note_id, err

}

func (n *Note) UpdateNote(a *axon_types.AxonContext, name *string, description *string, folder_id string, note_id string) (*string, error) {

	// Create the DynamoDB client
	db, err := axon_coredb.NewDb(n.AwsSession)
	if err != nil {
		return nil, errors.New("could not fetch note detail - " + err.Error())
	}
	
	// Create a map to store the updated attributes
	updatedAttributes := make(map[string]*dynamodb.AttributeValue)

	// Check if the name field is provided and update it
	if name != nil {
		updatedAttributes["name"] = &dynamodb.AttributeValue{S: name}
	}
 
	// Check if the description field is provided and update it
	if description != nil {
		updatedAttributes["description"] = &dynamodb.AttributeValue{S: description}
	}

	// Update the LastEdited field with the current timestamp
	updatedAttributes["last_edited"] = &dynamodb.AttributeValue{
		S: jsii.String(time.Now().Format(time.RFC3339)),
	}

	err = db.UpdateRecord(axon_coredb.AXON_TABLE, fmt.Sprintf("NOTE#%s#%s", n.Session.SessionData.User.Email, folder_id), note_id, updatedAttributes)


	return &note_id, err

}
