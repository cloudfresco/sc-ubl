package invoicecontrollers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	invoiceproto "github.com/cloudfresco/sc-ubl/internal/protogen/invoice/v1"
	"github.com/cloudfresco/sc-ubl/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvoice(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	invoiceHeader := invoiceproto.CreateInvoiceRequest{}
	invoiceHeader.IssueDate = "06/21/2022"
	invoiceHeader.DueDate = "06/25/2022"
	invoiceHeader.TaxPointDate = "07/25/2022"
	invoiceHeader.InvoiceTypeCode = "SalesInvoice"
	invoiceHeader.Note = "sample"
	invoiceHeader.InvoicePeriodStartDate = "06/21/2022"
	invoiceHeader.InvoicePeriodEndDate = "07/01/2022"
	invoiceHeader.TaxExDate = "06/21/2022"
	invoiceHeader.PricingExDate = "06/22/2022"
	invoiceHeader.PaymentExDate = "06/23/2022"
	invoiceHeader.PaymentAltExDate = "06/24/2022"
	invoiceHeader.ChargeTotalAmount = float64(0)
	invoiceHeader.PrepaidAmount = float64(0)
	invoiceHeader.PayableRoundingAmount = float64(0)
	invoiceHeader.PayableAmount = float64(200)

	invoiceLine := invoiceproto.CreateInvoiceLineRequest{}
	invoiceLine.Note = "Ordered in our booth at the convention."
	invoiceLine.InvoicedQuantity = float64(1)
	invoiceLine.LineExtensionAmount = float64(200)
	invoiceLine.TaxPointDate = "11/25/2022"
	invoiceLine.AccountingCost = "Code002"
	invoiceLine.ItemId = uint32(7)
	invoiceLine.InvoicePeriodStartDate = "06/21/2022"
	invoiceLine.InvoicePeriodEndDate = "07/01/2022"
	invoiceLine.PriceValidityPeriodStartDate = "07/01/2022"
	invoiceLine.PriceValidityPeriodEndDate = "08/01/2022"

	invoiceLines := []*invoiceproto.CreateInvoiceLineRequest{}
	invoiceLines = append(invoiceLines, &invoiceLine)
	invoiceHeader.InvoiceLines = invoiceLines

	data, _ := json.Marshal(&invoiceHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/invoices", bytes.NewBuffer(data))
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

func TestUpdateInvoice(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	form := invoiceproto.UpdateInvoiceRequest{}
	form.Note = "Ordered"
	form.InvoiceTypeCode = "Sales"
	form.ChargeTotalAmount = float64(200)
	form.PrepaidAmount = float64(100)
	form.PayableRoundingAmount = float64(150)
	form.PayableAmount = float64(400)

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/invoices/834dd586-a1b8-4e3f-89cf-bfd25edd9a60", bytes.NewBuffer(data))
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
