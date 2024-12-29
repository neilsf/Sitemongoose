# Sitemongoose

Sitemongoose is a simple, lightweight and zero-dependency site monitoring tool written in Go. It is useful for

- Monitoring Website Availability: check if your websites are up and running by sending periodic HTTP requests. You can define expected Status Codes and send alerts if the server's response differs.
- Tracking Response Times: measure the response time of your websites to ensure they are performing optimally, or send alerts otherwise.
- Evaluating JSON Responses: define rules to evaluate JSON responses from your APIs and trigger alerts based on conditions such as value comparisons or key existence.
- Alerting: Sitemongoose can send alerts based on specific conditions. The following alerting channels are currently available:
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
            condition: "gt"
            value: 500
        alerts:
          - type: custom_cmd
            command: "/path/to/alert_script.sh"
```

## Run

To start monitoring, invoke Sitemongoose's start command and specify the location of the configuration file:

    ./sitemongoose start -c /path/to/config.yaml

In a production environment, you may want to run it as a service, using systemd or Supervisor.

## Full Config Specification

TBD

