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
      <view v-for="(item, index) in items" :key="index" class="preview-item">
        <view class="item-content">
          <textarea
            v-if="item.type === 'text'"
            v-model="item.content"
            class="text-input"
            :placeholder="t('putin_text_placeholder')"
            autoHeight
          ></textarea>

          <!-- 图片上传功能 -->
          <view v-if="item.type === 'image'">
            <image
              :src="item.content"
              class="image-preview"
              :style="{ height: item.imageHeight + 'px' }"
              @click="handleImageChange(index)"
              mode="widthFix"
            />
            <textarea
              v-model="item.imageDescription"
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
        <image src="@/pagesNews/static/addtext.svg" alt="Add Text" class="icon" />
      </button>
      <button v-if="showfunctions" @click="addImage" class="function-btn">
        <image src="@/pagesNews/static/addpicture.svg" alt="Add Image" class="icon" />
      </button>
      <button v-if="showfunctions" @click="publish" class="push-btn">
        <image src="@/pagesNews/static/share.svg" alt="Publish" class="icon" />
      </button>
      <button v-if="showfunctions" @click="saveDraft" class="function-btn">
        <image src="@/pagesNews/static/save.svg" alt="Save" class="icon" />
      </button>
      <button v-if="showfunctions" @click="changefunction" class="function-btn">
        <image src="@/pagesNews/static/minus.svg" alt="-" class="icon" />
      </button>
      <button v-if="hidefunctions" @click="changefunction" class="add-btn">
        <image src="@/pagesNews/static/plus.svg" alt="+" class="icon" />
      </button>
    </view>

    <!-- 发布确认弹窗 -->
    <view v-if="showModal" class="modal">
      <view class="popup-content">
        <view class="popup-header">
          <image :src="authorAvatar" class="avatar" />
          <span class="nickname">{{ authorNickname }}</span>
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
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'

/* ----------------- Setup ----------------- */
const userStore = useUserStore()
const { t } = useI18n()

const BASE_URL = 'https://xcxcs.uwdjl.cn'
const PageId = ref('')

const authorNickname = computed(() => userStore.user.nickName)
const authorAvatar = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
)
const token = computed(() => userStore.user.token)

const title = ref('')
const description = ref('')
const items = ref([])
const showModal = ref(false)
const showfunctions = ref(true)
const hidefunctions = ref(false)
const isPublished = ref(false)

/* ----------------- Methods ----------------- */
function addText() {
  items.value.push({ type: 'text', content: '' })
}

function addImage() {
  items.value.push({
    type: 'image',
    content: '',
    itemHeight: 280,
    imageHeight: 200,
    imageDescription: ''
  })
}

function removeItem(index) {
  items.value.splice(index, 1)
}

function changefunction() {
  showfunctions.value = !showfunctions.value
  hidefunctions.value = !hidefunctions.value
}

function publish() {
  showModal.value = true
}

function confirmPublish() {
  const pageIdInt = parseInt(PageId.value, 10)
  if (isNaN(pageIdInt)) {
    uni.showToast({
      title: 'Invalid PageId',
      icon: 'none',
      duration: 2000
    })
    return
  }
  uni.request({
    url: `${BASE_URL}/news/convert_draft`,
    method: 'POST',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json'
    },
    data: {
      draft_id: pageIdInt
    },
    success: (res) => {
      if (res.data.message === 'Draft converted to news successfully.') {
        uni.showToast({
          title: '已发布',
          icon: 'success',
          duration: 2000
        })
        isPublished.value = true
        setTimeout(() => {
          uni.navigateBack()
        }, 2000)
      } else {
        uni.showToast({
          title: '发布失败',
          icon: 'none',
          duration: 2000
        })
        console.error('后端错误信息:', res.data.message)
      }
    },
    fail: (err) => {
      uni.showToast({
        title: '请求失败',
        icon: 'none',
        duration: 2000
      })
      console.error('请求失败', err)
    }
  })
}

function cancelPublish() {
  showModal.value = false
}

