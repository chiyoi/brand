package data

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
)

var (
	EndpointCosmos = os.Getenv("ENDPOINT_COSMOS")
	NameDatabase   = "atri"
	NameContainer  = "brand_data"

	KeyData = "data"
)

var Data DataModel

func Container(databaseID, containerID string) (c *azcosmos.ContainerClient, err error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return
	}
	client, err := azcosmos.NewClient(EndpointCosmos, cred, nil)
	if err != nil {
		return
	}
	return client.NewContainer(databaseID, containerID)
}

func Load() {
	logs.Info("Load data.")
	c, err := Container(NameDatabase, NameContainer)
	if err != nil {
		logs.Panic(err)
	}

	var v struct {
		Data DataModel
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := cosmos.PointRead(ctx, c, KeyData, &v); err != nil {
		if az.IsNotFound(err) {
			logs.Warning("No data saved.")
		} else {
			logs.Panic(err)
		}
	}
	Data = v.Data
	logs.Info("Data loaded.")
}

func Save() {
	logs.Info("Save data.")
	c, err := Container(NameDatabase, NameContainer)
	if err != nil {
		logs.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := cosmos.PointDelete(ctx, c, KeyData); err != nil {
		if az.IsNotFound(err) {
			logs.Warning("New data.")
		} else {
			logs.Panic(err)
		}
	}
	if err := cosmos.PointCreate(ctx, c, KeyData, struct {
		ID   string `json:"id"`
		Data DataModel
	}{KeyData, Data}); err != nil {
		logs.Panic(err)
	}
	logs.Info("Data saved.")
}

type DataModel struct {
	LatestID string `json:"latest_id"`
}
