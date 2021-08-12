package com.example.demo.mapper;

import java.util.List;

import org.apache.ibatis.annotations.Param;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.example.demo.entry.TestEntry;

public interface DemoMapper extends BaseMapper<TestEntry>
{
	public List<TestEntry> getList();
	
	
	public void updateById(@Param("id") String id, @Param("name") String name);
	
	public void insert(@Param("id") String id, @Param("name") String name);
}
