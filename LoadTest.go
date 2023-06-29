package main

import (
	"fmt";
	"net/http";
	"io/ioutil";
	"net/url";
	"strings";
	"encoding/json";
	"bytes";
	"mime/multipart";
	"os";
	"path/filepath";
	"io";
)


//_______________________________________________//
//Entry Point//

func handleLoadTest() {
	argument := "";
	if (len(args) > 2) {
		argument = args[2];
	}

	switch argument {
	case "simple":
		simpleLoadTest();
	case "jmeter":
		jmeterLoadTest();
	case "logfile":
		logFileReplayTest();
	case "custom":
		customLoadTest();
	case "testJSON":
		fmt.Println(parseLoadTestJSON("serverData.json"));
	default:
		fmt.Println("Unknown test type provided");
	}
}


//_______________________________________________//
//Core Functions//

func simpleLoadTest() {
	form := url.Values{};
	//required flags : url, numUsers, numIterations, minDelayMs, maxDelayMs
	//additional flags : name, desc, storeOutput, rampUpSec, loadResources //// params
	
	form.Add("testType", "simple");

	url := getFlag("-url", "");
	form.Add("url", url);

	numUsers := getFlag("-numUsers", "");
	form.Add("numUsers", numUsers);

	numIterations := getFlag("-numIterations", "");
	form.Add("numIterations", numIterations);

	minDelayMs := getFlag("-minDelayMs", "0");
	form.Add("minDelayMs", minDelayMs);

	maxDelayMs := getFlag("-maxDelayMs", "10000");
	form.Add("maxDelayMs", maxDelayMs);

	if (url == "" || numUsers == "" || numIterations == "") {
		fmt.Println("Missing required flag(s)");
		return
	} 

	name := getFlag("-name", "");
	form.Add("name", name);

	desc := getFlag("-desc", "");
	form.Add("desc", desc);

	storeOutput := getFlag("-storeOutput", "");
	form.Add("storeOutput", storeOutput);

	rampUpSec := getFlag("-rampUpSec", "");
	form.Add("rampUpSec", rampUpSec);

	loadResources := getFlag("-loadResources", "");
	form.Add("loadResources", loadResources);

	serverData, err := parseLoadTestJSON("serverData.json");
	if err != nil {
		fmt.Println("Error parsing serverData JSON");
		return;
	}
	for i := 0; i < len(serverData); i++ {
		form.Add(serverData[i][0], serverData[i][1]);
	} 
	fmt.Println(form);


	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = addFieldsToBody(form, writer);
	if err != nil {
		fmt.Println("Error adding form fields:", err);
		return;
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}

	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
}

func jmeterLoadTest() {
	//required flags : numServers, version
	//additional flags : name, desc, storeOutput //// params
	form := url.Values{};

	form.Add("testType", "jmeter-test");

	filePath := getFileArg(".jmx");
	if (filePath == "") {
		fmt.Println("Please provide a Jmeter test file");
		return;
	}

	numServers := getFlag("-numServers", "1");
	form.Add("numServers", numServers);

	version := getFlag("-version", "5.5");
	form.Add("version", version);

	name := getFlag("-name", "");
	form.Add("name", name);

	desc := getFlag("-desc", "");
	form.Add("desc", desc);

	storeOutput := getFlag("-storeOutput", "");
	form.Add("storeOutput", storeOutput);

	serverData, err := parseLoadTestJSON("serverData.json");
	if err != nil {
		fmt.Println("Error parsing serverData JSON");
		return;
	}
	for i := 0; i < len(serverData); i++ {
		form.Add(serverData[i][0], serverData[i][1]);
	} 

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = addFieldsToBody(form, writer);
	if err != nil {
		fmt.Println("Error adding form fields:", err);
		return;
	}
	
	err = addFileToBody(filePath, writer);
	if err != nil {
		return;
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}
	//fmt.Println(writer.FormDataContentType());
	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
}

