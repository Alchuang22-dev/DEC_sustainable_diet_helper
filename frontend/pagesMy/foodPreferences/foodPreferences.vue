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
import { ref } from 'vue';
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
  { name: t('glutenFree'), key: 'glutenFree', icon: 'https://via.placeholder.com/50' },
  { name: t('alcoholFree'), key: 'alcoholFree', icon: 'https://via.placeholder.com/50' },
  { name: t('dairyFree'), key: 'dairyFree', icon: 'https://via.placeholder.com/50' },
]);

const showModal = ref(false);

const removePreference = (index) => {
  const preferenceToRemove = preferences.value[index];
  uni.request({
    url: 'http://122.51.231.155:8080/preferences',
    method: 'DELETE',
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
  uni.request({
    url: 'http://122.51.231.155:8080/preferences',
    method: 'POST',
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
    },
    fail: (err) => {
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
