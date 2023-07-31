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
)

var bodyData map[string]interface{};

//_______________________________________________//
//Entry Point//


func handleLoadTest() {
	argument := "";
	if len(args) > 2 {
		argument = args[2];
	}

	if (getFileArg(".jmx") != "") {
		argument = "jmeter";
	} else if (getFileArg(".scala") != "") || (getFileArg(".tar.gz") != "") || (getFileArg(".tar") != "") || (getFileArg(".jar") != "") {
		argument = "gatling"
	} else if (getFileArg(".php") != "") || (getFileArg(".py") != "") || (getFileArg(".js") != "") {
		argument = "custom"
	}

	switch argument {
	case "jmeter":
		jmeterLoadTest();
	case "gatling":
		gatlingLoadTest();
	case "custom":
		customLoadTest();
	case "test":
		testLoadTest();
	default:
		fmt.Println("Unknown test type provided");
	}
}

//_______________________________________________//
//Core Functions//


func customLoadTest() {
	var customFileTypes []string = []string{
		".php",
		".py",
		".js",
	}
	var customSingleValueFlags []string = []string {
		// Required
		"lang",
		"numUsers",
		// Optional
		"loadResources",
	}

	setInitialBodyData();

	fileSet := false;
	for _, fileType := range customFileTypes {
		if filePath := getFileArg(fileType); filePath != "" {
			bodyData["file"] = filePath;
			fileSet = true;
		}
	}
	if !fileSet {
		fmt.Println("please provide a valid Custom test");
		return;
	} 

	for _, flag := range customSingleValueFlags {
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

	body := &bytes.Buffer{};
	writer := multipart.NewWriter(body);

	writer.WriteField("testType", "custom-test");

	writeFromBodyData(writer, bodyData);

	err := writer.Close();
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}

	responseData, statusCode := httpPostRequest(body, writer.FormDataContentType());

	id := "";
	if (statusCode == 200) {
		var idStorage map[string]float64;
		err = json.Unmarshal(responseData, &idStorage);
		if err != nil {
			println("Error parsing http response", err);
		} else {
			id = fmt.Sprintf("%.0f", idStorage["loadTestId"]);
			fmt.Println("loadTestId:", id);
			if getFlagExist("-o") {
				redirectBrowserToUrl(id);
			}
		}
	}
}

func gatlingLoadTest() {
	var gatlingFileTypes []string = []string {
		".scala",
		".tar.gz",
		".tar",
		".jar",
	}
	var gatlingSingleValueFlags []string = []string{
		// Required
		"version",
		"numServers",
		// Optional
		"name",
		"desc",
		"storeOutput",
		"opts",
	}

	setInitialBodyData();

	fileSet := false 
	for _, fileType := range gatlingFileTypes {
		if filePath := getFileArg(fileType); filePath != "" {
			bodyData["file"] = filePath;
			fileSet = true
		}
	}
	if !fileSet {
		fmt.Println("please provide a valid Gatling test");
		return;
	} 

	for _, flag := range gatlingSingleValueFlags {
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

	body := &bytes.Buffer{};
	writer := multipart.NewWriter(body);

	writer.WriteField("testType", "gatling-test");

	writeFromBodyData(writer, bodyData);

	err := writer.Close();
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}

	responseData, statusCode := httpPostRequest(body, writer.FormDataContentType());

	id := "";
	if (statusCode == 200) {
		var idStorage map[string]float64;
		err = json.Unmarshal(responseData, &idStorage);
		if err != nil {
			println("Error parsing http response", err);
		} else {
			id = fmt.Sprintf("%.0f", idStorage["loadTestId"]);
			fmt.Println("loadTestId:", id);
			if getFlagExist("-o") {
				redirectBrowserToUrl(id);
			}
		}
	}
}