func logFileReplayTest() {
	form := url.Values{};

	form.Add("testType", "replay");

	filePath := getFileArg(".log");
	if (filePath == "") {
		fmt.Println("Please provide a Jmeter test file");
		return;
	}

	url := getFlag("-url", "");
	form.Add("url", url)

	numUsers := getFlag("-numUsers", "");
	form.Add("numUsers", numUsers);

	numIterations := getFlag("-numIterations", "");
	form.Add("numIterations", numIterations);

	minDelayMs := getFlag("-minDelayMs", "0");
	form.Add("minDelayMs", minDelayMs);

	maxDelayMs := getFlag("-maxDelayMs", "10000");
	form.Add("maxDelayMs", maxDelayMs);

	if (url == "" || numUsers == "" || numIterations == "") {
		fmt.Println("Missing required flag(s)");
		return;
	} 

	rampUpSec := getFlag("-rampUpSec", "");
	form.Add("rampUpSec", rampUpSec);

	loadResources := getFlag("-loadResources", "");
	form.Add("loadResources", loadResources);

	logFormat := getMultiFlag("-format");
	if !(len(logFormat) == 0) {
		str := strings.Join(logFormat, " ");
		form.Add("log_format", str);
	} 
	
	serverData, err := parseLoadTestJSON("serverData.json");
	if err != nil {
		fmt.Println("Error parsing serverData JSON");
		return;
	}
	for i := 0; i < len(serverData); i++ {
		form.Add(serverData[i][0], serverData[i][1]);
	}
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = addFieldsToBody(form, writer);
	if err != nil {
		fmt.Println("Error adding form fields:", err);
		return;
	}
	
	err = addFileToBody(filePath, writer);
	if err != nil {
		return;
	}

	// Close the multipart form
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart form:", err)
		return
	}
	//fmt.Println(writer.FormDataContentType());
	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
}

func customLoadTest() {
	form := url.Values{};

	form.Add("testType", "custom-test");

	lang := getFlag("-lang", "")
	if lang == "" {
		fmt.Println("You must provide a -lang (python, php, nodejs) flag");
		return;
	}
	fileType := "";
	switch lang {
	case "python":
		fileType = ".py"
	case "php":
		fileType = ".php"
	case "nodejs":
		// Don't know nodejs filetype yet
		return;
	}
		
	numUsers := getFlag("-numUsers", "1");
	form.Add("numUsers", numUsers);

	serverData, err := parseLoadTestJSON("serverData.json");
	if err != nil {
		fmt.Println("Error parsing serverData JSON");
		return;
	}
	for i := 0; i < len(serverData); i++ {
		form.Add(serverData[i][0], serverData[i][1]);
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = addFieldsToBody(form, writer);
	if err != nil {
		fmt.Println("Error adding form fields:", err);
		return;
	}

	filePath := getFileArg(fileType)
	if (filePath == "") {
		fmt.Println("Please provide a custom loadtest file");
		return;
	}

	err = addFileToBody(filePath, writer);
	if err != nil {
		return;
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart form:", err)
		return
	}
}


//http requests//
func httpPostRequest(body *bytes.Buffer, content string) string {
	client := http.Client{};
	
	//(form url.Values)
	//req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/Api/LoadTest", build), strings.NewReader(form.Encode()));
	
	//(body *bytes.Buffer)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/Api/LoadTest", build), body);

	fmt.Println(body);
	if err != nil {
		fmt.Println("Error creating request: " + err.Error());
		return "";
	}
	req.Header.Set("X-Redline-Auth", getAPIKEY());
	req.Header.Add("Content-Type", content);

	resp, err := client.Do(req);
	if err != nil {
		fmt.Println("Error sending request: " + err.Error());
		return "";
	}
	defer resp.Body.Close();

	fmt.Println("Response status:", resp.Status);

	responseBody, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		fmt.Println("Error reading response body: " + err.Error());
		return "";
	}

	return "Response body:" + string(responseBody);
}

//_______________________________________________//
//Miscellaneous//

func printLoadTestInfo() {
	fmt.Println("	loadTest - Starts a load test");
	fmt.Println("	    Simple Test");
	fmt.Println("	    JMeter Test");
	fmt.Println("	    LogFileReplay Test");
	fmt.Println("	    Gatling Test");
	fmt.Println("	    Custom Test");
}

func parseLoadTestJSON(path string) ([][]string, error) {
	jsonData, err := ioutil.ReadFile(path);
	if err != nil {
		fmt.Println("Error reading JSON file:", err);
		return nil, err;
	}

	var data map[string]string
	err = json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err);
		return nil, err;
	}

	var arr [][]string
	for key, value := range data {
		arr = append(arr, []string{key, value});
	}
	return arr, nil;
}

func addFieldsToBody(form url.Values, writer *multipart.Writer) error {
	for key, values := range form {
		for _, value := range values {
			err := writer.WriteField(key, value);
			if err != nil {
				return err;
			}
		}
	}
	return nil;
}

func addFileToBody(filePath string, writer *multipart.Writer) error {
	//filePath := "/Downloads/MyJmeterTest.log"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err;
	}
	defer file.Close()

	// Create the file part in the multipart form
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file part:", err)
		return err;
	}

	// Copy the file contents to the file part
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file contents to part:", err)
		return err;
	}
	return nil;
}

func getFileArg(fileType string) string {
	for i := 0; i < len(args); i++ {
		if (strings.Contains(args[i], fileType)) {
			return args[i];
		}
	}
	return "";
}