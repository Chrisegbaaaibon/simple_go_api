package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"auth-api/models/project"
)


func GetAllProjects(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	projects := []project.Project{}
	db.Find(&projects)
	respondJSON(w, http.StatusOK, projects)
}

func getProjectOr404(db *mongo.Database, id string, w http.ResponseWriter, r *http.Request) *project.Project {
	project := project.Project{}
	if err := db.First(&project, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return nil
	}
	return &project
}

func CreateProject(db *mongo.Database, w http.ResponseWriter, r *http.Request){
	project := project.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, project)	
}

func GetProject(db *mongo.Database, w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}

	respondJSON(w, http.StatusOK, project)
}

func UpdateProject(db *mongo.Database, w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *mongo.Database, w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}