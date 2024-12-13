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
          <view v-else @click="handleUsernameClick" class="username-container">
            <text class="greeting">{{ isLoggedIn ? uid : $t('profile_greeting') }}</text>
            <!-- 只在用户登录后显示编辑图标 -->
            <image
                v-if="isLoggedIn"
                src="@/pages/static/editor.svg"
                class="edit-icon"
                @click="handleUsernameClick"
            />
          </view>
          <view>
            <text class="login-prompt">
              {{ isLoggedIn ? $t('profile_logged_in') : $t('profile_login_prompt') }}
            </text>
          </view>
        </view>
      </view>
      <button v-if=!isLoggedIn class="login-button" @click="handleLoginButtonClick">
        {{ isLoggedIn ? $t('profile_switch_account') : $t('profile_register_login') }}
      </button>
    </view>

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
<!--      <view-->
<!--          v-if="isLoggedIn"-->
<!--          class="menu-item"-->
<!--          @click="navigateTo('myFamily')"-->
<!--      >-->
<!--        <image src="@/pages/static/family.svg" class="icon_svg"></image>-->
<!--        <text class="menu-text">{{$t('menu_my_family')}}</text>-->
<!--      </view>-->
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
import { ref, computed, nextTick, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '../../stores/user'; // 引入 Pinia 用户存储
import { onShow } from '@dcloudio/uni-app';
const BASE_URL = 'http://122.51.231.155:8080';

// 国际化
const { t } = useI18n();

// Pinia 用户存储
const userStore = useUserStore();

// 计算属性从 Pinia store 获取用户状态
const isLoggedIn = computed(() => userStore.user.isLoggedIn);
const uid = computed(() => userStore.user.nickName);
const avatarSrc = computed(() =>
    userStore.user.avatarUrl
        ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
        : '/static/images/index/background_img.jpg'
);

onShow(async () => {
  console.log('onShow triggered');  // 添加这行来确认是否触发
  console.log('token', userStore.user.token);
  console.log('avatarSrc', avatarSrc);

  if (isLoggedIn.value) {
    try {
      await userStore.fetchBasicDetails();
      console.log('用户基本信息已刷新');
    } catch (error) {
      console.error('错误详情:', error);  // 添加这行来查看具体错误
      userStore.reset();
      uni.showToast({
        title: error.message || '刷新用户信息失败',
        icon: 'none',
        duration: 2000,
      });
    }
  }
});

// 新增的响应式变量
const isEditingUsername = ref(false);
const newUsername = ref('');
const usernameInput = ref(null);

// 导航到指定页面
function navigateTo(page) {
  uni.navigateTo({
    url: `/pagesMy/${page}/${page}`,
  });
}

// 处理登录按钮点击
function handleLoginButtonClick() {
  if (isLoggedIn.value) {
    // 切换账户的逻辑，调用登出并跳转到登录页面
    uni.navigateTo({
      url: '/pagesMy/login/login',
    });
  } else {
    // 跳转到登录页面
    navigateTo('login');
  }
}

// 处理登出
async function handleLogout() {
  try {
    await userStore.logout();
    uni.showToast({
      title: t('profile_logout_success'),
      icon: 'success',
      duration: 2000,
    });
    // 可选择在登出后跳转到某个页面
  } catch (error) {
    uni.showToast({
      title: error.message || t('profile_logout_fail'),
      icon: 'none',
      duration: 2000,
    });
  }
}

// 处理头像点击
function handleAvatarClick() {
  if (isLoggedIn.value) {
    // 选择图片作为头像
    uni.chooseImage({
      count: 1,
      sourceType: ['album', 'camera'],
      success: async (res) => {
        const tempFilePath = res.tempFilePaths[0];
        try {
          console.log('tempFilePath', tempFilePath);
          await userStore.setAvatar(tempFilePath);
          uni.showToast({
            title: t('profile_avatar_update_success'),
            icon: 'success',
            duration: 2000,
          });
        } catch (error) {
          uni.showToast({
            title: error.message || t('profile_avatar_update_fail'),
            icon: 'none',
            duration: 2000,
          });
        }
      },
      fail: (err) => {
        console.error('选择头像失败', err);
        uni.showToast({
          title: t('profile_avatar_select_fail'),
          icon: 'none',
          duration: 2000,
        });
      },
    });
  }
}

// 处理用户名点击，进入编辑模式
function handleUsernameClick(event) {
  if (isLoggedIn.value) {
    event.stopPropagation(); // 阻止事件冒泡，避免触发其他点击事件
    isEditingUsername.value = true;
    newUsername.value = uid.value;
    nextTick(() => {
      if (usernameInput.value) {
        usernameInput.value.focus();
      }
    });
  }
}

// 提交新的用户名
async function submitUsername() {
  const trimmedUsername = newUsername.value.trim();
  if (trimmedUsername) {
    try {
      await userStore.setNickname(trimmedUsername);
      uni.showToast({
        title: t('profile_username_update_success'),
        icon: 'success',
        duration: 2000,
      });
    } catch (error) {
      uni.showToast({
        title: error.message || t('profile_username_update_fail'),
        icon: 'none',
        duration: 2000,
      });
    }
  } else {
    uni.showToast({
      title: t('profile_username_empty'),
      icon: 'none',
      duration: 2000,
    });
  }
  isEditingUsername.value = false;
}

// 页面显示时的逻辑
onShow(() => {
  // Pinia store 已经在初始化时加载了用户数据，无需额外操作
  // 更新页面标题和 TabBar 项目
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

/* 用户名编辑图标样式 */
.username-container {
  display: flex;
  align-items: center;
  position: relative;
}

.edit-icon {
  width: 24rpx;
  height: 24rpx;
  margin-left: 10rpx;
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
  width: 60%;
  margin-top: 10rpx;
}

.logout-button {
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

.logout-button:hover {
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
  transform: scale(1.2); /* 悬停时放大 20% */
  fill: #45a049; /* 悬停时改变图标颜色 */
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
