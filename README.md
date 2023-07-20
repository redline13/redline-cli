# Redline13 CLI
A simple CLI written in Golang for load testing through Redline13

# Description
This CLI utilizes Redline13's API Infustructure for cloud based load testing to run JMeter, ~~Gatling~~, or ~~custom~~ ~~code~~ load test plans at scale using low cost instance pricing.

# Installation
Simply fork or clone repo to desired location. Execute "go build -o {alias}" in the project directory to build a working executable. 

# Usage
    redline [command]
    
Available Commands:

    run - Run a load test on redline13
    viewtest - View all tests or specific load test(s)
    statsdownload - Download load test stats in CSV or compressed file formats
    config - Set up local config with API Key and defaults
    version - Show CLI version information
    help - [Command] show information about a command

