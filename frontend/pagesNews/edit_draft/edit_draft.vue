<template>
  <view class="container">
    <!-- 文章标题 -->
    <view class="title-input-container">
      <textarea
        v-model="title"
        class="title-input"
        :placeholder="t('putin_title')"
        autoHeight
      ></textarea>
    </view>

    <!-- 文章简介 -->
    <view class="description-input-container">
      <textarea
        v-model="description"
        class="description-input"
        :placeholder="t('putin_introduction')"
        autoHeight
      ></textarea>
    </view>

    <!-- 预览区 -->
    <view class="preview">
      <view v-for="(item, index) in post.components" :key="item.id" class="preview-item">
        <view
          class="item-content"
          :style="{ height: item.style === 'image' ? item.imageHeight + 'px' : 'auto' }"
        >
          <textarea
            v-if="item.style === 'text'"
            v-model="item.content"
            class="text-input"
            :placeholder="t('putin_text_placeholder')"
            autoHeight
          ></textarea>

          <!-- 图片预览功能 -->
          <view v-if="item.style === 'image'">
            <image
              :src="item.content"
              class="image-preview"
              :style="{ height: item.imageHeight + 'px' }"
              @click="handleImageChange(index)"
            />
            <textarea
              v-model="item.description"
              class="image-description-input"
              :placeholder="t('add_description')"
              autoHeight
            ></textarea>
          </view>

          <!-- 删除按钮 -->
          <button @click="removeItem(index)" class="remove-btn">🗑️</button>
        </view>
      </view>
    </view>

    <!-- 功能区 -->
    <view class="functions">
      <button v-if="showfunctions" @click="addText" class="function-btn">
        <image src="@/pagesNews/static/addtext.svg" alt="Add Text" class="icon"></image>
      </button>
      <button v-if="showfunctions" @click="addImage" class="function-btn">
        <image src="@/pagesNews/static/addpicture.svg" alt="Add Image" class="icon"></image>
      </button>
      <button v-if="showfunctions" @click="publish" class="push-btn">
        <image src="@/pagesNews/static/share.svg" alt="Publish" class="icon"></image>
      </button>
      <button v-if="showfunctions" @click="saveDraft" class="function-btn">
        <image src="@/pagesNews/static/save.svg" alt="Save" class="icon"></image>
      </button>
      <button v-if="showfunctions" @click="changefunction" class="function-btn">
        <image src="@/pagesNews/static/minus.svg" alt="-" class="icon"></image>
      </button>
      <button v-if="hidefunctions" @click="changefunction" class="add-btn">
        <image src="@/pagesNews/static/plus.svg" alt="+" class="icon"></image>
      </button>
    </view>

    <!-- 发布确认弹窗 -->
    <view v-if="showModal" class="modal">
      <view class="popup-content">
        <view class="popup-header">
          <image :src="authorAvatar" class="avatar" />
          <span class="nickname">{{ authorNickname }}</span>
        </view>

        <view class="popup-body">
          <button @click="confirmPublish" class="confirm-btn">{{ t('ano_issue') }}</button>
        </view>

        <view class="popup-footer">
          <button @click="confirmPublish" class="confirm-btn">{{ t('confirm_issue') }}</button>
          <button @click="cancelPublish" class="cancel-btn">{{ t('cancel') }}</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useDraftStore } from '@/stores/draft';
import { useUserStore } from '@/stores/user'; // 引入 Pinia 用户存储

const draftStore = useDraftStore();
const userStore = useUserStore();

const BASE_URL = 'http://122.51.231.155:8080';

const authorNickname = computed(() => userStore.user.nickName);
const authorAvatar = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
);

const { t } = useI18n();

// 草稿数据初始化
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
});

const title = ref(post.value.title);
const description = ref(post.value.description);

// 根据post组件内容动态渲染
const addText = () => {
  post.value.components.push({ type: 'text', content: '' });
};

const addImage = () => {
  post.value.components.push({ type: 'image', content: '', imageHeight: 200, description: '' });
};

const removeItem = (index) => {
  post.value.components.splice(index, 1);
};

const handleImageChange = (index) => {
  uni.chooseImage({
    count: 1,
    success: (res) => {
      const imagePath = res.tempFilePaths[0];
      post.value.components[index].content = imagePath;
      
      // 调整图片高度
      uni.getImageInfo({
        src: imagePath,
        success: (info) => {
          const aspectRatio = info.width / info.height;
          const newHeight = uni.getSystemInfoSync().windowWidth / aspectRatio;
          post.value.components[index].imageHeight = newHeight;
        },
      });
    },
  });
};

