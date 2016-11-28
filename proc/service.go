package proc

import (
	"io/ioutil"
	procinfo "github.com/c9s/goprocinfo/linux"
	log "github.com/Sirupsen/logrus"
	"time"
	"strings"
	"os"
	"syscall"
)

const (
	procDirectory = "/proc/"
	statusFile = "/status"
	abc = "abcdefghijklmnopqrstuvwxyz"

)

type service struct{
	Pid uint64
	Name string
	TaskName string
	Ppid int64
	Executor bool
	ChaosTimeStamp int
}

func (s *service) Kill () (err error){
	log.Debug("proc.service.Kill")
	s.ChaosTimeStamp = time.Now().UTC().Nanosecond()
	log.Infof("proc.service.Kill - '%v' '%v' '%v' '%v' '%v'", s.Pid, s.Name, s.Ppid, s.TaskName, s.ChaosTimeStamp)
	err = syscall.Kill(int(s.Pid), 9)
	if err != nil {
		log.Infof("proc.service.Kill - ERROR: '%v'", err.Error())
	}
	return
}

func ReadAllChildProcess (daemons []daemon) (aux []service, err error){
	var ser []service
	for _, d := range daemons {
		if d.Pid > 0 {
			ser, err = readAllChildServices(int64(d.Pid), []string{mesosAgentLogrotate}, true)
			for _ , s := range ser {
				aux = append(aux, s)
			}
		}
	}
	for _, s := range ser {
		ser, err = readAllChildServices(int64(s.Pid), []string{}, false)
		for _ , s := range ser {
			aux = append(aux, s)
		}
	}
	return
}

// Read all child process for parent pid
func readAllChildServices(ppid int64, blackListServices []string, setExecutor bool) (res []service, err error){
	log.Debug("proc.service.ReadAllServices")
	if files, err := ioutil.ReadDir(procDirectory); err == nil {
		for _, file := range files {
			if ! strings.ContainsAny(file.Name(), abc){
				status, err := procinfo.ReadProcessStatus(procDirectory + file.Name() + statusFile)
				if err != nil {
					log.Infof("Error reading file: '%v'. ERROR: '%v'", file.Name(), err.Error())
				} else {
					if ppid == status.PPid && !isInBlackList(status.Name, blackListServices) {
						link, _ := os.Readlink(procDirectory + file.Name() + "/cwd")
						taskName := strings.Split(link, "/")[10]
						res = append(res, service{Pid: status.Pid, Name: status.Name, Ppid: status.PPid, Executor: setExecutor, TaskName: taskName})
						log.Debugf("proc.service.ReadAllServices - append - '%v' '%v' '%v' '%v'", taskName, status.Pid, status.Name, status.PPid)
					}
				}
			}

		}
	}
	if err != nil {
		log.Infof("proc.service.Kill - ERROR: '%v'", err.Error())
	}
	log.Debugf("proc.service.ReadAllService - lenService: '%v'", len(res))
	return
}

func isInBlackList (name string, blackListServices []string) (res bool){
	log.Debug("proc.service.isInBlackList")
	for _, blackService := range blackListServices {
		if strings.Compare(name, blackService) == 0 {res = true}
	}
	log.Debugf("proc.service.isInBlackList - '%v' blcakList '%v'", name, res)
	return
}

