## 项目说明
1. 项目使用gin框架编写

## 接口响应设计
1. 接口统一使用json响应格式：
   ```json
   {
      "code": 0,
      "msg": "成功",
      "data": {
         
      }
   }
   ```
2. code字段说明
   1. 0:成功
   2. >0:失败
3. msg：提示信息
4. data：返回数据

## 接口鉴权设计
1. 鉴权使用http Authorization头，格式为：Bearer <token>
2. token使用jwt
3. token过期时间：1天


## 目录结构
1. main.go: 主程序文件
2. config: 配置文件目录
   1. config.yaml: 配置文件
3. route: 路由配置目录
    1. web.go 路由文件
4. route/middleware: 中间件文件目录
5. model: 模型文件目录
    - db 数据库模型
    - dto 数据传输对象
6. controller: 控制器文件目录
7. repo: 业务仓库文件目录
8. utils: 工具文件目录
9. lib: 库文件目录

## 项目依赖
1. gin lib: github.com/gin-gonic/gin
2. mysql lib: gorm.io/gorm
3. jwt lib: github.com/golang-jwt/jwt
4. captcha lib: github.com/dchest/captcha
5. log lib: go.uber.org/zap
6. yaml lib: gopkg.in/yaml.v3
