package main

import (
	"bytes";
	"encoding/json";
	"errors";
	"fmt";
	"io";
	"io/ioutil";
	"mime/multipart";
	"net/http";
	"net/url";
	"os";
	"os/exec";
	"path/filepath";
	"runtime";
	"strings";
	"time";
)

var bodyData map[string]interface{};

//_______________________________________________//
//Entry Point//

func handleLoadTest(shortCall bool) {
	argument := "";
	if shortCall {
		argument = args[1];
	} else if len(args) > 2 {
		argument = args[2];
	}

	if (getFileArg(".jmx") != "") {
		argument = "jmeter";
	}


	switch argument {
	case "simple":
		//simpleLoadTest();
	case "jmeter":
		jmeterLoadTest();
	case "logfile":
		//logFileReplayTest();
	case "custom":
		//customLoadTest();
	case "test":
		testLoadTest();
	default:
		fmt.Println("Unknown test type provided");
	}
}

//_______________________________________________//
//Core Functions//

// func simpleLoadTest() {
// 	form := url.Values{};
// 	//required flags : url, numUsers, numIterations, minDelayMs, maxDelayMs
// 	//additional flags : name, desc, storeOutput, rampUpSec, loadResources //// params

// 	form.Add("testType", "simple");

// 	url := getFlag("-url", "");
// 	form.Add("url", url);

// 	numUsers := getFlag("-numUsers", "");
// 	form.Add("numUsers", numUsers);

// 	numIterations := getFlag("-numIterations", "");
// 	form.Add("numIterations", numIterations);

// 	minDelayMs := getFlag("-minDelayMs", "0");
// 	form.Add("minDelayMs", minDelayMs);

// 	maxDelayMs := getFlag("-maxDelayMs", "10000");
// 	form.Add("maxDelayMs", maxDelayMs);

// 	if (url == "" || numUsers == "" || numIterations == "") {
// 		fmt.Println("Missing required flag(s)");
// 		return
// 	}

// 	name := getFlag("-name", "");
// 	form.Add("name", name);

// 	desc := getFlag("-desc", "");
// 	form.Add("desc", desc);

// 	storeOutput := getFlag("-storeOutput", "");
// 	form.Add("storeOutput", storeOutput);

// 	rampUpSec := getFlag("-rampUpSec", "");
// 	form.Add("rampUpSec", rampUpSec);

// 	loadResources := getFlag("-loadResources", "");
// 	form.Add("loadResources", loadResources);

// 	serverData, err := parseLoadTestJSON("serverData.json");
// 	if err != nil {
// 		fmt.Println("Error parsing serverData JSON");
// 		return;
// 	}
// 	for i := 0; i < len(serverData); i++ {
// 		form.Add(serverData[i][0], serverData[i][1]);
// 	}
// 	fmt.Println(form);

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	err = addFieldsToBody(form, writer);
// 	if err != nil {
// 		fmt.Println("Error adding form fields:", err);
// 		return;
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		fmt.Println("Error closing multipart form:", err);
// 		return;
// 	}

// 	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
// }

// func jmeterLoadTest() {
// 	//required flags : numServers, version
// 	//additional flags : name, desc, storeOutput //// params
// 	form := url.Values{};

// 	form.Add("testType", "jmeter-test");

// 	filePath := getFileArg(".jmx");
// 	if (filePath == "") {
// 		fmt.Println("Please provide a Jmeter test file");
// 		return;
// 	}

// 	numServers := getFlag("-numServers", "1");
// 	form.Add("numServers", numServers);

// 	version := getFlag("-version", "5.5");
// 	form.Add("version", version);

// 	name := getFlag("-name", "");
// 	form.Add("name", name);

// 	desc := getFlag("-desc", "");
// 	form.Add("desc", desc);

// 	storeOutput := getFlag("-storeOutput", "");
// 	form.Add("storeOutput", storeOutput);

// 	serverData, err := parseLoadTestJSON("serverData.json");
// 	if err != nil {
// 		fmt.Println("Error parsing serverData JSON");
// 		return;
// 	}
// 	for i := 0; i < len(serverData); i++ {
// 		form.Add(serverData[i][0], serverData[i][1]);
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	err = addFieldsToBody(form, writer);
// 	if err != nil {
// 		fmt.Println("Error adding form fields:", err);
// 		return;
// 	}

