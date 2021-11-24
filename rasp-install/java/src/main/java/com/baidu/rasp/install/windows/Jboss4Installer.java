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

package com.baidu.rasp.install.windows;

import com.baidu.rasp.RaspError;
import com.baidu.rasp.install.BaseStandardInstaller;

import java.util.Scanner;
import java.util.regex.Pattern;

import static com.baidu.rasp.RaspError.E10001;

/**
 * Created by OpenRASP on 5/15/17.
 */
public class Jboss4Installer extends BaseStandardInstaller {

    private static String OPENRASP_CONFIG =
            "rem BEGIN OPENRASP - DO NOT MODIFY" + LINE_SEP +
                    "set JAVA_OPTS=-javaagent:%JBOSS_HOME%\\rasp\\rasp.jar %JAVA_OPTS%" + LINE_SEP +
                    "rem END OPENRASP" + LINE_SEP;
    private static Pattern OPENRASP_REGEX = Pattern.compile(".*(\\s*OPENRASP\\s*|JAVA_OPTS.*\\\\rasp\\\\).*");

    public Jboss4Installer(String serverName, String serverRoot) {
        super(serverName, serverRoot);
    }

    @Override
    protected String getInstallPath(String serverRoot) {
        return serverRoot + "\\rasp";
    }

    @Override
    protected String getScript(String installPath) {
        return installPath + "\\..\\bin\\run.bat";
    }

    @Override
    protected String modifyStartScript(String content) throws RaspError {
        int modifyConfigState = NOTFOUND;
        StringBuilder sb = new StringBuilder();
        Scanner scanner = new Scanner(content);
        while (scanner.hasNextLine()) {
            String line = scanner.nextLine();
            if (FOUND == modifyConfigState) {
                sb.append(OPENRASP_CONFIG);
                modifyConfigState = DONE;
            }
            if (DONE == modifyConfigState) {
                if (OPENRASP_REGEX.matcher(line).matches()) {
                    continue;
                }
            }
            if (line.startsWith("rem Setup JBoss specific properties") && NOTFOUND == modifyConfigState) {
                modifyConfigState = FOUND;
            }

            sb.append(line).append(LINE_SEP);
        }
        if (NOTFOUND == modifyConfigState) {
            throw new RaspError(E10001 + "\"rem Setup JBoss specific properties\"");
        }
        return sb.toString();
    }

}
