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

package com.gxdun.corerasp.hook.file;

import com.gxdun.corerasp.HookHandler;
import com.gxdun.corerasp.config.Config;
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.tool.FileUtil;
import com.gxdun.corerasp.tool.annotation.HookAnnotation;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;

/**
 * Created by zhuming01 on 5/31/17. All rights reserved
 */
@HookAnnotation
public class FileInputStreamHook extends AbstractClassHook {
    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#getType()
     */
    @Override
    public String getType() {
        return "readFile";
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#isClassMatched(String)
     */
    @Override
    public boolean isClassMatched(String className) {
        return "java/io/FileInputStream".equals(className);
    }

    /**
     * (none-javadoc)
     *
     * @see com.gxdun.corerasp.hook.AbstractClassHook#hookMethod(CtClass)
     */
    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String src = getInvokeStaticSrc(FileInputStreamHook.class, "checkReadFile", "$1", File.class);
        insertBefore(ctClass.getConstructor("(Ljava/io/File;)V"), src);
    }

    /**
     * 文件读取hook点
     *
     * @param file 文件对象
     */
    public static void checkReadFile(File file) {
        boolean checkSwitch = Config.getConfig().getPluginFilter();
        if (file != null) {
            if (checkSwitch && !file.exists()) {
                return;
            }
            HashMap<String, Object> params = new HashMap<String, Object>();
            params.put("path", file.getPath());

            String path;
            try {
                path = file.getCanonicalPath();
            } catch (Exception e) {
                path = file.getAbsolutePath();
            }
            if (path.endsWith(".class")) {
                return;
            }
            params.put("realpath", FileUtil.getRealPath(file));

            HookHandler.doCheck(CheckParameter.Type.READFILE, params);
        }
    }
}
