// +build integration

package client_test

import (
	"github.com/hashicorp/vault/api"
	. "github.com/jasoet/vault-client/pkg/client"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type kvTestCtx struct {
	vaultClient *api.Client
}

func (ctx *kvTestCtx) setup(t *testing.T) {
	config := &api.Config{Address: os.Getenv("TEST_VAULT_ADDR")}
	client, err := api.NewClient(config)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	client.SetToken(os.Getenv("TEST_VAULT_TOKEN"))
	ctx.vaultClient = client
}

func TestKV(t *testing.T) {
	ctx := new(kvTestCtx)
	ctx.setup(t)

	kvEnginePath := "kv-path"
	kv, err := NewKV(ctx.vaultClient, kvEnginePath)
	assert.Nil(t, err)
	assert.NotNil(t, kv)

	_ = kv.Enable()

	t.Run("status should return correct result", func(t *testing.T) {
		status, err := kv.Status()
		assert.Nil(t, err)
		assert.NotNil(t, status)
	})

	config := KVConfig{
		MaxVersions:        20,
		CasRequired:        false,
		DeleteVersionAfter: "10000000s",
	}

	t.Run("write config should return non error", func(t *testing.T) {
		err = kv.WriteConfig(config)
		assert.Nil(t, err)
	})

	t.Run("read config should return correct value", func(t *testing.T) {
		kvConfig, err := kv.ReadConfig()
		assert.Nil(t, err)
		assert.NotNil(t, kvConfig)

		assert.Equal(t, config.MaxVersions, kvConfig.MaxVersions)
		assert.Equal(t, config.CasRequired, kvConfig.CasRequired)
		expected, err := time.ParseDuration(config.DeleteVersionAfter)
		assert.Nil(t, err)
		actual, err := time.ParseDuration(kvConfig.DeleteVersionAfter)
		assert.Nil(t, err)
		assert.Equal(t, expected.Seconds(), actual.Seconds())
	})

	sampleData := DatabaseConfig{
		Type:          MySQL,
		ConnectionUrl: "{{username}}:{{password}}@tcp(db:3306)/",
		Username:      "root",
		Password:      "localhost",
		AllowedRoles:  []string{"Super"},
	}

	dataPath := "sample/first"
	t.Run("write secret data should success", func(t *testing.T) {
		metadata, err := kv.Write(dataPath, sampleData)
		assert.Nil(t, err)
		assert.NotNil(t, metadata)
	})

	t.Run("read secret data should produce correct result", func(t *testing.T) {
		output := new(DatabaseConfig)
		metadata, err := kv.Read(dataPath, output)
		assert.Nil(t, err)
		assert.NotNil(t, metadata)
		assert.NotNil(t, output)
	})

}
