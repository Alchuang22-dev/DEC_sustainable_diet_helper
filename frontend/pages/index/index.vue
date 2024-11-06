<template>
	<view class="container">
		<image src="/static/images/index/login.png" class="background-image"></image>
		<!-- 头部 -->
		<view class="header">
			<image src="/static/images/index/logo_wide.png" alt="DEC logo" class="logo"></image>
			<text class="title">欢迎来到我们的站点！</text>
		</view>

		<!-- 碳排放信息 -->
		<view class="carbon-info">
			<view class="carbon-progress">
				<text class="carbon-description">您已和我们一起保护地球</text>
				<text class="carbon-number">{{ days }}天</text>
			</view>
			<view class="charts">
				<view class="chart history">
					<text class="chart-title">碳排放历史曲线</text>
					<!-- <text>{{chartData}}</text> -->
					<qiun-data-charts :canvas2d='false' canvas-id="carbonHistoryChart" type="line"
						:chartData="chartData" />
				</view>
				<view class="chart today">
					<text class="chart-title">今日碳排放</text>
					<!-- <canvas canvas-id="carbonTodayChart" class="chart-canvas"></canvas> -->
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
				<navigator url="/pages/my/my" class="nav-item">我的</navigator>
			</view>
		</view>
	</view>
</template>


<script setup>
	import {
		ref,
		onMounted
	} from 'vue'

	// const echarts = require('../../uni_modules/lime-echart/static/echarts.min');

	// 控制查看更多内容的显示
	const showMore = ref(false)
	const days = ref(9) // 可以根据实际情况动态更新

	const toggleMoreContent = () => {
		showMore.value = !showMore.value
	}

	const chartData = {
		categories: ["七天前", "六天前", "五天前", "四天前", "三天前", "两天前", "昨天"],
		series: [{
				name: "目标值",
				data: [35, 36, 37, 38, 20, 88, 55]
			},
			{
				name: "实际值",
				data: [55, 27, 22, 45, 68, 22, 19]
			}
		]
	};

	// 页面跳转方法
	const navigateTo = (page) => {
		uni.navigateTo({
			url: `/pages/${page}/${page}`
		})
	}

	// 初始化图表
	// onMounted(() => {
	// 	// 碳排放历史曲线图
	// 	const historyCtx = uni.createCanvasContext('carbonHistoryChart', this)
	// 	const historyData = [50, 80, 65, 90, 70, 85, 60, 75, 95, 80]
	// 	const historyWidth = 400
	// 	const historyHeight = 200
	// 	const padding = 5
	// 	const maxData = Math.max(...historyData) + 10 // 增加10以避免数据点触碰边界
	// 	const step = (historyWidth - 2 * padding) / (historyData.length - 1)

	// 	// 绘制网格线
	// 	historyCtx.setStrokeStyle('#e0e0e0')
	// 	historyCtx.setLineWidth(1)
	// 	for (let i = 0; i <= 5; i++) {
	// 		let y = padding + i * (historyHeight - 2 * padding) / 5
	// 		historyCtx.moveTo(padding, y)
	// 		historyCtx.lineTo(historyWidth - padding, y)
	// 	}
	// 	historyCtx.stroke()

	// 	// 绘制坐标轴
	// 	historyCtx.setStrokeStyle('#333')
	// 	historyCtx.setLineWidth(2)
	// 	// Y轴
	// 	historyCtx.moveTo(padding, padding)
	// 	historyCtx.lineTo(padding, historyHeight - padding)
	// 	// X轴
	// 	historyCtx.lineTo(historyWidth - padding, historyHeight - padding)
	// 	historyCtx.stroke()

	// 	// 绘制数据曲线
	// 	const gradient = historyCtx.createLinearGradient(0, padding, 0, historyHeight - padding)
	// 	gradient.addColorStop(0, '#4CAF50')
	// 	gradient.addColorStop(1, '#81C784')
	// 	historyCtx.setStrokeStyle(gradient)
	// 	historyCtx.setLineWidth(3)
	// 	historyCtx.beginPath()
	// 	historyCtx.moveTo(padding, historyHeight - padding - (historyData[0] / maxData) * (historyHeight - 2 *
	// 		padding))
	// 	historyData.forEach((point, index) => {
	// 		const x = padding + step * index
	// 		const y = historyHeight - padding - (point / maxData) * (historyHeight - 2 * padding)
	// 		historyCtx.lineTo(x, y)
	// 		// 绘制数据点
	// 		historyCtx.beginPath()
	// 		historyCtx.arc(x, y, 4, 0, 2 * Math.PI)
	// 		historyCtx.setFillStyle('#4CAF50')
	// 		historyCtx.fill()
	// 		historyCtx.closePath()
	// 	})
	// 	historyCtx.stroke()

	// 	// 添加数据标签
	// 	historyCtx.setFillStyle('#333')
	// 	historyCtx.setFontSize(20)
	// 	historyData.forEach((point, index) => {
	// 		const x = padding + step * index
	// 		const y = historyHeight - padding - (point / maxData) * (historyHeight - 2 * padding) - 10
	// 		historyCtx.fillText(`${point}g`, x, y)
	// 	})

	// 	historyCtx.draw()

	// 	// 今日碳排放圆形图
	// 	const todayCtx = uni.createCanvasContext('carbonTodayChart', this)
	// 	const centerX = 300
	// 	const centerY = 150
	// 	const radius = 100
	// 	const todayEmission = 200 // 这里可以替换为动态数据
	// 	const maxEmission = 300 // 假设最大排放量

	// 	// 绘制背景圆
	// 	todayCtx.beginPath()
	// 	todayCtx.arc(centerX, centerY, radius, 0, 2 * Math.PI)
	// 	todayCtx.setFillStyle('#e0e0e0')
	// 	todayCtx.fill()
	// 	todayCtx.closePath()

	// 	// 绘制进度圆
	// 	const emissionRatio = todayEmission / maxEmission
	// 	const endAngle = emissionRatio * 2 * Math.PI
	// 	const gradientCircle = todayCtx.createLinearGradient(centerX - radius, centerY - radius, centerX + radius,
	// 		centerY + radius)
	// 	gradientCircle.addColorStop(0, '#4CAF50')
	// 	gradientCircle.addColorStop(1, '#81C784')
	// 	todayCtx.beginPath()
	// 	todayCtx.arc(centerX, centerY, radius, -Math.PI / 2, endAngle - Math.PI / 2)
	// 	todayCtx.setStrokeStyle(gradientCircle)
	// 	todayCtx.setLineWidth(15)
	// 	todayCtx.setLineCap('round')
	// 	todayCtx.stroke()
	// 	todayCtx.closePath()

	// 	// 在圆中心绘制文本
	// 	todayCtx.setFillStyle('#333')
	// 	todayCtx.setFontSize(24)
	// 	todayCtx.setTextAlign('center')
	// 	todayCtx.setTextBaseline('middle')
	// 	todayCtx.fillText(`${todayEmission}g`, centerX, centerY)
	// 	todayCtx.draw()
	// })
