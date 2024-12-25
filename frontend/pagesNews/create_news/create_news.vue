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
        <view
          class="item-content"
          :style="{ height: item.type === 'image' ? item.itemHeight + 'px' : 'auto' }"
        >
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
            />
            <textarea
              v-model="item.imageDescription"
              class="image-description-input"
              :placeholder="t('add_description')"
              autoHeight
            ></textarea>
          </view>

          <!-- 删除按钮改为小图标 -->
          <button @click="removeItem(index)" class="remove-btn">🗑️</button>
        </view>
      </view>
    </view>

    <!-- 功能区 -->
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
          <!-- 显示作者头像和昵称 -->
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
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n'
import { onHide, onUnload } from '@dcloudio/uni-app'; // 导入 onUnload
import { useDraftStore } from '../stores/draft';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储
const draftStore = useDraftStore();
const userStore = useUserStore();

const BASE_URL = 'http://122.51.231.155:8080';
const BASE_URL_SH = 'http://122.51.231.155';

const authorNickname = computed(() => userStore.user.nickName);
const authorAvatar = computed(() =>
    userStore.user.avatarUrl
        ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);
const user_id = computed(() => userStore.user.uid);
const token = computed(() => userStore.user.token);
const { t } = useI18n()

const title = ref('') // 文章标题
const description = ref('') // 文章简介
const items = ref([]) // 预览区的内容
const showModal = ref(false) // 控制发布确认弹窗的显示与否
const showfunctions = ref(true)
const hidefunctions = ref(false)
const PageId = ref(0)//草稿编号

// 添加文字
const addText = () => {
  items.value.push({ type: 'text', content: '' })
}

const changefunction = () => {
	if(showfunctions.value === true){
		showfunctions.value = false;
		hidefunctions.value = true;
	}
	else{
		showfunctions.value = true;
		hidefunctions.value = false;
	}
}

// 添加图片
const addImage = () => {
  items.value.push({ type: 'image', content: '', itemHeight: 280, imageHeight: 200, imageDescription: '' }) // 初始化图片项
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
const confirmPublish = async () => {
  console.log('确认发布');
  try {
    const pageId = await saveDraft(); // 等待 saveDraft 完成并获取 pageId
    console.log('草稿保存编号:', pageId);

    const pageIdInt = parseInt(pageId, 10); // 转换为整数
    if (isNaN(pageIdInt)) {
      uni.showToast({
        title: 'Invalid PageId',
        icon: 'none',
        duration: 2000,
      });
      return;
    }

    // 调用 API 将草稿转换为新闻
    uni.request({
      url: `${BASE_URL}/news/convert_draft`, // 转换草稿的 API 路径
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token.value}`, // 使用当前 token
        'Content-Type': 'application/json',
      },
      data: {
        draft_id: pageIdInt, // 使用转换后的整数 PageId
      },
      success: (res) => {
        if (res.data.message === 'Draft converted to news successfully.') {
          uni.showToast({
            title: '草稿已发布为新闻',
            icon: 'success',
            duration: 2000,
          });
          showModal.value = false; // 关闭弹窗
        } else {
          uni.showToast({
            title: '发布失败',
            icon: 'none',
            duration: 2000,
          });
          console.error('后端错误信息:', res.data.message);
        }
      },
      fail: (err) => {
        uni.showToast({
          title: '请求失败',
          icon: 'none',
          duration: 2000,
        });
        console.error('请求失败', err);
      }
    });
    showModal.value = false;
  } catch (error) {
    uni.showToast({
      title: '保存草稿失败，请重试',
      icon: 'none',
      duration: 2000,
    });
    console.error('保存草稿失败:', error);
  }
};



// 取消发布
const cancelPublish = () => {
  showModal.value = false
}

//上传图片
const uploadImage = (filePath) => {
	console.log(token.value);
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${BASE_URL}/news/upload_image`, // 上传图片的 API 地址
      method: 'POST',
      header: {
        "Authorization": `Bearer ${token.value}`, // 替换为实际的 Token 变量
        "Content-Type": "application/json", // 设置请求类型
      },
      filePath: filePath,
      name: 'image', // form-data 中字段名
      success: (res) => {
        console.log('上传图片返回结果:', res); // 打印响应内容用于调试
        try {
          const data = JSON.parse(res.data); // 解析返回的 JSON 数据
          if (data.message === 'Image uploaded successfully') {
            resolve(data.path); // 返回图片相对路径
			console.log(data.path);
          } else {
            reject(data.error); // 上传失败，返回错误信息
          }
        } catch (error) {
          reject(`JSON 解析错误: ${error.message}`); // 解析失败时的错误提示
        }
      },
      fail: (err) => {
        reject(err); // 请求失败，返回错误信息
      }
    });
  });
};


