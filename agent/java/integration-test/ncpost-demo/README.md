## 1.启动配置代理
配置代理包，路径需要改成自己的目录路径
-javaagent:E:\workspace\IdeaProjects\openrasp\agent\java\integration-test\rasp\rasp.jar

## 2.配置先编译后启动项目
clean package -Dmaven.test.skip=true

## 3.git子模块 openrasp-v8
1、删除子模块
（1）rm -rf 子模块目录 删除子模块目录及源码
（2）vi .gitmodules 删除项目目录下.gitmodules文件中子模块相关条目
（3）vi .git/config 删除配置项中子模块相关条目
（4）rm .git/module/* 删除模块下的子模块目录，每个子模块对应一个目录，注意只删除对应的子模块目录即可
(5) git rm --cached 子模块名称
（6）commit 这些修改–这一步很关键，否则会报错：already exists in the index

2、重新添加子模块
git submodule add git1@47.100.225.14:lvxi/openrasp-v8.git openrasp-v8

