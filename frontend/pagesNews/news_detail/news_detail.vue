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

      <!-- Dynamic content -->
      <view class="components">
        <view
          v-for="component in post.components"
          :key="component.id"
          class="component-item"
        >
          <!-- Text component -->
          <view
            v-if="component.style === 'text'"
            class="text-block"
          >
            <text>{{ component.content }}</text>
          </view>

          <!-- Image component -->
          <view
            v-if="component.style === 'image'"
            class="image-block"
          >
            <image
              :src="component.content"
              class="content-image"
              mode="widthFix"
            />
            <text v-if="component.description" class="image-caption">
              {{ component.description }}
            </text>
          </view>
        </view>
      </view>

      <!-- Article metadata -->
      <view class="metadata">
        <text class="timestamp">{{ formattedSaveTime }}</text>
        <text class="views">阅读量：{{ post.viewCount }}</text>
      </view>

      <!-- Interaction buttons -->
      <view class="interaction-bar">
        <view
          class="interaction-btn"
          :class="{ 'active': ifLike }"
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
          :class="{ 'active': ifFavourite }"
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
          :class="{ 'active': ifDislike }"
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

        <!-- Comments list -->
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
              <rich-text :nodes="renderCommentText(comment.text)" class="comment-text"></rich-text>
            </view>

            <view class="comment-actions">
              <view
                class="action-btn like-btn"
                :class="{ 'active': comment.liked }"
                @click="toggleCommentLike(index)"
              >
                <image
                  :src="comment.liked ? '/pagesNews/static/liked.svg' : '/pagesNews/static/like.svg'"
                  class="action-icon"
                />
                <text>{{ formatCount(comment.likecount) }}</text>
              </view>

              <view class="action-btn reply-btn" @click="replyToComment(index)">
                <uni-icons type="chat" size="16"></uni-icons>
                <text>解决</text>
              </view>
            </view>

            <!-- Reply input -->
            <view v-if="replyingTo === index" class="reply-input">
              <input
                type="text"
                v-model="newReply"
                placeholder="回复..."
                class="reply-field"
              />
              <view class="send-btn" @click="addReply(index)">发送</view>
            </view>

            <!-- Replies list -->
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
                <text>{{ !comment.showAllReplies ?
                  `还有 ${comment.replies.length - 3} 条附加说明` :
                  '收起附加说明' }}
                </text>
              </view>
            </view>
          </view>

          <!-- Show more comments button -->
          <view
            v-if="comments.length > 5"
            class="show-more-comments"
            @click="toggleComments"
          >
            <text>{{ !showAllComments ?
              `还有 ${comments.length - 5} 条注释` :
              '收起注释' }}
            </text>
          </view>
        </view>

        <!-- Add comment section -->
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

        <!-- Mentions popup -->
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
const user_id = computed(() => userStore.user.uid);
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

const switchtoStranger = (id) => {
	uni.navigateTo({
	  url: `/pagesMy/stranger/stranger?id=${id}`,
	});
}

const toggleInteraction = (type) => {
  // 确保 post 是一个 ref，并正确访问其属性
  const authorName = post.value.authorName; // 确保 post 对象中有 authorName 属性

  // 获取当前系统时间
  const systemDate = new Date();
  const systemDateStr = systemDate.toISOString().slice(0, 10); // YYYY-MM-DD

  // 处理操作
  if (type === "like") {
    if (ifLike.value === false) {
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
	if(!comments[index].liked){
		uni.request({
		  url: `http://122.51.231.155:8080/news/${comments[index].id}/comment_like`,
		  method: "POST",
		  header: {
		    "Content-type": "application/json",
		    "Authorization": `Bearer ${jwtToken.value}`,
		  },
		  data: {},
		  success: () => {
		    comments[index].liked = true;
			comments[index].likecount++;
		  },
		  fail: (err) => {
		    console.error("Error canceling dislike on comments:", err);
		  },
		});
	}
   else{
	   uni.request({
	     url: `http://122.51.231.155:8080/news/${comments[index].id}/comment_like`,
	     method: "DELETE",
	     header: {
	       "Content-type": "application/json",
	       "Authorization": `Bearer ${jwtToken.value}`,
	     },
	     data: {},
	     success: () => {
	       comments[index].liked = false;
		   comments[index].likecount--;
	     },
	     fail: (err) => {
	       console.error("Error canceling dislike on comments:", err);
	     },
	   });
   }
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
            title: '解决说明失败',
            icon: 'none',
            duration: 2000,
          });
        }
      },
      fail: (err) => {
        console.error("Error adding reply:", err);
        uni.showToast({
          title: '结局说明失败',
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
			likecount: 0,
			authorName: uid.value,
			authorAvatar: avatarSrc_sh.value,
			authorid: user_id.value,
			publish_time: formatPublishTime(newCommentData.publish_time), // Format time
            replies: [],
			showAllReplies: false,
          });
          newComment.value = "";
        } else {
          console.error("Unexpected response:", res);
          uni.showToast({
            title: '添加说明失败',
            icon: 'none',
            duration: 2000,
          });
        }
      },
      fail: (err) => {
        console.error("Error adding comment:", err);
        uni.showToast({
          title: '添加说明失败',
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
		  likecount: comment.like_count,
          publish_time: formattedTime,
		  authorName: comment.author.nickname,
		  authorAvatar: comment.author.avatar_url,
		  authorid: comment.author.id,
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
    return res.data;
  } catch (error) {
    console.error('Error fetching article details', error);
    return null;
  }
};
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
  border-left: 8rpx solid #4CAF50;
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
  color: #4CAF50;
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
  background-color: #4CAF50;
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
  color: #4CAF50;
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
  background-color: #4CAF50;
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
