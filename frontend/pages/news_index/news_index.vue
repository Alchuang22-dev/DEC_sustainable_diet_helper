<template>
	<view class="container">
		<!-- 头部导航 -->
		<view class="header">
			<button class="header-button" @click="showSection('全部')">全部</button>
			<button class="header-button" @click="showSection('环保科普')">环保科普</button>
			<button class="header-button" @click="showSection('环保要闻')">环保要闻</button>
		</view>

		<!-- 新闻列表 -->
		<scroll-view class="news-section" scroll-y :scroll-top="scrollTop">
			<view v-for="(item, index) in filteredNews" :key="index" class="news-item" @click="navigateToDetail(item)">
				<image v-if="item.image" :src="item.image" class="news-image"></image>
				<view class="news-content">
					<text class="news-title">{{ item.title }}</text>
					<text class="news-description">{{ item.category }} | {{ item.source }} | {{ item.time }}</text>
				</view>
			</view>
		</scroll-view>

	</view>
</template>

<script setup>
	import {
		ref,
		computed
	} from 'vue';

	const newsList = ref([{
			title: '国际氢能联盟和麦肯锡联合发布《氢能洞察2024》',
			category: '环保要闻',
			source: '双碳指挥',
			time: '刚刚',
			image: '',
			section: '环保要闻',
			url: '/pages/newsDetail/newsDetail?id=1',
		},
		{
			title: '把自然讲给你听 | 什么是森林？',
			category: '环保科普',
			source: '环保科普365',
			time: '1小时前',
			image: '/static/images/nature.jpg',
			section: '环保科普',
			url: '/pages/newsDetail/newsDetail?id=2',
		},
		{
			title: '视频 | 垃圾分类',
			category: '环保科普',
			source: '环保科普365',
			time: '2024-10-14',
			image: '',
			section: '环保科普',
			url: '/pages/videoDetail/videoDetail?id=3',
		},
		// 可以继续添加更多新闻项
	]);

	const currentSection = ref('全部');
	const scrollTop = ref(0);

	const showSection = (section) => {
		currentSection.value = section;
		scrollTop.value = 0; // 回到顶部
	};

	const filteredNews = computed(() => {
		if (currentSection.value === '全部') {
			return newsList.value;
		} else {
			return newsList.value.filter((item) => item.section === currentSection.value);
		}
	});

	const navigateToDetail = (item) => {
		uni.navigateTo({
			url: item.url,
		});
	};
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

	/* 头部导航 */
	.header {
		display: flex;
		align-items: center;
		padding: 20rpx;
		background-color: #ffffff;
		border-bottom: 1rpx solid var(--border-color);
		justify-content: space-around;
	}

	.header-button {
		border: none;
		background-color: transparent;
		font-size: 32rpx;
		cursor: pointer;
		color: var(--text-color);
	}

	.header-button:active {
		color: var(--primary-color);
	}

	/* 新闻列表 */
	.news-section {
		flex: 1;
		padding: 20rpx;
	}

	.news-item {
		background-color: #ffffff;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.1);
		padding: 20rpx;
		margin-bottom: 20rpx;
		display: flex;
		flex-direction: row;
		align-items: center;
		cursor: pointer;
	}

	.news-image {
		width: 160rpx;
		height: 120rpx;
		border-radius: 10rpx;
		margin-right: 20rpx;
	}

	.news-content {
		flex: 1;
	}

	.news-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--text-color);
		margin-bottom: 10rpx;
	}

	.news-description {
		font-size: 28rpx;
		color: var(--secondary-text-color);
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
</style>