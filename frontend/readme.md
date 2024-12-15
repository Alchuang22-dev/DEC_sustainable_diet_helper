# DEC前端文档

> 最新修改: 2014/12/12 (from Zeyu Zhang)

Using Vue3 in uni-app

![](C:\Users\25848\Desktop\软工原型\login.png)

## 文件结构介绍

```c++
│  App.vue
│  categories.txt
│  index.html
│  main.js //主入口
│  manifest.json
│  pages.json //页面路由
│  project.config.json //系统设置
│  project.private.config.json
│  uni.scss
│  
├─.hbuilderx
│      launch.json //启动器
│      
├─locales //内置语言包
│      en.json //英语包
│      index.ts //包管理文件
│      readme.md
│      zh-Hans.json //简体中文包
│      
├─pages //一级页面
│  ├─index
│  │      index.vue //主页
│  │      
│  ├─my_index
│  │      my_index.vue //用户页
│  │      
│  ├─news_index
│  │      news_index.vue //社区页
│  │      
│  ├─static //用户页svg文件
│  │      close.svg
│  │      editor.svg
│  │      family.svg
│  │      favorites.svg
│  │      food.svg
│  │      historicaldata.svg
│  │      logout.svg
│  │      mywork.svg
│  │      search.svg
│  │      setgoals.svg
│  │      setting.svg
│  │      user.svg
│  │      
│  └─tools_index //功能页
│          tools_index.vue
│          
├─pagesMy //用户页二级页面
│  ├─appSettings 
│  │      appSettings.vue //软件设置
│  │      
│  ├─favorites
│  │      favorites.vue //用户收藏
│  │      
│  ├─foodPreferences
│  │      foodPreferences.json
│  │      foodPreferences.vue //用户饮食偏好
│  │      
│  ├─historyData
│  │      historyData.vue //历史数据
│  │      
│  ├─login
│  │      login.vue //登录页
│  │      
│  ├─myFamily
│  │      myFamily.vue //用户家庭管理
│  │      
│  ├─my_home
│  │      my_home.vue //用户创作主页
│  │      
│  ├─searchTools
│  │      searchTools.vue //搜索工具
│  │      
│  ├─setGoals
│  │      setGoals.vue //设置目标
│  │      
│  ├─static
│  │      delete.svg
│  │      edit.svg
│  │      view.svg
│  │      
│  ├─userSettings
│  │      userSettings.vue //用户设置
│  │      
│  └─wechatLogin
│      │  wechatLogin.vue //微信登录专门页
│      │  
│      └─static
│              avatar.png
│              
├─pagesNews //社区二级页面
│  ├─create_news //创作编辑器（未完成）
│  │      create_news.vue
│  │      
│  ├─news_detail
│  │      news_detail.vue //图文详情 （未完成）
│  │      
│  ├─static
│  │      addpicture.svg
│  │      addtext.svg
│  │      minus.svg
│  │      plus.svg
│  │      save.svg
│  │      share.svg
│  │      
│  ├─video_detail
│  │      video_detail.vue //视频播放器（暂停开发）
│  │      
│  └─web_detail
│          web_detail.vue //网页适配器 （暂停开发）
│          
├─pagesSetting
│  ├─Authorizations
│  │      Authorizations.vue //用户协议
│  │      
│  ├─Bend
│  │      Bend.vue //账号绑定（此版本不可用）
│  │      
│  ├─ConnectUs
│  │      ConnectUs.vue //开发团队信息
│  │      
│  ├─DeleteData
│  │      DeleteData.vue //删除数据 （未完成）
│  │      
│  ├─DeleteId
│  │      DeleteId.vue //注销账号（未完成）
│  │      
│  ├─DevelopingFunc
│  │      DevelopingFunc.vue //开发中功能
│  │      
│  ├─FeedBack
│  │      FeedBack.vue//反馈
│  │      
│  ├─Fonts
│  │      Fonts.vue//字体调节（此版本不可用，请使用微信适老化工具）
│  │      
│  ├─Idinfo
│  │      Idinfo.vue//用户信息
│  │      
│  ├─InfoShared
│  │      InfoShared.vue//用户权限
│  │      
│  ├─language
│  │      language.vue//多语言设置
│  │      
│  ├─newsSetting
│  │      newsSetting.vue//消息设置（未完成）
│  │      
│  ├─Permissions
│  │      Permissions.vue//个人信息收集清单
│  │      
│  ├─PIConnected
│  │      PIConnected.vue//个人信息共享清单
│  │      
│  ├─recoSetting
│  │      recoSetting.vue//推荐设置（未完成）
│  │      
│  ├─SoftwareInfo
│  │      SoftwareInfo.vue//软件信息
│  │      
│  └─Storage
│          Storage.vue//存储
│          
├─pagesTool //功能二级页面
│  ├─add_food
│  │      add_food.vue //添加食材
│  │      
│  ├─carbon_calculator
│  │      carbon_calculator.vue //碳计算器
│  │      
│  ├─food_recommend
│  │      food_recommend.vue //食谱推荐
│  │      
│  ├─home_servant
│  │      home_servant.vue //家庭管理
│  │      
│  ├─modify_food
│  │      modify_food.vue //编辑食材
│  │      
│  ├─nutrition_calculator
│  │      nutrition_calculator.vue //营养日历
│  │      
│  ├─recipe
│  │      recipe.vue //食谱详情
│  │      
│  └─static
│          background_img.jpg
│          recipe_default.png
│          toufu.png
│          
├─static
│  │  c1.png
│  │  c2.png
│  │  c3.png
│  │  c4.png
│  │  c5.png
│  │  c6.png
│  │  c7.png
│  │  c8.png
│  │  c9.png
│  │  customicons.css
│  │  customicons.ttf
│  │  loading.gif
│  │  logo.png
│  │  uni.png
│  │  
│  └─images
│      ├─index
│      │      background_img.jpg
│      │      background_index_new.png
│      │      logo_wide.png
│      │      weidenglu.png
│      │      
│      └─settings
│              logo.png
│              mail.png
│              mail_open.png
│              trash.png
│              
├─stores //pinia库
│      carbon_and_nutrition_data.js
│      draft.js
│      family.js
│      food_list.js
│      news_list.js
│      user.js
│      
├─[releasepackages]
```

## 开发计划

近期：

+ 营养日历前后端对接
+ 主页数据图表前后端对接
+ 家庭管理完善
+ 食材图片加入数据库
+ 菜谱描述加入数据库

下一开发周期：

+ 社区和草稿功能前后端连接
+ 添加账号注销、数据删除支持
+ 添加历史数据支持

后期计划：

+ 大模型开发功能引入
+ 正式新闻功能
