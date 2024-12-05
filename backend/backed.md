
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
4. 食物
   

## 食物以及食谱部分设计思路
食物：
  用户输入类型，下拉选取（基于已有的数据库）
  用户购买的重量
  用户购买的价格
当用户输入类型之后，点击确定之后，检索数据库，查看平均碳排放和平均价格，营养成分

食物集合：
  关联的食物
  总碳排放
  总营养成分
  关联的菜谱（需要某种算法来从菜谱数据库来查询）

菜谱：（应该是来自数据库）
  制作方法
  关联的制作视频
  相关联的食物类型

# 可持续饮食助手 API 文档 v1.0

## 目录
- [通用说明](#通用说明)
- [用户相关](#用户相关)
- [新闻相关](#新闻相关)
- [家庭相关](#家庭相关)

## 通用说明

### Base URL
```
https://api.example.com
```

### 认证方式
大多数API端点需要JWT认证。在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

### 状态码说明
- 200: 成功
- 201: 创建成功
- 400: 请求参数错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误

### 响应格式
所有响应均为JSON格式：

**成功响应 (HTTP 2xx)**
```json
{
    "message": "操作成功信息",
    "data": {
        // 具体数据
    }
}
```

**错误响应 (HTTP 4xx/5xx)**
```json
{
    "error": "错误信息"
}
```

## 用户相关 API

### 1. 微信登录/注册
**POST** `/users/auth`

使用微信登录凭证进行登录或注册。

**请求体**：
```json
{
    "code": "string",     // 必需，微信登录凭证
    "nickname": "string"  // 可选，用户昵称
}
```

**成功响应 (200)**：
```json
{
    "token": "jwt_token",
    "user": {
        "id": 1,
        "nickname": "用户昵称",
        "avatar_url": "头像地址"
    }
}
```

### 2. 设置用户昵称
**PUT** `/users/{id}/set_nickname`

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "nickname": "string"  // 必需，新昵称
}
```

**成功响应 (200)**：
```json
{
    "message": "Nickname updated successfully",
    "nickname": "新昵称"
}
```

### 3. 设置用户头像
**PUT** `/users/{id}/set_avator`

**请求头**：
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**表单参数**：
- `avatar`: File (必需，jpg格式图片)

**成功响应 (200)**：
```json
{
    "message": "Avatar updated successfully",
    "avatar_url": "头像路径"
}
```

## 新闻相关 API

### 1. 上传新闻
**POST** `/news/upload`

**请求头**：
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**：
```json
{
    "title": "string",          // 必需，新闻标题
    "description": "string",    // 必需，新闻描述
    "news_type": "string",      // 必需，新闻类型（video/regular/external）
    "video": {                  // 视频类型新闻必需
        "video_url": "string"
    },
    "paragraphs": [            // 普通新闻可选
        {
            "content": "string",
            "order": number
        }
    ],
    "resources": [             // 普通新闻可选
        {
            "url": "string",
            "type": "string",
            "description": "string"
        }
    ],
    "external_link": "string", // 外部新闻必需
    "tags": ["string"]         // 可选，标签列表
}
```

**成功响应 (201)**：
```json
{
    "id": 1,
    "title": "标题",
    "description": "描述",
    "upload_time": "2024-12-05T00:00:00Z",
    "view_count": 0,
    "like_count": 0,
    "favorite_count": 0,
    "dislike_count": 0,
    "share_count": 0,
    "news_type": "新闻类型",
    "tags": ["标签1", "标签2"]
}
```

### 2. 获取新闻详情
**GET** `/news/{id}`

无需认证。

**路径参数**：
- `id`: number (必需，新闻ID)

**响应 (200)**：
```json
{
    "id": 1,
    "title": "标题",
    "description": "描述",
    "upload_time": "时间",
    "view_count": 0,
    "like_count": 0,
    "favorite_count": 0,
    "dislike_count": 0,
    "share_count": 0,
    "news_type": "类型",
    "author_name": "作者",
    "author_avatar": "作者头像",
    "tags": ["标签"],
    "is_liked": 0,
    "is_favorited": 0,
    "is_disliked": 0,
    "comments": [],
    // 根据news_type会有额外字段：
    "video": {},          // 视频类型
    "paragraphs": [],     // 普通类型
    "resources": [],      // 普通类型
    "external_link": ""   // 外部链接类型
}
```

### 3. 新闻互动 API
以下所有互动API都需要认证头：
```
Authorization: Bearer <token>
```

#### 点赞相关
- **点赞**: POST `/news/like`
- **取消点赞**: POST `/news/cancel_like`

#### 收藏相关
- **收藏**: POST `/news/favourite`
- **取消收藏**: POST `/news/cancel_favourite`

#### 点踩相关
- **点踩**: POST `/news/dislike`
- **取消点踩**: POST `/news/cancel_dislike`

#### 查看新闻
- **记录浏览**: POST `/news/view`

所有互动API的响应格式：
```json
{
    "message": "操作成功消息",
    "xxx_count": 新数量  // like_count/favorite_count/dislike_count
}
```

### 4. 评论 API
**POST** `/news/comment`

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "news_id": 1,        // 必需，新闻ID
    "content": "string", // 必需，评论内容
    "is_reply": false,   // 必需，是否为回复
    "parent_id": null    // 可选，父评论ID（回复评论时必需）
}
```

**成功响应 (201)**：
```json
{
    "message": "Comment added successfully",
    "comment": {
        "id": 1,
        "content": "评论内容",
        "user_id": 1,
        "publish_time": "2024-12-05T00:00:00Z",
        "like_count": 0
    }
}
```

## 家庭相关 API

所有家庭相关API都需要以下认证头：
```
Authorization: Bearer <token>
```

### 1. 创建家庭
**POST** `/families/create`

**请求体**：
```json
{
    "name": "string"  // 必需，家庭名称
}
```

**成功响应 (201)**：
```json
{
    "message": "Family created successfully",
    "family": {
        "id": 1,
        "name": "家庭名称",
        "family_id": "唯一标识"
    }
}
```

### 2. 查看家庭详情
**GET** `/families/details`

**成功响应 (200)**：
```json
{
    "id": 1,
    "name": "名称",
    "family_id": "标识",
    "member_count": 1,
    "admins": [
        {
            "id": 1,
            "nickname": "昵称",
            "avatar_url": "头像"
        }
    ],
    "members": [
        {
            "id": 2,
            "nickname": "昵称",
            "avatar_url": "头像"
        }
    ]
}
```

### 3. 搜索家庭
**GET** `/families/search`

**查询参数**：
- `family_id`: string (必需，家庭唯一标识)

**成功响应 (200)**：
```json
{
    "id": 1,
    "name": "名称",
    "family_id": "标识",
    "member_count": 1
}
```

### 4. 家庭成员管理 API

#### 申请加入家庭
**POST** `/families/{id}/join`

**路径参数**：
- `id`: number (必需，家庭ID)

**成功响应 (200)**：
```json
{
    "message": "Join request sent successfully"
}
```

#### 批准加入申请
**POST** `/families/admit`

**请求体**：
```json
{
    "user_id": 1  // 必需，申请用户的ID
}
```

#### 拒绝加入申请
**POST** `/families/reject`

**请求体**：
```json
{
    "user_id": 1  // 必需，申请用户的ID
}
```

#### 取消加入申请
**POST** `/families/cancel_join`

#### 查看待处理的家庭申请
**POST** `/families/pending_family_details`

**成功响应 (200)**：
```json
{
    "id": 1,
    "name": "家庭名称",
    "token": "家庭标识"
}
```

### 错误处理
所有API在遇到错误时会返回对应的HTTP状态码和错误信息：

```json
{
    "error": "具体错误信息"
}
```

常见错误：
- 未认证：401 Unauthorized
- 参数错误：400 Bad Request
- 权限不足：403 Forbidden
- 资源不存在：404 Not Found
- 服务器错误：500 Internal Server Error

### 注意事项
1. 所有需要认证的API必须在请求头中包含有效的JWT token
2. 文件上传的容量限制为5MB
3. 所有时间字段均使用ISO 8601格式
4. API的版本信息包含在URL路径中
5. 分页相关的参数默认值：page=1, size=10