package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var OpenWeatherApiKey = ""

func main() {
	//fmt.Println(getTitleAuthor("..\\public-domain-poetry\\poems\\W-M-MacKeracher-Milton.txt"))
	poemlist := populatePoemList()
	fmt.Println(len(poemlist))
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", weatherHandler)
	http.HandleFunc("/poems", listPoemsHandler)
	fmt.Println("Starting server...")
	fmt.Println("The Server lives! [CTRL+C to exit]")
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func weatherIn(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + OpenWeatherApiKey + "&q=" + city)
	if err != nil {
		fmt.Println(err)
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var responseWeather weatherData

	if err := json.NewDecoder(resp.Body).Decode(&responseWeather); err != nil {
		fmt.Println(err)
		return weatherData{}, err
	}

	return responseWeather, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := weatherIn(city)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

type poem struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func populatePoemList() []poem {
	files, err := ioutil.ReadDir("../public-domain-poetry/poems/")
	if err != nil {
		fmt.Println(err)
	}

	p := make([]poem, len(files))

	for i, ele := range files {
		p[i] = visit("../public-domain-poetry/poems/"+ele.Name(), ele)

		if i%10 == 0 || (i+1) == len(files) {
			fmt.Printf("Processed %d/%d", i+1, len(files))
			fmt.Println()
		}
	}
	return p
}

func listPoemsHandler(w http.ResponseWriter, r *http.Request) {
	data := "data"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func visit(path string, f os.FileInfo) poem {
	fmt.Printf("Processing %s \n", f.Name())
	t, a := getTitleAuthor(path)

	p := poem{Title: t, Author: a}
	return p
}

//Very specific function for the data we have.
func getTitleAuthor(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	//Skip the first --- of the file
	Readln(reader)
	//Extract the author
	line, err := Readln(reader)
	if err != nil {
		fmt.Println(err)
	}
	author := line[strings.Index(line, ":")+2:]
	//Extract the title
	line, err = Readln(reader)
	if err != nil {
		fmt.Println(err)
	}
	title := line[strings.Index(line, ":")+2:]

	return strings.TrimSpace(title), strings.TrimSpace(author)
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
