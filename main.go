package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/nelsonli33/go-parser/models"
	"github.com/nelsonli33/go-parser/repository"
	"github.com/tealeg/xlsx"
)

func uploadFile(r *http.Request) (*os.File, error) {
	fmt.Println("---- File Upload Start ----")
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		panic(err)
		return nil, err
	}

	file, header, err := r.FormFile("originalFile")
	if err != nil {
		log.Fatal("Error Retrieving the File")
		panic(err)
		return nil, err
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	// Create a temporary file within out tmp dir
	tempFile, err := ioutil.TempFile("test_data", strings.TrimSuffix(header.Filename, path.Ext(header.Filename))+"-*"+path.Ext(header.Filename))
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()

	io.Copy(tempFile, file)
	fmt.Println("---- File Upload End ----")
	return tempFile, nil
}

type TrafficAccRecord struct {
	Year        string `xlsx:"0"`
	Month       string `xlsx:"1"`
	Day         string `xlsx:"2"`
	Hour        string `xlsx:"3"`
	Minute      string `xlsx:"4"`
	DeathCount  string `xlsx:"7"`
	InjuryCount string `xlsx:"8"`
	Latitude    string `xlsx:"47"`
	Longitude   string `xlsx:"48"`
}

func readExcelData(tempFilePath string) ([]models.TrafficAccident, error) {
	xlFile, err := xlsx.OpenFile(tempFilePath)
	if err != nil {
		return nil, err
	}

	trafficAccidents := make([]models.TrafficAccident, 0)
	trafficAccRecord := new(TrafficAccRecord)

	for _, sheet := range xlFile.Sheets {
		for i := 1; i < sheet.MaxRow; i++ {
			err := sheet.Rows[i].ReadStruct(trafficAccRecord)
			if err != nil {
				return nil, err
			}

			if trafficAccRecord.DeathCount == "" || trafficAccRecord.InjuryCount == "" {
				continue
			}

			deathCount, _ := strconv.Atoi(trafficAccRecord.DeathCount)
			injuryCount, _ := strconv.Atoi(trafficAccRecord.InjuryCount)
			lat, _ := strconv.ParseFloat(trafficAccRecord.Latitude, 32)
			lon, _ := strconv.ParseFloat(trafficAccRecord.Longitude, 32)

			accident := models.TrafficAccident{
				Date:        parseTime(trafficAccRecord),
				DeathCount:  deathCount,
				InjuryCount: injuryCount,
				Latitude:    float32(lat),
				Longitude:   float32(lon),
			}

			trafficAccidents = append(trafficAccidents, accident)
		}
	}

	return trafficAccidents, nil
}

func parseTime(record *TrafficAccRecord) string {
	year, _ := strconv.Atoi(record.Year)
	monthInt, _ := strconv.Atoi(record.Month)
	month := time.Month(monthInt)
	day, _ := strconv.Atoi(record.Day)
	hour, _ := strconv.Atoi(record.Hour)
	minute, _ := strconv.Atoi(record.Minute)

	return time.Date(year, month, day, hour, minute, 00, 0, time.UTC).Format("2006-01-02T15:04:05")
}

type ResponseMsg struct {
	Message string
}

func importFileToDBHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("index.html"))

	tempFile, err := uploadFile(r)
	if err != nil {
		panic(err)
		tmpl.Execute(w, ResponseMsg{Message: "Your file import failed."})
	}

	trafficAccidents, err := readExcelData(tempFile.Name())
	if err != nil {
		panic(err)
		tmpl.Execute(w, ResponseMsg{Message: "Your file import failed."})
	}

	defer os.Remove(tempFile.Name())

	db, err := models.InitDB()
	if err != nil {
		panic(err)
		tmpl.Execute(w, err)
	}

	accRepository := repository.NewAccidentRepository(db)

	accRepository.ClearAccidentTable()
	accRepository.InsertAccidents(trafficAccidents)

	tmpl.Execute(w, ResponseMsg{Message: "Your file import already succesful."})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, "Your file import already succesful.")
}

func runSever() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/importFile", importFileToDBHandler)

	fmt.Printf("Start Server at %s\n", "localhost:8000")
	// server is listening on port
	http.ListenAndServe(":8000", nil)

}

func main() {
	runSever()
}
