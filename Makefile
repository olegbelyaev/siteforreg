SHELL := /bin/bash

mariadb-run:
	[[ ! -f "mariadb_secret.txt" ]] && echo "mariadb_secret.txt not found! Createit and try again." && exit; \
	[[ ! -d `pwd`/database/mariadb ]] && echo "NEW `pwd`/database/mariadb will be created!"; \
	mkdir -p `pwd`/database/mariadb || echo "Cant create  `pwd`/database/mariadb! Create it and try again."; \
	docker  run --name site-forreg-mariadb --hostname site-forreg-mariadb -p 3306:3306 \
	-v `pwd`/database/mariadb/init:/docker-entrypoint-initdb.d \
	-v `pwd`/database/mariadb:/var/lib/mysql \
	-it -e MYSQL_ROOT_PASSWORD=`cat mariadb_secret.txt` -d \
	mariadb --character-set-server=utf8mb4 \
	--collation-server=utf8mb4_unicode_ci ;\
	# от юзера: запуск mariadb


mariadb-start:
	docker start site-forreg-mariadb \
	# запуск остановленного контейнера


mariadb-stop:
	make mariadb-dump; \
	docker stop site-forreg-mariadb ;\
	# от юзера остановка mariadb


mariadb-exec-shell:
	docker exec -it site-forreg-mariadb mysql -p siteforeg --default-character-set=utf8 ;\
	# от юзера: подключение к mariadb для выполнения команд

mariadb-exec-sql-input:
	[[ ! -f "mariadb_secret.txt" ]] && echo "mariadb_secret.txt not found! Createit and try again." && exit; \
	docker exec -i site-forreg-mariadb mysql -p`cat mariadb_secret.txt` --default-character-set=utf8 ;\
	# sample: cat ./migrations/mariadb/init-database.sql | make mariadb-exec-sql-input

mariadb-dump:
	mkdir -p backup; \
	docker exec site-forreg-mariadb sh -c 'exec mysqldump --all-databases -uroot -p"$$MYSQL_ROOT_PASSWORD"' | gzip > ./backup/all-db.sql.gz


mariadb-restore:
	zcat ./backup/all-db.sql.gz | docker exec -i site-forreg-mariadb sh -c 'mysql -uroot -p"$$MYSQL_ROOT_PASSWORD"'


arango-run:
	[[ ! -f "arangodb_secret.txt" ]] && echo "arangodb_secret.txt not found! Createit and try again." && exit; \
	[[ ! -d `pwd`/database/arangodb ]] && echo " New `pwd`/database/arangodb will be created!"; \
	mkdir -p `pwd`/database/arangodb || echo "Cant create `pwd`/database/arangodb"; \
	docker run -e ARANGO_ROOT_PASSWORD=`cat arangodb_secret.txt` -d \
	-v `pwd`/database/arangodb:/var/lib/arangodb3 \
	-v `pwd`/database/arangodb/dump:/dump \
	-v `pwd`/database/arangodb/import:/import \
	-v `pwd`/database/arangodb/export:/export \
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


arango-import:
	docker exec -it site-forreg-arango arangoimport --overwrite; \
	echo "see `pwd`/arangodb_data/import/"



siteforreg-fork-run:
	./fork.pl -pf=siteforreg.pid --single "make siteforreg-run" >> siteforrreg.log  2>&1 ;\
	# запуск сервера siteforreg через форк


siteforreg-fork-kill:
	./fork.pl -pf=siteforreg.pid --kila >> siteforreg.log  2>&1 ;\
	# от рута: остановка процесса сервера siteforreg через форк


git-pull:
	git pull ;\

siteforreg-fork-release:
	make siteforreg-fork-kill && make git-pull && make siteforreg-fork-run ;\


siteforreg-run:
	source ./shell-scripts/export-siteforreg-vars.sh; \
	export PORT=8081; \
	go run tester.go; \
	# запуск сервера siteforreg от юзера на юзерском порту в консоли с выводом на косоль


webhook-server-run:
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


webhook-server-fork-run:
	./fork.pl -pf=webhook-server.pid --single "make webhook-server-run" >> webhook-server.log  2>&1 ;\


webhook-server-fork-kill:
	./fork.pl -pf=webhook-server.pid --kila


echo-ok:
	echo "test-ok"; \
	тестовый рецепт



