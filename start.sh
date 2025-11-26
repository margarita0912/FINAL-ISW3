#!/bin/sh

# El backend ahora sirve también el frontend (./dist) y escucha en $PORT (Render)
exec ./backend-app
