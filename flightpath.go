package main
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"errors"
)

/*
 * Config to bootstrap the url monitor
 */
 type Config struct {
	// Address to listen on
	Address      string        `json:"address"`
	// Logfile for stdout/debug
	LogFile      string        `json:"logfile"`
}

/*
 * FlightPath struct implements TODO 
 */
 type FlightPath struct {
	Cfg    *Config
	Client *http.Client
}

/*
 * ParseConfig parses the json config file and picks the URLs to be monitored
 */
 func ParseConfig(fileName string) (*Config, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	raw, err := ioutil.ReadAll(jsonFile)
	cfg := &Config{}
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		fmt.Println("Unmarshal failed", err)
		return nil, err
	}
	return cfg, nil
}

/*
 * GetFlighgPath returns a flightPath based on given inputs
 * Topological sort is used to make sure the children destinations are NOT ahead of parent destinations
 */
 func (f *FlightPath) GetFlightPath(inputMap map[string]string) (error, string) {

	// Detect cycles or for anamolies
	// Note : Only one way is taken care, no roundtrips (cycles) or disjoint paths
	reverseMap := map[string]string{}
	for k, v := range inputMap {
		reverseMap[v] = k
	}

	// find the starting location
	startingLocation := ""
	for k, _ := range inputMap {
		if _, exists := reverseMap[k]; !exists {
			startingLocation = k
			break
		}
	}
	if startingLocation == "" {
		fmt.Println("Invalid input")
		return errors.New("Invalid input"), startingLocation
	}
	flightPath:= flights(inputMap)

	return nil, flightPath
}

/*
 * FlightPathHandler that handles the incoming requests
 */
 func FlightPathHandler(f *FlightPath) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside the Handler")
		if r.URL.Path != "/calculate" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "POST":
			reqBody, err := ioutil.ReadAll(r.Body)
              		if err != nil {
                       		log.Fatal(err)
              		}
	                fmt.Printf("%s\n", reqBody)
			inputMap := map[string]string{}
			err = json.Unmarshal(reqBody, &inputMap)
			if err != nil {
				fmt.Println("Unmarshal failed", err)
				//TODO check this
				w.Write([]byte("Unmarshal failed"))
			}
			err, mResp := f.GetFlightPath(inputMap)
			w.Write([]byte(mResp))
		default:
			w.WriteHeader(http.StatusNotImplemented)
			v := http.StatusText(http.StatusNotImplemented) + "\n"
			w.Write([]byte(v))

		}
	}

}


func main() {

	cfg, err := ParseConfig("/tmp/config.json")
	if err != nil {
		fmt.Println("ParseConfig returned error", err)
		return
	}
	log.Println("++++++++++ Welcome to the flight path detection ++++++")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: 5 * time.Second,
	}

	f := &FlightPath{cfg, client}

	http.HandleFunc("/calculate", FlightPathHandler(f))
	http.ListenAndServe(cfg.Address, nil)

}
