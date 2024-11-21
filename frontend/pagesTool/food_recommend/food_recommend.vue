<template>
	<view class="container">
		<!-- 全屏背景图片 -->
		<image src="/static/images/index/background_img.jpg" class="background-image"></image>

		<!-- 推荐区域 -->
		<view class="recommendation-section">
			<!-- 替换图片为文字 -->
			<text class="recommend-title">{{ $t('recommendation_title') }}</text>

			<!-- 推荐菜品列表 -->
			<view class="dishes">
				<view class="dish" v-for="(dish, index) in dishes" :key="index"
					:class="'fade-in-up delay-' + (index + 1)">
					<image :src="dish.image" :alt="dish.name" class="dish-image"></image>
					<view class="dish-title">{{ dish.name }}</view>
				</view>
			</view>

			<!-- 生成菜谱按钮 -->
			<button class="generate-button fade-in-up delay-5" @click="generateRecipe">
				{{$t('generate_recipe')}}
			</button>
		</view>


		<!-- 推荐菜谱 -->
		<view class="recipe-boxes" v-if="showRecipeBoxes">
			<view class="box fade-in-up delay-6" @click="goToRecipe('dapanji')">
				<image src="/static/images/dapanji.png" alt="大盘鸡" class="box-image"></image>
				<view class="box-description">
					<text class="box-title">{{$t('recommended_recipe')}}</text>
					<text class="box-text">{{$t('recommended_recipe_info')}}</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref
	} from 'vue'
	import {
		useI18n
	} from 'vue-i18n'

	// 初始化 i18n
	const {
		t
	} = useI18n()

	// 响应式数据
	const showRecipeBoxes = ref(false)
	const dishes = ref([{
			name: t('dish_1'),
			image: 'https://cdn.pixabay.com/photo/2016/03/17/23/30/salad-1264107_1280.jpg'
		},
		{
			name: t('dish_2'),
			image: 'https://cdn.pixabay.com/photo/2016/03/05/19/02/dish-1238243_1280.jpg'
		},
		{
			name: t('dish_3'),
			image: 'https://cdn.pixabay.com/photo/2016/11/18/17/44/chicken-1835703_1280.jpg'
		},
		{
			name: t('dish_4'),
			image: 'https://cdn.pixabay.com/photo/2016/11/18/14/40/pasta-1836457_1280.jpg'
		}
	])

	// 方法
	const generateRecipe = () => {
		showRecipeBoxes.value = true
	}

	const goToRecipe = (recipeName) => {
		// 跳转到对应的菜谱页面
		uni.navigateTo({
			url: `/pages/recipes/${recipeName}`,
		})
	}
</script>

<style scoped>
	/* 通用变量 */
	:root {
		--primary-color: #4CAF50;
		--secondary-color: #2fc25b;
		--background-color: #f5f5f5;
		--card-background: rgba(255, 255, 255, 0.8);
		--text-color: #333;
		--shadow-color: rgba(0, 0, 0, 0.1);
		--font-size-title: 32rpx;
		--font-size-subtitle: 24rpx;
		--transition-duration: 0.5s;
	}

	/* 容器 */
	.container {
		display: flex;
		flex-direction: column;
		background-color: var(--background-color);
		min-height: 100vh;
		padding-bottom: 80rpx;
		position: relative;
		overflow: hidden;
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
		opacity: 0.1;
	}

	/* 推荐区域 */
	.recommendation-section {
		display: flex;
		flex-direction: column;
		align-items: center;
		background-color: rgba(76, 175, 80, 0.1);
		/* 半透明绿色背景 */
		backdrop-filter: blur(2rpx);
		/* 高斯模糊 */
		padding: 30rpx;
		margin: 20rpx;
		border-radius: 15rpx;
		box-shadow: 0 4rpx 10rpx var(--shadow-color);
		z-index: 1;
		position: relative;
		animation: fadeInUp 1s ease;
	}

  .recommend-title {
		text-align: center;
		margin-bottom: 15rpx;
		font-size: 28rpx;
		color: var(--primary-color);
		font-weight: bold;
	}

	/* 菜品宫格布局 */
	.dishes {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 20rpx;
		justify-items: center;
		width: 100%;
	}

	.dish {
		width: 100%;
		max-width: 300rpx;
		background-color: rgba(255, 255, 255, 0.9);
		color: #333;
		border-radius: 10rpx;
		overflow: hidden;
		box-shadow: 0 4rpx 10rpx var(--shadow-color);
		display: flex;
		flex-direction: column;
		align-items: center;
		opacity: 0;
		transform: translateY(20px);
		animation: fadeInUp 0.5s forwards;
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
		opacity: 0;
		transform: translateY(20px);
		animation: fadeInUp 0.5s forwards;
	}

	/* 推荐菜谱 */
	.recipe-boxes {
		background-color: rgba(255, 255, 255, 0.9);
		padding: 30rpx;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
		margin: 30rpx 20rpx;
		animation: fadeInUp 1s ease;
	}

	.box {
		display: flex;
		align-items: center;
		cursor: pointer;
		opacity: 0;
		transform: translateY(20px);
		animation: fadeInUp 0.5s forwards;
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
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* 动画延迟 */
	.delay-1 {
		animation-delay: 0.3s;
	}

	.delay-2 {
		animation-delay: 0.6s;
	}

	.delay-3 {
		animation-delay: 0.9s;
	}

	.delay-4 {
		animation-delay: 1.2s;
	}

	.delay-5 {
		animation-delay: 1.5s;
	}

	.delay-6 {
		animation-delay: 1.8s;
	}

	/* 响应式设计 */
	@media (max-width: 600px) {
		.dec_header {
			flex-direction: column;
			align-items: center;
		}

		.dec_logo {
			width: 80%;
			margin-bottom: 10rpx;
		}

		.title {
			width: 100%;
			text-align: center;
			margin-left: 0;
		}

		.recommendation-section {
			padding: 20rpx;
		}

		.recommend-image {
			width: 150rpx;
			height: 150rpx;
		}

		.dish-title {
			font-size: 24rpx;
		}

		.generate-button {
			font-size: 28rpx;
			padding: 15rpx 30rpx;
		}

		.box-title {
			font-size: 32rpx;
		}

		.box-text {
			font-size: 24rpx;
		}

		.dishes {
			gap: 15rpx;
		}
	}
</style>