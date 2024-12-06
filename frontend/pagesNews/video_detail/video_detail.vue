<template>
  <view v-if="videoData.length > 0" class="container">
    <!-- 主滚动视图 -->
    <scroll-view
      class="main-scroll"
      scroll-y
      @scroll="handleScroll"
      :scroll-with-animation="true"
      :style="{ height: '100vh' }"
    >
      <!-- Header 部分 -->

      <!-- 视频标题 -->
      <view :class="['video-header', { 'video-header-shrink': isTabSticky }]">
        <text class="content">{{ videoTitle }}</text>
      </view>

      <!-- 视频内容 -->
      <view ref="videoContent" :class="['video-content', { 'video-content-shrink': isTabSticky }]">
        <video
          ref="videoPlayer"
          class="video-container"
          :src="videoData[0].newsSrc"
          controls
          autoplay
          @play="onPlay"
          @pause="onPause"
          @waiting="onWaiting"
          @canplay="onCanPlay"
          @canplaythrough="onCanPlayThrough"
          @error="onError"
          @seeking="onSeeking"
          @seeked="onSeeked"
          @playing="onPlaying"
          @fullscreenchange="onFullScreenChange"
          @enterpictureinpicture="onEnterPiP"
          @leavepictureinpicture="onLeavePiP"
        >
        </video>
        <!-- 加载动画覆盖层 -->
        <view v-if="loading" class="loading-overlay">
          <image src="@/static/loading.gif" class="loading-spinner"></image>
        </view>
      </view>

      <!-- Tab 选择部分 -->
      <view class="tab-selection" :class="{ 'tab-selection-sticky': isTabSticky }">
        <view class="tab-container">
          <button @click="selectTab('简介')" :class="{ active: selectedTab === '简介' }">简介</button>
          <button @click="selectTab('评论')" :class="{ active: selectedTab === '评论' }">评论</button>
        </view>
      </view>

      <!-- Tab 内容部分 -->
      <view class="tab-content">
        <template v-if="selectedTab === '简介'">
          <!-- 简介内容 -->
          <view class="author-info">
            <view class="author-details">
              <view class="author-header">
                <view class="author-avatar"></view>
                <text class="author-username">{{ videoData[0].authorName }}</text>
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
                <view class="video_info"> 
                  {{ videoData[0].newsinfo }}
                  <!-- 为测试滚动添加更多内容 -->
                  <text v-for="n in 20" :key="n"> 这是更多的简介内容，用于测试滚动功能。 </text>
                </view>
              </view>
              <view class="author-interactions">
                <button @click="toggleInteraction('like')">👍 {{ formatCount(videoData[0].likeCount) }}</button>
                <button @click="toggleInteraction('favorite')">⭐ {{ formatCount(videoData[0].favoriteCount) }}</button>
                <button @click="toggleInteraction('share')">🔄 {{ formatCount(videoData[0].shareCount)}}</button>
                <button @click="toggleInteraction('dislike')" :style="{ color: ifDislike ? 'green' : 'black' }">👎 dis</button>
              </view>
            </view>
          </view>

          <!-- 相关推荐部分 -->
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
        </template>

        <template v-else-if="selectedTab === '评论'">
          <!-- 评论内容 -->
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

                <!-- 回复输入部分 -->
                <view v-if="replyingTo === index" class="add-reply">
                  <input type="text" v-model="newReply" placeholder="回复..." />
                  <button @click="addReply(index)">发送</button>
                </view>

                <!-- 回复内容部分 -->
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
        </template>
      </view>
    </scroll-view>
  </view>
  <!-- Loading State -->
  <view v-else-if="loadingError" class="loading-container">
    <text>加载失败，请稍后重试</text>
  </view>
  <view v-else class="loading-container">
    <text>加载中...</text>
  </view>
</template>


<script setup>
import { ref, reactive, onMounted } from "vue";
import { onLoad } from "@dcloudio/uni-app";

// 视频相关数据
const videoTitle = ref("");
const videoData = ref([]);

