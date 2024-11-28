<template>
  <view>
    <!-- 条件渲染，确保数据加载完成后才渲染内容 -->
    <view v-if="newsData.length > 0" class="content-container">
      <!-- Main Content Section -->
      <view class="main-content">
        <view class="news-content">
          <view class="news-title">{{ webTitle }}</view>
          <view class="author-header">
            <view class="author-avatar"></view>
            <text class="author-username">{{ newsData[0].authorName }}</text>
            <button
              class="stable-button"
              @click="toggleInteraction('follow')"
              :style="{ 
                color: ifFollowed ? 'black' : 'white', 
                backgroundColor: ifFollowed ? 'lightgrey' : 'green' 
              }"
            >
              {{ ifFollowed ? '已关注' : '关注' }}
            </button>
          </view>
          <view class="news-body">
            {{ newsData[0].newsbody }}
          </view>

          <!-- Interaction Buttons -->
          <view class="inline-interaction-buttons">
            <button @click="toggleInteraction('like')">👍 {{ formatCount(newsData[0].likeCount) }}</button>
            <button @click="toggleInteraction('favorite')">⭐ {{ formatCount(newsData[0].favoriteCount) }}</button>
            <button @click="toggleInteraction('share')">🔄 {{ formatCount(newsData[0].shareCount)}}</button>
            <button @click="toggleInteraction('dislike')" :style="{ color: ifDislike ? 'green' : 'black' }">👎 dis</button>
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
          <view class="recommendation-title" @click="goRecommend(recommendation.title, recommendation.form, recommendation.id)">
            {{ recommendation.title }}
          </view>
          <view class="recommendation-info">{{ recommendation.info }}</view>
        </view>
      </view>
    </view>

    <!-- Loading State -->
       <!-- Loading State -->
       <view v-else-if="loadingError" class="loading-container">
         <text>加载失败，请稍后重试</text>
       </view>
       <view v-else class="loading-container">
         <text>加载中...</text>
       </view>
  </view>
</template>


<script setup>
import { ref, reactive, onMounted } from "vue";
import { onLoad } from "@dcloudio/uni-app";

const webTitle = ref("");
const newsData = ref([]);
const comments = reactive([
  { text: "这篇文章非常有用！", liked: false, replies: [] },
]);
const newComment = ref("");
const replyingTo = ref(null); // 当前正在回复的评论的索引
const newReply = ref(""); // 回复内容
const recommendations = ref([]);
const loadingError = ref(false); // 加载错误标志
const timeout = 15000; // 超时时间：15秒

const ifLike = ref(false);
const ifFavourite = ref(false);
const ifDislike = ref(false);
const ifShare = ref(false);
const ifFollowed = ref(false);

