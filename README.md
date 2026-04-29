#Для запуска проекта
student@yu01edu-78:~/social-network$ docker ps
student@yu01edu-78:~/social-network$ mkdir -p ~/.docker/cli-plugins/
student@yu01edu-78:~/social-network$ curl -SL https://github.com/docker/compose/releases/download/v2.26.1/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
student@yu01edu-78:~/social-network$ chmod +x ~/.docker/cli-plugins/docker-compose
student@yu01edu-78:~/social-network$ docker compose up --build -V запуск
docker compose down для сброса

docker-compose up --build

docker-compose down