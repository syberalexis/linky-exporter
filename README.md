# Linky-exporter

[![Build Status](https://travis-ci.com/syberalexis/linky-exporter.svg?branch=master)](https://travis-ci.com/syberalexis/linky-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/linky-exporter)](https://goreportcard.com/report/github.com/syberalexis/linky-exporter)

This exporter get and expose French remote electrical information (Linky from EDF)


## Install

### From binary

Download binary from [releases page](https://github.com/syberalexis/linky-exporter/releases)

Example :
```bash
curl -L https://github.com/syberalexis/linky-exporter/releases/download/v1.0.1/linky-exporter-1.0.1-linux-amd64 -o /usr/local/bin/linky-exporter
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
  -p, --port=9901              Listen port
```


## Metrics example

```
# HELP linky_hchc Index heure creuse en Wh
# TYPE linky_hchc counter
linky_hchc 1.174768e+06
# HELP linky_hchp Index heure pleine en Wh
# TYPE linky_hchp counter
linky_hchp 3.523819e+06
# HELP linky_hhphc Horaire Heures Pleines Heures Creuses
# TYPE linky_hhphc untyped
linky_hhphc{name="A"} 1
# HELP linky_iinst Intensité instantanée en A
# TYPE linky_iinst counter
linky_iinst 10
# HELP linky_imax Intensité max
# TYPE linky_imax counter
linky_imax 90
# HELP linky_isousc Intensité souscrite en A
# TYPE linky_isousc counter
linky_isousc 30
# HELP linky_motdetat Mot d'état du compteur
# TYPE linky_motdetat untyped
linky_motdetat{name="000000"} 0
# HELP linky_optarif Option tarifaire
# TYPE linky_optarif untyped
linky_optarif{contrat="HC.."} 1
# HELP linky_papp Puissance Apparente, en VA
# TYPE linky_papp counter
linky_papp 2390
# HELP linky_ptec Période tarifaire en cours
# TYPE linky_ptec gauge
linky_ptec{option="hc"} 0
linky_ptec{option="hp"} 1
```
