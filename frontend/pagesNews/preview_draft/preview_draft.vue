<template>
  <view class="container">
    <view class="author-header">
      <image :src="post.authoravatar" class="author-avatar" />
      <text class="author-username">{{ post.authorname }}</text>
    </view>

    <view class="title-container">
      <h1 class="article-title">{{ post.title }}</h1>
      <p class="article-description">{{ post.description }}</p>
    </view>

    <view class="components-container">
      <view v-for="component in post.components" :key="component.id">
        <view v-if="component.style === 'text'" class="text-content">
          <p>{{ component.content }}</p>
        </view>
        <view v-if="component.style === 'image'" class="image-content">
          <image :src="component.content" class="image" mode="widthFix" />
          <p class="image-description">{{ component.description }}</p>
        </view>
      </view>
    </view>

    <view class="post-time">{{ post.savetime }}</view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useUserStore } from '@/stores/user'

/* ----------------- Setup ----------------- */
const userStore = useUserStore()
const BASE_URL = 'https://dechelper.com'
const PageId = ref('')

const post = ref({ components: [] })
const jwtToken = computed(() => userStore.user.token)
const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
)

/* ----------------- Lifecycle ----------------- */
onLoad(async (options) => {
  const articleId = options.id
  PageId.value = articleId

  const details = await getArticleDetails(PageId.value)
  if (!details) return

  post.value = {
    id: details.id,
    authoravatar: avatarSrc.value,
    authorname: details.author.nickname,
    authorid: details.author.id,
    savetime: details.updated_at,
    title: details.title,
    description: details.paragraphs[0]?.text || '',
    components: []
  }

  const totalItems = Math.max(details.paragraphs.length, details.images.length)
  for (let i = 1; i < totalItems; i++) {
    if (details.paragraphs[i]?.text) {
      post.value.components.push({
        id: i + 1,
        content: details.paragraphs[i].text,
        style: 'text'
      })
    }
    if (details.images[i]?.url) {
      post.value.components.push({
        id: i + 1,
        content: details.images[i].url,
        style: 'image',
        description: details.images[i].description || ''
      })
    }
  }
})

function getArticleDetails(id) {
  return new Promise((resolve) => {
    uni.request({
      url: `${BASE_URL}/news/details/draft/${id}`,
      method: 'GET',
      header: {
        Authorization: `Bearer ${jwtToken.value}`
      },
      success: (res) => {
        resolve(res.data)
      },
      fail: (err) => {
        console.error('Error fetching article details', err)
        resolve(null)
      }
    })
  })
}
</script>

<style scoped>
.container {
  padding: 20px;
}

.author-header {
  display: flex;
  margin-bottom: 10px;
}

.author-avatar {
  width: 50px;
  height: 50px;
  background-color: #ccc;
  border-radius: 50%;
  margin-right: 15px;
}

.author-username {
  font-weight: bold;
  margin-right: 20px;
  display: flex;
  align-items: center;
}

.title-container {
  margin-bottom: 20px;
}

.article-title {
  font-size: 26px;
  font-weight: bold;
  margin-bottom: 10px;
}

.article-description {
  font-size: 18px;
  color: #666;
}

.components-container {
  margin-top: 20px;
  margin-bottom: 20px;
}

.text-content p {
  margin-top: 10px;
  font-size: 16px;
  line-height: 1.5;
  margin-bottom: 10px;
}

.image-content {
  margin-top: 10px;
  margin-bottom: 20px;
}

.image {
  width: 100%;
  border-radius: 8px;
}

.image-description {
  font-size: 14px;
  color: #777;
  margin-top: 10px;
}

.post-time {
  font-size: 14px;
  color: #888;
  text-align: right;
  margin-top: 20px;
}
</style>