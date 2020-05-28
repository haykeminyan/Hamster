package scenario

import (
	"encoding/json"

	"github.com/matscus/Hamster/Package/Clients/client/postgres"
)

//Scenario - struct for scenario
type Scenario struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Gun          string             `json:"gun"`
	LastModified string             `json:"lastmodified"`
	Projects     string             `json:"projects"`
	ThreadGroups []ThreadGroup      `json:"threadGroups"`
	DBClient     *postgres.PGClient `json:",omitempty"`
}

//InsertToDB - func for insert new scenario values in table
func (s *Scenario) InsertToDB() (err error) {
	err = s.DBClient.NewScenario(s.Name, s.Type, s.Gun, s.Projects, paramsToString(s.ThreadGroups))
	return err
}

//CheckScenario - func for delete secenario(rows db and files)
func (s *Scenario) CheckScenario() (res bool, err error) {
	res, err = s.DBClient.CheckScenario(s.Name, s.Gun, s.Projects)
	if err != nil {
		return res, err
	}
	return res, nil
}

//GetNameForID - func for insert new scenario values in table
func (s *Scenario) GetNameForID() (res string, err error) {
	res, err = s.DBClient.GetScenarioName(s.ID)
	return res, err
}

//Update - func for update scenario values in table
func (s *Scenario) Update() (err error) {
	return s.DBClient.UpdateScenario(s.ID, s.Name, s.Type, s.Gun, s.Projects, paramsToString(s.ThreadGroups))
}

//DeleteScenario - func for delete secenario(rows db and files)
func (s *Scenario) DeleteScenario() (err error) {
	err = s.DBClient.DeleteScenario(s.ID)
	return err
}

func paramsToString(params []ThreadGroup) string {
	res, _ := json.MarshalIndent(params, "", " ")
	return string(res)
}
