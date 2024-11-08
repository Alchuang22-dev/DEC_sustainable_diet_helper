<template>
  <view>
    <!-- Header Section -->
    <view class="header">
      <text class="back-button" @click="goBack">&larr;</text>
    </view>

    <!-- Video Header -->
    <view class="video-header">
      <text class="content">{{videoTitle}}</text>
    </view>

    <!-- Video Content -->
    <view class="video-content">
      <video
        class="video-container"
        :src="videoSrc"
        controls
        autoplay
        id="video"
        @play="onPlay"
        @pause="onPause"
      >
      </video>
    </view>

    <!-- Interaction Buttons -->
    <view class="interaction-buttons">
      <button @click="toggleInteraction('like')">👍 {{ likeText }}</button>
      <button @click="toggleInteraction('favorite')">⭐ {{ favoriteText }}</button>
      <button @click="toggleInteraction('share')">🔄 {{ shareText }}</button>
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
</template>

<script>
export default {
  data() {
    return {
	  videoTitle: '',
      videoSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
      likeText: '点赞',
      favoriteText: '收藏',
      shareText: '分享',
      comments: [
        { text: '这个视频非常有用！', liked: false, replies: [] },
      ],
      newComment: '',
      replyingTo: null,
      newReply: '',
      recommendations: [
        {
          image: '',
          title: '把自然讲给你听 | 什么是森林？',
          info: '阅读量: 1234 | 点赞量: 456',
        },
        {
          image: '',
          title: '全球氢能发展最新动态',
          info: '阅读量: 987 | 点赞量: 321',
        },
        {
          image: '',
          title: '如何做好垃圾分类',
          info: '阅读量: 789 | 点赞量: 123',
        },
      ],
    };
  },
  onLoad(options) {
    if (options.title) {
      this.videoTitle = decodeURIComponent(options.title);
    }
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
    },
    onPlay() {
      console.log('Video is playing');
    },
    onPause() {
      console.log('Video is paused');
    },
  },
};
</script>

<style scoped>
/* Header Section */
.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.back-button {
  font-size: 24px;
  cursor: pointer;
}

/* Video Header */
.video-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
}

.content {
  font-size: 18px;
  font-weight: bold;
}

/* Video Content */
.video-content {
  padding: 20px;
  background-color: #ffffff;
  margin: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.video-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  border-radius: 10px;
  overflow: hidden;
}

/* Interaction Buttons */
.interaction-buttons {
  display: flex;
  justify-content: space-around;
  margin: 20px;
  padding: 10px;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.interaction-buttons button {
  border: none;
  background-color: #ffffff;
  padding: 10px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.interaction-buttons button:hover {
  background-color: #f0f0f0;
}

/* Comments Section */
.comments-section {
  padding: 20px;
  background-color: #ffffff;
  margin: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.comments-header {
  font-size: 22px;
  margin-bottom: 15px;
  color: #333;
}

.comment {
  border-bottom: 1px solid #e0e0e0;
  padding: 10px 0;
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
