<template>
  <!-- 用户信息区 -->
  <view class="profile-header">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <view class="profile-info">
      <image :src="avatarSrc" class="avatar" alt="用户头像" />
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
        @click="switchtoPost(index)"
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
/* ----------------- Imports ----------------- */
import {ref, computed} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import {useI18n} from 'vue-i18n'
import {useUserStore} from '../../stores/user'

/* ----------------- Setup ----------------- */
const {t} = useI18n()
const userStore = useUserStore()

/* ----------------- Reactive & State ----------------- */
const BASE_URL = ref('https://xcxcs.uwdjl.cn')
const articles = ref([])

const uid = computed(() => userStore.user.nickName)
const user_id = computed(() => userStore.user.uid)
const avatarSrc = computed(() =>
    userStore.user.avatarUrl
        ? userStore.user.avatarUrl
        : '/static/images/index/default_avatar.jpg'
)

/* ----------------- Computed ----------------- */
const publishedCount = computed(() => articles.value.filter(a => a.status === '已发布').length)
const draftCount = computed(() => articles.value.filter(a => a.status === '草稿').length)

const publishedArticles = computed(() =>
    articles.value.filter(a => a.status === '已发布')
)

/* ----------------- Methods ----------------- */
/**
 * 格式化发布时间
 */
function formatPublishTime(publishTime) {
  const date = new Date(publishTime)
  const now = new Date()

  const isSameDay =
      date.getFullYear() === now.getFullYear() &&
      date.getMonth() === now.getMonth() &&
      date.getDate() === now.getDate()

  if (isSameDay) {
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    return `今天 ${hours}:${minutes}`
  } else {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    return `${year}-${month}-${day}`
  }
}

function switchtoPost(index) {
  const article = articles.value[index]
  uni.navigateTo({
    url: `/pagesNews/news_detail/news_detail?id=${article.id}`
  })
}

/**
 * 获取文章详情
 */
async function getArticleDetails(id) {
  const jwtToken = userStore.user.token
  const url = `${BASE_URL.value}/news/details/news/${id}`

  const res = await uni.request({
    url,
    method: 'GET',
    header: {
      'Authorization': `Bearer ${jwtToken}`
    }
  })

  if (res.statusCode === 200) {
    return res.data
  } else {
    console.error(`获取文章详情失败: ${res.statusCode}`)
    return null
  }
}

/**
 * 获取用户和发布的文章列表
 */
async function fetchData() {
  const userId = getUserIdFromRoute()
  if (!userId) {
    uni.showToast({
      title: '用户ID未找到',
      icon: 'none',
      duration: 2000
    })
    return
  }

  const jwtToken = userStore.user.token
  const res = await uni.request({
    url: `${BASE_URL.value}/users/${userId}/profile`,
    method: 'GET',
    header: {
      'Authorization': `Bearer ${jwtToken}`
    }
  })

  if (res.statusCode === 200) {
    const data = res.data
    // 更新用户头像
    userStore.user.nickName = data.nickname
    userStore.user.avatarUrl = `${BASE_URL.value}/static/${data.avatar_url}`

    // 获取并拼装新闻详情
    const newsDetailsPromises = data.news.map(newsItem => getArticleDetails(newsItem.id))
    const newsDetails = await Promise.all(newsDetailsPromises)
    const validNewsDetails = newsDetails.filter(detail => detail !== null)

    articles.value = validNewsDetails.map(detail => ({
      ...detail,
      publishTime: detail.upload_time || detail.updated_at,
      likes: detail.like_count || 0,
      favorites: detail.favorite_count || 0,
      shares: detail.share_count || 0,
      status: detail.status || '已发布',
      bgColor:
          detail.status === '已发布'
              ? 'rgba(0, 123, 255, 0.1)'
              : 'rgba(255, 193, 7, 0.1)'
    }))
  } else if (res.statusCode === 401) {
    uni.showToast({
      title: '未授权，请重新登录',
      icon: 'none',
      duration: 2000
    })
  } else if (res.statusCode === 404) {
    uni.showToast({
      title: '用户未找到',
      icon: 'none',
      duration: 2000
    })
  } else {
    uni.showToast({
      title: '获取用户数据失败',
      icon: 'none',
      duration: 2000
    })
  }
}

/**
 * 从路由参数中获取 userId
 */
function getUserIdFromRoute() {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  return currentPage.options.id
}

/* ----------------- Lifecycle ----------------- */
onShow(() => {
  fetchData()
})
</script>

<style scoped>
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

/* 全屏背景图片 */
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

.profile-header {
  margin-top: 20px;
  position: relative;
  width: 100%;
  height: 220px;
  background-color: #f5f5f5;
  overflow: hidden;
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

.container {
  padding: 20px;
  margin-top: 20px;
}

.card-list {
  display: flex;
  flex-direction: column;
}

.card {
  margin-bottom: 10px;
  border-radius: 8px;
  padding: 10px;
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
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
  margin-bottom: 10px;
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
  gap: 5px;
}

.publish-time {
  font-size: 12px;
  color: #777;
}

.like-count,
.favorite-count,
.share-count {
  font-size: 12px;
  color: #777;
}
</style>