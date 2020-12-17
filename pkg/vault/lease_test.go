// +build integration

package vault_test

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	. "github.com/jasoet/run-vault/pkg/vault"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type leaseTestCtx struct {
	vaultClient *api.Client
}

func (ctx *leaseTestCtx) setup(t *testing.T) {
	config := &api.Config{Address: os.Getenv("TEST_VAULT_ADDR")}
	client, err := api.NewClient(config)
	assert.Nil(t, err)

	client.SetToken(os.Getenv("TEST_VAULT_TOKEN"))

	ctx.vaultClient = client
}

func TestLease(t *testing.T) {
	ctx := new(leaseTestCtx)
	ctx.setup(t)

	engine, err := NewLease(ctx.vaultClient)
	assert.Nil(t, err)
	assert.NotNil(t, engine)

	// Setup Database engine to generate creds
	databaseEnginePath := "lease-path"
	database, err := NewDatabase(ctx.vaultClient, databaseEnginePath)
	assert.Nil(t, err)
	assert.NotNil(t, database)

	_ = database.Enable()

	roleName := "lease-role"
	connectionName := "lease-test"
	databaseConfig := DatabaseConfig{
		Type:          MySQL,
		ConnectionUrl: "{{username}}:{{password}}@tcp(db:3306)/",
		Username:      "root",
		Password:      "localhost",
		AllowedRoles:  []string{roleName},
	}

	err = database.CreateConnection(connectionName, databaseConfig)
	assert.Nil(t, err)

	roleConfig := DatabaseRole{
		DatabaseName:         connectionName,
		DefaultTtl:           120,
		MaxTtl:               600,
		CreationStatements:   []string{"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';", "GRANT SELECT ON *.* TO '{{name}}'@'%';"},
		RevocationStatements: []string{"DROP USER '{{name}}'@'%';"},
	}

	err = database.CreateRole(roleName, roleConfig)
	assert.Nil(t, err)

	// Database prefix, used to check leases
	databaseLeasePrefix := fmt.Sprintf("%v/creds/%v/", databaseEnginePath, roleName)

	// Generate 3 Database Creds for test
	creds, err := database.GenerateCreds(roleName)
	assert.Nil(t, err)
	assert.NotNil(t, creds)

	credsTwo, err := database.GenerateCreds(roleName)
	assert.Nil(t, err)
	assert.NotNil(t, credsTwo)

	credsThree, err := database.GenerateCreds(roleName)
	assert.Nil(t, err)
	assert.NotNil(t, credsThree)

	t.Run("lookup should return non nil result", func(t *testing.T) {
		detail, err := engine.Lookup(creds.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)

		detail, err = engine.Lookup(credsTwo.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)
	})

	t.Run("list should return result at lease 3", func(t *testing.T) {
		list, err := engine.List(databaseLeasePrefix)
		assert.Nil(t, err)
		assert.NotNil(t, list)
		assert.GreaterOrEqual(t, len(list), 3)
	})

	t.Run("creds lease ttl should greater than previous after renew", func(t *testing.T) {
		// Ttl on Database Creds stored as LeaseDuration
		detail, err := engine.Lookup(credsTwo.LeaseId)
		assert.Nil(t, err)
		previousTtl := detail.Ttl

		err = engine.Renew(credsTwo.LeaseId, 1000)
		assert.Nil(t, err)

		detail, err = engine.Lookup(credsTwo.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)
		assert.GreaterOrEqual(t, detail.Ttl, previousTtl)
	})

	t.Run("looking up revoked creds should return err", func(t *testing.T) {
		// Make sure lease still valid
		detail, err := engine.Lookup(credsTwo.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)

		err = engine.Revoke(credsTwo.LeaseId)
		assert.Nil(t, err)

		// lease should invalid
		detail, err = engine.Lookup(credsTwo.LeaseId)
		assert.NotNil(t, err)
		assert.Nil(t, detail)
	})

	t.Run("all lease should invalid after revoke-prefix", func(t *testing.T) {
		// Make sure creds and credsThree lease still valid
		detail, err := engine.Lookup(creds.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)

		detail, err = engine.Lookup(credsThree.LeaseId)
		assert.Nil(t, err)
		assert.NotNil(t, detail)

		err = engine.RevokePrefix(databaseLeasePrefix)
		assert.Nil(t, err)

		// Make sure all creds invalid
		detail, err = engine.Lookup(creds.LeaseId)
		assert.NotNil(t, err)
		assert.Nil(t, detail)

		detail, err = engine.Lookup(credsTwo.LeaseId)
		assert.NotNil(t, err)
		assert.Nil(t, detail)

		detail, err = engine.Lookup(credsThree.LeaseId)
		assert.NotNil(t, err)
		assert.Nil(t, detail)

	})

	t.Run("tidy should not return err", func(t *testing.T) {
		err = engine.Tidy()
		assert.Nil(t, err)

	})

}
