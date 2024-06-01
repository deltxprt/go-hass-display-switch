# go-hass-display-switch

## Description

This is a small piece of program that watch an MQTT topic and paste the value of the message into ddccontrol for linux (Windows is TBD)

## How does it work?

- The program will listen to the messages that you send on MQTT.

    *The value needs to be a number and nothing else, otherwise it won't work!*
  
- When the message is received it will execute the ddccontrol command to switch the display

    *it doesn't check the current value since the value is inacurate, check the #QnA section for details*
- repeat the steps!

## Requirements

You need to have [`ddccontrol`](https://github.com/ddccontrol/ddccontrol?tab=readme-ov-file) installed

## Installation

Download the latest package in the [release page](https://github.com/deltxprt/go-hass-display-switch/releases)

## Configuration

the program only needs 2 parameters:

1. url:
    - Format: `mqtt://username:password@homeassistant.example.com:1883/topic`
2. Device:
    - Format: /dev/i2c-2
    - Hint: using `ddccontrol -p` will show the i2c path of the display you wish to control

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
EnvironmentFile=/etc/default/go-hass-ds
ExecStart=/usr/bin/go-hass-display-switch -url ${URL} -device ${DEVICE}
Type=notify
Restart=always


[Install]
WantedBy=default.target
RequiredBy=network.target
```

# QnA

## How to know what's the right values for each of my inputs?

Standards for DDC/CI are not well defined and documented so it does require some trials and errors to find the right values.

A good example was for my Samsung Odyssey Neo G9. 
when you probe the displays with `ddccontrol -p` it show the possible value for each settings ( VGA-1 = 1, DVI = 3 and maximum was 3)
Yet those values are not the values my monitor were looking for.
In the end, my monitor inputs numbers where 15 for DP, 16 for HDMI 2 and 17 for HDMI 1.

So you need to test it yourself to figure each inputs number.

Good Luck!


# Home Assistant Setup

All you need to do is create scripts that will send the value associated with the input you want to display.

## example

``` yaml
alias: Display DP
sequence:
  - service: mqtt.publish
    metadata: {}
    data:
      qos: "0"
      retain: false
      topic: display_ctrl
      payload: "15"
mode: single
icon: mdi:monitor
```
