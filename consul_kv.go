package onlinelab

import (
	"github.com/hashicorp/consul/api"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var kv *api.KV

func NewConsulConfigStorage(consulAddress string) (*ConsulConfigStorage, error) {
	//var err error
	client, err := api.NewClient(&api.Config{Address: consulAddress})
	if err != nil {
		return nil, err
	}

	// Get a handle to the KV API
	kv = client.KV()
	return &ConsulConfigStorage{}, nil
}

type ConsulConfigStorage struct {
	config Config
}

func (cs *ConsulConfigStorage) GetConfig(labName string) (Config, error) {
	pair, _, err := kv.Get(labName, nil)
	if err != nil {
		return cs.config, err
	}
	cs.config.Name = pair.Key
	if err := json.Unmarshal(pair.Value, &cs.config.Treatments); err != nil {
		return cs.config, err
	}
	return cs.config, nil
}

func (cs *ConsulConfigStorage) SetConfig(config Config) error {
	// PUT a new KV pair
	value, err := json.Marshal(config.Treatments)
	if err != nil {
		return err
	}
	p := &api.KVPair{Key: config.Name, Value: value}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}
	cs.config = config
	return nil
}
