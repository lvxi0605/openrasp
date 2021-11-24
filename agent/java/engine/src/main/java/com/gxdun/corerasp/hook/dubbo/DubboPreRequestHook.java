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

package com.gxdun.corerasp.hook.dubbo;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.IOException;

/**
 * @author anyang
 * @Description: dubbo请求预处理hook点
 * @date 2018/8/1315:13
 */
@HookAnnotation
public class DubboPreRequestHook extends AbstractClassHook {

    public DubboPreRequestHook() {

        couldIgnore = false;
    }

    @Override
    public boolean isClassMatched(String className) {
        return "com/alibaba/dubbo/rpc/filter/GenericFilter".equals(className);
    }

    @Override
    public String getType() {
        return "dubbo_preRequest";
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {

        String src = getInvokeStaticSrc(HookHandler.class, "onDubboExit", "");
        insertBefore(ctClass, "invoke", null, src);

    }
}
