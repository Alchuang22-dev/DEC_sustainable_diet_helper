<template> 
  <!-- 用户信息区 -->
  <view class="profile-header">
    <view class="profile-info">
      <image :src="avatarSrc" class="avatar" alt="用户头像"></image>
      <text class="nickname">{{ uid }}</text>
      <text class="userid">uid：{{ user_id || 'test_user' }}</text>

      <!-- 创作统计：使用 t('creation') -->
      <view class="stats">
        <view class="stats-item">
          <text>{{ t('creation') }}</text>
          <text>{{ publishedCount }}</text>
        </view>
        <view class="stats-item">
          <!-- 使用 t('draft') -->
          <text>{{ t('draft') }}</text>
          <text>{{ draftCount }}</text>
        </view>
      </view>
    </view>
  </view>

  <!-- 分割线 -->
  <view class="separator"></view>

  <view class="container">
    <!-- 已发布列表 -->
    <view class="card-list">
      <!-- 如果 publishedArticles 为空，显示占位符，否则显示卡片 -->
      <view v-if="publishedArticles.length === 0" class="empty-placeholder">
        {{ t('emptyArticles') }}
      </view>
      <view
        v-else
        v-for="(item, index) in publishedArticles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
		@click = "switchtoPost(index)"
      >
        <view class="card-header">
          <view class="title">{{ item.title }}</view>
          <!-- 使用 t('pieces') 代替“作品” -->
          <view class="status">{{ t('pieces') }}</view>
        </view>

        <view class="card-body">
          <view class="description">{{ item.description }}</view>
          <view class="info">
            <text class="publish-time">{{ formatPublishTime(item.publishTime) }}</text>
            <view class="stats">
              <text class="like-count">👍 {{ item.likes }}</text>
              <text class="favorite-count">⭐ {{ item.favorites }}</text>
              <text class="share-count">🔗 {{ item.shares }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>



<script setup>
import { onMounted, ref, reactive, computed, watch } from 'vue';
import { useNewsStore } from '@/stores/news_list';
import { useI18n } from 'vue-i18n';
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app';
import { storeToRefs } from 'pinia';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储
import { useRouter } from 'vue-router'; // 如果使用 vue-router

const newsStore = useNewsStore();
const userStore = useUserStore(); // 使用用户存储

const router = useRouter();

// 用来获取本地时间和日期
const systemDate = new Date();
const systemDateStr = systemDate.toISOString().slice(0, 10); // 获取当前系统日期，格式：YYYY-MM-DD
const BASE_URL = ref('http://122.51.231.155:8080');

const activeIndex = ref(null);
// 计算属性从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);
const uid = computed(() => userStore.user.nickName);
const user_id = computed(() => userStore.user.uid);

const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? userStore.user.avatarUrl // 假设 avatar_url 是完整的 URL
    : '/static/images/index/default_avatar.jpg'
);

// 从 Store 获取数据和方法
const { filteredNewsItems, selectedSection, isRefreshing } = storeToRefs(newsStore);
const { setSection, refreshNews, fetchNews } = newsStore;

// 模拟数据
const articles = ref([]);
const { t, locale, messages } = useI18n();
const jwtToken = computed(() => userStore.user.token); // 用户的 JWT Token

// 转换时间函数
const formattedSaveTime = computed((time) => {
  const postDate = time.slice(0, 10); // 提取日期部分

  if (postDate === systemDateStr) {
    // 如果日期相同，显示 "今天 HH:mm:ss"
    const postTime = new Date(time); // 转换为 Date 对象
    const hours = postTime.getHours().toString().padStart(2, '0');
    const minutes = postTime.getMinutes().toString().padStart(2, '0');
    const seconds = postTime.getSeconds().toString().padStart(2, '0');
    return `今天 ${hours}:${minutes}:${seconds}`;
  } else {
    // 否则显示完整日期
    return postDate;
  }
});

