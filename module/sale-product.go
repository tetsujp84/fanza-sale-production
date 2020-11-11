package module

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/dmmlabo/dmm-go-sdk"
	"github.com/dmmlabo/dmm-go-sdk/api"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxLoopCount = 10
	fetchLength = 100
)


func GetProductions(isWriteToFile bool) []*Production {
	var productionList []*Production
	var wg sync.WaitGroup
	wg.Add(maxLoopCount)
	for count := 0; count < maxLoopCount; count++ {
		fmt.Printf("%d回目の取得\n", count)
		c := count
		dmmapi := initializeAPI(c)
		result, err := dmmapi.Execute()
		if err != nil {
			fmt.Println(err)
			wg.Done()
			return nil
		}
		productionList = append(productionList, getProductionList(result)...)
		wg.Done()
	}
	wg.Wait()
	sort.Sort(sort.Reverse(ProductionList(productionList)))

	if isWriteToFile {
		writeToFile(productionList)	
	}
	return productionList
}

func writeToFile(productionList []*Production)  {
	file, err := os.Create("result.txt")
	if err != nil {
		fmt.Println("生成失敗")
	}
	defer file.Close()
	for _, p := range productionList {
		file.Write(([]byte)(p.getPrintStr()))
	}
}

// 検索結果からProduction化
func getProductionList(result *api.ProductResponse) []*Production {
	itemCount := int(result.ResultCount)
	fmt.Printf("itemCount %d\n", itemCount)
	var items = make([]*Production, 0)
	var wg sync.WaitGroup
	wg.Add(itemCount)
	for i := 0; i < itemCount; i++ {
		var k = i
		go func() {
			p, err := getProductin(&result.Items[k])
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}
			items = append(items, p)
			wg.Done()
		}()
	}
	wg.Wait()
	return items
}

func initializeAPI(count int) *api.ProductService {
	client := dmm.New("fordmmte-999", "※ApiIDを設定")
	dmmapi := client.Product
	dmmapi.SetSite(api.SiteAdult)
	dmmapi.SetService("doujin")
	dmmapi.SetFloor("digital_doujin")
	dmmapi.SetSort("rank")
	dmmapi.SetOffset(int64(count * fetchLength))
	dmmapi.SetLength(int64(fetchLength))
	return dmmapi
}

type Production struct {
	title        string
	url          string
	publisher    string
	price        int
	subPrice     int
	discountRate int
}

func (p Production) print() {
	fmt.Println(p.title + "," + strconv.Itoa(p.price) + "," + strconv.Itoa(p.discountRate))
	return
}
func (p Production) getPrintStr() string {
	return p.title + "," + p.publisher + "," + strconv.Itoa(p.price) + "," + strconv.Itoa(p.discountRate) + "," + p.url + "\n"
}

type ProductionList []*Production

func (p ProductionList) Len() int {
	return len(p)
}

func (p ProductionList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ProductionList) Less(i, j int) bool {
	return p[i].discountRate < p[j].discountRate
}

// ストアのプロダクト取得
func getProductin(item *api.Item) (*Production, error) {
	// ストア値取得
	price, ep := strconv.Atoi(item.PriceInformation.Price)
	if ep != nil {
		//fmt.Println(item.Title + "のストア価格取得失敗")
		return nil, errors.New(item.Title + "のストア価格取得失敗")
	}
	subPriceStr := item.PriceInformation.RetailPrice //getSubPrice(item.URL)
	// urlで元値取得
	subPrice, es := strconv.Atoi(subPriceStr)
	if es != nil {
		//fmt.Println(item.Title + "の元価格取得失敗")
		subPrice = 0
	}

	rate := 0
	if subPrice != 0 {
		rate = 100 - int(float64(price)/float64(subPrice)*100.0)
	}

	return &Production{item.Title, item.URL, item.ItemInformation.Maker[0].Name, price, subPrice, rate}, nil
}

// 割引前価格の取得（旧取得方法、スクレイピングの利用）
func getSubPrice(url string) string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	subPrice := doc.Find(".priceList__sub.priceList__sub--big").Text()
	subPrice = strings.Replace(subPrice, "\n", "", -1)
	subPrice = strings.Replace(subPrice, " ", "", -1)
	subPrice = strings.Replace(subPrice, "円", "", -1)
	subPrice = strings.Replace(subPrice, ",", "", -1)
	return subPrice
}