// 	err = addFileToBody(filePath, writer, "file");
// 	if err != nil {
// 		return;
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		fmt.Println("Error closing multipart form:", err);
// 		return;
// 	}
// 	//fmt.Println(writer.FormDataContentType());
// 	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
// }

// func logFileReplayTest() {
// 	form := url.Values{};

// 	form.Add("testType", "replay");

// 	filePath := getFileArg(".log");
// 	if (filePath == "") {
// 		fmt.Println("Please provide a Jmeter test file");
// 		return;
// 	}

// 	url := getFlag("-url", "");
// 	form.Add("url", url)

// 	numUsers := getFlag("-numUsers", "");
// 	form.Add("numUsers", numUsers);

// 	numIterations := getFlag("-numIterations", "");
// 	form.Add("numIterations", numIterations);

// 	minDelayMs := getFlag("-minDelayMs", "0");
// 	form.Add("minDelayMs", minDelayMs);

// 	maxDelayMs := getFlag("-maxDelayMs", "10000");
// 	form.Add("maxDelayMs", maxDelayMs);

// 	if (url == "" || numUsers == "" || numIterations == "") {
// 		fmt.Println("Missing required flag(s)");
// 		return;
// 	}

// 	rampUpSec := getFlag("-rampUpSec", "");
// 	form.Add("rampUpSec", rampUpSec);

// 	loadResources := getFlag("-loadResources", "");
// 	form.Add("loadResources", loadResources);

// 	logFormat := getMultiFlag("-format");
// 	if !(len(logFormat) == 0) {
// 		str := strings.Join(logFormat, " ");
// 		form.Add("log_format", str);
// 	}

// 	serverData, err := parseLoadTestJSON("serverData.json");
// 	if err != nil {
// 		fmt.Println("Error parsing serverData JSON");
// 		return;
// 	}
// 	for i := 0; i < len(serverData); i++ {
// 		form.Add(serverData[i][0], serverData[i][1]);
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	err = addFieldsToBody(form, writer);
// 	if err != nil {
// 		fmt.Println("Error adding form fields:", err);
// 		return;
// 	}

// 	err = addFileToBody(filePath, writer, "file");
// 	if err != nil {
// 		return;
// 	}

// 	// Close the multipart form
// 	err = writer.Close();
// 	if err != nil {
// 		fmt.Println("Error closing multipart form:", err);
// 		return;
// 	}
// 	//fmt.Println(writer.FormDataContentType());
// 	fmt.Println(httpPostRequest(body, writer.FormDataContentType()));
// }

// func customLoadTest() {
// 	form := url.Values{};

// 	form.Add("testType", "custom-test");

// 	lang := getFlag("-lang", "")
// 	if lang == "" {
// 		fmt.Println("You must provide a -lang (python, php, nodejs) flag");
// 		return;
// 	}
// 	fileType := "";
// 	switch lang {
// 	case "python":
// 		fileType = ".py"
// 	case "php":
// 		fileType = ".php"
// 	case "nodejs":
// 		// Don't know nodejs filetype yet
// 		return;
// 	}

// 	numUsers := getFlag("-numUsers", "1");
// 	form.Add("numUsers", numUsers);

// 	serverData, err := parseLoadTestJSON("serverData.json");
// 	if err != nil {
// 		fmt.Println("Error parsing serverData JSON");
// 		return;
// 	}
// 	for i := 0; i < len(serverData); i++ {
// 		form.Add(serverData[i][0], serverData[i][1]);
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	err = addFieldsToBody(form, writer);
// 	if err != nil {
// 		fmt.Println("Error adding form fields:", err);
// 		return;
// 	}

// 	filePath := getFileArg(fileType)
// 	if (filePath == "") {
// 		fmt.Println("Please provide a custom loadtest file");
// 		return;
// 	}

// 	err = addFileToBody(filePath, writer, "file");
// 	if err != nil {
// 		return;
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		fmt.Println("Error closing multipart form:", err)
// 		return
// 	}
// }

