<template>
  <!-- 外层容器 -->
  <view class="container">
    
    <!-- 作者信息及关注按钮 -->
   <!-- 作者信息，无关注按钮 -->
   <view class="author-header">
     <image 
       :src="formatAvatar(post.authoravatar)" 
       class="author-avatar"
     ></image>
     <text class="author-username">{{ post.authorname }}</text>
   </view>

    <!-- 文章标题和描述 -->
    <view class="title-container">
      <!-- 替换 h1 为 view 或 text -->
      <view class="article-title">{{ post.title }}</view>
      <!-- 替换 p 为 view 或 text -->
      <view class="article-description">{{ post.description }}</view>
    </view>

    <!-- 内容组件展示区 -->
    <view class="components-container">
      <view v-for="component in post.components" :key="component.id">
        
        <!-- 文本组件 -->
        <view v-if="component.style === 'text'" class="text-content">
          <text>{{ component.content }}</text>
        </view>

        <!-- 图片组件 -->
        <view v-if="component.style === 'image'" class="image-content">
          <image
            :src="component.content"
            class="image"
            mode="widthFix"
          ></image>
          <text class="image-description">{{ component.description }}</text>
        </view>
      </view>
    </view>
    
    <!-- 显示发布时间和阅读量 -->
    <view class="post-time">{{ formattedSaveTime }}</view>
    <view class="post-time">阅读量：{{ post.viewCount }}</view>

    <!-- 操作按钮：点赞、收藏、分享、踩等 -->
	<view class="inline-interaction-buttons">
	  
	  <!-- 点赞按钮 -->
	<button
	  class="action-button"
	  :class="{ active: ifLike }"
	  @click="toggleInteraction('like')"
	>
	  <!-- 根据 ifLike 状态动态切换图标 -->
	  <image 
		:src="ifLike 
		  ? '/pagesNews/static/liked.svg' 
		  : '/pagesNews/static/like.svg'"
		alt="Like" 
		class="icon"
	  />
	  <text class="count-text">{{ formatCount(post.likeCount) }}</text>
	</button>

	  <!-- 收藏按钮 -->
	  <button
		class="action-button"
		:class="{ active: ifFavourite }"
		@click="toggleInteraction('favorite')"
	  >
		<image 
		  :src="ifFavourite
		    ? '/pagesNews/static/favorited.svg' 
		    : '/pagesNews/static/favorite.svg'"
		  alt="Save" 
		  class="icon"
		></image>
		<text class="count-text">{{ formatCount(post.favoriteCount) }}</text>
	  </button>

	  <!-- 踩（dislike）按钮 -->
	  <button
		class="action-button"
		@click="toggleInteraction('dislike')"
	  >
		<image 
		  :src="ifDislike
		    ? '/pagesNews/static/disliked.svg' 
		    : '/pagesNews/static/dislike.svg'"
		  alt="Dislike" 
		  class="icon"
		></image>
		<!-- “dis” 或者直接显示 formatCount(post.dislikeCount) 也可 -->
		<text class="count-text">dis</text>
	  </button>

	</view>

    <!-- 评论区域 -->
    <view class="comments-section">
      <view class="comments-header">评论</view>
      <view id="comments-container">
        <view
          v-for="(comment, index) in limitedComments"
          :key="comment.id"
          class="comment"
        >
          <view class="comment-content">
            <image
              class="comment-avatar"
              :src="formatAvatar(comment.authorAvatar)"
            ></image>
            <view>
              <text class="comment-username">{{ comment.authorName }}: </text>
            </view>
			<view>
				<rich-text :nodes="renderCommentText(comment.text)" class="comment-text"></rich-text>
			</view>
          </view>
          <view class="comment-time">{{ comment.publish_time }}</view>
          
          <!-- 评论交互 -->
          <view class="comment-interactions">
            <button @click="toggleCommentLike(index)">
              👍 {{ comment.liked ? '已点赞' : '点赞' }}
            </button>
            <button @click="replyToComment(index)">
              💬 回复
            </button>
          </view>

          <!-- 回复输入区域 -->
          <view v-if="replyingTo === index" class="add-reply">
            <input
              type="text"
              v-model="newReply"
              placeholder="回复..."
            />
            <button @click="addReply(index)">发送</button>
          </view>

          <!-- 回复内容列表 -->
          <view
            v-if="comment.replies.length > 0"
            class="replies"
          >
			<view
				v-for="(reply, replyIndex) in limitedReplies(comment)"
				:key="reply.id"
				class="reply"
			  >
				<!-- 这里是显示每条回复的内容 -->
				<text class="comment-username">{{ reply.authorName }}</text>
				<text class="comment-text">:{{ reply.text }}</text>
				<text class="comment-time">{{ reply.publish_time }}</text>
			</view>
			<view
			    v-if="comment.replies.length > 3"
			    class="show-more-replies"
			    @click="toggleReplies(comment)"
			  >
			    <text v-if="!comment.showAllReplies" class="comment-time">
			      还有 {{ comment.replies.length - 3 }} 条回复
			    </text>
			    <text v-else class="comment-time">收起回复</text>
			</view>
			
          </view>
		  
        </view>
		<!-- 折叠/展开 按钮 -->
		<view v-if="comments.length > 5" class="show-more-comments" @click="toggleComments">
			<text v-if="!showAllComments" class="comment-time">还有 {{ comments.length - 5 }} 条评论</text>
			<text v-else class="comment-time">收起评论</text>
		</view>
      </view>
	  
	  
	  
      <!-- 发表评论 -->
      <view class="add-comment">
          <!-- 评论输入框 -->
          <input
            type="text"
            v-model="newComment"
            @input="handleCommentInput"
            placeholder="发表评论..."
          />
          <button @click="addComment">评论</button>
      
          <!-- 当检测到输入框中包含 '@' 时，弹出 popup  -->
          <uni-popup ref="mentionPopup" type="bottom" :mask="false" class="mention-popup">
            <view class="mention-list">
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
import { ref, reactive, onMounted, computed} from "vue";
import { getCurrentInstance } from 'vue';
import { onLoad } from "@dcloudio/uni-app";
import { useNewsStore } from '@/stores/news_list';
import { useI18n } from 'vue-i18n';
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app';
import { storeToRefs } from 'pinia';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储

