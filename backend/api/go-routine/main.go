package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var publicUrl = "frontend/public/"

var apiUrls = map[string]string{
	"base": "https://groupietrackers.herokuapp.com/api",
}

type ArtistsData struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type LocationsDataLevel2 struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationsDataLevel1 struct {
	Index []LocationsDataLevel2 `json:"index"`
}

type DatesDataLevel2 struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesDataLevel1 struct {
	Index []DatesDataLevel2 `json:"index"`
}

type RelationsDataLevel2 struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsDataLevel1 struct {
	Index []RelationsDataLevel2 `json:"index"`
}

type ErrorPageData struct {
	Name       string
	Code       string
	CodeNumber int
	Info       string
}

var PredefinedErrors = map[string]ErrorPageData{
	"BadRequestError": {
		Name:       "BadRequestError",
		Code:       strconv.Itoa(http.StatusBadRequest),
		CodeNumber: http.StatusBadRequest,
		Info:       "Bad request",
	},
	"NotFoundError": {
		Name:       "NotFoundError",
		Code:       strconv.Itoa(http.StatusNotFound),
		CodeNumber: http.StatusNotFound,
		Info:       "Page not found",
	},
	"MethodNotAllowedError": {
		Name:       "MethodNotAllowedError",
		Code:       strconv.Itoa(http.StatusMethodNotAllowed),
		CodeNumber: http.StatusMethodNotAllowed,
		Info:       "Method not allowed",
	},
	"InternalServerError": {
		Name:       "InternalServerError",
		Code:       strconv.Itoa(http.StatusInternalServerError),
		CodeNumber: http.StatusInternalServerError,
		Info:       "Internal server error",
	},
}

var (
	BadRequestError       = PredefinedErrors["BadRequestError"]
	NotFoundError         = PredefinedErrors["NotFoundError"]
	MethodNotAllowedError = PredefinedErrors["MethodNotAllowedError"]
	InternalServerError   = PredefinedErrors["InternalServerError"]
)

type RequestTask struct {
	url      string
	dataObj  interface{}
	response chan error
}

var workerPoolSize = 5

func worker(wg *sync.WaitGroup, tasks <-chan RequestTask) {
	defer wg.Done()
	for task := range tasks {
		err := sendGetRequest(task.url, task.dataObj, nil)
		task.response <- err
	}
}

func sendGetRequest(url string, dataObj interface{}, client *http.Client) error {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(body, &dataObj)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}

func getApiUrls() {
	type ApiInfo struct {
		Artists   string `json:"artists"`
		Locations string `json:"locations"`
		Dates     string `json:"dates"`
		Relations string `json:"relation"`
	}

	var dataObj ApiInfo
	sendGetRequest(apiUrls["base"], &dataObj, nil)

	apiUrls["artists"] = dataObj.Artists
	apiUrls["dates"] = dataObj.Dates
	apiUrls["locations"] = dataObj.Locations
	apiUrls["relations"] = dataObj.Relations
}

