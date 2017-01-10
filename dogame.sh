#!/bin/bash

./debmod -dist="mmrepo" -dbhost="localhost" -dbport="5432" -dbpwd="mypassword" -source ~/temp -build ~/build -deploy="/usr/share/games/minetest/games" -buildmode="games" -pkg $1

#./debmod -dist="mmrepo" -pkg="$2" -source ~/temp
