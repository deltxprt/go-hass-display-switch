# go-hass-display-switch

## Description

This is a small piece of program that watch an MQTT topic and paste the value of the message into ddccontrol for linux (Windows is TBD)

## Requirements

You need to have [`ddccontrol`]("https://github.com/ddccontrol/ddccontrol?tab=readme-ov-file") installed

## Configuration

the program only needs 2 parameters:

1. url:
    Format: mqtt://username:password@homeassistant.example.com:1883/topic
2. Device:
    Format: /dev/i2c-2
    Hint: using `ddccontrol -p` will show the i2c path of the display you wish to control

## Adhoc usage

A simple one liner: `./go-hass-display-switch -url mqtt://username:password@homeassistant.example.com:1883/topic -device "/dev/i2c-2"`

## As a Service

### Systemd

File: `/default/go-hass-ds`
```
URL=mqtt://username:password@homeassistant.example.com:1883/topic
DEVICE=/dev/i2c-2
```

File: `/etc/systemd/system/go-hass-display-switch.service`
```
[Unit]
Description=Display Switcher service made in GO
After=network.target

[Service]
EnvironmentFile=/default/go-hass-ds
ExecStart=/usr/bin/go-hass-display-switch -url ${URL} -device ${DEVICE}
Type=notify
Restart=always


[Install]
WantedBy=default.target
RequiredBy=network.target
```

## Installation

`curl "https://raw.githubusercontent.com/deltxprt/go-hass-display-switch/v1.0.2/install.sh" | sudo bash`
