## 题目

发现世界的美

## 需求分析与开发计划

写一个类小红书的app

### 登录注册模块

- [ ] 用户登录，需要使用短信，邮箱验证，未注册用户当场自动注册
- [ ] 选择性别，年龄，可见
- [ ] 选择至少3个感兴趣的内容（设置用户tag初始权重）
- [ ] 手机，邮箱登录
### 个人信息模块

- [ ] 可以查看个人和他人的信息
- [ ] 可编辑信息：头像和昵称，简介，性别，地区
- [ ] 发过的贴子
- [ ] 收藏的帖子 
- [ ] 赞过的帖子
### 首页模块

- [ ] 瀑布流展示热门帖子
- [ ] 按tag对帖子分类，用户可以新增tag，系统也有默认tag
- [ ] 利用推荐算法，推荐感兴趣tag相关的帖子
- [ ] 可以推荐附近的人发的贴子
### 帖子模块

- [ ] 用户点的“+”号可以创建帖子，包含图片和文字，富媒体（语音和视频etc）
- [ ] 富文本框显示内容，显示用户发帖的IP属地
- [ ] 用户点赞
- [ ] 嵌套评论，显示评论用户的IP属地
### 扩展：AIGC

- [ ] 接入openai 可以利用chat3.5生成生成内容（文章和图片）
- [ ] ai生成预测tag
## API

### 用户鉴权

#### 登录注册

- [x] POST /session
**请求样例**

```json
{
  "mode":"<email | sms>",
  // option1 email auth
  "auth": {
     "email": "i@zhangzqs.cn",
     "code": "<code>"
  },
  // option2 sms auth
  "auth": {
     "sms": "1234567890",
     "code": "<code>"
  }
}
```
**响应样例**
```json
200 OK 
{
  "token": "xxxxxxxx.xxxxxx.xxxx"
}
```
jwt过期时间设置为7天过期
#### 发送验证码

- [ ] POST /session/authcode
**请求样例**

```json
{
  "mode":"<email | sms>", //
  // option1 email auth
  "auth": {
     "email": "i@zhangzqs.cn"
  },
  // option2 sms auth
  "auth": {
     "sms": "1234567890"
  }
}
```
**响应样例**
```json
200 OK
{
  "message": "<xxx>"
}
```
### 个人信息

#### 获取个人信息

- [x] GET /users/{uid}/info
**请求样例**

uid: int

**响应样例**

```json
{
  "uid": <int>,
  "email": "" | null,
  "sms": "" | null,
  "nickname": "",
  "avatar": "",
  "introduction": ""  
}
```
#### 更新个人信息

- [x] POST /users/{uid}/info
```json
{
  "email": "",
  "sms": "",
  "nickname": "",
  "avatar": "",
  "introduction": ""  
}
```
注意未填写的字段表示不更新的字段，不要全部更新数据库全部字段为null
#### 获取足迹

- [ ] GET  /users/{uid}/history/{type}?pageNum={}&pageSize={}
获取发过的，赞过的，收藏过的帖子历史记录

**请求样例**

type: <my | like | favorite>

pageNum: 页号

pageSize: 页数

**响应样例**

```json
[
  {
    "pid": <pid>,
    "title": "<title>",
    "cover": "<image url>",
    "nickname": "<nickname>",
    "avatar": "<avatar url>",
  },
  {
    ...
  }
]
```
### 标签模块

#### 创建标签

- [ ] POST /tags
```json
{
  "tid": 1, 
  "name": "<name>"
}
```
#### 检索标签

- [ ] GET /tags?q={}&pageNum={}&pageSize={}
**请求示例**

```json
Empty
```
**响应示例**
```json
[
  {"tid": 1, name: "<name>"},
  ...
]
```
#### 获取所有基本标签

- [ ] GET /basic_tags
**请求示例**

**响应示例**

```json
[
  {
    "name": "<name>", // 名称
    "cover": "<cover url>", // 封皮图片
    "tid": <tag id>
  }
]
```
#### 更新用户初始感兴趣的标签

- [ ] POST /users/{uid}/my_tags
**请求示例**

```json
[<tagId1>, <tagId2>, ...]
```
### **响应示例**

### 首页&帖子模块

#### 列举帖子

- [x] GET /posts?tag=<tagId>
- [ ] 待测试
列举出首页推荐的帖子/选择相关tag下的帖子

**请求样例**

**响应样例**

```json
[
  {
    "pid": <pid>,
    "title": "<title>",
    "cover": "<image url>",
    "nickname": "<nickname>",
    "avatar": "<avatar url>",
    "isLike": false,
  },
  {
    ...
  }
]
```

#### 搜索帖子

