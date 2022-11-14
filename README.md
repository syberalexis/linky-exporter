# Linky-exporter

[![Build Status](https://travis-ci.com/syberalexis/linky-exporter.svg?branch=master)](https://travis-ci.com/syberalexis/linky-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/linky-exporter)](https://goreportcard.com/report/github.com/syberalexis/linky-exporter)

This exporter get and expose French remote electrical information (Linky from EDF).

## Summary

- [Install](#install)
  - [From binary](#from-binary)
  - [From docker](#from-docker)
  - [From sources](#from-sources)
- [Install as a service](#install-as-a-service)
  - [Systemd](#systemd)
  - [OpenBSD](#openbsd)
- [Help](#help)
- [Metrics example](#metrics-example)
  - [Historical](#historical)
  - [Standard](#standard)
- [How to make all installation on Raspberry Pi Zero](#how-to-make-all-installation-on-raspberry-pi-zero)
  - [Without USB on Linky](#without-usb-on-linky)
  - [With USB on Linky](#with-usb-on-linky)
- [Good links](#good-links)

## Install

### From binary

Download binary from [releases page](https://github.com/syberalexis/linky-exporter/releases)

Example :
```bash
curl -L https://github.com/syberalexis/linky-exporter/releases/download/v3.0.0/linky-exporter-3.0.0-linux-amd64 -o /usr/local/bin/linky-exporter
chmod +x /usr/local/bin/linky-exporter
/usr/local/bin/linky-exporter
```

### From docker

```bash
docker pull syberalexis/linky-exporter
docker run -d -p 9901:9901 -v /dev/serial0:/dev/serial0 syberalexis/linky-exporter:3.0.0 --device /dev/serial0
```

### From sources

```bash
git clone git@github.com:syberalexis/linky-exporter.git
cd linky-exporter
go build cmd/linky-exporter/main.go -o linky-exporter
./linky-exporter --device /dev/serial0
```

or

```bash
git clone git@github.com:syberalexis/linky-exporter.git
cd linky-exporter
GOOS=linux GOARCH=amd64 VERSION=3.0.0 make clean build
./dist/linky-exporter-3.0.0-linux-amd64 --device /dev/serial0
```

## Install as a service

In file `/lib/systemd/system/linky_exporter.service` :
### Systemd
```
[Unit]
Description=Linky Exporter service
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/linky-exporter --device /dev/serial0
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable linky-exporter
systemctl start linky-exporter
```

### OpenBSD
In file `/etc/rc.d/linky_exporter` :
```
#!/bin/ksh

daemon="/usr/local/bin/linky_exporter --device /dev/cuaU0" 
daemon_logger="daemon.info"
daemon_user="_nodeexporter"
. /etc/rc.d/rc.subr

pexp="${daemon}.*"
rc_bg=YES
rc_reload=NO

rc_cmd $1
```

```shell
usermod -G dialer _nodeexporter
chmod 555 /etc/rc.d/linky_exporter
rcctl enable linky_exporter
```

## Help

```
usage: linky-exporter --device=DEVICE [<flags>]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --version            Show application version.
      --debug              Enable debug mode.
  -a, --address="0.0.0.0"  Listen address
  -p, --port=9901          Listen port
      --auto               Automatique mode
      --historical         Historical mode
      --standard           Standard mode
  -d, --device=DEVICE      Device to read
  -b, --baud=BAUD          Baud rate
      --size=SIZE          Serial frame size
      --parity=PARITY      Serial parity
      --stopbits=STOPBITS  Serial stopbits
```

## Metrics example

### Historical
```
# HELP linky_energy Energie en Wh
# TYPE linky_energy counter
linky_energy{index="F1",linky_id="XXXX",mode="used"} 2.345675e+06
linky_energy{index="F2",linky_id="XXXX",mode="used"} 6.662251e+06
# HELP linky_energy_total Total Energie en Wh
# TYPE linky_energy_total counter
linky_energy_total{linky_id="XXXX",mode="used"} 9.007926e+06
# HELP linky_intensity Courant efficace en A
# TYPE linky_intensity gauge
linky_intensity{linky_id="XXXX",phase="1"} 11
# HELP linky_power Puissance apparente en VA
# TYPE linky_power gauge
linky_power{linky_id="XXXX",mode="used",phase="1"} 2530
# HELP linky_power_reference Puissance apparente de référence en kVA
# TYPE linky_power_reference counter
linky_power_reference{linky_id="XXXX",type="subscribed"} 6
# HELP linky_timestamp Synchronized timestamp in Linky
# TYPE linky_timestamp counter
linky_timestamp{contract="HC..",linky_id="XXXX",pricing="HP..",version="1"} 0
# HELP linky_voltage Tension efficace en V
# TYPE linky_voltage gauge
```

### Standard
```
# HELP linky_energy Energie en Wh
# TYPE linky_energy counter
linky_energy{index="D1",linky_id="XXXX",mode="used"} 4.1585532e+07
linky_energy{index="F1",linky_id="XXXX",mode="used"} 4.1352473e+07
linky_energy{index="F2",linky_id="XXXX",mode="used"} 233059
# HELP linky_energy_total Total Energie en Wh
# TYPE linky_energy_total counter
linky_energy_total{linky_id="XXXX",mode="used"} 4.1585532e+07
# HELP linky_intensity Courant efficace en A
# TYPE linky_intensity gauge
linky_intensity{linky_id="XXXX",phase="1"} 6
# HELP linky_movable_peak Pointe mobile
# TYPE linky_movable_peak gauge
linky_movable_peak{linky_id="XXXX",phase="1",type="end"} 0
linky_movable_peak{linky_id="XXXX",phase="1",type="start"} 0
linky_movable_peak{linky_id="XXXX",phase="2",type="end"} 0
linky_movable_peak{linky_id="XXXX",phase="2",type="start"} 0
linky_movable_peak{linky_id="XXXX",phase="3",type="end"} 0
linky_movable_peak{linky_id="XXXX",phase="3",type="start"} 0
# HELP linky_power Puissance apparente en VA
# TYPE linky_power gauge
linky_power{linky_id="XXXX",mode="used",phase="1"} 1420
# HELP linky_power_last_year Puissance apparente n-1 en VA
# TYPE linky_power_last_year gauge
linky_power_last_year{linky_id="XXXX",mode="used",phase="1"} 3080
# HELP linky_power_max Puissance apparente en VA
# TYPE linky_power_max gauge
linky_power_max{linky_id="XXXX",mode="used",phase="1"} 2860
# HELP linky_power_reference Puissance apparente de référence en kVA
# TYPE linky_power_reference counter
linky_power_reference{linky_id="XXXX",type="breaking"} 6
linky_power_reference{linky_id="XXXX",type="subscribed"} 6
# HELP linky_status Statut issu du registre
# TYPE linky_status gauge
linky_status{linky_id="XXXX",name="Contact sec"} 0
linky_status{linky_id="XXXX",name="Couleur du jour pour le contrat historique tempo"} 0
linky_status{linky_id="XXXX",name="Couleur du lendemain pour le contrat historique tempo"} 0
linky_status{linky_id="XXXX",name="Dépassement de la puissance de référence"} 0
linky_status{linky_id="XXXX",name="Fonctionnement producteur/consommateur"} 0
linky_status{linky_id="XXXX",name="Mode dégradée de l horloge"} 0
linky_status{linky_id="XXXX",name="Organe de coupure"} 0
linky_status{linky_id="XXXX",name="Pointe mobile (PM)"} 0
linky_status{linky_id="XXXX",name="Préavis pointes mobiles"} 0
linky_status{linky_id="XXXX",name="Sens de l énergie active"} 0
linky_status{linky_id="XXXX",name="Statut du CPL"} 0
linky_status{linky_id="XXXX",name="Surtension sur une des phases"} 0
linky_status{linky_id="XXXX",name="Synchronisation CPL"} 0
linky_status{linky_id="XXXX",name="Tarif en cours sur le contrat distributeur"} 0
linky_status{linky_id="XXXX",name="Tarif en cours sur le contrat fourniture"} 0
linky_status{linky_id="XXXX",name="État de la sortie communication Euridis"} 0
linky_status{linky_id="XXXX",name="État de la sortie télé-information"} 0
linky_status{linky_id="XXXX",name="État du cache-bornes distributeur"} 0
# HELP linky_timestamp Timestamp en seconde
# TYPE linky_timestamp counter
linky_timestamp{contract="BASE",linky_id="XXXX",pricing="BASE",version="02"} 1668350147
# HELP linky_voltage Tension efficace en V
# TYPE linky_voltage gauge
linky_voltage{linky_id="XXXX",phase="1"} 229
# HELP linky_voltage_average Tension moyenne en V
# TYPE linky_voltage_average gauge
linky_voltage_average{linky_id="XXXX",phase="1"} 230
```

## How to make all installation on Raspberry Pi Zero

### Without USB on Linky

It's my case.  

I followed this french tutorial https://www.jonathandupre.fr/articles/24-logiciel-scripts/208-suivi-consommation-electrique-compteur-edf-linky-avec-raspberry-pi-zero-w/  
To use this exporter, you don't need install the PHP script, [follow this](#install).

And I bought this convertor https://www.tindie.com/products/Hallard/pitinfo/

### With USB on Linky

If you have an USB on your Linky, follow this french tutorial https://sebastienreuiller.fr/blog/monitorer-son-compteur-linky-avec-grafana-cest-possible-et-ca-tourne-sur-un-raspberry-pi/

## Good links

- https://www.enedis.fr/media/2035/download
- https://www.capeb.fr/www/capeb/media/vaucluse/document/FicheSeQuelecN17TIC.pdf
