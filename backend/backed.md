# TODO
 - [ ] 设置Go项目，创建测试路由
 - [ ] 编写简单的数据库模型和接口
 - [ ] 配置文件，数据库连接
 - [ ] dokcer-compose部署测试
 - [ ] 数据库模型
 - [ ] 编写api文档
 - [ ] 用户部分
   - [ ] 用户登陆
 - [ ] 新闻部分
 - [ ] 食谱部分
 - [ ] 单元测试
 - [ ] 性能测试
 - [ ] 
  

# 启动说明

# API设计
提供给前端的API接口，包括请求方法，请求路径，请求参数，请求体，响应体等。



## 新闻及评论
### 1. 视频新闻
#### 1.1 创建视频新闻
* **URL**: `/news/video`
* **Method**: `POST`
* **Content-Type**: `application/json`

##### 请求体
```json
{
  "title": "视频新闻标题",
  "description": "视频新闻简介",
  "video_url": "http://example.com/video.mp4",
  "upload_time": "2024-11-05T12:00:00Z"
}
```
##### 响应示例
* **Status** 201 Created
```json
{
  "id": 1,
  "title": "视频新闻标题",
  "description": "视频新闻简介",
  "upload_time": "2024-11-05T12:00:00Z",
  "like_count": 0,
  "favorite_count": 0,
  "dislike_count": 0,
  "view_count": 0,
  "comments": null,
  "video_url": "http://example.com/video.mp4",
}
```

#### 获取视频新闻详情
* **URL**: `/news/{id}`
* **Method**: `GET`

##### 路径参数
* id: 新闻唯一的标识符

##### 响应示例
* **Status** 200 OK
```json
{
}
```


# 模型设计
1. 用户表
2. 新闻表
3. 回复表


