package types

import (
	"time"
)

type AxonData struct {
	Email   string   `json:"email"`
	Folders []Folder `json:"folders"`
}

type FolderList struct {
	UserId      string `json:"user_id"`
	FolderID    string `json:"folder_id"`
	FolderName        string             `json:"folder_name"`
	DateCreated time.Time          `json:"date_created"`
	LastEdited  time.Time          `json:"last_edited"`
	Notes       []Note             `json:"notes"`
}

type NoteDetail struct {
	UserId      string `json:"user_id"`
	FolderID    string `json:"folder_id"`
	NoteID      string `json:"note_id"`
	NoteName        string             `json:"note_name"`
	Description string             `json:"description"`
	DateCreated time.Time          `json:"date_created"`
	LastEdited  time.Time          `json:"last_edited"`
	Nodes       []Node             `json:"nodes"`
	Edges       []Edge             `json:"edges"`
}

type Folder struct {
	UserId      string 			   `json:"user_id"`
	FolderID    string 			   `json:"folder_id"`
	FolderName        string             `json:"folder_name" `
	DateCreated time.Time          `json:"date_created"`
	LastEdited  time.Time          `json:"last_edited"`
}

type Note struct {
	UserId      string `json:"user_id"`
	FolderID    string `json:"folder_id"`
	NoteID      string `json:"note_id"`
	NoteName        string             `json:"note_name"`
	Description string             `json:"description"`
	DateCreated time.Time          `json:"date_created"`
	LastEdited  time.Time          `json:"last_edited"`
}

type Node struct {
	UserId     string `json:"user_id"`
	FolderID   string `json:"folder_id"`
	NoteID     string `json:"note_id"`
	NodeID     string `json:"node_id"`
	Data       NodeData           `json:"data"`
	Position   Position           `json:"position"`
	Content    NodeContent        `json:"node_content"`
	Styles     NodeStyles         `json:"node_styles"`
	LastEdited time.Time          `json:"last_edited"`
}

type NodeData struct {
	Label        string `json:"label"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	NodeCategory string `json:"node_category"`
}

type Position struct {
	X int `json:"x" bson:"x"`
	Y int `json:"y" bson:"y"`
}

type NodeContent struct {
	MarkDown string `json:"markdown" bson:"markdown"`
}

type NodeStyles struct {
	BackgroundStyles  map[string]interface{} `json:"background_styles"`
	LabelStyles       map[string]interface{} `json:"label_styles"`
	DescriptionStyles map[string]interface{} `json:"description_styles"`
}

type Edge struct {
	UserId     string `json:"user_id"`
	FolderID   string `json:"folder_id"`
	NoteID     string `json:"note_id"`
	EdgeID     string `json:"edge_id"`
	SourceID   string `json:"source"`
	TargetID   string `json:"target"`
	Animated   bool               `json:"animated"`
	Label      string             `json:"label"`
	EdgeType   string             `json:"edge_type"`
	LastEdited time.Time          `json:"last_edited"`
}