const newsStore = useNewsStore();
const userStore = useUserStore(); // 使用用户存储

const BASE_URL = 'http://122.51.231.155:8080';
const BASE_URL_SH = 'http://122.51.231.155';
const PageId = ref('');

// 用来获取本地时间和日期
const systemDate = new Date();
const systemDateStr = systemDate.toISOString().slice(0, 10); // 获取当前系统日期，格式：YYYY-MM-DD

const jwtToken = computed(() => userStore.user.token);; // Replace with actual token

const newsData = ref([]);
const comments = reactive([]); // Initialize as empty array
const newComment = ref("");
const replyingTo = ref(null); // 当前正在回复的评论的索引
const newReply = ref(""); // 回复内容
const recommendations = ref([]);
const loadingError = ref(false); // 加载错误标志
const showMentionList = ref(false);
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
const avatarSrc_sh =  computed(() =>
    userStore.user.avatarUrl
        ? `${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);

const userListForMentions = computed(() => {
  // 从评论和回复中收集到的所有用户名
  let allNames = [];
  comments.forEach((c) => {
    allNames.push(c.authorName);
    c.replies.forEach((r) => allNames.push(r.authorName));
  });
  // 去重
  return Array.from(new Set(allNames));
});

/**
 * 将评论文本中的 “@xx” 高亮显示
 * 如果没有 @，则直接返回原文本
 * @param {string} text
 * @returns {string}
 */
const renderCommentText = (text) => {
  console.log('renderCommentText input:', text);

  // 如果 text 中没有 '@'，则直接返回原字符串
  if (!text.includes('@')) {
	  console.log('no @');
    return text;
  }

  // 如果包含 '@'，执行替换
  return text.replace(
    /@(\S+)\s/g,
    '<span style="color:blue;">@$1</span> '
  );
};

// 传入的post数据
const post = ref({
  components: [
  ],
});

//转换头像路径
const formatAvatar = (path) => {
	//console.log('解析的头像路径：',`${BASE_URL}/static/${path}`);
	return `${BASE_URL}/static/${path}`;
}

//转换数字
const formatCount = (count) => {
  return count < 10000 ? count : (count / 1000).toFixed(1) + 'k';
};

//转换时间
const formattedSaveTime = computed(() => {
  const postDate = post.value.savetime.slice(0, 10); // 提取日期部分

  if (postDate === systemDateStr) {
    // 如果日期相同，显示 "today" + 时间
    const postTime = new Date(post.value.savetime); // 转换为 Date 对象
    const hours = postTime.getHours().toString().padStart(2, '0');
    const minutes = postTime.getMinutes().toString().padStart(2, '0');
    const seconds = postTime.getSeconds().toString().padStart(2, '0');
    return `今天 ${hours}:${minutes}:${seconds}`;
  } else {
    // 否则显示完整日期
    return postDate;
  }
});

/**
 * @param {string} publishTime - ISO 格式或其他可被 Date 解析的字符串
 * @returns {string} - 格式化后的时间字符串
 */
const formatPublishTime = (publishTime) => {
  const date = new Date(publishTime);
  const now = new Date();

  // 判断是否是同一天
  const isSameDay =
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate();

  if (isSameDay) {
    // 如果是同一天，显示“今天 HH:mm”
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `今天 ${hours}:${minutes}`;
  } else {
    // 否则显示 YYYY-MM-DD 或者你想要的其他格式
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
  }
};

// 是否展开所有评论
const showAllComments = ref(false);

// 计算属性：根据 showAllComments 状态返回前5条或全部
const limitedComments = computed(() => {
	console.log('获取少量评论:',comments)
  if (showAllComments.value) {
    return comments;
  } else {
    return comments.slice(0, 5);
  }
});

// 点击切换的方法
const toggleComments = () => {
  showAllComments.value = !showAllComments.value;
};

// 计算“有限回复”的方法
const limitedReplies = (comment) => {
  if (comment.showAllReplies) {
    return comment.replies;
  } else {
    return comment.replies.slice(0, 3);
  }
};

// 切换展开/折叠
const toggleReplies = (comment) => {
  comment.showAllReplies = !comment.showAllReplies;
};

const { proxy } = getCurrentInstance(); // 在 <script setup> 下获取实例

const handleCommentInput = (e) => {
  const val = e.detail.value;
  // 如果输入含有 "@"
  if (val.includes("@")) {
    showMentionList.value = true;
    // 打开 popup
    proxy.$refs.mentionPopup.open();
  } else {
    // 关闭 popup
    showMentionList.value = false;
    proxy.$refs.mentionPopup.close();
  }
};

// 3. 点击某个昵称将其插入输入框
const insertMention = (name) => {
  // 简易示例：将第一个 "@" 替换成 "@User "
  newComment.value = newComment.value.replace("@", `@${name} `);
  proxy.$refs.mentionPopup.close();
  showMentionList.value = false;
};

const toggleInteraction = (type) => {
  // 确保 post 是一个 ref，并正确访问其属性
  const authorName = post.value.authorName; // 确保 post 对象中有 authorName 属性

  // 获取当前系统时间
  const systemDate = new Date();
  const systemDateStr = systemDate.toISOString().slice(0, 10); // YYYY-MM-DD
 
  // 处理操作
  if (type === "like") {
    if (ifLike.value === false) {
      console.log('点赞新闻');
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/like`,
        method: "POST",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: (res) => {
                  if (res.statusCode === 200) {
                    post.value.likeCount = res.data.like_count; // 使用后端返回的like_count
                    ifLike.value = true;
                    console.log('点赞状态:', ifLike.value);
                  } else {
                    console.error("Unexpected response:", res);
					uni.showToast({
					  title: '已经点过赞了~',
					  icon: 'none',
					  duration: 2000,
					});
					ifLike.value = true;
                  }
                },
                fail: (err) => {
                  console.error("Error liking news:", err);
                },
      });
    } else {
      console.log('取消点赞新闻');
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/like`,
        method: "DELETE",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: () => {
          post.value.likeCount--;
          ifLike.value = false;
          console.log('点赞状态:', ifLike.value);
		  uni.showToast({
		    title: '点赞已取消',
		    icon: 'none',
		    duration: 2000,
		  });
        },
        fail: (err) => {
          console.error("Error canceling like on news:", err);
        },
      });
    }
  }

  else if (type === "favorite") {
    if (ifFavourite.value === false) {
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/favorite`,
        method: "POST",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: (res) => {
                  if (res.statusCode === 200) {
                    post.value.favoriteCount = res.data.favorite_count; // 使用后端返回的favorite_count
                    ifFavourite.value = true;
                  } else {
                    console.error("Unexpected response:", res);
					uni.showToast({
					  title: '已经收藏了~',
					  icon: 'none',
					  duration: 2000,
					});
					ifFavourite.value = true;
                  }
                },
                fail: (err) => {
                  console.error("Error favoriting news:", err);
                },
      });
    } else {
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/favorite`,
        method: "DELETE",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: () => {
          post.value.favoriteCount--;
          ifFavourite.value = false;
		  uni.showToast({
		    title: '已取消收藏',
		    icon: 'none',
		    duration: 2000,
		  });
        },
        fail: (err) => {
          console.error("Error canceling favorite on news:", err);
        },
      });
    }
  }
  else if (type === "dislike") {
    if (ifDislike.value === false) {
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/dislike`,
        method: "POST",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: () => {
          post.value.dislikeCount++;
          ifDislike.value = true;
        },
        fail: (err) => {
          console.error("Error disliking news:", err);
        },
      });
    } else {
      uni.request({
        url: `http://122.51.231.155:8080/news/${PageId.value}/dislike`,
        method: "DELETE",
        header: {
          "Content-type": "application/json",
          "Authorization": `Bearer ${jwtToken.value}`,
        },
        data: {},
        success: () => {
          post.value.dislikeCount--;
          ifDislike.value = false;
        },
        fail: (err) => {
          console.error("Error canceling dislike on news:", err);
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

/**
 * 添加回复的函数，修改为使用新的后端接口
 */
const addReply = (index) => {
  if (newReply.value.trim()) {
    const parentComment = comments[index];
    uni.request({
      url: `${BASE_URL}/news/comments`,
      method: "POST",
      header: {
        "Content-type": "application/json",
        "Authorization": `Bearer ${jwtToken.value}`,
      },
      data: {
        news_id: parseInt(PageId.value), // 转换为 int
        content: newReply.value,
        is_reply: true,
        parent_id: parentComment.id, // 使用被回复评论的 id
      },
      success: (res) => {
        if (res.statusCode === 201) {
          const newReplyComment = res.data.comment;
          comments[index].replies.push({
            id: newReplyComment.id,
            text: newReplyComment.content,
            liked: newReplyComment.like_count > 0, // 根据需要调整
			authorName: uid.value,
			publish_time: formatPublishTime(newReplyComment.publish_time), // Format time
          });
          newReply.value = "";
          replyingTo.value = null; // 回复完成后取消回复状态
        } else {
          console.error("Unexpected response:", res);
          uni.showToast({
            title: '回复失败',
            icon: 'none',
            duration: 2000,
          });
        }
      },
      fail: (err) => {
        console.error("Error adding reply:", err);
        uni.showToast({
          title: '回复失败',
          icon: 'none',
          duration: 2000,
        });
      },
    });
  }
};

/**
 * 添加评论的函数，修改为使用新的后端接口
 */
const addComment = () => {
  if (newComment.value.trim()) {
    uni.request({
      url: `${BASE_URL}/news/comments`,
      method: "POST",
      header: {
        "Content-type": "application/json",
        "Authorization": `Bearer ${jwtToken.value}`,
      },
      data: {
        news_id: parseInt(PageId.value), // 转换为 int
        content: newComment.value,
        is_reply: false,
        parent_id: 0, // 设为 0 或 null 表示顶级评论
      },
      success: (res) => {
        if (res.statusCode === 201) {
          const newCommentData = res.data.comment;
          comments.push({
            id: newCommentData.id,
            text: newCommentData.content,
            liked: newCommentData.like_count > 0, // 根据需要调整
			authorName: uid.value,
			authorAvatar: avatarSrc_sh.value,
			publish_time: formatPublishTime(newCommentData.publish_time), // Format time
            replies: [],
			showAllReplies: false,
          });
          newComment.value = "";
        } else {
          console.error("Unexpected response:", res);
          uni.showToast({
            title: '发表评论失败',
            icon: 'none',
            duration: 2000,
          });
        }
      },
      fail: (err) => {
        console.error("Error adding comment:", err);
        uni.showToast({
          title: '发表评论失败',
          icon: 'none',
          duration: 2000,
        });
      },
    });
  }
};



// Simulate fetching data from backend
onLoad(async (options) => {
  const articleId = options.id;
  PageId.value = articleId;
  console.log('接收到的文章 ID:', articleId);

  // 根据 articleId 获取文章详情等操作
  const details = await getArticleDetails(PageId.value, false);
  console.log('获取的文章内容:', details);

  // 更新 post 对象
  post.value = {
    id: PageId.value,
    authoravatar: details.author.avatar_url,
    authorname: details.author.nickname,
    authorid: details.author.id,
    savetime: details.upload_time,
    title: details.title,
    description: details.paragraphs[0].text,
    components: [] ,// 初始化组件数组
	likeCount: details.like_count,
	shareCount: details.share_count,
	favoriteCount: details.favorite_count,
	followCount: 189,
	dislikeCount: details.dislike_count,
	viewCount: details.view_count,
	type: 'main',
  };

  // 更新 title 和 description
  //title.value = post.value.title;
  //description.value = post.value.description;

  // 遍历 paragraphs 和 images 填充 components
  const totalItems = Math.max(details.paragraphs.length, details.images.length);
  for (let index = 1; index < totalItems; index++) {
    // 处理段落文本
    if (details.paragraphs[index] && details.paragraphs[index].text) {
      post.value.components.push({
        id: index + 1, // 确保 id 从 1 开始
        content: details.paragraphs[index].text,
        style: 'text',
      });
    }

    // 处理图片
    if (details.images[index] && details.images[index].url) {
      post.value.components.push({
        id: index + 1, // 确保 id 从 1 开始
        content: details.images[index].url,
        style: 'image',
        description: details.images[index].description || '', // 如果没有描述，则为空
      });
    }
  }

  console.log('更新后的组件内容:', post.value.components);

  // 将 post 中的组件内容添加到 items 中
  // 处理评论
    if (details.comments && Array.isArray(details.comments)) {
      details.comments.forEach((comment) => {
        // Format the publish_time
        const formattedTime = formatPublishTime(comment.publish_time);
  
        // Construct the comment object
        const commentObj = {
          id: comment.id,
          text: comment.content,
          liked: comment.like_count > 0,
          publish_time: formattedTime,
		  authorName: comment.author.nickname,
		  authorAvatar: comment.author.avatar_url,
          replies: [],
		  showAllReplies: false,
        };
  
        // Process replies if any
        if (comment.replies && Array.isArray(comment.replies)) {
          comment.replies.forEach((reply) => {
            // Format the publish_time for replies
            const formattedReplyTime = formatPublishTime(reply.publish_time);
  
            // Construct the reply object
            const replyObj = {
              id: reply.id,
              text: reply.content,
              liked: reply.like_count > 0,
			  authorName: reply.author.nickname,
              publish_time: formattedReplyTime,
            };
  
            commentObj.replies.push(replyObj);
          });
        }
  
        // Add the comment to the comments array
        comments.push(commentObj);
      });
    } else {
      console.warn('No comments found in details.');
    }
	uni.request({
	  url: `http://122.51.231.155:8080/news/${PageId.value}/view`,
	  method: "POST",
	  header: {
	    "Content-type": "application/json",
	    "Authorization": `Bearer ${jwtToken.value}`, // 直接使用 jwtToken
	  },
	  data: {},
	  success: () => {
	    console.log("News view recorded successfully");
	  },
	  fail: (err) => {
	    console.error("Error viewing news:", err);
	  },
	});
	uni.request({
	  url: `http://122.51.231.155:8080/news/${PageId.value}/status`,
	  method: "GET",
	  header: {
	    "Content-type": "application/json",
	    "Authorization": `Bearer ${jwtToken.value}`, // 直接使用 jwtToken
	  },
	  data: {},
	  success: (res) => {
		  if(res.statusCode === 200){
			  console.log(res.data);
			  ifLike.value = res.data.liked;
			  ifDislike.value = res.data.disliked;
			  ifFavourite.value = res.data.favorited;
		  }
		  else{
			  console.log("Error getting status");
		  }
	    
	  },
	  fail: (err) => {
	    console.error("Error getting status:", err);
	  },
	});
});

// Function to get news or draft details
const getArticleDetails = async (id, isDraft = false) => {
  const url = isDraft
    ? `${BASE_URL}/news/details/draft/${id}`
    : `${BASE_URL}/news/details/news/${id}`;
  try {
    const res = await uni.request({
      url: url,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${jwtToken.value}`
      }
    });
    console.log('获取详细信息');
    console.log(res.data);
    return res.data;
  } catch (error) {
    console.error('Error fetching article details', error);
    return null;
  }
};
</script>

<style scoped>
.container {
  padding: 20px;
}

/*author part form video_detail*/
.author-header {
  display: flex;          /* 使头像与用户名排列在同一行 */
  align-items: center;    /* 垂直方向居中对齐 */
  margin-bottom: 10px;    /* 根据需求设置下方间距 */
}

.author-avatar {
  width: 50px;
  height: 50px;
  background-color: #ccc;
  border-radius: 50%;
  margin-right: 10px;      /* 头像和用户名之间留出合适间距 */
}

.author-username {
  font-weight: bold;
  /* 若需要在用户名和头像之间再留出一些距离，也可在这里增加 margin-left */
  /* margin-left: 10px; */
  font-size: 16px;        /* 根据需求设置文字大小 */
  color: #333;            /* 文字颜色可自行调整 */
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

/* 公共样式：确保按钮固定大小并进行水平排列 */
.action-button {
  width: 70px;          /* 固定宽度，按需调整 */
  height: 40px;         /* 固定高度，按需调整 */
  display: inline-flex; 
  align-items: center;  
  justify-content: center;
  border: none;
  border-radius: 5px;
  background-color: #f0f0f0; /* 默认背景色，可自行调整 */
  color: #333;               /* 默认文字颜色 */
  cursor: pointer;
  margin: 0 5px;            /* 每个按钮左右留出一点空隙 */
  overflow: hidden;         /* 保证文字不溢出 */
  transition: background-color 0.3s, color 0.3s;
}

/* 激活状态下的样式：如果按钮被点击过（ifLike/ifFavourite/ifShare/ifDislike 为 true）就添加该样式 */
.action-button.active {
  background-color: #4caf50;  /* 激活时的背景色示例，绿色 */
  color: #ffffff;            /* 激活时文字为白色 */
}

/* 图标样式：统一大小并与文字在一行 */
.icon {
  width: 16px;
  height: 16px;
  margin-right: 4px; /* 图标与数字之间的间距 */
}

/* 数字或文字部分 */
.count-text {
  font-size: 14px;
  line-height: 1;    /* 与 icon 高度相配合，避免偏移 */
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

.comment-time {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
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

/* popup 容器的样式，如高度、背景等 */
.mention-popup {
  height: auto;
  background-color: #ffffff;
}

/* 列表区域 */
.mention-list {
  display: flex;
  flex-direction: column;
  padding: 16px;
}

/* 单项 */
.mention-item {
  padding: 8px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}
</style>
