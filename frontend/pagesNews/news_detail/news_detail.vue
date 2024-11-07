<template>
  <view>
    <!-- Header Section -->
    <view class="header">
      <text class="back-button" @click="goBack">&larr;</text>
    </view>

    <!-- Content Section with Sidebar -->
    <view class="content-container">
      <!-- Main Content Section -->
      <view class="main-content">
        <view class="news-content">
          <view class="news-title">国际氢能联盟和麦肯锡联合发布《氢能洞察2024》</view>
          <view class="news-body">
            9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。
          </view>

          <!-- Interaction Buttons - Merged into News Content -->
          <view class="inline-interaction-buttons">
            <button @click="toggleInteraction('like')">👍 {{ likeText }}</button>
            <button @click="toggleInteraction('favorite')">⭐ {{ favoriteText }}</button>
            <button @click="toggleInteraction('share')">🔄 {{ shareText }}</button>
          </view>
        </view>

        <!-- Comments Section -->
        <view class="comments-section">
          <view class="comments-header">评论</view>
          <view id="comments-container">
            <view v-for="(comment, index) in comments" :key="index" class="comment">
              <view class="comment-content">
                <view class="comment-avatar"></view>
                <view>
                  <text class="comment-username">user_test:</text>
                  <text class="comment-text">{{ comment.text }}</text>
                </view>
              </view>
              <view class="comment-interactions">
                <button @click="toggleCommentLike(index)">👍 {{ comment.liked ? '已点赞' : '点赞' }}</button>
                <button @click="replyToComment(index)">💬 回复</button>
              </view>

              <!-- Reply Input Section -->
              <view v-if="replyingTo === index" class="add-reply">
                <input type="text" v-model="newReply" placeholder="回复..." />
                <button @click="addReply(index)">发送</button>
              </view>

              <!-- Replies Section -->
              <view v-if="comment.replies.length > 0" class="replies">
                <view v-for="(reply, replyIndex) in comment.replies" :key="replyIndex" class="reply">
                  <text class="comment-username">user_test:</text>
                  <text class="comment-text">{{ reply.text }}</text>
                </view>
              </view>
            </view>
          </view>
          <view class="add-comment">
            <input type="text" v-model="newComment" placeholder="发表评论..." />
            <button @click="addComment">评论</button>
          </view>
        </view>
      </view>

      <!-- Sidebar Section -->
      <view class="sidebar">
        <view class="sidebar-header">相关推荐</view>
        <view v-for="(recommendation, index) in recommendations" :key="index" class="recommendation-item">
          <image :src="recommendation.image" mode="widthFix" />
          <view class="recommendation-title">{{ recommendation.title }}</view>
          <view class="recommendation-info">{{ recommendation.info }}</view>
        </view>
      </view>
    </view>
  </view>
</template>


<script>
export default {
  data() {
    return {
      comments: [
        { text: "这篇文章非常有用！", liked: false, replies: [] },
      ],
      newComment: '',
      replyingTo: null,  // 当前正在回复的评论的索引
      newReply: '',      // 回复内容
      likeText: '点赞',
      favoriteText: '收藏',
      shareText: '分享',
      recommendations: [
        {
          image: "",
          title: "把自然讲给你听 | 什么是森林？",
          info: "阅读量: 1234 | 点赞量: 456"
        },
        {
          image: "",
          title: "全球氢能发展最新动态",
          info: "阅读量: 987 | 点赞量: 321"
        },
        {
          image: "",
          title: "如何做好垃圾分类",
          info: "阅读量: 789 | 点赞量: 123"
        }
      ]
    };
  },
  methods: {
    goBack() {
      uni.navigateBack();
    },
    toggleInteraction(type) {
      if (type === 'like') {
        this.likeText = this.likeText === '点赞' ? '已点赞' : '点赞';
      } else if (type === 'favorite') {
        this.favoriteText = this.favoriteText === '收藏' ? '已收藏' : '收藏';
      } else if (type === 'share') {
        this.shareText = this.shareText === '分享' ? '已分享' : '分享';
      }
    },
    toggleCommentLike(index) {
      this.$set(this.comments[index], 'liked', !this.comments[index].liked);
    },
    replyToComment(index) {
      this.replyingTo = index;
      this.newReply = ''; // 清空之前的回复内容
    },
    addReply(index) {
      if (this.newReply.trim()) {
        this.comments[index].replies.push({ text: this.newReply });
        this.newReply = '';
        this.replyingTo = null; // 回复完成后取消回复状态
      }
    },
    addComment() {
      if (this.newComment.trim()) {
        this.comments.push({ text: this.newComment, liked: false, replies: [] });
        this.newComment = '';
      }
    }
  }
};
</script>

<style scoped>
/* Body */
body {
  font-family: 'Arial', sans-serif;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

/* Header Section */
.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
}

.back-button {
  font-size: 24px;
  cursor: pointer;
}

/* Main Content Section */
.main-content {
  padding: 20px;
}

.news-content {
  padding: 20px;
  background-color: #ffffff;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.news-title {
  font-size: 26px;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
}

.news-body {
  font-size: 16px;
  line-height: 1.8;
  color: #555;
}

/* Inline Interaction Buttons - Combined to News Content */
.inline-interaction-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 10px;
  padding: 5px 0;
}

.inline-interaction-buttons button {
  border: none;
  background-color: transparent;
  cursor: pointer;
  font-size: 14px;
  color: #333;
  transition: color 0.3s;
}

.inline-interaction-buttons button:hover {
  color: #4caf50;
}

/* Comments Section */
.comments-section {
  padding: 20px;
  background-color: #ffffff;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.comment {
  border-bottom: 1px solid #e0e0e0;
  padding: 10px 0;
}

.comment:last-child {
  border-bottom: none;
}

.comment-content {
  display: flex;
  align-items: center;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  background-color: #ccc;
  border-radius: 50%;
  margin-right: 10px;
}

.comment-username {
  font-weight: bold;
  color: #4caf50;
}

.comment-text {
  font-size: 14px;
  color: #555;
}

.comment-interactions {
  display: flex;
  margin-top: 10px;
}

.comment-interactions button {
  border: none;
  background-color: transparent;
  cursor: pointer;
  font-size: 14px;
  color: #888;
  margin-right: 10px;
  transition: color 0.3s;
}

.comment-interactions button:hover {
  color: #4caf50;
}

.add-comment,
.add-reply {
  margin-top: 20px;
  display: flex;
}

.add-comment input,
.add-reply input {
  flex: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  margin-right: 10px;
  font-size: 14px;
}

.add-comment button,
.add-reply button {
  padding: 10px 20px;
  border: none;
  background-color: #4caf50;
  color: #ffffff;
  font-size: 14px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s;
}

.add-comment button:hover,
.add-reply button:hover {
  background-color: #45a049;
}

/* Replies Section */
.replies {
  margin-top: 10px;
  padding-left: 20px;
  border-left: 2px solid #e0e0e0;
}

.reply {
  margin-top: 10px;
}

/* Sidebar Section */
.sidebar {
  padding: 20px;
  background-color: #ffffff;
  margin: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.sidebar-header {
  font-size: 22px;
  margin-bottom: 15px;
  color: #333;
}

.recommendation-item {
  margin-bottom: 15px;
}

.recommendation-item image {
  width: 100%;
  height: auto;
  border-radius: 5px;
}

.recommendation-title {
  font-size: 16px;
  font-weight: bold;
  margin-top: 10px;
  color: #4caf50;
}

.recommendation-info {
  font-size: 14px;
  color: #555;
}

</style>