//////////////////////
func jmeterLoadTest() {
	var jmeterSingleValueFlags []string = []string{
		"name",
		"desc",
		"numServers",
		"version",
		"storeOutput",
		"webdriver-width",
		"webdriver-height",
		"webdriver-depth",
		"opts",
	}
	// remaining jmeter flags
	// servers
	// plugins
	// opts
	// jvm_args
	// [plugin-name]_[KEY]

	bodyData = make(map[string]interface{});

	data, err := parseTestJSON(defaultConfigPath);
	if err != nil {
		fmt.Println("Error parsing JSON");
		return;
	}
	for key, value := range data {
		bodyData[key] = value;
	}

	jsonPath := getFlag("-cfg", "");
	if jsonPath != "" {
		data, err := parseTestJSON(jsonPath);
		if err != nil {
			fmt.Println("Error parsing JSON");
			return;
		}
		for key, value := range data {
			bodyData[key] = value;
		}
	}

	filePath := getFileArg(".jmx");
	if filePath == "" {
		fmt.Println("Please provide a Jmeter test file");
		//return;
	} else {
		bodyData["file"] = filePath;
	}

	for _, flag := range jmeterSingleValueFlags {
		if flagValue := getFlag(fmt.Sprintf("-%s", flag), ""); flagValue != "" {
			bodyData[flag] = flagValue;
		}
	}

	servers := getFlag("-servers", "");
	if servers != "" {
		var serverMapArray []map[string]string;
		err := json.Unmarshal([]byte(servers), &serverMapArray);
		if err != nil {
			fmt.Println("Error parsing value for servers from CLI", err);
		} else {
			count := 0;
			for _, item := range serverMapArray {
				for innerKey := range item {
					sKey := fmt.Sprintf("servers[%d][%s]", count, innerKey);
					//fmt.Println(sKey, item[innerKey]);
					bodyData[sKey] = item[innerKey];
				}
				count++;
			}
		}
	}

	plugins := getFlag("-plugins", "");
	if plugins != "" {
		var pluginsMapArray []map[string]json.RawMessage;
		var pluginName string;
		var pluginOptions map[string]string;
		err := json.Unmarshal([]byte(plugins), &pluginsMapArray);
		if err != nil {
			fmt.Println("Error parsing value for plugins from CLI", err);
		} else {
			for _, item := range pluginsMapArray {
				for key, value := range item {
					if key == "plugin" {
						err := json.Unmarshal(value, &pluginName);
						if err != nil {
							fmt.Println("Error:", err);
							return;
						}
						addToInterfaceMap("plugin[]", pluginName, &bodyData);
					} else if key == "options" && pluginName != "" {
						err := json.Unmarshal(value, &pluginOptions);
						if err != nil {
							fmt.Println("Error:", err);
							return;
						}
						for innerKey, innerValue := range pluginOptions {
							sKey := fmt.Sprintf("%s_%s", pluginName, innerKey);
							addToInterfaceMap(sKey, innerValue, &bodyData);
						}
					}
				}
			}
		}
	}

	jvm_args := getMultiFlag("-jvm_args");
	if len(jvm_args) != 0 {
		sValue := strings.Join(jvm_args, " ");
		bodyData["jvm_args"] = sValue;
	}

	extras := getMultiFlag("-extras");
	if len(extras) != 0 {
		for _, file := range extras {
			addToInterfaceMap("extras[]", file, &bodyData);
		}
	}

	split := getMultiFlag("-split");
	if len(split) != 0 {
		sValue := strings.Join(split, " ");
		bodyData["split[]"] = sValue;
	}

	//fmt.Println(bodyData);

	body := &bytes.Buffer{};
	writer := multipart.NewWriter(body);

	writer.WriteField("testType", "jmeter-test");

	writeFromBodyData(writer, bodyData);

	err = writer.Close();
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}

	responseData := httpPostRequest(body, writer.FormDataContentType());

	if getFlagExist("-o") {
		var idStorage map[string]float64;
		err = json.Unmarshal(responseData, &idStorage);
		if err != nil {
			println("Error parsing http response", err);
		} else {
			// Ensure that our redirect to webpage doesn't beat its creation
			timeSeconds := (1/4) * time.Second;
			time.Sleep(timeSeconds);

			id := fmt.Sprintf("%.0f", idStorage["loadTestId"]);
			redirectBrowserToUrl(id);
		}
	}
}

// http requests//
func httpPostRequest(body *bytes.Buffer, content string) []byte {
	client := http.Client{};

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/Api/LoadTest", build), body);
	if err != nil {
		fmt.Println("Error creating request: " + err.Error());
		return nil;
	}
	req.Header.Set("X-Redline-Auth", getAPIKEY());
	req.Header.Add("Content-Type", content);

	resp, err := client.Do(req);
	if err != nil {
		fmt.Println("Error sending request: " + err.Error());
		return nil;
	}
	defer resp.Body.Close();

	fmt.Println("Response status:", resp.Status);

	responseBody, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		fmt.Println("Error reading response body: " + err.Error());
		return nil;
	}

	return responseBody;
}

