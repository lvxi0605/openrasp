package com.example.demo.biz.impl;

import java.io.*;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.UUID;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.example.demo.biz.IDemoBiz;
import com.example.demo.entry.TestEntry;
import com.example.demo.mapper.DemoMapper;
import org.springframework.util.ResourceUtils;
import org.springframework.util.StringUtils;

@Service
public class DemoBizImpl extends ServiceImpl<DemoMapper,TestEntry> implements IDemoBiz
{

	@Override
	public List<TestEntry> sayHello() 
	{
		try {
			return baseMapper.getList(); 
		} catch (Exception e) {
			return null;
		}
		
	}


	@Override
	public TestEntry getListById(String id)
	{
			return baseMapper.getListById(id);
	}
	
	@Transactional
	@Override
	public void txTest(String id, String name) throws Exception
	{
		baseMapper.updateById(id,name);
		baseMapper.insert(UUID.randomUUID().toString(), name);
		if ("123".equals(id)){
			throw new Exception("呵呵呵");
		}
	}

	@Transactional
	@Override
	public void testSQLXSS() throws IOException {

		File file = ResourceUtils.getFile("classpath:xss_test.text");
		try(
		FileInputStream fi = new FileInputStream(file);
		BufferedReader bufferedReader = new BufferedReader((new InputStreamReader(fi)));
		) {
			int i=0;
			String str = null;
			while ((str = bufferedReader.readLine()) != null) {
				i++;
				String testValue = str.trim();
				if (!StringUtils.hasText(testValue)) {
					continue;
				}

				baseMapper.insert(LocalDate.now().format(DateTimeFormatter.BASIC_ISO_DATE)+String.format("%05d",i), testValue);

			}
		}


	}



	
}
