### 技术选型
框架: gin

日志处理: zap

orm: gorm

优雅重启HTTP服务: gracehttp

测试框架: ginkgo


### 接口风格

 RESTFUL

###  项目结构
* api - handler函数
* route - 路由注册
* model - 数据模型以及数据库操作
* doc - 接口文档
* log - 日志处理
* config - 配置文件
* tools - 公共工具类
* vendor - 项目依赖其他开源项目目录
* dist - 打包生成安装包的存放路径
* main.go - 程序执行入口
* Makefile 提供编译、打包、测试等功能的脚本文件
* ginkgo 二进制文件，容器内执行测试用例的时候需要使用的命令，请勿移除
* junit.xml 测试报告
* coverprofile.txt 通过的测试的覆盖率的概要

### 语法检查

    make govet

### 检查是否符合官方统一标准的风格

    make gofmt

### 编译

    make build

### 打包

    make package

### 测试

#### 先创建一个新的Docker网络

    docker network create -d bridge my-net

#### 执行测试用例
创建mysql容器实例会占用3306，确保该端口未被其他应用使用
  
    make test

最后生成的测试报告junit.xml和覆盖率coverprofile.txt文件在对应package的目录下

make test包括下面三个步骤

*  移除容器,排除历史测试数据干扰

    make clean-docker

*  启动docker，创建一个mysql容器实例

    make start-docker

*  开始执行测试用例

   make test-server

### 测试用例覆盖率可视化

  go tool cover -html=coverprofile.txt -o coverprofile.html

  可以很清楚地看到测试用例覆盖的代码和未曾覆盖到的代码

### 管理依赖包工具(Vendor Tool)

	govendor是依赖管理利器 

##### install govendor

	go get -u github.com/kardianos/govendor


##### command

通过指定包类型，可以过滤仅对指定包进行操作。


| 命令	|  功能   |
|-------| --------|
| govendor init |	初始化 vendor 目录 |
| govendor list |	列出所有的依赖包 |
| govendor add |	添加包到 vendor 目录，如 govendor add +external 添加所有外部包 |
| govendor add PKG_PATH |	添加指定的依赖包到 vendor 目录|
| govendor update |	从 $GOPATH 更新依赖包到 vendor 目录|
| govendor remove |	从 vendor 管理中删除依赖  |
| govendor status |	列出所有缺失、过期和修改过的包  |
| govendor fetch |	添加或更新包到本地 vendor 目录 |
| govendor sync | 	本地存在 vendor.json 时候拉去依赖包，匹配所记录的版本 |
| govendor get |	类似 go get 目录，拉取依赖包到 vendor 目录 |


### Create Database
* 以root用户登录MySQL
  
  mysql -u root -p

* 创建luckgo用户'luckgo'
  
  mysql> create user 'luckgo'@'%' identified by 'luckgo-password';
   其中%表示网上的所有机器都可以连接上，使用具体的IP地址更安全点
  mysql> create user 'luckgo'@'0.0.0.0' identified by 'luckgo-password';


* 创建critic数据库

  mysql>  create database luckgo;


* 允许luckgo用户的访问权限

  mysql> grant all privileges on luckgo.* to 'luckgo'@'%';


* 退出MySQL

  mysql> exit
  
### 平滑升级

 kill -USR2  PID

### 应用日志切割
使用Linux系统默认安装的logrotate工具,在目录/etc/logrotate.d/下增加配置文件luckgo.conf

    /path/to/log1 /path/to/log2 {
        compress
        dateext
        maxage 365
        rotate 60
        size = 10M
        notifempty
        missingok
        create 644
        copytruncate
    }
