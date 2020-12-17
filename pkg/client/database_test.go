// +build integration

package client_test

import (
	"github.com/hashicorp/vault/api"
	. "github.com/jasoet/vault-client/pkg/client"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type databaseTestCtx struct {
	vaultClient *api.Client
}

func (ctx *databaseTestCtx) setup(t *testing.T) {
	config := &api.Config{Address: os.Getenv("TEST_VAULT_ADDR")}
	client, err := api.NewClient(config)
	assert.Nil(t, err)

	client.SetToken(os.Getenv("TEST_VAULT_TOKEN"))

	ctx.vaultClient = client
}

func TestDatabase(t *testing.T) {
	ctx := new(databaseTestCtx)
	ctx.setup(t)

	databaseEnginePath := "db-path"
	database, err := NewDatabase(ctx.vaultClient, databaseEnginePath)
	assert.Nil(t, err)
	assert.NotNil(t, database)

	_ = database.Enable()

	t.Run("status should return correct result", func(t *testing.T) {
		status, err := database.Status()
		assert.Nil(t, err)
		assert.NotNil(t, status)
	})

	roleName := "test-database"
	connectionName := "test-db"

	databaseConfig := DatabaseConfig{
		Type:          MySQL,
		ConnectionUrl: "{{username}}:{{password}}@tcp(db:3306)/",
		Username:      "root",
		Password:      "localhost",
		AllowedRoles:  []string{roleName},
	}

	t.Run("create config should success", func(t *testing.T) {
		err = database.CreateConnection(connectionName, databaseConfig)
		assert.Nil(t, err)
	})

	t.Run("read config should produce correct result", func(t *testing.T) {
		config, err := database.ReadConnection(connectionName)

		assert.Nil(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, databaseConfig.Type, config.Type)
		assert.Equal(t, databaseConfig.ConnectionUrl, config.ConnectionUrl)
		assert.Equal(t, databaseConfig.Username, config.Username)
	})

	t.Run("reset config should not produce error", func(t *testing.T) {
		err := database.ResetConnection(connectionName)
		assert.Nil(t, err)
	})

	secondConnectionName := "test-dbx"
	t.Run("list config should produce non-zero result", func(t *testing.T) {
		err = database.CreateConnection(secondConnectionName, databaseConfig)
		assert.Nil(t, err)

		configs, err := database.ListConnection()
		assert.Nil(t, err)
		assert.NotNil(t, configs)
		assert.NotEmpty(t, configs)
		assert.Contains(t, configs, connectionName)
		assert.Contains(t, configs, secondConnectionName)
	})

	t.Run("should not able to fetch deleted connection", func(t *testing.T) {
		err = database.DeleteConnection(secondConnectionName)
		assert.Nil(t, err)

		config, err := database.ReadConnection(secondConnectionName)
		assert.NotNil(t, err)
		assert.Nil(t, config)
	})

	roleConfig := DatabaseRole{
		DatabaseName:         connectionName,
		DefaultTtl:           60,
		MaxTtl:               600,
		CreationStatements:   []string{"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';", "GRANT SELECT ON *.* TO '{{name}}'@'%';"},
		RevocationStatements: []string{"DROP USER '{{name}}'@'%';"},
	}

	t.Run("create role should not produce error", func(t *testing.T) {
		err = database.CreateRole(roleName, roleConfig)
		assert.Nil(t, err)
	})

	t.Run("fetch role should return correct values", func(t *testing.T) {
		detail, err := database.ReadRole(roleName)
		assert.Nil(t, err)
		assert.NotNil(t, detail)
		assert.Equal(t, roleConfig.DatabaseName, detail.DatabaseName)
		assert.Equal(t, roleConfig.DefaultTtl, detail.DefaultTtl)
		assert.Equal(t, roleConfig.MaxTtl, detail.MaxTtl)
	})

	roleNameSecond := "test-database-rolex"
	t.Run("list role should return non-zero result", func(t *testing.T) {
		err = database.CreateRole(roleNameSecond, roleConfig)
		assert.Nil(t, err)

		list, err := database.ListRole()
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.NotEmpty(t, list)
		assert.Contains(t, list, roleName)
		assert.Contains(t, list, roleNameSecond)
	})

	t.Run("cannot fetch deleted role", func(t *testing.T) {
		err := database.DeleteRole(roleNameSecond)
		assert.Nil(t, err)

		detail, err := database.ReadRole(roleNameSecond)
		assert.NotNil(t, err)
		assert.Nil(t, detail)
	})

	t.Run("generate creds on valid role should success", func(t *testing.T) {
		cred, err := database.GenerateCreds(roleName)
		assert.Nil(t, err)
		assert.NotNil(t, cred)
		assert.NotNil(t, cred.Username)
		assert.NotNil(t, cred.Password)
		assert.NotNil(t, cred.LeaseDuration)
		assert.NotNil(t, cred.LeaseId)
	})

	t.Run("list lease should return non-zero result", func(t *testing.T) {
		cred, err := database.GenerateCreds(roleName)
		assert.Nil(t, err)
		assert.NotNil(t, cred)
		assert.NotNil(t, cred.Username)
		assert.NotNil(t, cred.Password)
		assert.NotNil(t, cred.LeaseDuration)
		assert.NotNil(t, cred.LeaseId)

		leases, err := database.ListLease(roleName)
		assert.Nil(t, err)
		assert.NotNil(t, leases)
		assert.NotEmpty(t, leases)
		assert.Contains(t, leases, cred.LeaseId)

	})

}
