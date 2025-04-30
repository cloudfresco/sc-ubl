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

func TestCreateDespatchHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	despatchHeader := logisticsproto.CreateDespatchHeaderRequest{}
	despatchHeader.IssueDate = "06/20/2022"
	despatchHeader.DocumentStatusCode = "NoStatus"
	despatchHeader.DespatchAdviceTypeCode = "delivery"
	despatchHeader.Note = "sample"
	despatchHeader.OrderId = uint32(1)

	despatchLine := logisticsproto.CreateDespatchLineRequest{}
	despatchLine.Note = "Mrs Green agreed to waive charge"
	despatchLine.LineStatusCode = "NoStatus"
	despatchLine.DeliveredQuantity = float64(2)
	despatchLine.BackorderQuantity = float64(1)
	despatchLine.ItemId = uint32(7)

	despatchLines := []*logisticsproto.CreateDespatchLineRequest{}
	despatchLines = append(despatchLines, &despatchLine)
	despatchHeader.DespatchLines = despatchLines

	data, _ := json.Marshal(&despatchHeader)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/despatches", bytes.NewBuffer(data))
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

func TestUpdateDespatchHeader(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	form := logisticsproto.UpdateDespatchHeaderRequest{}
	form.DocumentStatusCode = "Status"
	form.DespatchAdviceTypeCode = "delivery"
	form.Note = "sample1"

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&form)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/despatches/234fd566-9451-4e3e-8318-4b713e688960", bytes.NewBuffer(data))
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
