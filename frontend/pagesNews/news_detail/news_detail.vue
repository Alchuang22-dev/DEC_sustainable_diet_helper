<template>
  <view class="page-container">
    <!-- Author header section -->
    <view class="author-section">
      <view class="author-info">
        <image
          :src="formatAvatar(post.authoravatar)"
          class="author-avatar"
          @click="switchtoStranger(post.authorid)"
        />
        <text class="author-name">{{ post.authorname }}</text>
      </view>
    </view>

    <!-- Article content section -->
    <view class="content-section">
      <view class="title-wrapper">
        <text class="title">{{ post.title }}</text>
        <text class="description">{{ post.description }}</text>
      </view>

      <view class="components">
        <view v-for="component in post.components" :key="component.id" class="component-item">
          <view v-if="component.style === 'text'" class="text-block">
            <text>{{ component.content }}</text>
          </view>
          <view v-if="component.style === 'image'" class="image-block">
            <image :src="component.content" class="content-image" mode="widthFix" />
            <text v-if="component.description" class="image-caption">{{ component.description }}</text>
          </view>
        </view>
      </view>

      <view class="metadata">
        <text class="timestamp">{{ formattedSaveTime }}</text>
        <text class="views">阅读量：{{ post.viewCount }}</text>
      </view>

      <view class="interaction-bar">
        <view
          class="interaction-btn"
          :class="{ active: ifLike }"
          @click="toggleInteraction('like')"
        >
          <image
            :src="ifLike ? '/pagesNews/static/liked.svg' : '/pagesNews/static/like.svg'"
            class="interaction-icon"
          />
          <text>{{ formatCount(post.likeCount) }}</text>
        </view>

        <view
          class="interaction-btn"
          :class="{ active: ifFavourite }"
          @click="toggleInteraction('favorite')"
        >
          <image
            :src="ifFavourite ? '/pagesNews/static/favorited.svg' : '/pagesNews/static/favorite.svg'"
            class="interaction-icon"
          />
          <text>{{ formatCount(post.favoriteCount) }}</text>
        </view>

        <view
          class="interaction-btn"
          :class="{ active: ifDislike }"
          @click="toggleInteraction('dislike')"
        >
          <image
            :src="ifDislike ? '/pagesNews/static/disliked.svg' : '/pagesNews/static/dislike.svg'"
            class="interaction-icon"
          />
          <text>dis</text>
        </view>
      </view>

      <!-- Comments section -->
      <view class="comments-section">
        <text class="section-title">注释与说明</text>

        <view class="comments-list">
          <view
            v-for="(comment, index) in limitedComments"
            :key="comment.id"
            class="comment-card"
          >
            <view class="comment-header">
              <view class="commenter-info">
                <image
                  :src="formatAvatar(comment.authorAvatar)"
                  class="commenter-avatar"
                  @click="switchtoStranger(comment.authorid)"
                />
                <text class="commenter-name">{{ comment.authorName }}</text>
              </view>
              <text class="comment-time">{{ comment.publish_time }}</text>
            </view>

            <view class="comment-body">
              <rich-text :nodes="renderCommentText(comment.text)" class="comment-text" />
            </view>

            <view class="comment-actions">
              <view
                class="action-btn like-btn"
                :class="{ active: comment.liked }"
                @click="toggleCommentLike(index)"
              >
                <image
                  :src="comment.liked ? '/pagesNews/static/liked.svg' : '/pagesNews/static/like.svg'"
                  class="action-icon"
                />
                <text>{{ formatCount(comment.likecount) }}</text>
              </view>
              <view class="action-btn reply-btn" @click="replyToComment(index)">
                <uni-icons type="chat" size="16" />
                <text>解决</text>
              </view>
            </view>

            <view v-if="replyingTo === index" class="reply-input">
              <input
                type="text"
                v-model="newReply"
                placeholder="回复..."
                class="reply-field"
              />
              <view class="send-btn" @click="addReply(index)">发送</view>
            </view>

            <view v-if="comment.replies.length > 0" class="replies-list">
              <view
                v-for="reply in limitedReplies(comment)"
                :key="reply.id"
                class="reply-item"
              >
                <text class="reply-author">{{ reply.authorName }}</text>
                <text class="reply-content">{{ reply.text }}</text>
                <text class="reply-time">{{ reply.publish_time }}</text>
              </view>

              <view
                v-if="comment.replies.length > 3"
                class="show-more"
                @click="toggleReplies(comment)"
              >
                <text>
                  {{
                    !comment.showAllReplies
                      ? `还有 ${comment.replies.length - 3} 条附加说明`
                      : '收起附加说明'
                  }}
                </text>
              </view>
            </view>
          </view>

          <view
            v-if="comments.length > 5"
            class="show-more-comments"
            @click="toggleComments"
          >
            <text>
              {{
                !showAllComments
                  ? `还有 ${comments.length - 5} 条注释`
                  : '收起注释'
              }}
            </text>
          </view>
        </view>

        <view class="add-comment">
          <input
            type="text"
            v-model="newComment"
            @input="handleCommentInput"
            placeholder="在此处添加说明"
            class="comment-input"
          />
          <view class="submit-btn" @click="addComment">注释</view>
        </view>

        <uni-popup ref="mentionPopup" type="bottom" :mask="false">
          <view class="mentions-list">
            <view
              v-for="(name, idx) in userListForMentions"
              :key="idx"
              class="mention-item"
              @click="insertMention(name)"
            >
              @{{ name }}
            </view>
          </view>
        </uni-popup>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, reactive, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useUserStore } from '../../stores/user'

