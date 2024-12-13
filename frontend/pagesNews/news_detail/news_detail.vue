<template>
  <view class="container">
    <view class="author-header">
      <image :src="post.authoravatar" class="author-avatar"></image>
      <text class="author-username">{{ post.authorname }}</text>
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

    <!-- 文章标题和描述 -->
    <view class="title-container">
      <h1 class="article-title">{{ post.title }}</h1>
      <p class="article-description">{{ post.description }}</p>
    </view>

    <!-- 内容组件展示区 -->
    <view class="components-container">
      <view v-for="component in post.components" :key="component.id">
        <!-- 文本组件 -->
        <view v-if="component.style === 'text'" class="text-content">
          <p>{{ component.content }}</p>
        </view>

        <!-- 图片组件 -->
        <view v-if="component.style === 'image'" class="image-content">
          <image src="https://cdn.pixabay.com/photo/2017/04/09/07/25/honey-pomelo-2215031_1280.jpg" class="image"></image>
          <p class="image-description">{{ component.description }}</p>
        </view>
      </view>
    </view>
	
	<!-- Display the post time -->
	<view class="post-time">{{ post.savetime }}</view>

    <!-- 操作按钮 -->
    <view class="inline-interaction-buttons">
      <button @click="toggleInteraction('like')">👍 {{ formatCount(post.likeCount) }}</button>
      <button @click="toggleInteraction('favorite')">⭐ {{ formatCount(post.favoriteCount) }}</button>
      <button @click="toggleInteraction('share')">🔄 {{ formatCount(post.shareCount)}}</button>
      <button @click="toggleInteraction('dislike')" :style="{ color: ifDislike ? 'green' : 'black' }">👎 dis</button>
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
</template>

<script setup>
import { ref, reactive, onMounted, computed} from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { useNewsStore } from '@/stores/news_list';
import { useI18n } from 'vue-i18n';
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app';
import { storeToRefs } from 'pinia';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储

const newsStore = useNewsStore();
const userStore = useUserStore(); // 使用用户存储

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

const activeIndex = ref(null);
// 计算属性从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);
const uid = computed(() => userStore.user.nickName);
const avatarSrc = computed(() =>
    userStore.user.avatarUrl
        ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);

// 模拟传入的post数据
const post = ref({
  id: '1',
  authoravatar: 'https://example.com/avatar.jpg',
  authorname: 'John Doe',
  authorid: '123',
  savetime: '2024-12-13',
  title: 'Sample Article Title',
  description: 'This is a description of the article.',
  components: [
    { id: 1, content: 'This is a text component', style: 'text' },
    { id: 2, content: 'https://cdn.pixabay.com/photo/2017/04/09/07/25/honey-pomelo-2215031_1280.jpg', style: 'image', description: 'This is an image' },
  ],
  likeCount: 1001,
  shareCount: 37,
  favoriteCount: 897,
  followCount: 189,
  dislikeCount: 199,
  type: 'main',
});

//转换数字
const formatCount = (count) => {
  return count < 10000 ? count : (count / 1000).toFixed(1) + 'k';
};

//处理操作
const toggleInteraction = (type) => {
  const userId = uni.getStorageSync('UserId');
  if (type === "like") {
		if(ifLike.value === false) {
			uni.request({
			url: `http://122.51.231.155:8080/news/${post.id}/like`,
			method: "POST",
			header: {
				"Content-type": "application/json",
			},
			data: {
				user_id: userId,
			},
			success: () => {
				post.likeCount++;
				ifLike.value = true;
			},
			fail: (err) => {
				console.error("Error liking news:", err);
			},
		});
	  }
	  else{
		    uni.request({
		    	url: `http://122.51.231.155:8080/news/${post.id}/cancel_like`,
		    	method: "POST",
		    	header: {
		    		"Content-type": "application/json",
		    	},
		    	data: {
		    		user_id: userId,
		    	},
		    	success: () => {
		    		post.likeCount--;
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
		    url: `http://122.51.231.155:8080/news/${post.id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      post.favoriteCount++;
			  ifFavourite.value = true;
		    },
		    fail: (err) => {
		      console.error("Error favoriting news:", err);
		    },
		  });
	  }
	  else{
		  uni.request({
		    url: `http://122.51.231.155:8080/news/${post.id}/favourite`,
		    method: "POST",
		    header: {
		      "Content-type": "application/json",
		    },
		    data: {
		      user_id: userId,
		    },
		    success: () => {
		      post.favoriteCount--;
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
        url: `http://122.51.231.155:8080/user/${uid.value}/follow`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          target_id: post.authorName, // 示例参数
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
        url: `http://122.51.231.155:8080/user/${uid.value}/unfollow`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          target_id: post.authorName, // 示例参数
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
    post.shareCount++;
  } else if (type === "dislike"){
	  if(ifDislike.value === false){
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${post.id}/dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		post.dislikeCount++;
		  		ifDislike.value = true;
		  	},
		  	fail: (err) => {
		  		console.error("Error liking news:", err);
		  	},
		  });
	  }
	  else{
		  uni.request({
		  	url: `http://122.51.231.155:8080/news/${post.id}/cancel_dislike`,
		  	method: "POST",
		  	header: {
		  		"Content-type": "application/json",
		  	},
		  	data: {
		  		user_id: userId,
		  	},
		  	success: () => {
		  		post.dislikeCount--;
		  		ifDislike.value = false;
		  	},
		  	fail: (err) => {
		  		console.error("Error liking news:", err);
		  	},
		  });
	  }
  }
};

//评论相关方法
const toggleCommentLike = (index) => {
  comments[index].liked = !comments[index].liked;
};

const replyToComment = (index) => {
  replyingTo.value = index;
  newReply.value = ""; // 清空之前的回复内容
};

const addReply = (index) => {
  if (newReply.value.trim()) {
    uni.request({
      url: `http://122.51.231.155:8080/news/${post.id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newReply.value,
        publish_time: "2024-11-05T12:35:00Z",
        user_id: uid.value,
        parent_id: 1,
        news_id: post.id,
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
    uni.request({
      url: `http://122.51.231.155:8080/news/${post.id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newComment.value,
        publish_time: "2024-11-05T12:30:00Z",
        user_id: uid.value,
        news_id: post.id,
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

</script>

<style scoped>
.container {
  padding: 20px;
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

/* Title and Description styles */
.article-title {
  font-family: 'Arial', sans-serif;
  font-size: 26px;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
}

.article-description {
  font-family: 'Verdana', sans-serif;
  font-size: 18px;
  color: #666;
}

/*关注按钮*/
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

/* 交互按钮 */
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


/* Content Section */
.components-container {
  margin-top: 20px;
  margin-bottom: 20px;
}

.text-content p {
  margin-top: 10px; 
  font-size: 16px;
  line-height: 1.5;
  margin-bottom: 10px; /* Add space between text components */
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

.extra-info {
  font-size: 14px;
  color: #777;
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

/* Post Time */
.post-time {
  font-size: 14px;
  color: #888;
  text-align: right;
  margin-top: 20px;
}
</style>
