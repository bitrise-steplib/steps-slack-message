# Send a Slack message

[![Step changelog](https://shields.io/github/v/release/bitrise-io/steps-slack-message?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/steps-slack-message/releases)

Send a [Slack](https://slack.com/) message to a channel or group.

<details>
<summary>Description</summary>


Send a [Slack](https://slack.com/) message to a Slack user, a group, or a channel. Create and customize the messages however you see fit. Among other things, you can:

- Set a different text for failed and successful builds.
- Add an icon and/or emojis to your messages. 
- Set the bot user's name for the messages.
- Linkify channel names and usernames.
- Add and customize attachments. 

### Configuring the Step 

To use this Step, you need either a configured Slack Integration in your workspace, an incoming Slack webhook or a Slack bot user with an API token. For the former see your Workspace settings, for the latter two, you can set them up in Slack:

- [Incoming webhooks](https://api.slack.com/incoming-webhooks).
- [Bot user with an API token](https://api.slack.com/bot-users).

Once you're ready with those, come back to Bitrise and configure the Step itself:

1. Create a [Secret Env Var](https://devcenter.bitrise.io/builds/env-vars-secret-env-vars/) for either your Slack webhook URL or your Slack API token.
1. Add the Secret to either the **Slack Webhook URL** or the **Slack API token** input.
1. Toggle the **Run if previous Step failed** option on - you should see a white checkmark on green background next to it. This allows Slack messages to be sent for failed builds, too.
1. In the **Target Slack channel, group or username**, set where the Slack message should be sent. 
1. Customize your messages as you'd like. For the details, see the respective inputs.


In case of the Slack Integration usecase you can copy the ID in your Workspace settings, on the Integrations page. This ID is not senstive, you can use it as a step input as-is, or put it into a regular environment variable.

Note that this step always sends a message (either to `channel` or `channel_on_error`). If your use case is to send a message only on success or on failure, then you can [run the entire step conditionally](https://devcenter.bitrise.io/en/steps-and-workflows/introduction-to-steps/enabling-or-disabling-a-step-conditionally.html).

### Troubleshooting 

If the Step fails, check your Slack settings, the incoming webhook or the API token, and your Slack channel permissions. 

### Useful links 

- [Integrating with Slack](https://devcenter.bitrise.io/builds/configuring-notifications/#integrating-with-slack)
- [Slack attachments](https://api.slack.com/messaging/composing/layouts#attachments)

### Related Steps 

- [Send email with Mailgun](https://www.bitrise.io/integrations/steps/email-with-mailgun)
- [Post Jira Comment](https://www.bitrise.io/integrations/steps/post-jira-comment-with-build-details)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

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

To register an Incoming WebHook integration visit: https://api.slack.com/incoming-webhooks.

```yaml
steps:
- slack:
    title: Notify team
    inputs:
    - webhook_url: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
    - message: This is a test notification, please ignore
    - emoji: ":bitrise:"
```

#### Using Block Kit

Please check the format guideline [https://api.slack.com/methods/chat.postMessage#arg_blocks](https://api.slack.com/methods/chat.postMessage#arg_blocks)

```yaml
steps:
- slack:
    title: Notify team
    inputs:
    - workspace_slack_integration_id: example
    - message: This is a test notification, please ignore
    - blocks: |-
            [
                {
                    "type": "section",
                    "text": {
                        "type": "plain_text",
                        "text": "Hello world"
                    }
                }
            ]

```

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `is_debug_mode` | Step prints additional debug information if this option is enabled  |  | `no` |
| `webhook_url` | **One of workspace\_integration\_id, webhook\_url or api\_token input is required.** To register an **Incoming WebHook integration** visit: https://api.slack.com/incoming-webhooks  | sensitive |  |
| `webhook_url_on_error` | **One of workspace\_integration\_id, webhook\_url or api\_token input is required.** To register an **Incoming WebHook integration** visit: https://api.slack.com/incoming-webhooks  | sensitive |  |
| `workspace_slack_integration_id` | **One of workspace\_integration\_id, webhook\_url or api\_token input is required.** To register a **Workspace Slack Integration** see the Integration page in your Workspace settings  |  |  |
| `workspace_slack__integration_id_on_error` | **One of workspace\_integration\_id, webhook\_url or api\_token input is required.** To register a **Workspace Slack Integration** see the Integration page in your Workspace settings  |  |  |
| `api_token` | **One of workspace\_integration\_id, webhook\_url or api\_token input is required.**  To setup a **bot with an API token** visit: https://api.slack.com/bot-users  | sensitive |  |
| `channel` | Can be an encoded ID, or the channel's name.  Examples:  * channel ID: C024BE91L  * channel: #general  * username: @username  |  |  |
| `channel_on_error` | * channel example: #general * username example: @username  |  |  |
| `text` | Text of the message to send. Required unless you wish to send attachments only.  |  |  |
| `text_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  |  |
| `emoji` | Optionally you can specify a Slack emoji as the sender icon. You can use the Ghost icon for example if you specify `:ghost:` here as an input. **If you specify an Icon URL then this Emoji input will be ignored!**  |  |  |
| `emoji_on_error` | **This option will be used if the build failed.** If you leave this option empty then the default one will be used.  |  |  |
| `icon_url` | Optionally, you can specify a custom icon image URL which will be presented as the sender icon. Slack recommends an image a square image, which can't be larger than 128px in either width or height, and it must be smaller than 64K in size. Slack custom emoji guideline: [https://slack.zendesk.com/hc/en-us/articles/202931348-Using-emoji-and-emoticons](https://slack.zendesk.com/hc/en-us/articles/202931348-Using-emoji-and-emoticons) If you specify this input, the **Emoji** input will be ignored!  |  | `https://github.com/bitrise-io.png` |
| `icon_url_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  | `https://github.com/bitrise-io.png` |
| `link_names` | Linkify names in the message such as `@slackbot` or `#random`  |  | `yes` |
| `from_username` | The username of the bot user which will be presented as the sender of the message  |  | `Bitrise` |
| `from_username_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  | `Bitrise` |
| `thread_ts` | Sends the message as a reply to the message with the given ts if set (in a thread). |  |  |
| `thread_ts_on_error` | Sends the message as a reply to the message with the given ts if set (in a thread) if the build failed. |  |  |
| `ts` | Timestamp of the message to be updated.  When **Message Timestamp** is provided an existing Slack message will be updated, identified by the provided timestamp.   Example: `"1405894322.002768"`. |  |  |
| `ts_on_error` | Timestamp of the message to be updated if the build failed.  When **Message Timestamp if the build failed** is provided an existing Slack message will be updated, identified by the provided timestamp.   Example: `"1405894322.002768"`. |  |  |
| `reply_broadcast` | Used in conjunction with thread_ts and indicates whether reply should be made visible to everyone in the channel or conversation |  | `no` |
| `reply_broadcast_on_error` | Used in conjunction with thread_ts and indicates whether reply should be made visible to everyone in the channel or conversation |  | `no` |
| `color` | Color is used to color the border along the left side of the attachment. Can either be one of good, warning, danger, or any hex color code (eg. #439FE0). You can find more info about the color and other text formatting in [Slack's documentation](https://api.slack.com/docs/message-attachments).  | required | `#3bc3a3` |
| `color_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  | `#f0741f` |
| `pretext` | An optional text that appears above the attachment block. |  | `*Build Succeeded!*` |
| `pretext_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  | `*Build Failed!*` |
| `author_name` | A small text used to display the author's name. |  | `$GIT_CLONE_COMMIT_AUTHOR_NAME` |
| `title` | Title is displayed as larger, bold text near the top of a attachment. |  | `$GIT_CLONE_COMMIT_MESSAGE_SUBJECT` |
| `title_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  |  |
| `title_link` | A URL that will hyperlink the title. |  |  |
| `message` | Text is the main text of the attachment, and can contain standard message markup. The content will automatically collapse if it contains 700+ characters or 5+ linebreaks, and will display a "Show more..." link to expand the content.  |  | `$GIT_CLONE_COMMIT_MESSAGE_BODY` |
| `message_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  | `$GIT_CLONE_COMMIT_MESSAGE_BODY` |
| `image_url` | A URL to an image file that will be displayed inside the attachment.  Supported formats: GIF, JPEG, PNG, and BMP. Large images will be resized to a maximum width of 400px or a maximum height of 500px.  |  |  |
| `image_url_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  |  |
| `thumb_url` | A URL to an image file that will be displayed as a thumbnail on the right side of a attachment.  Supported formats: GIF, JPEG, PNG, and BMP. The thumbnail's longest dimension will be scaled down to 75px.  |  |  |
| `thumb_url_on_error` | This option will be used if the build failed. If you leave this option empty then the default one will be used.  |  |  |
| `footer` | The footer adds some brief text to help contextualize and identify an attachment.  Limited to 300 characters.  |  | `Bitrise` |
| `footer_on_error` | The footer adds some brief text to help contextualize and identify an attachment.  Limited to 300 characters.  |  | `Bitrise` |
| `footer_icon` | Renders a small icon beside the footer text It will be scaled down to 16px by 16px.  |  | `https://github.com/bitrise-io.png?size=16` |
| `footer_icon_on_error` | Renders a small icon beside the footer text It will be scaled down to 16px by 16px.  |  | `https://github.com/bitrise-io.png?size=16` |
| `timestamp` | Show the current time as part of the attachment's footer? |  | `yes` |
| `fields` | Fields separated by newlines and each field contains a `title` and a `value`. The `title` and the `value` fields are separated by a pipe `\|` character.  The *title* shown as a bold heading above the `value` text. The *value* is the text value of the field.  Supports multiline text with escaped newlines. Example: `Release notes\| - Line1 \n -Line2`.  Empty lines and lines without a separator are omitted.  |  | `App\|${BITRISE_APP_TITLE} Branch\|${BITRISE_GIT_BRANCH} Pipeline\|${BITRISEIO_PIPELINE_TITLE} Workflow\|${BITRISE_TRIGGERED_WORKFLOW_ID} ` |
| `buttons` | Buttons separated by newlines and each field contains a `text` and a `url`. The `text` and the `url` fields are separated by a pipe `\|` character. Empty lines and lines without a separator are omitted.  The *text* is the label for the button. The *url* is the fully qualified http or https url to deliver users to. An attachment may contain 1 to 5 buttons.  |  | `View App\|${BITRISE_APP_URL} View Pipeline Build\|${BITRISEIO_PIPELINE_BUILD_URL} View Workflow Build\|${BITRISE_BUILD_URL} Install Page\|${BITRISE_PUBLIC_INSTALL_PAGE_URL} ` |
| `pipeline_build_status` | This status will be used to help choosing between _on_error inputs and normal ones when sending the slack message.  |  | `$BITRISEIO_PIPELINE_BUILD_STATUS` |
| `build_status` | This status will be used to help choosing between _on_error inputs and normal ones.  |  | `$BITRISE_BUILD_STATUS` |
| `output_thread_ts` | Will export the created thread's timestamp to the environment with the supplied name (if not already in thread) |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/steps-slack-message/pulls) and [issues](https://github.com/bitrise-io/steps-slack-message/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
