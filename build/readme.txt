install mongo
mkdir /root/mongo
mkdir /root/mongo/data
mkdir /root/mongo/dump
mkdir /root/mongo/conf
docker run --name mongo -v /root/mongo/data:/data/db -v /root/mongo/dump:/dump -v /root/mongo/conf:/etc/mongo -p 37017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=AAA@111 -d mongo:5.0

install mysql
mkdir /root/mysql
mkdir /root/mysql/conf
mkdir /root/mysql/data
mkdir /root/mysql/log	

上传mysql配置文件mysql.cnf到服务器目录/root/mysql/conf下
docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -v /root/mysql/data:/var/lib/mysql -v /root/mysql/log:/var/log/mysql -p 4306:3306 -v /root/mysql/conf:/etc/mysql/conf.d -d  mysql:8.0.18

install redis
mkdir /root/redis
mkdir /root/redis/data
mkdir /root/redis/conf
touch /root/redis/conf/redis.conf

docker run --name redis -p 6479:6379 -v /root/redis/data:/data -v /root/redis/conf/redis.conf:/etc/redis/redis.conf --privileged=true --restart=always -d redis

install mosquitto
mkdir /root/mosquitto
mkdir /root/mosquitto/config
mkdir /root/mosquitto/data
mkdir /root/mosquitto/log
上传mosquitto.conf和password_file到/root/mosquitto/config目录下

docker run --name mosquitto -p 1983:1883 -p 9101:9001 -v /root/mosquitto/config:/mosquitto/config -v /root/mosquitto/data:/mosquitto/data -v /root/mosquitto/log:/mosquitto/log -d eclipse-mosquitto

install node
wget https://nodejs.org/dist/v16.15.1/node-v16.15.1-linux-x64.tar.xz
tar -xzvf node-v16.15.1-linux-x64.tar.xz
mv node-v16.15.1-linux-x64 node

vi /etc/profile  增加以下内容
export NODE_HOME=/root/node
export PATH=$NODE_HOME/bin:$PATH
让配置生效
source /etc/profile

install go
wget https://golang.google.cn/dl/go1.18.3.linux-amd64.tar.gz
tar -xzf go1.18.3.linux-amd64.tar.gz

vi /etc/profile  增加以下内容
export PATH=$PATH:/root/go/bin
让配置生效
source /etc/profile

docker run -p80:80 -v /root/crvframe/appfile:/frame_service/appfile -v /root/crvframe/apps:/frame_service/apps -v /root/crvframe/conf:/frame_service/conf -v /root/diagnosis/conf:/diagnosis/conf digimatrix/diagnosis:0.1.0