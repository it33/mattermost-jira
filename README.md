# JIRA Mattermost Webhook Connector

## For Mattermost administrator
You need incoming webhook URL in `http://<mattermost_server>/hooks/<web_hook_id>` format. Can be copied from Mattermost config.

 - System console
 - INTEGRATIONS → Custom Integrations
    - Enable Incoming Webhooks: true
    - Enable integrations to override usernames: true
    - Enable integrations to override profile picture icons: true
    
 - Team menu (3 dots near the Team name in top-left corner at the team-screen)
 - Integrations → Incoming Webhooks
 - Add Incoming Webhook
 
## For JIRA administrator
 - JIRA Administration → System
 - ADVANCED → WebHooks
 - Create a WebHook:
    - URL:  https://_**heroku_app_name**_.herokuapp.com?mattermost_hook_url=_**mattermost_hook_url**_
    - Issue:
        - created: true
        - updated: true
        - deleted: true

## Other info
        