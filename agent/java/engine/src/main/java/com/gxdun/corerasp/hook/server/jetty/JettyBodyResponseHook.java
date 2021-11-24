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

package com.gxdun.corerasp.hook.server.jetty;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.server.ServerResponseBodyHook;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.response.HttpServletResponse;
import com.gxdun.corerasp.tool.Reflection;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import com.gxdun.corerasp.tool.model.ApplicationModel;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;
import java.util.HashMap;

/**
 * @author anyang
 * @Description: jetty xss检测的hook点
 * @date 2018/8/1315:13
 */
@HookAnnotation
public class JettyBodyResponseHook extends ServerResponseBodyHook {

    @Override
    public boolean isClassMatched(String className) {
        return "org/eclipse/jetty/server/AbstractHttpConnection".equals(className) ||
                "org/eclipse/jetty/server/Utf8HttpWriter".equals(className) ||
                "org/eclipse/jetty/server/Iso88591HttpWriter".equals(className) ||
                "org/eclipse/jetty/server/EncodingHttpWriter".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src1 = getInvokeStaticSrc(JettyBodyResponseHook.class, "getJettyOutputBuffer", "_generator", Object.class);
        insertBefore(ctClass, "completeResponse", "()V", src1);
        String src2 = getInvokeStaticSrc(JettyBodyResponseHook.class, "getJetty9OutputBuffer", "$1,$2,$3", char[].class, int.class, int.class);
        insertBefore(ctClass, "write", null, src2);
    }


    public static void getJettyOutputBuffer(Object object) {
        boolean isCheckXss = isCheckXss();
        boolean isCheckSensitive = isCheckSensitive();
        boolean isCheckRequest404 = isCheckRequest404();
        if (HookHandler.isEnableXssHook() && (isCheckXss || isCheckSensitive||isCheckRequest404)) {
            HookHandler.disableBodyXssHook();
            HashMap<String, Object> params = new HashMap<String, Object>();
            try {
                if(isCheckRequest404){
                    checkResponseStatus404();
                }

                Object buffer = Reflection.getSuperField(object, "_buffer");
                if (buffer != null) {
                    String content = buffer.toString();
                    params.put("content", content);
                    HttpServletResponse res = HookHandler.responseCache.get();
                    if (res != null) {
                        params.put("content_type", res.getContentType());
                    }
                }
            } catch (Exception e) {
                LogTool.traceHookWarn(ApplicationModel.getServerName() + " xss detectde failed: " +
                        e.getMessage(), e);
            }
            if (HookHandler.requestCache.get() != null && !params.isEmpty()) {
                checkBody(params, isCheckXss, isCheckSensitive);
            }
        }
    }

    public static void getJetty9OutputBuffer(char[] buffer, int offset, int length) {
        boolean isCheckXss = isCheckXss();
        boolean isCheckSensitive = isCheckSensitive();
        boolean isCheckRequest404 = isCheckRequest404();
        if (HookHandler.isEnableXssHook() && (isCheckXss || isCheckSensitive || isCheckRequest404)) {
            HookHandler.disableBodyXssHook();
            if(isCheckRequest404){
                checkResponseStatus404();
            }

            if (buffer != null && length > 0) {
                HashMap<String, Object> params = new HashMap<String, Object>();
                try {
                    char[] temp = new char[length];
                    System.arraycopy(buffer, offset, temp, 0, length);
                    String content = new String(temp);
                    params.put("content", content);
                    HttpServletResponse res = HookHandler.responseCache.get();
                    if (res != null) {
                        params.put("content_type", res.getContentType());
                    }
                } catch (Exception e) {
                    LogTool.traceHookWarn(ApplicationModel.getServerName() + " xss detectde failed: " +
                            e.getMessage(), e);
                }
                if (!params.isEmpty()) {
                    checkBody(params, isCheckXss, isCheckSensitive);
                }
            }
        }
    }
}
