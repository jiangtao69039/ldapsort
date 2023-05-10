
## LDAP SORT

#### 描述

  使用ldapsearch命令从ldap服务器中查出的数据是无顺序的,查出的数据不会以DN为顺序存放.
  所以再导入回去的时候,通常会失败. 因为DN之间有先后依赖的顺序.
  
  例子:
  DN: ou=OU1  
   
  DN: people=admin001,ou=OU1   
  
  上面两条数据在导入ldap时,必须先导入第一条,再导入第二条.
  
#### 功能
  对一个 .ldif 数据文件按照 DN 顺序排序, 生成的文件可正确的导入 ldap