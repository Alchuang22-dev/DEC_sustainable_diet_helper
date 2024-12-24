<template>
  <view class="login-container">
    <!-- 背景图 -->
    <image 
      src="https://cloud.tsinghua.edu.cn/thumbnail/cf9dba3a498247469fd4/1024/login.png"
      class="background-image"
      mode="widthFix"
    ></image>

    <!-- 放在正中央的按钮 -->
    <view class="wechat-login">
      <button class="wechat-button" @tap="testLogin">
        <span>{{ $t('loginWithWeChat') }}</span>
      </button>
    </view>

    <!-- 页面底部文字：其他方式登录（仅示例，不需实现具体逻辑） -->
    <view class="other-login-text" @click="login">
      {{ $t('loginWithOthers') }}
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useStore } from 'vuex' 

const { t, locale, messages } = useI18n()
const store = useStore()

console.log('Current locale:', locale.value)
console.log('Available messages:', messages)

const testLogin = () => {
  uni.navigateTo({
    url: "/pagesMy/wechatLogin/wechatLogin",
  })
}

const login = () => {
  uni.navigateTo({
    url: "/pagesMy/login_old/login_old",
  })
}
</script>

<style scoped>
/* 让 body 不出现滚动条 */
body {
  margin: 0;
  padding: 0;
  overflow: hidden; /* 关键：取消滚动条 */
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
}

/* 背景图充满屏幕 */
.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
}

/* 外层容器：居中对齐，撑满屏幕 */
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center; /* 垂直居中 */
  position: relative;
  height: 100vh; /* 撑满可视高度 */
  background-color: rgba(248, 251, 247, 0.8); 
  /* 加一点半透明背景，以便看清内容；可按需去掉 */
  z-index: 1;
}

/* 微信按钮容器，可根据需要再调整 */
.wechat-login {
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 按钮样式：背景绿色、文字白色 */
.wechat-button {
  background-color: transparent; 
  color: #ffffff;       /* 文字白色 */
  border: 2px solid #ffffff;  	
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}

/* 在页面底部的可点击文字 */
.other-login-text {
  position: absolute;
  bottom: 20px;
  left: 50%;               /* 水平居中 */
  transform: translateX(-50%);
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
  z-index: 2;
}
</style>
