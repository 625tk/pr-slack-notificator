# pr-slack-notificator

PRのdescription中の最初に見つかった引用
\`\`\`で囲まれたもの
をslackに通知する君

必要な環境変数
|env name|内容|example|note|
| ------------- | ------------- | ------------- | ------------- |
|REPOSITORY|リポジトリ|625tk/pr-slack-notificator||
|GIT_PR_RELEASE_TOKEN|githubへのアクセストークン|xxxxxxxxxx||
|PR_NUMBER|prの番号|12||
|SLACK_WEBHOOK_URL|slackのwebhook url|https://hooks.slack.com/xxxxxxxxxxxxxx||
|SLACK_CHANNEL|飛ばしたい先のchannel|test_625tk||
|SLACK_API_TOKEN|slack api使うようのtoken|xoxb-aaaaaaaa|SLACK_WWBHOOK_URLとどちらか|


