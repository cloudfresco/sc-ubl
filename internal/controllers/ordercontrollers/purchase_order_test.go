package ordercontrollers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	orderproto "github.com/cloudfresco/sc-ubl/internal/protogen/order/v1"
	"github.com/cloudfresco/sc-ubl/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreatePurchaseOrderHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	purchaseOrderHeader := orderproto.CreatePurchaseOrderHeaderRequest{}
	purchaseOrderHeader.IssueDate = "02/12/2022"
	purchaseOrderHeader.ValidityPeriod = "02/22/2022"
	purchaseOrderHeader.OrderTypeCode = "ABCFES"
	purchaseOrderHeader.Note = "Sample"
	purchaseOrderHeader.LineExtensionAmount = float64(100)
	purchaseOrderHeader.PayableAmount = float64(100)
	purchaseOrderHeader.TaxExDate = "08/10/2022"
	purchaseOrderHeader.PricingExDate = "08/10/2022"
	purchaseOrderHeader.PaymentExDate = "08/11/2022"

	purchaseOrderLine := orderproto.CreatePurchaseOrderLineRequest{}
	purchaseOrderLine.Note = "Mrs Green agreed to waive charge"
	purchaseOrderLine.LineStatusCode = "ABCFES"
	purchaseOrderLine.OriginatorPartyId = uint32(1)
	purchaseOrderLine.Quantity = float64(100)
	purchaseOrderLine.LineExtensionAmount = float64(100)
	purchaseOrderLine.TotalTaxAmount = float64(17.5)
	purchaseOrderLine.ItemId = uint32(7)
	purchaseOrderLine.PriceAmount = float64(17.5)
	purchaseOrderLine.PriceBaseQuantity = float64(100)
	purchaseOrderLine.PriceValidityPeriodStartDate = "02/22/2022"
	purchaseOrderLine.PriceValidityPeriodEndDate = "02/27/2022"

	purchaseOrderLines := []*orderproto.CreatePurchaseOrderLineRequest{}
	purchaseOrderLines = append(purchaseOrderLines, &purchaseOrderLine)
	purchaseOrderHeader.PurchaseOrderLines = purchaseOrderLines

	data, _ := json.Marshal(&purchaseOrderHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/purchase-orders", bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
		return
	}

	req = common.SetEmailToken(req, tokenString, email)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}
	assert.NotNil(t, w.Body.String())
}

func TestUpdatePurchaseOrderHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	form := orderproto.UpdatePurchaseOrderHeaderRequest{}
	form.OrderTypeCode = "BCFESD"
	form.Note = "sample2"
	form.RequestedInvoiceCurrencyCode = "EUR"
	form.DocumentCurrencyCode = "EUR"
	form.PricingCurrencyCode = "EUR"
	form.TaxCurrencyCode = "EUR"
	form.AccountingCostCode = ""
	form.AccountingCost = "BookingCode001"

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/purchase-orders/413a40b5-5f7b-40c5-bbaf-d6e025543fde", bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
		return
	}

	req = common.SetEmailToken(req, tokenString, email)
	mux.ServeHTTP(w, req)

	resp := w.Result()
	// Check the status code is what we expect.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
		return
	}

	expected := string(`"Updated Successfully"` + "\n")
	assert.Equal(t, w.Body.String(), expected, "they should be equal")
}
