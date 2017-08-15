package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/donnie4w/dom4g"
	//"bytes"
	//"html"
	//"strings"
	//"text/template"
	//"CICD_REST2/sonar"
	//"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"time"
)

func Handle_xml_python(jobname string, path string, filename string) (string, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	configString := string(buf)

	el, _ := dom4g.LoadByXml(configString)

	//el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value = githuburl
	//fmt.Println(el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value)
	//fmt.Println(el.ToString())

	old_commands := `cp -r src_path/* /root/.jenkins/workspace/dst_path/`
	new_commands_1 := strings.Replace(old_commands, "src_path", path, 1)
	new_commands_2 := strings.Replace(new_commands_1, "dst_path", jobname, 1)
	el.Node("builders").Node("hudson.tasks.Shell").Node("command").Value = new_commands_2

	old_sonar_properties := `sonar.projectKey=jobname_temp
sonar.projectName=jobname_temp
sonar.projectVersion=1.0
sonar.sources=./
sonar.language=py
sonar.python.xunit.reportPath=nosetests.xml
sonar.python.coverage.reportPath=coverage.xml
sonar.sourceEncoding=UTF-8`
	new_sonar_properties := strings.Replace(old_sonar_properties, "jobname_temp", jobname, 2)

	el.Node("builders").Node("hudson.plugins.sonar.SonarRunnerBuilder").Node("properties").Value = new_sonar_properties
	return el.ToString(), nil
}

func Handle_xml_golang(jobname string, path string, filename string) (string, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	configString := string(buf)

	el, _ := dom4g.LoadByXml(configString)

	//el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value = githuburl
	//fmt.Println(el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value)
	//fmt.Println(el.ToString())

	old_commands := `cp -r src_path/* /root/.jenkins/workspace/dst_path/
/root/go/bin/gometalinter.v1 --checkstyle &gt; report.xml || true`
	new_commands_1 := strings.Replace(old_commands, "src_path", path, 1)
	new_commands_2 := strings.Replace(new_commands_1, "dst_path", jobname, 1)
	el.Node("builders").Node("hudson.tasks.Shell").Node("command").Value = new_commands_2

	old_sonar_properties := `sonar.projectKey=jobname_temp
sonar.projectName=jobname_temp
sonar.projectVersion=1.0
sonar.golint.reportPath=report.xml 
sonar.coverage.reportPath=coverage.xml
sonar.coverage.dtdVerification=false
sonar.test.reportPath=test.xml
sonar.sources=./`
	new_sonar_properties := strings.Replace(old_sonar_properties, "jobname_temp", jobname, 2)

	el.Node("builders").Node("hudson.plugins.sonar.SonarRunnerBuilder").Node("properties").Value = new_sonar_properties
	return el.ToString(), nil
}

func Handle_xml_java(jobname string, path string, filename string) (string, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	configString := string(buf)

	el, _ := dom4g.LoadByXml(configString)

	//el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value = githuburl
	//fmt.Println(el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value)
	//fmt.Println(el.ToString())

	old_commands := `cp -r src_path/* /root/.jenkins/workspace/dst_path/
mvn clean org.jacoco:jacoco-maven-plugin:prepare-agent install -Dmaven.test.failure.ignore=true`
	new_commands_1 := strings.Replace(old_commands, "src_path", path, 1)
	new_commands_2 := strings.Replace(new_commands_1, "dst_path", jobname, 1)
	el.Node("builders").Node("hudson.tasks.Shell").Node("command").Value = new_commands_2

	old_sonar_properties := `# Required metadata
sonar.projectKey=jobname_temp
sonar.projectName=jobname_temp
sonar.projectVersion=1.0
sonar.sources=src
sonar.language=java
sonar.sourceEncoding=UTF-8`
	new_sonar_properties := strings.Replace(old_sonar_properties, "jobname_temp", jobname, 2)

	el.Node("builders").Node("hudson.plugins.sonar.SonarRunnerBuilder").Node("properties").Value = new_sonar_properties
	return el.ToString(), nil
}

func GetXmlConfig(jobname string, language string, path string) (string, error) {
	if language == "python" {
		xmlconfigstring, err := Handle_xml_python(jobname, path, "template/python_template.xml")
		if err != nil {
			fmt.Println("get configstring error.")
			return "Error in GetXmlConfig", nil
		}
		return xmlconfigstring, nil
	}
	if language == "golang" {
		xmlconfigstring, err := Handle_xml_golang(jobname, path, "template/golang_template.xml")
		if err != nil {
			fmt.Println("get configstring error.")
			return "Error in GetXmlConfig", nil
		}
		return xmlconfigstring, nil
	}
	if language == "java" {
		xmlconfigstring, err := Handle_xml_java(jobname, path, "template/java_template.xml")
		if err != nil {
			fmt.Println("get configstring error.")
			return "Error in GetXmlConfig", nil
		}
		return xmlconfigstring, nil
	}
	return "Error in GetXmlConfig", nil
}

