#!/bin/bash

./debmod -dist="mmrepo" -dbhost="localhost" -dbport="5432" -dbpwd="Hendk387yrhdjue!" -source ~/temp -build ~/build -buildmode="meta" -pkg $1
