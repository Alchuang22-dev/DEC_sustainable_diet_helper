<template>
  <view class="container">
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

    <!-- 日期导航部分 -->
    <uni-section title="" padding>
      <view class="header">
        <text class="title">{{ $t('nutrition_calendar') }}</text>
        <scroll-view
            class="date-selector"
            scroll-x="true"
            scroll-with-animation="true"
            :scroll-into-view="'date-' + currentDateIndex"
            :scroll-left="scrollPosition"
        >
          <view class="date-buttons">
            <view
                :id="'date-' + index"
                v-for="(date, index) in dateTabs"
                :key="date.dateString"
                :class="['date-button', { 'selected': currentDateIndex === index }]"
                @click="onDateChange(index)"
            >
              <text class="day">{{ date.day }}</text>
              <text class="date">{{ date.date }}</text>
            </view>
          </view>
        </scroll-view>
      </view>
    </uni-section>

    <!-- 营养环形图部分 -->
    <uni-section title="" padding>
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
    </uni-section>

    <!-- 营养摄入详情部分 -->
    <uni-section :title="$t('nutrition_details')" padding>
      <view class="nutrition-details">
        <view
            class="nutrition-item"
            v-for="item in summaryNutrition"
            :key="item.label"
        >
          <view class="nutrition-card">
            <view
                class="color-square"
                :style="{ backgroundColor: item.over ? getRedShade(item.nutrient) : item.color }"
            ></view>
            <view class="nutrition-text">
              <text class="label">{{ item.label }}</text>
              <text class="value" :class="{ 'over': item.over }">
                {{ item.intake }} / {{ item.plan }}
              </text>
            </view>
          </view>
        </view>
      </view>
      <!-- 仅在查看今天的结果时显示“设置目标”按钮 -->
      <view class="set-goals-button-wrapper" v-if="isTodaySelected">
        <button class="set-goals-button" @click="navigateToSetGoals">
          {{ $t('set_nutrition_goals') }}
        </button>
      </view>
    </uni-section>

    <!-- 膳食记录部分 -->
    <uni-section :title="$t('meal_detail_records')" padding>
      <view class="meal-records">
        <uni-collapse ref="collapse" v-model="activeMeal" @change="onMealChange">
          <uni-collapse-item :title="$t('breakfast')" name="breakfast">
            <view class="meal-content">
              <view
                  v-for="nutrient in breakfastNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text class="nutrient-label">{{ nutrient.label }}:</text>
                <text class="nutrient-value">{{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="$t('lunch')" name="lunch">
            <view class="meal-content">
              <view
                  v-for="nutrient in lunchNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text class="nutrient-label">{{ nutrient.label }}:</text>
                <text class="nutrient-value">{{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="$t('dinner')" name="dinner">
            <view class="meal-content">
              <view
                  v-for="nutrient in dinnerNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text class="nutrient-label">{{ nutrient.label }}:</text>
                <text class="nutrient-value">{{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>

          <uni-collapse-item :title="$t('other')" name="other">
            <view class="meal-content">
              <view
                  v-for="nutrient in otherNutrients"
                  :key="nutrient.label"
                  class="nutrient-item"
              >
                <text class="nutrient-label">{{ nutrient.label }}:</text>
                <text class="nutrient-value">{{ nutrient.intake }}</text>
              </view>
            </view>
          </uni-collapse-item>
        </uni-collapse>
      </view>
    </uni-section>
  </view>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { onShow } from '@dcloudio/uni-app'
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

    // 获取星期几
    const day = getWeekdayKey(date.getDay())
    const dateNumber = date.getDate()

    // 格式化日期为 YYYY-MM-DD
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const dayOfMonth = String(date.getDate()).padStart(2, '0')
    const dateString = `${year}-${month}-${dayOfMonth}`

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
    return defaultEmptyData()
  }
  const selectedDate = dateTabs.value[currentDateIndex.value].dateString
  const dateData = getDataByDate(selectedDate)

  if (dateData) {
    return dateData
  } else {
    // 如果获取不到数据，就返回空数据结构
    return defaultEmptyData()
  }
}

// 当找不到数据时，返回这个空结构，避免 undefined 出错
const defaultEmptyData = () => {
  return {
    nutrients: {
      actual: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0},
      target: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0}
    },
    meals: {
      breakfast: {nutrients: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0}},
      lunch: {nutrients: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0}},
      dinner: {nutrients: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0}},
      other: {nutrients: {calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0}}
    }
  }
}

