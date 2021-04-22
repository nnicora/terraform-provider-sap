package sap

import (
	"github.com/nnicora/sap-sdk-go/sap/session"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
	"github.com/nnicora/sap-sdk-go/service/btpprovisioning"
)

type SAPClient struct {
	session                 *session.RuntimeSession
	btpAccountsV1Client     *btpaccounts.AccountsV1
	btpEntitlementsV1Client *btpentitlements.EntitlementsV1
	btpProvisioningV1Client *btpprovisioning.ProvisioningV1
}