// 评论相关数据
const comments = reactive([
  { text: "这个视频非常有用！", liked: false, replies: [] },
]);
const newComment = ref("");
const replyingTo = ref(null); // 当前正在回复的评论的索引
const newReply = ref(""); // 回复内容
const loadingError = ref(false); // 加载错误标志
const timeout = 15000; // 超时时间：15秒

// Tab 相关
const selectedTab = ref("简介");

// 推荐相关数据
const recommendations = ref([]);

// 交互状态
const ifLike = ref(false);
const ifFavourite = ref(false);
const ifDislike = ref(false);
const ifShare = ref(false);
const ifFollowed = ref(false);

// 加载状态
const loading = ref(false); // 初始状态为 false，避免在数据未加载前显示

// 视频播放器引用
const videoPlayer = ref(null);

// Sticky 状态
const isTabSticky = ref(false);

// 动态视频高度
const videoHeight = ref(0);

// 引用视频内容的 DOM 元素
const videoContent = ref(null);

// 获取用户ID（假设已存储在本地存储中）
const userId = ref(uni.getStorageSync('UserId'));

// 数据获取函数
const fetchData = async () => {
	const timer = setTimeout(() => {
	    loadingError.value = true; // 超时后显示加载失败
	  }, timeout);
  try {
    uni.request({
      url: "https://122.51.231.155/news/{id}", // 后端接口URL
      method: "GET",
      data: {
        id: 1,
      },
      success: (res) => {
        // 使用模拟数据
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
        console.log("视频数据加载成功:", videoData.value);
        recommendations.value = [];
        videoData.value.forEach((video) => convertnewsToRecommendation(video));
      },
      fail: (err) => {
        console.error("请求失败，使用模拟数据:", err);
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
        console.log("请求失败，使用模拟数据:", videoData.value);
        recommendations.value = [];
        videoData.value.forEach((video) => convertnewsToRecommendation(video));
      },
    });
  } catch (error) {
    console.error("异常错误，使用模拟数据:", error);
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
    console.log("异常错误，使用模拟数据:", videoData.value);
    recommendations.value = [];
    videoData.value.forEach((video) => convertnewsToRecommendation(video));
  }
};

// 格式化数字显示
const formatCount = (count) => {
  return count < 10000 ? count : (count / 1000).toFixed(1) + 'k';
};

// 将新闻数据转换为推荐项
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

// 返回上一页
const goBack = () => {
  uni.navigateBack();
};

// 切换 Tab
const selectTab = (tab) => {
  selectedTab.value = tab;
};

