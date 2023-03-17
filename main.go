package main

import (
	"log"

	"encoding/json"
	"fmt"

	// "math"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type Receipt struct {
	Devid         string   `json:"devid"`
	Datetime      JSONTime `json:"datetime"`
	ReceiptNumber string   `json:"receipt_number"`
	NonTax        int      `json:"non_tax"`
	Subtotal      int      `json:"subtotal"`
	Discount      int      `json:"discount"`
	Service       int      `json:"service"`
	Tax           int      `json:"tax"`
	GrandTotal    int      `json:"grand_total"`
}

func getReceiptNumber() string {
	val := rand.Intn(1000000)
	valText := "000000" + strconv.Itoa(val)
	return valText[len(valText)-6:]
}
func randate() time.Time {
	min := time.Date(2019, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2023, 3, 18, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
func getSubTotal() int {
	return 10000 + rand.Intn(23)*10000
}
func createRandomReceipt() (string, error) {
	randIn := rand.Intn(150)
	randSuffix := ""
	if randIn < 100 {
		randSuffix = "0" + strconv.Itoa(randIn)
	} else {
		randSuffix = strconv.Itoa(randIn)
	}
	var devId string = "ND4X504E310200" + randSuffix
	var subtotal = getSubTotal()
	var tax = subtotal / 10
	var grandTotal = subtotal + tax
	// var datetimex =
	randomReceipt := &Receipt{Devid: devId, Datetime: JSONTime{randate()},
		ReceiptNumber: getReceiptNumber(),
		Subtotal:      subtotal,
		Tax:           tax,
		GrandTotal:    grandTotal}
	data, err := json.Marshal(randomReceipt)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Choice Widgets")
	x := time.Now()
	log.Println(x)
	rand.Seed(time.Now().UnixNano())
	y, _ := createRandomReceipt()
	log.Println(y)
	y, _ = createRandomReceipt()
	log.Println(y)
	entry := widget.NewEntry()
	// entry.Disable()

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				entry.SetText(uri.Path())
			}, myWindow)
			// myWindow.Close()
		},
	}

	// we can also append items
	form.Append("Entry", entry)

	myWindow.SetContent(form)

	myWindow.ShowAndRun()
}
