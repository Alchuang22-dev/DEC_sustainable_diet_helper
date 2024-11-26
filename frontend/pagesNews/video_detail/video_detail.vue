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
        :src="videoData[0].newsSrc"
        controls
        autoplay
        id=1
        @play="onPlay"
        @pause="onPause"
      >
      </video>
    </view>

    <!-- Tab Selection -->
    <!-- Tab Selection -->
    <view class="tab-selection">
      <view class="tab-container">
        <button @click="selectTab('简介')" :class="{ active: selectedTab === '简介' }">简介</button>
        <button @click="selectTab('评论')" :class="{ active: selectedTab === '评论' }">评论</button>
      </view>
    </view>

    <!-- Tab Content -->
    <view v-if="selectedTab === '简介'">
      <!-- Video Author Info and Interactions -->
      <view class="author-info">
        <view class="author-details">
          <view class="author-header">
			<view class="author-avatar"></view>
            <text class="author-username">{{videoData[0].authorName}}</text>
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
		  <view class="video_content">
			<view class="video_info"> {{videoData[0].newsinfo}}</view>
		  </view>
          <view class="author-interactions">
            <button @click="toggleInteraction('like')">👍 {{ formatCount(videoData[0].likeCount) }}</button>
            <button @click="toggleInteraction('favorite')">⭐ {{ formatCount(videoData[0].favoriteCount) }}</button>
            <button @click="toggleInteraction('share')">🔄 {{ formatCount(videoData[0].shareCount)}}</button>
			<button @click="toggleInteraction('dislike')" :style="{ color: ifDislike ? 'green' : 'black' }">👎 dis</button>
          </view>
        </view>
      </view>

      <!-- Sidebar Section (Recommendations) -->
      <view class="sidebar">
        <view class="sidebar-header">相关推荐</view>
        <view v-for="(recommendation, index) in recommendations" :key="index" class="recommendation-item">
          <image :src="recommendation.image" mode="widthFix" />
          <view class="recommendation-title" @click="goRecommend(recommendation.title, recommendation.form, recommendation.id)">{{ recommendation.title }}</view>
          <view class="recommendation-info">{{ recommendation.info }}</view>
        </view>
      </view>
    </view>

    <!-- Comments Section -->
    <view v-if="selectedTab === '评论'" class="comments-section">
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
        <input type="text" v-model="newComment" placeholder="才，才不是在等你的评论呢！" />
        <button @click="addComment">评论</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { onLoad } from "@dcloudio/uni-app";

const videoTitle = ref("");
const videoData = ref([]);
const comments = reactive([
  { text: "这个视频非常有用！", liked: false, replies: [] },
]);
const newComment = ref("");
const replyingTo = ref(null); // 当前正在回复的评论的索引
const selectedTab = ref("简介");
const newReply = ref(""); // 回复内容
const recommendations = ref([]);

const ifLike = ref(false);
const ifFavourite = ref(false);
const ifDislike = ref(false);
const ifShare = ref(false);
const ifFollowed = ref(false);

const fetchData = async () => {
  try {
    uni.request({
      url: "https://122.51.231.155/news/{id}", // 模拟的后端接口URL
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
        videoData.value = mockResponse.data;
        recommendations.value = [];
        videoData.value.forEach((video) => convertnewsToRecommendation(video));
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
        videoData.value = mockResponse.data;
        recommendations.value = [];
        videoData.value.forEach((video) => convertnewsToRecommendation(video));
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
    videoData.value = mockResponse.data;
    recommendations.value = [];
    videoData.value.forEach((video) => convertnewsToRecommendation(video));
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
			url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/like`,
			method: "POST",
			header: {
				"Content-type": "application/json",
			},
			data: {
				user_id: userId,
			},
			success: () => {
				videoData.value[0].likeCount++;
				ifLike.value = true;
			},
			fail: (err) => {
				console.error("Error liking news:", err);
			},
		});
	  }
	  else{
		    uni.request({
		    	url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/cancel_like`,
		    	method: "POST",
		    	header: {
		    		"Content-type": "application/json",
		    	},
		    	data: {
		    		user_id: userId,
		    	},
		    	success: () => {
		    		videoData.value[0].likeCount--;
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
		    url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      videoData.value[0].favoriteCount++;
			  ifFavourite.value = true;
		    },
		    fail: (err) => {
		      console.error("Error favoriting news:", err);
		    },
		  });
	  }
	  else{
		  uni.request({
		    url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      videoData.value[0].favoriteCount--;
		  	  ifFavourite.value = false;
		    },
		    fail: (err) => {
		      console.error("Error favoriting news:", err);
		    },
		  });
	  }
  } else if (type === "follow") {
	  if(ifFollowed.value === false){
		  //console.log("followed");
		  //向后端发送关注数++
		  ifFollowed.value = true;
	  }
	  else{
		  //向后端发送关注数--
		  ifFollowed.value = false;
	  }
  } else if (type === "share") {
    videoData.value[0].shareCount++;
  } else if (type === "dislike"){
	  if(ifDislike.value === false){
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		videoData.value[0].dislikeCount++;
		  		ifDislike.value = true;
		  	},
		  	fail: (err) => {
		  		console.error("Error liking news:", err);
		  	},
		  });
	  }
	  else{
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/cancel_dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		videoData.value[0].dislikeCount--;
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
      url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newReply.value,
        publish_time: "2024-11-05T12:35:00Z",
        user_id: userId,
        parent_id: 1,
        news_id: videoData.value[0].id,
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
      url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newComment.value,
        publish_time: "2024-11-05T12:30:00Z",
        user_id: userId,
        news_id: videoData.value[0].id,
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

const onPlay = () => {
    console.log('Video is playing');
};
const onPause = () => {
    console.log('Video is paused');
};
const selectTab = (tab) => {
    selectedTab.value = tab;
};

onMounted(async () => {
  await fetchData();
});

onLoad((options) => {
  if (options.title) {
    videoTitle.value = decodeURIComponent(options.title);
  }
});
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

.video-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  border-radius: 10px;
  overflow: hidden;
}

/* Tab Selection */
.tab-selection {
  display: flex;
  justify-content: center;
  margin: 20px;
}

.tab-container {
  display: flex;
  width: 100%;
  background-color: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
}

.tab-container button {
  flex: 1;
  border: none;
  padding: 10px;
  cursor: pointer;
  font-size: 16px;
  background-color: #ffffff;
  transition: background-color 0.3s;
}

.tab-container button.active {
  background-color: #4caf50;
  color: #ffffff;
  font-weight: bold;
}

/* Author Info */
.author-info {
  display: flex;
  flex-direction: column;
  padding: 20px;
  background-color: #ffffff;
  margin: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

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

.author-interactions {
  display: flex;
  gap: 10px;
}

/* Interaction Buttons */
.interaction-buttons {
  display: flex;
  justify-content: space-around;
  margin: 20px 20px 0 20px;
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

/* Info Section */
.video_content{
	padding: 10px;
	background-color: #ffffff;
	margin-bottom: 20px;
}

.video_info {
  font-size: 16px;
  line-height: 1.8;
  color: #555;
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