func generateUrl(path string, desiredUrl string) (string, string, string) {
	var url string
	if path == "/"+desiredUrl {
		url = apiUrls[desiredUrl]
		return url, "", ""
	} else if strings.HasPrefix(path, "/"+desiredUrl+"/") {
		id := strings.TrimPrefix(path, "/"+desiredUrl+"/")
		url = apiUrls[desiredUrl+"s"]
		return url, id, ""
	} else {
		return "", "", "not found"
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/" {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"index.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/artists_swiper.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	var dataObj []ArtistsData
	response := make(chan error)
	task := RequestTask{url: apiUrls["artists"], dataObj: &dataObj, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	tmpl.Execute(w, dataObj)
}

func handleArtists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	url, _, errUrl := generateUrl(r.URL.Path, "artists")
	if errUrl == "not found" {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"artists.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/artist_filter.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	var dataObjArray []ArtistsData
	response := make(chan error)
	task := RequestTask{url: url, dataObj: &dataObjArray, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Marshal JSON for the template
	jsonData, err := json.Marshal(dataObjArray)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Prepare structured data for the template
	type ArtistsDataForPass struct {
		Artists         []ArtistsData
		ArtistsJsonData string
	}

	dataObjSender := ArtistsDataForPass{
		Artists:         dataObjArray,
		ArtistsJsonData: string(jsonData),
	}

	// Execute template with structured data
	err = tmpl.Execute(w, dataObjSender)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}
}

func handleArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	url, id, errUrl := generateUrl(r.URL.Path, "artist")
	if errUrl == "not found" {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"artist.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/artist_info.html",
		publicUrl+"templates/artist_dates.html",
		publicUrl+"templates/artist_locations.html",
		publicUrl+"templates/artist_relation.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Initialize the response channel for asynchronous requests
	response := make(chan error)

	// Start the request for artist data
	var dataObj ArtistsData
	task := RequestTask{url: url + "/" + id, dataObj: &dataObj, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Fetch the related data asynchronously
	var dateDataObj DatesDataLevel2
	task = RequestTask{url: dataObj.ConcertDates, dataObj: &dateDataObj, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	var locationDataObj LocationsDataLevel2
	task = RequestTask{url: dataObj.Locations, dataObj: &locationDataObj, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	var relationDataObj RelationsDataLevel2
	task = RequestTask{url: dataObj.Relations, dataObj: &relationDataObj, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Prepare template data
	templateData := struct {
		ArtistInfo      ArtistsData
		ArtistDates     DatesDataLevel2
		ArtistLocations LocationsDataLevel2
		Relation        RelationsDataLevel2
	}{
		ArtistInfo:      dataObj,
		ArtistDates:     dateDataObj,
		ArtistLocations: locationDataObj,
		Relation:        relationDataObj,
	}

	// Render the template with the fetched data
	err = tmpl.Execute(w, templateData)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		handleErrorPage(w, r, InternalServerError)
		return
	}
}

func handleLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	url, _, errUrl := generateUrl(r.URL.Path, "locations")
	if errUrl == "not found" {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	// Parse the template files
	tmpl, err := template.ParseFiles(
		publicUrl+"locations.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/locations_list.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Create a response channel to synchronize goroutines
	artistsData := make([]ArtistsData, 0)
	locationsData := LocationsDataLevel1{}
	errCh := make(chan error, 2)

	// Fetch artists data asynchronously
	go func() {
		err := sendGetRequest(apiUrls["artists"], &artistsData, nil)
		errCh <- err
	}()

	// Fetch locations data asynchronously
	go func() {
		err := sendGetRequest(url, &locationsData, nil)
		errCh <- err
	}()

	// Wait for both requests to complete
	var err1, err2 error
	err1 = <-errCh
	err2 = <-errCh

	// Check if there were any errors
	if err1 != nil || err2 != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Prepare the data to be passed to the template
	templateData := struct {
		ArtistsData   []ArtistsData
		LocationsData LocationsDataLevel1
	}{
		ArtistsData:   artistsData,
		LocationsData: locationsData,
	}

	// Execute the template with the fetched data
	err = tmpl.Execute(w, templateData)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}
}

func handleDates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	url, _, errUrl := generateUrl(r.URL.Path, "dates")
	if errUrl == "not found" {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"dates.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/dates_list.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Create channels to capture data and errors
	artistsDataChannel := make(chan []ArtistsData)
	datesDataChannel := make(chan DatesDataLevel1)
	errorChannel := make(chan error, 2)

	// Goroutine for fetching artists data
	go func() {
		var dataObj []ArtistsData
		err := sendGetRequest(apiUrls["artists"], &dataObj, nil)
		if err != nil {
			errorChannel <- err // Use err here
			return
		}
		artistsDataChannel <- dataObj
	}()

	// Goroutine for fetching dates data
	go func() {
		var dataObj DatesDataLevel1
		err := sendGetRequest(url, &dataObj, nil)
		if err != nil {
			errorChannel <- err // Use err here
			return
		}
		datesDataChannel <- dataObj
	}()

	// Wait for both goroutines to finish and check for errors
	var artistsData []ArtistsData
	var datesData DatesDataLevel1

	select {
	case artistsData = <-artistsDataChannel:
	case err = <-errorChannel:
		handleErrorPage(w, r, InternalServerError)
		return
	}

	select {
	case datesData = <-datesDataChannel:
	case err = <-errorChannel:
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Render the template with the fetched data
	templateData := struct {
		ArtistsData []ArtistsData
		DatesData   DatesDataLevel1
	}{
		ArtistsData: artistsData,
		DatesData:   datesData,
	}

	tmpl.Execute(w, templateData)
}

func handleRelations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"relations.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/relations_list.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Create channels to capture data and errors
	artistsDataChannel := make(chan []ArtistsData)
	relationsDataChannel := make(chan RelationsDataLevel1)
	errorChannel := make(chan error, 2)

	// Goroutine for fetching artists data
	go func() {
		var dataObj []ArtistsData
		err := sendGetRequest(apiUrls["artists"], &dataObj, nil)
		if err != nil {
			errorChannel <- err // Use err here
			return
		}
		artistsDataChannel <- dataObj
	}()

	// Goroutine for fetching relations data
	go func() {
		var dataObj RelationsDataLevel1
		err := sendGetRequest(apiUrls["relations"], &dataObj, nil)
		if err != nil {
			errorChannel <- err // Use err here
			return
		}
		relationsDataChannel <- dataObj
	}()

	// Wait for both goroutines to finish and check for errors
	var artistsData []ArtistsData
	var relationsData RelationsDataLevel1

	select {
	case artistsData = <-artistsDataChannel:
	case err = <-errorChannel:
		handleErrorPage(w, r, InternalServerError)
		return
	}

	select {
	case relationsData = <-relationsDataChannel:
	case err = <-errorChannel:
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Render the template with the fetched data
	templateData := struct {
		ArtistsData   []ArtistsData
		RelationsData RelationsDataLevel1
	}{
		ArtistsData:   artistsData,
		RelationsData: relationsData,
	}

	tmpl.Execute(w, templateData)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleErrorPage(w, r, MethodNotAllowedError)
		return
	}

	searchText := r.URL.Query().Get("search_text")
	if searchText == "" {
		handleIndex(w, r) // Show the index page if no search term is provided.
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"artists.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/artist_filter.html", // Added missing template
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Use a channel to fetch data
	var artistsData []ArtistsData
	response := make(chan error)
	task := RequestTask{url: apiUrls["artists"], dataObj: &artistsData, response: response}
	tasks <- task
	err = <-response
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Filter the data
	var filteredData []ArtistsData
	for _, artist := range artistsData {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchText)) {
			filteredData = append(filteredData, artist)
		}
	}

	// Handle no matches
	if len(filteredData) == 0 {
		handleErrorPage(w, r, NotFoundError)
		return
	}

	// Marshal filtered data into JSON
	jsonData, err := json.Marshal(filteredData)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}

	// Prepare structured data for template
	type ArtistsDataForPass struct {
		Artists         []ArtistsData
		ArtistsJsonData string
	}

	dataObjSender := ArtistsDataForPass{
		Artists:         filteredData,
		ArtistsJsonData: string(jsonData),
	}

	// Execute the template
	err = tmpl.Execute(w, dataObjSender)
	if err != nil {
		handleErrorPage(w, r, InternalServerError)
		return
	}
}

func handleErrorPage(w http.ResponseWriter, r *http.Request, errorPageData ErrorPageData) {
	tmpl, err := template.ParseFiles(
		publicUrl+"errors.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/menu.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(errorPageData.CodeNumber)
	tmpl.Execute(w, errorPageData)
}

var tasks chan RequestTask

func main() {
	// fs := http.FileServer(http.Dir(publicUrl + "assets"))
	// http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/public/static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./frontend/public/img"))))

	tasks = make(chan RequestTask, workerPoolSize)

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go worker(&wg, tasks)
	}

	getApiUrls()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/artists", handleArtists)
	http.HandleFunc("/artist/", handleArtist)
	http.HandleFunc("/locations", handleLocations)
	http.HandleFunc("/dates", handleDates)
	http.HandleFunc("/tours", handleRelations)
	http.HandleFunc("/search", handleSearch)

	// Start the server on port 8082
	fmt.Println("Starting server on 0.0.0.0:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))

	close(tasks)
	wg.Wait()
}
