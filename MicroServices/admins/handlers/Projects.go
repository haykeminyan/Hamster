package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/matscus/Hamster/Package/Projects/projects"
)

//GetAllProjects -  return all projects
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	allproject, err := pgClient.GetAllProjects()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
		}
		return
	}
	err = json.NewEncoder(w).Encode(allproject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
		}
	}
}

//Projects -  create new Project, update Project and delete Project
func Projects(w http.ResponseWriter, r *http.Request) {
	project := projects.Project{}
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
		}
		return
	}
	project.DBClient = pgClient
	switch r.Method {
	case "POST":
		err = project.Create()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
		os.MkdirAll(os.Getenv("DIRPROJECTS")+"/"+project.Name+"/jmeter", 0777)
		if os.IsExist(err) {
			//todo
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, errWrite := w.Write([]byte("{\"Message\":\"Project created \"}"))
		if errWrite != nil {
			log.Printf("[ERROR]User created, but  Not Writing to ResponseWriter due: %s", errWrite.Error())
		}
	case "PUT":
		name, err := project.DBClient.GetProjectName(project.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
		os.Rename(os.Getenv("DIRPROJECTS")+"/"+name, os.Getenv("DIRPROJECTS")+"/"+project.Name)
		err = project.Update()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, errWrite := w.Write([]byte("{\"Message\":\"Project updated \"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Project updated, but Not Writing to ResponseWriter due: %s", errWrite.Error())
		}
	case "DELETE":
		os.Remove(os.Getenv("DIRPROJECTS") + "/" + project.Name)
		err = project.Delete()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, errWrite := w.Write([]byte("{\"Message\":\"Project deleted \"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Project deleted, but Not Writing to ResponseWriter due: %s", errWrite.Error())
		}
	}
}
