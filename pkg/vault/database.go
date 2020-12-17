package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/jasoet/run-vault/util"
)

type databaseEngine struct {
	vaultClient *api.Client
	path        string
}

type usernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d databaseEngine) GenerateCreds(roleName string) (creds *Creds, err error) {
	result, err := d.vaultClient.Logical().Read(fmt.Sprintf("%v/creds/%v", d.path, roleName))
	if err != nil {
		return
	}

	if result == nil {
		err = fmt.Errorf("%v is not found", roleName)
		return
	}

	userPass := new(usernamePassword)
	err = util.MapToStruct(result.Data, userPass)
	if err != nil {
		return
	}

	creds = new(Creds)
	creds.LeaseId = result.LeaseID
	creds.LeaseDuration = result.LeaseDuration
	creds.Renewable = result.Renewable
	creds.Username = userPass.Username
	creds.Password = userPass.Password

	return
}

func (d databaseEngine) ListConnection() (list []string, err error) {
	result, err := d.vaultClient.Logical().List(fmt.Sprintf("%v/config", d.path))
	if err != nil || result == nil {
		return []string{}, nil
	}

	if val, ok := result.Data["keys"]; ok {
		list = util.ToArrStr(val.([]interface{}))
	}
	return
}

func (d databaseEngine) ListRole() (list []string, err error) {
	result, err := d.vaultClient.Logical().List(fmt.Sprintf("%v/roles", d.path))
	if err != nil || result == nil {
		return []string{}, nil
	}

	if val, ok := result.Data["keys"]; ok {
		list = util.ToArrStr(val.([]interface{}))
	}
	return
}

func (d databaseEngine) CreateConnection(name string, config DatabaseConfig) (err error) {
	_, err = d.vaultClient.Logical().Write(fmt.Sprintf("%v/config/%v", d.path, name), util.StructToMap(config))
	return
}

func (d databaseEngine) ResetConnection(name string) (err error) {
	_, err = d.vaultClient.Logical().Write(fmt.Sprintf("%v/reset/%v", d.path, name), map[string]interface{}{})
	return
}

func (d databaseEngine) DeleteConnection(name string) (err error) {
	_, err = d.vaultClient.Logical().Delete(fmt.Sprintf("%v/config/%v", d.path, name))
	return
}

type configConnectionDetail struct {
	ConnectionUrl string `json:"connection_url"`
	Username      string `json:"username"`
}

type databaseConfigDetail struct {
	Type                   DatabaseType           `json:"plugin_name"`
	ConnectionDetails      configConnectionDetail `json:"connection_details"`
	AllowedRoles           []string               `json:"allowed_roles"`
	RootRotationStatements []string               `json:"root_credentials_rotate_statements"`
	PasswordPolicy         string                 `json:"password_policy"`
}

func (d databaseEngine) ReadConnection(name string) (config *DatabaseConfig, err error) {
	result, err := d.vaultClient.Logical().Read(fmt.Sprintf("%v/config/%v", d.path, name))
	if err != nil {
		return
	}

	if result == nil {
		err = fmt.Errorf("%v is not found", name)
		return
	}

	configDetail := new(databaseConfigDetail)
	err = util.MapToStruct(result.Data, configDetail)
	if err != nil {
		return
	}

	config = new(DatabaseConfig)
	config.Type = configDetail.Type
	config.Username = configDetail.ConnectionDetails.Username
	config.ConnectionUrl = configDetail.ConnectionDetails.ConnectionUrl
	config.AllowedRoles = configDetail.AllowedRoles
	config.PasswordPolicy = configDetail.PasswordPolicy
	config.RootRotationStatements = configDetail.RootRotationStatements
	return
}

func (d databaseEngine) CreateRole(name string, config DatabaseRole) (err error) {
	_, err = d.vaultClient.Logical().Write(fmt.Sprintf("%v/roles/%v", d.path, name), util.StructToMap(config))
	return
}

func (d databaseEngine) DeleteRole(name string) (err error) {
	_, err = d.vaultClient.Logical().Delete(fmt.Sprintf("%v/roles/%v", d.path, name))
	return
}

func (d databaseEngine) ReadRole(name string) (role *DatabaseRole, err error) {
	result, err := d.vaultClient.Logical().Read(fmt.Sprintf("%v/roles/%v", d.path, name))
	if err != nil {
		return
	}

	if result == nil {
		err = fmt.Errorf("%v is not found", name)
		return
	}

	role = new(DatabaseRole)
	err = util.MapToStruct(result.Data, role)
	return
}

func (d databaseEngine) Path() string {
	return d.path
}

func (d databaseEngine) Enable() (err error) {
	data := map[string]interface{}{"type": "database"}
	_, err = d.vaultClient.Logical().Write(fmt.Sprintf("/sys/mounts/%v", d.path), data)
	return
}

func (d databaseEngine) Status() (status *DatabaseSecretStatus, err error) {
	result, err := d.vaultClient.Logical().Read(fmt.Sprintf("/sys/mounts/%v/tune", d.path))
	if err != nil || result == nil {
		return
	}

	status = new(DatabaseSecretStatus)
	err = util.MapToStruct(result.Data, status)

	return
}

func (d databaseEngine) ListLease(roleName string) (list []string, err error) {
	prefix := fmt.Sprintf("%v/creds/%v/", d.path, roleName)
	result, err := d.vaultClient.Logical().List(fmt.Sprintf("/sys/leases/lookup/%v", prefix))
	if err != nil || result == nil {
		return []string{}, nil
	}

	if val, ok := result.Data["keys"]; ok {
		list = util.ToArrStrPrefixPath(val.([]interface{}), prefix)
	}
	return
}

func DefaultDatabase() (Database, error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &databaseEngine{vaultClient: vaultClient, path: "database"}, nil
}

func NewDatabase(vaultClient *api.Client, path string) (Database, error) {
	return &databaseEngine{vaultClient: vaultClient, path: path}, nil
}

func NewDatabaseWithPath(path string) (Database, error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &databaseEngine{vaultClient: vaultClient, path: path}, nil
}
