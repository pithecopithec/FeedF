package main

import (
	"FeedFucker/structure"
	"encoding/xml"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var pf = fmt.Printf

func main() {

	//Для установки иконки
	//rsrc -ico="1.ico"

	dir, err := os.Getwd()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("cant find myself")
		fmt.Println("press ENTER")
		fmt.Scanln()
	}
	dir = dir + "\\"

	processDir(dir)

	//fmt.Println("press ENTER")
	fmt.Println("Нажми Enter")
	fmt.Scanln()
}

func processDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, f := range files {
		filename := f.Name()
		dim := filepath.Ext(filename)
		if dim == ".xml" {
			getdata(filename)
		}
	}
}

func getdata(filename string) {
	doc := structure.Yml_catalog{}
	xmlFile, err := os.Open(filename)
	if err != nil {
		pf("error Open xml %v", err)
	}

	defer xmlFile.Close()
	//pf("Start read %v \r\n", filename)
	pf("В работе - %v \r\n", filename)
	b := xml.NewDecoder(xmlFile)
	b.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	err = b.Decode(&doc)
	if err != nil {
		pf("error Decode %v", err)
	}
	//pf("read %v - OK \r\n",filename)

	esskeetit(filename, doc)

	pf("%s завершено!  \r\n", filename)
}

func esskeetit(filename string, doc structure.Yml_catalog) {
	xlsx := excelize.NewFile()
	catMap := make(map[string]string)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		file, err = os.Create(filename)
	}
	file.Close()

	catMap = getCatalog(xlsx, &doc, catMap)
	getOffer(xlsx, &doc, catMap)
	getParam(xlsx, &doc)
	getPhotoSheet(xlsx, &doc)

	xlsx.DeleteSheet("Sheet1")

	pf("Сохранение. Не закрывайте программу! \r\n")
	err = xlsx.SaveAs(string(filename[:len(filename)-4]) + ".xlsx")
	if err != nil {
		fmt.Printf("error save %v", err)
	}
}

// Создается лист с каталогом

func getCatalog(xlsx *excelize.File, doc *structure.Yml_catalog, catMap map[string]string) map[string]string {
	//pf("Start write Category \r\n")

	xlsx.NewSheet("Category")

	xlsx.SetCellValue("Category", "a1", "cat id")
	xlsx.SetCellValue("Category", "b1", "par id")
	xlsx.SetCellValue("Category", "c1", "cat_name")

	i := 1
	for _, cat := range doc.Shop.Categories.Categories {
		xlsx.SetCellValue("Category", "a"+row(i), cat.Id)
		xlsx.SetCellValue("Category", "b"+row(i), cat.ParentId)
		xlsx.SetCellValue("Category", "c"+row(i), cat.Name)
		catMap[cat.Id] = cat.Name
		i++
	}
	//pf("Category done \r\n")
	return catMap
}

func getOffer(xlsx *excelize.File, doc *structure.Yml_catalog, catMap map[string]string) {
	//pf("Start write Offer \r\n")
	xlsx.NewSheet("Offer")

	xlsx.SetCellValue("Offer", "a1", "Offer ID")
	xlsx.SetCellValue("Offer", "b1", "Name")
	xlsx.SetCellValue("Offer", "c1", "VendorCode")
	xlsx.SetCellValue("Offer", "d1", "Model")
	xlsx.SetCellValue("Offer", "e1", "Vendor")
	xlsx.SetCellValue("Offer", "f1", "Url")
	xlsx.SetCellValue("Offer", "g1", "CategoryId")
	xlsx.SetCellValue("Offer", "h1", "CategoryName")
	xlsx.SetCellValue("Offer", "i1", "Description")
	xlsx.SetCellValue("Offer", "j1", "Dimensions")
	xlsx.SetCellValue("Offer", "k1", "Weight")
	xlsx.SetCellValue("Offer", "l1", "Barcode")
	xlsx.SetCellValue("Offer", "m1", "Other_Barcode")
	xlsx.SetCellValue("Offer", "n1", "typePrefix")

	i := 1

	for _, offer := range doc.Shop.Offers.Offers {
		xlsx.SetCellValue("Offer", "a"+row(i), offer.Id)
		xlsx.SetCellValue("Offer", "b"+row(i), offer.Name)
		xlsx.SetCellValue("Offer", "c"+row(i), offer.VendorCode)
		xlsx.SetCellValue("Offer", "d"+row(i), offer.Model)
		xlsx.SetCellValue("Offer", "e"+row(i), offer.Vendor)
		xlsx.SetCellValue("Offer", "f"+row(i), offer.Url)
		xlsx.SetCellValue("Offer", "g"+row(i), offer.CategoryId)
		xlsx.SetCellValue("Offer", "h"+row(i), catMap[offer.CategoryId])
		xlsx.SetCellValue("Offer", "i"+row(i), offer.Description)
		xlsx.SetCellValue("Offer", "j"+row(i), offer.Dimensions)
		xlsx.SetCellValue("Offer", "k"+row(i), offer.Weight)
		if offer.Barcodes != nil {
			xlsx.SetCellValue("Offer", "l"+row(i), offer.Barcodes[0])
			xlsx.SetCellValue("Offer", "m"+row(i), zahyiar(offer.Barcodes))
		}
		xlsx.SetCellValue("Offer", "n"+row(i), offer.TypePrefix)
		i++
	}
	//pf("Offer done \r\n")
}

