<template>
	<view class="container">
		<!-- <image src="/static/images/index/login.png" class="background-image"></image> -->

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
						<!-- 						<qiun-data-charts :canvas2d="true" canvas-id="carbonPieChart" type="pie" :opts="pieOpts"
							:chartData="chartPieData" /> -->
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
			<view class="tool" @click="navigateTo('calculator')">
				<image src="/static/images/index/caculator.png" alt="Carbon Calculator" class="tool-icon"></image>
				<view class="tool-description">
					<text class="tool-name">碳计算器</text>
					<text class="tool-info">计算您的碳足迹</text>
				</view>
			</view>
			<view class="tool" @click="navigateTo('recommend')">
				<image src="/static/images/index/food.png" alt="Diet Recommendation" class="tool-icon"></image>
				<view class="tool-description">
					<text class="tool-name">饮食推荐</text>
					<text class="tool-info">优化您的饮食计划</text>
				</view>
			</view>
			<text class="view-more" @click="toggleMoreContent">
				{{ showMore ? '收起' : '查看更多' }}
			</text>
			<view class="view-more-content" v-if="showMore">
				<view class="tool" @click="navigateTo('nutrition')">
					<image src="/static/images/index/nutrition.png" alt="Nutrition Calculator" class="tool-icon">
					</image>
					<view class="tool-description">
						<text class="tool-name">营养计算</text>
						<text class="tool-info">了解您的营养需求</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('family')">
					<image src="/static/images/index/family.jpg" alt="Family Recipe" class="tool-icon"></image>
					<view class="tool-description">
						<text class="tool-name">家庭菜谱</text>
						<text class="tool-info">规划家人饮食</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 底部导航 -->
		<view class="footer">
			<view class="footer-nav">
				<navigator url="/pages/index/index" class="nav-item">主页</navigator>
				<navigator url="/pages/news/news" class="nav-item">资讯</navigator>
				<navigator url="/pagesMy/my_index/my_index" class="nav-item">我的</navigator>
			</view>
		</view>
	</view>
</template>



<script setup>
	import {
		ref
	} from 'vue'

	// 控制查看更多内容的显示
	const showMore = ref(false)
	const days = ref(9) // 可以根据实际情况动态更新

	const toggleMoreContent = () => {
		showMore.value = !showMore.value
	}

	// 示例的历史数据
	const chartHistoryData = {
		categories: ["六天前", "五天前", "四天前", "三天前", "两天前", "昨天"],
		series: [{
				name: "目标值",
				data: [36, 37, 38, 20, 88, 55]
			},
			{
				name: "实际值",
				data: [27, 22, 45, 68, 22, 19]
			}
		]
	};

	// 示例仪表盘配置
	// const gaugeOpts = {
	// 	color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
	// 	padding: undefined,
	// 	title: {
	// 		name: "3.6Kg",
	// 		fontSize: 25,
	// 		color: "#2fc25b",
	// 		offsetY: 50
	// 	},
	// 	subtitle: {
	// 		name: "总额",
	// 		fontSize: 15,
	// 		color: "#666666",
	// 		offsetY: -50,
	// 	},
	// 	extra: {
	// 		gauge: {
	// 			type: "default",
	// 			width: 50,
	// 			labelColor: "#666666",
	// 			startAngle: 0.75,
	// 			endAngle: 0.25,
	// 			startNumber: 0,
	// 			endNumber: 20,
	// 			labelFormat: "",
	// 			splitLine: {
	// 				fixRadius: 0,
	// 				splitNumber: 10,
	// 				width: 30,
	// 				color: "#FFFFFF",
	// 				childNumber: 5,
	// 				childWidth: 12
	// 			},
	// 			pointer: {
	// 				width: 18,
	// 				color: "auto"
	// 			}
	// 		}
	// 	}
	// }

	// // 示例今日碳排放数据
	// const chartTodayData = {
	// 	categories: [{
	// 		"value": 0.2,
	// 		"color": "#1890ff"
	// 	}, {
	// 		"value": 0.8,
	// 		"color": "#2fc25b"
	// 	}, {
	// 		"value": 1,
	// 		"color": "#f04864"
	// 	}],
	// 	series: [{
	// 		name: "完成率",
	// 		data: 0.81
	// 	}]
	// };

	// 新增饼图数据和配置
	// const chartPieData = {
	// 	categories: ["早餐", "中餐", "晚餐", "其他"],
	// 	series: [{
	// 			name: "早餐",
	// 			data: 25
	// 		},
	// 		{
	// 			name: "中餐",
	// 			data: 35
	// 		},
	// 		{
	// 			name: "晚餐",
	// 			data: 30
	// 		},
	// 		{
	// 			name: "其他",
	// 			data: 10
	// 		}
	// 	]
	// };

	// const pieOpts = {
	// 	color: ["#FF6384", "#36A2EB", "#FFCE56", "#8A2BE2"],
	// 	legend: {
	// 		position: 'bottom'
	// 	},
	// 	tooltip: {
	// 		show: true
	// 	}
	// };

	// 环形图配置
	// TODO: 请求数据，数据格式和下面的类似

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
		uni.navigateTo({
			url: `/pages/${page}/${page}`
		})
	}