func jmeterLoadTest() {
	var jmeterSingleValueFlags []string = []string{
		// Required
		"numServers",
		"version",
		// Opitional
		"name",
		"desc",
		"storeOutput",
		"webdriver-width",
		"webdriver-height",
		"webdriver-depth",
		"opts",
	}

	setInitialBodyData();

	filePath := getFileArg(".jmx");
	if filePath == "" {
		fmt.Println("Please provide a Jmeter test file");
		return;
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

	body := &bytes.Buffer{};
	writer := multipart.NewWriter(body);

	writer.WriteField("testType", "jmeter-test");

	writeFromBodyData(writer, bodyData);

	err := writer.Close();
	if err != nil {
		fmt.Println("Error closing multipart form:", err);
		return;
	}

	responseData, statusCode := httpPostRequest(body, writer.FormDataContentType());

	id := "";
	if (statusCode == 200) {
		var idStorage map[string]float64;
		err = json.Unmarshal(responseData, &idStorage);
		if err != nil {
			println("Error parsing http response", err);
		} else {
			id = fmt.Sprintf("%.0f", idStorage["loadTestId"]);
			fmt.Println("loadTestId:", id);
			if getFlagExist("-o") {
				redirectBrowserToUrl(id);
			}
		}
	}
}

func httpPostRequest(body *bytes.Buffer, content string) ([]byte, int) {
	client := http.Client{};

	url := fmt.Sprintf("%s/Api/LoadTest", build);
	req, err := http.NewRequest(http.MethodPost, url, body);
	if err != nil {
		fmt.Println("Error creating request: " + err.Error());
		return nil, -1;
	}
	req.Header.Set("X-Redline-Auth", getAPIKEY());
	req.Header.Add("Content-Type", content);

	resp, err := client.Do(req);
	if err != nil {
		fmt.Println("Error sending request: " + err.Error());
		return nil, resp.StatusCode;
	}
	defer resp.Body.Close();

	
	responseBody, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		fmt.Println("Error reading response body: " + err.Error());
		return nil, resp.StatusCode;
	}

	if (resp.StatusCode != 200) {
		fmt.Println("Response status:", resp.Status);
		fmt.Println(string(responseBody));
	}

	return responseBody, resp.StatusCode;
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


// Must be called in all loadTests so bodyData map is made
func setInitialBodyData() {
	bodyData = make(map[string]interface{});

	data, err := parseLoadTestJSON(defaultConfigPath);
	if err != nil {
		fmt.Println("Error parsing JSON");
		return;
	}
	for key, value := range data {
		bodyData[key] = value;
	}

	jsonPath := getFlag("-cfg", "");
	if jsonPath != "" {
		data, err := parseLoadTestJSON(jsonPath);
		if err != nil {
			fmt.Println("Error parsing JSON");
			return;
		}
		for key, value := range data {
			bodyData[key] = value;
		}
	}
}

func printLoadTestInfo() {
	fmt.Println("Usage:");
	fmt.Println("    redline run [testType/fileForTest] [flags]");
	fmt.Println("\nTest Types:");
	fmt.Println("    jmeter - can ommit by providing .jmx file type");
	fmt.Println("    gatling - can ommit by providing .scala, .tar, .tar.gz, or .jar file type");
	fmt.Println("    custom - can ommit by providing .py, .php, or .js file type");
	fmt.Println("\nGlobal Flags:");
	fmt.Println("    -cfg - Takes additional .json file to overwrite values in config");
	fmt.Println("    -name - Name of loadTest");
	fmt.Println("    -desc - Description of loadTest")
	fmt.Println("    -storeOutput - Boolean, of whether test output should be saved");
	fmt.Println("    -servers - Specify servers as Array of Json objects [{\"size\":\"m5.large\", \"location\":\"us-east-1\"}, {}]");
	fmt.Println("    -plugins - Specify plugins as Array of Json objects [{\"plugin\": \"myPlugin\", \"options\": {\"myOption\": \"optionValue\", \"\": \"\"}}]");
	fmt.Println("    -extras - Extra file(s) to be included in loadTest");
	fmt.Println("    -split - Filename(s) will be split across all the test servers, Usually used for splitting CSV files.");
	fmt.Println("    -o - follows loadTest to browser");
	fmt.Println("\nJmeter Flags:");
	fmt.Println("    -version - jmeter version of test to run");
	fmt.Println("    -opts - Specify JMeter options as string \"-Jkey=value -Jkey=value\"");
	fmt.Println("    -jvm_args - Specify JVM Options such as \"Xms256m Xmx256m\"");
	fmt.Println("    -numServers - number of servers to run test on");
	fmt.Println("    -webdriver-width - width of screen for simulated browser");
	fmt.Println("    -webdriver-height - height of screen for simulated browser");
	fmt.Println("    -webdriver-depth - screen depth for simulated browser");
	fmt.Println("\nGatling Flags:");
	fmt.Println("    -version - gatling version of test to run");
	fmt.Println("    -opts - Specify Gatling options as string \"-Dkey=value -Dkey=value\"");
	fmt.Println("    -numServers - number of servers to run test on");
	fmt.Println("\nCustom Flags:");
	fmt.Println("    -lang - Language of loadTest");
	fmt.Println("    -loadResources - Should resources in the returned HTML be loaded? T or F");
	fmt.Println("    -numUsers - Number of users to simulate in the test on per server basis");
	fmt.Println("\nExamples:");
	fmt.Println("    redline run test.jmx -cfg myCustomConfig.json -name CLILoadTest -desc \"my desc\" -extras extra1.csv extra2.csv");
	fmt.Println("    redline run customTest.py -cfg myCustomConfig.json -name CLILoadTest -desc \"my desc\"");
	fmt.Println("    redline run gatlingTest.scala -cfg myCustomConfig.json -name CLILoadTest -version 3.3.1");
	fmt.Println("For further information on loadTest flags, visit http://redline13.com/ApiDoc/LoadTest/Post");
}

func parseLoadTestJSON(path string) (map[string]interface{}, error) {
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
	file, err := os.Open(filePath);
	if err != nil {
		fmt.Println("Error opening file:", err);
		return err;
	}
	defer file.Close();

	part, err := writer.CreateFormFile(key, filepath.Base(filePath));
	if err != nil {
		fmt.Println("Error creating form file part:", err);
		return err;
	}

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
