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

sudo-go-run-fork:
	./fork.pl -pf=siteforreg.pid --single "make sudo-go-run" >> siteforrreg.log  2>&1

sudo-go-run:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	source ./export-siteforreg-vars.sh; \
	export PATH="${HOME}/bin:${HOME}/.local/bin:${PATH}:/home/dima/bin/go/bin"; \
	export GOPATH=/home/dima/go; \
	export PORT=80; \
	go run tester.go 

docker-run-go:
	

go-run:
	source ./export-siteforreg-vars.sh; \
	export PORT=8081; \
	go run tester.go


	

run-webhook-server:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	if [ -r webhook_secret.txt ]; then \
		WEBHOOK_SECRET=`cat webhook_secret.txt`; \
	else \
		echo "WARNING: file webhook_secret.txt NOT FOUND"; \
	fi; \
	export WEBHOOK_SECRET; \
	export PATH="${HOME}/bin:$HOME/.local/bin:${PATH}:${HOME}/bin/go/bin"; \
	export GOPATH=/home/dima/go; \
	go run webhook.go


echo-ok:
	echo "test-ok"