/**
 * 格式化发布时间
 * @param {string} publishTime - ISO 格式或其他可被 Date 解析的字符串
 * @returns {string} - 格式化后的时间字符串
 */
const formatPublishTime = (publishTime) => {
  const date = new Date(publishTime);
  const now = new Date();

  // 判断是否是同一天
  const isSameDay =
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate();

  if (isSameDay) {
    // 如果是同一天，显示“今天 HH:mm”
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `今天 ${hours}:${minutes}`;
  } else {
    // 否则显示 YYYY-MM-DD 或者你想要的其他格式
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
};

// 用户信息部分的逻辑
const backgroundImageUrl = ref('/static/images/index/background_img.jpg');

const editBackgroundImage = () => {
  uni.chooseImage({
    count: 1,
    success: (res) => {
      // 假设取第一张图片并上传
      const tempFilePath = res.tempFilePaths[0];
      // ...上传逻辑省略
      // 上传成功后更新
      backgroundImageUrl.value = tempFilePath;
    }
  });
};

const switchtoPost = (index) => {
	const article = articles.value[index];
		  uni.navigateTo({
		  	url: `/pagesNews/news_detail/news_detail?id=${article.id}`,
		  });
}

// 创作统计数量
const publishedCount = computed(() =>
  articles.value.filter(a => a.status === '已发布').length
);
const draftCount = computed(() =>
  articles.value.filter(a => a.status === '草稿').length
);
const favoriteCount = ref(0);  // 你可以根据后端返回的数据赋值
const followerCount = ref(0);  // 同上

// 当前标签（如果需要切换）
const currentTab = ref('published');

// 根据 status 筛选
const publishedArticles = computed(() => {
  return articles.value.filter(a => a.status === '已发布');
});
const draftArticles = computed(() => {
  return articles.value.filter(a => a.status === '草稿');
});

// Function to get news or draft details
const getArticleDetails = async (id, isDraft = false) => {
  const url = isDraft
    ? `${BASE_URL.value}/news/details/draft/${id}`
    : `${BASE_URL.value}/news/details/news/${id}`;
    const res = await uni.request({
      url: url,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });
    if (res.statusCode === 200) {
      console.log('获取详细信息', res.data);
      return res.data;
    } else {
      console.error(`获取文章详情失败: ${res.statusCode}`);
      return null;
    }
};

// Function to fetch user profile and articles
const fetchData = async () => {
  // 获取用户 ID 作为路径参数
  const userId = getUserIdFromRoute(); // 实现该函数
  console.log("寻找用户：",userId);

  if (!userId) {
    uni.showToast({
      title: '用户ID未找到',
      icon: 'none',
      duration: 2000
    });
    return;
  }
    // 请求用户个人主页数据
    const res = await uni.request({
      url: `${BASE_URL.value}/users/${userId}/profile`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });

    if (res.statusCode === 200) {
      const data = res.data;
      console.log('用户个人主页数据:', data);

      // 更新用户信息
      userStore.user.nickName = data.nickname;
      userStore.user.avatarUrl = `${BASE_URL.value}/static/${data.avatar_url}`;

      // 获取新闻详细信息
      const newsDetailsPromises = data.news.map(newsItem => getArticleDetails(newsItem.id));
      const newsDetails = await Promise.all(newsDetailsPromises);

      // 过滤掉获取失败的新闻
      const validNewsDetails = newsDetails.filter(detail => detail !== null);

      // 更新 articles
      articles.value = validNewsDetails.map(detail => ({
        ...detail,
        publishTime: detail.upload_time || detail.updated_at, // 根据实际字段调整
        likes: detail.like_count || 0,
        favorites: detail.favorite_count || 0,
        shares: detail.share_count || 0,
        status: detail.status || '已发布', // 根据需要设置 status
        bgColor: detail.status === '已发布' ? 'rgba(0, 123, 255, 0.1)' : 'rgba(255, 193, 7, 0.1)'
      }));


    } else if (res.statusCode === 401) {
      uni.showToast({
        title: '未授权，请重新登录',
        icon: 'none',
        duration: 2000
      });
      // 可能需要跳转到登录页
    } else if (res.statusCode === 404) {
      uni.showToast({
        title: '用户未找到',
        icon: 'none',
        duration: 2000
      });
    } else {
      uni.showToast({
        title: '获取用户数据失败',
        icon: 'none',
        duration: 2000
      });
    }
};

// 获取路由参数中的用户 ID
const getUserIdFromRoute = () => {
  // 如果使用 vue-router
  // return router.currentRoute.value.params.id;

  // 如果使用 uni-app 的页面参数
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  return currentPage.options.id;
};

// Lifecycle hook to load articles
onShow(() => {
  fetchData();
});
</script>

<style scoped>
/* Body */
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

.profile-header {
  margin-top: 20px;
  position: relative;
  width: 100%;
  height: 220px; /* 适当加大 */
  background-color: #f5f5f5; /* 如果没有背景图时的底色 */
  overflow: hidden;
}

/* 背景图可编辑：点击后替换 */
.profile-bg {
  width: 100%;
  height: 160px;
  object-fit: cover;
  z-index: -1;
}

/* 编辑背景按钮，如果你想单独做一个icon，也可绝对定位到右下角 */
.edit-bg-btn {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background-color: rgba(255,255,255,0.5);
  border: none;
  border-radius: 4px;
  padding: 6px 10px;
}

.profile-info {
  position: absolute;
  left: 20px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2px solid #fff;
  object-fit: cover;
  margin-bottom: 5px;
}

.nickname {
  font-weight: bold;
  font-size: 16px;
  margin-bottom: 2px;
  color: #333;
}

.userid {
  font-size: 12px;
  margin-top: 4px;
  color: #666;
  z-index: 10;
}

.separator {
  margin: 10px 0;
  width: 100%;
  height: 1px;
  background-color: #e0e0e0;
}

/* 创作、草稿统计 */
.stats {
  margin-top: 8px;
  display: flex;
  gap: 20px; /* 控制间距 */
  /* 移除 height */
}

.stats-item text:nth-child(1) {
  font-size: 12px;
  color: black;
}
.stats-item text:nth-child(2) {
  font-size: 14px;
  font-weight: bold;
  margin-left: 4px;
}

/* 修改资料按钮 */
.edit-profile-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: #f5f5f5;
  border: 1px solid #ccc;
  border-radius: 5px;
  padding: 4px 8px;
}

