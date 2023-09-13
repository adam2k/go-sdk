//go:build integration
// +build integration

package regtest_test

import (
	"github.com/lightsparkdev/go-sdk/objects"
	"github.com/lightsparkdev/go-sdk/services"
	servicestest "github.com/lightsparkdev/go-sdk/services/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateInvoice(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID, env.ApiClientSecret, &env.ApiClientEndpoint)
	invoice, err := servicestest.CreateInvoiceForNode(client, env.NodeID)
	require.NoError(t, err)
	t.Log(invoice)
}

// TODO: We should add a check for funds on the nodes and automatically add funds if needed.

// Create test invoice from routing node and pay it from node 1.
func TestCreateTestInvoiceNode1(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID, env.ApiClientSecret, &env.ApiClientEndpoint)
	invoice, err := client.CreateTestModeInvoice(env.NodeID, 50_000, nil, nil)
	require.NoError(t, err)
	servicestest.LoadSeedAsSigningKey(t, env.NodeID, env.MasterSeedHex, objects.BitcoinNetworkRegtest, client)
	payment, err := client.PayInvoice(env.NodeID, *invoice, 60, 10_000_000, nil)
	require.NoError(t, err)
	for payment.Status != objects.TransactionStatusSuccess && payment.Status != objects.TransactionStatusFailed {
		entity, err := client.GetEntity(payment.Id)
		require.NoError(t, err)
		castPayment, didCast := (*entity).(objects.OutgoingPayment)
		require.True(t, didCast)
		payment = &castPayment
	}
	if payment.Status == objects.TransactionStatusFailed {
		t.Error("Payment failed")
	}
	t.Log(payment)
}

// Create test invoice from routing node and pay it from node 2.
func TestCreateTestInvoiceNode2(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID2, env.ApiClientSecret2, &env.ApiClientEndpoint)
	invoice, err := client.CreateTestModeInvoice(env.NodeID2, 50_000, nil, nil)
	require.NoError(t, err)
	servicestest.LoadSeedAsSigningKey(t, env.NodeID2, env.MasterSeedHex2, objects.BitcoinNetworkRegtest, client)
	payment, err := client.PayInvoice(env.NodeID2, *invoice, 60, 10_000_000, nil)
	require.NoError(t, err)
	for payment.Status != objects.TransactionStatusSuccess && payment.Status != objects.TransactionStatusFailed {
		entity, err := client.GetEntity(payment.Id)
		require.NoError(t, err)
		castPayment, didCast := (*entity).(objects.OutgoingPayment)
		require.True(t, didCast)
		payment = &castPayment
	}
	if payment.Status == objects.TransactionStatusFailed {
		t.Error("Payment failed")
	}
	t.Log(payment)
}

// Create an invoice from node 1, pay it from node 2
func TestPayInvoiceNode1ToNode2(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID, env.ApiClientSecret, &env.ApiClientEndpoint)
	invoice, err := servicestest.CreateInvoiceForNode(client, env.NodeID)
	require.NoError(t, err)

	t.Log(invoice)
	client_2 := services.NewLightsparkClient(env.ApiClientID2, env.ApiClientSecret2, &env.ApiClientEndpoint)
	servicestest.LoadSeedAsSigningKey(t, env.NodeID2, env.MasterSeedHex2, objects.BitcoinNetworkRegtest, client_2)
	payment, err := client_2.PayInvoice(env.NodeID2, invoice.Data.EncodedPaymentRequest, 60, 1000000, nil)
	require.NoError(t, err)

	for payment.Status != objects.TransactionStatusSuccess && payment.Status != objects.TransactionStatusFailed {
		entity, err := client_2.GetEntity(payment.Id)
		require.NoError(t, err)
		castPayment, didCast := (*entity).(objects.OutgoingPayment)
		require.True(t, didCast)
		payment = &castPayment
	}
	if payment.Status == objects.TransactionStatusFailed {
		t.Errorf("Payment failed: %s", payment.FailureReason.StringValue())
	}
	t.Log(payment)
}

// Create an invoice from node 2, pay it from node 1
func TestPayInvoiceNode2ToNode1(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID2, env.ApiClientSecret2, &env.ApiClientEndpoint)
	invoice, err := servicestest.CreateInvoiceForNode(client, env.NodeID2)
	require.NoError(t, err)

	t.Log(invoice)
	client_2 := services.NewLightsparkClient(env.ApiClientID, env.ApiClientSecret, &env.ApiClientEndpoint)
	servicestest.LoadSeedAsSigningKey(t, env.NodeID, env.MasterSeedHex, objects.BitcoinNetworkRegtest, client_2)
	payment, err := client_2.PayInvoice(env.NodeID, invoice.Data.EncodedPaymentRequest, 60, 1000000, nil)
	require.NoError(t, err)

	for payment.Status != objects.TransactionStatusSuccess && payment.Status != objects.TransactionStatusFailed {
		entity, err := client_2.GetEntity(payment.Id)
		require.NoError(t, err)
		castPayment, didCast := (*entity).(objects.OutgoingPayment)
		require.True(t, didCast)
		payment = &castPayment
	}
	if payment.Status == objects.TransactionStatusFailed {
		t.Errorf("Payment failed: %s", payment.FailureReason.StringValue())
	}
	t.Log(payment)
}

func TestGetChannelUtxos(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID, env.ApiClientSecret, &env.ApiClientEndpoint)
	utxos, err := client.GetNodeChannelUtxos(env.NodeID)
	require.NoError(t, err)
	t.Log(utxos)
}

func TestGetFundingAddress(t *testing.T) {
	env := servicestest.NewConfig()
	client := services.NewLightsparkClient(env.ApiClientID2, env.ApiClientSecret2, &env.ApiClientEndpoint)
	servicestest.LoadSeedAsSigningKey(t, env.NodeID2, env.MasterSeedHex2, objects.BitcoinNetworkRegtest, client)
	address, err := client.CreateNodeWalletAddress(env.NodeID2)
	require.NoError(t, err)
	t.Log(address)
}
