package com.gxdun.corerasp.plugin.checker.local;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;
import com.google.gson.JsonObject;
import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.messaging.ErrorType;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.plugin.checker.AttackChecker;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.plugin.info.AttackInfo;
import com.gxdun.corerasp.plugin.info.EventInfo;

import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

public class RequestPathScanChecker extends AttackChecker {

    private final static Cache<String, RequestInfo> REQUEST_NOT_FOUND_CACHE = CacheBuilder.newBuilder().expireAfterAccess(1, TimeUnit.MINUTES).maximumSize(1000).softValues().build();
    private final static String PLUGIN_NAME ="requestPathScan_check404";
    public RequestPathScanChecker() {
        super();
    }

    public RequestPathScanChecker(boolean canBlock) {
        super(canBlock);
    }

    public static RequestInfo addRequestNotFoundInfo(final String ip, final String url){
        if(isIgnore()){
            return null;
        }
        RequestInfo requestInfo = null;
        try {
            requestInfo = REQUEST_NOT_FOUND_CACHE.get(ip, new Callable<RequestInfo>() {
                @Override
                public RequestInfo call() throws Exception {
                    return new RequestInfo();
                }
            });
        } catch (ExecutionException e) {
            LogTool.error(ErrorType.PLUGIN_ERROR,"",e);
        }
        if(requestInfo==null){
            return requestInfo;
        }

        AtomicBoolean attacked = requestInfo.attacked;
        if (!requestInfo.urls.containsKey(url) && !attacked.get()) {
            synchronized (attacked) {
                if(!attacked.get()) {
                    requestInfo.count.incrementAndGet();
                    requestInfo.urls.put(url, Boolean.TRUE);
                }
            }
        }
        return requestInfo;
    }

    public static RequestInfo getRequestNotFoundInfo(final String ip){
        return REQUEST_NOT_FOUND_CACHE.getIfPresent(ip);
    }

    @Override
    public List<EventInfo> checkParam(CheckParameter checkParameter) {
        LinkedList<EventInfo> result = new LinkedList<EventInfo>();

        String action = ConfigurableChecker.getActionElement(Config.getConfig().getAlgorithmConfig(), PLUGIN_NAME);
        if("ignore".equals(action)){
            return result;
        }

        //String ip = (String) checkParameter.getParam("ip");
        Integer count = (Integer) checkParameter.getParam("count");
        //Set<String> urls = (Set<String>) checkParameter.getParam("urls");

        JsonObject pluginConfig = Config.getConfig().getAlgorithmConfig().getAsJsonObject(PLUGIN_NAME);
        int countPerMinute = ConfigurableChecker.getIntElement(pluginConfig, "countPerMinute");
        if(count!=null && count == countPerMinute){
            result.add(AttackInfo.createLocalAttackInfo(checkParameter, action,"遍历服务器目录" , ConfigurableChecker.getStringElement(pluginConfig, "algorithm")));
        }
        HookHandler.dataThreadHook.set(result);
        return result;
    }

    public static boolean isIgnore(){
        String action = ConfigurableChecker.getActionElement(Config.getConfig().getAlgorithmConfig(), PLUGIN_NAME);
        if("ignore".equals(action)){
            return true;
        }
        return false;
    }

    public static class RequestInfo{
        private final AtomicBoolean attacked = new AtomicBoolean(Boolean.FALSE);
        private final AtomicInteger count = new AtomicInteger(0);
        private final ConcurrentHashMap<String,Boolean> urls = new ConcurrentHashMap<String,Boolean>();
        public AtomicBoolean getAttacked() {
            return attacked;
        }

        public AtomicInteger getCount() {
            return count;
        }

        public ConcurrentHashMap<String, Boolean> getUrls() {
            return urls;
        }
    }

}