function uploadImage(filePath) {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${BASE_URL}/news/upload_image`,
      method: 'POST',
      header: {
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      },
      filePath: filePath,
      name: 'image',
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (data.message === 'Image uploaded successfully') {
            resolve(data.path)
          } else {
            reject(data.error)
          }
        } catch (error) {
          reject(`JSON 解析错误: ${error.message}`)
        }
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

async function saveDraft() {
  const postData = {
    title: title.value,
    paragraphs: [],
    image_paths: [],
    image_descriptions: []
  }

  // 简介作为第一段
  postData.paragraphs.push(description.value)
  postData.image_paths.push('')
  postData.image_descriptions.push('')

  // 组装组件内容
  for (const item of items.value) {
    if (item.type === 'image') {
      postData.paragraphs.push('')
      postData.image_paths.push(item.content)
      postData.image_descriptions.push(item.imageDescription || '')
    } else if (item.type === 'text') {
      postData.paragraphs.push(item.content || '')
      postData.image_paths.push('')
      postData.image_descriptions.push('')
    }
  }

  console.log("pageid.value", PageId.value)
  uni.request({
    url: `${BASE_URL}/news/drafts/${PageId.value}`,
    method: 'PUT',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json'
    },
    data: {
      title: postData.title,
      paragraphs: postData.paragraphs,
      image_descriptions: postData.image_descriptions,
      image_paths: postData.image_paths
    },
    success: (res) => {
      if (res.data.message === 'Draft updated successfully') {
        uni.showToast({
          title: '草稿已保存',
          icon: 'success',
          duration: 2000
        })
        PageId.value = res.data.draft_id
      } else {
        uni.showToast({
          title: '保存草稿失败',
          icon: 'none',
          duration: 2000
        })
        console.error('后端错误信息:', res.data.message)
      }
    },
    fail: (err) => {
      uni.showToast({
        title: '请求失败',
        icon: 'none',
        duration: 2000
      })
      console.error('请求失败', err)
    }
  })
}

function handleImageChange(index) {
  uni.chooseImage({
    count: 1,
    sourceType: ['album'],
    success: (res) => {
      const imagePath = res.tempFilePaths[0]
      items.value[index].content = imagePath

      // 动态计算图片尺寸
      uni.getImageInfo({
        src: imagePath,
        success: (info) => {
          const aspectRatio = info.width / info.height
          const newHeight = uni.getSystemInfoSync().windowWidth / aspectRatio
          items.value[index].imageHeight = newHeight
          items.value[index].itemHeight = newHeight + 80
        }
      })

      // 上传到服务器
      uploadImage(imagePath)
        .then((uploadedPath) => {
          const fullImageUrl = `${BASE_URL}/static/${uploadedPath}`
          items.value[index].content = fullImageUrl
        })
        .catch((error) => {
          console.error('图片上传服务器失败', error)
        })
    }
  })
}

/* ----------------- Lifecycle ----------------- */
onLoad(async (options) => {
  const articleId = options.id
  PageId.value = articleId

  const details = await getArticleDetails(PageId.value)
  if (!details) return

  title.value = details.title
  description.value = details.paragraphs[0]?.text || ''

  // 拼装组件
  const totalItems = Math.max(details.paragraphs.length, details.images.length)
  for (let i = 1; i < totalItems; i++) {
    if (details.paragraphs[i] && details.paragraphs[i].text) {
      items.value.push({
        type: 'text',
        content: details.paragraphs[i].text
      })
    }
    if (details.images[i] && details.images[i].url) {
      items.value.push({
        type: 'image',
        content: details.images[i].url,
        imageDescription: details.images[i].description || '',
        itemHeight: 280,
        imageHeight: 200
      })
    }
  }
})

function getArticleDetails(id) {
  return new Promise((resolve) => {
    uni.request({
      url: `${BASE_URL}/news/details/draft/${id}`,
      method: 'GET',
      header: {
        Authorization: `Bearer ${token.value}`
      },
      success: (res) => {
        resolve(res.data)
      },
      fail: (err) => {
        console.error('获取文章详情失败', err)
        resolve(null)
      }
    })
  })
}
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
  resize: none;
  overflow: auto;
}

.title-input {
  font-size: 24px;
  font-weight: bold;
  min-height: 80px;
  max-height: 200px;
}

.description-input {
  font-size: 16px;
  color: #555;
  min-height: 80px;
  max-height: 150px;
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

.text-input {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: none;
  outline: none;
  resize: none;
  min-height: 80px;
  max-height: 200px;
  overflow: auto;
  padding-right: 30px;
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
  margin-top: 12px;
  box-sizing: border-box;
  min-height: 40px;
  max-height: 40px;
}

.remove-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  z-index: 2;
}

.functions {
  position: fixed;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  background-color: rgba(0, 0, 0, 0.5);
  padding: 10px;
  border-radius: 8px;
  box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.3);
  z-index: 10;
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

.icon {
  width: 24px;
  height: 24px;
}

.icon:hover {
  transform: scale(1.2);
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
  z-index: 3;
}

.popup-content {
  background: white;
  padding: 20px;
  width: 70%;
  max-width: 350px;
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