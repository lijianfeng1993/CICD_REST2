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
	//"net/http"
	"os"
	"strings"
	//"time"
)

func Handle_xml_python_cicd(jobname string, url string, filename string) (string, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	configString := string(buf)

	el, _ := dom4g.LoadByXml(configString)

	el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value = url
	//fmt.Println(el.Node("scm").Node("userRemoteConfigs").Node("hudson.plugins.git.UserRemoteConfig").Node("url").Value)
	//fmt.Println(el.ToString())

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

	old_command := `ssh root@host &apos;docker rm -f jobname_temp || true&apos;
ssh root@host &apos;cd /data/jenkins/workspace/jobname_temp &amp;&amp; docker build -t jobname_temp:1.0 .&apos;
ssh root@host &apos;docker run -d --name jobname_temp -p 4501:4501 jobname_temp:1.0&apos;`
	new_command1 := strings.Replace(old_command, "jobname_temp", jobname, 5)
	new_command2 := strings.Replace(new_command1, "host", "10.132.47.15", 3)
	el.Node("builders").Node("hudson.tasks.Shell").Node("command").Value = new_command2

	return el.ToString(), nil
}

func GetXmlConfig_cicd(jobname string, language string, url string) (string, error) {
	if language == "python" {
		xmlconfigstring, err := Handle_xml_python_cicd(jobname, url, "template/python_template_cicd.xml")
		if err != nil {
			fmt.Println("get configstring error.")
			return "Error in GetXmlConfig", nil
		}
		return xmlconfigstring, nil
	}
	/*
		if language == "golang" {
			xmlconfigstring, err := Handle_xml_golang(jobname, url, "template/golang_template.xml")
			if err != nil {
				fmt.Println("get configstring error.")
				return "Error in GetXmlConfig", nil
			}
			return xmlconfigstring, nil
		}
		if language == "java" {
			xmlconfigstring, err := Handle_xml_java(jobname, url, "template/java_template.xml")
			if err != nil {
				fmt.Println("get configstring error.")
				return "Error in GetXmlConfig", nil
			}
			return xmlconfigstring, nil
		}
	*/
	return "Error in GetXmlConfig", nil
}

func CreateJob_cicd(jobname string, language string, url string) map[string]string {

	configString, err := GetXmlConfig_cicd(jobname, language, url)
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
