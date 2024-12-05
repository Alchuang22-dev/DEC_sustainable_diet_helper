
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

# 可持续饮食助手 API 文档

## 目录
- [通用说明](#通用说明)
- [用户相关](#用户相关)
- [新闻相关](#新闻相关)
- [家庭相关](#家庭相关)

## 通用说明

### Base URL
```
https://122.51.231.155/
```

### 认证方式
大多数API端点需要JWT认证。在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

### 响应格式
所有响应均为JSON格式，包含以下基本结构：
- 成功响应：HTTP 2xx
```json
{
    "message": "操作成功信息",
    "data": {
        // 具体数据
    }
}
```
- 错误响应：HTTP 4xx/5xx
```json
{
    "error": "错误信息"
}
```

## 用户相关

### 微信登录/注册
**POST** `/auth/wechat`

用微信Code进行登录或注册。

**请求体**：
```json
{
    "code": "string",     // 必需，微信登录凭证
    "nickname": "string"  // 可选，用户昵称
}
```

**响应**：
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

### 设置用户昵称
**PUT** `/user/nickname`

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

**响应**：
```json
{
    "message": "Nickname updated successfully",
    "nickname": "新昵称"
}
```

### 设置用户头像
**PUT** `/user/avatar`

**请求头**：
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求参数**：
- `avatar`: 文件(jpg格式)

**响应**：
```json
{
    "message": "Avatar updated successfully",
    "avatar_url": "头像路径"
}
```

## 新闻相关

### 上传新闻
**POST** `/news`

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

**响应**：
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

### 获取新闻详情
**GET** `/news/{id}`

**请求头**：
```
Authorization: Bearer <token>
```

**响应**：
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

### 获取新闻列表
**GET** `/news`

**请求头**：
```
Authorization: Bearer <token>
```

**查询参数**：
- `category`: string (可选，新闻类别)
- `page`: number (可选，默认1)
- `size`: number (可选，默认10)
- `sort_by`: string (可选，排序字段：upload_time/view_count/like_count)
- `order`: string (可选，排序方式：asc/desc)

**响应**：
```json
{
    "total": 100,
    "page": 1,
    "size": 10,
    "data": [
        {
            "id": 1,
            "title": "标题",
            "description": "描述",
            "upload_time": "时间",
            "view_count": 0,
            "share_count": 0,
            "like_count": 0,
            "favorite_count": 0,
            "dislike_count": 0,
            "news_type": "类型",
            "author_name": "作者",
            "author_avatar": "头像",
            "tags": ["标签"],
            "extra_field": "预览信息"
        }
    ]
}
```

### 新闻互动API

#### 点赞新闻
**POST** `/news/{id}/like`

#### 取消点赞
**DELETE** `/news/{id}/like`

#### 收藏新闻
**POST** `/news/{id}/favorite`

#### 取消收藏
**DELETE** `/news/{id}/favorite`

#### 点踩新闻
**POST** `/news/{id}/dislike`

#### 取消点踩
**DELETE** `/news/{id}/dislike`

所有互动API都需要以下请求头：
```
Authorization: Bearer <token>
```

响应格式：
```json
{
    "message": "操作成功消息",
    "xxx_count": 新数量  // like_count/favorite_count/dislike_count
}
```

### 添加评论
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

**响应**：
```json
{
    "message": "Comment added successfully",
    "comment": {
        "id": 1,
        "content": "内容",
        "user_id": 1,
        "publish_time": "时间",
        "like_count": 0
    }
}
```

## 家庭相关

### 创建家庭
**POST** `/family`

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "name": "string"  // 必需，家庭名称
}
```

**响应**：
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

### 查看家庭详情
**GET** `/family`

**请求头**：
```
Authorization: Bearer <token>
```

**响应**：
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

### 搜索家庭
**GET** `/family/search`

**请求头**：
```
Authorization: Bearer <token>
```

**查询参数**：
- `family_id`: string (必需，家庭唯一标识)

**响应**：
```json
{
    "id": 1,
    "name": "名称",
    "family_id": "标识",
    "member_count": 1
}
```

### 申请加入家庭
**POST** `/family/{id}/join`

**请求头**：
```
Authorization: Bearer <token>
```

**响应**：
```json
{
    "message": "Join request sent successfully"
}
```

### 管理员审批加入申请
**POST** `/family/admit`

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "user_id": 1  // 必需，申请用户的ID
}
```

**响应**：
```json
{
    "message": "User successfully admitted to the family",
    "family_id": 1,
    "user_id": 1
}
```

### 拒绝加入申请
**POST** `/family/reject`

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "user_id": 1  // 必需，申请用户的ID
}
```

**响应**：
```json
{
    "message": "User's join request rejected successfully",
    "family_id": 1,
    "user_id": 1
}
```

### 取消加入申请
**DELETE** `/family/join`

**请求头**：
```
Authorization: Bearer <token>
```

**响应**：
```json
{
    "message": "Join request canceled successfully",
    "family_id": 1,
    "user_id": 1
}
```

### 管理员权限设置
#### 设置为普通成员
**POST** `/family/set-member`

#### 设置为管理员
**POST** `/family/set-admin`

两个API都需要以下：

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：
```json
{
    "user_id": 1  // 必需，目标用户ID
}
```

**响应**：
```json
{
    "message": "Successfully set user to member/admin",
    "family_id": 1,
    "user_id": 1
}
```
