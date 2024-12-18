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

				<!-- 今日碳排放环形图 -->
				<view class="chart today">
					<text class="chart-title">{{$t('carbon_today')}}</text>
					<view class="today-charts">
						<!-- 注意这里的 :opts="ringOpts" 若不同时更新，可以保留对 opts 的监听 -->
						<qiun-data-charts :canvas2d="true" type="ring" :opts="ringOpts"
							:chartData="chartTodayData" />
					</view>
				</view>

				<!-- 历史碳排放曲线图 -->
				<view class="chart history">
					<text class="chart-title">{{$t('carbon_history')}}</text>
					<qiun-data-charts :canvas2d="true" canvas-id="carbonHistoryChart" type="line" :opts="historyOpts"
						:chartData="chartHistoryData" />
				</view>

				<!-- 今日营养情况图表（柱状图） -->
				<view class="chart nutrition">
					<text class="chart-title">{{$t('nutrition_today')}}</text>
					<view class="nutrition-charts">
						<!-- 若需要同时变动 nutritionOpts 和 chartNutritionData，请参考官方文档 optsWatch=false 的做法 -->
						<qiun-data-charts :canvas2d="true" canvas-id="xvmMWWFdeOdEnvVDPjotdobEUaWVmvav" type="column" :ontouch="true"
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
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCarbonAndNutritionStore } from '@/stores/carbon_and_nutrition_data.js'
import { onShow } from '@dcloudio/uni-app';

const { t } = useI18n()
const showMore = ref(false)
const days = ref(9)

// 获取 Pinia 存储
const carbonNutritionStore = useCarbonAndNutritionStore()

// 图表数据（注意仅进行初始化）
const chartHistoryData = ref({
	categories: [],
	series: [
		{ name: t('target_value'), data: [] },
		{ name: t('actual_value'), data: [] }
	]
})

const chartNutritionData = ref({
	categories: [t('energy_unit'), t('protein_unit'), t('fat_unit'), t('carbohydrates_unit'), t('sodium_unit')],
	series: [
		{ name: t('intake'), data: [] },
		{ name: t('target_value'), data: [] }
	]
})

const chartTodayData = ref({
	series: [
		{ data: [] }
	]
})

// 样式设定
const nutritionOpts = {
	color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE"],
	padding: [15, 15, 0, 15],
	// enableScroll: true,
	xAxis: { disableGrid: false, axisLine: true, itemCount: 4, rotateLabel: true, rotateAngle: 60},
	yAxis: {disabled:true},
	extra: {
    column: {
      width: 20,
      type: 'group',
      seriesGap: 5,
    }
	}
}

const historyOpts = {
  yAxis: {
    disabled: true
  }
}

