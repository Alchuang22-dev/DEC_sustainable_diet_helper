<template>
	<head>
	    <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap" rel="stylesheet">
	</head>
	<view class="container">
		<!-- 全屏背景图片 -->
		<image src="/static/images/index/background_index_new.png" class="background-image"></image>

		<!-- 头部 -->
		<view class="dec_header">
			<image src="/static/images/index/logo_wide.png" :alt="$t('dec_logo_alt')" class="dec_logo" mode="aspectFit">
			</image>
			<text class="title">{{$t('welcome_title')}}</text>
		</view>

		<!-- 实用工具 -->
		<view class="useful-tools">
			<text class="tools-title">{{$t('tools_title')}}</text>
			<view class="tools-list">
				<view class="tool" @click="navigateTo('calculator')" animation="fadeInUp">
					<image src="https://cdn.pixabay.com/photo/2017/07/06/17/13/calculator-2478633_1280.png"
						:alt="$t('tool_carbon_calculator')" class="tool-icon" mode="aspectFill">
					</image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_carbon_calculator')}}</text>
						<text class="tool-info">{{$t('tool_carbon_calculator_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('recommend')" animation="fadeInUp" animation-delay="0.2s">
					<image src="https://cdn.pixabay.com/photo/2020/03/12/18/37/dish-4925892_1280.png"
						:alt="$t('tool_diet_recommendation')" class="tool-icon" mode="aspectFill">
					</image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_diet_recommendation')}}</text>
						<text class="tool-info">{{$t('tool_diet_recommendation_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('nutrition')" animation="fadeInUp" animation-delay="0.4s">
					<image src="https://cdn.pixabay.com/photo/2016/11/14/15/42/calendar-1823848_1280.png"
						:alt="$t('tool_nutrition_calculator')" class="tool-icon" mode="aspectFill">
					</image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_nutrition_calculator')}}</text>
						<text class="tool-info">{{$t('tool_nutrition_calculator_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('family')" animation="fadeInUp" animation-delay="0.6s">
					<image src="https://cdn.pixabay.com/photo/2016/01/04/14/24/terminal-board-1120961_1280.png"
						:alt="$t('tool_family_recipe')" class="tool-icon" mode="aspectFill">
					</image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_family_recipe')}}</text>
						<text class="tool-info">{{$t('tool_family_recipe_info')}}</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref,
		onMounted
	} from 'vue'
	import {
		useI18n
	} from 'vue-i18n'
  import { onShow } from '@dcloudio/uni-app';

	// 初始化 i18n
	const {
		t,
	} = useI18n()

  onShow(() => {
    uni.setNavigationBarTitle({
      title: t('tools_index')
    })
    uni.setTabBarItem({
      index: 0,
      text: t('index')
    })
    uni.setTabBarItem({
      index: 1,
      text: t('tools_index')
    })
    uni.setTabBarItem({
      index: 2,
      text: t('news_index')
    })
    uni.setTabBarItem({
      index: 3,
      text: t('my_index')
    })
  });

	// 页面跳转方法
	const navigateTo = (page) => {
		if (page === 'recommend') {
			uni.navigateTo({
				url: "/pagesTool/food_recommend/food_recommend",
			});
		} else if (page === 'nutrition') {
			uni.navigateTo({
				url: "/pagesTool/nutrition_calculator/nutrition_calculator",
			});
		} else if (page === 'family') {
			uni.navigateTo({
				url: "/pagesTool/home_servant/home_servant",
			});
		} else {
			uni.navigateTo({
				url: "/pagesTool/carbon_calculator/carbon_calculator",
			});
		}
	};
</script>

<style scoped>
	/* 通用变量 */
	:root {
		--primary-color: #4CAF50;
		--secondary-color: #2fc25b;
		--background-color: #f5f5f5;
		--card-background: rgba(255, 255, 255, 0.8);
		/* 半透明背景 */
		--text-color: #333;
		--shadow-color: rgba(0, 0, 0, 0.1);
		--font-size-title: 32rpx;
		--font-size-subtitle: 24rpx;
		--font-family: 'Roboto', sans-serif;
	}

	.title, .tools-title, .tool-name, .tool-info {
	    font-family: var(--font-family);
	}
	/* 容器 */
	.container {
		display: flex;
		flex-direction: column;
		background-color: var(--background-color);
		min-height: 100vh;
		padding-bottom: 80rpx;
		/* 为底部导航预留空间 */
		position: relative;
		overflow: hidden;
		/* 防止动画溢出 */
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
		opacity: 0.3;
		/* 调整透明度以不干扰内容 */
	}

	/* 头部 */
	.dec_header {
		display: flex;
		align-items: center;
		background-color: var(--card-background);
		padding: 20rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		animation: fadeInDown 1.2s cubic-bezier(0.25, 0.8, 0.25, 1);
	}

	.dec_logo {
		height: 80rpx;
		width: 60%;
		/* padding-right: 20rpx; */
	}

	.title {
		font-size: var(--font-size-title);
		font-weight: bold;
		width: 50%;
		color: var(--primary-color);
		margin-left: 20rpx;
	}

	/* 实用工具 */
	.useful-tools {
		background-color: rgba(33, 255, 6, 0.15);
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		margin: 30rpx; /* 增加工具区域外边距 */
		padding: 30rpx; /* 增加工具区域内边距 */
		animation: fadeInUp 1s ease;
		backdrop-filter: blur(10rpx);
	}

	.tools-title {
		font-size: 30rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 20rpx;
		text-align: center;
	}

	.tools-list {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}

	.tool {
		display: flex;
		align-items: center;
		background-color: rgba(255, 255, 255, 0.9);
		border-radius: 10rpx;
		padding: 20rpx; /* 增加工具项内边距 */
		margin-bottom: 20rpx; /* 保持工具项间距 */
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		cursor: pointer;
		transition: transform 0.3s, box-shadow 0.3s;
		animation: fadeInUp 1s ease;
	}

	.tool:hover {
		transform: translateY(-5rpx) scale(1.02);
		box-shadow: 0 4rpx 10rpx var(--shadow-color);
	}

	.tool-icon {
		width: 100rpx;
		height: 100rpx;
		margin-right: 20rpx;
		border-radius: 10rpx;
		object-fit: cover;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		transition: transform 0.3s;
	}

	.tool:hover .tool-icon {
		transform: rotate(5deg);
	}

	.tool-description {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.tool-name {
		font-size: 26rpx;
		color: var(--primary-color);
		font-weight: bold;
		margin-bottom: 5rpx;
	}

	.tool-info {
		font-size: 22rpx;
		color: #666;
	}
	
	.tool:active {
	    transform: scale(0.98);
	    box-shadow: 0 2rpx 8rpx var(--shadow-color);
	}


	/* 动画效果 */
	@keyframes fadeInDown {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}

		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}

		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* 响应式调整 */
	@media (min-width: 600rpx) {
		.tool {
			flex-direction: row;
		}
	}

	@media (max-width: 599rpx) {
		.tool {
			flex-direction: column;
			align-items: center;
			text-align: center;
		}

		.tool-icon {
			margin-right: 0;
			margin-bottom: 10rpx;
		}

		.tool-description {
			align-items: center;
		}
	}
</style>