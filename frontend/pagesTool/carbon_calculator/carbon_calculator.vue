<template>
	<view class="container" @load="handleLoad">
		<!-- 全屏背景图片 -->
		<image src="../static/background_img.jpg" class="background-image"></image>

		<!-- 头部标题 -->
		<view class="header">
			<text class="header-title">碳足迹计算器</text>
		</view>

		<!-- 已添加的食物标题 -->
		<text class="list-title">您已添加的食物</text>

		<!-- 可滑动的食物列表 -->
		<scroll-view scroll-y="true" class="food-list scroll-view">
			<view v-for="(food, index) in foodList" :key="index" class="card-container">
				<uni-card :title="food.name || '西红柿'"
					:thumbnail="food.image || 'https://cdn.pixabay.com/photo/2015/05/16/15/03/tomatoes-769999_1280.jpg'"
					:sub-title="`重量: ${food.weight || '1.2kg'} 价格: ${food.price || '5元'}`" shadow=1
					@click="animateCard(index)" :class="{ clicked: food.isAnimating }"
					:extra="`${food.transportMethod} ${food.foodSource}`"
					:style="{ animationDelay: `${index * 0.1}s` }">

					<view class="card-actions">
						<button class="delete-button" @click="deleteFood(index)">删除</button>
						<button class="edit-button">修改</button>
					</view>
				</uni-card>
			</view>
		</scroll-view>

		<!-- 按钮区 -->
		<view class="button-group">
			<button class="primary-button small-button" @click="navigateTo('add_food')">
				添加食品
			</button>
			<button class="secondary-button small-button" @click="saveData">
				保存添加
			</button>
			<button class="calculate-button small-button" @click="calculateEmission">
				计算碳排放
			</button>
		</view>

		<!-- 碳排放结果环形图 -->
		<view class="result" v-if="showResult">
			<text class="result-title">您本次的碳足迹是：</text>
			<qiun-data-charts :canvas2d="true" canvas-id="carbonEmissionChart" type="ring" :opts="ringOpts"
				:chartData="chartEmissionData" />
			<button class="save-button" @click="saveEmissionData">保存</button>
		</view>


		<!-- 实用工具 -->
		<view class="useful-tools">
			<view class="tool" @click="navigateTo('recommendMenu')">
				<image src="../static/toufu.png" class="tool-image" alt="推荐菜单"></image>
				<view class="tool-description">
					<text class="tool-title">推荐菜单</text>
					<text class="tool-text">试试我们的可持续推荐菜单！</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref,
		reactive,
		onMounted
	} from 'vue';
	import {
		onLoad
	} from '@dcloudio/uni-app';

	// 食物列表
	const foodList = reactive([{
			name: "西红柿",
			weight: "1kg",
			price: "5元",
			transportMethod: "陆运",
			foodSource: "本地",
			image: "",
			isAnimating: false, // 用于跟踪动画状态
			emission: 0
		},
		{
			name: "西红柿",
			weight: "1kg",
			price: "5元",
			transportMethod: "陆运",
			foodSource: "本地",
			image: "",
			isAnimating: false,
			emission: 0
		},
		{
			name: "西红柿",
			weight: "1kg",
			price: "5元",
			transportMethod: "陆运",
			foodSource: "本地",
			image: "",
			isAnimating: false,
			emission: 0
		},
		{
			name: "西红柿",
			weight: "1kg",
			price: "5元",
			transportMethod: "陆运",
			foodSource: "本地",
			image: "",
			isAnimating: false,
			emission: 0
		},
		// 可以预先添加更多食物项
	]);

	// 碳排放数据，仅包含CO2
	const emission = ref({
		CO2: 0,
	});

	const showResult = ref(false);

	// 环形图数据和配置
	const chartEmissionData = ref({
		series: [{
			name: "CO₂排放",
			data: []
		}]
	});

	const ringOpts = ref({
		rotate: false,
		rotateLock: false,
		color: ["#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0", "#9966FF"],
		// padding: [10, 10, 10, 10],
		dataLabel: true,
		enableScroll: false,
		legend: {
			show: true,
			position: "right",
			lineHeight: 25
		},
		title: {
			name: "总排放量",
			fontSize: 15,
			color: "#666666"
		},
		subtitle: {
			name: "", // 中心显示总排放量
			fontSize: 25,
			color: "#4CAF50"
		},
		extra: {
			ring: {
				ringWidth: 10,
				activeOpacity: 0.5,
				activeRadius: 20,
				offsetAngle: 0,
				labelWidth: 15,
				border: false,
				borderWidth: 3,
				borderColor: "#FFFFFF"
			}
		}
	});

	// 页面加载时处理动画
	const handleLoad = () => {
		foodList.forEach((food, index) => {
			setTimeout(() => {
				food.isAnimating = true;
				setTimeout(() => {
					food.isAnimating = false;
				}, 500);
			}, index * 100);
		});
	};

	// 保存数据到本地存储
	const saveData = () => {
		uni.setStorageSync('foodDetails', foodList);
		uni.showToast({
			title: '保存成功',
			icon: 'success',
			duration: 2000,
		});
	};

	// 删除食物项
	const deleteFood = (index) => {
		foodList.splice(index, 1);
		uni.setStorageSync('foodDetails', foodList);
	};

	// 计算碳排放
	const calculateEmission = () => {
		// 模拟向后端发送请求
		uni.request({
			url: 'https://mock-api.com/calculateEmission', // 模拟的后端接口URL
			method: 'POST',
			data: {
				foodList: foodList.map(food => ({
					name: food.name,
					weight: food.weight,
					// 其他需要发送的字段
				}))
			},
			success: (res) => {
				// 模拟后端返回的数据
				const mockResponse = {
					data: [{
							name: "西红柿",
							emission: 2
						},
						{
							name: "苹果",
							emission: 3
						},
						{
							name: "牛肉",
							emission: 10
						},
						{
							name: "豆腐",
							emission: 1.5
						},
						// 根据foodList的实际内容添加更多食物项
					]
				};

				// 假设后端返回的数据结构与mockResponse相同
				const emissionData = mockResponse.data;

				// 更新foodList中的每个食物项，添加emission字段
				emissionData.forEach((item, index) => {
					if (foodList[index]) {
						foodList[index].emission = item.emission;
					}
				});

				// 更新环形图的数据和总排放量
				let totalCO2 = 0;
				chartEmissionData.value.series[0].data = emissionData.map(item => {
					totalCO2 += item.emission;
					return {
						name: item.name,
						value: item.emission
					};
				});

				// 更新环形图中心显示的总排放量
				ringOpts.value.subtitle.name = `${totalCO2} kg`;

				// 显示结果
				showResult.value = true;

				// 初始化并绘制环形图
				uni.createSelectorQuery().select('#carbonEmissionChart').fields({
					node: true,
					size: true
				}, (res) => {
					const canvas = res.node;
					const ctx = canvas.getContext('2d');
					const chart = new qCharts({
						canvas: ctx,
						type: 'ring',
						data: chartEmissionData.value,
						options: ringOpts.value
					});
					chart.draw();
				}).exec();
			},
			fail: (err) => {
				// 模拟后端返回的数据
				const mockResponse = {
					data: [{
							name: "西红柿",
							emission: 2
						},
						{
							name: "苹果",
							emission: 3
						},
						{
							name: "牛肉",
							emission: 10
						},
						{
							name: "豆腐",
							emission: 1.5
						},
						// 根据foodList的实际内容添加更多食物项
					]
				};

				// 假设后端返回的数据结构与mockResponse相同
				const emissionData = mockResponse.data;

				// 更新foodList中的每个食物项，添加emission字段
				emissionData.forEach((item, index) => {
					if (foodList[index]) {
						foodList[index].emission = item.emission;
					}
				});

				// 更新环形图的数据和总排放量
				let totalCO2 = 0;
				chartEmissionData.value.series[0].data = emissionData.map(item => {
					totalCO2 += item.emission;
					return {
						name: item.name,
						value: item.emission
					};
				});

				// 更新环形图中心显示的总排放量
				ringOpts.value.subtitle.name = `${totalCO2} kg`;

				// 显示结果
				showResult.value = true;

				// 初始化并绘制环形图
				uni.createSelectorQuery().select('#carbonEmissionChart').fields({
					node: true,
					size: true
				}, (res) => {
					const canvas = res.node;
					const ctx = canvas.getContext('2d');
					const chart = new qCharts({
						canvas: ctx,
						type: 'ring',
						data: chartEmissionData.value,
						options: ringOpts.value
					});
					chart.draw();
				}).exec();
			},
			// TODO: 修改请求失败的逻辑为处理失败逻辑，现在只是为了数据能显示。
			// fail: (err) => {
			// 	// 处理请求失败的情况
			// 	console.error('请求失败', err);
			// 	uni.showToast({
			// 		title: '计算失败，请稍后再试',
			// 		icon: 'none',
			// 		duration: 2000,
			// 	});
			// }
		});

		// 保存碳排放数据到后端
		const saveEmissionData = () => {
			uni.request({
				url: 'https://mock-api.com/saveEmissionData', // 模拟的后端保存接口URL
				method: 'POST',
				data: {
					foodList: foodList.map(food => ({
						name: food.name,
						weight: food.weight,
						price: food.price,
						transportMethod: food.transportMethod,
						foodSource: food.foodSource,
						image: food.image,
						emission: food.emission || 0, // 添加emission字段，默认为0
					}))
				},
				success: (res) => {
					uni.showToast({
						title: '保存成功',
						icon: 'success',
						duration: 2000,
					});
				},
				fail: (err) => {
					console.error('保存失败', err);
					uni.showToast({
						title: '保存失败，请稍后再试',
						icon: 'none',
						duration: 2000,
					});
				}
			});
		};
	};


	// 页面跳转方法
	const navigateTo = (page) => {
		if (page === 'add_food') {
			uni.navigateTo({
				url: '/pagesTool/add_food/add_food',
			});
		} else if (page === 'recommendMenu') {
			uni.navigateTo({
				url: '/pages/recommendMenu/recommendMenu',
			});
		}
	};

	const animateCard = (index) => {
		const food = foodList[index];
		if (!food) return;

		food.isAnimating = true;

		setTimeout(() => {
			food.isAnimating = false;
		}, 300);
	};

	// 页面加载时执行
	// onLoad(() => {
	// 	let storedFoodList = uni.getStorageSync('foodDetails');
	// 	if (!storedFoodList || storedFoodList.length === 0) {
	// 		// 将初始的 foodList 保存到本地存储
	// 		uni.setStorageSync('foodDetails', foodList);
	// 		storedFoodList = foodList;
	// 	} else {
	// 		// 用本地存储的数据更新 foodList
	// 		foodList.splice(0, foodList.length, ...storedFoodList.map(food => ({
	// 			...food,
	// 			isAnimating: false,
	// 		})));
	// 	}
	// 	console.log(foodList)
	// 	handleLoad();
	// });
