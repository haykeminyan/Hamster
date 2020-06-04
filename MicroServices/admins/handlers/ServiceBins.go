package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/matscus/Hamster/Package/JWTToken/jwttoken"
	"github.com/matscus/Hamster/Package/Services/service"
	"github.com/matscus/Hamster/Package/httperror"
)

//GetAllServiceBins -  handle function, for return ALL servicebins info
func GetAllServiceBins(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	projects := jwttoken.GetUserProjects(strings.TrimSpace(splitToken[1]))
	projectsID, err := pgClient.GetProjectsIDtoString(projects)
	if err != nil {
		httperror.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	bins, err := pgClient.GetAllServiceBinsByOwner(projectsID)
	if err != nil {
		httperror.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = json.NewEncoder(w).Encode(bins)
	if err != nil {
		httperror.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

//GetAllServiceBinsType -  handle function, for return ALL servicebins info
func GetAllServiceBinsType(w http.ResponseWriter, r *http.Request) {
	bins, err := pgClient.GetAllServiceBinsType()
	if err != nil {
		httperror.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = json.NewEncoder(w).Encode(bins)
	if err != nil {
		httperror.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

//ServiceBins -  handle function, for new,update and delete host
func ServiceBins(w http.ResponseWriter, r *http.Request) {
	var s service.Service
	s.Type = r.FormValue("serviceType")
	s.RunSTR = r.FormValue("runSTR")
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	own := jwttoken.GetUser(strings.TrimSpace(splitToken[1]))
	s.Projects = strings.Split(r.FormValue("projects"), ";")
	s.DBClient = pgClient
	switch r.Method {
	case "POST":
		file, header, err := r.FormFile("uploadFile")
		if err != nil {
			httperror.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		s.Name = header.Filename
		err = s.CreateBin(file, own)
		if err != nil {
			httperror.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, errWrite := w.Write([]byte("{\"Message\":\"Bins created \"}"))
		if errWrite != nil {
			log.Printf("[ERROR] Bins created, but Not Writing to ResponseWriter due: %s", errWrite.Error())
		}
	case "PUT":
		if own == s.Owner || own == "god" {
			id, err := strconv.Atoi(r.FormValue("serviceID"))
			if err == nil {
				s.ID = int64(id)
			}
			err = s.UpdateBin()
			if err != nil {
				httperror.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, errWrite := w.Write([]byte("{\"Message\":\"Host updated \"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Host updated, but Not Writing to ResponseWriter due: %s", errWrite.Error())
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			_, errWrite := w.Write([]byte("{\"Message\": You are not a owner for this service}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error  due: %s", errWrite.Error())
			}
			return
		}
	case "DELETE":
		s.Name = r.FormValue("fileName")
		id, err := strconv.Atoi(r.FormValue("serviceID"))
		if err == nil {
			s.ID = int64(id)
		}
		if own == s.Owner || own == "god" {
			err := s.DeleteBin()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, errWrite := w.Write([]byte("{\"Message\":\"" + err.Error() + "\"}"))
				if errWrite != nil {
					log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
				}
				return
			}
			w.WriteHeader(http.StatusOK)
			_, errWrite := w.Write([]byte("{\"Message\":\"Host deleted \"}"))
			if errWrite != nil {
				log.Printf("[ERROR] Host deleted, but Not Writing to ResponseWriter due: %s", errWrite.Error())
			}
		} else {
			w.WriteHeader(http.StatusOK)
			_, errWrite := w.Write([]byte("{\"Message\": You are not a owner for this service}"))
			if errWrite != nil {
				log.Printf("[ERROR] Not Writing to ResponseWriter error %s due: %s", err.Error(), errWrite.Error())
			}
			return
		}
	}
}
