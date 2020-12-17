package client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/jasoet/vault-client/pkg/util"
)

type leaseEngine struct {
	vaultClient *api.Client
}

func (l leaseEngine) Lookup(leaseId string) (detail *LeaseDetail, err error) {
	payload := map[string]interface{}{
		"lease_id": leaseId,
	}
	result, err := l.vaultClient.Logical().Write("/sys/leases/lookup", payload)
	if err != nil {
		return
	}

	if result == nil {
		err = fmt.Errorf("%v is not found", leaseId)
		return
	}

	detail = new(LeaseDetail)
	err = util.MapToStruct(result.Data, detail)
	return
}

func (l leaseEngine) List(prefix string) (list []string, err error) {
	result, err := l.vaultClient.Logical().List(fmt.Sprintf("/sys/leases/lookup/%v", prefix))
	if err != nil || result == nil {
		return []string{}, nil
	}

	if val, ok := result.Data["keys"]; ok {
		list = util.ToArrStrPrefixPath(val.([]interface{}), prefix)
	}
	return
}

func (l leaseEngine) Renew(leaseId string, increment int) (err error) {
	payload := map[string]interface{}{
		"lease_id":  leaseId,
		"increment": increment,
	}
	_, err = l.vaultClient.Logical().Write("/sys/leases/renew", payload)
	return
}

func (l leaseEngine) Revoke(leaseId string) (err error) {
	payload := map[string]interface{}{
		"lease_id": leaseId,
	}
	_, err = l.vaultClient.Logical().Write("/sys/leases/revoke", payload)
	return
}

func (l leaseEngine) RevokePrefix(prefix string) (err error) {
	_, err = l.vaultClient.Logical().Write(fmt.Sprintf("/sys/leases/revoke-prefix/%v", prefix), map[string]interface{}{})
	return
}

func (l leaseEngine) Tidy() (err error) {
	_, err = l.vaultClient.Logical().Write("/sys/leases/tidy", map[string]interface{}{})
	return
}

func DefaultLease() (lease Lease, err error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return
	}

	lease = &leaseEngine{vaultClient: vaultClient}
	return
}

func NewLease(vaultClient *api.Client) (lease Lease, err error) {
	lease = &leaseEngine{vaultClient: vaultClient}
	return
}
