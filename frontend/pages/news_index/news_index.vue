<template>
  <view class="container">
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    />

    <view class="header">
      <!-- 抽屉触发按钮 -->
      <uni-icons
        @click="toggleDrawer"
        type="more-filled"
        size="24"
        color="#666"
        class="drawer-trigger"
      />

      <!-- 搜索框 -->
      <uni-search-bar
        v-model="searchText"
        :placeholder="t('text_search')"
        class="search-bar"
        @confirm="onSearch"
        cancelButton="none"
      />

      <!-- "去信"按钮：仅登录可见 -->
      <view v-if="isLoggedIn" class="create-news-wrapper">
        <uni-icons
          @click="createNews"
          type="compose"
          size="24"
          color="#4caf50"
          class="create-icon"
        />
      </view>
    </view>

    <!-- Loading 状态 -->
    <uni-load-more v-if="isRefreshing" status="loading" :content-text="loadMoreText" />

    <!-- 新闻列表 -->
    <view class="news-section">
      <uni-card
        v-for="(item, index) in filteredNewsItems"
        :key="index"
        :title="item.title"
        :extra="item.info"
        :thumbnail="item.image"
        @click="navigateTo(item.id)"
        :class="['news-card', { active: activeIndex === index }]"
        @touchstart="pressFeedback(index)"
        @touchend="releaseFeedback"
      >
        <text class="news-description">{{ item.description }}</text>
      </uni-card>
    </view>

    <!-- 抽屉：排序选项 -->
    <uni-drawer
      ref="drawer"
      mode="bottom"
      :mask="true"
      :maskClick="true"
      @close="handleDrawerClose"
    >
      <view class="drawer-content">
        <uni-list>
          <uni-list-item
            v-for="(sort, idx) in sortOptions"
            :key="idx"
            :title="t(sort.text)"
            :showArrow="true"
            clickable
            @click="handleSort(sort.value)"
            :class="{ 'active-sort': currentSort === sort.value }"
          />
        </uni-list>
      </view>
    </uni-drawer>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useNewsStore } from '@/stores/news_list'
import { useUserStore } from '../../stores/user'

/* ----------------- Setup ----------------- */
const newsStore = useNewsStore()
const userStore = useUserStore()
const { filteredNewsItems, isRefreshing } = storeToRefs(newsStore)
const { refreshNews, fetchNews } = newsStore
const { t } = useI18n()

/* ----------------- Reactive & State ----------------- */
const isLoggedIn = computed(() => userStore.user.isLoggedIn)
const drawer = ref(null)
const searchText = ref('')
const currentSort = ref('top-views')
const activeIndex = ref(null)

/* ----------------- Computed ----------------- */
// 加载更多文字
const loadMoreText = {
  contentdown: t('pull_down_text'),
  contentrefresh: t('loading_text'),
  contentnomore: t('no_more_data')
}

// 排序选项
const sortOptions = [
  { text: 'sort_by_time', value: 'latest' },
  { text: 'sort_by_views', value: 'top-views' },
  { text: 'sort_by_likes', value: 'top-likes' }
]

/* ----------------- Lifecycle ----------------- */
onShow(() => {
  uni.setNavigationBarTitle({ title: t('news_index') })
  uni.setTabBarItem({ index: 0, text: t('index') })
  uni.setTabBarItem({ index: 1, text: t('tools_index') })
  uni.setTabBarItem({ index: 2, text: t('news_index') })
  uni.setTabBarItem({ index: 3, text: t('my_index') })

  fetchNews()
})

onPullDownRefresh(async () => {
  try {
    await refreshNews()
    uni.stopPullDownRefresh()
  } catch (error) {
    console.error('Error during refresh:', error)
    uni.stopPullDownRefresh()
  }
})

/* ----------------- Methods ----------------- */
function toggleDrawer() {
  if (drawer.value) {
    drawer.value.open()
  }
}

function handleDrawerClose() {
  if (drawer.value) {
    drawer.value.close()
  }
}

function handleSort(sortType) {
  currentSort.value = sortType
  fetchNews(1, sortType)
  handleDrawerClose()
}

function onSearch() {
  fetchNews(1, 'search', searchText.value)
}

function createNews() {
  uni.navigateTo({
    url: "/pagesNews/create_news/create_news"
  })
}

function navigateTo(id) {
  uni.navigateTo({
    url: `/pagesNews/news_detail/news_detail?id=${id}`
  })
}

function pressFeedback(index) {
  activeIndex.value = index
}

function releaseFeedback() {
  activeIndex.value = null
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
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
  pointer-events: none;
}

.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background-color: #ffffff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.drawer-trigger {
  margin-right: 12px;
}

.search-bar {
  flex: 1;
}

.create-news-wrapper {
  margin-left: 12px;
}

.news-section {
  padding: 80px 16px 16px;
}

.news-card {
  margin-bottom: 16px;
  transition: transform 0.2s ease;
}

.news-card.active {
  transform: scale(0.98);
}

.news-description {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.drawer-content {
  padding: 16px 0;
}

.active-sort {
  background-color: #f0f0f0;
}

:deep(.uni-searchbar) {
  padding: 0 !important;
}

:deep(.uni-searchbar__box) {
  height: 36px !important;
}

:deep(.uni-card) {
  margin: 0 0 16px 0 !important;
  padding: 0 !important;
  border-radius: 8px !important;
}

:deep(.uni-card__header) {
  padding: 12px !important;
}

:deep(.uni-card__content) {
  padding: 12px !important;
}

:deep(.uni-list-item) {
  padding: 12px 16px !important;
}

:deep(.uni-list-item__content-title) {
  font-size: 16px !important;
}

:deep(.uni-card__header-content-title) {
  font-weight: bold !important;
}
</style>