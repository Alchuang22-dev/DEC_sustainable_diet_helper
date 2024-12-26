<template>
  <view class="container">
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    ></image>

    <!-- 个人信息 -->
    <view class="profile-section">
      <view class="profile-top">
        <!-- 若已登录，可点击头像选择新的微信头像 -->
        <button
          v-if="isLoggedIn"
          class="avatar-button"
          open-type="chooseAvatar"
          @chooseavatar="onChooseAvatar"
        >
          <image :src="avatarSrc" class="avatar" />
        </button>

        <!-- 未登录时，仅显示默认头像 -->
        <image v-else :src="avatarSrc" class="avatar" />

        <view class="profile-text">
          <!-- 显示/编辑昵称 -->
          <view v-if="!isEditingNickname" class="greeting-container">
            <text class="greeting" @click="enableNicknameEdit">
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
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { onShow } from '@dcloudio/uni-app'

/* ----------------- Setup ----------------- */
const { t } = useI18n()
const userStore = useUserStore()

/* ----------------- Reactive & State ----------------- */
const BASE_URL = 'https://xcxcs.uwdjl.cn'
const isLoggedIn = computed(() => userStore.user.isLoggedIn)
const avatarSrc = computed(() =>
  userStore.user.avatarUrl
    ? `${BASE_URL}/static/${userStore.user.avatarUrl}`
    : '/static/images/index/background_img.jpg'
)

const nickname = ref(userStore.user.nickName || '')
const isEditingNickname = ref(false)

/* ----------------- Lifecycle ----------------- */
onShow(async () => {
  uni.setNavigationBarTitle({ title: t('my_index') })
  uni.setTabBarItem({ index: 0, text: t('index') })
  uni.setTabBarItem({ index: 1, text: t('tools_index') })
  uni.setTabBarItem({ index: 2, text: t('news_index') })
  uni.setTabBarItem({ index: 3, text: t('my_index') })

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
 * 点击昵称 => 允许编辑
 */
function enableNicknameEdit() {
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