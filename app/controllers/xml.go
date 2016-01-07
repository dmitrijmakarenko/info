package controllers

import (
	"github.com/robfig/revel"
	"encoding/xml"
	"os"
	"io"
	"io/ioutil"
)

type xmlProperty struct {
	XMLName xml.Name `xml:"property"`
	Name string `xml:"name"`
	Type string `xml:"type"`
	Desc string `xml:"desc"`
}

type xmlEntity struct {
	XMLName xml.Name `xml:"entity"`
	ID string `xml:"id"`
	Name string `xml:"name"`
	Properties []xmlProperty  `xml:"property"`
}

type xmlConfig struct {
	XMLName xml.Name `xml:"config"`
	Entities []xmlEntity  `xml:"entity"`
}

func generateXML(config Config) {
	revel.INFO.Println("generate xml file")
	file, err := os.Create("gocode/src/infosystem/configs/config.xml")
	if err != nil {
		revel.ERROR.Println("create file error", err)
	}
	defer file.Close()

	data := &xmlConfig{}

	for _, ent := range config.Entities {
		entity := xmlEntity{}
		entity.ID = ent.Id
		entity.Name = ent.Name
		for _, prop := range ent.Properties {
			entity.Properties = append(entity.Properties, xmlProperty{Name: prop.Name, Type: prop.Type, Desc: prop.Desc})
		}
		data.Entities = append(data.Entities, entity)
	}

	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(data); err != nil {
		revel.ERROR.Println("error", err)
	}
}

func readXML() Config {
	file, err := os.Open("gocode/src/infosystem/configs/config.xml")
	if err != nil {
		revel.ERROR.Println("read file error", err)
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	var cfgXML xmlConfig
	xml.Unmarshal(data, &cfgXML)
	var cfg Config
	for _, ent := range cfgXML.Entities {
		entity := Entity{}
		entity.Id = ent.ID
		entity.Name = ent.Name
		for _, prop := range ent.Properties {
			entity.Properties = append(entity.Properties, Property{Name: prop.Name, Type: prop.Type, Desc: prop.Desc})
		}
		cfg.Entities = append(cfg.Entities, entity)
	}
	return cfg
}