</script>

<style scoped>
	/* 全局样式变量 */
	:root {
		--primary-color: #4CAF50;
		--secondary-color: #8BC34A;
		--accent-color: #FF9800;
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
		padding-bottom: 50rpx;
		animation: fadeIn 1s ease-in-out;
	}

	/* 全屏背景图片 */
	.background-image {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
		z-index: -1;
		opacity: 0.05;
	}

	/* 头部标题 */
	.header {
		padding: 40rpx 20rpx 20rpx;
		text-align: center;
	}

	.header-title {
		font-size: 48rpx;
		color: var(--primary-color);
		font-weight: bold;
		animation: slideDown 1s ease-out;
	}

	/* 已添加的食物列表 */
	.food-list {
		max-height: 600rpx;
		margin: 20rpx 0rpx;
		padding: 20rpx 0rpx;
		background-color: #ffffff;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
		overflow: hidden;
	}

	.list-title {
		margin-left: 10rpx;
		font-size: 30rpx;
		font-weight: bold;
		color: var(--text-color);
		margin-bottom: 20rpx;
		text-align: center;
	}

	/* uni-card 相关样式 */
	.card-actions {
		display: flex;
		flex-direction: row;
		justify-content: flex-start;
		position: relative;
		width: 100%;
	}

	.delete-button,
	.edit-button {
		font-size: 18rpx;
		cursor: pointer;
		transition: color 0.3s ease;
	}

	.edit-button:hover {
		color: #8BC34A;
	}

	.delete-button:hover {
		color: #f44336;
	}

	/* 按钮区 */
	.button-group {
		display: flex;
		justify-content: flex-start;
		margin: 20rpx 20rpx;
		gap: 20rpx;
	}

	.small-button {
		padding: 15rpx 30rpx;
		border-radius: 30rpx;
		border: none;
		font-size: 20rpx;
		color: #ffffff;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.primary-button {
		background-color: var(--primary-color);
	}

	.primary-button:hover {
		transform: translateY(-2rpx);
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
	}

	.secondary-button {
		background-color: var(--secondary-color);
	}

	.secondary-button:hover {
		transform: translateY(-2rpx);
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
	}

	/* 计算部分 */
	.calculate-section {
		text-align: center;
		margin: 10rpx 0;
	}

	.calculate-button {
		border-radius: 30rpx;
		background-color: var(--accent-color);
		border: none;
		color: #ffffff;
		font-size: 20rpx;
		cursor: pointer;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.calculate-button:hover {
		transform: translateY(-2rpx);
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
	}

	/* 结果显示（环形图） */
	.result {
		position: relative;
		margin: 20rpx 20rpx;
		padding: 20rpx 30rpx;
		background-color: #ffffff;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
		font-size: 32rpx;
		color: var(--text-color);
		text-align: center;
		animation: fadeIn 1s ease-in-out;
	}

	.result-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 15rpx;
	}

	/* 保存按钮样式 */
	.save-button {
		position: absolute;
		bottom: 10rpx;
		right: 10rpx;
		padding: 10rpx 20rpx;
		background-color: var(--primary-color);
		color: #ffffff;
		border: none;
		border-radius: 20rpx;
		font-size: 24rpx;
		cursor: pointer;
		transition: background-color 0.3s ease, transform 0.2s ease;
	}

	.save-button:hover {
		background-color: var(--secondary-color);
		transform: translateY(-2rpx);
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
	}


	/* 实用工具 */
	.useful-tools {
		background-color: #ffffff;
		padding: 20rpx 30rpx;
		border-radius: 20rpx;
		box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
		margin: 20rpx;
		animation: fadeInUp 1s ease-out;
	}

	.tool {
		display: flex;
		align-items: center;
		cursor: pointer;
		transition: transform 0.3s ease;
	}

	.tool:hover {
		transform: translateY(-5rpx);
	}

	.tool-image {
		width: 140rpx;
		height: 140rpx;
		margin-right: 20rpx;
		border-radius: 10rpx;
		object-fit: cover;
	}

	.tool-description {
		display: flex;
		flex-direction: column;
	}

	.tool-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 5rpx;
	}

	.tool-text {
		font-size: 28rpx;
		color: var(--text-color);
	}

	/* 动画 */
	@keyframes fadeIn {
		from {
			opacity: 0;
		}

		to {
			opacity: 1;
		}
	}

	@keyframes slideDown {
		from {
			transform: translateY(-20rpx);
			opacity: 0;
		}

		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	@keyframes popIn {
		0% {
			transform: scale(0.95);
			opacity: 0;
		}

		60% {
			transform: scale(1.05);
			opacity: 1;
		}

		100% {
			transform: scale(1);
		}
	}

	@keyframes fadeInUp {
		from {
			transform: translateY(20rpx);
			opacity: 0;
		}

		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	/* 点击动画效果 */
	@keyframes clickEffect {
		0% {
			transform: scale(1);
		}

		50% {
			transform: scale(1.05);
		}

		100% {
			transform: scale(1);
		}
	}
</style>