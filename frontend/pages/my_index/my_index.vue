<template>
  <view class="container">
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 个人信息部分 -->
    <view class="profile-section">
      <view class="profile-top">

        <!-- 修改后的头像按钮，透明且无边框 -->
        <button
          v-if="isLoggedIn"
          class="avatar-button"
          open-type="chooseAvatar"
          @chooseavatar="onChooseAvatar"
        >
          <image :src="avatarSrc" class="avatar" />
        </button>

        <!-- 若未登录，仍显示默认头像，但不可更换 -->
        <image
          v-else
          :src="avatarSrc"
          class="avatar"
        />

        <view class="profile-text">
          <!-- 昵称显示/编辑 -->
          <view v-if="!isEditingNickname" class="greeting-container">
            <text
              class="greeting"
              @click="enableNicknameEdit"
            >
              {{ isLoggedIn ? nickname : $t('profile_greeting') }}
            </text>
            <image src="@/pages/static/editor.svg" class="edit-icon"></image>
          </view>
          <input
            v-else
            type="nickname"
            v-model="nickname"
            class="nickname-input-inline"
            @blur="onNicknameBlur"
            @input="onNicknameInput"
            placeholder="请输入昵称"
            focus
          />
          <view>
            <text class="login-prompt">
              {{ isLoggedIn ? $t('profile_logged_in') : $t('profile_login_prompt') }}
            </text>
          </view>
        </view>
      </view>

      <!-- 未登录时的登录按钮 -->
      <button v-if="!isLoggedIn" class="login-button" @click="handleLoginButtonClick">
        {{ $t('profile_register_login') }}
      </button>
    </view>

    <!-- 删除了底部的 nickname-input-section -->

    <!-- 菜单部分 -->
    <view class="menu-section">
      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="navigateTo('searchTools')"
      >
        <image src="@/pages/static/search.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_search_tools')}}</text>
      </view>

      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="navigateTo('setGoals')"
      >
        <image src="@/pages/static/setgoals.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_set_goals')}}</text>
      </view>

      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="navigateTo('foodPreferences')"
      >
        <image src="@/pages/static/food.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_food_preferences')}}</text>
      </view>

      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="navigateTo('favorites')"
      >
        <image src="@/pages/static/favorites.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_favorites')}}</text>
      </view>

      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="navigateTo('my_home')"
      >
        <image src="@/pages/static/mywork.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_creations')}}</text>
      </view>

      <view class="menu-item" @click="navigateTo('appSettings')">
        <image src="@/pages/static/setting.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_app_settings')}}</text>
      </view>

      <view class="menu-item" @click="navigateTo('userSettings')">
        <image src="@/pages/static/user.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('menu_user_settings')}}</text>
      </view>

      <view
        v-if="isLoggedIn"
        class="menu-item"
        @click="handleLogout"
      >
        <image src="@/pages/static/logout.svg" class="icon_svg"></image>
        <text class="menu-text">{{$t('profile_logout')}}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储
import { onShow } from '@dcloudio/uni-app';

const BASE_URL = 'http://122.51.231.155:8080';

// 国际化
const { t } = useI18n();
// Pinia 用户存储
const userStore = useUserStore();

// 计算属性，从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);

/**
 * 用户头像
 * 若后端给出的 userStore.user.avatarUrl 是相对路径,
 * 则拼接: `${BASE_URL}/static/${userStore.user.avatarUrl}`
 * 否则使用默认背景
 */
const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
);

/**
 * 用户昵称
 * 默认显示 store 内的 nickName
 * 也可以在 onShow 或 onMounted 时进行赋值
 */
const nickname = ref(userStore.user.nickName || '');

// 控制昵称编辑状态
const isEditingNickname = ref(false);

// 页面显示时的逻辑
onShow(async () => {
  console.log('onShow triggered');
  console.log('token', userStore.user.token);

  if (isLoggedIn.value) {
    try {
      await userStore.fetchBasicDetails();
      // 更新 nickname 显示
      nickname.value = userStore.user.nickName || '';
      console.log('用户基本信息已刷新');
    } catch (error) {
      console.error('错误详情:', error);
      userStore.reset();
      uni.showToast({
        title: error.message || '刷新用户信息失败',
        icon: 'none',
        duration: 2000,
      });
    }
  }

  // 设置页面标题和底部 TabBar
  uni.setNavigationBarTitle({
    title: t('my_index'),
  });
  uni.setTabBarItem({
    index: 0,
    text: t('index'),
  });
  uni.setTabBarItem({
    index: 1,
    text: t('tools_index'),
  });
  uni.setTabBarItem({
    index: 2,
    text: t('news_index'),
  });
  uni.setTabBarItem({
    index: 3,
    text: t('my_index'),
  });
});

/** 导航到指定页面 */
function navigateTo(page) {
  uni.navigateTo({
    url: `/pagesMy/${page}/${page}`,
  });
}

/** 处理登录按钮点击 */
function handleLoginButtonClick() {
  // 跳转到登录页面
  navigateTo('login');
}

/** 处理登出 */
async function handleLogout() {
  try {
    await userStore.logout();
    uni.showToast({
      title: t('profile_logout_success'),
      icon: 'success',
      duration: 2000,
    });
  } catch (error) {
    uni.showToast({
      title: error.message || t('profile_logout_fail'),
      icon: 'none',
      duration: 2000,
    });
  }
}

/**
 * 用户点击 “选择头像” 按钮（open-type="chooseAvatar"）
 * 当用户手动选择了头像后，触发该回调
 */
