package main

import (
	"fmt";
	"os";
	"io/ioutil";
	"strings";
	"encoding/json";
)

var production string = "https://www.redline13.com"
var localHost string = "http://localhost";
var build string = production;

var defaultConfigPath string = "config.json";

var shortCallTestTypes []string = []string{"simple", "jmeter", "logfile", "custom", "test"};

var args = os.Args;

//_______________________________________________//
//Entry Point//

func main() {
	//displayArgs(args);

	//fmt.Println(len(args))


	//switch for first arg parameter,
	//could be loadTest, viewTest or apikey

	argument := "";
	if (len(args) > 1) {
		argument = args[1];
	}

	switch argument {
	case "loadTest", "simple", "jmeter", "logifle", "custom", "test" :
		loadTest();
	case "viewTest":
		viewTest();
	case "statsDownload":
		statsDownload();
	case "apikey":
		apiKey();
	default:
		if (getAPIKEY() == "") {
			fmt.Println("*important* You have no saved apikey, please visit https://www.redline13.com/Account/apikey to generate an api key");
		}
		fmt.Println(
			" Usage:\n",
			"	Redline [arguments]\n",
			"Arguments:",
		);
		printAPIKEYInfo();
		printLoadTestInfo();
		printViewTestInfo();
		printStatsDownloadInfo();
	}
}


//_______________________________________________//
//Core Functions//

func apiKey() {
	apikey := getAPIKEY();
	noArg := len(args) < 3;
	if (noArg && apikey == "") {
		fmt.Println("You have no saved apikey, please visit https://www.redline13.com/Account/apikey to generate an api key");
	} else if flag := getFlag("-set", ""); flag != "" {
		key := setAPIKEYJson(flag);
		fmt.Println("Key set: " + key);
	} else if getFlagExist("-show") {
		fmt.Println("Key: " + apikey);
	} else {
		printAPIKEYInfo();
	}
}

func loadTest() {
	// Create and handle loadTest
	argument := args[1];
	shortCall := (argument == "simple" || argument == "jmeter" || argument == "logfile" || argument == "custom" || argument == "test");
	handleLoadTest(shortCall);
}

func viewTest() {
	// Create and handle viewTest
	handleViewTest();
}

func statsDownload() {
	// Create and handle statsDownload
	handleStatsDownload();
}


//_______________________________________________//
//Miscellaneous//

func printAPIKEYInfo() {
	fmt.Println("	apikey - Set/Display your API key");
	fmt.Println("	    Flags:");
	fmt.Println("	        -set {your apikey} : Sets your API key");
	fmt.Println("	        -show : Displays API key");
}

func setAPIKEY(apikey string) string {
	err := ioutil.WriteFile("key.txt", []byte(apikey), 0644);
	if err != nil {
		fmt.Println(err);
		return "";
	}
	return apikey;
}

// func getAPIKEY() string {
// 	content, err := ioutil.ReadFile("key.txt");
// 	if err != nil {
// 		fmt.Println(err);
// 		return "";
// 	}
// 	return string(content);
// }

func setAPIKEYJson(apikey string) string {
	jsonData, err := ioutil.ReadFile(defaultConfigPath);
	if err != nil {
		fmt.Println("Error reading JSON file:", err);
		return "";
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return "";
	}

	data["apikey"] = apikey;

	updatedJson, err := json.Marshal(data)
	if err != nil {
		return "";
	}

	err = ioutil.WriteFile(defaultConfigPath, updatedJson, 0644)
	if err != nil {
		return "";
	}

	return apikey;

}

func getAPIKEY() string {
	var path string = defaultConfigPath;
	var apikey string = "";

	jsonData, err := ioutil.ReadFile(path);
	if err != nil {
		fmt.Println("Error reading JSON file:", err);
		return "";
	}

	var data map[string]json.RawMessage;
	err = json.Unmarshal([]byte(jsonData), &data);
	if err != nil {
		fmt.Println("Error:", err);
		return "";
	}
	
	for key, value := range data {
		if (key == "apikey") {
			err := json.Unmarshal(value, &apikey);
			if err != nil {
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err)
			}
		}
	}
	fmt.Println()
	return apikey;
}

func displayArgs() {
	for i := 1; i < len(args); i++ {
		fmt.Println(args[i]);
	}
}


//Flag Functions//
func getFlag(flag string, defaultFlag string) string {
	for i := 0; i < (len(args) - 1); i++ {
		if (args[i] == flag) {
			return args[i+1];
		}
	}
	return defaultFlag;
}

func getFlagExist(flag string) bool {
	for i := 0; i < (len(args)); i++ {
		if (args[i] == flag) {
			return true;
		}
	}
	return false;
}

func getMultiFlag(flag string) []string {
	ret := []string{};

	active := false
	for i := 0; i < len(args); i++ {
		//fmt.Println(args[i]);
		typeFlag := (args[i] == flag);
		if (typeFlag) {
			active = true;
			//fmt.Println("actived");
		}
		isFlag := strings.Contains(args[i], "-")
		if (!isFlag && active) {
			ret = append(ret, args[i]);
		} else if (isFlag && !typeFlag) {
			active = false;
			//fmt.Println("unactived");
		}
	}
	return ret;
}