const ringOpts = {
	rotate: false,
	rotateLock: false,
	color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
	animation: { duration: 0 },
	padding: [5, 5, 5, 5],
	dataLabel: true,
	enableScroll: false,
	legend: { show: true, position: "bottom", lineHeight: 25 },
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
			ringWidth: 10,
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

// 根据日期从 store 数据中计算每日数据
const getDataByDate = (dateString) => {
  console.log('getDataByDate', dateString)
	// 从 store 中查找对应日期的营养目标/碳排放目标
	const nutritionGoal = carbonNutritionStore.state.nutritionGoals.find(g => g.date.startsWith(dateString))
	const carbonGoal = carbonNutritionStore.state.carbonGoals.find(g => g.date.startsWith(dateString))

	// 汇总该日期四餐的实际摄入
	const dailyNutritionIntakes = carbonNutritionStore.state.nutritionIntakes.filter(i => i.date.startsWith(dateString))
	const dailyCarbonIntakes = carbonNutritionStore.state.carbonIntakes.filter(i => i.date.startsWith(dateString))

	// 按餐分类实际摄入，预定义四种餐类型
	const meals = { breakfast: {}, lunch: {}, dinner: {}, other: {} }

	// 初始化营养实际总值
	const totalNutrients = { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
	let totalCarbonEmission = 0

	for (const intake of dailyNutritionIntakes) {
		const mealType = intake.meal_type || 'other'
		if (!meals[mealType].nutrients) {
			meals[mealType].nutrients = { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
		}
		meals[mealType].nutrients.calories += intake.calories || 0
		meals[mealType].nutrients.protein += intake.protein || 0
		meals[mealType].nutrients.fat += intake.fat || 0
		meals[mealType].nutrients.carbohydrates += intake.carbohydrates || 0
		meals[mealType].nutrients.sodium += intake.sodium || 0

		totalNutrients.calories += intake.calories || 0
		totalNutrients.protein += intake.protein || 0
		totalNutrients.fat += intake.fat || 0
		totalNutrients.carbohydrates += intake.carbohydrates || 0
		totalNutrients.sodium += intake.sodium || 0
	}

	for (const cIntake of dailyCarbonIntakes) {
		let mealType = cIntake.meal_type || 'other'
		if (!meals[mealType].carbonEmission) {
			meals[mealType].carbonEmission = 0
		}
		meals[mealType].carbonEmission += cIntake.emission || 0
		totalCarbonEmission += cIntake.emission || 0
	}

	return {
		nutrients: {
			actual: totalNutrients,
			target: nutritionGoal ? {
				calories: nutritionGoal.calories,
				protein: nutritionGoal.protein,
				fat: nutritionGoal.fat,
				carbohydrates: nutritionGoal.carbohydrates,
				sodium: nutritionGoal.sodium
			} : { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
		},
		carbonEmission: {
			actual: totalCarbonEmission,
			target: carbonGoal ? carbonGoal.emission : 0
		},
		meals
	}
}

onShow(async () => {
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
  console.log('onMounted')
  // 获取数据
  await carbonNutritionStore.getNutritionGoals()
  await carbonNutritionStore.getCarbonGoals()
  await carbonNutritionStore.getNutritionIntakes()
  await carbonNutritionStore.getCarbonIntakes()

  // 更新今日数据

  const today = new Date()
  console.log('today', today)
  const dateString = today.getFullYear() + '-'
      + String(today.getMonth() + 1).padStart(2, '0') + '-'
      + String(today.getDate()).padStart(2, '0');
  console.log('dateString', dateString)
  const todayData = getDataByDate(dateString)
  console.log('todayData', todayData)

  if (todayData) {
    // 1. 先构造今日环形图数据的临时变量
    const mealTypes = ['breakfast', 'lunch', 'dinner', 'other']
    let mealData = []
    let totalCarbonEmission = 0
    mealTypes.forEach(mealType => {
      const emission = todayData.meals[mealType].carbonEmission || 0
      mealData.push({
        name: t(mealType),
        value: emission
      })
      totalCarbonEmission += emission
    })

    // 今日环形图 chartData 的最终结构
    let tempChartTodayData = {
      series: [
        {data: mealData}
      ]
    }
    // 深拷贝后再赋值，避免局部多次变动
    chartTodayData.value = JSON.parse(JSON.stringify(tempChartTodayData))

    // 如果需要更新 opts 中的 subtitle，同样建议复制一份
    let tempRingOpts = JSON.parse(JSON.stringify(ringOpts))
    tempRingOpts.subtitle.name = `${totalCarbonEmission.toFixed(1)}Kg`
    // 这里是直接覆盖 ringOpts，若不想组件重复初始化可考虑 :optsWatch="false"
    Object.assign(ringOpts, tempRingOpts)

    // 2. 更新今日营养柱状图（同样用临时变量整体赋值）
    const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
    const intakeData = nutrients.map(n => todayData.nutrients.actual[n] || 0)
    const targetData = nutrients.map(n => todayData.nutrients.target[n] || 0)
    let tempNutritionData = {
      categories: chartNutritionData.value.categories,  // 保留原先的 categories
      series: [
        {name: t('intake'), data: intakeData},
        {name: t('target_value'), data: targetData}
      ]
    }
    chartNutritionData.value = JSON.parse(JSON.stringify(tempNutritionData))
  }

  // 更新历史碳排放曲线
  const categories = []
  const targetData = []
  const actualData = []
  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(d.getDate() - i)

    // 获取日期字符串格式化为 YYYY-MM-DD
    const year = d.getFullYear()
    const month = String(d.getMonth() + 1).padStart(2, '0') // 月份是从0开始的，所以加1
    const day = String(d.getDate()).padStart(2, '0') // 保证日期是两位数

    // 构造 ds 为 YYYY-MM-DD 格式的字符串
    const ds = `${year}-${month}-${day}`
    console.log('date string (ds)', ds)

    // 添加到categories
    categories.push(`${month}/${day}`)

    // 根据 ds 获取数据
    const dailyData = getDataByDate(ds)
    targetData.push(dailyData ? dailyData.carbonEmission.target : 0)
    actualData.push(dailyData ? dailyData.carbonEmission.actual : 0)
  }

  // 用临时变量构造历史曲线图数据
  let tempHistoryData = {
    categories: categories,
    series: [
      {name: t('target_value'), data: targetData},
      {name: t('actual_value'), data: actualData}
    ]
  }
  chartHistoryData.value = JSON.parse(JSON.stringify(tempHistoryData))
})

// 页面跳转方法
const navigateTo = (page) => {
  if (page === 'recommend') {
    uni.navigateTo({
      url: "/pagesTool/food_recommend/food_recommend",
    });
  } else if (page === 'nutrition') {
    uni.navigateTo({
      url: "/pagesTool/nutrition_calendar/nutrition_calendar",
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

.charts {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.chart {
  background-color: rgba(255, 255, 255, 0.9);
  padding: 20rpx;
  border-radius: 15rpx;
  box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
  margin-bottom: 40rpx;
  width: 90%;
  transition: transform 0.3s, box-shadow 0.3s;
}

.chart.today {
  background-color: rgba(24, 144, 255, 0.1);
}

.chart.history {
  background-color: rgba(255, 193, 7, 0.1);
}

.chart.nutrition {
  background-color: rgba(76, 175, 80, 0.1);
}

.chart:hover {
  transform: translateY(-5rpx);
  box-shadow: 0 8rpx 20rpx rgba(0, 0, 0, 0.2);
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
  transition: transform 0.3s, box-shadow 0.3s, background-color 0.3s;
  animation: fadeInUp 1s ease;
}

.tool:hover {
  transform: translateY(-5rpx) scale(1.05);
  box-shadow: 0 4rpx 10rpx var(--shadow-color);
  background-color: rgba(255, 255, 255, 1);
}

.tool:active {
  transform: translateY(0) scale(1);
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
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
