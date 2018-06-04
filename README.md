# debmod

Related Minetest forum thread:  https://forum.minetest.net/viewtopic.php?f=14&t=13051

Tool to maintain an APT repository of Minetest mods and games

This tool packages Minetest mods and games pulling from Git based repositories that are maintained by plugin authors.  

Git repositories of upstream mod authors are scanned for changes and any changes detected will trigger a new package build. 

Tested on: Debian Jessie 8.5

**Requirements:**

Install postgresql 9.3 using postgresql apt repo: https://wiki.postgresql.org/wiki/Apt

apt-get install git build-essential reprepro equivs

Also requires:

Bash / GNU Utils

Mybatis Migrations was used to manage the database creation scripts:  http://www.mybatis.org/migrations

**debmod Usage**

The included bash scripts dobuild.sh and dobuildmeta.sh are required and need to be located in the folder specified in the -build argument when calling debmod. 
```
Usage of ./debmod:
  -build="/tmp/bbb": base folder to use for builds
  -buildmode="mods": Build mode can be mods, games, or meta
  -compat="8": Package compatibility
  -dbhost="localhost": Host location of PostgresSQL database.
  -dbname="mmrepodb": Name of mmrepo databaes.
  -dbport="5432": Port used by PostgresSQL database.
  -dbpwd="xxxx": Password to connect to mmrepo database.
  -dbuser="myuser": User to connect to mmrepo database.
  -debhelper=">= 8.0.0": Version of Debhelper to put into control file
  -deploy="/usr/share/games/minetest/mods": folder to install mod
  -dist="wheezy": Debian distribution
  -maintainer="myname <myemail>": Maintainer name and email
  -pkg="mmod-zzz": desired package name
  -repo="/var/opt/mmrepo": Path to APT repository
  -source="/tmp/zzz": base folder containing source
```
If you pass *mmod-zzz* to the -pkg argument then debmod will package all mods defined in the database.  When packaging mods you need to specify *mods* for the -buildmode argument which is the default. 

If you pass *mgame-zzz* to the -pkg argument then debmod will package all subgames defined in the database.  When packaging games you need to specify *games* for the -buildmode argument.

Scripts domod.sh and dogame.sh show examples of calling debmod.

A collection of SQL scripts also are included that demonstrate how mods and added to the database.
