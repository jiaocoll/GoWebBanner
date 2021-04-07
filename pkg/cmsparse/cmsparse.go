package cmsparse

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

type Cms struct {
	Name    string
	Path    string
	Option  string
	Content string
}


type Coup struct {
	Path   string
	length int
}

type CoupList []Coup

func (c CoupList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CoupList) Len() int {
	return len(c)
}

func (c CoupList) Less(i, j int) bool {
	return c[i].length > c[j].length
}

func (c CoupList) More(i, j int) bool {
	return c[i].length < c[j].length
}

func sortMapByValue(m map[string]int) CoupList {
	c := make(CoupList, len(m))
	i := 0
	for k, v := range m {
		c[i] = Coup{k, v}
		i++
	}
	sort.Sort(c)
	return c
}


func Parsecms(filename string) (CoupList, map[string][]Cms){
	var unjsonfile map[string][]Cms
	cmsfile ,err := ioutil.ReadFile(filename)
	if err != nil{
		log.Println(err)
	}
	err = json.Unmarshal(cmsfile,&unjsonfile)
	webdata := make(map[string][]Cms)

	for k, v := range unjsonfile{
		for _, data := range v{
			path := data.Path
			_,ok := webdata[path]
			if !ok {
				webdata[path] = make([]Cms,0)
			}
			data.Name = k
			webdata[path] = append(webdata[path], data)
		}
	}
	Data := make(map[string]int)
	for k, v := range webdata {
		Data[k] = len(v)
	}
	sortCoups := sortMapByValue(Data)
	return sortCoups, webdata

}


