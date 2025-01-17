<template>
  <view class="container">
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    ></image>

    <!-- 个人信息 -->
    <view class="profile-section">
      <view class="profile-top">
        <!-- 若已登录，可点击头像选择新的微信头像；加上 @click="onAvatarClick" 事件 -->
        <button
          v-if="isLoggedIn"
          class="avatar-button"
          :open-type="hasPermission ? 'chooseAvatar' : ''"
          @click="onAvatarClick"
          @chooseavatar="onChooseAvatar"
        >
          <image :src="avatarSrc" class="avatar" />
        </button>

        <!-- 未登录时，仅显示默认头像 -->
        <image v-else :src="avatarSrc" class="avatar" />

        <view class="profile-text">
          <!-- 显示/编辑昵称；加上 @click="onNicknameClick" 事件 -->
          <view v-if="!isEditingNickname" class="greeting-container">
            <text class="greeting" @click="onNicknameClick">
              {{ isLoggedIn ? nickname : t('profile_greeting') }}
            </text>
            <image
              v-if="isLoggedIn"
              src="@/pages/static/editor.svg"
              class="edit-icon"
            />
          </view>
          <input
            v-else
            type="nickname"
            v-model="nickname"
            class="nickname-input-inline"
            @blur="onNicknameBlur"
            @input="onNicknameInput"
            :placeholder="t('placeholder_nickname')"
          />
          <view>
            <text class="login-prompt">
              {{ isLoggedIn ? t('profile_logged_in') : t('profile_login_prompt') }}
            </text>
          </view>
        </view>
      </view>

      <!-- 未登录时显示登录/注册按钮 -->
      <button v-if="!isLoggedIn" class="login-button" @click="handleLoginButtonClick">
        {{ t('profile_register_login') }}
      </button>
    </view>

    <!-- 菜单 -->
    <view class="menu-section">
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('searchTools')">
        <image src="@/pages/static/search.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_search_tools') }}</text>
      </view>

      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('setGoals')">
        <image src="@/pages/static/setgoals.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_set_goals') }}</text>
      </view>

      <view v-if="isLoggedIn" class="menu-item" @click="navigateToFoodPreferences('foodPreferences')">
        <image src="@/pages/static/food.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_food_preferences') }}</text>
      </view>

      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('favorites')">
        <image src="@/pages/static/favorites.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_favorites') }}</text>
      </view>

      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('my_home')">
        <image src="@/pages/static/mywork.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_creations') }}</text>
      </view>

      <view class="menu-item" @click="navigateTo('appSettings')">
        <image src="@/pages/static/setting.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_app_settings') }}</text>
      </view>

      <view class="menu-item" @click="navigateTo('userSettings')">
        <image src="@/pages/static/user.svg" class="icon_svg" />
        <text class="menu-text">{{ t('menu_user_settings') }}</text>
      </view>

      <view v-if="isLoggedIn" class="menu-item" @click="handleLogout">
        <image src="@/pages/static/logout.svg" class="icon_svg" />
        <text class="menu-text">{{ t('profile_logout') }}</text>
      </view>
    </view>

    <!-- 如果需要弹窗提示用户是否同意权限 -->
    <view v-if="showPermissionModal" class="permission-modal">
      <view class="permission-content">
        <text>需要获取您的头像和昵称，是否同意？</text>
        <view class="btn-row">
          <button class="btn-agree" @click="acceptPermission">同意</button>
          <button class="btn-reject" @click="rejectPermission">拒绝</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import {ref, computed, onMounted} from 'vue'
import {useI18n} from 'vue-i18n'
import {useUserStore} from '@/stores/user'
import {onShow} from '@dcloudio/uni-app'
import Login from "@/pagesMy/login/login.vue";

/* ----------------- Setup ----------------- */
const {t} = useI18n()
const userStore = useUserStore()

/* ----------------- Reactive & State ----------------- */
const BASE_URL = 'https://xcxcs.uwdjl.cn'
const isLoggedIn = computed(() => userStore.user.isLoggedIn)
const avatarSrc = computed(() =>
    userStore.user.avatarUrl && userStore.user.avatarUrl !== 'avatars/default.jpg'
        ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
        : '../static/default.jpg'
)

// 昵称编辑相关
const nickname = ref(userStore.user.nickName || '')
const isEditingNickname = ref(false)

// 是否已同意使用头像与昵称的权限
const hasPermission = ref(false)

// 控制是否展示自定义弹窗
const showPermissionModal = ref(false)

/* ----------------- Lifecycle ----------------- */
onShow(async () => {
  uni.setNavigationBarTitle({title: t('my_index')})
  uni.setTabBarItem({index: 0, text: t('index')})
  uni.setTabBarItem({index: 1, text: t('tools_index')})
  uni.setTabBarItem({index: 2, text: t('news_index')})
  uni.setTabBarItem({index: 3, text: t('my_index')})

  // 从 localStorage 中读取用户是否已经同意权限
  const storedPermission = uni.getStorageSync('hasPermission')
  if (storedPermission === 'true') {
    hasPermission.value = true
  }

  if (isLoggedIn.value) {
    try {
      await userStore.fetchBasicDetails()
      nickname.value = userStore.user.nickName || ''
    } catch (error) {
      console.error('获取用户信息失败:', error)
      userStore.reset()
      uni.showToast({
        title: error.message || t('fetch_user_info_error'),
        icon: 'none',
        duration: 2000
      })
    }
  }
})

