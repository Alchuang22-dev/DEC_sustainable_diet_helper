<template>
  <view class="login-container">
    <!-- 背景图 -->
    <image 
      src="https://cloud.tsinghua.edu.cn/thumbnail/cf9dba3a498247469fd4/1024/login.png"
      class="background-image"
      mode="widthFix"
    />

    <!-- 正中央按钮：微信登录示例 -->
    <view class="wechat-login">
      <button class="wechat-button" @tap="testLogin">
        <span>{{ t('loginWithWeChat') }}</span>
      </button>
    </view>

    <!-- 页面底部文字：其他方式登录（仅示例） -->
    <view class="other-login-text" @click="login">
      {{ t('loginWithOthers') }}
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { useI18n } from 'vue-i18n'
import { useUserStore } from "@/stores/user"

/* ----------------- Setup ----------------- */
const { t, locale, messages } = useI18n()
const userStore = useUserStore()

console.log('Current locale:', locale.value)
console.log('Available messages:', messages)

/* ----------------- Methods ----------------- */
async function testLogin() {
  try {
    uni.showLoading({ title: '登录中' })
    await userStore.login()
    uni.showToast({
      title: '登录成功',
      icon: 'success',
      duration: 2000,
    })
    uni.switchTab({
      url: '/pages/my_index/my_index',
    })
  } catch (error) {
    uni.showToast({
      title: error.message || '登录失败，请重试',
      icon: 'none',
      duration: 2000,
    })
  } finally {
    uni.hideLoading()
  }
}

async function login() {
  try {
    uni.showLoading({title: '登录中'})
    await userStore.login()
    uni.showToast({
      title: '登录成功',
      icon: 'success',
      duration: 2000,
    })
    uni.switchTab({
      url: '/pages/my_index/my_index',
    })
  } catch (error) {
    uni.showToast({
      title: error.message || '登录失败，请重试',
      icon: 'none',
      duration: 2000,
    })
  } finally {
    uni.hideLoading()
  }
}
</script>

<style scoped>
body {
  margin: 0;
  padding: 0;
  overflow: hidden;
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
}

.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  height: 100vh;
  background-color: rgba(248, 251, 247, 0.8);
  z-index: 1;
}

.wechat-login {
  display: flex;
  justify-content: center;
  align-items: center;
}

.wechat-button {
  background-color: transparent;
  color: #ffffff;
  border: 2px solid #ffffff;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}

.other-login-text {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
  z-index: 2;
}
</style>