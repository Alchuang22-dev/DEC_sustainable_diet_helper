<template>
	<view class="container">
		<!-- 头部 -->
		<view class="header">
			<text class="title">今日您的地区适宜</text>
		</view>

		<!-- 推荐区域 -->
		<view class="recommendation-section">
			<image src="/static/images/index/recommand.png" alt="蒸白菜" class="recommend-image"></image>
			<text class="section-title">今日您的地区适宜</text>
			<view class="dishes">
				<view class="dish" v-for="(dish, index) in dishes" :key="index">
					<!-- 如果有图片，可以添加 <image> 标签 -->
					<img src="https://cdn.pixabay.com/photo/2015/05/16/15/03/tomatoes-769999_1280.jpg" alt="" />
					<view class="dish-title">{{ dish }}</view>
				</view>
			</view>
			<button class="generate-button" @click="generateRecipe">生成菜谱</button>
		</view>

		<!-- 无数据提示 -->
		<view class="no-data" v-if="showNoData">
			<image src="/static/images/no_data.png" alt="No Data" class="no-data-image"></image>
			<text class="no-data-text">暂无推荐菜谱</text>
		</view>

		<!-- 推荐菜谱 -->
		<view class="recipe-boxes" v-if="showRecipeBoxes">
			<view class="box" @click="goToRecipe('dapanji')">
				<image src="/static/images/dapanji.png" alt="大盘鸡" class="box-image"></image>
				<view class="box-description">
					<text class="box-title">推荐菜谱</text>
					<text class="box-text">根据您所在的地区，推荐大盘鸡！</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref
	} from 'vue';

	// 响应式数据
	const showNoData = ref(true);
	const showRecipeBoxes = ref(false);
	const dishes = ref(['蒸白菜', '大地锅', '大盘鸡', '通心粉']);

	// 方法
	const navigateBack = () => {
		uni.navigateBack();
	};

	const generateRecipe = () => {
		showNoData.value = false;
		showRecipeBoxes.value = true;
	};

	const goToRecipe = (recipeName) => {
		// 跳转到对应的菜谱页面
		uni.navigateTo({
			url: `/pages/recipes/${recipeName}`,
		});
	};
</script>

<style scoped>
	/* 全局样式变量 */
	:root {
		--primary-color: #4caf50;
		--secondary-color: #8bc34a;
		--text-color: #333;
		--background-color: #f5f5f5;
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

	/* 头部 */
	.header {
		display: flex;
		align-items: center;
		padding: 20rpx;
		background-color: #ffffff;
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
	}

	.back-button {
		font-size: 36rpx;
		margin-right: 20rpx;
		color: var(--primary-color);
		cursor: pointer;
	}

	.title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--text-color);
		flex-grow: 1;
	}

	.home-button {
		font-size: 36rpx;
		color: var(--primary-color);
		cursor: pointer;
	}

	/* 推荐区域 */
	.recommendation-section {
		text-align: center;
		padding: 30rpx;
		/* 		background-image: url('/static/images/background-image.jpg');
		background-size: cover; */
		color: #ffffff;
	}

	.recommend-image {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 20rpx;
	}

	.section-title {
		font-size: 48rpx;
		margin-bottom: 20rpx;
	}

	.dishes {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
	}

	.dish {
		width: 200rpx;
		margin: 10rpx;
		background-color: #ffffff;
		color: #333;
		border-radius: 10rpx;
		overflow: hidden;
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
	}

	.dish-title {
		padding: 10rpx;
		font-size: 28rpx;
		background-color: #ffe082;
		text-align: center;
		font-weight: bold;
	}

	/* 生成菜谱按钮 */
	.generate-button {
		background-color: var(--primary-color);
		color: #ffffff;
		padding: 20rpx 40rpx;
		border: none;
		border-radius: 30rpx;
		font-size: 32rpx;
		cursor: pointer;
		margin: 30rpx auto 0;
	}

	/* 无数据提示 */
	.no-data {
		text-align: center;
		margin-top: 60rpx;
	}

	.no-data-image {
		width: 200rpx;
		height: 200rpx;
	}

	.no-data-text {
		margin-top: 20rpx;
		font-size: 28rpx;
		color: #999;
	}

	/* 推荐菜谱 */
	.recipe-boxes {
		background-color: #ffffff;
		padding: 30rpx;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
		margin: 30rpx 20rpx;
	}

	.box {
		display: flex;
		align-items: center;
		cursor: pointer;
	}

	.box-image {
		width: 160rpx;
		height: 160rpx;
		margin-right: 30rpx;
		border-radius: 10rpx;
	}

	.box-description {
		flex-grow: 1;
	}

	.box-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 10rpx;
	}

	.box-text {
		font-size: 28rpx;
		color: var(--text-color);
	}
</style>