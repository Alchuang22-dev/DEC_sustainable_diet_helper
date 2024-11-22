<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 合并后的顶部导航和日期选择器 -->
    <view class="header">
      <text class="title">{{ $t('nutrition_calendar') }}</text>
      <scroll-view class="date-selector" scroll-x="true" scroll-with-animation="true">
        <view class="date-buttons">
          <button v-for="(date, index) in dateTabs" :key="date.dateString"
                  :class="['date-button', { 'selected': currentDateIndex === index }]"
                  @click="onDateChange(index)">
            <text class="day">{{ date.day }}</text>
            <text class="date">{{ date.date }}</text>
          </button>
        </view>
      </scroll-view>
    </view>


    <!-- 多重环形图 -->
    <view class="charts-box">
      <qiun-data-charts type="arcbar" :opts="chartOpts" :chartData="chartData" :canvas2d="true"
                        canvasId="nutritionChart" />
      <view class="chart-center-text">
        <text class="center-title">{{ $t('nutrition\noverview') }}</text>
      </view>
    </view>

    <!-- 五大营养详细信息 -->
    <view class="nutrition-details">
      <view class="nutrition-item" v-for="item in summaryNutrition" :key="item.label">
        <view class="color-square"
              :style="{ backgroundColor: item.over ? getRedShade(item.label) : item.color }"></view>
        <view class="nutrition-text">
          <text class="intake" :class="{ 'over': item.over }">
            {{ item.label }}:
            <text :class="{ 'over': item.over }">{{ item.intake }}</text> / {{ item.plan }}
          </text>
        </view>
      </view>
    </view>

    <!-- 设置营养目标按钮 -->
    <view class="set-goals-button-wrapper">
      <button class="set-goals-button" @click="navigateToSetGoals">
        {{ $t('set_nutrition_goals') }}
      </button>
    </view>

    <!-- 饮食记录 -->
    <view class="meal-records">
      <uni-section :title="t('meal_detail_records')">
        <uni-collapse ref="collapse" v-model="activeMeal" @change="onMealChange">
          <uni-collapse-item :title="t('breakfast')" :name="'breakfast'">
            <view class="content">
              <view v-for="nutrient in breakfastNutrients" :key="nutrient.label" class="nutrient-item">
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('lunch')" :name="'lunch'">
            <view class="content">
              <view v-for="nutrient in lunchNutrients" :key="nutrient.label" class="nutrient-item">
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('dinner')" :name="'dinner'">
            <view class="content">
              <view v-for="nutrient in dinnerNutrients" :key="nutrient.label" class="nutrient-item">
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('others')" :name="'others'">
            <view class="content">
              <view v-for="nutrient in otherNutrients" :key="nutrient.label" class="nutrient-item">
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>
        </uni-collapse>
      </uni-section>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCarbonAndNutritionStore } from '@/stores/carbon_and_nutrition_data.js'

// 国际化
const { t } = useI18n()

// 当前日期索引，默认选中今天（索引为3）
const currentDateIndex = ref(3)

// 生成日期标签：三天前到三天后
const dateTabs = ref([])

const generateDateTabs = () => {
  const tabs = []
  const today = new Date()
  for (let i = -3; i <= 3; i++) {
    const date = new Date()
    date.setDate(today.getDate() + i)
    const day = t(getWeekdayKey(date.getDay())) // 国际化后的星期
    const dateNumber = date.getDate()
    const dateString = date.toISOString().split('T')[0] // e.g., 2024-04-27
    tabs.push({
      day: day,
      date: dateNumber,
      dateString: dateString
    })
  }
  dateTabs.value = tabs
}

// Helper to get weekday key for i18n
const getWeekdayKey = (dayIndex) => {
  const weekdays = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday']
  return weekdays[dayIndex]
}

// 初始化日期标签
onMounted(() => {
  generateDateTabs()
  updateSummaryNutrition()
  updateChartData()
})

// 获取 Pinia 存储
const carbonNutritionStore = useCarbonAndNutritionStore()

// 五大营养摄入与计划
const summaryNutrition = ref([])

// 当前展开的折叠项
const activeMeal = ref('breakfast')

// 每餐的营养成分
const breakfastNutrients = ref([])
const lunchNutrients = ref([])
const dinnerNutrients = ref([])
const otherNutrients = ref([])

