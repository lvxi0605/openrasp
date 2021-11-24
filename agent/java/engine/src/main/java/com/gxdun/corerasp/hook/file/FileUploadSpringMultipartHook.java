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
import com.gxdun.corerasp.hook.AbstractClassHook;
import com.gxdun.corerasp.messaging.LogTool;
import com.gxdun.corerasp.plugin.checker.CheckParameter;
import com.gxdun.corerasp.tool.Reflection;
import javassist.CannotCompileException;
import javassist.CtClass;
import javassist.NotFoundException;
import org.apache.commons.io.IOUtils;

import java.io.IOException;
import java.io.InputStream;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * @description: Jersey 文件上传
 * @author: anyang
 * @create: 2019/01/28 12:43
 */
//@HookAnnotation
public class FileUploadSpringMultipartHook extends AbstractClassHook {
    @Override
    public boolean isClassMatched(String className) {
        return "org/springframework/web/multipart/support/StandardMultipartHttpServletRequest$StandardMultipartFile".equals(className);
    }

    @Override
    public String getType() {
        return "fileUpload";
    }

    @Override
    protected void hookMethod(CtClass ctClass) throws IOException, CannotCompileException, NotFoundException {
        String srcGetBytes = getInvokeStaticSrc(FileUploadSpringMultipartHook.class, "checkGetBytes", "$0,$_", Object.class);
        insertAfter(ctClass, "getBytes", "()[B", srcGetBytes);
//
//        String srcGetInputStream = getInvokeStaticSrc(SpringMultipart.class, "checkGetInputStream", "$0,$_", Object.class);
//        insertAfter(ctClass, "getInputStream", "()Ljava/io/InputStream;", srcGetInputStream);
    }

    public static void checkGetBytes(Object that,byte[] content) {
        String filename = Reflection.invokeStringMethod(that, "getOriginalFilename", new Class[]{});
        if (filename == null) {
            return;
        }
        String name = Reflection.invokeStringMethod(that, "getName", new Class[]{});
        HashMap<String, Object> params = new HashMap<String, Object>();
        params.put("filename", filename);
        if (content.length > 4 * 1024) {
            content = Arrays.copyOf(content, 4 * 1024);
        }
        params.put("content", new String(content));
        params.put("name", name);
        HookHandler.doCheck(CheckParameter.Type.FILEUPLOAD, params);
    }

    public static void checkGetInputStream(Object that,InputStream inputStream) {
        String filename = Reflection.invokeStringMethod(that, "getOriginalFilename", new Class[]{});
        if (filename == null) {
            return;
        }
        String name = Reflection.invokeStringMethod(that, "getName", new Class[]{});
        HashMap<String, Object> params = new HashMap<String, Object>();
        params.put("filename", filename);

        try {
            byte[] content = IOUtils.toByteArray(inputStream);
            if (content.length > 4 * 1024) {
                content = Arrays.copyOf(content, 4 * 1024);
            }
            params.put("content", new String(content));
        } catch (IOException e) {
            LogTool.traceHookWarn(e.getMessage(), e);
            params.put("content", "[rasp error:" + e.getMessage() + "]");
        }
        params.put("name", name);
        HookHandler.doCheck(CheckParameter.Type.FILEUPLOAD, params);
    }


    public static void checkFileUpload0(Object result) {
        if (result != null) {
            Map<String, List<Object>> map = (Map<String, List<Object>>) result;
            for (Map.Entry<String, List<Object>> entry : map.entrySet()) {
                Object o = entry.getValue().get(0);
                Object contentDisposition = Reflection.invokeMethod(o, "getFormDataContentDisposition", new Class[]{});
                String name = Reflection.invokeStringMethod(contentDisposition, "getFileName", new Class[]{});
                if (name != null) {
                    HashMap<String, Object> params = new HashMap<String, Object>();
                    params.put("filename", name);
                    InputStream inputStream = (InputStream) Reflection.invokeMethod(o, "getValueAs", new Class[]{Class.class}, InputStream.class);
                    try {
                        byte[] content = IOUtils.toByteArray(inputStream);
                        if (content.length > 4 * 1024) {
                            content = Arrays.copyOf(content, 4 * 1024);
                        }
                        params.put("content", new String(content));
                    } catch (IOException e) {
                        LogTool.traceHookWarn(e.getMessage(), e);
                        params.put("content", "[rasp error:" + e.getMessage() + "]");
                    }
                    params.put("name", "");
                    HookHandler.doCheck(CheckParameter.Type.FILEUPLOAD, params);
                }
            }
        }
    }
}