func redirectBrowserToUrl(id string) {
	url := fmt.Sprintf("%s/LoadTest/View/%s", build, id);

	var cmd *exec.Cmd;
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url);
	case "darwin":
		cmd = exec.Command("open", url);
	case "linux":
		cmd = exec.Command("xdg-open", url);
	default:
		println("Unsupported operating system: %s", runtime.GOOS);
		return;
	}

	err := cmd.Start();
	if err != nil {
		println("Error executing redirect through exec", err);
		return;
	}
}

func testLoadTest() {
	fmt.Println("test loadtest function");
}

//_______________________________________________//
//Miscellaneous//

func printLoadTestInfo() {
	// fmt.Println("	loadTest - Starts a load test");
	// fmt.Println("	    Simple Test");
	// fmt.Println("	    JMeter Test");
	// fmt.Println("	    LogFileReplay Test");
	// fmt.Println("	    Gatling Test");
	// fmt.Println("	    Custom Test");
	fmt.Println("Usage:");
	fmt.Println("    redline13 run [testType/fileForTest] [flags]");
	fmt.Println("\nTest Types:");
	fmt.Println("    jmeter - can ommit by providing .jmx file type");
	fmt.Println("    More to be added...");
	fmt.Println("\nFlags:")
	fmt.Println("    -cfg - Takes additional .json file to overwrite values in config")
	fmt.Println("    -name - Name of loadtest")
	fmt.Println("    -desc - Description of loadtest")
	fmt.Println("    -version - jmeter version of test to run");
	fmt.Println("    -numServers - number of servers to run test on");
	fmt.Println("    -storeOutput - Boolean, of whether test output should be saved");
	fmt.Println("    -webdriver-width - width of screen for simulated browser");
	fmt.Println("    -webdriver-height - height of screen for simulated browser");
	fmt.Println("    -webdriver-depth - screen depth for simulated browser");
	fmt.Println("    -opts - Specify JMeter options as string \"-Jkey=value -Jkey=value\"");
	fmt.Println("    -servers - Specify servers as Array of Json objects [{\"size\":\"m5.large\", \"location\":\"us-east-1\"}, {}]");
	fmt.Println("    -plugins - Specify plugins as Array of Json objects [{\"plugin\": \"myPlugin\", \"options\": {\"myOption\": \"optionValue\", \"\": \"\"}}]");
	fmt.Println("    -jvm_args - Specify JVM Options such as \"Xms256m Xmx256m\"");
	fmt.Println("    -extras - Extra file(s) to be included in loadtest");
	fmt.Println("    -split - Filename(s) will be split across all the test servers, Usually used for splitting CSV files.");
	fmt.Println("    -o - follows loadtest to browser");
	fmt.Println("\nExamples:");
	fmt.Println("    redline13 run test.jmx -cfg myConfig.json -name CLILoadTest -desc \"my desc\" -extras extra1.csv extra2.csv");
}

