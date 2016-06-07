package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
)

func getIcalFromUrl(url string, user string, password string) string{
  fmt.Printf("Start getting ical from [%s].\n", url)
  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  req.SetBasicAuth(user, password)

  resp, err := client.Do(req)

  if err != nil {
      panic(err)
  }

  defer resp.Body.Close()
  htmlData, err := ioutil.ReadAll(resp.Body)
 	if err != nil {
 		panic(err)
 	}

  fmt.Printf("[%d] bytes retrieved.\n", len(htmlData))

  return string(htmlData)
}

func main() {
  if len(os.Args) != 4 {
    fmt.Printf("Usage: %s URL user password\n", os.Args[0])
    os.Exit(1)
  }

  icalData := getIcalFromUrl(os.Args[1], os.Args[2], os.Args[3])
  _ = icalData

}
