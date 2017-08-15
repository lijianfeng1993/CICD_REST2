package main

import (
	"CICD_REST2/jenkins"
	"CICD_REST2/sonar"
	//"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/v1/sonarresult/:jobname", GetSonar),
		rest.Post("/v1//createjenkinsjob", CreateJenkinsJob),
		rest.Get("/v1/buildjenkinsjob/:jobname", BuildJenkinsJob),
		rest.Delete("/v1/deletejenkinsjob/:jobname", DeleteJenkinsJob),
		rest.Get("/v1/jenkinsconsole/:jobname", GetJenkinsConsole),
		rest.Post("/v2/createjenkinsjob", CreateJenkinsJob_v2),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8888", api.MakeHandler()))
}

type Jobinfo struct {
	Jobname  string
	Language string
	Path     string
	Url      string
}

var lock = sync.RWMutex{}

func GetSonar(w rest.ResponseWriter, r *rest.Request) {
	jobname := r.PathParam("jobname")
	lock.RLock()
	sonarresult := sonar.GetSonarResult(jobname)
	lock.RUnlock()
	w.WriteJson(sonarresult)
}

func CreateJenkinsJob(w rest.ResponseWriter, r *rest.Request) {
	jobinfo := Jobinfo{}
	err := r.DecodeJsonPayload(&jobinfo)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(jobinfo.Jobname)

	if jobinfo.Jobname == "" {
		rest.Error(w, "jobinfo jobname required", 400)
		return
	}
	if jobinfo.Language == "" {
		rest.Error(w, "jobinfo language required", 400)
		return
	}
	if jobinfo.Path == "" {
		rest.Error(w, "jobinfo path required", 400)
		return
	}

	lock.Lock()
	result := jenkins.CreateJob(jobinfo.Jobname, jobinfo.Language, jobinfo.Path)
	lock.Unlock()
	w.WriteJson(result)
}

func CreateJenkinsJob_v2(w rest.ResponseWriter, r *rest.Request) {
	jobinfo := Jobinfo{}
	err := r.DecodeJsonPayload(&jobinfo)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(jobinfo.Jobname)

	if jobinfo.Jobname == "" {
		rest.Error(w, "jobinfo jobname required", 400)
		return
	}
	if jobinfo.Language == "" {
		rest.Error(w, "jobinfo language required", 400)
		return
	}
	if jobinfo.Url == "" {
		rest.Error(w, "jobinfo url required", 400)
		return
	}
	lock.Lock()
	result := jenkins.CreateJob_cicd(jobinfo.Jobname, jobinfo.Language, jobinfo.Url)
	lock.Unlock()
	w.WriteJson(result)
}

func BuildJenkinsJob(w rest.ResponseWriter, r *rest.Request) {
	jobname := r.PathParam("jobname")
	lock.RLock()
	result := jenkins.BuildJob(jobname)
	lock.RUnlock()
	w.WriteJson(result)
}

func DeleteJenkinsJob(w rest.ResponseWriter, r *rest.Request) {
	jobname := r.PathParam("jobname")
	lock.RLock()
	result := jenkins.DeleteJob(jobname)
	lock.RUnlock()
	w.WriteJson(result)
}

func GetJenkinsConsole(w rest.ResponseWriter, r *rest.Request) {
	jobname := r.PathParam("jobname")
	lock.RLock()
	result := jenkins.GetConsole(jobname)
	lock.RUnlock()
	w.WriteJson(result)
}
