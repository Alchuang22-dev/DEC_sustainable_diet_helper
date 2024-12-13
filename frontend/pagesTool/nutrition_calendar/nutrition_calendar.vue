<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 合并后的顶部导航和日期选择器 -->
    <view class="header">
      <text class="title">{{ $t('nutrition_calendar') }}</text>
      <scroll-view class="date-selector" scroll-x="true" scroll-with-animation="true">
        <view class="date-buttons">
          <button
              v-for="(date, index) in dateTabs"
              :key="date.dateString"
              :class="['date-button', { 'selected': currentDateIndex === index }]"
              @click="onDateChange(index)"
          >
            <text class="day">{{ date.day }}</text>
            <text class="date">{{ date.date }}</text>
          </button>
        </view>
      </scroll-view>
    </view>

    <!-- 多重环形图 -->
    <view class="charts-box">
      <qiun-data-charts
          type="arcbar"
          :opts="chartOpts"
          :chartData="chartData"
          :canvas2d="true"
          canvasId="nutritionChart"
      />
      <view class="chart-center-text">
        <text class="center-title">{{ $t('nutrition\noverview') }}</text>
      </view>
    </view>

    <!-- 五大营养详细信息 -->
    <view class="nutrition-details">
      <view
          class="nutrition-item"
          v-for="item in summaryNutrition"
          :key="item.label"
      >
        <view
            class="color-square"
            :style="{ backgroundColor: item.over ? getRedShade(item.nutrient) : item.color }"
        ></view>
        <view class="nutrition-text">
          <text class="intake" :class="{ 'over': item.over }">
            {{ item.label }}:
            <text :class="{ 'over': item.over }">{{ item.intake }}</text>
            / {{ item.plan }}
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
          <uni-collapse-item :title="t('breakfast')" name="breakfast">
            <view class="content">
              <view
                  v-for="nutrient in breakfastNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('lunch')" name="lunch">
            <view class="content">
              <view
                  v-for="nutrient in lunchNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('dinner')" name="dinner">
            <view class="content">
              <view
                  v-for="nutrient in dinnerNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text>{{ nutrient.label }}: {{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="t('other')" name="other">
            <view class="content">
              <view
                  v-for="nutrient in otherNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
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

// 获取 Pinia 存储
const carbonNutritionStore = useCarbonAndNutritionStore()

// 当前日期索引
const currentDateIndex = ref(6)
const dateTabs = ref([])

// 营养概览数据
const summaryNutrition = ref([])

// 每餐的详细营养数据
const breakfastNutrients = ref([])
const lunchNutrients = ref([])
const dinnerNutrients = ref([])
const otherNutrients = ref([])
const activeMeal = ref('breakfast')

// 图表数据和配置（五大营养arcbar图）
const chartData = ref({ series: [] })
const chartOpts = ref({
  title: { name: "", fontSize: 35, color: "#1890ff" },
  subtitle: { name: "", fontSize: 15, color: "#666666" },
  extra: {
    arcbar: {
      type: "circle",
      width: 10,
      backgroundColor: "#E9E9E9",
      startAngle: 1.5,
      endAngle: 0.25,
      gap: 2
    }
  }
})

// 生成日期标签：今天及以前6天加上明天一天，共8天
const generateDateTabs = () => {
  const tabs = []
  const today = new Date()
  for (let i = -6; i <= 1; i++) {
    const date = new Date()
    date.setDate(today.getDate() + i)
    const day = getWeekdayKey(date.getDay())
    const dateNumber = date.getDate()
    const dateString = date.toLocaleDateString('en-CA', { timeZone: 'Asia/Shanghai' }).replace(/\//g, '-').split('T')[0]
    tabs.push({
      day: t(day),
      date: dateNumber,
      dateString: dateString
    })
  }
  dateTabs.value = tabs
}

// 返回星期对应的 key
const getWeekdayKey = (dayIndex) => {
  const weekdays = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday']
  return weekdays[dayIndex]
}

// 根据日期从 store 获取数据
const getDataByDate = (dateString) => {
  return carbonNutritionStore.getDataByDate(dateString)
}

// 获取当前选中日期的数据
const getDataForSelectedDate = () => {
  if (currentDateIndex.value < 0 || currentDateIndex.value >= dateTabs.value.length) {
    return {
      nutrients: {
        actual: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
        target: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
      },
      meals: {
        breakfast: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        lunch: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        dinner: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        other: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } }
      }
    }
  }

  const selectedDate = dateTabs.value[currentDateIndex.value].dateString
  const dateData = getDataByDate(selectedDate)

  if (dateData) {
    return dateData
  } else {
    return {
      nutrients: {
        actual: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
        target: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
      },
      meals: {
        breakfast: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        lunch: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        dinner: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } },
        other: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 } }
      }
    }
  }
}

