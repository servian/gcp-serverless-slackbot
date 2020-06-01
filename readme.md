# GCP Serverless Slackbot Template

Minimalistic code to create a serverless slash handler for a Slack bot. Further information on this repository can be found in the article [Google Cloud Platform - Serverless Slack Bot](https://medium.com/weareservian/google-cloud-platform-serverless-slack-bot-c3b3d1c43330)

## Directory contents

### [v_alpha](v_alpha/)
Initial functional version of the code included for completeness

### [v1](v1/)
Pure `python` implementation - functions as expected, but presents the issue related to cold starts

### [v2](v2/)
Implemented in `go` and `python` to alleviate the cold start issue
