# lma-dispatcher

Go backend microservice, queue of the tasks, load balancer.
Implements parallel mode of the requests to different endpoints.

![](https://img.shields.io/github/languages/code-size/postnikovmu/lma-dispatcher)
![](https://img.shields.io/github/directory-file-count/postnikovmu/lma-dispatcher)
![](https://img.shields.io/github/languages/count/postnikovmu/lma-dispatcher)
![](https://img.shields.io/github/languages/top/postnikovmu/lma-dispatcher)

## General Info

This is a microservice that can be used to aggregate data from other microservices,
parsers, analyzers.
It implements a queue of tasks and parallel mode of the requests to different endpoints.

Example of request (local usage):
http://localhost:8080/?area=<Example_City>&text=<Example_developer>

Responce is the data in JSON format.

## Install

### Prerequisites:
microservices: \
./lma-extractor-hh \
(...lma-analyzer is not on github yet...)
### Deploy:
Deploy is on SAP BTP CloudFoundry platform: the settings are in manifest.yml. (or another CloudFoundry)

## Technologies

Go (Golang)
