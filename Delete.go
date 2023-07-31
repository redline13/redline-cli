package main

import (
	"fmt";
	"net/http";
	"io/ioutil";
	"encoding/json";
)

//_______________________________________________//
//Entry Point//


func handleDelete() {
	id := getFlag("-id", "")
	if (id != "") {
		deleteTest(id);
	} else {
		fmt.Println("Must provide loadTest id with -id flag");
	}
}

//_______________________________________________//
//Core Functions//


func deleteTest(id string) {
	response := httpRequestDelete(id);
	parsedResponse := parseDeleteJSON(response);

	if (parsedResponse["delete"] == "T") {
		fmt.Println("Deleted loadTest with id:", parsedResponse["loadTestId"]);
	} else {
		fmt.Println("LoadTest did not delete for an unknown reason");
	}
}

func httpRequestDelete(id string) []byte {
	client := http.Client{};

	url := fmt.Sprintf("%s/Api/LoadTest?loadTestId=%s", build, id);

	req, err := http.NewRequest(http.MethodDelete, url, nil);
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

//_______________________________________________//
//Miscellaneous//


func parseDeleteJSON(jsonData []byte) map[string]string {
	var data map[string]string;
	err := json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err);
		return nil;
	}
	return data;
}

func printDeleteInfo() {
	fmt.Println("Usage:");
	fmt.Println("    redline delete [flags]");
	fmt.Println("\nFlags:");
	fmt.Println("    -id - ID of loadTest to delete or cancel");
	fmt.Println("\nExample:");
	fmt.Println("    redline delete -id 123321");
}