# GCP Serverless Slackbot Template - version 2
Minimalistic code to create a serverless slash handler for a Slack bot. Further information on this repository can be found in the article [Google Cloud Platform - Serverless Slack Bot](https://medium.com/weareservian/google-cloud-platform-serverless-slack-bot-c3b3d1c43330). This version more gracefully handles cold starts to provide a more consistent acknowledgement response to Slack within the expected timeframs.

# Deploy
Both the `go` and `python` cloud functions need to be deployed to create the slash command handler.

`Python` deployment:

`gcloud functions deploy hello_bot --runtime python37 --trigger-http`

The `python` function permissions need to grant the `<project_id>@appspot.gserviceaccount.com` user the Cloud Functions Invoker role. This makes the function available to be invoked from another Cloud function.

`go` deployment:

`gcloud functions deploy Gobotween --runtime go111 --trigger-http`

The `go` function permissions will need to be updated to grant `allUsers` the `Cloud Functions Invoker` role. This makes the function publicly available to be invoked from Slack.

## Directory contents

### config.json
Configuration information used at runtime. Update the `SLACK_TOKEN` and `WEBHOOK_URL` here for your environment/project. Use the endpoint URL for the `python` Cloud Function to update the value for `SLASH_HANDLER`.

### gobotween.go
Main `go` code. Entry point is:

```go
func Gobotween(w http.ResponseWriter, r *http.Request) {
```

Execution is then:
 - return a `HTTP 200` response
 - load config if null (for a cold start)
 - validate request and signature
 - create `POST` request and forward to `SLASH_HANDLER`
