# GCP Serverless Slackbot Template

Minimalistic code to create a serverless slash handler for a Slack bot. Further information on this repository can be found in the article [Google Cloud Platform - Serverless Slack Bot](https://medium.com/tbc)

## Directory contents

### v_alpha
This is the initial functional version of the code included for completeness

### v1
This version of the bot is implemented completely in `python` and functions as expected, but presents the issue related to cold starts

### v2
This version of the bot is implemented in `go` and `python` to alleviate the cold start issue
