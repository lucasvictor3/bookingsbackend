#!/bin/bash
echo $DBHOST && 
ls && go build -o bookings cmd/web/*.go && 
./bookings 