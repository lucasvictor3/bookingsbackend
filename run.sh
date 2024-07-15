#!/bin/bash

go build -o bookings cmd/web/*.go && 
./bookings -dbname=bookings-1 -dbuser=postgres -cache=false -production=false -dbpass=test -dbhost=172.25.176.1