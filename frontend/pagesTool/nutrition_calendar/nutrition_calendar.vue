<template>
  <view class="container">
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    />

    <!-- 日期导航 -->
    <uni-section title="" padding>
      <view class="header">
        <text class="title">{{ t('nutrition_calendar') }}</text>
        <scroll-view
          class="date-selector"
          scroll-x
          scroll-with-animation
          :scroll-into-view="'date-' + currentDateIndex"
          :scroll-left="scrollPosition"
        >
          <view class="date-buttons">
            <view
                v-for="(date, index) in dateTabs"
                :key="date.dateString"
                :id="'date-' + index"
                :class="['date-button', { selected: currentDateIndex === index }]"
                @click="onDateChange(index)"
            >
              <text class="day">{{ date.day }}</text>
              <text class="date">{{ date.date }}</text>
            </view>
          </view>
        </scroll-view>
      </view>
    </uni-section>

    <!-- 营养环形图 -->
    <uni-section title="" padding>
      <view class="charts-box">
        <qiun-data-charts
            type="arcbar"
            :opts="chartOpts"
            :chartData="chartData"
            :canvas2d="true"
            canvas-id="nutritionChart"
        />
        <view class="chart-center-text">
          <text class="center-title">
            {{ t('nutrition\noverview') }}
          </text>
        </view>
      </view>
    </uni-section>

    <!-- 营养摄入详情 -->
    <uni-section :title="t('nutrition_details')" padding>
      <view class="nutrition-details">
        <view
            v-for="item in summaryNutrition"
            :key="item.label"
            class="nutrition-item"
        >
          <view class="nutrition-card">
            <view
                class="color-square"
                :style="{
                backgroundColor: item.over ? getRedShade(item.nutrient) : item.color
              }"
            />
            <view class="nutrition-text">
              <text class="label">{{ item.label }}</text>
              <text class="value" :class="{ over: item.over }">
                {{ item.intake }} / {{ item.plan }}
              </text>
            </view>
          </view>
        </view>
      </view>

      <!-- 仅今天可设置目标 -->
      <view class="set-goals-button-wrapper" v-if="isTodaySelected">
        <button class="set-goals-button" @click="navigateToSetGoals">
          {{ t('set_nutrition_goals') }}
        </button>
      </view>
    </uni-section>

    <!-- 膳食记录部分 -->
    <uni-section :title="t('meal_detail_records')" padding>
      <view class="meal-records">
        <uni-collapse v-model="activeMeal" @change="onMealChange">
          <uni-collapse-item :title="t('breakfast')" name="breakfast">
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
          <uni-collapse-item :title="t('lunch')" name="lunch">
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
          <uni-collapse-item :title="t('dinner')" name="dinner">
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
          <uni-collapse-item :title="t('other')" name="other">
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
/* ----------------- Imports ----------------- */
import {ref, watch, computed, nextTick} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import {useI18n} from 'vue-i18n'
import {useCarbonAndNutritionStore} from '@/stores/carbon_and_nutrition_data.js'

/* ----------------- Setup ----------------- */
const {t} = useI18n()
const carbonNutritionStore = useCarbonAndNutritionStore()

/* ----------------- Reactive & State ----------------- */
const currentDateIndex = ref(6)
const dateTabs = ref([])
const summaryNutrition = ref([])

const breakfastNutrients = ref([])
const lunchNutrients = ref([])
const dinnerNutrients = ref([])
const otherNutrients = ref([])
const activeMeal = ref('breakfast')

const chartData = ref({series: []})
const chartOpts = ref({
  title: {name: '', fontSize: 35, color: '#1890ff'},
  subtitle: {name: '', fontSize: 15, color: '#666666'},
  extra: {
    arcbar: {
      type: 'circle',
      width: 10,
      backgroundColor: '#E9E9E9',
      startAngle: 1.5,
      endAngle: 0.25,
      gap: 2
    }
  }
})

const scrollPosition = ref(0)

