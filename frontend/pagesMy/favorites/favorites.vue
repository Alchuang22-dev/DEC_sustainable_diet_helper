<template>
  <view>
    <!-- Header Section -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
	  <button @click="toggleDrawer" class="drawer-button">
		  >
	  </button>
      <input
        class="search-box"
        v-model="searchText"
        @input="onSearchInput"
        :placeholder="placeholderText"
      />
      <button @click="onSearch" class="search-button">
        {{$t('text_search')}}
      </button>
    </view>

    <!-- Loading Indicator -->
    <view v-if="isRefreshing" class="loading-overlay">
      <text class="loading-text">正在加载...</text>
      <!-- 可选：添加动画 -->
      <view class="loading-spinner"></view>
    </view>

    <!-- News Section -->
    <view class="news-section">
      <view
        v-for="(item, index) in filteredNewsItems"
        :key="index"
        :class="['news-item', { active: activeIndex === index }]"
        @click="navigateTo(item.id, item.title)"
        @touchstart="pressFeedback(index)"
        @touchend="releaseFeedback()"
      >
        <view class="news-title">{{ item.title }}</view>
        <view v-if="item.image" class="news-image">
          <image :src="item.image" :alt="item.title" mode="widthFix" />
        </view>
        <view class="news-description">{{ item.description }}</view>
		<view class="news-description">{{ item.info }}</view>
      </view>
    </view>

    <!-- Drawer Component -->
    <uni-drawer
      ref="drawer"
      placement="bottom"
      :mask="mask"
      :width="drawWid"
      :mask-closable="maskClick"
      @close="handleDrawerClose"
      :mask-style="'background-color: rgba(0, 0, 0, 0.5);'"
      :style="'background-color: rgba(255, 255, 255, 0.9);'"
    >
      <view class="drawer-content">
        <button 
          @click="handleSort('favorite')" 
          :class="['nav-item', { active: currentSort === 'favorite' }]"
        >
          我收藏的
        </button>
        <button 
          @click="handleSort('viewed')" 
          :class="['nav-item', { active: currentSort === 'viewed' }]"
        >
          我看过的
        </button>
      </view>
    </uni-drawer>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useNewsStore } from '@/stores/news_list';
import { useI18n } from 'vue-i18n';
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app';
import { storeToRefs } from 'pinia';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储

const newsStore = useNewsStore();
const userStore = useUserStore(); // 使用用户存储

const activeIndex = ref(null);
const currentSort = ref('top-views'); // 默认排序类型

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

const { t } = useI18n();

// 新增：Drawer显示状态通过ref控制
const drawer = ref(null);
const mask = true;
const drawWid = '50%'; // 根据需要调整宽度
const maskClick = true;

// 切换Drawer显示/隐藏的方法
function toggleDrawer() {
  if (drawer.value) {
    drawer.value.open();
    isDrawerVisible.value = true;
  }
}

// 隐藏Drawer的方法
function hideDrawer() {
  if (drawer.value) {
    drawer.value.close();
    isDrawerVisible.value = false;
  }
}

// 处理Drawer关闭事件
function handleDrawerClose() {
  // 可以在这里处理关闭后的逻辑，如重置状态等
  isDrawerVisible.value = false;
  console.log('Drawer closed');
}

// 新增：isDrawerVisible 状态管理
const isDrawerVisible = ref(false);

// 页面跳转方法
function navigateTo(link, name) {
  console.log('跳转至：', link);
  setTimeout(() => {
    uni.navigateTo({
      url: `/pagesNews/news_detail/news_detail?id=${link}`,
    });
  }, 100);
}

// 触摸反馈
function pressFeedback(index) {
  activeIndex.value = index;
}

function releaseFeedback() {
  activeIndex.value = null;
}

// 跳转至新建图文页面
function createNews() {
  uni.navigateTo({
    url: "/pagesNews/create_news/create_news",
  });
}

// 排序功能
function handleSort(sortType) {
  currentSort.value = sortType;
  fetchNews(1, sortType); // 根据排序类型获取新闻
  hideDrawer(); // 排序后关闭抽屉
}

