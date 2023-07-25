package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"path/filepath"
)

//_______________________________________________//
//Entry Point//

// required flags: -id -type
func handleStatsDownload() {
	// Get ID flag
	id := getFlag("-id", "");
	if (id == "") {
		fmt.Println("You must provide a loadTestId with -id flag");
		return;
	}

	// Get Path flag
	path := getFlag("-path", "");
	if (path == "") {
		path = getDownloadDirPath();
	}
	response := httpRequestStatsDownload(id);
	parsedResponse := parseStatsDownloadJSON(response);

	// Get download types with -type flag
	downloadTypes := getMultiFlag("-type");
	if (len(downloadTypes) == 0) {
		fmt.Println("Must specify download type with -type flag, available types:");
		for _, dType := range parsedResponse {
			fmt.Println(dType[0]);
		}
		return;
	}

	for _, downloadType := range downloadTypes {
		foundType := false;
		for i := 0; i < len(parsedResponse); i++ {
			localDownloadType := parsedResponse[i][0];
			//fmt.Println(localDownloadType);
			if (localDownloadType == downloadType) {
				foundType = true;
				url := parsedResponse[i][1];

				fileType := ""
				if (localDownloadType == "merged" || localDownloadType == "report") {
					fileType = ".tar.gz";
				} else if (strings.Contains(localDownloadType, "outputFile")) {
					fileType = ".tgz";
				} else {
					fileType = ".csv";
				}
	
				downloadFile(url, path, fmt.Sprintf("%s_%s", id, localDownloadType), fileType);
			}	
		}
		if (!foundType) {
			fmt.Println("Unknown download type found");
		}
	}
}

//_______________________________________________//
//Core Functions//

func httpRequestStatsDownload(id string) []byte {
	client := http.Client{};

	url := fmt.Sprintf("%s/Api/StatsDownloadUrls?loadTestId=%s", build, id);

	req, err := http.NewRequest("GET", url, nil);
	if err != nil {
		fmt.Println("Error creating GET request: " + err.Error());
		return nil;
	}
	req.Header.Set("X-Redline-Auth", getAPIKEY());

	resp, err := client.Do(req);
	if err != nil {
		fmt.Println("Error sending GET request: " + err.Error());
		return nil;
	}
	defer resp.Body.Close();

	//fmt.Println("Response status:", resp.Status);

	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		fmt.Println("Error reading response body: " + err.Error());
		return nil;
	}

	if (resp.StatusCode != 200) {
		fmt.Println("Response status:", resp.Status);
		fmt.Println(string(body));
	}

	return body;
}

func downloadFile(url string, directory string, fileName string, fileType string) {
	fileName += fileType;

	err := os.MkdirAll(directory, os.ModePerm);
	if err != nil {
		fmt.Println(err);
		return; 
	}

	filePath := filepath.Join(directory, fileName);
	file, err := os.Create(filePath);
	if err != nil {
		fmt.Println(err);
		return; 
	}
	defer file.Close();

	response, err := http.Get(url);
	if err != nil {
		fmt.Println(err);
		return; 
	}
	defer response.Body.Close();

	_, err = io.Copy(file, response.Body);
	if err != nil {
		fmt.Println(err);
		return;
	}

	fmt.Println("File downloaded successfully:", filePath);
}

//_______________________________________________//
//Miscellaneous//

func printStatsDownloadInfo() {
	// fmt.Println("	statsDownload - Download loadTest stats as CSV");
	// fmt.Println("	    Flags:");	
	// fmt.Println("	        -id {load_test_id} : ID of loadTest to download data");
	// fmt.Println("	        -type {type1Value} {type2Value} : Specific download types to download");
	fmt.Println("Usage:")
	fmt.Println("    redline statsdownload [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("    -id - ID of loadTest to download data")
	fmt.Println("    -type {type1Value} {type2Value}... - Specific download types to download")
	fmt.Println("\nExample:")
	fmt.Println("    redline statsdownload -id 123321 -type cpuUsage netIn netOut")
}

func parseStatsDownloadJSON(jsonData []byte) [][]string {
	// Set up handling array in response
	var data map[string]json.RawMessage;
	err := json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err);
		return nil;
	}

	var arr [][]string;
	for key, value := range data {
		if key == "outputFiles" {
			var outputFiles []map[string]interface{};
			err := json.Unmarshal(value, &outputFiles)
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				for i, file := range outputFiles {
					sKey := fmt.Sprintf("outputFile%d", i);
					arr = append(arr, []string{sKey, file["url"].(string)});
				}
			}
		} else {
			var val string;
			err := json.Unmarshal(value, &val);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				arr = append(arr, []string{key, val});
			}
		}
	}

	return arr;
}

func getDownloadDirPath() string {
	// usr, err := user.Current();
	// if err != nil {
	// 	fmt.Println(err);
	// 	return "";
	// }

	// homeDir := usr.HomeDir;
	// downloadsDir := filepath.Join(homeDir, "Downloads");

    //return downloadsDir;

	// _____________________________________ //

	// Downlaod to current working directory 
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	return cwd;
}