const fetchData = async () => {
	const timer = setTimeout(() => {
	    loadingError.value = true; // 超时后显示加载失败
	  }, timeout);
  try {
    uni.request({
      url: "https://122.51.231.155/news/1", // 模拟的后端接口URL
      method: "GET",
      data: {
        id: 1,
      },
      success: (res) => {
        const mockResponse = {
          data: [
            {
              id: 1,
              form: "news",
              newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
              imgsSrc: "",
              tabs: ["环境保护", "环保要闻"],
              time: "2024-4-17",
              newsName: "垃圾分类",
              authorName: "user_test",
              authorAvatar: "",
              newsinfo: "测试测试测试测试测试",
              newsbody:
                "9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。",
              likeCount: 10010,
              shareCount: 37,
              favoriteCount: 897,
              followCount: 189,
              dislikeCount: 100,
              type: "main",
            },
            {
              id: 2,
              form: "news",
              newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
              imgsSrc: "",
              tabs: ["环境保护", "环保要闻"],
              time: "2024-4-17",
              newsName: "把自然讲给你听",
              authorName: "中野梓",
              authorAvatar: "",
              newsinfo: "测试测试测试测试测试",
              newsbody: "",
              likeCount: 1001,
              shareCount: 37,
              favoriteCount: 897,
              followCount: 189,
              dislikeCount: 100,
              type: "reco",
            },
          ],
        };
        newsData.value = mockResponse.data;
        recommendations.value = [];
        newsData.value.forEach((video) => convertnewsToRecommendation(video));
      },
      fail: (err) => {
        const mockResponse = {
          data: [
            {
              id: 1,
              form: "news",
              newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
              imgsSrc: "",
              tabs: ["环境保护", "环保要闻"],
              time: "2024-4-17",
              newsName: "垃圾分类",
              authorName: "user_test",
              authorAvatar: "",
              newsinfo: "测试测试测试测试测试",
              newsbody:
                "9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。",
              likeCount: 10010,
              shareCount: 37,
              favoriteCount: 897,
              followCount: 189,
              dislikeCount: 100,
              type: "main",
            },
            {
              id: 2,
              form: "news",
              newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
              imgsSrc: "",
              tabs: ["环境保护", "环保要闻"],
              time: "2024-4-17",
              newsName: "把自然讲给你听",
              authorName: "中野梓",
              authorAvatar: "",
              newsinfo: "测试测试测试测试测试",
              newsbody: "",
              likeCount: 1001,
              shareCount: 37,
              favoriteCount: 897,
              followCount: 189,
              dislikeCount: 100,
              type: "reco",
            },
          ],
        };
        newsData.value = mockResponse.data;
        recommendations.value = [];
        newsData.value.forEach((video) => convertnewsToRecommendation(video));
      },
    });
  } catch (error) {
    const mockResponse = {
      data: [
        {
          id: 1,
          form: "news",
          newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
          imgsSrc: "",
          tabs: ["环境保护", "环保要闻"],
          time: "2024-4-17",
          newsName: "垃圾分类",
          authorName: "user_test",
          authorAvatar: "",
          newsinfo: "测试测试测试测试测试",
          newsbody:
            "9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。",
          likeCount: 10010,
          shareCount: 37,
          favoriteCount: 897,
          followCount: 189,
          dislikeCount: 100,
          type: "main",
        },
        {
          id: 2,
          form: "news",
          newsSrc: "http://vjs.zencdn.net/v/oceans.mp4",
          imgsSrc: "",
          tabs: ["环境保护", "环保要闻"],
          time: "2024-4-17",
          newsName: "把自然讲给你听",
          authorName: "中野梓",
          authorAvatar: "",
          newsinfo: "测试测试测试测试测试",
          newsbody: "",
          likeCount: 1001,
          shareCount: 37,
          favoriteCount: 897,
          followCount: 189,
          dislikeCount: 100,
          type: "reco",
        },
      ],
    };
    newsData.value = mockResponse.data;
    recommendations.value = [];
    newsData.value.forEach((video) => convertnewsToRecommendation(video));
  }
};

const formatCount = (count) => {
  return count < 10000 ? count : (count / 1000).toFixed(1) + 'k';
};

const convertnewsToRecommendation = (news) => {
  if (news.type === "reco") {
    recommendations.value.push({
      id: news.id,
      src: news.newsSrc,
      image: "",
      title: news.authorName + " | " + news.newsName,
      info: "阅读量: " + news.followCount + " | 点赞量: " + news.likeCount,
      form: news.form,
    });
  }
};

const goBack = () => {
  uni.navigateBack();
};