func getParam(xlsx *excelize.File, doc *structure.Yml_catalog) {
	//pf("Start write PARAM \r\n")
	xlsx.NewSheet("PARAM")
	sheetname := "PARAM"

	xlsx.SetCellValue("PARAM", "a1", "Offer ID")
	xlsx.SetCellValue("PARAM", "b1", "Attribute")
	xlsx.SetCellValue("PARAM", "c1", "Value")
	xlsx.SetCellValue("PARAM", "d1", "Unit")

	i := 1
	sheeti := 1
	for _, offer := range doc.Shop.Offers.Offers {
		for _, atr := range offer.Params {
			xlsx.SetCellValue(sheetname, "a"+row(i), offer.Id)
			xlsx.SetCellValue(sheetname, "b"+row(i), atr.Name)
			xlsx.SetCellValue(sheetname, "c"+row(i), atr.Value)
			xlsx.SetCellValue(sheetname, "d"+row(i), atr.Unit)
			i = i + 1
			if i == 1048575 {
				sheeti++
				xlsx.NewSheet("PARAM" + strconv.Itoa(sheeti))
				sheetname = "PARAM" + strconv.Itoa(sheeti)
				xlsx.SetCellValue("PARAM"+strconv.Itoa(sheeti), "a1", "Offer ID")
				xlsx.SetCellValue("PARAM"+strconv.Itoa(sheeti), "b1", "Attribute")
				xlsx.SetCellValue("PARAM"+strconv.Itoa(sheeti), "c1", "Value")
				xlsx.SetCellValue("PARAM"+strconv.Itoa(sheeti), "d1", "Unit")
				i = 1
			}
		}
	}

	//pf("PARAM done \r\n")
}

func getPhotoSheet(xlsx *excelize.File, doc *structure.Yml_catalog) {
	i := 1
	sheeti := 1

	xlsx.NewSheet("Photo")
	sheetname := "Photo"

	xlsx.SetCellValue("Photo", "a1", "offer.Id")
	xlsx.SetCellValue("Photo", "b1", "type")
	xlsx.SetCellValue("Photo", "c1", "url")

	for _, offer := range doc.Shop.Offers.Offers {
		ii := 0
		for _, photo := range offer.Picture {

			xlsx.SetCellValue(sheetname, "a"+row(i), offer.Id)
			xlsx.SetCellValue(sheetname, "b"+row(i), "b"+strconv.Itoa(ii))
			xlsx.SetCellValue(sheetname, "c"+row(i), photo)
			ii++
			i++
			if i == 1048575 {
				sheeti++
				xlsx.NewSheet("Photo" + strconv.Itoa(sheeti))
				sheetname = "Photo" + strconv.Itoa(sheeti)
				i = 1
			}

		}
	}
}

//ахахахахахха
// это ж JOIN

func zahyiar(arr []string) string {
	var str string
	var delim string
	if len(arr) > 1 {
		for _, i := range arr[1:] {
			str = str + delim + i
			delim = ";"
		}
	}
	return str
}

func row(i int) (str string) {
	return strconv.Itoa(i + 1)
}
