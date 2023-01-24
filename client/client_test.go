package client_test

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/openapi/api/excel"
	"github.com/weplanx/openapi/client"
	"os"
	"testing"
)

type TestEnv struct {
	Url        string `env:"URL"`
	Key        string `env:"KEY"`
	Secret     string `env:"SECRET"`
	client.Cos `envPrefix:"COS_"`
}

var x *client.Client

func TestMain(m *testing.M) {
	var err error
	var e TestEnv
	if err := env.Parse(&e); err != nil {
		panic(err)
	}
	if x, err = client.New(e.Url,
		client.SetApiGateway(e.Key, e.Secret),
		client.SetCos(e.Cos.Url, e.Cos.SecretID, e.Cos.SecretKey),
	); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestClient_Ping(t *testing.T) {
	data, err := x.Ping(context.TODO())
	assert.NoError(t, err)
	t.Log(data)
}

func TestClient_GetIp(t *testing.T) {
	data, err := x.GetIp(context.TODO(), "119.41.207.227")
	assert.NoError(t, err)
	t.Log(data)
}

func TestClient_GetCountries(t *testing.T) {
	data, err := x.GetCountries(context.TODO(), []string{"iso3"})
	assert.NoError(t, err)
	t.Log(data)
}

func TestClient_GetStates(t *testing.T) {
	data, err := x.GetStates(context.TODO(), "CN", []string{"type"})
	assert.NoError(t, err)
	t.Log(data)
}

func TestClient_GetCities(t *testing.T) {
	data, err := x.GetCities(context.TODO(), "CN", "AH", []string{"latitude"})
	assert.NoError(t, err)
	t.Log(data)
}

func TestClient_CreateExcel(t *testing.T) {
	data := [][]interface{}{
		{"Name", "CCType", "CCNumber", "Century", "Currency", "Date", "Email", "URL"},
	}
	for n := 0; n < 100000; n++ {
		data = append(data, []interface{}{
			faker.Name(), faker.CCType(), faker.CCNumber(), faker.Century(), faker.Currency(), faker.Date(), faker.Email(), faker.URL(),
		})
	}
	r, err := x.CreateExcel(context.TODO(), excel.CreateDto{
		Sheets: []excel.Sheet{
			{
				Name: "Sheet1",
				Data: data,
			},
		},
	})
	assert.NoError(t, err)
	t.Log(r)
}

func TestClient_Excel(t *testing.T) {
	sheets := []client.Sheet{
		{
			Name: "Sheet1",
			Data: []interface{}{"Name", "CCType", "CCNumber", "Century", "Currency", "Date", "Email", "URL"},
		},
	}
	for n := 0; n < 100000; n++ {
		sheets = append(sheets, client.Sheet{
			Name: "Sheet1",
			Data: []interface{}{
				faker.Name(), faker.CCType(), faker.CCNumber(), faker.Century(), faker.Currency(), faker.Date(), faker.Email(), faker.URL(),
			},
		})
	}
	err := x.Excel(context.TODO(), "test.excel", sheets)
	assert.NoError(t, err)
}
