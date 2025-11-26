#!/bin/sh

echo "Iniciando backend..."
./backend-app &

echo "Iniciando frontend..."
serve -s dist -l 3000
