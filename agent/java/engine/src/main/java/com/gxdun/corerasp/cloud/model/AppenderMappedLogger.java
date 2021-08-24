package com.gxdun.corerasp.cloud.model;


public enum AppenderMappedLogger {
    ROOT("root", "RASP", "/logs/rasp/rasp.log"),
    ALARM("com.gxdun.corerasp.plugin.checker.alarm", "ALARM", "/logs/alarm/alarm.log"),
    POLICY_ALARM("com.gxdun.corerasp.plugin.checker.policy_alarm", "POLICY_ALARM", "/logs/policy_alarm/policy_alarm.log"),
    JS("com.gxdun.corerasp.plugin.js.log", "PLUGIN", "/logs/plugin/plugin.log"),
    HTTP_ROOT("root", "HTTPEXCEPTIONAPPENDER", ""),
    HTTP_ALARM("com.gxdun.corerasp.plugin.checker.alarm", "HTTPALARMAPPENDER", ""),
    HTTP_POLICY_ALARM("com.gxdun.corerasp.plugin.checker.policy_alarm", "HTTPPOLICYAPPENDER", "");
    private String logger;
    private String appender;
    private String targetPath;


    AppenderMappedLogger(String logger, String appender, String targetPath) {
        this.logger = logger;
        this.appender = appender;
        this.targetPath = targetPath;
    }

    public String getLogger() {
        return logger;
    }

    public String getAppender() {
        return appender;
    }

    public String getTargetPath() {
        return targetPath;
    }
}
