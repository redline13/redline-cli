package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//_______________________________________________//
//Entry Point//

func handleViewTest() {
	id := getFlag("-id", "")
	if (id != "") {
		showTest(id);
	} else {
		showTests();
	}
}


//_______________________________________________//
//Core Functions//

func showTest(id string) {
	request := httpRequestViewTest();
	parsedJSON := parseViewTestJSON(request);
	for _, test := range parsedJSON {
		isTest := false;
		for i := 0; i < len(test); i++ {
			if (test[i][0] == "load_test_id" && test[i][1] == id) {
				isTest = true;
			}
		}
		if (isTest) {
			for i := 0; i < len(test); i++ {
				fmt.Println(test[i][0] + ": " + test[i][1]);
			}
		}
	}
}

func showTests() {
	request := httpRequestViewTest();
	parsedJSON := parseViewTestJSON(request);
	for _, test := range parsedJSON {
		for i := 0; i < len(test); i++ {
			if (test[i][0] == "load_test_id" || test[i][0] == "load_test_name") {
				fmt.Println(test[i][0] + ": " + test[i][1]);
			}
		}
		fmt.Println();
	}
}

func httpRequestViewTest() []byte {
	client := http.Client{};

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Api/LoadTest", build), nil);
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


//_______________________________________________//
//Miscellaneous//

func printViewTestInfo() {
	fmt.Println("	viewTest - View existing load tests");
	fmt.Println("	    Flags:");	
	fmt.Println("	        -id {load_test_id} : Displays all info on specific loadtest");
}

func parseViewTestJSON(jsonData []byte) [][][]string {
	// Dump json into map
	var data []map[string]interface{};
	err := json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err);
		return nil;
	}

	// Convert map into 3d array [loadTest][dataEntry][(datapair: [key, value])]
	var Rarr [][][]string;
	for _, obj := range data {
		var arr [][]string;
		for key, value := range obj {
			strValue := "";
			switch v := value.(type) {
			case string:
				strValue = v;
			case int, int64, float64:
				strValue = fmt.Sprintf("%v", v);
			}
			arr = append(arr, []string{key, strValue});
		}
		Rarr = append(Rarr, arr);
	}
	return Rarr;
}
