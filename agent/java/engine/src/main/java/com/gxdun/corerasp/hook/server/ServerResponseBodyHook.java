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

package com.gxdun.corerasp.hook.server;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.exceptions.SecurityException;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.plugin.checker.local.RequestPathScanChecker;
import com.gxdun.corerasp.plugin.info.EventInfo;
import com.gxdun.corerasp.request.AbstractRequest;
import com.gxdun.corerasp.response.HttpServletResponse;

import java.util.HashMap;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

import com.gxdun.corerasp.tool.RequestUtil;
import com.gxdun.corerasp.tool.Sampler;

/**
 * @author anyang
 * @Description: xss检测基类
 * @date 2018/8/15 15:37
 */
public abstract class ServerResponseBodyHook extends AbstractClassHook {
    private static Sampler sampler = new Sampler();

    @Override
    public String getType() {
        return "xss";
    }

    protected static boolean isCheckXss() {
        if (HookHandler.requestCache.get() != null && HookHandler.responseCache.get() != null) {
            String contentType = HookHandler.responseCache.get().getContentType();
            return contentType == null || contentType.startsWith(HttpServletResponse.CONTENT_TYPE_HTML_VALUE);
        }
        return false;
    }

    protected static boolean isCheckSensitive() {
        Config config = Config.getConfig();
        if (config.getResponseSamplerInterval() <= 0 || config.getResponseSamplerBurst() <= 0) {
            return false;
        }
        if (HookHandler.responseCache.get() != null) {
            String contentType = HookHandler.responseCache.get().getContentType();
            if (contentType != null && (contentType.contains("video") || contentType.contains("audio")
                    || contentType.contains("image"))) {
                return false;
            }
        }
        sampler.update(config.getResponseSamplerInterval(), config.getResponseSamplerBurst());
        // 限速检测
        return sampler.check();
    }

    protected static boolean isCheckRequest404() {
        return !RequestPathScanChecker.isIgnore();
    }

    protected static void checkResponseStatus404(){
        final HttpServletResponse httpServletResponse = HookHandler.responseCache.get();
        if(httpServletResponse!=null){
            final int status = httpServletResponse.getStatus();
            if(status==404){
                AbstractRequest request = HookHandler.requestCache.get();
                String ip = RequestUtil.getIpAddr(request);
                RequestPathScanChecker.RequestInfo requestInfo = RequestPathScanChecker.addRequestNotFoundInfo(ip,request.getRequestURI());
                AtomicBoolean attacked = requestInfo.getAttacked();
                if(requestInfo!=null && !attacked.get()) {
                    HashMap<String, Object> checkParams = new HashMap<String, Object>();
                    checkParams.put("ip", ip);
                    checkParams.put("count", requestInfo.getCount().get());
                    checkParams.put("urls", requestInfo.getUrls().keySet());
                    synchronized (attacked) {
                        if(attacked.get()){
                            return;
                        }
                        try {
                            HookHandler.doCheck(CheckParameter.Type.REQUEST_PATH_SCAN, checkParams);
                        } finally {
                            List<EventInfo> eventInfos = HookHandler.dataThreadHook.get();
                            if (eventInfos != null && !eventInfos.isEmpty()) {
                                //被拦截或者记录日志了,不要重复记录日志
                                attacked.set(true);
                            }
                        }
                    }
                }
            }
        }
    }

    protected static void checkBody(HashMap<String, Object> params, boolean isCheckXss, boolean isCheckSensitive) {
        if (isCheckXss) {
            HookHandler.doCheck(CheckParameter.Type.XSS_USERINPUT, params);
        }
        if (isCheckSensitive) {
            params.remove("buffer");
            params.remove("encoding");
            params.remove("content_length");
            HookHandler.doCheck(CheckParameter.Type.RESPONSE, params);
        }




    }
}
