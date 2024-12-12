<template>
  <view class="container">
    <!-- 页面功能提示 -->
    <view class="header">
      <text class="title">用户反馈</text>
      <text class="description">请填写您的反馈意见，我们将不断改进。</text>
    </view>

    <!-- 文本输入框 -->
    <view class="input-section">
      <textarea
        v-model="feedbackText"
        :maxlength="maxLength"
        placeholder="请输入您的反馈"
        class="textarea"
      ></textarea>
      <text class="char-count">{{ feedbackText.length }}/{{ maxLength }}</text>
    </view>

    <!-- 图片上传部分 -->
    <view class="image-section">
      <view class="images-container">
        <view v-for="(image, index) in images" :key="index" class="image-item">
          <image :src="image" class="uploaded-image"></image>
          <button class="remove-image-button" @click="removeImage(index)">x</button>
        </view>
      </view>
      <button type="primary" @click="addImage" class="add-image-button">添加图片</button>
    </view>

    <!-- 提交反馈按钮 -->
    <button type="primary" class="submit-button" @click="submitFeedback">
      提交反馈
    </button>
  </view>
</template>

<script setup>
import { ref } from 'vue'

const feedbackText = ref('')
const maxLength = 500
const images = ref([])

// 添加图片的方法
const addImage = () => {
  uni.chooseImage({
    count: 9 - images.value.length, // 限制最多上传9张图片
    sizeType: ['original', 'compressed'], // 可以选择原图或压缩图
    sourceType: ['album', 'camera'], // 可以选择相册或拍照
    success: (res) => {
      // 将选择的图片路径添加到images数组中
      images.value = images.value.concat(res.tempFilePaths)
    },
    fail: (err) => {
      console.error('选择图片失败:', err)
    }
  })
}

// 删除图片的方法
const removeImage = (index) => {
  images.value.splice(index, 1)
}

// 提交反馈的方法框架
const submitFeedback = () => {
  if (!feedbackText.value.trim()) {
    uni.showToast({
      title: '请输入反馈内容',
      icon: 'none'
    })
    return
  }

  // TODO: 实现提交反馈的逻辑，例如发送到服务器
  console.log('提交的反馈内容:', feedbackText.value)
  console.log('上传的图片:', images.value)

  uni.showToast({
    title: '反馈已提交',
    icon: 'success'
  })

  // 清空表单
  feedbackText.value = ''
  images.value = []
}
</script>

<style scoped>
.container {
  padding: 20px;
  background-color: #ffffff;
  min-height: 100vh;
}

.header {
  margin-bottom: 20px;
}

.title {
  font-size: 24px;
  font-weight: bold;
}

.description {
  font-size: 14px;
  color: #666666;
  margin-top: 5px;
}

.input-section {
  margin-bottom: 20px;
}

.textarea {
  width: 100%;
  height: 150px;
  border: 1px solid #dddddd;
  border-radius: 5px;
  padding: 10px;
  resize: none;
}

.char-count {
  text-align: right;
  font-size: 12px;
  color: #999999;
  margin-top: 5px;
}

.image-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  margin-bottom: 10px;
}

.images-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.image-item {
  position: relative;
  width: 100px;
  height: 100px;
}

.uploaded-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 5px;
}

.remove-image-button {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: rgba(0, 0, 0, 0.5);
  color: #ffffff;
  border: none;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  text-align: center;
  line-height: 20px;
}

.add-image-button {
  width: 100%;
  margin-top: 10px;
}

.submit-button {
  width: 100%;
  padding: 15px;
  background-color: #007AFF;
  color: #ffffff;
  border: none;
  border-radius: 5px;
  font-size: 16px;
}
</style>