</script>

<style scoped>
	.container {
		display: flex;
		flex-direction: column;
		background-color: #f5f5f5;
		min-height: 100vh;
		padding-bottom: 80rpx;
		/* 为底部导航预留空间 */
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
	.header {
		display: flex;
		align-items: center;
		background-color: #ffffff;
		padding: 20rpx;
		box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.1);
	}

	.logo {
		width: 80rpx;
		height: 80rpx;
		margin-right: 20rpx;
	}

	.title {
		font-size: 32rpx;
		font-weight: bold;
		color: #4CAF50;
	}

	/* 碳排放信息 */
	.carbon-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 40rpx;
		background-color: #ffffff;
		margin: 20rpx;
		border-radius: 15rpx;
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
		width: 90%;
		/* 增加宽度 */
		max-width: 700rpx;
		/* 增加最大宽度，避免过大 */
	}

	.carbon-progress {
		text-align: center;
		margin-bottom: 30rpx;
	}

	.carbon-number {
		font-size: 60rpx;
		color: #4CAF50;
		font-weight: bold;
	}

	.carbon-description {
		font-size: 24rpx;
		color: #333;
	}

	.charts {
		display: flex;
		flex-direction: column;
		width: 100%;
		height: 300px;
		align-items: center;
	}

	.chart {
		width: 100%;
		height: 300px;
		margin-bottom: 40rpx;
	}

	.chart-title {
		text-align: center;
		margin-bottom: 15rpx;
		font-size: 28rpx;
		color: #4CAF50;
		font-weight: bold;
	}

	.chart-canvas {
		width: 100%;
		height: 100rpx;
		/* 增加高度以适应内容 */
		/* 		background-color: #ffffff;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.1); */
	}

	/* 实用工具 */
	.useful-tools {
		background-color: #ffffff;
		padding: 20rpx;
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.1);
		margin: 20rpx;
	}

	.tools-title {
		font-size: 28rpx;
		font-weight: bold;
		color: #333;
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
		box-shadow: 0 2rpx 5rpx rgba(0, 0, 0, 0.1);
	}

	.tool-description {
		flex: 1;
	}

	.tool-name {
		font-size: 24rpx;
		color: #4CAF50;
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
		color: #4CAF50;
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
		background-color: #ffffff;
		padding: 20rpx 0;
		box-shadow: 0 -2rpx 5rpx rgba(0, 0, 0, 0.1);
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
		color: #4CAF50;
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