// 更新五大营养的过量标志
const updateSummaryNutrition = () => {
  // 获取当前日期的数据
  const dateData = getDataForSelectedDate()

  // 处理实际与目标的安全取值
  const actualNutrients = dateData?.nutrients?.actual || {}
  const targetNutrients = dateData?.nutrients?.target || {}

  const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
  const tempSummary = nutrients.map(nutrient => {
    const intake = actualNutrients[nutrient] || 0
    const plan = targetNutrients[nutrient] || 0
    let ratio = 0

    if (plan > 0) {
      ratio = intake / plan
      if (ratio > 1) {
        ratio = 1 // 确保比例不超过 1
      }
    } else {
      ratio = intake > 0 ? 1 : 0
    }

    const over = plan > 0 && intake > plan
    const color = over ? getRedShade(nutrient) : getNutrientColor(nutrient)

    return {
      nutrient: nutrient, // 用于颜色映射
      label: t(`${nutrient}_unit`),
      intake: intake,
      plan: plan,
      ratio: ratio,
      color: color,
      over: over
    }
  })

  summaryNutrition.value = JSON.parse(JSON.stringify(tempSummary))

  // 每餐的营养数据
  const m = dateData.meals || {}
  breakfastNutrients.value = mapMealNutrients(m.breakfast?.nutrients || {})
  lunchNutrients.value = mapMealNutrients(m.lunch?.nutrients || {})
  dinnerNutrients.value = mapMealNutrients(m.dinner?.nutrients || {})
  otherNutrients.value = mapMealNutrients(m.other?.nutrients || {})
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
    data: item.ratio,
    color: item.color
  }))

  const tempChartData = {series: tempSeries}
  chartData.value = JSON.parse(JSON.stringify(tempChartData))
}

// Helper to get 红色不同深浅基于营养成分
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

const scrollPosition = ref(0)

