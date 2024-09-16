#!/bin/bash

colima start --cpu 4 --memory 8 --disk 100 --vm-type qemu

docker run -itd --rm --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydb -p 5432:5432 postgres
