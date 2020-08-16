# Devops Task

## What is it

slack chatbot

## Available Commands
```console
/k-bot pods
```
<img src="./resources/logs-command.png"  width="300"> 
<img src="./resources/logs-response.png"  width="800"> 

```console
/k-bot logs [service] [tail]
```
<img src="./resources/pods-command.png"  width="300"> 
<img src="./resources/pods-response.png"  width="800"> 

## Installation 
Via the [helm chart](./chart)
in `values.yaml`, you will have to provide the ingress host name, and slack signing secret.

## Meeting Task Requirements
> 1.Production-readiness: code should be reliable, tested and clean.

The code is separated into small reusable functions and follows clean code principles  

> 2.Developer Experience (DX): deliver easy-to-use, self-service experience.

The code uses markdown to format answers for easy usage

> 3.Security.

The code uses slack built-in token-based authentication to secure the communication using a [signing secret](https://api.slack.com/authentication/verifying-requests-from-slack#about)

>4.Observability: easily investigate and learn how the bot is being used.

- Various log levels (`debug`, `info`, `warn` and `error`)
- Clear error messages
- usage metrics exported in prometheus-readable format

## Projects Steps
- [X] k8s client: 
  - get name, age, and /version of each running pod
  - get x log lines
- [X] slack server
  - accept slash command
- [X] expose metrics
  - requests counter
- [X] tests
- [X] errors handling
- [X] logs
