package main

import (
	"encoding/json";
	"fmt";
	"net/http";
	"io/ioutil";
	"io";
	"os";
	"os/user";
	"path/filepath";
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
		path = getDownloadsFolderPath();
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
				url := parsedResponse[0][1];
				downloadFile(url, path, fmt.Sprintf("%s_%s", id, localDownloadType));
			}	
		}
		if (!foundType) {
			fmt.Println("Unknown download type found");
		}
	}
}

//_______________________________________________//
//Core Functions//

func httpRequestStatsDownload(id string) []byte{
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

	return body;
}

func downloadFile(url string, directory string, fileName string) {
	fileName += ".csv";
	
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println(err);
		return; 
	}

	filePath := filepath.Join(directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err);
		return; 
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err);
		return; 
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err);
		return;
	}

	fmt.Println("File downloaded successfully:", filePath)
}

//_______________________________________________//
//Miscellaneous//

func printStatsDownloadInfo() {
	fmt.Println("	statsDownload - Download loadTest stats as CSV");
	fmt.Println("	    Flags:");	
	fmt.Println("	        -id {load_test_id} : ID of loadTest to download data");
}

func parseStatsDownloadJSON(jsonData []byte) [][]string{
	var data map[string]string
	err := json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil;
	}

	var arr [][]string
	for key, value := range data {
		arr = append(arr, []string{key, value})
	}
	return arr;
}

func getDownloadsFolderPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err);
		return "";
	}

	homeDir := usr.HomeDir

	downloadsDir := filepath.Join(homeDir, "Downloads")
	return downloadsDir
}

// func getDownloadTypes() []string {
// 	ret := []string{};

// 	active := false
// 	for i := 0; i < len(args); i++ {
// 		//fmt.Println(args[i]);
// 		typeFlag := args[i] == "-type"
// 		if (typeFlag) {
// 			active = true;
// 			//fmt.Println("actived");
// 		}
// 		isFlag := strings.Contains(args[i], "-")
// 		if (!isFlag && active) {
// 			ret = append(ret, args[i]);
// 		} else if (isFlag && !typeFlag) {
// 			active = false;
// 			//fmt.Println("unactived");
// 		}
// 	}
// 	return ret;
// }