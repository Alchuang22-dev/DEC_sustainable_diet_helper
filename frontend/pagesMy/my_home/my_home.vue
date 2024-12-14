<template>
  <image src="/static/images/index/background_img.jpg" class="background-image"></image>
  <view class="header">
    <text class="title">{{$t('menu_creations')}}</text>
    <view class="header-actions">
      <button class="menu-icon"></button>
      <button class="camera-icon"></button>
    </view>
  </view>
  <view class="container">
    <!-- 图文卡片列表 -->
    <view class="card-list">
      <view
        v-for="(item, index) in articles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
      >
        <view class="card-header">
          <view class="title">{{ item.title }}</view>
          <view class="status">{{ item.status }}</view>
        </view>
        <view class="card-body">
          <view class="description">{{ item.description }}</view>
          <view class="info">
            <text class="publish-time">{{ item.publishTime }}</text>
            <view class="stats">
              <text class="like-count">👍 {{ item.likes }}</text>
              <text class="favorite-count">⭐ {{ item.favorites }}</text>
              <text class="share-count">🔗 {{ item.shares }}</text>
            </view>
          </view>
        </view>
        <view class="card-footer">
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

const activeIndex = ref(null);
// 计算属性从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);
const uid = computed(() => userStore.user.nickName);
const avatarSrc = computed(() =>
    userStore.user.avatarUrl
        ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);

// 从 Store 获取数据和方法
const { filteredNewsItems, selectedSection, isRefreshing } = storeToRefs(newsStore);
const { setSection, refreshNews, fetchNews } = newsStore;


const BASE_URL = ref('http://122.51.231.155:8080');

// 模拟数据
const articles = ref([]);
const { t, locale, messages } = useI18n();
const jwtToken = computed(() => userStore.user.token);; // Replace with actual token

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
        status: '草稿',
        bgColor: 'rgba(255, 193, 7, 0.1)', // Draft color
      });
    }
  }

  articles.value = allArticles;
};

// Lifecycle hook to load articles
onMounted(() => {
  fetchArticles();
});

// View article function
const viewArticle = (index) => {
  console.log('查看文章:', articles.value[index]);
  // 跳转到文章详情页
};

// Edit article function
const editArticle = (index) => {
  const article = articles.value[index];
  console.log('编辑文章:', article);

  // 将文章的 ID 作为查询参数传递到新页面
  uni.navigateTo({
    url: `/pagesNews/edit_draft/edit_draft?id=${article.id}`,
  });
};

// Delete article function
// Delete article function
const deleteArticle = async (index) => {
  const article = articles.value[index];

  if (article.status === '草稿') {
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
        console.log('草稿删除成功');
        // 从数据中删除该文章
        articles.value.splice(index, 1); 
      } else {
        console.error('删除失败:', res.data.message);
      }
    } catch (error) {
      console.error('Error deleting draft article', error);
    }
  } else if (article.status === '已发布') {
    // 如果状态是已发布，提示用户联系管理员
    console.log('请联系管理员删除');
    uni.showToast({
      title: '请联系管理员删除',
      icon: 'none',
      duration: 2000
    });
  }
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

.container {
  padding: 20px;
}

.card-list {
  display: flex;
  flex-direction: column;
}

.card {
  margin-bottom: 20px;
  border-radius: 8px;
  padding: 15px;
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
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
