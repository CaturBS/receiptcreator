package main

import (
	// "log"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	// "io/ioutil"

	"subaga.com/receiptcreator/custwidget"

	// "math"
	// "bufio"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/data/validation"
	// "fyne.io/fyne/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type JSONTime struct {
	time.Time
}

type myTheme struct{}

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
	Filename      string   `json:"filename"`
	Stamptime     string   `json:"stamptime"`
	Valid         int      `json:"valid"`
	Edit          int      `json:"edit"`
	ReceiptHead   string   `json:"receipt_head"`
	Forward       int      `json:"forward"`
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
		GrandTotal:    grandTotal,
		Valid:         1,
		Filename:      "BBS087741671156090",
		ReceiptHead:   "NULL",
		Stamptime:     "2023-03-10 16:38:13.0"}
	data, err := json.Marshal(randomReceipt)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func main() {
	myApp := app.New()
	// myApp.Settings().SetTheme(theme.Vari)
	myWindow := myApp.NewWindow("Receipt Creator")
	myWindow.Resize(fyne.NewSize(1200, 300))
	rand.Seed(time.Now().UnixNano())
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter text...")
	button := widget.NewButton("select folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			entry.SetText(uri.Path())
		}, myWindow)
	})
	centered := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), button)
	// entry.Disable()
	entry1 := custwidget.NewIntegerEntry()
	entry1.SetText("1000")
	entry2 := custwidget.NewIntegerEntry()
	entry2.SetText("5")

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() { // optional, handle form submission
			loop, _ := strconv.Atoi(entry2.Text)
			for i := 0; i < loop; i++ {
				f, _ := os.Create(entry.Text + "/receipt")
				defer f.Close()
				count, _ := strconv.Atoi(entry1.Text)
				f.WriteString("[")
				for j := 0; j < count; j++ {
					if j > 0 {
						f.WriteString(",\n")
					}
					y, _ := createRandomReceipt()
					// log.Println(strconv.Itoa(i) + ": " + y)
					f.WriteString("\t" + y)
				}
				f.WriteString("\n]")
				url := "http://localhost:8000/transaction?debug=true"
				client := &http.Client{
					Timeout: time.Second * 30,
				}
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				fw, _ := writer.CreateFormField("device_id")
				io.Copy(fw, strings.NewReader("bbs08774"))

				fw, _ = writer.CreateFormFile("trans_file", "receipt")
				file, _ := os.Open(entry.Text + "/receipt")
				io.Copy(fw, file)
				writer.Close()
				// log.Println(body)
				// log.Println(url)
				req, _ := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
				req.Header.Set("Content-Type", writer.FormDataContentType())
				req.Header.Set("Authorization", "Bearer YbrB5HQX2zicn8A3Ffif5ycoEMmd4JMhTjpFUM2V")
				client.Do(req)

				// b, _ := ioutil.ReadAll(resp.Body)
				// log.Println(string(b))
				// log.Println(resp.StatusCode)
			}

		},
	}

	// we can also append items
	form.Append("Folder", entry)
	form.Append("", centered)
	form.Append("Limit", entry1)
	form.Append("Loop", entry2)

	myWindow.SetContent(form)

	myWindow.ShowAndRun()
}