const toggleInteraction = (type) => {
  const userId = uni.getStorageSync('UserId');
  if (type === "like") {
		if(ifLike.value === false) {
			uni.request({
			url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/like`,
			method: "POST",
			header: {
				"Content-type": "application/json",
			},
			data: {
				user_id: userId,
			},
			success: () => {
				newsData.value[0].likeCount++;
				ifLike.value = true;
			},
			fail: (err) => {
				console.error("Error liking news:", err);
			},
		});
	  }
	  else{
		    uni.request({
		    	url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/cancel_like`,
		    	method: "POST",
		    	header: {
		    		"Content-type": "application/json",
		    	},
		    	data: {
		    		user_id: userId,
		    	},
		    	success: () => {
		    		newsData.value[0].likeCount--;
		    		ifLike.value = false;
		    	},
		    	fail: (err) => {
		    		console.error("Error Cancel liking news:", err);
		    	},
		    });
	  }
  } else if (type === "favorite") {
	  if(ifFavourite.value === false){
		  uni.request({
		    url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      newsData.value[0].favoriteCount++;
			  ifFavourite.value = true;
		    },
		    fail: (err) => {
		      console.error("Error favoriting news:", err);
		    },
		  });
	  }
	  else{
		  uni.request({
		    url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      newsData.value[0].favoriteCount--;
		  	  ifFavourite.value = false;
		    },
		    fail: (err) => {
		      console.error("Error favoriting news:", err);
		    },
		  });
	  }
  } else if (type === "follow") {
    if (ifFollowed.value === false) {
      // 向后端发送关注请求
      uni.request({
        url: `http://122.51.231.155:8080/user/${userId.value}/follow`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          target_id: newsData.value[0].authorName, // 示例参数
        },
        success: () => {
          ifFollowed.value = true;
        },
        fail: (err) => {
          console.error("Error following user:", err);
        },
      });
    } else {
      // 向后端发送取消关注请求
      uni.request({
        url: `http://122.51.231.155:8080/user/${userId.value}/unfollow`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          target_id: newsData.value[0].authorName, // 示例参数
        },
        success: () => {
          ifFollowed.value = false;
        },
        fail: (err) => {
          console.error("Error unfollowing user:", err);
        },
      });
    }
  } else if (type === "share") {
    newsData.value[0].shareCount++;
  } else if (type === "dislike"){
	  if(ifDislike.value === false){
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		newsData.value[0].dislikeCount++;
		  		ifDislike.value = true;
		  	},
		  	fail: (err) => {
		  		console.error("Error liking news:", err);
		  	},
		  });
	  }
	  else{
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/cancel_dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		newsData.value[0].dislikeCount--;
		  		ifDislike.value = false;
		  	},
		  	fail: (err) => {
		  		console.error("Error liking news:", err);
		  	},
		  });
	  }
  }
};

const toggleCommentLike = (index) => {
  comments[index].liked = !comments[index].liked;
};

const replyToComment = (index) => {
  replyingTo.value = index;
  newReply.value = ""; // 清空之前的回复内容
};

const addReply = (index) => {
  if (newReply.value.trim()) {
    const userId = uni.getStorageSync('UserId');
    uni.request({
      url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newReply.value,
        publish_time: "2024-11-05T12:35:00Z",
        user_id: userId,
        parent_id: 1,
        news_id: newsData.value[0].id,
        is_reply: true,
      },
      success: () => {
        comments[index].replies.push({ text: newReply.value });
        newReply.value = "";
        replyingTo.value = null; // 回复完成后取消回复状态
      },
      fail: (err) => {
        console.error("Error adding reply:", err);
      },
    });
  }
};

const addComment = () => {
  if (newComment.value.trim()) {
    const userId = uni.getStorageSync('UserId');
    uni.request({
      url: `http://122.51.231.155:8080/news/${newsData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newComment.value,
        publish_time: "2024-11-05T12:30:00Z",
        user_id: userId,
        news_id: newsData.value[0].id,
        is_reply: false,
        is_liked: false,
      },
      success: () => {
        comments.push({ text: newComment.value, liked: false, replies: [] });
        newComment.value = "";
      },
      fail: (err) => {
        console.error("Error adding comment:", err);
      },
    });
  }
};

// 页面跳转方法
const goRecommend = (title, form, id) => {
  setTimeout(() => {
    if (form === "news") {
      // 图文页面跳转
      uni.navigateTo({
        url: `/pagesNews/news_detail/news_detail?title=${title}}`,
      });
    } else if (form === "video") {
      // 视频页面跳转
      uni.navigateTo({
        url: `/pagesNews/video_detail/video_detail?title=${title}`,
      });
    } else {
      uni.navigateTo({
        url: `/pagesNews/web_detail/web_detail?url=${encodeURIComponent(id)}`,
      });
    }
  }, 100); // 延迟 100 毫秒
};

onMounted(async () => {
  await fetchData();
});

onLoad((options) => {
  if (options.title) {
    webTitle.value = decodeURIComponent(options.title);
  }
});
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

.stable-button {
  width: 100px; /* 固定宽度 */
  height: 40px; /* 固定高度 */
  display: inline-flex; /* 使内容居中对齐 */
  align-items: center; /* 垂直居中 */
  justify-content: center; /* 水平居中 */
  border: 1px solid #ccc; /* 可选：边框样式 */
  border-radius: 5px; /* 可选：圆角 */
  background-color: #f5f5f5; /* 可选：背景颜色 */
  cursor: pointer; /* 鼠标悬浮时的样式 */
  overflow: hidden; /* 防止内容溢出 */
  text-align: center; /* 文本居中 */
  font-size: 14px; /* 可选：字体大小 */
  box-sizing: border-box; /* 包括 padding 和 border */
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
