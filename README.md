# Sitemongoose

Sitemongoose is a simple, lightweight and zero-dependency site monitoring tool written in Go. It is useful for

- **Monitoring Website Availability**: check if your websites are up and running by sending periodic HTTP requests. You can define expected Status Codes and send alerts if the server's response differs.
- **Tracking Response Times**: measure the response time of your websites to ensure they are performing optimally, or send alerts otherwise.
- **Evaluating JSON Responses**: define rules to evaluate JSON responses from your APIs and trigger alerts based on conditions such as value comparisons or key existence.
- **Alerting**: Sitemongoose can send alerts based on specific conditions. The following alerting channels are currently available:
  - Email
  - Slack
  - Pushover
  - Custom command: run a shell command in case of an alert is riggered

## Install

Sitemongoose is a single binary executable without any dependencies. Just download, extract and mark it as executable:

    wget https://github.com/neilsf/sitemongoose/???
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
          - type: email
            from: alerts@example.com
            to: admin@example.com
            alert_message: "Example.com is down!"
            resolution_message: "Example.com is back up."
      - evaluate: response_time
        expected_response_time_ms: 1000
        alerts:
          - type: pushover
            alert_message: "Example.com is slow!"
            resolution_message: "Example.com is fast again."
      - evaluate: json_rule
        json_rule:
            json_path: "database.connections"
            condition: "lt"
            value: 500
        alerts:
          - type: custom_cmd
            alert_command: "/path/to/alert_script.sh"
            resolution_command: "/path/to/resolution.sh"
```

## Run

To start monitoring, invoke Sitemongoose's start command and specify the location of the configuration file:

    ./sitemongoose start -c /path/to/config.yaml

In a production environment, you may want to run it as a service, using systemd or Supervisor.

## Full Config Specification

A Sitemongoose configuration consists of three main building blocks:

1. The configuration must have one or more _Monitors_
2. A _Monitor_ has zero or more _Events_
3. An _Event_ has zero or more _Alerts_

### Monitors

A _Monitor_ is the top level building block of the configuration. It defines a service that runs in a loop, sending periodical HTTP requests to an URL and firing _Events_. You can define as many _Monitors_ as you wish. A _Monitor_ has the following configuration options:

- `name` (mandatory): an arbitrary string that is unique, e.g no other monitors can have the same name.
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



