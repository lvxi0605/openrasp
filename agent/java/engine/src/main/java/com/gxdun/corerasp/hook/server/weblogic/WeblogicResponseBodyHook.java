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

package com.gxdun.corerasp.hook.server.weblogic;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.server.ServerResponseBodyHook;
import com.gxdun.corerasp.response.HttpServletResponse;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import com.gxdun.corerasp.tool.model.ApplicationModel;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;
import java.nio.CharBuffer;
import java.util.HashMap;

/**
 * @description: weblogic的xss检测hook点
 * @author: anyang
 * @create: 2018/09/05 15:06
 */
@HookAnnotation
public class WeblogicResponseBodyHook extends ServerResponseBodyHook {
    @Override
    public boolean isClassMatched(String className) {
        return "weblogic/servlet/internal/CharsetChunkOutput".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(WeblogicResponseBodyHook.class, "getWeblogicOutputBuffer", "$1", CharBuffer.class);
        insertBefore(ctClass, "write", "(Ljava/nio/CharBuffer;)V", src);
    }

    public static void getWeblogicOutputBuffer(CharBuffer buffer) {
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
            if (!params.isEmpty()) {
                checkBody(params, isCheckXss, isCheckSensitive);
            }
        }
    }
}