// 更新五大营养的过量标志
const updateSummaryNutrition = () => {
  const selectedDate = dateTabs.value[currentDateIndex.value].dateString
  const dateData = carbonNutritionStore.getDataByDate(selectedDate)

  if (dateData) {
    const nutrients = ['energy', 'protein', 'fat', 'carbohydrates', 'sodium']
    summaryNutrition.value = nutrients.map(nutrient => {
      const intake = dateData.nutrients.actual[nutrient] || 0
      const plan = dateData.nutrients.target[nutrient] || 0
      const over = intake > plan
      const color = getNutrientColor(nutrient)
      return {
        label: t(`${nutrient}_unit`),
        intake: intake,
        plan: plan,
        color: color,
        over: over
      }
    })

    // 更新每餐的营养成分
    const meals = ['breakfast', 'lunch', 'dinner', 'others']
    const mealNutrientsRefs = [breakfastNutrients, lunchNutrients, dinnerNutrients, otherNutrients]
    meals.forEach((mealType, index) => {
      const mealData = dateData.meals[mealType]
      if (mealData) {
        mealNutrientsRefs[index].value = nutrients.map(nutrient => {
          return {
            label: t(`${nutrient}_unit`),
            intake: mealData.nutrients[nutrient] || 0
          }
        })
      } else {
        mealNutrientsRefs[index].value = []
      }
    })
  } else {
    // 选定日期没有数据
    summaryNutrition.value = []
    breakfastNutrients.value = []
    lunchNutrients.value = []
    dinnerNutrients.value = []
    otherNutrients.value = []
  }
}

// 更新图表数据
const updateChartData = () => {
  // 根据 summaryNutrition 更新 chartData
  chartData.value.series = getChartSeries()
}

// 图表数据和选项
const chartData = ref({
  series: []
})

const chartOpts = ref({
  color: [], // 不再使用全局颜色
  padding: undefined,
  title: {
    name: "",
    fontSize: 35,
    color: "#1890ff"
  },
  subtitle: {
    name: "",
    fontSize: 15,
    color: "#666666"
  },
  extra: {
    arcbar: {
      type: "circle",
      width: 10, // 增加宽度确保环足够宽
      backgroundColor: "#E9E9E9",
      startAngle: 1.5,
      endAngle: 0.25,
      gap: 2
    }
  }
})

// 获取图表数据，处理超标情况
const getChartSeries = () => {
  return summaryNutrition.value.map(item => {
    const data = item.intake > item.plan ? 1 : item.intake / item.plan
    const color = item.over ? getRedShade(item.label) : item.color
    return {
      name: item.label,
      data: data,
      color: color,
      over: item.over
    }
  })
}

// Helper to get红色不同深浅基于标签
const getRedShade = (label) => {
  const shades = {
    [t('energy_unit')]: '#FF4D4F',
    [t('protein_unit')]: '#FF7875',
    [t('fat_unit')]: '#FFB3BA',
    [t('carbohydrates_unit')]: '#FFA39E',
    [t('sodium_unit')]: '#FF4D4F'
  }
  return shades[label] || '#FF4D4F'
}

// Helper to get nutrient color
const getNutrientColor = (nutrient) => {
  const colors = {
    'energy': '#1890FF',
    'protein': '#91CB74',
    'fat': '#FAC858',
    'carbohydrates': '#73C0DE',
    'sodium': '#3CA272'
  }
  return colors[nutrient] || '#000000'
}

// 监听 summaryNutrition 的变化以更新图表
watch(summaryNutrition, () => {
  updateChartData()
})

// 日期改变
const onDateChange = (index) => {
  currentDateIndex.value = index
  updateSummaryNutrition()
  updateChartData()
}

// 折叠项切换处理
const onMealChange = (name) => {
  console.log('当前展开的餐：', name)
}

// 跳转到设置营养目标页面
const navigateToSetGoals = () => {
  uni.navigateTo({
    url: "/pagesMy/setGoals/setGoals",
  })
}
</script>

