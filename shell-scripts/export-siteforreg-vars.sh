#!/bin/bash

	if [ -r session_secret.txt ]; then
	    SESSION_SECRET=`cat session_secret.txt`;
	else
	    echo "WARNING: NEW SESSION KEY GENERATED";
	    SESSION_SECRET=`date`;
	fi;
	export SESSION_SECRET;


	as=admin_secret.txt;
	if [ -r admin_secret.txt ]; then
	    ADMIN_SECRET=`cat admin_secret.txt`;
	else
	    ADMIN_SECRET=11;
	    echo "WARNING: file admin_secret.txt NOT FOUND. Use: 11";
	fi;
	export ADMIN_SECRET;
	
    
    if [ -r email_secret.txt ]; then
		EMAIL_SECRET=`cat email_secret.txt`;
	else
		echo "WARNING: file email_secret.txt NOT FOUND";
	fi;
    export EMAIL_SECRET;
    
    
    if [ -r mariadb_secret.txt ]; then
		MYSQL_SECRET=`cat mariadb_secret.txt`;
	else
		echo "WARNING: file mariadb_secret.txt NOT FOUND";
	fi;
    export MYSQL_SECRET;