/* ----------------- Computed ----------------- */
const todayDateString = computed(() => {
  const today = new Date()
  const year = today.getFullYear()
  const month = String(today.getMonth() + 1).padStart(2, '0')
  const day = String(today.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
})
const isTodaySelected = computed(() => {
  return dateTabs.value[currentDateIndex.value]?.dateString === todayDateString.value
})

/* ----------------- Lifecycle ----------------- */
onShow(async () => {
  generateDateTabs()
  await carbonNutritionStore.getNutritionGoals()
  await carbonNutritionStore.getNutritionIntakes()
  currentDateIndex.value = dateTabs.value.length - 2
  nextTick(() => {
    scrollPosition.value = 9999
  })
  updateSummaryNutrition()
  updateChartData()
})

/* ----------------- Methods ----------------- */
// 生成 8 天：今天前 6 天 + 今天 + 明天
function generateDateTabs() {
  const tabs = []
  const today = new Date()
  for (let i = -6; i <= 1; i++) {
    const date = new Date()
    date.setDate(today.getDate() + i)
    const day = getWeekdayKey(date.getDay())
    const dateNumber = date.getDate()
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const dayOfMonth = String(date.getDate()).padStart(2, '0')
    const dateString = `${year}-${month}-${dayOfMonth}`
    tabs.push({
      day: t(day),
      date: dateNumber,
      dateString
    })
  }
  dateTabs.value = tabs
}

function getWeekdayKey(dayIndex) {
  const weekdays = [
    'sunday',
    'monday',
    'tuesday',
    'wednesday',
    'thursday',
    'friday',
    'saturday'
  ]
  return weekdays[dayIndex]
}

// 根据当前选中日期，获取 store 数据
function getDataForSelectedDate() {
  if (
      currentDateIndex.value < 0 ||
      currentDateIndex.value >= dateTabs.value.length
  ) {
    return defaultEmptyData()
  }
  const selectedDate = dateTabs.value[currentDateIndex.value].dateString
  return carbonNutritionStore.getDataByDate(selectedDate) || defaultEmptyData()
}

function defaultEmptyData() {
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

// 更新摄入总体
function updateSummaryNutrition() {
  const dateData = getDataForSelectedDate()
  const actualNutrients = dateData.nutrients?.actual || {}
  const targetNutrients = dateData.nutrients?.target || {}
  const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
  const tempSummary = nutrients.map(nutrient => {
    const intake = actualNutrients[nutrient] || 0
    const plan = targetNutrients[nutrient] || 0
    let ratio = 0
    if (plan > 0) {
      ratio = intake / plan
      if (ratio > 1) ratio = 1
    } else {
      ratio = intake > 0 ? 1 : 0
    }
    const over = plan > 0 && intake > plan
    const color = over ? getRedShade(nutrient) : getNutrientColor(nutrient)
    return {
      nutrient,
      label: t(`${nutrient}_unit`),
      intake,
      plan,
      ratio,
      color,
      over
    }
  })
  summaryNutrition.value = JSON.parse(JSON.stringify(tempSummary))

  // 每餐
  const m = dateData.meals || {}
  breakfastNutrients.value = mapMealNutrients(m.breakfast?.nutrients || {})
  lunchNutrients.value = mapMealNutrients(m.lunch?.nutrients || {})
  dinnerNutrients.value = mapMealNutrients(m.dinner?.nutrients || {})
  otherNutrients.value = mapMealNutrients(m.other?.nutrients || {})
}

function mapMealNutrients(mealN) {
  const n = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
  return n.map(key => ({
    label: t(`${key}_unit`),
    intake: mealN[key] || 0
  }))
}

// 更新图表数据
function updateChartData() {
  const tempSeries = summaryNutrition.value.map(item => ({
    data: item.ratio,
    color: item.color
  }))
  chartData.value = {series: tempSeries}
}

function getRedShade(nutrient) {
  const shades = {
    calories: '#FF4D4F',
    protein: '#FF7875',
    fat: '#FFB3BA',
    carbohydrates: '#FFA39E',
    sodium: '#FF4D4F'
  }
  return shades[nutrient] || '#FF4D4F'
}

function getNutrientColor(nutrient) {
  const colors = {
    calories: '#1890FF',
    protein: '#91CB74',
    fat: '#FAC858',
    carbohydrates: '#73C0DE',
    sodium: '#3CA272'
  }
  return colors[nutrient] || '#000'
}

// 选择日期
function onDateChange(index) {
  currentDateIndex.value = index
  updateSummaryNutrition()
  updateChartData()
}

// 折叠面板事件
function onMealChange(name) {
  // console.log('当前展开的餐：', name)
}

// 设置目标
function navigateToSetGoals() {
  uni.navigateTo({url: '/pagesMy/setGoals/setGoals'})
}

/* ----------------- Watch ----------------- */
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
  color: #4caf50ff;
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
  background-color: #4caf50;
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
  color: #4caf50;
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
  background-color: #4caf50;
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
  color: #4caf50;
  font-weight: 500;
}

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

/* 隐藏滚动条 */
::-webkit-scrollbar {
  width: 0;
  height: 0;
  background: transparent;
}

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