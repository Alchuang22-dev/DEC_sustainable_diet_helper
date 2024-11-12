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
          <view class="news-title">{{ webTitle }}</view>
		  <view class="author-header">
		  	<view class="author-avatar"></view>
		    <text class="author-username">{{newsData[0].authorName}}</text>
		  </view>
          <view class="news-body">
            {{newsData[0].newsbody}}
          </view>

          <!-- Interaction Buttons - Merged into News Content -->
          <view class="inline-interaction-buttons">
            <button @click="toggleInteraction('like')">👍 {{ newsData[0].likeCount }}</button>
            <button @click="toggleInteraction('favorite')">⭐ {{ newsData[0].favoriteCount }}</button>
            <button @click="toggleInteraction('share')">🔄 {{ newsData[0].shareCount}}</button>
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
          <view class="recommendation-title" @click="goRecommend(recommendation.title, recommendation.form, recommendation.id)">{{ recommendation.title }}</view>
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
	  webTitle: '',
	  newsData: [],
      comments: [
        { text: "这篇文章非常有用！", liked: false, replies: [] },
      ],
      newComment: '',
      replyingTo: null,  // 当前正在回复的评论的索引
      newReply: '',      // 回复内容
      likeText: '点赞',
      favoriteText: '收藏',
      shareText: '分享',
      recommendations: []
    };
  },
  async created() {
    // 在组件创建时调用后端获取数据
    await this.fetchData();
  },
  onLoad(options) {
    if (options.title) {
      this.webTitle = decodeURIComponent(options.title);
    }
  },
  methods: {
	async fetchData() {
	    try {
	      // 模拟从后端获取数据
	      // 可以将此部分替换为实际的后端 API 调用，例如通过 axios:
	      // const response = await axios.get('your-api-endpoint');
	      
	      // 假设从后端获取的数据如下：
	      this.newsData = [{
			    id: 1,
			    form: 'news',
	        newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
			    imgsSrc: '',
			    tabs: ['环境保护','环保要闻'],
			    time: '2024-4-17',
			    newsName: '垃圾分类',
	        authorName: 'user_test',
	        authorAvatar: '',
	        newsinfo: '测试测试测试测试测试', 
			    newsbody: '9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。',
	        likeCount: 1001,
	        shareCount: 37,
	        favoriteCount: 897,
	        followCount: 189,
          dislikeCount: 100,
			    type: 'main'
	      },
	  	{
		    id: 2,
		    form: 'news',
	  	  newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
		    imgsSrc: '',
		    tabs: ['环境保护','环保要闻'],
		    time: '2024-4-17',
	  	  newsName: '把自然讲给你听',
	  	  authorName: '中野梓',
	  	  authorAvatar: '',
	  	  newsinfo: '测试测试测试测试测试', 
		    newsbody: '',
	  	  likeCount: 1001,
	  	  shareCount: 37,
	  	  favoriteCount: 897,
	  	  followCount: 189,
        dislikeCount: 100,
	  	  type: 'reco'
	  	}];
	  	this.recommendations = [
	  	  
	  	];
	  	this.newsData.forEach(news => this.convertnewsToRecommendation(news));
	    } catch (error) {
	      console.error('Error fetching data:', error);
	    }
	  },
	  convertnewsToRecommendation(news) {
	    if (news.type === 'reco') {
	      this.recommendations.push({
			id: news.id,
			src: news.newsSrc,
	        image: '',
	        title: news.authorName + ' | ' + news.newsName,
	        info: '阅读量: ' + news.followCount + ' | 点赞量: ' + news.likeCount,
			form: news.form,
	      });
	    }
	  },
    goBack() {
      uni.navigateBack();
    },
    toggleInteraction(type) {
      if (type === 'like') {
        this.newsData[0].likeCount++;
      } else if (type === 'favorite') {
        this.newsData[0].favoriteCount++;
      } else if (type === 'follow') {
        this.newsData[0].followCount++;
      } else if (type === 'share') {
        this.newsData[0].shareCount++;
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
	// 页面跳转方法
	goRecommend(title, form, id) {
	  setTimeout(() => {
	    if (form === 'news') {
	      // 图文页面跳转
	      uni.navigateTo({
	        url: `/pagesNews/news_detail/news_detail?title=${title}}`,
	      });
	    } else if(form === 'video'){
	      // 视频页面跳转
	      uni.navigateTo({
	        url: `/pagesNews/video_detail/video_detail?title=${name}`,
	      });
	    }
		else{
			uni.navigateTo({
			  url: `/pagesNews/web_detail/web_detail?url=${encodeURIComponent(id)}`,
			});
		}
	  }, 100); // 延迟 100 毫秒
	},
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
/*author part form video_detail*/
.author-avatar {
  width: 50px;
  height: 50px;
  background-color: #ccc;
  border-radius: 50%;
  margin-bottom: 10px;
}

.author-details {
  display: flex;
  flex-direction: column;
}

.author-header {
  display: flex;
  margin-bottom: 10px;
}

.author-username {
  font-weight: bold;
  margin-right: 20px;
}

</style>