const publish = () => {
  showModal.value = true;
};

const confirmPublish = () => {
  console.log('文章标题:', title.value);
  console.log('文章简介:', description.value);
  console.log('发布内容:', post.value.components);
  showModal.value = false;
};

const cancelPublish = () => {
  showModal.value = false;
};

const saveDraft = () => {
  draftStore.saveDraft();
  uni.showToast({
    title: '草稿已保存',
    icon: 'success',
    duration: 2000,
  });
};

const showModal = ref(false);
const showfunctions = ref(true);
const hidefunctions = ref(false);
</script>


<style scoped>
.container {
  padding: 20px;
}

.title-input-container,
.description-input-container {
  margin-bottom: 20px;
}

.title-input,
.description-input {
  width: 100%;
  padding: 15px;
  font-size: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-sizing: border-box;
  resize: none; /* 禁止用户手动调整大小 */
}

.title-input {
  font-size: 24px;
  font-weight: bold;
  min-height: 80px; /* 设置最小高度 */
  max-height: 200px; /* 设置最大高度 */
  overflow: auto; /* 超出部分可滚动 */
}

.description-input {
  font-size: 16px;
  color: #555;
  min-height: 80px; /* 设置最小高度 */
  max-height: 150px; /* 设置最大高度 */
  overflow: auto; /* 超出部分可滚动 */
}

.preview {
  margin-bottom: 20px;
}

.preview-item {
  margin-bottom: 15px;
}

.item-content {
  position: relative;
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 8px;
  box-sizing: border-box;
}

.text-input, {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: none;
  outline: none;
  resize: none; /* 禁止用户手动调整大小 */
  min-height: 80px; /* 设置最小高度 */
  max-height: 200px; /* 设置最大高度 */
  overflow: auto; /* 超出部分可滚动 */
  padding-right: 30px; /* 添加右内边距以避免覆盖删除按钮 */
}

.image-preview {
  width: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.image-description-input {
  width: 100%;
  padding: 8px;
  font-size: 12px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-top: 12px; /* 增加上边距，避免与其他组件重叠 */
  box-sizing: border-box;
  min-height: 40px; /* 设置最小高度 */
  max-height: 40px; /* 设置最大高度 */
}

.remove-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  z-index: 2; /* 确保删除按钮在最上层 */
}

/* 功能区固定左侧 */
.functions {
  position: fixed;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  background-color: rgba(0, 0, 0, 0.5); /* 半透明背景 */
  padding: 10px;
  border-radius: 8px;
  box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.3); /* 增加阴影效果 */
  z-index: 10; /* 确保按钮高于其他内容 */
  display: flex;
  flex-direction: column;
  align-items: center;
}

.function-btn,
.push-btn {
  margin-bottom: 10px;
  padding: 10px;
  background-color: #ffffff;
  color: black;
  border-radius: 50%;
  border: none;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
}

.add-btn {
  margin-bottom: 10px;
  padding: 10px;
  background-color: #ffffff;
  color: black;
  border-radius: 50%;
  border: none;
  font-size: 8px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
}

.push-btn {
  background-color: #4caf50;
}

.push-btn:hover {
  background-color: #45a049;
}

.function-btn:hover {
  background-color: #e6f0ff;
}

/* 按钮图标样式 */
.icon {
  width: 24px;
  height: 24px;
}

.icon:hover {
  transform: scale(1.2); /* 鼠标悬浮时放大图标 */
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 3; /* 更新 z-index */
}

.popup-content {
  background: white;
  padding: 20px;
  width: 70%; /* 减小弹窗宽度 */
  max-width: 350px; /* 设置最大宽度 */
  border-radius: 8px;
  box-sizing: border-box;
}

.popup-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.nickname {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.popup-body {
  margin-bottom: 20px;
}

.popup-footer {
  margin-top: 20px;
  text-align: right;
}

.confirm-btn,
.cancel-btn {
  padding: 8px 15px;
  border-radius: 5px;
  border: none;
  cursor: pointer;
}

.confirm-btn {
  background-color: #28a745;
  color: white;
  margin-right: 10px;
}

.cancel-btn {
  background-color: #dc3545;
  color: white;
  margin-right: 10px;
}

.confirm-btn:hover {
  background-color: #218838;
}

.cancel-btn:hover {
  background-color: #c82333;
}
</style>