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

func TestCreateCreditNoteHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	creditNoteHeader := invoiceproto.CreateCreditNoteHeaderRequest{}
	creditNoteHeader.IssueDate = "06/01/2022"
	creditNoteHeader.DueDate = "08/21/2022"
	creditNoteHeader.TaxPointDate = "08/01/2022"
	creditNoteHeader.Note = "Ordered in our booth at the convention"
	creditNoteHeader.DocumentCurrencyCode = ""
	creditNoteHeader.TaxCurrencyCode = "GBP"
	creditNoteHeader.InvoicePeriodStartDate = "08/01/2022"
	creditNoteHeader.InvoicePeriodEndDate = "09/01/2022"
	creditNoteHeader.TaxExDate = "08/10/2022"
	creditNoteHeader.PricingExDate = "08/10/2022"
	creditNoteHeader.PaymentExDate = "08/11/2022"
	creditNoteHeader.PaymentAltExDate = "08/10/2022"
	creditNoteHeader.ChargeTotalAmount = float64(0)
	creditNoteHeader.PrepaidAmount = float64(0)
	creditNoteHeader.PayableRoundingAmount = float64(0)
	creditNoteHeader.PayableAmount = float64(20)

	creditNoteLine := invoiceproto.CreateCreditNoteLineRequest{}
	creditNoteLine.Note = "Ordered in our booth at the convention"
	creditNoteLine.CreditedQuantity = float64(2)
	creditNoteLine.LineExtensionAmount = float64(20)
	creditNoteLine.TaxPointDate = "08/01/2022"
	creditNoteLine.InvoicePeriodStartDate = "08/01/2022"
	creditNoteLine.InvoicePeriodEndDate = "09/01/2022"
	creditNoteLine.ItemId = uint32(7)
	creditNoteLine.PriceValidityPeriodStartDate = "09/01/2022"
	creditNoteLine.PriceValidityPeriodEndDate = "09/02/2022"

	creditNoteLines := []*invoiceproto.CreateCreditNoteLineRequest{}
	creditNoteLines = append(creditNoteLines, &creditNoteLine)
	creditNoteHeader.CreditNoteLines = creditNoteLines

	data, _ := json.Marshal(&creditNoteHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/credit-notes", bytes.NewBuffer(data))
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

func TestUpdateCreditNoteHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	form := invoiceproto.UpdateCreditNoteHeaderRequest{}
	form.Note = "Ordered"
	form.TaxCurrencyCode = "GBP"
	form.ChargeTotalAmount = float64(10)
	form.PrepaidAmount = float64(20)
	form.PayableRoundingAmount = float64(100)
	form.PayableAmount = float64(200)

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/credit-notes/58b97fc8-b88e-4821-92f8-91265c4dc2fc", bytes.NewBuffer(data))
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
