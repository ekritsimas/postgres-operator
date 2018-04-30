package apiserver

/*
Copyright 2017-2018 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	crv1 "github.com/crunchydata/postgres-operator/apis/cr/v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ClusterStruct struct {
	CCPImagePrefix  string `yaml:"CCPImagePrefix"`
	CCPImageTag     string `yaml:"CCPImageTag"`
	Policies        string `yaml:"Policies"`
	Port            string `yaml:"Port"`
	User            string `yaml:"User"`
	Database        string `yaml:"Database"`
	PasswordAgeDays string `yaml:"PasswordAgeDays"`
	PasswordLength  string `yaml:"PasswordLength"`
	Strategy        string `yaml:"Strategy"`
	Replicas        string `yaml:"Replicas"`
}

type StorageStruct struct {
	AccessMode         string `yaml:"AccessMode"`
	Size               string `yaml:"Size"`
	StorageType        string `yaml:"StorageType"`
	StorageClass       string `yaml:"StorageClass"`
	Fsgroup            string `yaml:"Fsgroup"`
	SupplementalGroups string `yaml:"SupplementalGroups"`
}

type ContainerResourcesStruct struct {
	RequestsMemory string `yaml:"RequestsMemory"`
	RequestsCPU    string `yaml:"RequestsCPU"`
	LimitsMemory   string `yaml:"LimitsMemory"`
	LimitsCPU      string `yaml:"LimitsCPU"`
}

type PgoStruct struct {
	Audit         bool   `yaml:"Audit"`
	Metrics       bool   `yaml:"Metrics"`
	LSPVCTemplate string `yaml:"LSPVCTemplate"`
	LoadTemplate  string `yaml:"LoadTemplate"`
	COImagePrefix string `yaml:"COImagePrefix"`
	COImageTag    string `yaml:"COImageTag"`
}

type PgoConfig struct {
	BasicAuth                 string                              `yaml:"BasicAuth"`
	Cluster                   ClusterStruct                       `yaml:"Cluster"`
	Pgo                       PgoStruct                           `yaml:"Pgo"`
	ContainerResources        map[string]ContainerResourcesStruct `yaml:"ContainerResources"`
	PrimaryStorage            string                              `yaml:"PrimaryStorage"`
	BackupStorage             string                              `yaml:"BackupStorage"`
	ReplicaStorage            string                              `yaml:"ReplicaStorage"`
	Storage                   map[string]StorageStruct            `yaml:"Storage"`
	DefaultContainerResources string                              `yaml:"DefaultContainerResources"`
}

func (c *PgoConfig) validate() error {
	var err error
	_, ok := c.Storage[c.PrimaryStorage]
	if !ok {
		return errors.New("invalid PrimaryStorage setting")
	}
	_, ok = c.Storage[c.BackupStorage]
	if !ok {
		return errors.New("invalid BackupStorage setting")
	}
	_, ok = c.Storage[c.ReplicaStorage]
	if !ok {
		return errors.New("invalid ReplicaStorage setting")
	}

	if c.Pgo.LSPVCTemplate == "" {
		return errors.New("Pgo.LSPVCTemplate is required")
	}
	if c.Pgo.LoadTemplate == "" {
		return errors.New("Pgo.LoadTemplate is required")
	}
	if c.Pgo.COImagePrefix == "" {
		return errors.New("Pgo.COImagePrefix is required")
	}
	if c.Pgo.COImageTag == "" {
		return errors.New("Pgo.COImageTag is required")
	}

	if c.DefaultContainerResources != "" {
		_, ok = c.ContainerResources[c.DefaultContainerResources]
		if !ok {
			return errors.New("DefaultContainerResources setting invalid")
		}
	}

	if c.Cluster.CCPImagePrefix == "" {
		return errors.New("Cluster.CCPImagePrefix is required")
	}

	if c.Cluster.CCPImageTag == "" {
		return errors.New("Cluster.CCPImageTag is required")
	}
	return err
}

func (c *PgoConfig) getConf() *PgoConfig {

	yamlFile, err := ioutil.ReadFile("config/pgo.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

/**
var Pgo PgoConfig

func main() {
	Pgo.getConf()

	log.Println(Pgo.Cluster.CCPImageTag)
	log.Println(Pgo.PrimaryStorage)
	log.Println(Pgo.Pgo.COImageTag)
	log.Println(Pgo.ContainerResources)
	log.Println(Pgo.Storage)

	Pgo.validate()
}
*/

func (c *PgoConfig) GetStorageSpec(name string) (crv1.PgStorageSpec, error) {
	var err error
	storage := crv1.PgStorageSpec{}

	s, ok := c.Storage[name]
	if !ok {
		err = errors.New("invalid Storage name " + name)
		log.Error(err)
		return storage, err
	}

	storage.StorageClass = s.StorageClass
	storage.AccessMode = s.AccessMode
	storage.Size = s.Size
	storage.StorageType = s.StorageType
	storage.Fsgroup = s.Fsgroup
	storage.SupplementalGroups = s.SupplementalGroups

	return storage, err

}

func (c *PgoConfig) GetContainerResource(name string) (crv1.PgContainerResources, error) {
	var err error
	r := crv1.PgContainerResources{}

	s, ok := c.ContainerResources[name]
	if !ok {
		err = errors.New("invalid Container Resources name " + name)
		log.Error(err)
		return r, err
	}
	r.RequestsMemory = s.RequestsMemory
	r.RequestsCPU = s.RequestsCPU
	r.LimitsMemory = s.LimitsMemory
	r.LimitsCPU = s.LimitsCPU

	return r, err

}