package main

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestPositiveCaseTabung(t *testing.T) {
    client := resty.New()

    url := "http://localhost:8100/api/v1/transaction/tabung"
    reqBody := map[string]float64{
        "nominal": 11777.98,
    }

    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "563582419.111111").
        SetBody(reqBody).
        Post(url)

    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode(), "Expected status 200 OK")
    t.Logf("Response Status: %s", resp.Status())
    t.Logf("Response Body: %s", resp.String())
}

func TestNegativeCaseTabung(t *testing.T) {
    client := resty.New()

    url := "http://localhost:8100/api/v1/transaction/tabung"
    reqBody := map[string]float64{
        "nominal": 11777.98,
    }

    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "643007129.111111").
        SetBody(reqBody).
        Post(url)

    assert.NoError(t, err)
    assert.Equal(t, 404, resp.StatusCode(), "Expected status 404 Not Found")
    t.Logf("Response Status: %s", resp.Status())
    t.Logf("Response Body: %s", resp.String())
}

func TestPositiveCaseTarik(t *testing.T) {
    client := resty.New()

    url := "http://localhost:8100/api/v1/transaction/tarik"
    reqBody := map[string]float64{
        "nominal": 100.32,
    }

    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "563582419.111111").
        SetBody(reqBody).
        Post(url)

    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode(), "Expected status 200 OK")
    t.Logf("Response Status: %s", resp.Status())
    t.Logf("Response Body: %s", resp.String())
}

func TestNegativeCaseTarik(t *testing.T) {
    client := resty.New()

    url := "http://localhost:8100/api/v1/transaction/tarik"
    reqBody := map[string]float64{
        "nominal": 1000000.87,
    }

    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "563582419.111111").
        SetBody(reqBody).
        Post(url)

    assert.NoError(t, err)
    assert.Equal(t, 400, resp.StatusCode(), "Expected status 400 Bad Request")

    var response map[string]interface{}
    err = json.Unmarshal(resp.Body(), &response)
    assert.NoError(t, err)

    respMsg, ok := response["resp_msg"].(string)
    assert.True(t, ok, "Expected resp_msg to be a string")
    assert.Equal(t, "maaf, saldo tidak cukup", respMsg, "Expected response message 'maaf, saldo tidak cukup'")

    t.Logf("Response Status: %s", resp.Status())
    t.Logf("Response Body: %s", resp.String())
}