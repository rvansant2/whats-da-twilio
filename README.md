# What's da twilio

## Description
Small Golang HTTP Server that leverages Twilio's API for calling and SMS messaging

_**Note**: trial Twilio accounts and numbers can only send to verified numbers linked to the account_

## Required
- Go v1.15.x
- Twilio account and trial number
- ngrok for Twiml endpoint and default message and recording to play, where needed


## Get started
- Clone repo and run `cd whats-da-twilio`
- Then run `go mod tidy` to pull in all package depenendencies
- Copy or rename `default.example` as `default` and update/fill in the details of this configuration file with your Twilio account details.
- To start server, run `go run src/main.go` and go to your browser and type `http://localhost:3000/call` or `http://localhost:3000/sms`
- Verify a call or text message is made/sent and you hear the message recording or read the default text message sent

### To Do
- Add API Key verifcation
- Add REST-ful protocols to endpoints
- Add HTTP params support for SMS message bodies
