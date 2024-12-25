<template>
  <view class="settings">
    <view class="header">
      <view @click="goBack" class="back-icon">{{ t('back') }}</view>
      <text class="title">{{ userData.username }}</text>
      <view class="header-actions">
        <button class="menu-icon"></button>
        <button class="camera-icon"></button>
      </view>
    </view>

    <view class="list">
      <view class="list-item">
        <text>{{ t('settings_profile') }}</text>
        <text class="numbers right-align">{{ userData.username }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item">
        <text>{{ t('settings_id') }}</text>
        <text class="numbers right-align">{{ userData.accountId }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="Seal">
        <text>{{ t('settings_phonenumber') }}</text>
        <text class="numbers right-align">{{ userData.phoneNumber }}</text>
        <text class="arrow">></text>
      </view>
      <view class="list-item" @click="Seal">
        <text>{{ t('settings_email') }}</text>
        <text class="numbers right-align">{{ userData.email }}</text>
        <text class="arrow">></text>
      </view>
      <view class="divider"></view>
      <view class="list-item centered" @click="Seal">
        <text>{{ t('settings_security') }}</text>
      </view>
      <view class="list-item centered red-text" @click="Seal">
        <text>{{ t('settings_frozen') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'

/* ----------------- Setup ----------------- */
const userStore = useUserStore()
const { t } = useI18n()

const nickname = computed(() => userStore.user.nickName)
const user_id = computed(() => userStore.user.uid)
const email = computed(() => userStore.user.email)

const userData = ref({
  username: nickname.value,
  accountId: user_id.value,
  phoneNumber: '',
  email: email.value
})

/* ----------------- Methods ----------------- */
function goBack() {
  uni.navigateBack()
}

function Seal() {
  uni.showToast({
    title: '正在开发',
    icon: 'error',
    duration: 2000
  })
}

/* ----------------- Lifecycle ----------------- */
onMounted(() => {
  // 可以在此处发请求获取更多详细信息
})
</script>

<style scoped>
.settings {
  height: 100%;
  background: #f8f8f8;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.title {
  font-size: 18px;
  font-weight: bold;
}
.header-actions button {
  background: none;
  border: none;
}
.list {
  margin-top: 10px;
}
.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}
.list-item.centered {
  justify-content: center;
}
.numbers {
  color: #ccc;
}
.numbers.right-align {
  margin-left: auto;
}
.arrow {
  color: #ccc;
}
.divider {
  height: 1px;
  background-color: #ebebeb;
  margin: 10px 0;
}
.red-text {
  color: red;
}
</style>