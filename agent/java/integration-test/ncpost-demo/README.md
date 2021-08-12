## 1.启动配置代理
配置代理包，路径需要改成自己的目录路径
-javaagent:E:\workspace\IdeaProjects\openrasp\agent\java\integration-test\rasp\rasp.jar

## 2.配置先编译后启动项目
clean package -Dmaven.test.skip=true

