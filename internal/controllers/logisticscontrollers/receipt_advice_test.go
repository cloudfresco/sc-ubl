package logisticscontrollers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	logisticsproto "github.com/cloudfresco/sc-ubl/internal/protogen/logistics/v1"
	"github.com/cloudfresco/sc-ubl/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateReceiptAdviceHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	receiptAdviceHeader := logisticsproto.CreateReceiptAdviceHeaderRequest{}
	receiptAdviceHeader.IssueDate = "06/20/2022"
	receiptAdviceHeader.ReceiptAdviceTypeCode = "ABCFES"
	receiptAdviceHeader.Note = "sample"
	receiptAdviceHeader.OrderId = uint32(1)
	receiptAdviceHeader.DespatchId = uint32(1)

	receiptAdviceLine := logisticsproto.CreateReceiptAdviceLineRequest{}
	receiptAdviceLine.Note = "Mrs Green agreed to waive charge"
	receiptAdviceLine.ReceivedQuantity = uint32(2)
	receiptAdviceLine.ShortQuantity = uint32(1)
	receiptAdviceLine.ReceivedDate = "06/25/2022"
	receiptAdviceLine.OrderLineId = uint32(1)
	receiptAdviceLine.DespatchLineId = uint32(1)
	receiptAdviceLine.ItemId = uint32(7)
	receiptAdviceLine.ShipmentId = uint32(1)

	receiptAdviceLines := []*logisticsproto.CreateReceiptAdviceLineRequest{}
	receiptAdviceLines = append(receiptAdviceLines, &receiptAdviceLine)
	receiptAdviceHeader.ReceiptAdviceLines = receiptAdviceLines

	data, _ := json.Marshal(&receiptAdviceHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/receipt-advices", bytes.NewBuffer(data))
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

func TestUpdateReceiptAdviceHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	form := logisticsproto.UpdateReceiptAdviceHeaderRequest{}
	form.ReceiptAdviceTypeCode = "DFESAB"
	form.Note = "sample1"
	form.LineCountNumeric = uint32(10)

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/receipt-advices/234fd566-9451-4e3e-8318-4b713e688960", bytes.NewBuffer(data))
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
