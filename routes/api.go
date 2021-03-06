package routes

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/proc"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type response struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Proc   []proc.Process `json:"process,omitempty"`
}

type responseError struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Cause  string `json:"cause,omitempty"`
	Proc   []proc.Process `json:"processKill,omitempty"`
}

type ListRequest struct {
	Code string 		`json:"code,omitempty"`
	Daemon []proc.Daemon	`json:"daemon,omitempty"`
	Docker []proc.Docker	`json:"docker,omitempty"`
	Service []proc.Service	`json:"service,omitempty"`
}

type ShooterRequest struct {
	Name string
	ServiceType int
	KillExecutor int
}

// handleShooter kill process by name
func handleShooter(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	//.
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var sr ShooterRequest
	err := decoder.Decode(&sr)
	if err != nil {
		log.Warnf("routes.api.handleShooter - '%v'", err.Error())
		json.NewEncoder(w).Encode(responseError{Code: "400", Status: "Invalid params", Cause: err.Error()})
		return nil
	}
	if strings.EqualFold(sr.Name, "") {
		log.Warnf("routes.api.handleShooter - Name field is mandatory and should not be empty")
		json.NewEncoder(w).Encode(responseError{Code: "400", Status: "Invalid params", Cause: "Name field is mandatory and should not be empty"})
		return nil
	}
	//.
	var p = new(proc.Manager)
	p.ConfigManager()
	errLoad := p.LoadProcesses()
	if errLoad != nil {
		log.Errorf("routes.api.handleShooter - '%v'", err.Error())
		json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: err.Error()})
		return nil
	}
	//.
	proc, err := p.Shooter(sr.Name, sr.ServiceType, sr.KillExecutor)
	switch {
		case  len(proc) > 0 && err == nil:
			log.Infof("routes.api.handleShooter - %v", proc)
			json.NewEncoder(w).Encode(response{Code: "200", Status: "Succes", Proc: proc})
		case  len(proc) == 0 && err == nil:
			log.Warning("routes.api.handleShooter - 'TASK NOT FOUND'")
			json.NewEncoder(w).Encode(responseError{Code: "404", Status: "Empty result", Cause: "TASK NOT FOUND"})
		case  len(proc) > 0 && err != nil:
			log.Warningf("routes.api.handleShooter - %v", proc)
			json.NewEncoder(w).Encode(responseError{Code: "500", Status: "The kill process has been interrupted because errors have been encountered. Some tasks may be removed.", Cause: err.Error(), Proc: proc})
		case  err != nil:
			log.Errorf("routes.api.handleShooter - '%v'", err.Error())
			json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: err.Error()})
		default:
			log.Errorf("routes.api.handleShooter - Unknow option")
			json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: "Unknow option"})
	}
	return nil
}

// handleList return a list of process available
func handleList(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	//.
	var p = new(proc.Manager)
	p.ConfigManager()
	//.
	if err := p.LoadProcesses(); err == nil{
		log.Infof("routes.api.handleList - %v", p)
		json.NewEncoder(w).Encode(ListRequest{Code: "200", Daemon: p.Daemons, Docker: p.Dockers, Service: p.Services})
	} else {
		log.Errorf("routes.api.handleList - '%v'", err.Error())
		json.NewEncoder(w).Encode(responseError{Code: "500", Status: "Server error", Cause: err.Error()})
	}
	return nil
}
