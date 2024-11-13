<template>
	<view class="container">
		<!-- 全屏背景图片 -->
		<image src="/static/images/index/background_img.jpg" class="background-image"></image>

		<!-- 头部 -->
		<view class="dec_header">
			<image src="/static/images/index/logo_wide.png" alt="DEC logo" class="dec_logo" mode="aspectFit"></image>
			<text class="title">欢迎来到我们的站点！</text>
		</view>

		<!-- 碳排放信息 -->
		<view class="carbon-info">
			<view class="carbon-progress">
				<text class="carbon-description">您已和我们一起保护地球</text>
				<text class="carbon-number">{{ days }}天</text>
			</view>
			<view class="charts">

				<view class="chart today">
					<text class="chart-title">今日碳排放</text>
					<view class="today-charts">
						<qiun-data-charts :canvas2d="true" canvas-id="carbonTodayChart" type="ring" :opts="ringOpts"
							:chartData="chartTodayData" />
					</view>
				</view>
				<view class="chart history">
					<text class="chart-title">碳排放历史曲线</text>
					<qiun-data-charts :canvas2d="true" canvas-id="carbonHistoryChart" type="line"
						:chartData="chartHistoryData" />
				</view>
			</view>
		</view>

		<!-- 实用工具 -->
		<view class="useful-tools">
			<text class="tools-title">实用工具</text>
			<view class="tools-grid">
				<view class="tool" @click="navigateTo('calculator')" animation="fadeInUp">
					<image src="https://cdn.pixabay.com/photo/2015/12/04/17/07/co2-1076817_1280.jpg"
						alt="Carbon Calculator" class="tool-icon" mode="aspectFill"></image>
					<view class="tool-description">
						<text class="tool-name">碳计算器</text>
						<text class="tool-info">计算您的碳足迹</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('recommend')" animation="fadeInUp" animation-delay="0.2s">
					<image src="https://cdn.pixabay.com/photo/2020/03/12/18/37/dish-4925892_1280.png"
						alt="Diet Recommendation" class="tool-icon" mode="aspectFill"></image>
					<view class="tool-description">
						<text class="tool-name">饮食推荐</text>
						<text class="tool-info">优化您的饮食计划</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('nutrition')" animation="fadeInUp" animation-delay="0.4s">
					<image src="https://cdn.pixabay.com/photo/2017/07/06/17/13/calculator-2478633_1280.png"
						alt="Nutrition Calculator" class="tool-icon"></image>
					<view class="tool-description">
						<text class="tool-name">营养计算</text>
						<text class="tool-info">了解您的营养需求</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('family')" animation="fadeInUp" animation-delay="0.6s">
					<image src="https://cdn.pixabay.com/photo/2016/01/04/14/24/terminal-board-1120961_1280.png"
						alt="Family Recipe" class="tool-icon"></image>
					<view class="tool-description">
						<text class="tool-name">家庭菜谱</text>
						<text class="tool-info">规划家人饮食</text>
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

	// 控制查看更多内容的显示
	const showMore = ref(false)
	const days = ref(9) // 可以根据实际情况动态更新

	const toggleMoreContent = () => {
		showMore.value = !showMore.value
	}

	// 示例的历史数据
	const chartHistoryData = {
		categories: ["七天前", "六天前", "五天前", "四天前", "三天前", "两天前", "昨天"],
		series: [{
				name: "目标值",
				data: [40, 36, 37, 38, 20, 88, 55]
			},
			{
				name: "实际值",
				data: [39, 27, 22, 45, 68, 22, 19]
			}
		]
	};

	// 环形图配置
	const chartTodayData = {
		series: [{
			data: [{
				"name": "早餐",
				"value": 50
			}, {
				"name": "午餐",
				"value": 30
			}, {
				"name": "晚餐",
				"value": 20
			}, {
				"name": "其他",
				"value": 18
			}]
		}]
	}

	const ringOpts = {
		rotate: false,
		rotateLock: false,
		color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
		padding: [5, 5, 5, 5],
		dataLabel: true,
		enableScroll: false,
		legend: {
			show: true,
			position: "right",
			lineHeight: 25
		},
		title: {
			name: "总量",
			fontSize: 15,
			color: "#666666"
		},
		subtitle: {
			name: "5.3Kg",
			fontSize: 25,
			color: "#4CAF50"
		},
		extra: {
			ring: {
				ringWidth: 20,
				activeOpacity: 0.5,
				activeRadius: 10,
				offsetAngle: 0,
				labelWidth: 15,
				border: false,
				borderWidth: 3,
				borderColor: "#FFFFFF"
			}
		}
	}

	// 页面跳转方法
	const navigateTo = (page) => {

		if (page === 'recommend') {
			uni.navigateTo({
				url: "/pagesTool/food_recommend/food_recommend",
			})}
		if(page === 'nutrition') {
			uni.navigateTo({
				url: "/pagesTool/nutrition_calculator/nutrition_calculator",
			})}
		if(page === 'family') {
			uni.navigateTo({
				url: "/pagesTool/home_servant/home_servant",
			})}
		else{
			uni.navigateTo({
				url: "/pagesTool/carbon_calculator/carbon_calculator",
			})
		}
		
	}

	// 引入动画库（假设使用 Animate.css）
	// 如果使用其他动画库，请根据实际情况调整
	onMounted(() => {
		// 可在这里初始化动画
	})
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
		opacity: 0.1;
		/* 调整透明度以不干扰内容 */
	}

	/* 头部 */
	.dec_header {
		display: flex;
		align-items: center;
		background-color: var(--card-background);
		padding: 20rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		animation: fadeInDown 1s ease;
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

	/* 碳排放信息 */
	.carbon-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background-color: rgba(33, 255, 6, 0.1);
		max-width: 1000rpx;
		padding: 20rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		margin: 20rpx;
		backdrop-filter: blur(10rpx);
		animation: fadeInUp 1s ease;
	}

	.carbon-progress {
		text-align: center;
		margin-bottom: 30rpx;
	}

	.carbon-description {
		font-size: var(--font-size-subtitle);
		color: var(--text-color);
		padding-right: 20rpx;
	}

	.carbon-number {
		font-size: 60rpx;
		color: var(--primary-color);
		font-weight: bold;
	}

	/* 图表部分 */
	.charts {
		display: flex;
		flex-direction: column;
		width: 100%;
		align-items: center;
	}

	.chart {
		width: 100%;
		margin-bottom: 40rpx;
	}

	.chart-title {
		text-align: center;
		margin-bottom: 15rpx;
		font-size: 28rpx;
		color: var(--primary-color);
		font-weight: bold;
	}

	.today-charts {
		align-items: center;
		width: 100%;
	}

	/* 实用工具 */
	.useful-tools {
		background-color: rgba(33, 255, 6, 0.05);
		padding: 20rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		margin: 20rpx;
		animation: fadeInUp 1s ease;
		backdrop-filter: blur(10rpx);
	}

	.tools-title {
		font-size: 28rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 20rpx;
		text-align: center;
	}

	.tools-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 20rpx;
	}

	.tool {
		display: flex;
		flex-direction: column;
		align-items: center;
		background-color: rgba(255, 255, 255, 0.9);
		padding: 15rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		cursor: pointer;
		transition: transform 0.3s, box-shadow 0.3s;
		animation: fadeInUp 1s ease;
	}

	.tool:hover {
		transform: translateY(-5rpx) scale(1.05);
		box-shadow: 0 4rpx 10rpx var(--shadow-color);
	}

	.tool-icon {
		width: 120rpx;
		height: 120rpx;
		margin-bottom: 15rpx;
		border-radius: 10rpx;
		object-fit: cover;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		transition: transform 0.3s;
	}

	.tool:hover .tool-icon {
		transform: rotate(10deg);
	}

	.tool-description {
		text-align: center;
	}

	.tool-name {
		font-size: 24rpx;
		color: var(--primary-color);
		font-weight: bold;
		margin-bottom: 5rpx;
	}

	.tool-info {
		font-size: 20rpx;
		color: #666;
	}

	/* 底部导航（保留原样，如果需要可自行调整） */
	.footer {
		background-color: var(--card-background);
		padding: 20rpx 0;
		box-shadow: 0 -2rpx 5rpx var(--shadow-color);
		position: fixed;
		bottom: 0;
		width: 100%;
	}

	.footer-nav {
		display: flex;
		justify-content: space-around;
	}

	.nav-item {
		text-decoration: none;
		color: #333;
		font-weight: bold;
		transition: color 0.3s;
		font-size: 28rpx;
	}

	.nav-item:hover {
		color: var(--primary-color);
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
</style>