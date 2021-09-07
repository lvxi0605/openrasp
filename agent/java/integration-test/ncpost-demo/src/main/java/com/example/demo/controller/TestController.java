package com.example.demo.controller;

import com.baomidou.mybatisplus.extension.api.R;
import com.example.demo.biz.IDemoBiz;
import com.example.demo.entry.TestEntry;
import com.mysql.jdbc.PreparedStatement;
import javassist.ClassClassPath;
import javassist.LoaderClassPath;
import org.apache.ibatis.javassist.ClassPath;
import org.apache.ibatis.javassist.ClassPool;
import org.apache.ibatis.javassist.CtClass;
import org.apache.ibatis.javassist.CtMethod;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.UUID;

@RestController
@RequestMapping("/test")
public class TestController 
{
	@Autowired
	private IDemoBiz demoBiz;
	
	@GetMapping("/hello")
	@ResponseBody
	public Object hello()
	{
		return demoBiz.sayHello();
	}

	@GetMapping("/getListById")
	@ResponseBody
	public Object getListById(String id) throws Exception {
		demoBiz.txTest(UUID.randomUUID().toString(),"alert");
		return demoBiz.getListById(id);
	}

	@GetMapping("/testSQLXSS")
	@ResponseBody
	public void testSQLXSS() throws Exception {
		 demoBiz.testSQLXSS();
	}



	@PostMapping("/upload")
    public String SingleFileUpLoad(@RequestParam("myfile") MultipartFile file) 
    {
    	System.out.println(file.getName());
        return "1";
    }
    
	public static void main(String[] args) {
		try {
		  
			  ClassPool pool = ClassPool.getDefault();
		        //获取一个ctClass对象
		     CtClass ctClass = pool.makeClass("com.mysql.jdbc.PreparedStatement");
		        
		     CtMethod[] menths;
			
				menths = ctClass.getDeclaredMethods("executeInternal");
			
		for (CtMethod method : menths) 
			{
			CtClass[] params = method.getParameterTypes();
			for (CtClass param : params) 
			{ 
				System.out.println(param);
			}
		  
		 }
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
    @PostMapping("/insert")
    @ResponseBody
    public R insert(@RequestBody TestEntry entry)
    {
    	try 
    	{
			demoBiz.txTest(entry.getId(),entry.getName());
			/*
			 * ClassPool pool = ClassPool.getDefault();
			 * 
			 * ClassPool classPool = new ClassPool(); addLoader(classPool,
			 * PreparedStatement.class.getClassLoader());
			 * 
			 * 
			 * //获取一个ctClass对象
			 * 
			 * CtClass ctClass = pool.makeClass("com.mysql.jdbc.PreparedStatement");
			 * FileOutputStream out = new
			 * FileOutputStream("D:\\workspaces\\spring\\java\\debug\\zjsz.txt");
			 * out.write(ctClass.toBytecode()); out.flush(); out.close();
			 */
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
    	R<Object> kj = new R<>();
    	kj.ok();
    	return kj;
    	}
	
    private void addLoader(ClassPool classPool, ClassLoader loader) {
        classPool.appendSystemPath();
        classPool.appendClassPath((ClassPath) new ClassClassPath(PreparedStatement.class));
        if (loader != null) {
            classPool.appendClassPath((ClassPath) new LoaderClassPath(loader));
        }
    }
}
