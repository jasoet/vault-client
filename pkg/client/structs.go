package client

import "time"

type DatabaseType string

const (
	MySQL      DatabaseType = "mysql-database-plugin"
	PostgreSQL DatabaseType = "postgresql-database-plugin"
)

type SecretStatus struct {
	DefaultLeaseTtl int    `json:"default_lease_ttl"`
	MaxLeaseTtl     int    `json:"max_least_ttl"`
	Description     string `json:"description"`
	ForceNoCache    bool   `json:"force_no_cache"`
}

type DatabaseConfig struct {
	Type                   DatabaseType `json:"plugin_name"`
	ConnectionUrl          string       `json:"connection_url"`
	Username               string       `json:"username"`
	Password               string       `json:"password"`
	AllowedRoles           []string     `json:"allowed_roles"`
	RootRotationStatements []string     `json:"root_rotation_statements,omitempty"`
	PasswordPolicy         string       `json:"password_policy,omitempty"`
}

type DatabaseRole struct {
	ConnectionName       string   `json:"db_name"`
	DefaultTtl           int      `json:"default_ttl"`
	MaxTtl               int      `json:"max_ttl"`
	CreationStatements   []string `json:"creation_statements"`
	RevocationStatements []string `json:"revocation_statements"`
	RollbackStatements   []string `json:"rollback_statements,omitempty"`
	RenewStatements      []string `json:"renew_statements,omitempty"`
}

type Creds struct {
	LeaseId       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type LeaseDetail struct {
	LeaseId         string    `json:"id"`
	IssueTime       time.Time `json:"issue_time"`
	ExpiredTime     time.Time `json:"expired_time"`
	LastRenewalTime time.Time `json:"last_renewal_time"`
	Renewable       bool      `json:"renewable"`
	Ttl             int       `json:"ttl"`
}

type KVConfig struct {
	MaxVersions        int    `json:"max_versions"`
	CasRequired        bool   `json:"cas_required"`
	DeleteVersionAfter string `json:"delete_version_after"` //use go duration format https://golang.org/pkg/time/#ParseDuration
}

type KVMetadata struct {
	Version      int       `json:"version"`
	Destroyed    bool      `json:"destroyed"`
	DeletionTime time.Time `json:"deletion_time"`
	CreatedTime  time.Time `json:"created_time"`
}

type KVHistoryMetadata struct {
	CreatedTime    time.Time             `json:"created_time"`
	CurrentVersion int                   `json:"current_version"`
	MaxVersion     int                   `json:"max_version"`
	OldestVersion  int                   `json:"oldest_version"`
	UpdatedTime    time.Time             `json:"updated_time"`
	Versions       map[string]KVMetadata `json:"versions"`
}
