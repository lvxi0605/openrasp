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

package com.gxdun.corerasp.hook.xxe;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.tool.Reflection;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;

/**
 * @description: 关闭stax的XXE entity
 * @author: anyang
 * @create: 2019/04/30 15:25
 */
@HookAnnotation
public class DisableStaxXxeEntity extends DisableXxeEntity {
    @Override
    public boolean isClassMatched(String className) {
        return "com/sun/xml/internal/stream/XMLInputFactoryImpl".equals(className) ||
                "com/ctc/wstx/stax/WstxInputFactory".equals(className);
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(DisableStaxXxeEntity.class, "setFeature", "$0", Object.class);
        insertBefore(ctClass, "createXMLStreamReader", null, src);
        insertBefore(ctClass, "createXMLEventReader", null, src);
    }

    public static void setFeature(Object factory) {
        if (HookHandler.isEnableCurrThreadHook()) {
            String action = getAction();
            if (BLOCK_XXE_DISABLE_ENTITY.equals(action) && getStatus("java_stax")) {
                try {
                    String property = (String) factory.getClass().getField("SUPPORT_DTD").get(null);
                    if (property != null) {
                        Reflection.invokeMethod(factory, "setProperty", new Class[]{String.class, Object.class}, property, false);
                    }
                } catch (Throwable t) {
                    LogTool.traceHookWarn("Stax close xxe entity failed: " + t.getMessage(), t);
                }
            }
        }
    }
}
