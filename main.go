package main

import (
	"./module"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main()  {
	var region = os.Getenv("IsAWS")
	if  region == "true" {
		lambda.Start(entry)
		return
	}
	entry()
}

func entry() {
	fmt.Println("実行開始")
	productionList := module.GetProductions(false)
	fmt.Println("アイテム取得　完了")
	valueRange := module.ConvertFromProductionToSheet(productionList)
	fmt.Println("変換　完了")
	module.WriteToSpreadSheet(valueRange)
	fmt.Println("シート更新　完了")
}