- [x] GET /posts/search?q={}&pageNum={}&pageSize={}
- [ ] 待测试
**请求样例**

q = urlencode后的查询文本

**响应样例**

```json
[
  {
    "pid": <pid>,
    "title": "<title>",
    "cover": "<image url>",
    "nickname": "<nickname>",
    "avatar": "<avatar url>",
    "isLike": false,
    "LikeCount": 0,
     "FavoriteCount":0
  },
  {
    ...
  }
]
```

#### 附近的帖子

- [ ] GET /posts/local 附近的需要单独出来接口
**请求样例**

**响应样例**

```json
[
  {
    "pid": <pid>,
    "title": "<title>",
    "cover": "<image url>",
    "nickname": "<nickname>",
    "avatar": "<avatar url>",
    "isLike": false,
    "LikeCount": 0,
     "FavoriteCount":0
  },
  {
    ...
  }
]
```

#### **获取帖子内容**

- [ ] GET /posts/{pid} 获取具体一个帖子内容（包含正文富文本，点赞数）
**请求样例**

**响应样例**

```json
{
  "title": "<title>",
  "cover": "<image url>",
  "nickname": "<nickname>",
  "avatar": "<avatar url>",
  "content": "<html content>",
  "isLike": false,
  "LikeCount": 0,
  "FavoriteCount":0
 
}
```
#### 获取评论

- [ ] GET /posts/{pid}/comments?pageNum={}&pageSize={}
**请求样例**

pageNum: {min: 0, max: ?}

pageSize: {min: 10, max: 20}

**响应样例**

```json
[
  {
    "source": {
      "uid": <uid>,
      "nickname": <nickname>,
      "avatar": <avatar>
    },
    "target": {
      "uid": <uid>,
      "nickname": <nickname>,
      "avatar": <avatar>
    },
    "content": "content"
  },
  ...
]
```
#### 创建帖子

- [x] POST /posts 
- [ ] 待测试
**请求样例**

```json
{
  "title": "xxxx",
  "content": "<html content>",
}
```
**响应样例**
```json
200 OK
```
#### 编辑帖子

- [x] PUT /posts/{pid}
- [ ] 待测试
**请求样例**

```json
{
  "title": "xxxx",
  "cover": "<conver url>",
  "content": "<html content>",
}
```
**响应样例**
```json
200 OK
```
#### 删除帖子

- [x] DELETE /posts/{pid}
- [ ] 待测试
**请求样例**

**响应样例**

```json
200 OK
```
#### 点赞/收藏帖子

- [ ] POST /posts/{pid}/bookmark/{like, favorite} 
三个参数uid，pid，type；谁给什么帖子进行什么操作

uid为jwt中字段

**请求样例**

```json
Empty
```
**响应样例**
```json
200 OK
```
#### 取消点赞/收藏帖子

- [ ] DELETE /posts/{pid}/bookmark/{like, favorite} 
三个参数uid，pid，type；谁给什么帖子进行什么操作

uid为jwt中字段

**请求样例**

```json
Empty
```
**响应样例**
```json
200 OK
```

### 
### 附件模块

#### 上传附件

- [ ] POST /attachment
**请求样例**

```json
Content-Type: multipart/form-data

file: <binary>
```
**响应样例**
```json
{
  "url": "http://xxx.com/attachment/{uuid}.xxx"
}
```
#### 下载附件

- [ ] GET /attachment/{uuid}.xxx
**请求样例**

**响应样例**

```json
Content-Type: <contentType>
<binary-stream>
```
### 扩展：AIGC

## 数据库设计

### users表

存放了所有用户基本信息

1. uid : int (主键，自增)
2. email : text
3. sms : text
4. nickname : text
5. avatar : text
6. introduction : text
7. created_at: datetime
8. deleted_at: datetime
9. update_at: datetime
### posts表

存放了所有的文章信息

1. pid : INT (主键，自增)
2. uid : INT (外键, 对应users表中的uid)
3. title : text
4. content : text
5. like_count : int
6. favorite_count : int
7. create_at: datetime
8. update_at: datetime
9. deleted_at: datetime
### tags表

存放了所有的标签tags

1. tid : INT (主键，自增)
2. name : text
### post_tags表

存放了某个文章标注的所有tags的关系表，多对多关系必须创建新表

1. pid : int (外键, 对应posts表中的pid)
2. tid : int (外键, 对应tags表中的tid)
### bookmarks表

1. uid : INT (主键，外键, 对应users表中的uid)
2. pid : INT (主键，外键, 对应posts表中的pid)
3. type : ENUM('LIKE', 'FAVORITE') (主键)
4. create_at: datetime
5. update_at: datetime
6. deleted_at: datetime