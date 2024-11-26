## basic_test.go
运行测试
```bash
go test ./tests -v -run TestBasic
```

## food_controller_test
1. 首先确保有一个测试数据库：
```sql
CREATE DATABASE test_sustainable_diet;
CREATE USER 'test_user'@'localhost' IDENTIFIED BY 'test_password';
GRANT ALL PRIVILEGES ON test_sustainable_diet.* TO 'test_user'@'localhost';
FLUSH PRIVILEGES;
```

1. 调整测试数据库配置：
如果你的测试数据库配置与代码中的常量不同，请相应修改 `TestDBHost`、`TestDBPort`、`TestDBUser`、`TestDBPassword` 和 `TestDBName` 的值。

1. 运行测试：
```bash
# 运行所有测试
go test ./tests -v

# 运行特定测试
go test ./tests -v -run TestFoodNamesAPI
```