/* 下方卡片列表容器 */
.container {
  padding: 20px;
  margin-top: 20px; /* 给 profile-header 留出空间 */
}

/* 卡片列表 */
.card-list {
  display: flex;
  flex-direction: column;
}

/* 卡片样式 */
.card {
  margin-bottom: 10px;  /* 卡片之间的间距，改小比如 10px */
  border-radius: 8px;
  padding: 10px;        /* 卡片内容与边框的内边距，可改小比如 10px */
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px; /* 从10px改为5px */
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.status {
  font-size: 14px;
  color: #007bff;
}

.card-body {
  margin-bottom: 10px; /* 从15px改为10px */
  /* 移除固定高度 */
}

.description {
  font-size: 14px;
  color: #555;
  margin-bottom: 10px;
}

.info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #777;
}

.info .stats {
  display: flex;
  gap: 5px; /* 从10px改为5px */
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 5px; /* 从10px改为5px */
}

.action-btn {
  background-color: transparent;
  border: none;
  cursor: pointer;
  padding: 5px;
}

.icon {
  width: 24px;
  height: 24px;
  transition: transform 0.2s ease;
}

.icon:hover {
  transform: scale(1.2); /* 鼠标悬浮时放大图标 */
}

.publish-time {
  font-size: 12px;
  color: #777;
}

.like-count, .favorite-count, .share-count {
  font-size: 12px;
  color: #777;
}
</style>

