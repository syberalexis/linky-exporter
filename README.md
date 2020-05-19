# Linky-exporter

[![Build Status](https://travis-ci.com/syberalexis/linky-exporter.svg?branch=master)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/linky-exporter)](https://goreportcard.com/report/github.com/syberalexis/linky-exporter)

This exporter get and expose French remote electrical information (Linky from EDF)

## Example

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