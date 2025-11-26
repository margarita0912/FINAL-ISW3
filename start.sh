#!/bin/sh

# Iniciar backend en background
./backend-app &

# Iniciar frontend
serve -s dist -l 3000
