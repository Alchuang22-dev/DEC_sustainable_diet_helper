<template>
  <view>
    <!-- Header: 包含搜索、排序抽屉、以及“去信”按钮 -->
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    ></image>

    <view class="header">
      <!-- 打开抽屉按钮 -->
      <button @click="toggleDrawer" class="drawer-button">
        <image src="@/pagesNews/static/gengduo.png" alt=">" class="icon-news"></image>
      </button>

      <!-- 搜索框 -->
      <input
        class="search-box"
        v-model="searchText"
        :placeholder="t('text_search')"
      />
      <button @click="onSearch" class="search-button">
        {{ t('text_search') }}
      </button>

      <!-- “去信”按钮：仅已登录时可见 -->
      <button v-if="isLoggedIn" @click="createNews" class="create-button">
        {{ t('create_news_button') }}
      </button>
    </view>

    <!-- Loading 状态 -->
    <view v-if="isRefreshing" class="loading-overlay">
      <text class="loading-text">{{ t('loading_text') }}</text>
      <view class="loading-spinner"></view>
    </view>

    <!-- News 列表 -->
    <view class="news-section">
      <view
        v-for="(item, index) in filteredNewsItems"
        :key="index"
        :class="['news-item', { active: activeIndex === index }]"
        @click="navigateTo(item.id, item.title)"
        @touchstart="pressFeedback(index)"
        @touchend="releaseFeedback"
      >
        <view class="news-title">{{ item.title }}</view>
        <view v-if="item.image" class="news-image">
          <image :src="item.image" :alt="item.title" mode="widthFix" />
        </view>
        <view class="news-description">{{ item.description }}</view>
        <view class="news-description">{{ item.info }}</view>
      </view>
    </view>

    <!-- Drawer: 用于多种排序 -->
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
          @click="handleSort('latest')"
          :class="['nav-item', { active: currentSort === 'latest' }]"
        >
          {{ t('sort_by_time') }}
        </button>
        <button
          @click="handleSort('top-views')"
          :class="['nav-item', { active: currentSort === 'top-views' }]"
        >
          {{ t('sort_by_views') }}
        </button>
        <button
          @click="handleSort('top-likes')"
          :class="['nav-item', { active: currentSort === 'top-likes' }]"
        >
          {{ t('sort_by_likes') }}
        </button>
      </view>
    </uni-drawer>
  </view>
</template>

<script setup>
/**
 * 新闻页面示例：包括新闻列表、搜索、排序，以及创建新闻入口
 * - 使用下拉刷新获取新闻
 * - 抽屉式排序功能
 */

import { ref, computed } from 'vue'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useNewsStore } from '@/stores/news_list'
import { useUserStore } from '../../stores/user'

const newsStore = useNewsStore()
const userStore = useUserStore()

// 取出 store 的状态
const { filteredNewsItems, isRefreshing } = storeToRefs(newsStore)
const { refreshNews, fetchNews } = newsStore

// 国际化
const { t } = useI18n()

// 当前登录状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn)

// 抽屉相关
const drawer = ref(null)
const mask = true
const drawWid = '50%'
const maskClick = true
const isDrawerVisible = ref(false)

// 搜索相关
const searchText = ref('')

// 当前排序方式
const currentSort = ref('top-views')

// 触摸反馈
const activeIndex = ref(null)

/**
 * 切换抽屉
 */
function toggleDrawer() {
  if (drawer.value) {
    drawer.value.open()
    isDrawerVisible.value = true
  }
}

/**
 * 抽屉关闭
 */
function handleDrawerClose() {
  isDrawerVisible.value = false
}

/**
 * 抽屉：执行排序
 * @param {string} sortType
 */
function handleSort(sortType) {
  currentSort.value = sortType
  fetchNews(1, sortType)
  hideDrawer()
}

/**
 * 隐藏抽屉
 */
function hideDrawer() {
  if (drawer.value) {
    drawer.value.close()
    isDrawerVisible.value = false
  }
}

/**
 * 跳转到新闻详情
 * @param {string} link 新闻id
 * @param {string} name 新闻标题
 */
function navigateTo(link, name) {
  setTimeout(() => {
    uni.navigateTo({
      url: `/pagesNews/news_detail/news_detail?id=${link}`
    })
  }, 100)
}

/**
 * 按下时高亮
 */
function pressFeedback(index) {
  activeIndex.value = index
}

/**
 * 松开时取消高亮
 */
function releaseFeedback() {
  activeIndex.value = null
}

/**
 * 触发创建新闻
 */
function createNews() {
  uni.navigateTo({
    url: "/pagesNews/create_news/create_news"
  })
}

/**
 * 搜索按钮点击
 */
function onSearch() {
  fetchNews(1, 'search', searchText.value)
}

/**
 * 下拉刷新
 */
async function handlePullDownRefresh() {
  try {
    await refreshNews()
    uni.stopPullDownRefresh()
  } catch (error) {
    console.error('Error during refresh:', error)
    uni.stopPullDownRefresh()
  }
}

onPullDownRefresh(handlePullDownRefresh)

/**
 * 页面显示时：初始化标题，获取新闻等
 */
onShow(() => {
  // 设置页面标题
  uni.setNavigationBarTitle({ title: t('news_index') })
  // 更新底部Tab
  uni.setTabBarItem({ index: 0, text: t('index') })
  uni.setTabBarItem({ index: 1, text: t('tools_index') })
  uni.setTabBarItem({ index: 2, text: t('news_index') })
  uni.setTabBarItem({ index: 3, text: t('my_index') })

  // 默认获取一次数据
  fetchNews()
})
</script>

<style scoped>
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
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10;
  overflow-x: scroll;
  white-space: nowrap;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.header button {
  border: none;
  margin-left: 5px;
  font-size: 16px;
  cursor: pointer;
  transition: color 0.3s;
  white-space: nowrap;
  padding: 5px 15px;
  flex-shrink: 0;
}

/* 搜索框 */
.search-box {
  flex: 1;
  padding: 13px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

/* 抽屉按钮 */
.drawer-button {
  background-color: #ffffff;
  color: black;
  border: none;
  font-size: 16px;
  cursor: pointer;
  white-space: nowrap;
  margin-left: 0;
  flex-shrink: 0;
  align-items: center;
}

/* 搜索按钮 */
.search-button {
  background-color: #4caf50;
  color: white;
  border: none;
  font-size: 16px;
  cursor: pointer;
  white-space: nowrap;
  padding: 5px 15px;
  flex-shrink: 0;
}

/* “去信”按钮 */
.create-button {
  background-color: #ffffff;
  color: black;
  border: none;
  font-size: 16px;
  cursor: pointer;
  white-space: nowrap;
  padding: 5px 15px;
  margin-right: 20px;
  flex-shrink: 0;
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
  z-index: 10;
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
  padding-top: 70px;
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
  position: relative;
  z-index: 1;
}

.news-item.active {
  transform: scale(0.98);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  background-color: #e6f7ff;
}

.news-image {
  pointer-events: none;
}

/* 标题与描述 */
.news-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.news-description {
  font-size: 14px;
  margin-bottom: 10px;
}

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
  width: 80%;
  border-radius: 5px;
  text-align: center;
}

.nav-item:hover {
  color: #4caf50;
  background-color: rgba(76, 175, 80, 0.1);
}

.nav-item.active {
  color: #4caf50;
  border-bottom: 2px solid #4caf50;
}

.icon-news {
  width: 20px;
  height: 20px;
}
</style>