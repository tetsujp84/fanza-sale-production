package module

import (
	"google.golang.org/api/sheets/v4"
)

func ConvertFromProductionToSheet(productionList []*Production) sheets.ValueRange {
	var sheetValue sheets.ValueRange

	len := len(productionList)
	sheetValue.Values = make([][]interface{}, len)
	for i := 0; i < len; i++ {
		production := productionList[i]
		sheetValue.Values[i] = make([]interface{},5)
		sheetValue.Values[i][0] = production.title
		sheetValue.Values[i][1] = production.publisher
		sheetValue.Values[i][2] = production.price
		sheetValue.Values[i][3] = production.discountRate
		sheetValue.Values[i][4] = production.url
	}
	return sheetValue
}