// 异步函数处理下拉刷新
const handlePullDownRefresh = async () => {
  console.log('正在处理下拉刷新...');
  try {
    await newsStore.refreshNews();  // 等待 refreshNews 完成
    uni.stopPullDownRefresh();      // 完成后停止下拉刷新动画
  } catch (error) {
    console.error('Error during refresh:', error);
    uni.stopPullDownRefresh();      // 即使出错也停止刷新
  }
};

// 使用 uni.onPullDownRefresh() 将处理函数绑定到下拉刷新事件
onPullDownRefresh(handlePullDownRefresh);

onShow(() => {
  console.log("用户进入收藏");
  // 根据需求，这里设置为false，可能需要根据实际登录状态调整
  isLoggedIn.value = false; // 显式设置为未登录状态
  console.log("in onShow");
  fetchNews();
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

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: -1;
  opacity: 0.1;
}

/* Header Section */
.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  justify-content: flex-start;
  position: fixed; /* 固定头部 */
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10; /* 确保在页面的最上层 */
  overflow-x: scroll; /* 允许水平滚动 */
  white-space: nowrap; /* 防止内容换行 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); /* 可选，增加阴影效果 */
}

/* 防止按钮换行，确保每个按钮都保持在一行 */
.header button {
  border: none;
  margin-left: 5px;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.3s;
  white-space: nowrap; /* 防止按钮文本换行 */
  padding: 5px 15px;
  flex-shrink: 0; /* 防止按钮被压缩 */
}

.search-box {
  flex: 1;
  padding: 13px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.drawer-button {
  background-color: #ffffff;
  color: black;
  border: none;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.3s;
  white-space: nowrap; /* 防止按钮文本换行 */
  margin-left: 0;
  padding: 5px 15px;
  flex-shrink: 0; /* 防止按钮被压缩 */
}

.search-button {
  background-color: #4CAF50;
  color: white;
  border: none;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.3s;
  white-space: nowrap; /* 防止按钮文本换行 */
  padding: 5px 15px;
  flex-shrink: 0; /* 防止按钮被压缩 */
}

.create-button {
  background-color: #ffffff;
  color: black;
  border: none;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.3s;
  white-space: nowrap; /* 防止按钮文本换行 */
  padding: 5px 15px;
  margin-right: 20px;
  flex-shrink: 0; /* 防止按钮被压缩 */
}

/* 选中的按钮样式 */
.header button.active {
  color: #4caf50; /* 选中状态颜色 */
  font-weight: bold; /* 选中状态加粗 */
}

/* Loading Overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(240, 244, 247, 0.8);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 10; /* 确保在最上层 */
}

.loading-text {
  font-size: 18px;
  color: #333;
  margin-bottom: 10px;
}

/* 加载动画 */
.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #ccc;
  border-top-color: #4caf50;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* News Section */
.news-section {
  padding: 20px;
  padding-top: 70px; /* 根据header的高度调整，确保内容不被遮挡 */
  padding-bottom: 80px;
}

.news-item {
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.1s, box-shadow 0.1s;
  position: relative;              /* 确保其层级设置有效 */
  z-index: 1;                      /* 确保点击事件可以被接收 */
}

.news-item.active {
  transform: scale(0.98);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  background-color: #e6f7ff;
}

.news-image {
  pointer-events: none; /* 确保图片不会阻止父元素的点击事件 */
}

.news-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.news-description {
  font-size: 14px;
  margin-bottom: 10px;
}

/* Footer Toggle Button */

/* Drawer Content */
.drawer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.nav-item {
  text-decoration: none;
  color: #333;
  font-weight: bold;
  cursor: pointer;
  height: 40px;
  transition: color 0.3s, background-color 0.3s;
  margin-bottom: 10px;
  width: 80%; /* 适应抽屉宽度 */
  border-radius: 5px;
}

.nav-item:hover {
  color: #4caf50;
  background-color: rgba(76, 175, 80, 0.1);
}

.nav-item.active {
  color: #4caf50;
  border-bottom: 2px solid #4caf50;
}

/* 关闭按钮样式 */
.close-button {
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 5px;
  padding: 10px 20px;
  cursor: pointer;
  transition: background-color 0.3s;
  width: 80%;
  text-align: center;
}

.close-button:hover {
  background-color: #d32f2f;
}

/* 确保uni-drawer有足够的z-index */
.uni-drawer {
  transition: all 0.3s ease;
  z-index: 1000;
}
</style>
