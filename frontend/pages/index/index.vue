<template>
	<view class="container">
		<!-- 全屏背景图片 -->
		<image src="/static/images/index/background_index_new.png" class="background-image"></image>

		<!-- 头部 -->
		<view class="dec_header">
			<image src="/static/images/index/logo_wide.png" :alt="$t('dec_logo_alt')" class="dec_logo" mode="aspectFit">
			</image>
			<text class="title">{{$t('welcome_title')}}</text>
		</view>

		<!-- 碳排放信息 -->
		<view class="carbon-info">
			<view class="carbon-progress">
				<text class="carbon-description">{{$t('carbon_description')}}</text>
				<text class="carbon-number">{{ days }}{{$t('carbon_days')}}</text>
			</view>
			<view class="charts">

				<view class="chart today">
					<text class="chart-title">{{$t('carbon_today')}}</text>
					<view class="today-charts">
						<qiun-data-charts :canvas2d="true" canvas-id="carbonTodayChart" type="ring" :opts="ringOpts"
							:chartData="chartTodayData" />
					</view>
				</view>
				<view class="chart history">
					<text class="chart-title">{{$t('carbon_history')}}</text>
					<qiun-data-charts :canvas2d="true" canvas-id="carbonHistoryChart" type="line"
						:chartData="chartHistoryData" />
				</view>

				<!-- 新增今日营养情况图表 -->
				<view class="chart nutrition">
					<text class="chart-title">{{$t('nutrition_today')}}</text>
					<view class="nutrition-charts">
						<qiun-data-charts :canvas2d="true" canvas-id="nutritionChart" type="column"
							:opts="nutritionOpts" :chartData="chartNutritionData" />
					</view>
				</view>

			</view>
		</view>

		<!-- 实用工具 -->
		<view class="useful-tools">
			<text class="tools-title">{{$t('tools_title')}}</text>
			<view class="tools-grid">
				<view class="tool" @click="navigateTo('calculator')" animation="fadeInUp">
					<image src="https://cdn.pixabay.com/photo/2017/07/06/17/13/calculator-2478633_1280.png"
						:alt="$t('tool_carbon_calculator')" class="tool-icon" mode="aspectFill"></image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_carbon_calculator')}}</text>
						<text class="tool-info">{{$t('tool_carbon_calculator_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('recommend')" animation="fadeInUp" animation-delay="0.2s">
					<image src="https://cdn.pixabay.com/photo/2020/03/12/18/37/dish-4925892_1280.png"
						:alt="$t('tool_diet_recommendation')" class="tool-icon" mode="aspectFill"></image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_diet_recommendation')}}</text>
						<text class="tool-info">{{$t('tool_diet_recommendation_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('nutrition')" animation="fadeInUp" animation-delay="0.4s">
					<image src="https://cdn.pixabay.com/photo/2016/11/14/15/42/calendar-1823848_1280.png"
						:alt="$t('tool_nutrition_calculator')" class="tool-icon"></image>
					<view class="tool-description">
						<text class="tool-name">{{$t('tool_nutrition_calculator')}}</text>
						<text class="tool-info">{{$t('tool_nutrition_calculator_info')}}</text>
					</view>
				</view>
				<view class="tool" @click="navigateTo('family')" animation="fadeInUp" animation-delay="0.6s">
					<image src="https://cdn.pixabay.com/photo/2016/01/04/14/24/terminal-board-1120961_1280.png"
						:alt="$t('tool_family_recipe')" class="tool-icon"></image>
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
		onMounted,
		watch,
	} from 'vue'
	import {
		useI18n
	} from 'vue-i18n'
	import {
		useCarbonAndNutritionStore
	} from '@/stores/carbon_and_nutrition_data.js'
	import {
		onShow
	} from '@dcloudio/uni-app';

	// 初始化 i18n
	const {
		t,
		locale
	} = useI18n()

	onShow(() => {
    uni.setNavigationBarTitle({
      title: t('index')
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
  })

	// 控制查看更多内容的显示
	const showMore = ref(false)
	const days = ref(9) // 可以根据实际情况动态更新

	const toggleMoreContent = () => {
		showMore.value = !showMore.value
	}

	// 获取 Pinia 存储
	const carbonNutritionStore = useCarbonAndNutritionStore()

	// 历史碳排放数据，初始为空
	const chartHistoryData = ref({
		categories: [],
		series: [{
				name: t('target_value'),
				data: []
			},
			{
				name: t('actual_value'),
				data: []
			}
		]
	});

	// 在组件挂载时生成从七天前到今天的日期
	onMounted(() => {
		const categories = [];
		const targetData = [];
		const actualData = [];
		for (let i = 7; i > 0; i--) {
			const date = new Date();
			date.setDate(date.getDate() - i);
			const month = date.getMonth() + 1; // 月份从0开始
			const day = date.getDate();
			const dateString = date.toISOString().split('T')[0];
			categories.push(`${month}/${day}`);

			const dailyData = carbonNutritionStore.getDataByDate(dateString);

			if (dailyData) {
				targetData.push(dailyData.carbonEmission.target || 0);
				actualData.push(dailyData.carbonEmission.actual || 0);
			} else {
				targetData.push(0);
				actualData.push(0);
			}
		}
		chartHistoryData.value.categories = categories;
		chartHistoryData.value.series[0].data = targetData;
		chartHistoryData.value.series[1].data = actualData;
	});

	// 今日营养数据的条形图
	const chartNutritionData = ref({
		categories: [t('energy_unit'), t('protein_unit'), t('fat_unit'), t('carbohydrates_unit'), t(
			'sodium_unit')],
		series: [{
				name: t('intake'),
				data: []
			},
			{
				name: t('target_value'),
				data: []
			}
		]
	});

	// 营养条形图的配置选项
	const nutritionOpts = {
		color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE"],
		padding: [15, 15, 0, 5],
		enableScroll: false,
		legend: {},
		xAxis: {
			disableGrid: false,
			axisLine: true,
		},
		yAxis: {},
		extra: {
			bar: {
				type: "group",
				width: 30,
				meterBorde: 1,
				meterFillColor: "#FFFFFF",
				activeBgColor: "#000000",
				activeBgOpacity: 0.08,
				categoryGap: 2
			}
		}
	};

	// 环形图数据
	const chartTodayData = ref({
		series: [{
			data: []
		}]
	});

	// 环形图的配置选项
	const ringOpts = {
		rotate: false,
		rotateLock: false,
		color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
		animation: {
			duration: 0, // 设置为0可以避免动画延迟
		},
		padding: [5, 5, 5, 5],
		dataLabel: true,
		enableScroll: false,
		legend: {
			show: true,
			position: "right",
			lineHeight: 25
		},
		title: {
			name: t('total'),
			fontSize: 15,
			color: "#666666"
		},
		subtitle: {
			name: "",
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
	};

	// 更新今日图表
	onMounted(() => {
		const today = new Date();
		const dateString = today.toISOString().split('T')[0];
		const todayData = carbonNutritionStore.getDataByDate(dateString);

		if (todayData) {
			// 更新环形图数据
			const mealTypes = ['breakfast', 'lunch', 'dinner', 'others'];
			const mealData = [];
			let totalCarbonEmission = 0;
			mealTypes.forEach(mealType => {
				const meal = todayData.meals[mealType];
				mealData.push({
					name: t(mealType),
					value: meal.carbonEmission || 0
				});
				totalCarbonEmission += meal.carbonEmission || 0;
			});
			chartTodayData.value.series[0].data = mealData;

			// 更新环形图副标题
			ringOpts.subtitle.name = `${totalCarbonEmission.toFixed(1)}Kg`;

			// 更新营养条形图数据
			const nutrients = ['energy', 'protein', 'fat', 'carbohydrates', 'sodium'];
			const intakeData = [];
			const targetData = [];
			nutrients.forEach(nutrient => {
				intakeData.push(todayData.nutrients.actual[nutrient] || 0);
				targetData.push(todayData.nutrients.target[nutrient] || 0);
			});
			chartNutritionData.value.series[0].data = intakeData;
			chartNutritionData.value.series[1].data = targetData;
		} else {
			// 如果没有数据，设置为零
			chartTodayData.value.series[0].data = [];
			ringOpts.subtitle.name = `0Kg`;
			chartNutritionData.value.series[0].data = [0, 0, 0, 0, 0];
			chartNutritionData.value.series[1].data = [0, 0, 0, 0, 0];
		}
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
		will-change: contents;
	}

	.chart {
		transform: translateZ(0);
		will-change: transform;
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

	.today-charts,
	.nutrition-charts {
		align-items: center;
		width: 100%;
		height: 300px;
		position: relative;
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