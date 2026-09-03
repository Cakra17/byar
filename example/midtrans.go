package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cakra17/byar"
	"github.com/cakra17/byar/midtrans"
	"github.com/spf13/viper"
)

func loadConfig() (byar.Config, error) {
	var config byar.Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("Failed to load config: %s", err.Error())
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("Failed to unmarshal config: %s", err.Error())
	}

	return config, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

	service := midtrans.NewService(cfg)
	res, _ := service.CreateTransaction(context.Background(), &byar.Request{
		PaymentType: byar.BcaVa,
		TransactionDetails: byar.TransactionDetails{
			Orderid:     "0",
			GrossAmount: 200000,
		},
	})

	fmt.Println("Response :", res)

}
