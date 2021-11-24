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

package com.gxdun.corerasp.plugin.checker;

import com.gxdun.corerasp.plugin.event.CheckEventListener;
import com.gxdun.corerasp.plugin.info.AttackInfo;
import com.gxdun.corerasp.plugin.info.EventInfo;
import com.gxdun.corerasp.tool.thread.NameableThreadFactory;

import java.util.concurrent.*;

/**
 * Created by tyy on 17-11-22.
 *
 * 攻击检测事件监听器
 */
public class AttackCheckListener implements CheckEventListener {

    ThreadPoolExecutor threadPoolExecutor =  new ThreadPoolExecutor(0, 4, 1, TimeUnit.MINUTES, new LinkedBlockingQueue<Runnable>(), new NameableThreadFactory("attack_alarm_log_"), new ThreadPoolExecutor.CallerRunsPolicy());

    @Override
    public void onCheckUpdate(final EventInfo info) {
        //必须放在外面,拿request信息
        final String eventInfoString = info.toString();
        if (info instanceof AttackInfo) {
            threadPoolExecutor.execute(new Runnable() {
                @Override
                public void run() {
                    Checker.ATTACK_ALARM_LOGGER.info(eventInfoString);
                }
            });
        }

    }

}
