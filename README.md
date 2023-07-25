# Redline13 CLI
A simple CLI written in Golang for load testing through Redline13

## Description
The objective of this project was to create a versatile and powerful tool using Golang, designed as a standalone executable or CLI tool. This program enables users to execute JMeter, Gatling, or custom code load test plans at scale, all while benefiting from cost-effective instance pricing. The key to achieving this was by utilizing the capabilities of Redline13's API infrastructure. Through this API, users can effortlessly run load tests with clearly defined parameters and options, providing them with complete control and flexibility over their load testing process. 

## Installation
Simply fork or clone repo to desired location. Execute "go build -o redline" in the project directory to build a working executable. 

## Usage
    redline [command]
    
Available Commands:

    run - Run a load test on redline13
    viewtest - View all tests or specific load test(s)
    statsdownload - Download load test stats in CSV or compressed file formats
    delete - Delete or Cancel load test from loadTestId
    config - Set up local config with API Key and defaults
    version - Show CLI version information
    help - [Command] show information about a command
## Examples

    redline config -edit

    redline run test.jmx -cfg example.json -name MyCLILoadTest -desc "My test Description"
    redline run test.jmx -cfg example.json -name MyCLILoadTest -o
    redline run test.jmx -cfg example.json -servers '[{"location":"us-east-1","num":"15","onDemand":"T","size":"m5.large","subnetId":"subnet-00d66cc55ec4cf4bd","usersPerServer":"5"}]' -jvm_args arg1 arg2 -extras extra1.csv extra2.csv
    redline run gatlingTest.scala -cfg example.json -name MyCLIGatingTest -o
    redline run customTestCode.py -lang python -cfg example.json -name MyCustomLoadTest
    redline run customTestCode.js -lang nodejs -cfg example.json -name MyCustomLoadTest -o

    redline viewtest -id 123321

    redline statsdownload -id 123321 -type netIn netOut OutputFile1 OutputFile2

    redline help run
## About
This project was built by myself, [Mike Bugden](https://www.linkedin.com/in/mike-bugden-2b5196b0), during my summer internship as a developer at Redline13. The experience gained from this internship has been extremely valuable, as it allowed me to not only learn Golang as a new programming language, but to also expand my proficiency in other areas of software development. This project has also worked greatly to expose myself to the complex and fascinating world of load testing, and load testing software. Overall, this project has broadened my expertise in various aspects of software development, including Golang programming and the fascinating world of load testing.
