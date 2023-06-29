package main

import (
	"fmt";
	"os";
	"io/ioutil";
	"strings";
)

var production string = "https://www.redline13.com"
var localHost string = "http://localhost";
var build string = production;

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
	case "loadTest":
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
	} else if (noArg){
		fmt.Println("Key: " + apikey)
	} else {
		key := setAPIKEY(args[2]);
		fmt.Println("Key set: " + key);
	}
}

func loadTest() {
	// Create and handle loadTest
	handleLoadTest();
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
	fmt.Println("	    [operation] {command}:");
	fmt.Println("	        [Set] Redline apikey {your apikey}");
	fmt.Println("	        [Show] Redline apikey");
}

func setAPIKEY(apikey string) string {
	err := ioutil.WriteFile("key.txt", []byte(apikey), 0644);
	if err != nil {
		fmt.Println(err);
		return "";
	}
	return apikey;
}

func getAPIKEY() string {
	content, err := ioutil.ReadFile("key.txt");
	if err != nil {
		fmt.Println(err);
		return "";
	}
	return string(content);
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

func getFlagExist(flag string, defaultFlag bool) bool {
	for i := 0; i < (len(args) - 1); i++ {
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

