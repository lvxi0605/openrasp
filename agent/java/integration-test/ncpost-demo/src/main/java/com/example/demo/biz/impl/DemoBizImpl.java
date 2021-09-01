package com.example.demo.biz.impl;

import java.util.List;
import java.util.UUID;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.example.demo.biz.IDemoBiz;
import com.example.demo.entry.TestEntry;
import com.example.demo.mapper.DemoMapper;

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



	
}
