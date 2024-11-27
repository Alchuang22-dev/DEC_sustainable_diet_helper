<template>
  <view class="container">
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 个人信息部分 -->
    <view class="profile-section">
      <view class="profile-top">
        <image :src="avatarSrc" class="avatar" @click="handleAvatarClick"></image>
        <view class="profile-text">
          <template v-if="isEditingUsername">
            <input
              v-model="newUsername"
              class="username-input"
              @keyup.enter="submitUsername"
              @blur="submitUsername"
              ref="usernameInput"
            />
          </template>
          <view v-else @click="handleUsernameClick">
            <text class="greeting">{{ isLoggedIn ? uid : $t('profile_greeting') }}</text>
          </view>
		  <view>
			  <text class="login-prompt">{{ isLoggedIn ? $t('profile_logged_in') : $t('profile_login_prompt') }}</text>
		  </view>
        </view>
      </view>
      <button class="login-button" @click="handleLoginButtonClick">
        {{ isLoggedIn ? $t('profile_switch_account') : $t('profile_register_login') }}
      </button>
      <button v-if="isLoggedIn" class="login-button" @click="logout">
        {{ $t('profile_logout') }}
      </button>
    </view>

    <!-- 菜单部分 -->
    <view class="menu-section">
      <view class="menu-item" @click="navigateTo('setGoals')">
        <text class="icon">🎯</text>
        <text class="menu-text">{{$t('menu_set_goals')}}</text>
      </view>
      <view class="menu-item" @click="navigateTo('foodPreferences')">
        <text class="icon">🍲</text>
        <text class="menu-text">{{$t('menu_food_preferences')}}</text>
      </view>
      <view class="menu-item" @click="navigateTo('myFamily')">
        <text class="icon">👪</text>
        <text class="menu-text">{{$t('menu_my_family')}}</text>
      </view>
      <view class="menu-item" @click="navigateTo('favorites')">
        <text class="icon">❤️</text>
        <text class="menu-text">{{$t('menu_favorites')}}</text>
      </view>
      <view class="menu-item" @click="navigateTo('historyData')">
        <text class="icon">📊</text>
        <text class="menu-text">{{$t('menu_history_data')}}</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('appSettings')">
        <text class="icon">⚙️</text>
        <text class="menu-text">{{$t('menu_app_settings')}}</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('userSettings')">
        <text class="icon">👤</text>
        <text class="menu-text">{{$t('menu_user_settings')}}</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('searchTools')">
        <text class="icon">🔍</text>
        <text class="menu-text">{{$t('menu_search_tools')}}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, nextTick } from 'vue';
import { onShow } from '@dcloudio/uni-app';

const uid = ref('');
const avatarSrc = ref('/static/images/index/background_img.jpg');
const isLoggedIn = ref(false);

// 新增的响应式变量
const isEditingUsername = ref(false);
const newUsername = ref('');
const usernameInput = ref(null);

// 根据后端传递的id进行页面内容加载
function navigateTo(page) {
  uni.navigateTo({
    url: `/pagesMy/${page}/${page}`,
  });
}

function handleLoginButtonClick() {
  if (isLoggedIn.value) {
    // 切换账户的逻辑
    uni.navigateTo({
      url: '/pagesMy/login/login',
    });
  } else {
    // 跳转到登录页面
    navigateTo('login');
  }
}

function checkLoginStatus() {
  console.log("in check");
  const query = uni.getStorageSync('uid');
  console.log(query);
  if (query && query !== '') {
    uid.value = query;
    isLoggedIn.value = true;
  } else {
    isLoggedIn.value = false;
  }
}

function logout() {
  isLoggedIn.value = false;
  uni.removeStorageSync('uid');
}

function handleAvatarClick() {
  console.log("in changing avatar");
  if (isLoggedIn.value) {
    // 选择图片作为头像
    uni.chooseImage({
      count: 1,
      sourceType: ['album'],
	  // 省略了向后端的发信
      success: (res) => {
        avatarSrc.value = res.tempFilePaths[0];
      },
      fail: (err) => {
        console.error('选择头像失败', err);
      },
    });
  }
}

function handleUsernameClick() {
  if (isLoggedIn.value) {
    console.log("handleUsernameClick triggered");
    isEditingUsername.value = true;
    newUsername.value = uid.value;
    nextTick(() => {
      if (usernameInput.value) {
        usernameInput.value.focus();
      }
    });
  }
}

function submitUsername() {
  const trimmedUsername = newUsername.value.trim();
  if (trimmedUsername) {
    uid.value = trimmedUsername;
    uni.setStorageSync('uid', trimmedUsername);
	//省略了向后端的发信
    uni.showToast({
      title: '用户名已更新',
      icon: 'success',
    });
  } else {
    uni.showToast({
      title: '用户名不能为空',
      icon: 'none',
    });
  }
  isEditingUsername.value = false;
}

onShow(() => {
  console.log("in onShow");
  // 在页面显示时调用检查登录状态
  checkLoginStatus();
});
</script>


<style scoped>
  /* 全局样式变量 */
  :root {
    --primary-color: #4CAF50;
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
    /* 将背景图片置于最底层 */
    opacity: 0.1;
    /* 调整透明度以不干扰内容 */
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

  .avatar {
    width: 140rpx;
    height: 140rpx;
    border-radius: 70rpx;
    margin-right: 20rpx;
    cursor: pointer;
  }

  .profile-text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    flex: 1; /* 使其占据剩余空间 */
  }

  .greeting {
    font-size: 38rpx;
    margin: 10rpx 0;
    color: var(--text-color);
    cursor: pointer;
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
    width: 80%;
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

  .menu-text {
    font-size: 36rpx;
    color: var(--text-color);
  }

  /* 新增的输入框样式 */
  .username-input {
    font-size: 38rpx;
    padding: 5rpx 10rpx;
    border: 1rpx solid var(--border-color);
    border-radius: 5rpx;
    outline: none;
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
