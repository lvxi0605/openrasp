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

package com.gxdun.corerasp.hook.server.wildfly;

import com.gxdun.corerasp.hook.server.ServerPreRequestHook;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

/**
 * Created by izpz on 18-10-26.
 * <p>
 * wildfly 请求预处理hook点
 */
@HookAnnotation
public class UndertowPreRequestHook extends ServerPreRequestHook {

    /**
     * 用于判断类名与当前需要hook的类是否相同
     *
     * @param className 用于匹配的类名
     * @return 是否匹配
     */
    @Override
    public boolean isClassMatched(String className) {
        return "io/undertow/servlet/handlers/ServletInitialHandler".equals(className);
    }


    /**
     * hook 方法
     *
     * @param ctClass hook 点所在的类
     * @param src     加入 hook点的代码
     */
    @Override
    protected void hookMethod(CtClass ctClass, String src) throws NotFoundException, CannotCompileException {
        insertBefore(ctClass, "dispatchRequest", null, src);
    }
}
