<template> 
   
  <!-- 可编辑背景图 -->
  <view class="profile-header">
	<!--
	<image
      :src="backgroundImageUrl"
      class="profile-bg"
      @click="editBackgroundImage"
    ></image>
	-->
    <!-- 用户信息区 -->
    <view class="profile-info">
      <image :src="avatarSrc" class="avatar" alt="用户头像"></image>
      <text class="nickname">{{ uid }}</text>
      <text class="userid">uid：{{ user_id || 'test_user' }}</text>

      <!-- 创作统计 -->
      <view class="stats">
        <view class="stats-item">
          <text>创作</text>
          <text>{{ publishedCount }}</text>
        </view>
        <view class="stats-item">
          <text>草稿</text>
          <text>{{ draftCount }}</text>
        </view>
      </view>
    </view>
    
    <!-- 修改资料的按钮 -->
  </view>

  <!-- 标签（已发布/草稿）切换 -->
  <view class="tabs">
    <view
      :class="['tab', currentTab === 'published' ? 'active' : '']"
      @click="currentTab = 'published'"
    >
      作品栏
    </view>
    <view
      :class="['tab', currentTab === 'draft' ? 'active' : '']"
      @click="currentTab = 'draft'"
    >
      草稿箱
    </view>
  </view>

  <view class="container">
    <!-- 已发布列表：仅在 currentTab === 'published' 时显示 -->
    <view
      v-if="currentTab === 'published'"
      class="card-list"
    >
      <view
        v-for="(item, index) in publishedArticles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
      >
        <!-- 与原本的卡片展示基本一致 -->
        <view class="card-header">
          <view class="title">{{ item.title }}</view>
          <view class="status">作品</view>
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
        <view class="card-footer">
          <!-- 这里使用原先的按钮: 查看、编辑、删除等 -->
          <button @click="viewArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/view.svg" class="icon" alt="View" ></image>
          </button>
          <button @click="editArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/edit.svg" class="icon" alt="Edit" ></image>
          </button>
          <button @click="deleteArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/delete.svg" class="icon" alt="Delete" ></image>
          </button>
        </view>
      </view>
    </view>

    <!-- 草稿列表：仅在 currentTab === 'draft' 时显示 -->
    <view
      v-else
      class="card-list"
    >
      <view
        v-for="(item, index) in draftArticles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
      >
        <view class="card-header">
          <view class="title">{{ item.title }}</view>
          <view class="status">草稿</view>
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
        <view class="card-footer">
          <button @click="viewDraft(index)" class="action-btn">
            <image src="@/pagesMy/static/view.svg" class="icon" alt="View" ></image>
          </button>
          <button @click="editDraft(index)" class="action-btn">
            <image src="@/pagesMy/static/edit.svg" class="icon" alt="Edit" ></image>
          </button>
          <button @click="deleteDraft(index)" class="action-btn">
            <image src="@/pagesMy/static/delete.svg" class="icon" alt="Delete" ></image>
          </button>
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

const newsStore = useNewsStore();
const userStore = useUserStore(); // 使用用户存储

// 用来获取本地时间和日期
const systemDate = new Date();
const systemDateStr = systemDate.toISOString().slice(0, 10); // 获取当前系统日期，格式：YYYY-MM-DD
const BASE_URL = ref('http://122.51.231.155:8080');
const user_id = computed(() => userStore.user.uid);

const activeIndex = ref(null);
// 计算属性从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);
const uid = computed(() => userStore.user.nickName);

const avatarSrc = computed(() =>
    userStore.user.avatarUrl
        ? `${BASE_URL.value}/static/${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);

// 从 Store 获取数据和方法
const { filteredNewsItems, selectedSection, isRefreshing } = storeToRefs(newsStore);
const { setSection, refreshNews, fetchNews } = newsStore;

// 模拟数据
const articles = ref([]);
const { t, locale, messages } = useI18n();
const jwtToken = computed(() => userStore.user.token);; // Replace with actual token

//转换时间
const formattedSaveTime = computed((time) => {
  const postDate = time.slice(0, 10); // 提取日期部分

  if (postDate === systemDateStr) {
    // 如果日期相同，显示 "today" + 时间
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

//用户信息部分的逻辑
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

// 创作统计数量
const publishedCount = computed(() =>
  articles.value.filter(a => a.status === '已发布').length
);
const draftCount = computed(() =>
  articles.value.filter(a => a.status === '草稿').length
);
const favoriteCount = ref(0);  // 你可以根据后端返回的数据赋值
const followerCount = ref(0);  // 同上

// 当前标签
const currentTab = ref('published');

// 根据 status 筛选
const publishedArticles = computed(() => {
  return articles.value.filter(a => a.status === '已发布');
});
const draftArticles = computed(() => {
  return articles.value.filter(a => a.status === '草稿');
});

// Function to get published news IDs
const getPublishedNewsIds = async () => {
  console.log('获取已发布');
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/my_news`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });
    return res.data.news_ids || [];
  } catch (error) {
    console.error('Error fetching published news IDs', error);
    return [];
  }
};

