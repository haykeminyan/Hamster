package scn

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/matscus/Hamster/Package/Clients/client"
	"github.com/matscus/Hamster/Package/Clients/client/postgres"
	"github.com/matscus/Hamster/Package/Hosts/hosts"
	"github.com/matscus/Hamster/Package/Scenario/scenario"
)

var (
	//GetResponseAllData - struct for return information if scenarios and hosts.
	GetResponseAllData = GetResponse{}
	//GetState - struct for return state of scenario
	GetState = []scenario.State{}
	//LastRunsParams - sync map for last runs param from scenario
	LastRunsParams sync.Map
	//RunsGenerators - sync map for used generator
	RunsGenerators sync.Map
	//HostsAndUsers - sync map users from remote host
	HostsAndUsers sync.Map
	PgClient      *postgres.PGClient
)

//GetResponse -  struct for response
type GetResponse struct {
	Generators []hosts.Host
	Scenarios  []scenario.Scenario
}

//GeneratorState -  struct for generator state
type GeneratorState struct {
	Host  string `json:"name"`
	State string `json:"state"`
}

//InitData - func to init GetRes
func InitData() (err error) {
	scenarios, err := PgClient.GetAllScenarios()
	if err != nil {
		return err
	}
	l := len(*scenarios)
	GetResponseAllData.Scenarios = make([]scenario.Scenario, 0, l)
	for i := 0; i < l; i++ {
		var s scenario.Scenario
		t := (*scenarios)[i]
		s.ID = t.ID
		s.Name = t.Name
		s.Type = t.Type
		s.LastModified = t.LastModified
		s.Gun = t.Gun
		s.Projects = t.Projects
		var tgp []scenario.ThreadGroup

		err := json.Unmarshal([]byte(t.TreadGroups), &tgp)
		if err != nil {
			return (err)
		}
		s.ThreadGroups = tgp
		GetResponseAllData.Scenarios = append(GetResponseAllData.Scenarios, s)
	}
	gen, err := PgClient.GetAllGenerators()
	if err != nil {
		return err
	}
	l = len(gen)
	GetResponseAllData.Generators = make([]hosts.Host, 0, l)
	for i := 0; i < l; i++ {
		var g hosts.Host
		g.ID = gen[i].ID
		g.Host = gen[i].Host
		g.Type = gen[i].Type
		g.Projects = gen[i].Projects
		_, ok := RunsGenerators.Load(g.Host)
		if ok {
			g.State = "IsBusy"
		} else {
			g.State = "Free"
		}
		GetResponseAllData.Generators = append(GetResponseAllData.Generators, g)
	}
	hostsAndUsers, err := PgClient.GetUsersAndHosts()
	for k, v := range hostsAndUsers {
		HostsAndUsers.Store(k, v)
	}
	return nil
}

//SetState -  init state struct for ws
func SetState(s bool, id int64, n string, t string, d int64, gun string, g []hosts.Host) {
	if s {
		starttime := (time.Now().Unix() - time.Unix(10800, 0).Unix())
		endtime := (starttime + time.Unix(d, 0).Unix())
		var srv = scenario.State{}
		srv.RunID = id
		srv.Name = n
		srv.Type = t
		srv.StartTime = starttime
		srv.EndTime = endtime
		srv.Gun = gun
		srv.Generators = g
		GetState = append(GetState, srv)
		return
	}
	for i := 0; i < len(GetState); i++ {
		if id == GetState[i].RunID {
			removeState(i)
		}
	}
}

func removeState(s int) {
	GetState = append(GetState[:s], GetState[s+1:]...)
}

//FloatToString - convert type float  to type string
func FloatToString(i float64) string {
	return strconv.FormatFloat(i, 'f', 2, 64)
}

//СheckRun - fucn for check state scenario
func СheckRun() (res bool) {
	if len(GetState) > 0 {
		res = true
	}
	return res
}

//CheckGen - func fo check state generators
func CheckGen(g []hosts.Host) (res []GeneratorState, err error) {

	l := len(g)
	for i := 0; i < l; i++ {
		user, _ := HostsAndUsers.Load(g[i].Host)
		client, err := client.SSHClient{}.New(user.(string))
		if err != nil {
			return res, err
		}
		_, err = client.Ping(g[i].Host)
		if err != nil {
			var genstate GeneratorState
			genstate.Host = g[i].Host
			genstate.State = "NotAvailable"
			res = append(res, genstate)
			return res, err
		}
	}
	for i := 0; i < l; i++ {
		host, ok := RunsGenerators.Load(g[i].Host)
		if ok {
			var genstate GeneratorState
			genstate.Host = host.(string)
			genstate.State = "IsBusy"
			res = append(res, genstate)
		}
	}
	return res, err
}

//CheckRunsGen - func for chack state generators and change state
func CheckRunsGen() {
	l := len(GetResponseAllData.Generators)
	for i := 0; i < l; i++ {
		_, ok := RunsGenerators.Load(GetResponseAllData.Generators[i].Host)
		if ok {
			GetResponseAllData.Generators[i].State = "IsBusy"
		} else {
			GetResponseAllData.Generators[i].State = "Free"
		}
	}
}
