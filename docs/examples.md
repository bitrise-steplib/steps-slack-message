### Examples

#### Using OAuth integration (recommended)

First, set up the integration between your Bitrise Workspace and your Slack instance. Go to Bitrise _Workspace Settings > Integrations_, select _Slack_, then follow the instructions.

Once the two services are connected, copy the integration ID (under the `...` menu) and pass it as an input to this step. This ID is not senstive, you can use it as a step input as-is, or put it into a regular environment variable.

```yaml
steps:
- slack:
    title: Notify team
    inputs:
    - workspace_slack_integration_id: example
    - message: This is a test notification, please ignore
    - emoji: ":bitrise:"
```

#### Using webhooks

TODO