<style scoped>
	/* 使用主页中的通用变量 */
	:root {
		--primary-color: #4CAF50;
		--secondary-color: #2fc25b;
		--background-color: #f5f5f5;
		--section-background: rgba(144, 238, 144, 0.3);
		/* 淡绿色透明背景 */
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
		/* 将背景图片置于最底层 */
		opacity: 0.1;
		/* 调整透明度以不干扰内容 */
	}

	/* 合并后的头部 */
	.header {
		padding: 20rpx;
		background-color: var(--section-background);
		display: flex;
		flex-direction: column;
		align-items: center;
		position: sticky;
		top: 0;
		z-index: 10;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		backdrop-filter: blur(10rpx);
		border-radius: 10rpx;
	}

	/* 标题 */
	.title {
		font-size: var(--font-size-title);
		font-weight: bold;
		color: var(--primary-color);
		margin-bottom: 20rpx;
	}

	/* 日期选择器 */
	.date-selector {
		width: 100%;
		overflow-x: auto;
		white-space: nowrap;
	}

	.date-buttons {
		display: flex;
		flex-direction: row;
	}

	.date-button {
		flex: none;
		width: 150rpx;
		/* 固定宽度 */
		height: 120rpx;
		/* 固定高度 */
		padding: 10rpx;
		margin: 0 5rpx;
		border: none;
		background-color: #f0f0f0;
		border-radius: 20rpx;
		/* 圆角 */
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		transition: background-color 0.3s, color 0.3s;
	}

	.date-button.selected {
		background-color: var(--primary-color);
		color: #ffffff;
	}

	.day {
		font-size: 24rpx;
	}

	.date {
		font-size: 28rpx;
		font-weight: bold;
	}

	/* 多重环形图 */
	.charts-box {
		position: relative;
		width: 500rpx;
		height: 500rpx;
		margin: 40rpx auto;
		background-color: var(--section-background);
		border-radius: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		backdrop-filter: blur(10rpx);
		padding: 20rpx;
	}

	/* 环形图中心文字 */
	.chart-center-text {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		text-align: center;
	}

	.center-title {
		font-size: 36rpx;
		font-weight: bold;
		color: var(--text-color);
	}

	/* 五大营养详细信息 */
	.nutrition-details {
		padding: 20rpx;
		background-color: var(--section-background);
		border-radius: 10rpx;
		margin: 20rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		display: flex;
		flex-direction: column;
		backdrop-filter: blur(10rpx);
	}

	.nutrition-item {
		display: flex;
		align-items: center;
		margin-bottom: 10rpx;
	}

	.color-square {
		width: 30rpx;
		height: 30rpx;
		border-radius: 5rpx;
		margin-right: 10rpx;
	}

	.nutrition-text {
		font-size: 24rpx;
		color: #555;
		line-height: 36rpx;
		/* 增加行间距 */
		margin-left: 20rpx;
		/* 增加左边距 */
	}

	.nutrition-text.over {
		color: red;
	}

	.nutrition-text .intake {
		font-size: 24rpx;
		color: #555;
	}

	.nutrition-text .intake.over {
		color: red;
	}

	/* 设置营养目标按钮 */
	.set-goals-button-wrapper {
		display: flex;
		justify-content: center;
		margin: 20rpx;
	}

	.set-goals-button {
		background-color: var(--primary-color);
		color: #fff;
		border: none;
		border-radius: 30rpx;
		padding: 15rpx 30rpx;
		font-size: 28rpx;
		font-weight: bold;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		transition: background-color 0.3s, transform 0.3s;
	}

	.set-goals-button:hover {
		background-color: var(--secondary-color);
		transform: translateY(-2rpx);
	}

	/* 饮食记录 */
	.meal-records {
		padding: 20rpx;
		background-color: var(--section-background);
		border-radius: 10rpx;
		margin: 20rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		backdrop-filter: blur(10rpx);
	}

	.uni-section {
		background-color: transparent;
		/* 让背景继承父元素 */
	}

	.uni-collapse .content {
		padding: 20rpx;
		background-color: rgba(255, 255, 255, 0.5);
		/* 半透明背景 */
		border-radius: 10rpx;
		margin-top: 10rpx;
		box-shadow: 0 2rpx 5rpx var(--shadow-color);
		line-height: 36rpx;
		/* 增加行间距 */
		margin-left: 20rpx;
		/* 增加左边距 */
	}

	.nutrient-item {
		margin-bottom: 8rpx;
		font-size: 24rpx;
		color: #555;
	}
</style>