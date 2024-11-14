<template>
  <view class="container">
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 个人信息部分 -->
    <view class="profile-section">
      <view class="profile-top">
        <image src="/static/images/index/background_img.jpg" class="avatar"></image>
        <view class="profile-text">
          <text class="greeting">{{ isLoggedIn ? uid : 'Hello!' }}</text>
          <text class="login-prompt">{{ isLoggedIn ? '已登录' : '登录以享受更多服务' }}</text>
        </view>
      </view>
      <button class="login-button" @click="handleLoginButtonClick">
        {{ isLoggedIn ? '切换账号' : '注册/登录' }}
      </button>
	  <button v-if="isLoggedIn" class="login-button" @click="logout">
	    {{ '退出登录' }}
	  </button>
    </view>

    <!-- 菜单部分 -->
    <view class="menu-section">
      <view class="menu-item" @click="navigateTo('setGoals')">
        <text class="icon">🎯</text>
        <text class="menu-text">设置目标</text>
      </view>
      <view class="menu-item" @click="navigateTo('foodPreferences')">
        <text class="icon">🍲</text>
        <text class="menu-text">食物偏好</text>
      </view>
      <view class="menu-item" @click="navigateTo('myFamily')">
        <text class="icon">👪</text>
        <text class="menu-text">我的家庭</text>
      </view>
      <view class="menu-item" @click="navigateTo('favorites')">
        <text class="icon">❤️</text>
        <text class="menu-text">我的收藏</text>
      </view>
      <view class="menu-item" @click="navigateTo('historyData')">
        <text class="icon">📊</text>
        <text class="menu-text">历史数据</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('appSettings')">
        <text class="icon">⚙️</text>
        <text class="menu-text">软件设置</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('userSettings')">
        <text class="icon">👤</text>
        <text class="menu-text">用户设置</text>
      </view>
      <view v-if="isLoggedIn" class="menu-item" @click="navigateTo('searchTools')">
        <text class="icon">🔍</text>
        <text class="menu-text">搜索工具</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';
import { onShow } from '@dcloudio/uni-app';

const uid = ref('');
const isLoggedIn = ref(false);

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
		z-index: 0;
		/* 将背景图片置于最底层 */
		opacity: 0.1;
		/* 调整透明度以不干扰内容 */
	}

	/* 头部 */
	.header {
		display: flex;
		align-items: center;
		padding: 20rpx;
		background-color: #ffffff;
		border-bottom: 1rpx solid var(--border-color);
	}

	.header-button {
		border: none;
		background-color: transparent;
		font-size: 36rpx;
		cursor: pointer;
		font-weight: bold;
		margin-right: 20rpx;
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
	}

	.profile-text {
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.greeting {
		font-size: 38rpx;
		margin: 10rpx 0;
		color: var(--text-color);
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

	/* 底部导航 */
	.footer {
		background-color: #ffffff;
		padding: 20rpx 0;
		box-shadow: 0 -4rpx 10rpx rgba(0, 0, 0, 0.1);
	}

	.footer-nav {
		display: flex;
		justify-content: space-around;
	}

	.nav-item {
		text-decoration: none;
		color: var(--text-color);
		font-weight: bold;
		transition: color 0.3s;
		font-size: 36rpx;
	}

	.nav-item:hover {
		color: var(--primary-color);
	}

	/* 响应式调整 */
	@media screen and (max-width: 600px) {
		.header {
			flex-direction: column;
			align-items: center;
		}

		.header-button {
			margin-right: 0;
			margin-bottom: 10rpx;
		}

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

		.footer-nav {
			flex-direction: column;
			align-items: center;
		}

		.nav-item {
			font-size: 28rpx;
			margin-bottom: 10rpx;
		}
	}
</style>