/* ----------------- Methods ----------------- */
/**
 * 导航到指定页面
 */
function navigateTo(page) {
  uni.navigateTo({
    url: `/pagesMy/${page}/${page}`
  })
}

/** 导航到食物偏好设置页 */
function navigateToFoodPreferences(page) {
  uni.navigateTo({
    url: `/pagesTool/${page}/${page}`
  })
}

/** 登录按钮 */
function handleLoginButtonClick() {
  navigateTo('login')
}

/** 登出 */
async function handleLogout() {
  try {
    await userStore.logout()
    uni.showToast({
      title: t('success'),
      icon: 'success',
      duration: 2000
    })
  } catch (error) {
    uni.showToast({
      title: error.message || t('fail'),
      icon: 'none',
      duration: 2000
    })
  }
}

/**
 * 处理点击头像
 * 如果没权限，先显示弹窗；如果有权限，就执行上传逻辑
 */
function onAvatarClick(e) {
  if (!hasPermission.value) {
    showPermissionModal.value = true
    // 不执行任何头像上传逻辑
    return
  }
  // 已经同意权限则正常执行（在 button 上通过 @chooseavatar 绑定了 onChooseAvatar）
  // 这里不用再次调用 onChooseAvatar(e) 因为它会自动被触发
}

/**
 * 点击昵称 => 如果没权限则弹窗，否则允许编辑
 */
function onNicknameClick() {
  if (!hasPermission.value) {
    showPermissionModal.value = true
    return
  }
  if (isLoggedIn.value) {
    isEditingNickname.value = true
  }
}

/**
 * 用户手动输入昵称
 */
function onNicknameInput(e) {
  nickname.value = e.detail.value
}

/**
 * 当用户离开昵称输入框时提交更新
 */
async function onNicknameBlur() {
  if (!isEditingNickname.value) return

  if (!nickname.value.trim()) {
    uni.showToast({
      title: t('nickname_not_empty'),
      icon: 'none'
    })
    // 恢复原有昵称
    nickname.value = userStore.user.nickName
    isEditingNickname.value = false
    return
  }

  try {
    await userStore.setNickname(nickname.value.trim())
    uni.showToast({
      title: t('success'),
      icon: 'success',
      duration: 2000
    })
    isEditingNickname.value = false
  } catch (error) {
    uni.showToast({
      title: error.message || t('fail'),
      icon: 'none',
      duration: 2000
    })
  }
}

/**
 * 用户选择微信头像后，上传到后端
 */
async function onChooseAvatar(e) {
  try {
    const chosenAvatarUrl = e.detail.avatarUrl
    await userStore.setAvatar(chosenAvatarUrl)
    uni.showToast({
      title: t('success'),
      icon: 'success',
      duration: 2000
    })
  } catch (error) {
    console.error('上传头像失败:', error)
    uni.showToast({
      title: error.message || t('fail'),
      icon: 'none',
      duration: 2000
    })
  }
}

/**
 * 弹窗中用户点击“同意”
 */
function acceptPermission() {
  hasPermission.value = true
  uni.setStorageSync('hasPermission', 'true')
  showPermissionModal.value = false
}

/**
 * 弹窗中用户点击“拒绝”
 */
function rejectPermission() {
  showPermissionModal.value = false
}
</script>

<style scoped>
:root {
  --primary-color: #4caf50;
  --background-color: #f0f4f7;
  --text-color: #333;
  --secondary-text-color: #777;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--background-color);
  font-family: var(--font-family);
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
  pointer-events: none;
}

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

.avatar-button {
  background: transparent;
  padding: 0;
  margin-bottom: -40rpx;
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

.greeting-container {
  position: relative;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10rpx 0;
}

.greeting {
  font-size: 38rpx;
  color: var(--text-color);
  cursor: pointer;
  text-align: center;
}

.edit-icon {
  position: absolute;
  top: 20rpx;
  right: -40rpx;
  width: 30rpx;
  height: 30rpx;
  cursor: pointer;
}

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
  border-radius: 10rpx;
  width: 60%;
  margin-top: 10rpx;
}

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

/* 自定义弹窗示例 */
.permission-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.4);
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.permission-content {
  background-color: #fff;
  padding: 40rpx;
  border-radius: 20rpx;
  width: 80%;
  max-width: 600rpx;
  text-align: center;
}

.btn-row {
  display: flex;
  justify-content: space-around;
  margin-top: 40rpx;
}

.btn-agree {
  background-color: var(--primary-color);
  color: #fff;
  padding: 20rpx 40rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
  border: none;
}

.btn-reject {
  background-color: #ccc;
  color: #333;
  padding: 20rpx 40rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
  border: none;
}

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

  .icon_svg {
    margin-right: 20rpx;
  }

  .menu-text {
    font-size: 28rpx;
  }
}
</style>