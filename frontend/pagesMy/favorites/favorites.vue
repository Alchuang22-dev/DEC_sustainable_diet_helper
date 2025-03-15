<template>
  <view>
    <!-- 背景图 -->
    <image src="/static/images/index/background_img.jpg" class="background-image" />

    <!-- 加载中提示 -->
    <view v-if="isRefreshing" class="loading-overlay">
      <text class="loading-text">正在加载...</text>
      <view class="loading-spinner"></view>
    </view>

    <!-- 新闻/收藏 列表 -->
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

    <!-- 功能按钮 -->
    <view class="functions">
      <button @click="toggleDrawer" class="add-btn">
        <image src="@/pagesNews/static/gengduo.png" alt=">" class="icon" />
      </button>
    </view>

    <!-- 抽屉组件 -->
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
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { useNewsStore } from '@/stores/news_list'
import { useI18n } from 'vue-i18n'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import { storeToRefs } from 'pinia'
import { useUserStore } from '../../stores/user'

/* ----------------- Setup ----------------- */
const newsStore = useNewsStore()
const userStore = useUserStore()

/* ----------------- Reactive & State ----------------- */
const activeIndex = ref(null)
const currentSort = ref('top-views')
const drawer = ref(null)
const mask = true
const drawWid = '50%'
const maskClick = true
const isDrawerVisible = ref(false)

// 若需根据后端拼接头像，可自定义 BASE_URL；此处简化为空字符串或自行替换
const BASE_URL = 'https://dechelper.com'

/* ----------------- Computed ----------------- */
const isLoggedIn = computed(() => userStore.user.isLoggedIn)
const uid = computed(() => userStore.user.nickName)
const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
)
const { filteredNewsItems, isRefreshing } = storeToRefs(newsStore)
const { fetchNews } = newsStore

/* ----------------- Lifecycle ----------------- */
onShow(() => {
  console.log('用户进入收藏')
  // 仅示例，此处强制设置为 false
  isLoggedIn.value = false
  // 默认展示「收藏」
  handleSort('favorite')
})

/* ----------------- Methods ----------------- */
function toggleDrawer() {
  if (drawer.value) {
    drawer.value.open()
    isDrawerVisible.value = true
  }
}

function hideDrawer() {
  if (drawer.value) {
    drawer.value.close()
    isDrawerVisible.value = false
  }
}

function handleDrawerClose() {
  isDrawerVisible.value = false
  console.log('Drawer closed')
}

function navigateTo(link, name) {
  console.log('跳转至：', link)
  setTimeout(() => {
    uni.navigateTo({
      url: `/pagesNews/news_detail/news_detail?id=${link}`,
    })
  }, 100)
}

function pressFeedback(index) {
  activeIndex.value = index
}

function releaseFeedback() {
  activeIndex.value = null
}

function handleSort(sortType) {
  currentSort.value = sortType
  fetchNews(1, sortType)
  hideDrawer()
}

// 下拉刷新逻辑
async function handlePullDownRefresh() {
  console.log('正在处理下拉刷新...')
  try {
    await newsStore.refreshNews()
    uni.stopPullDownRefresh()
  } catch (error) {
    console.error('Error during refresh:', error)
    uni.stopPullDownRefresh()
  }
}

/* ----------------- Watchers ----------------- */
onPullDownRefresh(handlePullDownRefresh)
</script>

<style scoped>
/* 背景图 */
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

/* 加载提示 */
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
.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #ccc;
  border-top-color: #4caf50;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 新闻列表 */
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
.news-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}
.news-description {
  font-size: 14px;
  margin-bottom: 10px;
}

/* 悬浮功能按钮 */
.functions {
  position: fixed;
  top: 5%;
  left: 0;
  margin-left: 5px;
  transform: translateY(-50%);
  background-color: rgba(0, 0, 0, 0.25);
  padding: 5px;
  border-radius: 8px;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.add-btn {
  padding: 10px;
  background-color: #ffffff;
  color: black;
  border-radius: 50%;
  border: none;
  font-size: 8px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
}
.icon {
  width: 24px;
  height: 24px;
}

/* Drawer */
.drawer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}
.nav-item {
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

/* 确保 uni-drawer 在上方 */
.uni-drawer {
  z-index: 1000;
  transition: all 0.3s ease;
}
</style>