/* ----------------- Setup ----------------- */
const userStore = useUserStore()
const BASE_URL = 'http://xcxcs.uwdjl.cn:8080'
const PageId = ref('')

const post = ref({ components: [] })
const comments = reactive([])
const newComment = ref('')
const replyingTo = ref(null)
const newReply = ref('')
const showAllComments = ref(false)

const ifLike = ref(false)
const ifFavourite = ref(false)
const ifDislike = ref(false)

/* 用于当前用户信息 */
const uid = computed(() => userStore.user.nickName)
const user_id = computed(() => userStore.user.uid)
const jwtToken = computed(() => userStore.user.token)
const avatarSrc_sh = computed(() =>
  userStore.user.avatarUrl
    ? userStore.user.avatarUrl
    : '/static/images/index/background_img.jpg'
)

/* ----------------- Computed ----------------- */
const limitedComments = computed(() => {
  if (showAllComments.value) {
    return comments
  }
  return comments.slice(0, 5)
})

// 计算每个评论的有限回复
function limitedReplies(comment) {
  if (comment.showAllReplies) {
    return comment.replies
  } else {
    return comment.replies.slice(0, 3)
  }
}

function formatCount(count) {
  return count < 10000 ? count : (count / 1000).toFixed(1) + 'k'
}

const systemDate = new Date()
const systemDateStr = systemDate.toISOString().slice(0, 10)

const formattedSaveTime = computed(() => {
  const postDate = post.value.savetime?.slice(0, 10) || ''
  if (postDate === systemDateStr) {
    const postTime = new Date(post.value.savetime)
    const hours = String(postTime.getHours()).padStart(2, '0')
    const minutes = String(postTime.getMinutes()).padStart(2, '0')
    const seconds = String(postTime.getSeconds()).padStart(2, '0')
    return `今天 ${hours}:${minutes}:${seconds}`
  } else {
    return postDate
  }
})

/* ----------------- Methods ----------------- */
function formatAvatar(path) {
  return `${BASE_URL}/static/${path}`
}

function toggleComments() {
  showAllComments.value = !showAllComments.value
}

function toggleReplies(comment) {
  comment.showAllReplies = !comment.showAllReplies
}

function renderCommentText(text) {
  if (!text.includes('@')) {
    return text
  }
  return text.replace(/@(\S+)\s/g, '<span style="color:blue;">@$1</span> ')
}

function switchtoStranger(id) {
  uni.navigateTo({
    url: `/pagesMy/stranger/stranger?id=${id}`
  })
}