const saveDraft = async () => {
  // 生成草稿对象，包含标题、简介、组件内容等
  const post = {
    title: title.value, // 文章标题
    description: description.value, // 文章简介
    components: items.value.map((item, index) => {
      if (item.type === 'text') {
        return { id: index + 1, content: item.content, style: 'text' };
      } else if (item.type === 'image') {
        return { 
          id: index + 1, 
          content: item.content, 
          style: 'image', 
          description: item.imageDescription || '' 
        };
      }
    })
  };

  const data = {
    title: post.title,
    paragraphs: [], // 用于存放文本段落
    images: [], // 用于存放图片链接
    image_descriptions: [] // 用于存放图片描述
  };

  // 默认简介为第一个自然段
  data.paragraphs.push(description.value);
  data.images.push(''); // 先添加一个空的图片路径
  data.image_descriptions.push('');

  // 上传所有图片并填充图片路径
  const imagePaths = await Promise.all(
    post.components.map((item) => {
      if (item.style === 'image' && item.content) {
        data.paragraphs.push(''); // 添加空段落
        data.images.push(item.content); // 保存上传后的图片路径
        data.image_descriptions.push(item.description || ''); // 保存图片描述
      } else if (item.style === 'text') {
        data.paragraphs.push(item.content || ''); // 添加文字段落
        data.images.push(''); // 空白图片路径
        data.image_descriptions.push(''); // 空白图片描述
      }
    })
  );

  return new Promise((resolve, reject) => {
    // 提交草稿数据到服务器
    uni.request({
      url: `${BASE_URL}/news/create_draft`,
      method: 'POST',
      header: {
        'Authorization': `Bearer ${token.value}`,
        'Content-Type': 'application/json',
      },
      data: {
        title: data.title,
        paragraphs: data.paragraphs,
        image_descriptions: data.image_descriptions,
        image_paths: data.images,
      },
      success: (res) => {
        if (res.data.message === 'Draft created successfully.') {
          PageId.value = res.data.draft_id;
          resolve(PageId.value); // 成功时返回 draft_id
		  uni.showToast({
		    title: '草稿已保存',
		    icon: "fail",
		    duration: 2000,
		  });
        } else {
          reject('保存草稿失败'); // 失败时拒绝 Promise
        }
      },
      fail: (err) => {
        reject(err); // 请求失败，返回错误信息
      }
    });
  });
};



// 处理图片上传
const handleImageChange = (index) => {
  console.log("正在更改图片");

  uni.chooseImage({
    count: 1, // 选择一张图片
    sourceType: ['album'], // 只从相册中选择
    success: (res) => {
      const imagePath = res.tempFilePaths[0];
      items.value[index].content = imagePath;

      // 获取图片的宽高比
      uni.getImageInfo({
        src: imagePath,
        success: (info) => {
          const aspectRatio = info.width / info.height;
          const newHeight = uni.getSystemInfoSync().windowWidth / aspectRatio;
          items.value[index].imageHeight = newHeight;
          items.value[index].itemHeight = newHeight + 80;
        },
        fail: (err) => {
          console.error('获取图片信息失败', err);
        }
      });

      // 上传图片到服务器
      uploadImage(imagePath).then((uploadedPath) => {
        // 将上传返回的路径拼接成完整URL
        const fullImageUrl = `${BASE_URL}/static/${uploadedPath}`;
		console.log(fullImageUrl);
        items.value[index].content = fullImageUrl;
      }).catch((error) => {
        console.error('图片上传服务器失败', error);
      });
    },
    fail: (err) => {
      console.error('上传图片失败', err);
    }
  });
};


// Simulate fetching data from backend
onMounted(() => {
  console.log('进入新闻创作页');
});

// 在页面卸载时自动保存草稿
onHide(async () => {
  console.log('页面即将卸载，自动保存草稿');
  try {
    const pageId = await saveDraft();
    console.log('自动保存草稿成功，草稿编号:', pageId);
  } catch (error) {
    console.error('自动保存草稿失败:', error);
  }
});

onUnLoad(async () => {
  console.log('页面即将卸载，自动保存草稿');
  try {
    const pageId = await saveDraft();
    console.log('自动保存草稿成功，草稿编号:', pageId);
  } catch (error) {
    console.error('自动保存草稿失败:', error);
  }
});
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