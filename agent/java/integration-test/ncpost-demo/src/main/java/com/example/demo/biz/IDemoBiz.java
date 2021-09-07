package com.example.demo.biz;

import java.io.IOException;
import java.util.List;

import com.baomidou.mybatisplus.extension.service.IService;
import com.example.demo.entry.TestEntry;

public interface IDemoBiz extends IService<TestEntry>
{
	public List<TestEntry> sayHello();

	TestEntry getListById(String id);

	/**
	 * 事务测试
	 * @param id
	 * @param name
	 * @throws Exception
	 */
	public void txTest(String id, String name) throws Exception;

	void testSQLXSS() throws IOException;
}
