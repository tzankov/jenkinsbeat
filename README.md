# Jenkinsbeat

Welcome to Jenkinsbeat.

## Getting Started with Jenkinsbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Jenkinsbeat and also install the
dependencies, run the following command:

```
make setup
```


### Build

To build the binary for Jenkinsbeat run the command below. This will generate a binary
in the same directory with the name jenkinsbeat. Make sure to export the following variables and their respective required contents before building:

```
JENKINS_URL
```

```
JENKINS_USER
```

```
JENKINS_PASS
```

Then, run the build by running:

```
make
```


### Run

To run Jenkinsbeat with debugging output enabled, run:

```
./jenkinsbeat -c jenkinsbeat.yml -e -d "*"
```

To run Jenkinsbeat in standard mode, run:

```
./jenkinsbeat -e -d "*"
```
