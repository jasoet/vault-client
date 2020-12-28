package client

type Database interface {
	Path() string
	Enable() error
	Status() (*SecretStatus, error)

	CreateConnection(name string, config DatabaseConfig) error
	ReadConnection(name string) (*DatabaseConfig, error)
	ResetConnection(name string) error
	DeleteConnection(name string) error
	ListConnection() ([]string, error)

	CreateRole(name string, config DatabaseRole) error
	ReadRole(name string) (*DatabaseRole, error)
	DeleteRole(name string) error
	ListRole() ([]string, error)

	GenerateCreds(roleName string) (*Creds, error)
	ListLease(roleName string) ([]string, error)
}

type Lease interface {
	Lookup(leaseId string) (*LeaseDetail, error)
	List(prefix string) ([]string, error)
	Renew(leaseId string, increment int) error
	Revoke(leaseId string) error
	RevokePrefix(prefix string) error
	Tidy() error
}

type KV interface {
	Path() string
	Enable() error
	Status() (*SecretStatus, error)

	WriteConfig(config KVConfig) error
	ReadConfig() (*KVConfig, error)

	Write(path string, data interface{}) (*KVMetadata, error)
	Read(path string, result interface{}) (*KVMetadata, error)
}
