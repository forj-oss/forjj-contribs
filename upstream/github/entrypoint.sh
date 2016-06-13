#!/bin/sh

# Create a real tty, then run github driver
exec >/dev/tty 2>/dev/tty </dev/tty && github "$@"
