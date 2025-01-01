# Sitemongoose

Sitemongoose is a simple, lightweight and zero-dependency site monitoring tool written in Go. It is useful for

- **Monitoring Website Availability**: check if your websites are up and running by sending periodic HTTP requests. You can define expected Status Codes and send alerts if the server's response differs.
- **Tracking Response Times**: measure the response time of your websites to ensure they are performing optimally, or send alerts otherwise.
- **Evaluating JSON Responses**: define rules to evaluate JSON responses from your APIs and trigger alerts based on conditions such as value comparisons or key existence.
- **Alerting**: Sitemongoose can send alerts based on specific conditions. The following alerting channels are currently available:
  - Email
  - [Slack](https://slack.com/intl/en-gb/)
  - [Pushover](https://pushover.net/)
  - Custom command: run a shell command in case of an alert is riggered

## Install

Sitemongoose is a single binary executable without any dependencies. Just download, extract and mark it as executable:

    wget https://github.com/neilsf/sitemongoose/releases/download/v0.1.0/sitemongoose-0.1.0_linux_x86_64.tar.gz
    tar -xzvf sitemongoose-0.1.0_linux_x86_64.tar.gz
    rm sitemongoose-0.1.0_linux_x86_64.tar.gz
    chmod +x sitemongoose
    ./sitemongoose --help

## Configure

To use Sitemongoose, you must define Monitors, Events and Alerts in a configuration file, in YAML format.

### Example configuration

```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 60
    timeout_ms: 5000
    events:
      - evaluate: status_code
        expected_status_code: 200
        alerts:
          - channel: email
            from: alerts@example.com
            to: admin@example.com
            alert_message: "Example.com is down!"
            resolution_message: "Example.com is back up."
      - evaluate: response_time
        expected_response_time_ms: 1000
        alerts:
          - channel: pushover
            alert_message: "Example.com is slow!"
            resolution_message: "Example.com is fast again."
      - evaluate: json_rule
        json_rule:
            json_path: "database.connections"
            condition: "lt"
            value: 500
        alerts:
          - channel: custom_cmd
            alert_command: "/path/to/alert_script.sh"
            resolution_command: "/path/to/resolution.sh"
```

## Run

To start monitoring, invoke Sitemongoose's start command and specify the location of the configuration file:

    ./sitemongoose start -c /path/to/config.yaml

In a production environment, you may want to run it as a service, using [systemd](https://systemd.io/) or [Supervisor](https://supervisord.org/).

## Complete Configuration Reference

A Sitemongoose configuration consists of three main building blocks:

1. The configuration must have one or more _Monitors_
2. A _Monitor_ has zero or more _Events_
3. An _Event_ has zero or more _Alerts_

### Monitors

A _Monitor_ is the top level building block of the configuration. It defines a service that runs in a loop, sending periodical HTTP requests to an URL and firing _Events_. You can define as many _Monitors_ as you wish. A _Monitor_ has the following configuration options:

- `name` (mandatory): an arbitrary string that is unique, i. e. no other monitors can have the same name.
- `url` (mandatory): the URL that the monitor will send HTTP requests to.
- `interval_sec` (mandatory): an integer denoting how many seconds must elapse between two requests.
- `timeout_ms` (optional): an integer denoting the time in milliseconds before the request is considered to time out. The default value is `30000`, that is, 30 seconds.
- `events` (optional): an array of _Event_s, see below.

### Events

After sending an HTTP request, the _Monitor_ passes the response to all of its _Events_. It's the _Event_'s responsibility to evaluate the response and decide whether it should take any actions. An _Event_ has the following configuration options:

- `evaluate` (mandatory) defines what to evaluate. Valid values are: `status_code`, `response_time` and `json_rule`. Other options differ based on the setting.

#### Status Code Evaluation

For `evaluate: status_code`, the only other option is `expected_status_code`, an integer that's matched against the response's status code. If the numbers mismatch, the _Event_'s all _Alerts_ will be fired.

The following example sends an alert if the response's status code does not equal `200` (note: _Alerts_ will be discussed later):

```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 10
    events:
      - evaluate: status_code
        expected_status_code: 200
        alerts: ...
```

#### Reponse Time Evaluation

For `evaluate: response_time`, the only other option is `expected_resposne_time_ms`. If the server's response time is greater than this number, or a timeout occurs, the alerts will be fired.

The following example sends an alert if the response time is greater than `500` millisecs:
```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 10
    events:
      - evaluate: response_time
        expected_response_time_ms: 500
        alerts: ...
```

**Note**: there must be at least one _Event_ with `evaluate: response_time`, otherwise nothing will be fired in case of a timeout.

#### JSON Rule Evaluation

Sitemongoose can take actions based on the response body the server provides. If `evaluate: json_rule`, another block called `json_rule` must be defined with the following options:

- `condition` (required): the conditions that the returned JSON must satisfy to avoid an alert. Valid options are:
  - `valid`: the returned response must hold a valid JSON, otherwise an alert will be sent
  - `exists`: the key defined by `json_path` must exist
  - `eq`: the numeric value found at `json_path` must equal to the value defined in `value`
  - `ne`: the numeric value found at `json_path` must not equal to the value defined in `value`
  - `lt`: the numeric value found at `json_path` must be less than to the value defined in `value`
  - `gt`: the numeric value found at `json_path` must be less than to the value defined in `value`
  - `regexp`: the string value found at `json_path` must match the regular expression defined in `value`
- `json_path` (mandatory except if `condition: valid`): see [JsonPath](https://goessner.net/articles/JsonPath/) for valid syntax
- `value` (mandatory except if `condition: valid`): the value can be of any type (a number, string or boolean) but certain conversions will take place when comparing against the response.
  - If checked for equality or inequality (`eq` or `ne`), values will be compared as strings. For example, the value `null` will equal `null` or `"null"`, but `123.0` will not equal `"123"`.
  - If checked with `lt` or `gt`, `value` must be numeric
  - If checked with `regexp`, `value` must be a valid regular expression. For reference, see [Go regexp syntax](https://pkg.go.dev/regexp/syntax).

In the following example, the event fires, if the value in the JSON response `{"database" : {"connections": <value>}}` is greater than `300`:

```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 10
    events:
      - evaluate: json_rule
        json_rule:
          condition: lt
          json_path: database.connections
          value: 300
        alerts: ...
```

### Alerts

_Events_ can have multiple _Alerts_ attached to them. When an event condition is met (that is, the response does not meet the expectations), the _Event_ becomes active and all its _Alerts_ are fired. When the condition is over, the _Alerts_ are fired again, but this time, sending a resolution message, letting the user know that the _Event_ is no longer active.

Configuration values:

- `channel` (mandatory): an _Alert_'s configuration depends on what type we're using. The following channels are currently supported: `email`, `slack`, `pushover` and `command`.

#### Email Alerts

For `channel: email`, the following values are additionally required:

- `from` (mandatory): the sender's email address
- `to` (mandatory): the recipient's email address
- `alert_message` (mandatory): the message to send when the _Alert_ is fired
- `resolution_message` (mandatory): the message to send when the _Event_ is no longer active

Additionally, the following environment variables must be send to be able to send out emails:

- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USER`
- `SMTP_PASS`

The following example will send an email if the value in the JSON response `{"database" : {"connections": <value>}}` is greater than `300`:

```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 10
    events:
      - evaluate: json_rule
        json_rule:
          condition: lt
          json_path: database.connections
          value: 300
        alerts:
          - channel: email
            from: alerts@example.com
            to: admin@example.com
            alert_message: "The number of database connections is above 300"
            resolution_message: "The number of database connections is back to normal"
```

#### Slack Alerts

Sitemongoose can send alerts to Slack using [webhooks](https://api.slack.com/messaging/webhooks). All you need to do is set up a webhook and store the webhook URL in the `SLACK_WEBHOOK_URL` environment variable.

For `channel: slack`, the following values are additionally required:

- `alert_message` (mandatory): the message to send when the _Alert_ is fired
- `resolution_message` (mandatory): the message to send when the _Event_ is no longer active

Additionally, the following environment variable must be send to be able to send out emails:

- `SLACK_WEBHOOK_URL`

#### Pushover Alerts

[Pushover](https://pushover.net/) is a service that you can use to send push messages to your Android phone or iPhone. Sitemongoose supports Pushover integration out of the box.

For `channel: pushover`, the following values are additionally required:

- `alert_message` (mandatory): the message to send when the _Alert_ is fired
- `resolution_message` (mandatory): the message to send when the _Event_ is no longer active

Additionally, the following environment variables must be send to be able to send out emails:

- `PUSHOVER_APP_TOKEN`
- `PUSHOVER_USER_KEY`

See the [Pushover API documentation](https://pushover.net/api) for more details.

#### Triggering Custom Commands

Sitemongoose can execute shell commands when an _Alert_ is fired. This is essentially the way to implement your own alerting mechanisms.

For `channel:command`, the following values are required:

- `alert_command` (required): an array of strings including the command and its arguments
- `resolution_command` (required): similar to the above, but is executed when the _Event_ becomes inactive

For example, let's say you run Sitemongoose in a desktop environment, and you want to get desktop notifications:

```yaml
monitors:
  - name: Example Monitor
    url: http://example.com/health.json
    interval_sec: 10
    events:
      - evaluate: status_code
        expected_status_code: 200
        alerts:
          - channel: command
            alert_command: ["notify-send", "-a", "Sitemongoose", "Example.com is down!"]
            resolution_command: ["notify-send", "-a", "Sitemongoose", "Example.com back online"]
```