package main

import (
	"fmt"
	ics "github.com/PuloV/ics-golang"
	"io/ioutil"
	"net/http"
	"os"
)

func getIcalFromUrl(url string, user string, password string) (string, error) {
	fmt.Printf("Start getting ical from [%s].\n", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, password)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error code received [%d/%s].", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("[%d] bytes retrieved.\n", len(htmlData))

	return string(htmlData), nil
}

func writeStringToFile(str string, fileName string) error {
	return ioutil.WriteFile(fileName, []byte(str), 0644)
}

func parseCalendar(calendarFile string) (*ics.Calendar, error) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- calendarFile
	parser.Wait()
	parseErrors, err := parser.GetErrors()
	if err != nil {
		return nil, err
	}
	if len(parseErrors) != 0 {
		return nil, fmt.Errorf("Expected 0 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}
	calendars, errCal := parser.GetCalendars()
	if errCal != nil {
		return nil, err
	}

	fmt.Printf("[%d] calendar(s) found. Take first one!\n", len(calendars))
	if len(calendars) == 0 {
		return nil, fmt.Errorf("No calendars found!")
	}

	return calendars[0], nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s URL user password\n", os.Args[0])
		os.Exit(1)
	}

	icalData, err := getIcalFromUrl(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		fmt.Printf("Failed to retrieve calendar [%s].\n", err)
		os.Exit(1)
	}

	calendarFile := "calendar.ics"
	err = writeStringToFile(icalData, calendarFile)
	defer os.Remove(calendarFile)
	if err != nil {
		fmt.Printf("Failed to write calendar into file [%s] ( %s ) \n", calendarFile, err)
		os.Exit(1)
	}

	calendar, err := parseCalendar(calendarFile)
	if err != nil {
		fmt.Printf("Failed to parse ical file. (%s)\n", err)
	}

	fmt.Printf("The calendar contains [%d] events.\n", len(calendar.GetEvents()))
}
