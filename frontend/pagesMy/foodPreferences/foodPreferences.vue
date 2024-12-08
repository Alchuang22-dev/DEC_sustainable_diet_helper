<template>
  <view class="container">
    <view class="header">
      <text class="title">{{$t('preferences_title')}}</text>
    </view>
	<image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view group="preferences" class="preferences">
      <view v-for="(preference, index) in preferences" :key="index" :class="['preference-card', 'color-' + (index % 4)]">
        <image :src="preference.icon" class="preference-icon" />
        <text class="preference-name">{{ preference.name }}</text>
        <button class="delete-button" @click="removePreference(index)">{{$t('delete_button')}}</button>
      </view>
    </view>
    <view class="add-preference">
      <button @click="showPreferenceOptions">{{$t('add_preference_button')}}</button>
    </view>
    <view v-if="showModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{$t('modal_title')}}</text>
        <view v-for="(option, index) in preferenceOptions" :key="index" class="modal-option" @click="selectPreference(option)">
          <image :src="option.icon" class="option-icon" />
          <text class="option-name">{{ option.name }}</text>
        </view> 
      </view>
	  <view class="button-content">
		  <button class="close-button" @click="closeModal">{{$t('close_button')}}</button>
	  </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();

const preferences = ref([
  { name: t('foodpreference_greeting'), key: 'foodpreference_greeting', icon: 'https://via.placeholder.com/50' },
]);

const preferenceOptions = ref([
  { name: t('highProtein'), key: 'highProtein', icon: 'https://via.placeholder.com/50' },
  { name: t('highEnergy'), key: 'highEnergy', icon: 'https://via.placeholder.com/50' },
  { name: t('lowFat'), key: 'lowFat', icon: 'https://via.placeholder.com/50' },
  { name: t('lowCH'), key: 'lowCH', icon: 'https://via.placeholder.com/50' },
  { name: t('lowsodium'), key: 'lowsodium', icon: 'https://via.placeholder.com/50' },
  { name: t('vegan'), key: 'vegan', icon: 'https://via.placeholder.com/50' },
  { name: t('vegetarian'), key: 'vegetarian', icon: 'https://via.placeholder.com/50' },
  { name: t('glulenFree'), key: 'glulenFree', icon: 'https://via.placeholder.com/50' },
  { name: t('alcoholFree'), key: 'alcoholFree', icon: 'https://via.placeholder.com/50' },
  { name: t('dairyFree'), key: 'dairyFree', icon: 'https://via.placeholder.com/50' },
]);

const showModal = ref(false);

const user_id = ref('');

onMounted(() => {
  // 从 localStorage 中获取 token 信息作为 user_id
  const token = uni.getStorageSync('token');
  if (token) {
    user_id.value = token; // 如果存在，则将 token 存储为 user_id
  } else {
    console.warn('No tokens found in localStorage.');
  }

  // 请求偏好数据
  uni.request({
    url: 'http://122.51.231.155:8090/preferences',
    method: 'GET',
    header: {
      "Authorization": `Bearer ${user_id.value}`, // 替换为实际的 Token 变量
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {},
    success: (res) => {
      if (res.statusCode === 200) {
        // 处理返回的数据
		console.log('success to get preference');
        const data = res.data;  // 假设返回的数据格式是 { data: [...] }
        
        // 将后端数据转移到 preferences 数组中
        data.forEach(item => {
          preferences.value.push({
            name: t(item.name), // 通过翻译返回名称
            key: item.name,     // 设置 key 为后端的 name 字段
            icon: 'https://via.placeholder.com/50' // 默认 icon 地址
          });
        });
      } else {
        console.error('Failed to load preferences:', res.data);
      }
    },
    fail: (err) => {
      console.error('Error fetching preferences:', err);
    }
  });
});

const removePreference = (index) => {
  const preferenceToRemove = preferences.value[index];
  console.log(preferenceToRemove.key);
  uni.request({
    url: 'http://122.51.231.155:8090/preferences',
    method: 'DELETE',
	header: {
	      "Authorization": `Bearer ${user_id.value}`, // 替换为实际的 Token 变量
	      "Content-Type": "application/json", // 设置请求类型
	    },
    data: {
		preference_name: preferenceToRemove.key // 使用存储的 key 字段
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message); // 打印成功信息
        // 删除本地数组中的偏好
        preferences.value.splice(index, 1);
      }
    },
    fail: (err) => {
      console.error('Error removing preference:', err);
    }
  });
};

const showPreferenceOptions = () => {
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const selectPreference = (option) => {
	console.log(option.key);
  uni.request({
    url: 'http://122.51.231.155:8090/preferences',
    method: 'POST',
	header: {
	      "Authorization": `Bearer ${user_id.value}`, // 替换为实际的 Token 变量
	      "Content-Type": "application/json", // 设置请求类型
	    },
    data: {
		preference_name: option.key // 使用存储的 key 字段
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message); // 打印成功信息
        // 将新的偏好添加到本地
        preferences.value.push({ name: option.name, key: option.key, icon: option.icon });
        closeModal();
      }
	  console.log(res.data);
    },
    fail: (err) => {
		console.log(res.data);
		console.error('Error adding preference:', err);
    }
  });
};

</script>

<style scoped>
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.1;
}	

.container {
  padding: 20px;
}
.header {
  text-align: center;
  margin-bottom: 20px;
}
.title {
  font-size: 24px;
  font-weight: bold;
}
.preferences {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.preference-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
}
.color-0 {
  background-color: #4ca
}
.color-1 {
  background-color: #e0f7fa;
}
.color-2 {
  background-color: #ffe0b2;
}
.color-3 {
  background-color: #e1bee7;
}
.preference-icon {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}
.preference-name {
  flex: 1;
  font-size: 16px;
}
.delete-button {
  background-color: #ff4d4f;
  color: #fff;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}
.add-preference {
  margin-top: 20px;
  text-align: center;
}
.add-preference button {
  background-color: #4caf50;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column; /* 设置垂直排列 */
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 80%;
  max-width: 400px;
  max-height: 400px;
  overflow-y: auto; /* 允许垂直滚动 */
  margin-bottom: 20px; /* 添加底部间距 */
}

.modal-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 20px;
}

.modal-option {
  display: flex;
  align-items: center;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
}

.option-icon {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}

.option-name {
  font-size: 16px;
}

.button-content {
  display: flex;
  justify-content: center; /* 按钮居中 */
  width: 100%; /* 确保按钮容器宽度充满 */
}

.close-button {
  background-color: #ff4d4f;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 20px;
}

</style>
