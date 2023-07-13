package main

import (
	"fmt";
	"os";
	"io/ioutil";
	"strings";
	"encoding/json";
	"os/exec";
	"runtime";
)

var production string = "https://www.redline13.com"
var localHost string = "http://localhost";
var build string = localHost;

var defaultConfigPath string;

var shortCallTestTypes []string = []string{"simple", "jmeter", "logfile", "custom", "test"};

var args = os.Args;

//_______________________________________________//
//Entry Point//

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Could not find user config directory");
		return;
	}
	defaultConfigPath = userConfigDir + "/Redline13/config.json";

	argument := "";
	if (len(args) > 1) {
		argument = args[1];
	}

	switch argument {
	case "run", "simple", "jmeter", "logifle", "custom" :
		loadTest();
	case "viewtest":
		viewTest();
	case "statsdownload":
		statsDownload();
	case "config":
		configEdit();
	case "help":
		commandHelp();
	case "version":
		printCLIVersion();
	case "test":
		// for testing
		fmt.Println(getAPIKEY());
	default:
		fmt.Println()
		if (getAPIKEY() == "") {
			fmt.Println("*important* You have no saved apikey, please visit https://www.redline13.com/Account/apikey to generate an api key");
		}
		fmt.Println("Usage: ");
		fmt.Println("    redline13 [command]");
		fmt.Println("Available Commands:");
		fmt.Println("    run - Run a load test on redline13");
		fmt.Println("    viewtest - View all tests or specific load test(s)");
		fmt.Println("    statsdownload - Download load test stats in CSV");
		fmt.Println("    config - Set up local config with API Key and defaults");
		fmt.Println("    version - Show CLI version information");
		fmt.Println("    help - [Command] show information about a command");
		fmt.Println("Use redline13 help [command] to show flags for command.");
	}
}


//_______________________________________________//
//Core Functions//

func configEdit() {
	openConfig := func(path string) {
		var cmd *exec.Cmd;
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", "", path);
		case "darwin":
			cmd = exec.Command("open", path);
		case "linux":
			cmd = exec.Command("xdg-open", path);
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
	_, err := os.Stat(defaultConfigPath);
	if os.IsNotExist(err) {
		createConfigFile();
	}
	if (getFlagExist("-show")) {
		openConfig(defaultConfigPath);
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

func commandHelp() {
	arguement := "";
	if (len(args) > 2) {
		arguement = args[2];
	}
	
	switch arguement {
	case "run":
		printLoadTestInfo();
	case "viewtest":
		printViewTestInfo();
	case "statsdownload":
		printStatsDownloadInfo();
	case "config":
		printConfigInfo();
	default:
		fmt.Println("Unknown command, cannot display help");
	}
}


//_______________________________________________//
//Miscellaneous//

func printCLIVersion() {
	//ask about current version
	fmt.Println("Current version");
}

func printConfigInfo() {
	fmt.Println("Usage:")
	fmt.Println("    redline13 config [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("    -show - Brings config file to focus on screen")
	fmt.Println("\nExamples:")
	fmt.Println("    redline13 config -show")
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
				fmt.Printf("Error parsing value for key '%s': %s\n", key, err);
			}
		}
	}
	return apikey;
}

func displayArgs() {
	for i := 1; i < len(args); i++ {
		fmt.Println(args[i]);
	}
}

func createConfigFile() {
	createDefaultData := func() (map[string]interface{}, error) {
		var data map[string]interface{} = make(map[string]interface{});
		serversArray := []map[string]string{{}};
		serversArray[0] = make(map[string]string);
		jsonStr := `{"location":"us-east-1", "num":"1", "onDemand":"T", "size":"m5.large", "subnetId":"subnet-00d66cc55ec4cf4bd", "usersPerServer":"1"}`
		err := json.Unmarshal([]byte(jsonStr), &serversArray[0])
		if err != nil {
			fmt.Println("Error creating default data for config file: ", err);
			return nil, err;
		}

		data["keyPairId"] = "Your_Key_Pair_Id";
		data["apikey"] = "Your_Api_Key";
		data["servers"] = serversArray;
		
		return data, nil;
	}

	dir, err := os.UserConfigDir(); err = os.Mkdir(dir + "/Redline13", 0777);
	if err != nil {
		fmt.Println("Could not find user config directory: ", err);
		return;
	}

	file, err := os.Create(defaultConfigPath);
	if err != nil {
		fmt.Println("Error creating config file: ", err);
		return;
	}

	data, err := createDefaultData();
	if err != nil {
		fmt.Println("Error creating default data for config", err);
		return;
	}

	json, err := json.Marshal(data);
	if err != nil {
		fmt.Println("Error marshalling Json: ", err);
		return;
	}

	err = ioutil.WriteFile(defaultConfigPath, json, 0644);
	if err != nil {
		fmt.Println("Error writing to config file: ", err);
		return;
	}

	defer file.Close();
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