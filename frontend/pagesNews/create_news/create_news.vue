<template>
  <view class="container">
    <!-- 文章标题 -->
    <view class="title-input-container">
      <input v-model="title" class="title-input" :placeholder="$t('putin_title')" />
    </view>

    <!-- 文章简介 -->
    <view class="description-input-container">
      <input v-model="description" class="description-input" :placeholder="$t('putin_introduction')" />
    </view>

    <!-- 预览区 -->
    <view class="preview">
      <view v-for="(item, index) in items" :key="index" class="preview-item">
        <view class="item-content" :style="{ height: item.type === 'image' ? item.itemHeight + 'px' : 'auto' }">
          <input v-if="item.type === 'text'" v-model="item.content" class="text-input" />
          
          <!-- 图片上传功能 -->
          <view v-if="item.type === 'image'">
            <image :src="item.content" class="image-preview" :style="{ height: item.imageHeight + 'px' }" @click="handleImageChange(index)" />
            <input type="file" accept="image/*" class="image-upload-input" :placeholder="$t('add_description')" />
          </view>
          
          <!-- 删除按钮改为小图标 -->
          <button @click="removeItem(index)" class="remove-btn">🗑️</button>
        </view>
      </view>
    </view>

    <!-- 功能区 -->
    <view class="functions">
      <button @click="addText" class="function-btn">{{$t('add_words')}}</button>
      <button @click="addImage" class="function-btn">{{$t('add_image')}}</button>
      <button @click="publish" class="push-btn">{{$t('release_but')}}</button>
      <button @click="saveDraft" class="function-btn">{{$t('save_draft')}}</button>
    </view>

    <!-- 发布确认弹窗 -->
    <view v-if="showModal" class="modal">
      <view class="popup-content">
        <view class="popup-header">{{$t('issue_confirm')}}</view>
        <view class="popup-body">
          <button @click="confirmPublish" class="confirm-btn">{{$t('ano_issue')}}</button>
        </view>
        <view class="popup-footer">
          <button @click="confirmPublish" class="confirm-btn">{{$t('confirm_issue')}}</button>
          <button @click="cancelPublish" class="cancel-btn">{{$t('cancel')}}</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n';

const title = ref('')  // 文章标题
const description = ref('')  // 文章简介
const items = ref([])  // 预览区的内容
const showModal = ref(false)  // 控制发布确认弹窗的显示与否
const { t } = useI18n();

// 添加文字
const addText = () => {
  items.value.push({ type: 'text', content: '' })
}

// 添加图片
const addImage = () => {
  items.value.push({ type: 'image', content: '',ItemHeight: 220, imageHeight: 200, imageDescription: '' }) // 初始化图片项
}

// 删除项目
const removeItem = (index) => {
  items.value.splice(index, 1)
}

// 发布
const publish = () => {
  showModal.value = true
}

// 确认发布
const confirmPublish = () => {
  console.log('文章标题:', title.value)
  console.log('文章简介:', description.value)
  console.log('发布内容:', items.value)
  showModal.value = false
}

// 取消发布
const cancelPublish = () => {
  showModal.value = false
}

// 保存草稿
const saveDraft = () => {
  console.log('草稿已保存', { title: title.value, description: description.value, content: items.value })
  uni.showToast({
    title: '草稿已保存',
    icon: 'success',
    duration: 2000,
  });
}

// 处理图片上传
const handleImageChange = (index) => {
  console.log("正在更改图片");

  // 调用 uni.chooseImage 来选择图片
  uni.chooseImage({
    count: 1, // 选择一张图片
    sourceType: ['album'], // 只从相册中选择
    success: (res) => {
      const imagePath = res.tempFilePaths[0]
      items.value[index].content = imagePath

      // 获取图片的宽高比
      uni.getImageInfo({
        src: imagePath,
        success: (info) => {
          // 计算新的高度，保持图片的宽高比
          const aspectRatio = info.width / info.height
          const newHeight = uni.getSystemInfoSync().windowWidth / aspectRatio
          items.value[index].imageHeight = newHeight // 保存计算后的高度
		  items.value[index].ItemHeight = newHeight + 20
        },
        fail: (err) => {
          console.error('获取图片信息失败', err)
        }
      })
    },
    fail: (err) => {
      console.error('上传图片失败', err)
    }
  })
}
</script>

<style scoped>
.container {
  padding: 20px;
}

.title-input-container, .description-input-container {
  margin-bottom: 20px;
}

.title-input, .description-input {
  width: 100%;
  padding: 15px;
  font-size: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-sizing: border-box;
}

.title-input {
  font-size: 24px;
  font-weight: bold;
  height: 80px;  /* 增加高度 */
}

.description-input {
  font-size: 16px;
  color: #555;
  height: 80px;  /* 增加高度 */
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
}

.text-input {
  width: 100%;
  height: 100px;
  border: none;
  outline: none;
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
  margin-top: 12px;  /* 增加上边距，避免与其他组件重叠 */
  height: 20px;
  box-sizing: border-box;
}

.remove-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
}

.functions {
  margin-top: 20px;
}

.function-btn {
  margin-right: 10px;
  padding: 8px 15px;
  background-color: #ffffff;
  color: black;
  border-radius: 5px;
  border: none;
  font-size: 14px;
  cursor: pointer;
}

.function-btn:hover {
  background-color: #0056b3;
}

.push-btn {
  margin-right: 10px;
  padding: 8px 15px;
  background-color: #4caf50;
  color: white;
  border-radius: 5px;
  border: none;
  font-size: 14px;
  cursor: pointer;
}

.push-btn:hover {
  background-color: #0056b3;
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
  z-index: 1;
}

.popup-content {
  background: white;
  padding: 20px;
  width: 80%;
  max-width: 400px;
  border-radius: 8px;
}

.popup-header {
  font-weight: bold;
  margin-bottom: 10px;
}

.popup-footer {
  margin-top: 20px;
  text-align: right;
}

.confirm-btn, .cancel-btn {
  padding: 8px 15px;
  border-radius: 5px;
  border: none;
  cursor: pointer;
}

.confirm-btn {
  background-color: #28a745;
  color: white;
}

.cancel-btn {
  background-color: #dc3545;
  color: white;
}

.confirm-btn:hover {
  background-color: #218838;
}

.cancel-btn:hover {
  background-color: #c82333;
}
</style>
