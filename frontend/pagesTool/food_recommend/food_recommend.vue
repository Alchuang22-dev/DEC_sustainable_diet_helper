<template>
	<view class="container">
		<!-- 头部 -->
		<view class="header">
			<text class="title">今日您的地区适宜</text>
		</view>

		<!-- 推荐区域 -->
		<view class="recommendation-section">
			<image src="/static/images/index/recommand.png" alt="蒸白菜" class="recommend-image"></image>
			<!-- <text class="section-title">今日您的地区适宜</text> -->
			<view class="dishes">
				<view class="dish" v-for="(dish, index) in dishes" :key="index">
					<image src="https://cdn.pixabay.com/photo/2016/03/17/23/30/salad-1264107_1280.jpg" alt=""
						class="dish-image"></image>
					<view class="dish-title">{{ dish }}</view>
				</view>
			</view>
			<button class="generate-button" @click="generateRecipe">生成菜谱</button>
		</view>

		<!-- 无数据提示 -->
		<!-- 		<view class="no-data" v-if="showNoData">
			<image src="/static/images/no_data.png" alt="No Data" class="no-data-image"></image>
			<text class="no-data-text">暂无推荐菜谱</text>
		</view> -->

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
	const dishes = ref(['蒸白菜', '大地瓜', '大盘鸡', '通心粉']);

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

	.title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--text-color);
		flex-grow: 1;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* 推荐区域 */
	.recommendation-section {
		text-align: center;
		padding: 30rpx;
		/* background-image: url('/static/images/background-image.jpg');
    background-size: cover; */
		color: #ffffff;
	}

	.recommend-image {
		width: 200rpx;
		height: 200rpx;
		margin-bottom: 20rpx;
		object-fit: cover;
	}

	.section-title {
		font-size: 36rpx;
		margin-bottom: 20rpx;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* 菜品宫格布局 */
	.dishes {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		/* 每行两个菜品 */
		gap: 20rpx;
		/* 网格间距 */
		justify-items: center;
		/* 网格项水平居中 */
	}

	.dish {
		width: 100%;
		/* 使菜品项填满网格单元 */
		max-width: 300rpx;
		/* 设置最大宽度，防止过大 */
		background-color: #ffffff;
		color: #333;
		border-radius: 10rpx;
		overflow: hidden;
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.dish-image {
		width: 100%;
		height: 150rpx;
		object-fit: cover;
	}

	.dish-title {
		padding: 10rpx;
		font-size: 28rpx;
		background-color: #ffe082;
		text-align: center;
		font-weight: bold;
		width: 100%;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
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
		object-fit: cover;
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
		object-fit: cover;
	}

	.box-description {
		flex-grow: 1;
	}

	.box-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 10rpx;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.box-text {
		font-size: 28rpx;
		color: var(--text-color);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* 响应式设计 */
	@media (max-width: 600px) {
		.title {
			font-size: 28rpx;
		}

		.section-title,
		.box-title {
			font-size: 32rpx;
		}

		.dish-title {
			font-size: 24rpx;
		}

		.generate-button {
			font-size: 28rpx;
			padding: 15rpx 30rpx;
		}

		.box-image {
			width: 140rpx;
			height: 140rpx;
			margin-right: 20rpx;
		}

		/* 调整网格间距 */
		.dishes {
			gap: 15rpx;
		}
	}
</style>