function toggleInteraction(type) {
  if (type === 'like') {
    if (!ifLike.value) {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/like`,
        method: 'POST',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: (res) => {
          if (res.statusCode === 200) {
            post.value.likeCount = res.data.like_count
            ifLike.value = true
          } else {
            uni.showToast({
              title: '已经点过赞了~',
              icon: 'none',
              duration: 2000
            })
            ifLike.value = true
          }
        }
      })
    } else {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/like`,
        method: 'DELETE',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: () => {
          post.value.likeCount--
          ifLike.value = false
          uni.showToast({
            title: '点赞已取消',
            icon: 'none',
            duration: 2000
          })
        }
      })
    }
  } else if (type === 'favorite') {
    if (!ifFavourite.value) {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/favorite`,
        method: 'POST',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: (res) => {
          if (res.statusCode === 200) {
            post.value.favoriteCount = res.data.favorite_count
            ifFavourite.value = true
          } else {
            uni.showToast({
              title: '已经收藏了~',
              icon: 'none',
              duration: 2000
            })
            ifFavourite.value = true
          }
        }
      })
    } else {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/favorite`,
        method: 'DELETE',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: () => {
          post.value.favoriteCount--
          ifFavourite.value = false
          uni.showToast({
            title: '已取消收藏',
            icon: 'none',
            duration: 2000
          })
        }
      })
    }
  } else if (type === 'dislike') {
    if (!ifDislike.value) {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/dislike`,
        method: 'POST',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: () => {
          post.value.dislikeCount++
          ifDislike.value = true
        }
      })
    } else {
      uni.request({
        url: `${BASE_URL}/news/${PageId.value}/dislike`,
        method: 'DELETE',
        header: {
          'Content-type': 'application/json',
          Authorization: `Bearer ${jwtToken.value}`
        },
        success: () => {
          post.value.dislikeCount--
          ifDislike.value = false
        }
      })
    }
  }
}

function toggleCommentLike(index) {
  if (!comments[index].liked) {
    uni.request({
      url: `${BASE_URL}/news/${comments[index].id}/comment_like`,
      method: 'POST',
      header: {
        'Content-type': 'application/json',
        Authorization: `Bearer ${jwtToken.value}`
      },
      success: () => {
        comments[index].liked = true
        comments[index].likecount++
      }
    })
  } else {
    uni.request({
      url: `${BASE_URL}/news/${comments[index].id}/comment_like`,
      method: 'DELETE',
      header: {
        'Content-type': 'application/json',
        Authorization: `Bearer ${jwtToken.value}`
      },
      success: () => {
        comments[index].liked = false
        comments[index].likecount--
      }
    })
  }
}

function replyToComment(index) {
  replyingTo.value = index
  newReply.value = ''
}

function addReply(index) {
  if (!newReply.value.trim()) return
  const parentComment = comments[index]
  uni.request({
    url: `${BASE_URL}/news/comments`,
    method: 'POST',
    header: {
      'Content-type': 'application/json',
      Authorization: `Bearer ${jwtToken.value}`
    },
    data: {
      news_id: parseInt(PageId.value),
      content: newReply.value,
      is_reply: true,
      parent_id: parentComment.id
    },
    success: (res) => {
      if (res.statusCode === 201) {
        const newReplyComment = res.data.comment
        parentComment.replies.push({
          id: newReplyComment.id,
          text: newReplyComment.content,
          liked: newReplyComment.like_count > 0,
          authorName: uid.value,
          publish_time: formatPublishTime(newReplyComment.publish_time)
        })
        newReply.value = ''
        replyingTo.value = null
      } else {
        uni.showToast({
          title: '解决说明失败',
          icon: 'none',
          duration: 2000
        })
      }
    }
  })
}

function addComment() {
  if (!newComment.value.trim()) return
  uni.request({
    url: `${BASE_URL}/news/comments`,
    method: 'POST',
    header: {
      'Content-type': 'application/json',
      Authorization: `Bearer ${jwtToken.value}`
    },
    data: {
      news_id: parseInt(PageId.value),
      content: newComment.value,
      is_reply: false,
      parent_id: 0
    },
    success: (res) => {
      if (res.statusCode === 201) {
        const newCommentData = res.data.comment
        comments.push({
          id: newCommentData.id,
          text: newCommentData.content,
          liked: newCommentData.like_count > 0,
          likecount: 0,
          authorName: uid.value,
          authorAvatar: avatarSrc_sh.value,
          authorid: user_id.value,
          publish_time: formatPublishTime(newCommentData.publish_time),
          replies: [],
          showAllReplies: false
        })
        newComment.value = ''
      } else {
        uni.showToast({
          title: '添加说明失败',
          icon: 'none',
          duration: 2000
        })
      }
    }
  })
}

function formatPublishTime(publishTime) {
  const date = new Date(publishTime)
  const now = new Date()
  const isSameDay =
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  if (isSameDay) {
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `今天 ${hours}:${minutes}`
  } else {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
}

function handleCommentInput(e) {
  const val = e.detail.value
  if (!val.includes('@')) return
  // 这里也可以根据实际情况打开@用户的 popup
}

/* ----------------- Lifecycle ----------------- */
onLoad(async (options) => {
  const articleId = options.id
  PageId.value = articleId

  const details = await getArticleDetails(PageId.value)
  post.value = {
    id: PageId.value,
    authoravatar: details.author.avatar_url,
    authorname: details.author.nickname,
    authorid: details.author.id,
    savetime: details.upload_time,
    title: details.title,
    description: details.paragraphs[0]?.text || '',
    components: [],
    likeCount: details.like_count,
    shareCount: details.share_count,
    favoriteCount: details.favorite_count,
    followCount: 189,
    dislikeCount: details.dislike_count,
    viewCount: details.view_count
  }

  const totalItems = Math.max(details.paragraphs.length, details.images.length)
  for (let index = 1; index < totalItems; index++) {
    if (details.paragraphs[index]?.text) {
      post.value.components.push({
        id: index + 1,
        content: details.paragraphs[index].text,
        style: 'text'
      })
    }
    if (details.images[index]?.url) {
      post.value.components.push({
        id: index + 1,
        content: details.images[index].url,
        style: 'image',
        description: details.images[index].description || ''
      })
    }
  }

  if (details.comments && Array.isArray(details.comments)) {
    details.comments.forEach((comment) => {
      const formattedTime = formatPublishTime(comment.publish_time)
      const commentObj = {
        id: comment.id,
        text: comment.content,
        liked: comment.like_count > 0,
        likecount: comment.like_count,
        publish_time: formattedTime,
        authorName: comment.author.nickname,
        authorAvatar: comment.author.avatar_url,
        authorid: comment.author.id,
        replies: [],
        showAllReplies: false
      }
      if (comment.replies && Array.isArray(comment.replies)) {
        comment.replies.forEach((reply) => {
          const formattedReplyTime = formatPublishTime(reply.publish_time)
          commentObj.replies.push({
            id: reply.id,
            text: reply.content,
            liked: reply.like_count > 0,
            authorName: reply.author.nickname,
            publish_time: formattedReplyTime
          })
        })
      }
      comments.push(commentObj)
    })
  }
  // 记录浏览
  uni.request({
    url: `${BASE_URL}/news/${PageId.value}/view`,
    method: 'POST',
    header: {
      'Content-type': 'application/json',
      Authorization: `Bearer ${jwtToken.value}`
    }
  })

  // 获取初始点赞/收藏/不喜欢状态
  uni.request({
    url: `${BASE_URL}/news/${PageId.value}/status`,
    method: 'GET',
    header: {
      'Content-type': 'application/json',
      Authorization: `Bearer ${jwtToken.value}`
    },
    success: (res) => {
      ifLike.value = res.data.liked
      ifDislike.value = res.data.disliked
      ifFavourite.value = res.data.favorited
    }
  })
})

function getArticleDetails(id) {
  return new Promise((resolve) => {
    uni.request({
      url: `${BASE_URL}/news/details/news/${id}`,
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
.page-container {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding: 20rpx;
}

.author-section {
  background-color: #ffffff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.author-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.author-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background-color: #f0f0f0;
}

.author-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.content-section {
  background-color: #ffffff;
  padding: 30rpx;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.title-wrapper {
  margin-bottom: 40rpx;
  border-left: 8rpx solid #4caf50;
  padding-left: 20rpx;
}

.title {
  display: block;
  font-size: 44rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 16rpx;
  line-height: 1.4;
}

.description {
  display: block;
  font-size: 32rpx;
  color: #666666;
  line-height: 1.6;
}

.components {
  margin: 30rpx 0;
}

.component-item {
  margin-bottom: 30rpx;
}

.text-block {
  font-size: 30rpx;
  line-height: 1.6;
  color: #333333;
}

.image-block {
  margin: 20rpx 0;
}

.content-image {
  width: 100%;
  border-radius: 12rpx;
}

.image-caption {
  font-size: 26rpx;
  color: #999999;
  margin-top: 10rpx;
}

.metadata {
  display: flex;
  justify-content: space-between;
  margin: 30rpx 0;
  color: #999999;
  font-size: 26rpx;
}

.interaction-bar {
  display: flex;
  justify-content: space-around;
  padding: 20rpx 0;
  border-top: 2rpx solid #f0f0f0;
  border-bottom: 2rpx solid #f0f0f0;
}

.interaction-btn {
  display: flex;
  align-items: center;
  gap: 10rpx;
  padding: 10rpx 30rpx;
  border-radius: 30rpx;
  background-color: #f8f9fa;
  transition: all 0.3s ease;
}

.interaction-btn.active {
  background-color: #e6ffe6;
  color: #4caf50;
}

.interaction-icon {
  width: 36rpx;
  height: 36rpx;
}

.comments-section {
  margin-top: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}

.comment-card {
  padding: 20rpx;
  border-bottom: 2rpx solid #f0f0f0;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.commenter-info {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.commenter-avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
}

.commenter-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #333333;
}

.comment-time {
  font-size: 24rpx;
  color: #999999;
}

.comment-body {
  margin: 16rpx 0;
}

.comment-text {
  font-size: 28rpx;
  line-height: 1.6;
  color: #333333;
}

.comment-actions {
  display: flex;
  gap: 30rpx;
  margin-top: 16rpx;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: #666666;
}

.action-icon {
  width: 32rpx;
  height: 32rpx;
}

.reply-input {
  margin-top: 20rpx;
  display: flex;
  gap: 16rpx;
}

.reply-field {
  flex: 1;
  padding: 16rpx;
  border-radius: 8rpx;
  background-color: #f8f9fa;
  font-size: 28rpx;
}

.send-btn {
  padding: 12rpx 30rpx;
  background-color: #4caf50;
  color: #ffffff;
  border-radius: 8rpx;
  font-size: 28rpx;
}

.replies-list {
  margin-top: 20rpx;
  padding-left: 40rpx;
  border-left: 4rpx solid #f0f0f0;
}

.reply-item {
  margin-bottom: 16rpx;
}

.reply-author {
  font-weight: 500;
  color: #333333;
  font-size: 26rpx;
}

.reply-content {
  color: #666666;
  font-size: 26rpx;
  margin-left: 8rpx;
}

.reply-time {
  font-size: 24rpx;
  color: #999999;
  margin-left: 16rpx;
}

.show-more,
.show-more-comments {
  text-align: center;
  padding: 20rpx 0;
  color: #4caf50;
  font-size: 26rpx;
}

.add-comment {
  margin-top: 30rpx;
  display: flex;
  gap: 16rpx;
}

.comment-input {
  flex: 1;
  padding: 20rpx;
  border-radius: 8rpx;
  background-color: #f8f9fa;
  font-size: 28rpx;
}

.submit-btn {
  padding: 16rpx 40rpx;
  background-color: #4caf50;
  color: #ffffff;
  border-radius: 8rpx;
  font-size: 28rpx;
}

.mentions-list {
  background-color: #ffffff;
  border-radius: 16rpx 16rpx 0 0;
  padding: 20rpx;
  max-height: 400rpx;
  overflow-y: auto;
}

.mention-item {
  padding: 20rpx;
  border-bottom: 2rpx solid #f0f0f0;
  font-size: 28rpx;
}

.mention-item:last-child {
  border-bottom: none;
}
</style>