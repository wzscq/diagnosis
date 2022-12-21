install docker
yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
systemctl start docker

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

docker run -it --name mosquitto -p 1983:1883 -p 9101:9001 -v /root/mosquitto/config:/mosquitto/config -v /root/mosquitto/data:/mosquitto/data -v /root/mosquitto/log:/mosquitto/log -d eclipse-mosquitto

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
export GO111MODULE=on
export GOPROXY=https://goproxy.io
让配置生效
source /etc/profile

//go get 加速
# 配置 GOPROXY 环境变量，以下三选一
# 1. 七牛 CDN
GOPROXY=https://goproxy.cn,direct
# 2. 阿里云
GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
# 3. 官方
GOPROXY=https://goproxy.io,direct

docker run -d -p80:80 --name diagnosis -v /root/crvframe/appfile:/services/crvframe/appfile -v /root/crvframe/apps:/services/crvframe/apps -v /root/crvframe/conf:/services/crvframe/conf -v /root/diagnosis/conf:/services/diagnosis/conf digimatrix/diagnosis:0.1.0

待处理问题列表
1、通过mqtt接收设备心跳消息和下发参数接收成功的消息  ok
2、通过流程化配置实现勾选参数下发+选择下发车辆+参数下发  ok
3、车辆下发记录展示  ok
4、auto2.0登录协议实现
5、错误处理的通用逻辑优化  ok
6、Excel图片导出  ok
7、故障内容换行显示  ok

2022-12-17 同一个时间点的诊断数据合并到一个记录上
1、数据库增加一个视图对数据进行合并 diag_view_result   ok
2、配置和视图对应的模型    ok
3、增加新的页面对应的菜单  ok
4、增加新的报告页面，支持多个测试记录合并展示 ok
5、修改查询报告数据的后台逻辑 ok
6、修改下载报告逻辑，支持多个测试记录合并展示 
7、增加修改报告的后台逻辑，可通过数据流配置