func parseTestJSON(path string) (map[string]interface{}, error) {
	var ret map[string]interface{};
	ret = make(map[string]interface{});

	jsonData, err := ioutil.ReadFile(path);
	if err != nil {
		fmt.Println("Error reading JSON file:", err);
		return nil, err;
	}

	var data map[string]json.RawMessage;
	err = json.Unmarshal([]byte(jsonData), &data);
	if err != nil {
		fmt.Println("Error:", err);
		return nil, err;
	}

	for key, value := range data {
		switch key {
		case "servers":
			var serverMapArray []map[string]string;
			err := json.Unmarshal(value, &serverMapArray);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				count := 0;
				for _, item := range serverMapArray {
					for innerKey := range item {
						sKey := fmt.Sprintf("%s[%d][%s]", key, count, innerKey);
						//fmt.Println(sKey, item[innerKey]);
						ret[sKey] = item[innerKey];
					}
					count++;
				}
			}
		case "extras":
			var extrasArray []string;
			err := json.Unmarshal(value, &extrasArray);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				for _, value := range extrasArray {
					addToInterfaceMap("extras[]", value, &ret);
				}
			}
		case "split":
			var splitArray []string;
			err := json.Unmarshal(value, &splitArray);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				sValue := strings.Join(splitArray, " ");
				ret["split[]"] = sValue;
			}
		case "jvm_args":
			var argsArray []string;
			err := json.Unmarshal(value, &argsArray);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				sValue := strings.Join(argsArray, " ");
				ret[key] = sValue;
			}
		case "plugins":
			var pluginsMapArray []map[string]json.RawMessage;
			var pluginName string;
			var pluginOptions map[string]string;
			err := json.Unmarshal(value, &pluginsMapArray);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				for _, item := range pluginsMapArray {
					//pluginName = string(item["plugin"])
					for key, value := range item {
						if key == "plugin" {
							err := json.Unmarshal(value, &pluginName);
							if err != nil {
								fmt.Println("Error:", err);
								return nil, err;
							}
							addToInterfaceMap("plugin[]", pluginName, &ret);
						} else if key == "options" && pluginName != "" {
							err := json.Unmarshal(value, &pluginOptions);
							if err != nil {
								fmt.Println("Error:", err);
								return nil, err;
							}
							for innerKey, innerValue := range pluginOptions {
								sKey := fmt.Sprintf("%s_%s", pluginName, innerKey);
								addToInterfaceMap(sKey, innerValue, &ret);
							}
						}
					}
				}
			}
		case "opts":
			var optsMap map[string]string;
			var optsArgs []string;
			err := json.Unmarshal(value, &optsMap);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				for key, value := range optsMap {
					optsArgs = append(optsArgs, fmt.Sprintf("%s=%s", key, value));
				}
				ret[key] = strings.Join(optsArgs, " ");
			}
		default:
			if key == "apikey" {
				break;
			}
			var sValue string;
			err := json.Unmarshal(value, &sValue);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			} else {
				ret[key] = sValue;
			}
		}
	}

	return ret, nil;
}

func parseLoadTestJSON(path string) ([][]string, error) {
	jsonData, err := ioutil.ReadFile(path);
	if err != nil {
		fmt.Println("Error reading JSON file:", err);
		return nil, err;
	}

	var data map[string]string;
	err = json.Unmarshal(jsonData, &data);
	if err != nil {
		fmt.Println("Error parsing JSON:", err);
		return nil, err;
	}

	var arr [][]string;
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

func addFileToBody(filePath string, writer *multipart.Writer, key string) error {
	//filePath := "/Downloads/MyJmeterTest.jmx"
	file, err := os.Open(filePath);
	if err != nil {
		fmt.Println("Error opening file:", err);
		return err;
	}
	defer file.Close();

	// Create the file part in the multipart form
	part, err := writer.CreateFormFile(key, filepath.Base(filePath));
	if err != nil {
		fmt.Println("Error creating form file part:", err);
		return err;
	}

	// Copy the file contents to the file part
	_, err = io.Copy(part, file);
	if err != nil {
		fmt.Println("Error copying file contents to part:", err);
		return err;
	}
	return nil;
}

func isFile(path string) bool {
	info, err := os.Stat(path);
	if err != nil {
		if os.IsNotExist(err) {
			return false; // File does not exist
		}
		fmt.Println("Could not determine if value is file:", err);
		return false;
	}
	return !info.IsDir();
}

func getFileArg(fileType string) string {
	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], fileType) && isFile(args[i]) {
			return args[i];
		}
	}
	return "";
}

// takes a map[string]interface{} and a writer, writes to body based on interface type
func writeFromBodyData(writer *multipart.Writer, data map[string]interface{}) {
	writeKeyValuePair := func(writer *multipart.Writer, key string, value string) {
		if isFile(value) {
			addFileToBody(value, writer, key);
		} else {
			writer.WriteField(string(key), string(value));
		}
	}
	for key, value := range data {
		switch value.(type) {
		case []string:
			for _, val := range value.([]string) {
				writeKeyValuePair(writer, key, val);
			}
		case string:
			writeKeyValuePair(writer, key, value.(string));
		}
	}
}

// takes a key, value, and *map[stirng]interface{} and adds value based on key state (use if key might already exist and multiple are needed)
func addToInterfaceMap(key string, value string, Map *map[string]interface{}) error {
	if interfaceValue, ok := (*Map)[key]; ok {
		switch interfaceValue.(type) {
		case []string:
			(*Map)[key] = append((*Map)[key].([]string), value);
		case string:
			oldValue := (*Map)[key].(string);
			(*Map)[key] = []string{oldValue, value};
		default:
			return errors.New(fmt.Sprintf("Unexpected value type with key: %s", key));
		}
	} else {
		(*Map)[key] = value;
	}
	return nil;
}
