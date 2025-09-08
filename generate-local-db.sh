docker run -d --name jello-mark-database \
-p 127.0.0.1:3306:3306 \
--cpus="1" \
--memory="2g" \
-e MYSQL_USER=mad \
-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
-e MYSQL_DATABASE=main \
-v jello-mark-volume:/var/lib/mysql \
mysql:8.0