// 计算今天的日期字符串
const todayDateString = computed(() => {
  const today = new Date()
  const year = today.getFullYear()
  const month = String(today.getMonth() + 1).padStart(2, '0')
  const day = String(today.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
})

// 计算当前选中的日期是否为今天
const isTodaySelected = computed(() => {
  return dateTabs.value[currentDateIndex.value]?.dateString === todayDateString.value
})

// 在 onMounted 中初始化
onShow(async () => {
  generateDateTabs()
  await carbonNutritionStore.getNutritionGoals()
  await carbonNutritionStore.getNutritionIntakes()

  // 设置初始日期为今天（倒数第二个：i = 0, index = 6）
  currentDateIndex.value = dateTabs.value.length - 2

  // 让滚动条自动滚到最右
  nextTick(() => {
    scrollPosition.value = 9999
  })

  updateSummaryNutrition()
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

// 监听 summaryNutrition 的变化以更新图表
watch(summaryNutrition, () => {
  updateChartData()
})
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-bottom: 40rpx;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.08;
}

.header {
  padding: 20rpx;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
}

.title {
  font-size: 36rpx;
  font-weight: 600;
  color: #4CAF50FF;
  text-align: center;
  margin-bottom: 24rpx;
}

.date-selector {
  width: 100%;
  padding: 10rpx 0;
}

.date-buttons {
  display: flex;
  padding: 0 10rpx;
}

.date-button {
  flex: none;
  width: 140rpx;
  height: 120rpx;
  margin: 0 8rpx;
  padding: 16rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.date-button.selected {
  background-color: #4CAF50;
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 12rpx rgba(76, 175, 80, 0.2);
}

.date-button.selected .day,
.date-button.selected .date {
  color: #ffffff;
}

.day {
  font-size: 24rpx;
  color: #666;
  margin-bottom: 8rpx;
}

.date {
  font-size: 32rpx;
  font-weight: bold;
  color: #2c3e50;
}

.charts-box {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin: 20rpx 0;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
  min-height: 500rpx;
}

.chart-center-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  z-index: 2;
}

.center-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #2c3e50;
  white-space: pre-line;
  text-align: center;
}

.nutrition-details {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin: 0;
  position: relative;
  z-index: 1;
}

.nutrition-item {
  margin-bottom: 16rpx;
}

.nutrition-item:last-child {
  margin-bottom: 0;
}

.nutrition-card {
  display: flex;
  align-items: center;
  padding: 16rpx;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  transition: all 0.3s ease;
}

.nutrition-card:hover {
  transform: translateX(4rpx);
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.color-square {
  width: 32rpx;
  height: 32rpx;
  border-radius: 8rpx;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.nutrition-text {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label {
  font-size: 28rpx;
  color: #2c3e50;
  font-weight: 500;
}

.value {
  font-size: 28rpx;
  color: #4CAF50;
  font-weight: 500;
}

.value.over {
  color: #ff4d4f;
}

.set-goals-button-wrapper {
  padding: 20rpx;
  display: flex;
  justify-content: center;
}

.set-goals-button {
  width: 50%;
  height: 88rpx;
  background-color: #4CAF50;
  color: #ffffff;
  border-radius: 44rpx;
  font-size: 32rpx;
  font-weight: 500;
  box-shadow: 0 4rpx 12rpx rgba(76, 175, 80, 0.2);
  transition: all 0.3s ease;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.set-goals-button:active {
  transform: translateY(2rpx);
  box-shadow: 0 2rpx 6rpx rgba(76, 175, 80, 0.2);
  background-color: #43a047;
}

.meal-records {
  background-color: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

.meal-content {
  padding: 20rpx;
  background-color: #f8f9fa;
  border-radius: 8rpx;
}

.nutrient-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
  border-bottom: 1rpx solid rgba(0, 0, 0, 0.05);
}

.nutrient-item:last-child {
  border-bottom: none;
}

.nutrient-label {
  font-size: 28rpx;
  color: #2c3e50;
}

.nutrient-value {
  font-size: 28rpx;
  color: #4CAF50;
  font-weight: 500;
}

/* uni-section 样式优化 */
:deep(.uni-section) {
  padding: 0 20rpx;
  margin-bottom: 20rpx;
}

:deep(.uni-section-header) {
  padding: 20rpx 0;
}

:deep(.uni-section-header__content) {
  color: #2c3e50;
  font-size: 32rpx;
  font-weight: 600;
}

/* uni-collapse 样式优化 */
:deep(.uni-collapse-item) {
  margin-bottom: 12rpx;
}

:deep(.uni-collapse-item__title) {
  padding: 20rpx;
  background-color: #ffffff;
  border-radius: 12rpx;
}

:deep(.uni-collapse-item__title-text) {
  font-size: 30rpx;
  color: #2c3e50;
  font-weight: 500;
}

:deep(.uni-collapse-item__wrap) {
  background-color: transparent;
}

/* 滚动条样式优化 */
::-webkit-scrollbar {
  width: 0;
  height: 0;
  background: transparent;
}

/* 响应式适配 */
@media screen and (max-width: 375px) {
  .date-button {
    width: 120rpx;
    height: 100rpx;
  }

  .day {
    font-size: 22rpx;
  }

  .date {
    font-size: 28rpx;
  }

  .label,
  .value,
  .nutrient-label,
  .nutrient-value {
    font-size: 26rpx;
  }
}
</style>
