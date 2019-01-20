SHELL := /bin/bash

mysql-run:
	docker  run --name site-forreg-mysql --hostname site-forreg-mysql -p 3306:3306 \
	-v `pwd`:/docker-entrypoint-initdb.d \
	-it -e MYSQL_ROOT_PASSWORD=11 -d \
	mysql --character-set-server=utf8mb4 \
	--collation-server=utf8mb4_unicode_ci ;\
	# от юзера: запуск mysql


mysql-start:
	docker start site-forreg-mysql \
	# запуск остановленного контейнера


mysql-stop:
	make mysql-dump; \
	docker stop site-forreg-mysql ;\
	# от юзера остановка mysql


mysql-exec-shell:
	docker exec -it site-forreg-mysql mysql -p siteforeg --default-character-set=utf8 ;\
	# от юзера: подключение к mysql для выполнения команд


mysql-dump:
	mkdir -p backup; \
	docker exec site-forreg-mysql sh -c 'exec mysqldump --all-databases -uroot -p"$$MYSQL_ROOT_PASSWORD"' | gzip > ./backup/all-db.sql.gz


mysql-restore:
	zcat ./backup/all-db.sql.gz | docker exec -i site-forreg-mysql sh -c 'mysql -uroot -p"$$MYSQL_ROOT_PASSWORD"'


arango-run:
	[[ ! -d `pwd`/arangodb_data ]] && echo "`pwd`/arangodb_data not exists! Create it and try again." && exit; \
	[[ ! -f "arangodb_secret.txt" ]] && echo "arangodb_secret.txt not found! Createit and try again." && exit; \
	docker run -e ARANGO_ROOT_PASSWORD=`cat arangodb_secret.txt` -d -v `pwd`/arangodb_data:/var/lib/arangodb3 \
	-v `pwd`/arangodb_data/dump:/dump -v `pwd`/arangodb_data/import:/import -v `pwd`/arangodb_data/export:/export \
	--name site-forreg-arango --hostname site-forreg-arango -p 8529:8529 -v ARANGO_STORAGE_ENGINE=rocksdb arangodb


arango-start:
	docker start site-forreg-arango


arango-stop:
	docker stop site-forreg-arango


arango-dump:
	docker exec -it site-forreg-arango arangodump --overwrite; \
	echo "see `pwd`/arangodb_data/dump/"


arango-export:
	docker exec -it site-forreg-arango arangoexport --overwrite; \
	echo "see `pwd`/arangodb_data/export/"


arango-export:
	docker exec -it site-forreg-arango arangoimport --overwrite; \
	echo "see `pwd`/arangodb_data/import/"


sudo-siteforreg-fork-run:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	./fork.pl -pf=siteforreg.pid --single "make sudo-siteforreg-run" >> siteforrreg.log  2>&1 ;\
	# от рута: запуск сервера siteforreg через форк


sudo-siteforreg-fork-kill:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	./fork.pl -pf=siteforreg.pid --kila >> siteforreg.log  2>&1 ;\
	# от рута: остановка процесса сервера siteforreg через форк


git-pull:
	sudo -u dima git pull ;\
	# перевод [рута] на пользователя и стягивание из гита


sudo-siteforreg-fork-release:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	make sudo-siteforreg-fork-kill && make git-pull && make sudo-siteforreg-fork-run ;\
	# от рута: остановка сервера через форк, от юзера стягивание из гита обновлений и от рута запуск через форк (не трогает mysql)


sudo-siteforreg-run:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	source ./shell-scripts/export-siteforreg-vars.sh; \
	export PATH="${HOME}/bin:${HOME}/.local/bin:${PATH}:/home/dima/bin/go/bin"; \
	export GOPATH=/home/dima/go; \
	export PORT=80; \
	go run tester.go ;\
	# запуск сервера siteforreg от рута на порту 80


siteforreg-run:
	source ./shell-scripts/export-siteforreg-vars.sh; \
	export PORT=8081; \
	go run tester.go; \
	# запуск сервера siteforreg от юзера на юзерском порту в консоли с выводом на косоль


sudo-webhook-server-run:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	if [ -r webhook_secret.txt ]; then \
		WEBHOOK_SECRET=`cat webhook_secret.txt`; \
	else \
		echo "WARNING: file webhook_secret.txt NOT FOUND"; \
	fi; \
	export WEBHOOK_SECRET; \
	export PATH="${HOME}/bin:$HOME/.local/bin:${PATH}:${HOME}/bin/go/bin"; \
	export GOPATH=/home/dima/go; \
	export WEBHOOK_SCRIPT="./shell-scripts/webhook.sh"; \
	go run webhook.go; \
	# запускает сервер, ожидающий запросов на отдельный порт, реагирующий на веб-хуки


sudo-webhook-server-fork-run:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	./fork.pl -pf=webhook-server.pid --single "make sudo-webhook-server-run" >> webhook-server.log  2>&1 ;\
	# от рута: запуск сервера webhook через форк


sudo-webhook-server-fork-kill:
	[[ ${USER} != "root" ]] && echo 'this command should be run with sudo' && exit; \
	./fork.pl -pf=webhook-server.pid --kila


echo-ok:
	echo "test-ok"; \
	тестовый рецепт

