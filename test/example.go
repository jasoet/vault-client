package main

import (
	"fmt"
	"github.com/jasoet/vault-client/pkg/client"
	"os"
)

func main() {
	_ = os.Setenv("VAULT_ADDR", "http://127.0.0.1:18200")
	_ = os.Setenv("VAULT_TOKEN", "localhost")

	database, err := client.DefaultDatabase()
	if err != nil {
		panic(err)
	}

	status, err := database.Status()
	fmt.Printf("StatusResult: %#v, Error: %#v\n", status, err)

	err = database.Enable()

	roleName := "vault-database"
	connectionName := "vault-default"
	databaseConfig := client.DatabaseConfig{
		Type:          client.MySQL,
		ConnectionUrl: "{{username}}:{{password}}@tcp(db:3306)/",
		Username:      "root",
		Password:      "localhost",
		AllowedRoles:  []string{roleName},
	}

	err = database.CreateConnection(connectionName, databaseConfig)

	if err != nil {
		panic(err)
	}

	config, err := database.ReadConnection(connectionName)
	fmt.Printf("ReadConnection: %#v\n", config)

	roleConfig := client.DatabaseRole{
		ConnectionName:       connectionName,
		DefaultTtl:           60,
		MaxTtl:               600,
		CreationStatements:   []string{"CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';", "GRANT SELECT ON *.* TO '{{name}}'@'%';"},
		RevocationStatements: []string{"DROP USER '{{name}}'@'%';"},
	}

	err = database.CreateRole(roleName, roleConfig)
	fmt.Printf("error: %#v\n", err)

	role, err := database.ReadRole(roleName)
	fmt.Printf("Role: %#v, error: %#v\n", role, err)

	configs, err := database.ListConnection()
	fmt.Printf("configs: %#v, error: %#v\n", configs, err)

	roles, err := database.ListRole()
	fmt.Printf("Roles: %#v, error: %#v\n", roles, err)

	creds, err := database.GenerateCreds(roleName)
	fmt.Printf("Creds: %#v, error: %#v\n", creds, err)

	creds, err = database.GenerateCreds(roleName)
	fmt.Printf("Creds: %#v, error: %#v\n", creds, err)

	creds, err = database.GenerateCreds(roleName)
	fmt.Printf("Creds: %#v, error: %#v\n", creds, err)

	leases, err := database.ListLease(roleName)
	fmt.Printf("Leases: %#v, error: %#v\n", leases, err)

	lease, err := client.DefaultLease()
	if err != nil {
		panic(err)
	}
	listL, err := lease.List(fmt.Sprintf("database/%v", roleName))
	if err != nil {
		panic(err)
	}
	fmt.Printf("LeaseLL: %#v, error: %#v\n", listL, err)

}
