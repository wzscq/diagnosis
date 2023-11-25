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

docker run -d -p80:80 --name diagnosis -v /home/Digimatrix/project/saic/TotalData/data:/services/diagnosis/data -v /root/crvframe/appfile:/services/crvframe/appfile -v /root/crvframe/apps:/services/crvframe/apps -v /root/crvframe/conf:/services/crvframe/conf -v /root/diagnosis/conf:/services/diagnosis/conf digimatrix/diagnosis:0.1.0

//alpine支持kafka的编译，需要在alpine 3.15的镜像中实现
1、执行以下命令安装alpine的编译环境
   apk add alpine-sdk
2、执行以下命令编译程序
   go build -tags musl


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

2023-01-31 增加原始数据下载功能
1、模型diag_result，增加字段raw_data,数据行上增加下载原始文件按钮
2、模型diag_event，增加字段raw_data,数据行上增加下载原始文件按钮
3、模型diag_view_result，增加字段raw_data,数据行上增加下载原始文件按钮

2023-02-11 完善综合分析功能
1、修改diag_v2report，合并相同ECU物流和故障信息
2、修改diag_result和diag_view_result模型配置，支持综合分析页面的关闭功能
3、增加综合导出接口，导出合并的分析报告

2023-03-05 DTC维护页面中ECU改为多选
1、中间表diag_dtc_diag_ecu
2、修改diag_dtc的配置，ecu字段改为many2many类型，编辑使用穿梭框，listview中对多个ECU用逗号拼接显示
3、listview修改，增加many2many字段的显示，控制many2many字段不允许过滤和排序操作
4、formview修改，修改many2many字段级联下拉选择逻辑
5、修改diag_manual_fault配置，dtc的级联过滤配置更新

2023-05-20 修改诊断参数下发逻辑
1、查询dtclist时原来没有按照platformid做过滤，这次增加了这个过滤条件

2023-10-28 修改页面配置
1、修改菜单配置中的菜单名称：menu
1、项目数据库：diag_platform，
   a、修改配置
2、总线数据库：diag_domain
3、控制器数据库：diag_ecu
4、信号数据库： diag_signal   遗留问题，主键基于目前的逻辑存在冲突风险
5、故障码数据库：diag_dtc
   a、创建数据库视图：view_diag_dtc，这里注意修改了关联表的排序规则设置
   b、增加view_diag_dtc对应的相关配置
6、关联信号数据库：diag_manual_fault
7、诊断服务数据库：diag_logistics
8、诊断参数配置：diag_parameter
   a、数据表 vehiclemanagement 增加字段developPhase
   b、kafka接收车辆设备信息修改逻辑，对developPhase字段做处理
   c、修改了流程sendparameter配置
9、诊断下发记录：diag_param_sendrecord
   a、数据库增加字段: platform_id、test_specification
   b、后台下发逻辑中补充对platform_id和test_specification字段的处理
   c、数据库增加视图：view_vehicle
   d、增加相应的配置：view_vehicle
   e、修改diag_select_cars配置，关联到view_vehicle进行车辆先择
10、触发回传配置：diag_event_parameter
11、触发回传下发记录：diag_event_sendrecord
   a、数据库增加字段: platform_id、test_specification
   b、后台下发逻辑中补充对platform_id和test_specification字段的处理
12、车辆信息：vehiclemanagement   
13、设备上报记录：diag_device_heartbeat,目前上报信息中是否携带项目信息、试验规范和试验阶段信息
      a、增加字段：vehicle_management_code、project_num、test_specification、develop_phase
      b、修改后端心跳处理逻辑
      c、修改对应页面配置
14、故障仪表盘：
      a、页面展示效果修改
      b、查询逻辑修改
15、智能诊断分析：diag_result，部分字段目前数据库中没有 SAE故障码，需要诊断程序填充部分信息。
      增加字段：sae_code
16、智能诊断综合分析：diag_view_result
17、触发回传信号分析：diag_event  数据中缺少车辆编号字段、
      a、数据库表中增加车辆编号字段:vehicle_management_code
18、用户管理：core_user
    修改数据表，增加字段：email、department、job_number、dimission、disable
19、角色管理：core_role
20、访问统计：
    数据库中增加表：core_operation_log
    配置报表模块:reports/access_statics


2023-10-30 
1、IDM同步用户信息
   更新后台服务和相应的配置文件

2023-11-06
1、修改dashboard页面代码
2、修改配置diag_result.json
3、修改配置diag_event
4、修改配置diag_select_cars
5、修改配置diag_param_sendrecord
6、修改配置diag_event_parameter
7、修改配置diag_event_sendrecord
8、修改配置diag_platform
9、修改配置diag_device_heartbeat

2023-11-8
1、将原来换色显示的故障改为蓝色
   a、修改diag_report模块，
   b、修改diag_v2report模块
   c、后端服务导出Excel逻辑修改
2、报告页面打开问题
   a、修改配置diag_result
   b、修改配置diag_view_result
3、页面显示错误
   a、修改配置文件diag_signal_sendrecord.json

2023-11-19
1、修改dashboard页面代码
   a、故障控制器分布字体不一致问题
   b、年份默认显示当前年份
   c、项目故障分布横坐标部分项目没显示且不能横向拖动
   d、故障列表没有分页
2、进入首页时不跳转大蓝页
   项目配置中增加oauthback.png
3. 故障仪表盘改为首页
   修改菜单配置
   修改全局操作配置
4、所有页面序号固定放在操作列后，从1开始编号，更新人显示姓名，删除多余标题
   修改所有页面对应配置文件
   a、隐藏序列号字段
   b、更新人显示对应账号姓名
   c、去掉默认视图的名称
5、表vehiclemanagement增加字段vehicleConfiger
   修改对应的同步数据逻辑
   修改对应页面配置

2023-11-25
1、问题20，增加SAE故障码防止输入小写功能
   需改模型配置：diag_dtc
2、问题21，车辆信息页面隐藏序列号列
   修改模型配置：vehiclemanagement
3、问题40，故障码库所属控制器字段设置为必填
   需改模型配置：diag_dtc
4、问题42，故障仪表盘修改
5、问题46，设备上报记录序列号隐藏
   需改模型配置：diag_device_heartbeat
   修改获取项目信息的逻辑，原来通过vin匹配不到，改为通过device_id匹配
6、问题47，故障仪表盘修改
7、问题50，修改配置文件
   模型配置：core_user
   模型配置：core_operation_log
   模型配置：core_role
8、问题52，车辆信息页面增加解绑和绑定时间
   修改模型配置：vehiclemanagement
9、问题53，故障仪表盘修改
10、问题55，诊断下发记录发送人应显示名字
    模型配置：diag_param_sendrecord
11、问题56，修改后台下发逻辑
12、问题57，隐藏诊断执行配置功能
    需改菜单配置：menus.json
13、问题58，触发回传下发记录发送人应显示名字
   模型配置:diag_event_sendrecord
14、问题59，修改后台下发逻辑
15、问题60，隐藏信号下发记录菜单
   需改菜单配置：menus.json
16、问题62，修改密码功能隐藏
   修改框架处理逻辑
   增加应用配置：app.json
17、问题63，登录token的有效期延长已经设置了6000s


