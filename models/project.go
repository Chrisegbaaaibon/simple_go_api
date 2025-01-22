package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Title string `bson:"title" json:"title"`
	Archived bool `bson:"archived" json:"archived"`
	Tasks []Task `bson:"tasks" json:"tasks"`
}

type Priority string

const (
	Low Priority = "0"
	Medium Priority = "1"
	High Priority = "2"
	Urgent Priority = "3"
)

type Task struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id"`
	IsDone bool `bson:"is_done" json:"is_done"`
	Priority Priority `bson:"priority" json:"priority" default:"0"`
}

func (p *Project) Archive(){
	p.Archived = false
}

func (p *Project) Unarchive(){
	p.Archived = true
}
