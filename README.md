# CSV2SQL

将csv转换为sql插入语句

## 用途

此软件用于方便导入Excel里的数据

## 命令行参数

```
Usage:
  -i string
        输入CSV文件路径，默认使用stdin输入
  -o string
        输出SQL文件路径，默认使用stdout输出
  -t string
        要插入的表名，默认是`table_name`
```

## 输入CSR格式

第一行是对应数据库里的字段名，之后的每一行对应数据库里的数据

编码使用UTF-8

example:

```csv
id,name,password
1,tony,123456
2,root,abcdefg
3,admin,aaccbbaa
```

## 输出

目前只支持了mysql

example:

`go run main.go -i example.csv -t users`

```sql
INSERT INTO `users` (`id`, `name`, `password`)
VALUES
    ('1', 'tony', '123456'),
    ('2', 'root', 'abcdefg'),
    ('3', 'admin', 'aaccbbaa');
```