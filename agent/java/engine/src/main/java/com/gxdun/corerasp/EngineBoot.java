/*
 * Copyright 2021 CORE SHIELD Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.gxdun.corerasp;

import com.gxdun.corerasp.cloud.CloudManager;
import com.gxdun.corerasp.cloud.managehook.config.HookManageConfiguration;
import com.gxdun.corerasp.cloud.model.CloudCacheModel;
import com.gxdun.corerasp.cloud.utils.CloudUtils;
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.messaging.LogConfig;
import com.gxdun.corerasp.plugin.checker.CheckerManager;
import com.gxdun.corerasp.plugin.js.JS;
import com.gxdun.corerasp.tool.cpumonitor.CpuMonitorManager;
import com.gxdun.corerasp.tool.model.BuildRASPModel;
import com.gxdun.corerasp.transformer.CustomClassTransformer;
import com.gxdun.corerasp.v8.CrashReporter;
import com.gxdun.corerasp.v8.Loader;
import org.apache.commons.io.FileUtils;
import org.apache.log4j.Logger;

import java.io.File;
import java.lang.instrument.Instrumentation;
import java.lang.instrument.UnmodifiableClassException;
import java.util.Random;

/**
 * Created by tyy on 18-1-24.
 *
 * CoreRasp 引擎启动类
 */
public class EngineBoot implements Module {

    private CustomClassTransformer transformer;

    @Override
    public void start(String mode, Instrumentation inst) throws Exception {
        System.out.println("\033[0;32m" +
                "   ______                              ____     ___    _____    ____ \n" +
                "  / ____/  ____    _____  ___         / __ \\   /   |  / ___/   / __ \\\n" +
                " / /      / __ \\  / ___/ / _ \\       / /_/ /  / /| |  \\__ \\   / /_/ /\n" +
                "/ /___   / /_/ / / /    /  __/      / _, _/  / ___ | ___/ /  / ____/ \n" +
                "\\____/   \\____/ /_/     \\___/      /_/ |_|  /_/  |_|/____/  /_/      \033[0m\n");

        Agent.readVersion();
        BuildRASPModel.initRaspInfo(Agent.projectVersion, Agent.buildTime, Agent.gitCommit);

        String versionMessage = "\033[1;43m CoreRASP \033[0m \033[1;33mversion: " + Agent.projectVersion + " \033[0;32mbuild: GitCommit="
            + Agent.gitCommit + " date=" + Agent.buildTime + "\033[0m\n";
        System.out.println(versionMessage);

        try {
            Loader.load();
        } catch (Exception e) {
            System.out.println("[CoreRASP] Failed to load native library");
            e.printStackTrace();
            return;
        }
        if (!loadConfig()) {
            return;
        }
        //BuildRASPModel
        // 初始化插件系统
        if (!JS.Initialize()) {
            return;
        }
        CheckerManager.init();
        initTransformer(inst);
        if (CloudUtils.checkCloudControlEnter()) {
            CrashReporter.install(Config.getConfig().getCloudAddress() + "/v1/agent/crash/report",
                    Config.getConfig().getCloudAppId(), Config.getConfig().getCloudAppSecret(),
                    CloudCacheModel.getInstance().getRaspId());
            // TODO 001 开启代码层开关
            //initHookManagerConfig();
        }
        deleteTmpDir();
        String message = "[CoreRASP] Engine Initialized [" + Agent.projectVersion + " (build: GitCommit="
                + Agent.gitCommit + " date=" + Agent.buildTime + ")]";
        Logger.getLogger(EngineBoot.class.getName()).info(message);
    }

    @Override
    public void release(String mode) {
        CloudManager.stop();
        CpuMonitorManager.release();
        if (transformer != null) {
            transformer.release();
        }
        JS.Dispose();
        CheckerManager.release();
        String message = "[CoreRASP] Engine Released [" + Agent.projectVersion + " (build: GitCommit="
                + Agent.gitCommit + " date=" + Agent.buildTime + ")]";
        System.out.println(message);
    }

    private void deleteTmpDir() {
        try {
            File file = new File(Config.baseDirectory + File.separator + "jar_tmp");
            if (file.exists()) {
                FileUtils.deleteDirectory(file);
            }
        } catch (Throwable t) {
            Logger.getLogger(EngineBoot.class.getName()).warn("failed to delete jar_tmp directory: " + t.getMessage());
        }
    }

    /**
     * 初始化配置
     *
     * @return 配置是否成功
     */
    private boolean loadConfig() throws Exception {
        LogConfig.ConfigFileAppender();
        //单机模式下动态添加获取删除syslog
        if (!CloudUtils.checkCloudControlEnter()) {
            LogConfig.syslogManager();
        } else {
            System.out.println("[CoreRASP] RASP ID: " + CloudCacheModel.getInstance().getRaspId());
        }
        return true;
    }

    /**
     * 初始化类字节码的转换器
     *
     * @param inst 用于管理字节码转换器
     */
    private void initTransformer(Instrumentation inst) throws UnmodifiableClassException {
        transformer = new CustomClassTransformer(inst);
        transformer.retransform();
    }

    private void initHookManagerConfig()
    {
        for (AbstractClassHook hook : transformer.getHooks())
        {
            HookManageConfiguration.addHookCodeInformation(hook.getClass());
        }
        HookManageConfiguration.initCommandCahe();
    }
}
