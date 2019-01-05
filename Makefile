docker-run-mysql:
	docker  run --name mysql --hostname mysql -p 3306:3306 \
	-v `pwd`:/docker-entrypoint-initdb.d \
	--rm -it -e MYSQL_ROOT_PASSWORD=11 -d \
	mysql --character-set-server=utf8mb4 \
	--collation-server=utf8mb4_unicode_ci

docker-stop-mysql:
	docker stop mysql

docker-exec-mysql:
	docker exec -it mysql mysql -p siteforeg --default-character-set=utf8 

go-run:
	SESSION_SECRET=HGJHGJHGJHGJHGJHGJHGIUYOIUY go run tester.go