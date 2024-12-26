<template>
  <!-- 个人主页 -->
  <view class="profile-header">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <view class="profile-info">
      <image :src="avatarSrc" class="avatar" alt="用户头像" />
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
  </view>

  <!-- 标签切换：已发布 / 草稿 -->
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
    <!-- 已发布列表 -->
    <view v-if="currentTab === 'published'" class="card-list">
      <view
        v-for="(item, index) in publishedArticles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
      >
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
          <button @click="viewArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/view.svg" class="icon" alt="View" />
          </button>
          <button @click="editArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/edit.svg" class="icon" alt="Edit" />
          </button>
          <button @click="deleteArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/delete.svg" class="icon" alt="Delete" />
          </button>
        </view>
      </view>
    </view>

    <!-- 草稿列表 -->
    <view v-else class="card-list">
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
            <image src="@/pagesMy/static/view.svg" class="icon" alt="View" />
          </button>
          <button @click="editDraft(index)" class="action-btn">
            <image src="@/pagesMy/static/edit.svg" class="icon" alt="Edit" />
          </button>
          <button @click="deleteDraft(index)" class="action-btn">
            <image src="@/pagesMy/static/delete.svg" class="icon" alt="Delete" />
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { useNewsStore } from '@/stores/news_list'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../../stores/user'
import { onShow } from '@dcloudio/uni-app'

/* ----------------- Setup ----------------- */
const newsStore = useNewsStore()
const userStore = useUserStore()
const { t, locale, messages } = useI18n()

/* ----------------- Reactive & State ----------------- */
const articles = ref([])
const activeIndex = ref(null)
const backgroundImageUrl = ref('/static/images/index/background_img.jpg')
const currentTab = ref('published')

// 后端地址，若无需动态可直接写死；暂示例
const BASE_URL = ref('https://xcxcs.uwdjl.cn')

/* ----------------- Computed ----------------- */
const user_id = computed(() => userStore.user.uid)
const isLoggedIn = computed(() => userStore.user.isLoggedIn)
const uid = computed(() => userStore.user.nickName)
const jwtToken = computed(() => userStore.user.token)

const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL.value}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
)

const publishedCount = computed(() =>
  articles.value.filter(a => a.status === '已发布').length
)
const draftCount = computed(() =>
  articles.value.filter(a => a.status === '草稿').length
)
const publishedArticles = computed(() =>
  articles.value.filter(a => a.status === '已发布')
)
const draftArticles = computed(() =>
  articles.value.filter(a => a.status === '草稿')
)

/* ----------------- Lifecycle ----------------- */
onShow(() => {
  fetchArticles()
})

/* ----------------- Methods ----------------- */
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

function editBackgroundImage() {
  uni.chooseImage({
    count: 1,
    success: res => {
      const tempFilePath = res.tempFilePaths[0]
      backgroundImageUrl.value = tempFilePath
    }
  })
}

async function fetchArticles() {
  const publishedIds = await getPublishedNewsIds()
  const draftIds = await getDraftNewsIds()

  const allArticles = []

  for (const id of publishedIds) {
    const details = await getArticleDetails(id)
    if (details) {
      allArticles.push({
        ...details,
        publishTime: details.upload_time,
        likes: details.like_count,
        favorites: details.favorite_count,
        shares: details.share_count,
        status: '已发布',
        bgColor: 'rgba(0, 123, 255, 0.1)'
      })
    }
  }

  for (const id of draftIds) {
    const details = await getArticleDetails(id, true)
    if (details) {
      allArticles.push({
        ...details,
        publishTime: details.updated_at,
        status: '草稿',
        bgColor: 'rgba(255, 193, 7, 0.1)'
      })
    }
  }

  articles.value = allArticles
}

async function getPublishedNewsIds() {
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/my_news`,
      method: 'GET',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      }
    })
    return res.data.news_ids || []
  } catch (error) {
    console.error('Error fetching published news IDs', error)
    return []
  }
}

async function getDraftNewsIds() {
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/my_drafts`,
      method: 'GET',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      }
    })
    return res.data.draft_ids || []
  } catch (error) {
    console.error('Error fetching draft news IDs', error)
    return []
  }
}

async function getArticleDetails(id, isDraft = false) {
  const url = isDraft
    ? `${BASE_URL.value}/news/details/draft/${id}`
    : `${BASE_URL.value}/news/details/news/${id}`
  try {
    const res = await uni.request({
      url,
      method: 'GET',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      }
    })
    return res.data
  } catch (error) {
    console.error('Error fetching article details', error)
    return null
  }
}

function viewArticle(index) {
  const article = articles.value[index]
  uni.navigateTo({
    url: `/pagesNews/news_detail/news_detail?id=${article.id}`
  })
}

function viewDraft(index) {
  const draft = draftArticles.value[index]
  uni.navigateTo({
    url: `/pagesNews/preview_draft/preview_draft?id=${draft.id}`
  })
}

function editArticle(index) {
  uni.showToast({
    title: '发布后不可编辑',
    icon: 'none',
    duration: 2000
  })
}

function editDraft(index) {
  const draft = draftArticles.value[index]
  uni.navigateTo({
    url: `/pagesNews/edit_draft/edit_draft?id=${draft.id}`
  })
}

async function deleteArticle(index) {
  const article = publishedArticles.value[index]
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/${article.id}`,
      method: 'DELETE',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      }
    })
    if (res.data && res.data.message === 'News deleted successfully.') {
      uni.showToast({
        title: '删除成功',
        icon: 'none',
        duration: 2000
      })
      // 移除已删除的作品
      const articleId = article.id
      const targetIndexInAll = articles.value.findIndex(a => a.id === articleId)
      if (targetIndexInAll !== -1) {
        articles.value.splice(targetIndexInAll, 1)
      }
    } else {
      console.error('删除失败:', res.data)
      uni.showToast({
        title: '删除失败',
        icon: 'none',
        duration: 2000
      })
    }
  } catch (error) {
    console.error('Error deleting published article', error)
    uni.showToast({
      title: '删除出现错误',
      icon: 'none',
      duration: 2000
    })
  }
  await fetchArticles()
}

async function deleteDraft(index) {
  const article = draftArticles.value[index]
  try {
    const res = await uni.request({
      url: `${BASE_URL.value}/news/drafts/${article.id}`,
      method: 'DELETE',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      }
    })
    if (res.data && res.data.message === 'Draft deleted successfully.') {
      uni.showToast({
        title: '删除成功',
        icon: 'none',
        duration: 2000
      })
      articles.value.splice(index, 1)
    } else {
      console.error('删除失败:', res.data.message)
    }
  } catch (error) {
    console.error('Error deleting draft article', error)
  }
  await fetchArticles()
}
</script>

<style scoped>

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
  z-index: 10;
}
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

.container {
  padding: 20px;
  margin-top: 0;
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
  transform: scale(1.2);
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