// 切换交互（点赞、收藏、关注等）
const toggleInteraction = (type) => {
  if (type === "like") {
    if (ifLike.value === false) {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/like`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
        },
        success: () => {
          videoData.value[0].likeCount++;
          ifLike.value = true;
        },
        fail: (err) => {
          console.error("Error liking news:", err);
        },
      });
    } else {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/cancel_like`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
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
    if (ifFavourite.value === false) {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/favourite`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
        },
        success: () => {
          videoData.value[0].favoriteCount++;
          ifFavourite.value = true;
        },
        fail: (err) => {
          console.error("Error favoriting news:", err);
        },
      });
    } else {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/favourite`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
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
    if (ifFollowed.value === false) {
      // 向后端发送关注请求
      uni.request({
        url: `http://122.51.231.155:8080/user/${userId.value}/follow`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          target_id: videoData.value[0].authorName, // 示例参数
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
          target_id: videoData.value[0].authorName, // 示例参数
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
    videoData.value[0].shareCount++;
    // 这里可以添加分享功能的实现
  } else if (type === "dislike") {
    if (ifDislike.value === false) {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/dislike`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
        },
        success: () => {
          videoData.value[0].dislikeCount++;
          ifDislike.value = true;
        },
        fail: (err) => {
          console.error("Error disliking news:", err);
        },
      });
    } else {
      uni.request({
        url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/cancel_dislike`,
        method: "POST",
        header: {
          "Content-type": "application/json",
        },
        data: {
          user_id: userId.value,
        },
        success: () => {
          videoData.value[0].dislikeCount--;
          ifDislike.value = false;
        },
        fail: (err) => {
          console.error("Error canceling dislike:", err);
        },
      });
    }
  }
};

// 评论相关函数
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
      url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newReply.value,
        publish_time: new Date().toISOString(),
        user_id: userId.value,
        parent_id: index + 1, // 示例：假设 parent_id 是评论的索引加1
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
    uni.request({
      url: `http://122.51.231.155:8080/news/${videoData.value[0].id}/comment`,
      method: "POST",
      header: {
        "Content-type": "application/json",
      },
      data: {
        content: newComment.value,
        publish_time: new Date().toISOString(),
        user_id: userId.value,
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

// 跳转到推荐内容
const goRecommend = (title, form, id) => {
  setTimeout(() => {
    if (form === "news") {
      // 图文页面跳转
      uni.navigateTo({
        url: `/pagesNews/news_detail/news_detail?title=${encodeURIComponent(title)}`,
      });
    } else if (form === "video") {
      // 视频页面跳转
      uni.navigateTo({
        url: `/pagesNews/video_detail/video_detail?title=${encodeURIComponent(title)}`,
      });
    } else {
      uni.navigateTo({
        url: `/pagesNews/web_detail/web_detail?url=${encodeURIComponent(id)}`,
      });
    }
  }, 100); // 延迟 100 毫秒
};

// 全屏切换函数（可选）
const toggleFullScreen = () => {
  if (videoPlayer.value) {
    const videoElement = videoPlayer.value.$el.querySelector('video'); // 获取视频 DOM 元素
    if (videoElement.requestFullscreen) {
      videoElement.requestFullscreen();
    } else if (videoElement.webkitRequestFullscreen) { /* Safari */
      videoElement.webkitRequestFullscreen();
    } else if (videoElement.msRequestFullscreen) { /* IE11 */
      videoElement.msRequestFullscreen();
    }
  }
};

// 视频事件处理函数
const onPlay = () => {
  console.log('Video is playing');
  uni.hideLoading();
};

const onPause = () => {
  console.log('Video is paused');
};

const onWaiting = () => {
  console.log('Video is waiting to buffer');
  uni.showLoading({
    title: '加载中...',
    mask: true,
  });
};

const onCanPlay = () => {
  console.log('Video can play');
  uni.hideLoading();
};

const onCanPlayThrough = () => {
  console.log('Video can play through without stopping');
  uni.hideLoading();
};

const onError = () => {
  console.log('Video failed to load');
  uni.hideLoading();
  uni.showToast({
    title: '视频加载失败',
    icon: 'none',
  });
};

const onSeeking = () => {
  console.log('User is seeking');
  uni.showLoading({
    title: '加载中...',
    mask: true,
  });
};

const onSeeked = () => {
  console.log('User has finished seeking');
  setTimeout(() => {
    uni.hideLoading();
  }, 500);
};

const onFullScreenChange = (e) => {
  console.log('全屏状态改变', e);
};

const onEnterPiP = () => {
  console.log('进入画中画模式');
};

const onLeavePiP = () => {
  console.log('离开画中画模式');
};

// 页面挂载时获取数据和视频高度
onMounted(async () => {
  await fetchData();

  // 获取视频容器的高度
  uni.createSelectorQuery()
    .select('.video-content')
    .boundingClientRect((rect) => {
      if (rect) {
        videoHeight.value = rect.height; // 计算视频内容的高度
        console.log("视频内容高度:", videoHeight.value);
      }
    })
    .exec();
});

// 页面加载时获取标题
onLoad((options) => {
  if (options.title) {
    videoTitle.value = decodeURIComponent(options.title);
  }
});

// 处理滚动事件
const handleScroll = (e) => {
  const scrollTop = e.detail.scrollTop;
  console.log("当前滚动位置:", scrollTop, "视频高度阈值:", videoHeight.value);

  if (scrollTop >= videoHeight.value && !isTabSticky.value) {
    isTabSticky.value = true;
    console.log("设置 isTabSticky 为 true");
    // 缩小视频并暂停播放
    if (videoPlayer.value) {
      videoPlayer.value.pause();
    }
  } else if (scrollTop < videoHeight.value && isTabSticky.value) {
    isTabSticky.value = false;
    console.log("设置 isTabSticky 为 false");
    // 恢复视频大小（可选：恢复播放）
    // if (videoPlayer.value) {
    //   videoPlayer.value.play();
    // }
  }
};
</script>


<style scoped>
/* 容器使用相对定位 */
.container {
  position: relative;
  width: 100%;
  height: 100vh;
  overflow: hidden; /* 防止双重滚动 */
}

/* 主滚动视图样式 */
.main-scroll {
  width: 100%;
}

/* Header 和 Video Header 固定高度 */
.header, .video-header {
  flex: 0 0 auto;
}

/* Header 样式 */
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

/* 视频标题样式 */
.video-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  transition: all 0.3s ease;
}

.video-header-shrink {
  height: 50px; /* 缩小后的高度 */
}

.content {
  font-size: 18px;
  font-weight: bold;
  transition: font-size 0.3s ease;
}

.video-header-shrink .content {
  font-size: 16px; /* 缩小后的字体大小 */
}

/* 视频内容样式 */
.video-content {
  position: relative; /* 确保加载动画能正确定位 */
  flex: 0 0 auto;
  padding: 20px;
  background-color: #ffffff;
  margin: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.video-content-shrink {
  height: 200px; /* 缩小后的高度，根据需求调整 */
  width: 100%; /* 保持宽度不变，适应布局 */
  margin: 10px 20px; /* 缩小后的外边距 */
}

.video-content-shrink .video-container {
  max-width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Tab 选择部分样式 */
.tab-selection {
  flex: 0 0 auto;
  display: flex;
  justify-content: center;
  margin: 20px;
  transition: all 0.3s ease;
  position: sticky;
  top: 70px; /* header(50px) + margin(20px) */
  background-color: #f0f0f0;
  z-index: 999;
}

.tab-selection-sticky {
  /* 当 isTabSticky 为 true 时，保持在顶部 */
  top: 0;
}

/* Tab 容器样式 */
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
  transition: background-color 0.3s, color 0.3s;
}

.tab-container button.active {
  background-color: #4caf50;
  color: #ffffff;
  font-weight: bold;
}

/* Tab 内容样式 */
.tab-content {
  flex: 1 1 auto;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

/* 作者信息样式 */
.author-info {
  display: flex;
  flex-direction: column;
  padding: 20px;
  background-color: #ffffff;
  margin: 20px 0;
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
  align-items: center; /* 垂直居中 */
  margin-bottom: 10px;
}

.author-username {
  font-weight: bold;
  margin-right: 20px;
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

.author-interactions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

/* 相关推荐部分样式 */
.sidebar {
  padding: 20px;
  background-color: #ffffff;
  margin: 20px 0;
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
  cursor: pointer;
}

.recommendation-info {
  font-size: 14px;
  color: #555;
}

/* 评论部分样式 */
.comments-section {
  padding: 20px;
  background-color: #ffffff;
  margin: 20px 0;
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
  margin-right: 5px;
}

.comment-text {
  font-size: 14px;
  color: #555;
}

.comment-interactions {
  display: flex;
  margin-top: 10px;
  gap: 10px;
}

.comment-interactions button {
  border: none;
  background-color: transparent;
  cursor: pointer;
  font-size: 14px;
  color: #888;
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

/* 加载动画覆盖层样式 */
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  z-index: 10; /* 确保加载动画在最前 */
}

.loading-spinner {
  width: 50px;
  height: 50px;
}

/* 确保 scroll-view 内部内容有足够的空间 */
.author-info,
.sidebar,
.comments-section {
  padding: 20px;
}

/* 动态调整视频内容和标题的样式 */
.video-header-shrink {
  height: 50px; /* 缩小后的高度 */
}

.video-header-shrink .content {
  font-size: 16px; /* 缩小后的字体大小 */
}

.video-content-shrink {
  height: 200px; /* 缩小后的高度，根据需求调整 */
  width: 100%; /* 保持宽度不变，适应布局 */
  margin: 10px 20px; /* 缩小后的外边距 */
}

.video-content-shrink .video-container {
  max-width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
