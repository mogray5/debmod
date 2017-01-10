#!/bin/bash

./debmod -dist="mmrepo" -dbhost="localhost" -dbport="5432" -dbpwd="Hendk387yrhdjue!" -source ~/temp -build ~/build -pkg $1

#./debmod -dist="mmrepo" -pkg="$2" -source ~/temp
