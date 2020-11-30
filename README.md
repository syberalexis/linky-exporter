# Linky-exporter

[![Build Status](https://travis-ci.com/syberalexis/linky-exporter.svg?branch=master)](https://travis-ci.com/syberalexis/linky-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/linky-exporter)](https://goreportcard.com/report/github.com/syberalexis/linky-exporter)

This exporter get and expose French remote electrical information (Linky from EDF)


## Install

### From binary

Download binary from [releases page](https://github.com/syberalexis/linky-exporter/releases)

Example :
```bash
curl -L https://github.com/syberalexis/linky-exporter/releases/download/v2.0.0/linky-exporter-2.0.0-linux-amd64 -o /usr/local/bin/linky-exporter
chmod +x /usr/local/bin/linky-exporter
/usr/local/bin/linky-exporter
```

### From sources

```bash
git clone git@github.com:syberalexis/linky-exporter.git
cd linky-exporter
go build cmd/linky-exporter/main.go -o linky-exporter
./linky-exporter
```

or

```bash
git clone git@github.com:syberalexis/linky-exporter.git
cd linky-exporter
GOOS=linux GOARCH=amd64 VERSION=0.1.3 make clean build
./dist/linky-exporter-0.1.3-linux-amd64
```

## Install as a service

```
[Unit]
Description=Linky Exporter service
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/linky-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable linky-exporter
systemctl start linky-exporter
```

## Help

```
usage: main [<flags>]

Flags:
  -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
      --version                Show application version.
  -a, --address="0.0.0.0"      Listen address
  -b, --baud=1200              Baud rate
  -d, --device="/dev/serial0"  Device to read
      --parity="ParityNone"    Serial parity
  -p, --port=9901              Listen port
      --size=7                 Serial frame size
      --stopbits="Stop1"       Serial stopbits
```


## Metrics example

```
# HELP linky_hours_group_info Groupe horaire (tarif Tempo ou HPHC)
# TYPE linky_hours_group_info gauge
linky_hours_group_info{groupe="A",idcompteur="032914312110",tarif="heures creuses"} 1
# HELP linky_index_watthours_total Index en Wh
# TYPE linky_index_watthours_total counter
linky_index_watthours_total{idcompteur="032914312110",periode="HC",tarif="heures creuses"} 2.6907285e+07
# HELP linky_intensity_amperes Intensité en A
# TYPE linky_intensity_amperes gauge
linky_intensity_amperes{idcompteur="032914312110",tarif="heures creuses"} 3
# HELP linky_maximum_intensity_amperes Intensité maximale en A
# TYPE linky_maximum_intensity_amperes gauge
linky_maximum_intensity_amperes{idcompteur="032914312110",tarif="heures creuses"} 60
# HELP linky_power_voltamperes Puissance apparente en VA
# TYPE linky_power_voltamperes gauge
linky_power_voltamperes{idcompteur="032914312110",tarif="heures creuses"} 760
# HELP linky_subscribed_intensity_amperes Intensité souscrite en A
# TYPE linky_subscribed_intensity_amperes gauge
linky_subscribed_intensity_amperes{idcompteur="032914312110",tarif="heures creuses"} 30
```
