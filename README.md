# JIRA Mattermost Webhook Connector

## For Mattermost administrator
You need incoming webhook URL in `http://<mattermost_server>/hooks/<web_hook_id>` format. Can be copied from Mattermost config.

**Mattermost Config**
1. Go to `...` menu
2. `System console`
3. `Integrations` → `Custom Integrations`
4. Enable Settings
    - `Enable Incoming Webhooks`: `true`
    - `Enable integrations to override usernames`: `true`
    - `Enable integrations to override profile picture icons`: `true`

**Create Incoming Webhook**  
1. Go to `...` menu 
2. `Integrations` → `Incoming Webhooks`
3. `Add Incoming Webhook`
 
## For JIRA administrator
 - JIRA Administration → System
 - ADVANCED → WebHooks
 - Create a WebHook:
    - URL:  https://_**heroku_app_name**_.herokuapp.com?mattermost_hook_url=_**mattermost_hook_url**_
    - Issue:
        - created: true
        - updated: true
        - deleted: true

## Build Binary
`go build`

## Run
`./mattermost-jira -map=/absolute/path/to/mapping.json` map flag optional if you want room mapping
`./mattermost-jira -map=./mapping.json 2&1 >> data.log &`  pipe stderr and stdout to file and disown process

## Test
While server is running in background or different session, execute:
```
curl -X POST -H "Content-Type: application/json" --data @sample_hook.json localhost:5000?mattermost_hook_url=http://localhost:8065/hooks/bh9iwe5ezibepfcmqibgxpqs4c
```

## Room Map
The issue key will be split on the hypen and the project key will be looked up in a provided map.
example `JRA-121` -> `JRA`
example mapping.json

```
{
    "jiraprojectkey":"uchatroomname",
    "JRA":"off-topic"
}
```
If no key is found, the message will be sent to the default channel for the incoming webhook (configured when webhook is generated).
        