func CreateJob(jobname string, language string, path string) map[string]string {

	configString, err := GetXmlConfig(jobname, language, path)
	if err != nil {
		fmt.Println("get configString error.")
	}

	result := make(map[string]string)

	jenkins, err := gojenkins.CreateJenkins(jenkinsurl, jenkinsuser, jenkinspass).Init()
	if err != nil {
		result["status"] = "fail"
		result["info"] = "connect jenkins error"
		return result
	} else {
		fmt.Println("connect jenkins success.")
	}

	_, err = jenkins.CreateJob(configString, jobname)
	if err != nil {
		result["status"] = "fail"
		result["info"] = "create jenkinsjob error"
		return result
	} else {
		fmt.Println("create job success.")
	}

	result["status"] = "success"
	result["info"] = "job has been create"
	return result
}

func DeleteJob(jobname string) map[string]string {

	result := make(map[string]string)

	jenkins, err := gojenkins.CreateJenkins(jenkinsurl, jenkinsuser, jenkinspass).Init()
	if err != nil {
		result["status"] = "fail"
		result["info"] = "connect jenkins error"
		return result
	} else {
		fmt.Println("connect jenkins success.")
	}

	_, err = jenkins.DeleteJob(jobname)
	if err != nil {
		result["status"] = "fail"
		result["info"] = "delete jenkinsjob error"
		return result
	} else {
		fmt.Println("delete job success.")
	}

	result["status"] = "success"
	result["info"] = "delete jenkinsjob success"
	return result
}

func BuildJob(jobname string) map[string]string {

	result := make(map[string]string)

	jenkins, err := gojenkins.CreateJenkins(jenkinsurl, jenkinsuser, jenkinspass).Init()
	if err != nil {
		result["status"] = "fail"
		result["info"] = "connect jenkins error"
		return result
	} else {
		fmt.Println("connect jenkins success.")
	}

	_, err = jenkins.BuildJob(jobname)
	if err != nil {
		result["status"] = "fail"
		result["info"] = "build jenkinsjob error"
		return result
	} else {
		fmt.Println("job has been builted.")
	}

	result["status"] = "success"
	result["info"] = "job has been builted"
	return result
}

func GetConsole(jobname string) map[string]string {
	result := make(map[string]string)
	client := &http.Client{}
	base_url := jenkinsurl + "job/jobname_temp/lastBuild/consoleText"
	new_url := strings.Replace(base_url, "jobname_temp", jobname, 1)

	request, err := http.NewRequest("GET", new_url, nil)
	if err != nil {
		panic(err)
	}

	response, _ := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// handle error
	}
	result["Console Output"] = string(body)
	return result
}

/*
func GetBuildStatus() {


	json_info, err := ReadJsonFile()
	if err != nil {
		fmt.Println("get json_info error.")
	}
	jobname := json_info["jobname"]

	jenkins, err := gojenkins.CreateJenkins(jenkinsurl, jenkinsuser, jenkinspass).Init()
	if err != nil {
		fmt.Println("connect jenkins error.")
	} else {
		fmt.Println("connect jenkins success.")
	}

	builds, err := jenkins.GetAllBuildIds(jobname)
	if err != nil {
		panic(err)
	}
	lastbuild := builds[0]
	lastbuildId := lastbuild.Number
	data, err := jenkins.GetBuild(jobname, lastbuildId)
	if err != nil {
		panic(err)
	}

	for true {
		if data.GetResult() == "SUCCESS" {
			fmt.Println("The sonar result:")
			sonar.GetSonarResult(jobname)
			break
		} else {
			time.Sleep(time.Second)
			fmt.Println("The job is buiding...")
		}
	}
	//fmt.Println(data.GetResult())
}
*/

/*
func main() {
	//fmt.Println(ReadJsonFile())
	//DeleteJob()
	//CreateJob()
	//BuildJob()
	//time.Sleep(10 * time.Second)
	GetBuildStatus()
	//Handle_xml_python("python_template.xml")
	//GetXmlConfig
	//sonar.GetSonarResult("python-pattern")
}
*/
