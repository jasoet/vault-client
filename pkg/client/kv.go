package client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/jasoet/vault-client/pkg/util"
)

type kvEngine struct {
	vaultClient *api.Client
	path        string
}

func (k kvEngine) Write(path string, input interface{}) (metadata *KVMetadata, err error) {
	payload := map[string]interface{}{
		"data": util.StructToMap(input),
	}

	result, err := k.vaultClient.Logical().Write(fmt.Sprintf("%v/data/%v", k.path, path), payload)
	if err != nil {
		return
	}

	metadata = new(KVMetadata)
	err = util.MapToStruct(result.Data, metadata)
	return
}

func (k kvEngine) Read(path string, output interface{}) (metadata *KVMetadata, err error) {
	secret, err := k.vaultClient.Logical().Read(fmt.Sprintf("%v/data/%v", k.path, path))
	if err != nil || secret == nil {
		return
	}

	if val, ok := secret.Data["data"]; ok {
		err = util.MapToStruct(val, output)
		if err != nil {
			return
		}
	}

	if val, ok := secret.Data["metadata"]; ok {
		metadata = new(KVMetadata)
		err = util.MapToStruct(val, metadata)
		if err != nil {
			return
		}
	}

	return
}

func (k kvEngine) WriteConfig(config KVConfig) (err error) {
	_, err = k.vaultClient.Logical().Write(fmt.Sprintf("%v/config", k.path), util.StructToMap(config))
	return
}

func (k kvEngine) ReadConfig() (config *KVConfig, err error) {
	result, err := k.vaultClient.Logical().Read(fmt.Sprintf("%v/config", k.path))
	if err != nil {
		return
	}

	if result == nil {
		err = fmt.Errorf("%v secret config is not found", k.path)
		return
	}

	config = new(KVConfig)
	err = util.MapToStruct(result.Data, config)
	return
}

func (k kvEngine) Path() string {
	return k.path
}

func (k kvEngine) Enable() (err error) {
	data := map[string]interface{}{"type": "kv-v2"}
	_, err = k.vaultClient.Logical().Write(fmt.Sprintf("/sys/mounts/%v", k.path), data)
	return
}

func (k kvEngine) Status() (status *SecretStatus, err error) {
	result, err := k.vaultClient.Logical().Read(fmt.Sprintf("/sys/mounts/%v/tune", k.path))
	if err != nil || result == nil {
		return
	}

	status = new(SecretStatus)
	err = util.MapToStruct(result.Data, status)

	return
}

func DefaultKV() (KV, error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &kvEngine{vaultClient: vaultClient, path: "secret"}, nil
}

func NewKV(vaultClient *api.Client, path string) (KV, error) {
	return &kvEngine{vaultClient: vaultClient, path: path}, nil
}

func NewKVWithPath(path string) (KV, error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &kvEngine{vaultClient: vaultClient, path: path}, nil
}
