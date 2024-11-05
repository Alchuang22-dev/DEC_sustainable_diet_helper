


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