// Function to get draft news IDs
const getDraftNewsIds = async () => {
	console.log('获取草稿');
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/my_drafts`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });
    return res.data.draft_ids || [];
  } catch (error) {
    console.error('Error fetching draft news IDs', error);
    return [];
  }
};

// Function to get news or draft details
const getArticleDetails = async (id, isDraft = false) => {
  const url = isDraft
    ? `${BASE_URL.value}/news/details/draft/${id}`
    : `${BASE_URL.value}/news/details/news/${id}`;
  try {
    const res = await uni.request({
      url: url,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });
	console.log('获取详细信息');
	console.log(res.data);
    return res.data;
  } catch (error) {
    console.error('Error fetching article details', error);
    return null;
  }
};

// Function to fetch articles on page load
const fetchArticles = async () => {
  const publishedIds = await getPublishedNewsIds();
  const draftIds = await getDraftNewsIds();

  const allArticles = [];
  for (const id of publishedIds) {
    const details = await getArticleDetails(id);
    if (details) {
      allArticles.push({
        ...details,
		publishTime: details.upload_time,
		likes: details.like_count,
		favorites: details.favorite_count,
		shares: details.share_count,
        status: '已发布',
        bgColor: 'rgba(0, 123, 255, 0.1)', // Published color
      });
    }
  }

  for (const id of draftIds) {
    const details = await getArticleDetails(id, true);
    if (details) {
      allArticles.push({
        ...details,
		publishTime: details.updated_at,
        status: '草稿',
        bgColor: 'rgba(255, 193, 7, 0.1)', // Draft color
      });
    }
  }

  articles.value = allArticles;
};

// Lifecycle hook to load articles
onShow(() => {
  fetchArticles();
});

// View article function 仅限草稿
const viewArticle = (index) => {
  const article = articles.value[index];
	  uni.navigateTo({
	  	url: `/pagesNews/news_detail/news_detail?id=${article.id}`,
	  });
};

const viewDraft = (index) => {
	const draft = draftArticles.value[index];
	uni.navigateTo({
		url: `/pagesNews/preview_draft/preview_draft?id=${draft.id}`,
	});
};

// Edit article function
const editArticle = (index) => {
  const article = articles.value[index];
  	  uni.showToast({
  	    title: '发布后不可编辑',
  	    icon: 'none',
  	    duration: 2000
  	  });
};

const editDraft = (index) => {
	const draft = draftArticles.value[index];
	uni.navigateTo({
	  url: `/pagesNews/edit_draft/edit_draft?id=${draft.id}`,
	});
};

// Delete article function
// Delete article function
const deleteArticle = async (index) => {
    const article = articles.value[index];
    // 如果状态是已发布，提示用户联系管理员
    console.log('请联系管理员删除');
    uni.showToast({
      title: '请联系管理员删除',
      icon: 'none',
      duration: 2000
    });
};

const deleteDraft = async (index) => {
	const article = draftArticles.value[index];
	  // 如果状态是草稿，发送删除请求
	  try {
	    const res = await uni.request({
	      url: `${BASE_URL.value}/news/drafts/${article.id}`,
	      method: 'DELETE',
	      header: {
	        'Authorization': `Bearer ${jwtToken.value}`
	      }
	    });
	
	    if (res.data && res.data.message === 'Draft deleted successfully.') {
        uni.showToast({
          title: '删除成功',
          icon: 'none',
          duration: 2000
        });
	      // 从数据中删除该文章
	      articles.value.splice(index, 1); 
	    } else {
	      console.error('删除失败:', res.data.message);
	    }
	  } catch (error) {
	    console.error('Error deleting draft article', error);
	  }
    await fetchArticles();
};

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

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.1;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.header-actions button {
  background: none;
  border: none;
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

.bio {
  font-size: 14px;
  color: #666;
}

.userid {
  font-size: 12px;
  margin-top: 4px;
  color: #666;
  z-index: 10;
}

/* 创作、收藏、知音 */
.stats {
  margin-top: 8px;
  height: 140px;
  gap: 20px;
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

/* 标签切换 */
.tabs { 
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
  border-bottom: 1px solid #ebebeb;
}

.tab {
  padding: 10px 0;
  font-size: 16px;
  color: #666;
  position: relative;
  cursor: pointer;
}
.tab.active {
  font-weight: bold;
  color: #333;
}
.tab.active::after {
  content: "";
  display: block;
  width: 100%;
  height: 2px;
  background-color: #333;
  position: absolute;
  bottom: -1px;
  left: 0;
}

/* 下方卡片列表容器 */
.container {
  padding: 20px;
  margin-top: 0; /* 如果有需要可微调 */
}

/* 其余样式可沿用你原先的 .card, .card-header, .card-body 等... */


.card-list {
  display: flex;
  flex-direction: column;
}

.card {
  margin-bottom: 20px;  /* 卡片之间的间距，改小比如 10px */
  border-radius: 8px;
  padding: 15px;        /* 卡片内容与边框的内边距，可改小比如 10px */
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
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
  margin-bottom: 15px;
  height: 20px;
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

.stats {
  display: flex;
  gap: 10px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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
