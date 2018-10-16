package onlinelab

import (
	"github.com/hashicorp/consul/api"
	"github.com/json-iterator/go"
)

func NewConsulConfigStorage(consulAddress string) (*ConsulConfigStorage, error) {
	//var err error
	client, err := api.NewClient(&api.Config{Address: consulAddress})
	if err != nil {
		return nil, err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return &ConsulConfigStorage{kv: client.KV(), json: json}, nil
}

type ConsulConfigStorage struct {
	config Config
	kv     *api.KV
	json   jsoniter.API
}

func (cs *ConsulConfigStorage) GetConfig(labName string) (Config, error) {
	pair, _, err := cs.kv.Get(labName, nil)
	if err != nil {
		// get config failure from consul kv, cs.config default value is &Config{}
		return cs.config, err
	}
	if err := cs.json.Unmarshal(pair.Value, &cs.config.treatments); err != nil {
		// unmarshal json failure, cs.config default value is &Config{}
		return cs.config, err
	}
	cs.config.Name = pair.Key
	return cs.config, nil
}

func (cs *ConsulConfigStorage) SetConfig(config Config) {
	// PUT a new KV pair
	value, _ := cs.json.Marshal(config.treatments)
	p := &api.KVPair{Key: config.Name, Value: value}
	cs.kv.Put(p, nil)
	cs.config = config
}
