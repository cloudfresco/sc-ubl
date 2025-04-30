package taxcontrollers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfresco/sc-ubl/internal/common"
	taxproto "github.com/cloudfresco/sc-ubl/internal/protogen/tax/v1"
	"github.com/cloudfresco/sc-ubl/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaxScheme(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	taxScheme := taxproto.CreateTaxSchemeRequest{}
	taxScheme.TsId = "UK VAT"
	taxScheme.TaxSchemeName = "T"
	taxScheme.TaxTypeCode = "VAT"
	taxScheme.CurrencyCode = "EUR"

	data, _ := json.Marshal(&taxScheme)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/tax-schemes", bytes.NewBuffer(data))
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

func TestUpdateTaxScheme(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	taxScheme := taxproto.UpdateTaxSchemeRequest{}
	taxScheme.TaxSchemeName = "TaxSchemeName"
	taxScheme.TaxTypeCode = "VAT"
	taxScheme.CurrencyCode = "EUR"

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&taxScheme)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/tax-schemes/15fd8632-8674-462d-a08a-de1560d2d9e8", bytes.NewBuffer(data))
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

func TestCreateTaxCategory(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	taxCategory := taxproto.CreateTaxCategoryRequest{}
	taxCategory.TcId = ""
	taxCategory.TaxCategoryName = "TaxCategory"
	taxCategory.Percent = float32(20)
	taxCategory.BaseUnitMeasure = "EUR"
	taxCategory.PerUnitAmount = float64(10)
	taxCategory.TaxExemptionReasonCode = ""
	taxCategory.TaxExemptionReason = ""
	taxCategory.TierRange = ""
	taxCategory.TierRatePercent = float32(10)

	data, _ := json.Marshal(&taxCategory)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/tax-schemes/15fd8632-8674-462d-a08a-de1560d2d9e8/add-tax-category", bytes.NewBuffer(data))
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

func TestUpdateTaxCategory(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	taxCategory := taxproto.UpdateTaxCategoryRequest{}
	taxCategory.TaxCategoryName = "TaxCategory1"
	taxCategory.Percent = float32(10)
	taxCategory.BaseUnitMeasure = "EUR"
	taxCategory.PerUnitAmount = float64(10)
	taxCategory.TaxExemptionReasonCode = ""
	taxCategory.TaxExemptionReason = ""

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&taxCategory)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/tax-schemes/tax-categories/7036e24c-0ec5-48ac-bc97-4c2cbe75a54c", bytes.NewBuffer(data))
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

func TestCreateTaxTotal(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	taxTotal := taxproto.CreateTaxTotalRequest{}
	taxTotal.TaxAmount = float64(17.5)
	taxTotal.RoundingAmount = float64(18)
	taxTotal.TaxEvidenceIndicator = false
	taxTotal.TaxIncludedIndicator = false
	taxTotal.MasterFlag = "CNL"
	taxTotal.MasterId = uint32(1)

	data, _ := json.Marshal(&taxTotal)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/tax-schemes/tax-categories/7036e24c-0ec5-48ac-bc97-4c2cbe75a54c/add-tax-total", bytes.NewBuffer(data))
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

func TestUpdateTaxTotal(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	taxTotal := taxproto.UpdateTaxTotalRequest{}
	taxTotal.TaxAmount = float64(18.5)
	taxTotal.RoundingAmount = float64(19)

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&taxTotal)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/tax-schemes/tax-totals/a62007af-a514-45a5-967d-270ad8c34f91", bytes.NewBuffer(data))
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

func TestCreateTaxSubTotal(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	w := httptest.NewRecorder()

	taxSubTotal := taxproto.CreateTaxSubTotalRequest{}
	taxSubTotal.TaxableAmount = float64(100)
	taxSubTotal.TaxAmount = float64(17.5)
	taxSubTotal.CalculationSequenceNumeric = uint32(0)
	taxSubTotal.TransactionCurrencyTaxAmount = float64(0)
	taxSubTotal.Percent = float32(10)
	taxSubTotal.BaseUnitMeasure = "EUR"
	taxSubTotal.PerUnitAmount = float64(10)
	taxSubTotal.TierRange = ""
	taxSubTotal.TierRatePercent = float64(10)
	taxSubTotal.TaxCategoryId = uint32(1)

	data, _ := json.Marshal(&taxSubTotal)

	req, err := http.NewRequest("POST", backendServerAddr+"/v2.3/tax-schemes/tax-totals/ddee759b-936d-4935-a60c-b11c2109ffe9/add-tax-subtotal", bytes.NewBuffer(data))
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

func TestUpdateTaxSubTotal(t *testing.T) {
	err := test.LoadSQL(log, dbService)
	if err != nil {
		t.Error(err)
		return
	}

	tokenString, email, backendServerAddr := LoginUser()

	taxSubTotal := taxproto.UpdateTaxSubTotalRequest{}
	taxSubTotal.TaxableAmount = float64(200)
	taxSubTotal.TaxAmount = float64(18.5)
	taxSubTotal.CalculationSequenceNumeric = uint32(0)
	taxSubTotal.TransactionCurrencyTaxAmount = float64(0)
	taxSubTotal.Percent = float32(15)
	taxSubTotal.BaseUnitMeasure = "EUR"
	taxSubTotal.PerUnitAmount = float64(5)

	w := httptest.NewRecorder()

	data, _ := json.Marshal(&taxSubTotal)

	req, err := http.NewRequest("PUT", backendServerAddr+"/v2.3/tax-schemes/tax-subtotals/f7da0aaa-8068-495c-9176-b7681e651226", bytes.NewBuffer(data))
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
