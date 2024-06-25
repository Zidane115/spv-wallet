package main

import (
	"context"
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

func main() {
	const adminxPriv = "xprv..."
	serverURL := "http://localhost:3003"

	xprv, err := xpriv.Generate()

	if err != nil {
		panic(err)
	}

	// try to connect to the spv-wallet
	// get paymail domain from shared config
	//paymailDomain := "example.com"

	admWalletClient := walletclient.NewWithAdminKey(
		serverURL,
		adminxPriv,
	)

	if admWalletClient == nil {
		panic("failed to connect to the wallet")
	}

	respErr := admWalletClient.AdminNewXpub(context.Background(), xprv.XPub().String(), nil)
	if respErr != nil {
		panic(respErr)
	}

	// register new paymail
	// to verify: parameters: xpub, paymail, name, avatar
	_, respErr = admWalletClient.AdminCreatePaymail(context.Background(), xprv.XPub().String(), "", "leader", "")

	if respErr != nil {
		panic(respErr)
	}

	walletClient := walletclient.NewWithXPriv(serverURL, xprv.XPriv())

	if walletClient == nil {
		panic("failed to connect to the wallet")
	}

}
