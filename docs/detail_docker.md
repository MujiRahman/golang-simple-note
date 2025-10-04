docker run -d \
 --name mariadb-container \
 -e MYSQL_ROOT_PASSWORD=password123 \
 -e MYSQL_DATABASE=mySimpleNote \
 -e MYSQL_USER=myUser \
 -e MYSQL_PASSWORD=password123 \
 -p 3306:3306 \
 -v mariadb_data:/var/lib/mysql \
 mariadb:latest

docker logs mariadb-container

docker exec -it mariadb-container mysql -u root -p

# Masukkan password: password123

# Test dari dalam container

docker-compose -f config/docker-compose.yml exec mariadb mysql -u myUser -p mySimpleNote

# Test dari aplikasi (contoh connection string Go)

# DB_DSN = "myUser:password123@tcp(mariadb:3306)/mySimpleNote?charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta"

cd config
docker-compose --env-file ../.env up -d

docker-compose -f config/docker-compose.yml up -d

docker-compose -f config/docker-compose.yml up -d mariadb
