package com.gxdun.corerasp.tool.thread;

import java.util.concurrent.ThreadFactory;
import java.util.concurrent.atomic.AtomicInteger;

public class NameableThreadFactory implements ThreadFactory {
    final AtomicInteger threadNumber = new AtomicInteger(1);
    final String namePrefix;

    public NameableThreadFactory(String namePrefix) {
        this.namePrefix = namePrefix;
    }
    @Override
    public Thread newThread(Runnable r) {
       return new Thread( r,namePrefix + threadNumber.getAndIncrement());
    }
}