// 更新五大营养的过量标志
const updateSummaryNutrition = () => {
  const dateData = getDataForSelectedDate()

  const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
  const tempSummary = nutrients.map(nutrient => {
    const intake = dateData.nutrients.actual[nutrient] || 0
    const plan = dateData.nutrients.target[nutrient] || 0
    let ratio = 0
    if (plan > 0) {
      ratio = intake / plan
      if (ratio > 1) {
        ratio = 1 // 确保比例不超过1
      }
    } else {
      ratio = intake > 0 ? 1 : 0 // 如果计划为0且摄入大于0，设置为1，否则为0
    }
    const over = plan > 0 && intake > plan
    const color = over ? getRedShade(nutrient) : getNutrientColor(nutrient)
    return {
      nutrient: nutrient, // 新增字段，用于颜色映射
      label: t(`${nutrient}_unit`),
      intake: intake,
      plan: plan,
      ratio: ratio,
      color: color,
      over: over
    }
  })

  // 使用临时变量整体赋值
  summaryNutrition.value = JSON.parse(JSON.stringify(tempSummary))

  // 更新每餐营养
  const m = dateData.meals
  breakfastNutrients.value = mapMealNutrients(m.breakfast.nutrients)
  lunchNutrients.value = mapMealNutrients(m.lunch.nutrients)
  dinnerNutrients.value = mapMealNutrients(m.dinner.nutrients)
  otherNutrients.value = mapMealNutrients(m.other.nutrients)
}

// 映射每餐的营养数据
const mapMealNutrients = (mealN) => {
  const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
  return nutrients.map(nutrient => ({
    label: t(`${nutrient}_unit`),
    intake: mealN[nutrient] || 0
  }))
}

// 更新图表数据（五大营养arcbar图）
const updateChartData = () => {
  const tempSeries = summaryNutrition.value.map(item => ({
    data: item.ratio, // 已在 summaryNutrition 中处理过比例
    color: item.color
  }))

  // 使用临时变量整体赋值
  const tempChartData = {
    series: tempSeries
  }

  chartData.value = JSON.parse(JSON.stringify(tempChartData))
}

// Helper to get红色不同深浅基于营养成分
const getRedShade = (nutrient) => {
  const shades = {
    'calories': '#FF4D4F',
    'protein': '#FF7875',
    'fat': '#FFB3BA',
    'carbohydrates': '#FFA39E',
    'sodium': '#FF4D4F'
  }
  return shades[nutrient] || '#FF4D4F'
}

// Helper to get nutrient color
const getNutrientColor = (nutrient) => {
  const colors = {
    'calories': '#1890FF',
    'protein': '#91CB74',
    'fat': '#FAC858',
    'carbohydrates': '#73C0DE',
    'sodium': '#3CA272'
  }
  return colors[nutrient] || '#000000'
}

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

// 初始化数据
onMounted(async () => {
  generateDateTabs()
  await carbonNutritionStore.getNutritionGoals()
  await carbonNutritionStore.getNutritionIntakes()

  updateSummaryNutrition()
  updateChartData()
})

// 监听 summaryNutrition 的变化以更新图表
watch(summaryNutrition, () => {
  updateChartData()
})
</script>

<style scoped>
:root {
  --primary-color: #4CAF50;
  --secondary-color: #2fc25b;
  --background-color: #f5f5f5;
  --section-background: rgba(144, 238, 144, 0.3);
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --font-size-title: 32rpx;
  --font-size-subtitle: 24rpx;
}

.container {
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
  min-height: 100vh;
  padding-bottom: 80rpx;
  position: relative;
  overflow: hidden;
}

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

.title {
  font-size: var(--font-size-title);
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 20rpx;
}

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
  height: 120rpx;
  padding: 10rpx;
  margin: 0 5rpx;
  border: none;
  background-color: #f0f0f0;
  border-radius: 20rpx;
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
  margin-left: 20rpx;
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
}

.uni-collapse .content {
  padding: 20rpx;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 10rpx;
  margin-top: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  line-height: 36rpx;
  margin-left: 20rpx;
}

.nutrient-item {
  margin-bottom: 8rpx;
  font-size: 24rpx;
  color: #555;
}
</style>
