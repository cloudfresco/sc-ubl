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

func TestCreateDebitNoteHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	debitNoteHeader := invoiceproto.CreateDebitNoteHeaderRequest{}
	debitNoteHeader.IssueDate = "11/15/2022"
	debitNoteHeader.TaxPointDate = "11/25/2022"
	debitNoteHeader.Note = "Ordered in our booth at the convention"
	debitNoteHeader.DocumentCurrencyCode = "EUR"
	debitNoteHeader.AccountingCost = "Project cost code 123"
	debitNoteHeader.InvoicePeriodStartDate = "11/15/2022"
	debitNoteHeader.InvoicePeriodEndDate = "12/01/2022"
	debitNoteHeader.TaxExDate = "11/25/2022"
	debitNoteHeader.PricingExDate = "11/25/2022"
	debitNoteHeader.PaymentExDate = "11/25/2022"
	debitNoteHeader.PaymentAltExDate = "11/25/2022"
	debitNoteHeader.ChargeTotalAmount = float64(0)
	debitNoteHeader.PrepaidAmount = float64(0)
	debitNoteHeader.PayableRoundingAmount = float64(0)
	debitNoteHeader.PayableAmount = float64(1250)

	debitNoteLine := invoiceproto.CreateDebitNoteLineRequest{}
	debitNoteLine.Note = "Scratch on box"
	debitNoteLine.DebitedQuantity = float64(1)
	debitNoteLine.LineExtensionAmount = float64(1250)
	debitNoteLine.TaxPointDate = "11/25/2022"
	debitNoteLine.AccountingCost = "BookingCode002"
	debitNoteLine.ItemId = uint32(7)
	debitNoteLine.PriceValidityPeriodStartDate = "11/15/2022"
	debitNoteLine.PriceValidityPeriodEndDate = "12/01/2022"

	debitNoteLines := []*invoiceproto.CreateDebitNoteLineRequest{}
	debitNoteLines = append(debitNoteLines, &debitNoteLine)
	debitNoteHeader.DebitNoteLines = debitNoteLines

	data, _ := json.Marshal(&debitNoteHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/debit-notes", bytes.NewBuffer(data))
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

func TestUpdateDebitNoteHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	form := invoiceproto.UpdateDebitNoteHeaderRequest{}
	form.Note = "Ordered"
	form.DocumentCurrencyCode = "EUR"
	form.ChargeTotalAmount = float64(100)
	form.PrepaidAmount = float64(60)
	form.PayableRoundingAmount = float64(100)
	form.PayableAmount = float64(300)

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/debit-notes/930f8806-db24-4562-b8c7-72df75518355", bytes.NewBuffer(data))
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