</script>


<style scoped>
	/* 通用变量 */
	:root {
		--primary-color: #4CAF50;
		--secondary-color: #2fc25b;
		--background-color: #f5f5f5;
		--card-background: #ffffff;
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
	}

	/* 全屏背景图片 */
	.background-image {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		z-index: -1;
		/* 将背景图片置于最底层 */
	}

	/* 头部 */
	.dec_header {
		display: flex;
		align-items: center;
		background-color: var(--card-background);
		padding: 20rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
	}

	.dec_logo {
		/* 		width: 200rpx;
		height: 200rpx; */
		height: 80rpx;
		width: 60%;
		/* margin-right: 10rpx; */
	}

	.title {
		font-size: var(--font-size-title);
		font-weight: bold;
		color: var(--primary-color);
	}

	/* 碳排放信息 */
	.carbon-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background-color: var(--card-background);
		/* width: 90%; */
		max-width: 1000rpx;
		padding: 20rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		margin: 20rpx;
		/* 增加最大宽度 */
	}

	.carbon-progress {
		text-align: center;
		margin-bottom: 30rpx;
	}

	.carbon-description {
		font-size: var(--font-size-subtitle);
		color: var(--text-color);
		padding-right: 20rpx;
		/* align-self: center; */
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
		/* display: flex; */
		/* justify-content: space-between; */
		align-items: center;
		width: 100%;
		/* height: 600rpx; */
	}

	/* 	.today-charts qiun-data-charts {
		width: 100%;
	} */

	/* 实用工具 */
	.useful-tools {
		background-color: var(--card-background);
		padding: 20rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		margin: 20rpx;
	}

	.tools-title {
		font-size: 28rpx;
		font-weight: bold;
		color: var(--text-color);
		margin-bottom: 20rpx;
		text-align: center;
	}

	.tool {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
		cursor: pointer;
		transition: transform 0.2s;
	}

	.tool:hover {
		transform: scale(1.05);
	}

	.tool-icon {
		width: 120rpx;
		height: 120rpx;
		margin-right: 20rpx;
		border-radius: 10rpx;
		object-fit: cover;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
	}

	.tool-description {
		flex: 1;
	}

	.tool-name {
		font-size: 24rpx;
		color: var(--primary-color);
		font-weight: bold;
	}

	.tool-info {
		margin-top: 5rpx;
		font-size: 20rpx;
		color: #666;
	}

	.view-more {
		display: block;
		text-align: center;
		color: var(--primary-color);
		font-weight: bold;
		margin-top: 20rpx;
		cursor: pointer;
		font-size: 24rpx;
		transition: color 0.3s;
	}

	.view-more:hover {
		color: #388E3C;
	}

	.view-more-content {
		padding-top: 20rpx;
	}

	.view-more-content .tool {
		margin-bottom: 20rpx;
	}

	/* 底部导航 */
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

	/* 响应式调整 */
	@media screen and (max-width: 600px) {
		.header {
			padding: 15rpx;
		}

		.logo {
			width: 60rpx;
			height: 60rpx;
			margin-right: 15rpx;
		}

		.title {
			font-size: 28rpx;
		}

		.carbon-number {
			font-size: 50rpx;
		}

		.carbon-description {
			font-size: 22rpx;
		}

		.today-charts {
			flex-direction: column;
			height: auto;
		}

		.today-charts qiun-data-charts {
			width: 100%;
			height: 250rpx;
			margin-bottom: 20rpx;
		}

		.chart-canvas {
			width: 500rpx;
			height: 250rpx;
		}

		.tools-title {
			font-size: 24rpx;
		}

		.tool-icon {
			width: 100rpx;
			height: 100rpx;
			margin-right: 15rpx;
		}

		.tool-name {
			font-size: 22rpx;
		}

		.tool-info {
			font-size: 18rpx;
		}

		.view-more {
			font-size: 22rpx;
		}

		.nav-item {
			font-size: 24rpx;
		}
	}
</style>