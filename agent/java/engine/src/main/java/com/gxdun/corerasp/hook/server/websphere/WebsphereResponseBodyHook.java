/*
 * Copyright 2017-2018 Baidu Inc.
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

package com.gxdun.corerasp.hook.server.websphere;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.server.ServerResponseBodyHook;
import com.gxdun.corerasp.response.HttpServletResponse;
import com.gxdun.corerasp.messaging.LogTool;
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
 * @Description: websphere的xss检测hook点
 * @date 2018/8/15 14:18
 */
@HookAnnotation
public class WebsphereResponseBodyHook extends ServerResponseBodyHook {
    @Override
    public boolean isClassMatched(String className) {
        return "com/ibm/wsspi/webcontainer/util/BufferedWriter".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(WebsphereResponseBodyHook.class, "getWebsphereOutputBuffer", "$0", Object.class);
        insertBefore(ctClass, "flushChars", "()V", src);
    }

    public static void getWebsphereOutputBuffer(Object object) {
        boolean isCheckXss = isCheckXss();
        boolean isCheckSensitive = isCheckSensitive();
        boolean isCheckRequest404 = isCheckRequest404();
        if (HookHandler.isEnableXssHook() && (isCheckXss || isCheckSensitive||isCheckRequest404)) {
            HookHandler.disableBodyXssHook();
            if(isCheckRequest404){
                checkResponseStatus404();
            }
            HashMap<String, Object> params = new HashMap<String, Object>();
            try {
                char[] buffer = (char[]) Reflection.getField(object, "buf");
                int len = (Integer) Reflection.getField(object, "count");
                char[] temp = new char[len];
                if (buffer != null) {
                    System.arraycopy(buffer, 0, temp, 0, len);
                    String content = new String(temp);
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
            if (!params.isEmpty()) {
                checkBody(params, isCheckXss, isCheckSensitive);
            }
        }
    }
}