async function onChooseAvatar(e) {
  try {
    // e.detail.avatarUrl 就是用户最终选定的头像临时路径
    const chosenAvatarUrl = e.detail.avatarUrl;
    console.log('用户选择的微信头像:', chosenAvatarUrl);

    // 可直接将其上传到服务器
    await userStore.setAvatar(chosenAvatarUrl);

    uni.showToast({
      title: t('profile_avatar_update_success'),
      icon: 'success',
      duration: 2000,
    });
  } catch (error) {
    console.error('选择头像/上传头像失败:', error);
    uni.showToast({
      title: error.message || t('profile_avatar_update_fail'),
      icon: 'none',
      duration: 2000,
    });
  }
}

/**
 * 启用昵称编辑模式
 */
function enableNicknameEdit() {
  if (isLoggedIn.value) {
    isEditingNickname.value = true;
    // 延迟聚焦，确保 input 已渲染
    setTimeout(() => {
      const input = document.querySelector('.nickname-input-inline');
      if (input) {
        input.focus();
      }
    }, 100);
  }
}

/**
 * 用户在 input(type="nickname") 输入时触发
 * 注意：微信小程序中，当用户点击键盘上的微信昵称，input 会直接被替换成微信昵称
 */
function onNicknameInput(e) {
  nickname.value = e.detail.value;
}

/**
 * 用户离开昵称输入框时(blur 事件)触发，向后端更新最新昵称
 */
async function onNicknameBlur() {
  if (!isEditingNickname.value) return;

  try {
    if (!nickname.value.trim()) {
      // 如果用户没有输入内容
      uni.showToast({
        title: '昵称不能为空',
        icon: 'none',
      });
      // 恢复成原 store 中的昵称或做其他处理
      nickname.value = userStore.user.nickName;
      isEditingNickname.value = false;
      return;
    }

    // 如果有内容，则更新
    await userStore.setNickname(nickname.value.trim());
    uni.showToast({
      title: t('success'),
      icon: 'success',
      duration: 2000,
    });
    isEditingNickname.value = false;
  } catch (error) {
    uni.showToast({
      title: error.message || t('fail'),
      icon: 'none',
      duration: 2000,
    });
  }
}
</script>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-color: #4caf50;
  --background-color: #f0f4f7;
  --text-color: #333;
  --secondary-text-color: #777;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

/* 容器 */
.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--background-color);
  font-family: var(--font-family);
}

/* 全屏背景图片 */
.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: -1;
  opacity: 0.1;
}

/* 个人信息部分 */
.profile-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx;
  background-color: rgba(33, 255, 6, 0.1);
  margin: 40rpx 20rpx;
  border-radius: 20rpx;
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
}

.profile-top {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  width: 100%;
  margin-bottom: 30rpx;
}

/* 修改后的头像按钮，透明且无边框 */
.avatar-button {
  background: transparent;
}

.avatar-button:after {
  border: none;
}

.avatar {
  width: 140rpx;
  height: 140rpx;
  border-radius: 70rpx;
  margin-right: 20rpx;
}

.profile-text {
  display: flex;
  flex-direction: column;
  justify-content: center;
  flex: 1;
}

/* 新增的容器样式 */
.greeting-container {
  position: relative; /* 设置为相对定位 */
  width: 100%;
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  padding: 10rpx 0; /* 根据需要调整内边距 */
}

/* 唯一的 .greeting 类定义 */
.greeting {
  font-size: 38rpx;
  color: var(--text-color);
  cursor: pointer;
  text-align: center; /* 确保文本居中 */
}

/* 新增的编辑图标样式 */
.edit-icon {
  position: absolute;
  top: 20rpx;
  right: -40rpx; /* 向右移动10rpx */
  width: 30rpx; /* 根据需要调整大小 */
  height: 30rpx;
  cursor: pointer;
}

.greeting {
  font-size: 38rpx;
  margin: 10rpx 0;
  color: var(--text-color);
  cursor: pointer;
}

/* 透明且无边框的内联昵称输入框 */
.nickname-input-inline {
  font-size: 38rpx;
  margin: 10rpx 0;
  padding: 0;
  border: none;
  border-bottom: 1rpx solid var(--border-color);
  background: transparent;
  color: var(--text-color);
  outline: none;
}

.login-prompt {
  color: var(--secondary-text-color);
}

.login-button {
  padding: 20rpx 40rpx;
  border: none;
  background-color: var(--primary-color);
  color: #ffffff;
  font-size: 32rpx;
  cursor: pointer;
  border-radius: 10rpx;
  transition: background-color 0.3s;
  width: 60%;
  margin-top: 10rpx;
}

.login-button:hover {
  background-color: #45a049;
}

/* 菜单部分 */
.menu-section {
  background-color: rgba(33, 255, 6, 0.06);
  margin: 40rpx 20rpx;
  border-radius: 20rpx;
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10rpx);
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid var(--border-color);
  cursor: pointer;
  transition: background-color 0.3s;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:hover {
  background-color: #f9f9f9;
}

.icon {
  font-size: 48rpx;
  color: var(--primary-color);
  margin-right: 30rpx;
}

.icon_svg {
  width: 50rpx;
  height: 50rpx;
  margin-left: 10rpx;
  margin-right: 30rpx;
  cursor: pointer;
}

.icon_svg:hover {
  transform: scale(1.2);
  fill: #45a049;
}

.menu-text {
  font-size: 36rpx;
  color: var(--text-color);
}

/* 响应式调整 */
@media screen and (max-width: 600px) {
  .profile-top {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .avatar {
    margin-right: 0;
    margin-bottom: 20rpx;
  }

  .profile-text {
    align-items: center;
  }

  .login-button {
    width: 100%;
  }

  .menu-item {
    padding: 20rpx;
  }

  .icon {
    margin-right: 20rpx;
  }

  .menu-text {
    font-size: 28rpx;
  }
}
</style>