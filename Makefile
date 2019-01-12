SHELL := /bin/bash


docker-run-mysql:
	docker  run --name site-forreg-mysql --hostname site-forreg-mysql -p 3306:3306 \
	-v `pwd`:/docker-entrypoint-initdb.d \
	--rm -it -e MYSQL_ROOT_PASSWORD=11 -d \
	mysql --character-set-server=utf8mb4 \
	--collation-server=utf8mb4_unicode_ci

docker-stop-mysql:
	docker stop site-forreg-mysql

docker-exec-mysql:
	docker exec -it site-forreg-mysql mysql -p siteforeg --default-character-set=utf8 


sudo-go-run:
	@make _export-default-vars; \
	export PATH="$HOME/bin:$HOME/.local/bin:$PATH:/home/dima/bin/go/bin"; \
	export GOPATH=/home/dima/go; \
	export PORT=80; \
	go run tester.go

go-run:
	source ./export-default-vars.sh; \
	export PORT=8081; \
	go run tester.go

_export-default-vars:
	if [ -r session_secret.txt ]; then \
	    SESSION_SECRET=`cat session_secret.txt`; \
	else \
	    echo "WARNING: NEW SESSION KEY GENERATED"; \
	    SESSION_SECRET=`date`; \
	fi; \
	export SESSION_SECRET; \
	as=admin_secret.txt; \
	if [ -r admin_secret.txt ]; then \
	    ADMIN_SECRET=`cat admin_secret.txt`; \
	else \
	    ADMIN_SECRET=11; \
	    echo "WARNING: file admin_secret.txt NOT FOUND. Use: 11"; \
	fi; \
	export ADMIN_SECRET; \
	if [ -r email_secret.txt ]; then \
		EMAIL_SECRET=`cat email_secret.txt`; \
	else \
		echo "WARNING: file email_secret.txt NOT FOUND"; \
	fi; \


	

