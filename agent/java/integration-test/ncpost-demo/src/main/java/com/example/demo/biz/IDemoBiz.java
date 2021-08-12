package com.example.demo.biz;

import java.util.List;

import com.baomidou.mybatisplus.extension.service.IService;
import com.example.demo.entry.TestEntry;

public interface IDemoBiz extends IService<TestEntry>
{
	public List<TestEntry> sayHello();
	
	/**
	 * 事务测试
	 * @param id
	 * @param name
	 * @throws Exception
	 */
	public void txTest(String id, String name) throws Exception;
}
