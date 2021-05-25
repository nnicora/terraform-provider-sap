package sap

import (
	"github.com/nnicora/sap-sdk-go/sap/session"
	"github.com/nnicora/sap-sdk-go/service/btpaccounts"
	"github.com/nnicora/sap-sdk-go/service/btpentitlements"
)

type SAPClient struct {
	session                 *session.RuntimeSession
	btpAccountsV1Client     *btpaccounts.AccountsV1
	btpEntitlementsV1Client *btpentitlements.EntitlementsV1
	//btpProvisioningV1Client     *btpprovisioning.ProvisioningV1
	//btpSaaSProvisioningV1Client *btpsaasprovisioning.SaaSProvisioningV1
}
