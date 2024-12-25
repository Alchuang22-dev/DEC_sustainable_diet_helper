<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>

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
        <!-- 替换 days 为 registeredDays -->
        <text class="carbon-number">{{ registeredDays }}{{$t('carbon_days')}}</text>
      </view>
      <view class="charts">

        <!-- 今日碳排放环形图 -->
        <view class="chart today">
          <text class="chart-title">{{$t('carbon_today')}}</text>
          <view class="today-charts">
            <qiun-data-charts
              :canvas2d="true"
              type="ring"
              :opts="ringOpts"
              :chartData="chartTodayData"
            />
          </view>
        </view>

        <!-- 历史碳排放曲线图 -->
        <view class="chart history">
          <text class="chart-title">{{$t('carbon_history')}}</text>
          <qiun-data-charts
            :canvas2d="true"
            canvas-id="carbonHistoryChart"
            type="line"
            :opts="historyOpts"
            :chartData="chartHistoryData"
          />
        </view>

        <!-- 今日营养情况图表（柱状图） -->
        <view class="chart nutrition">
          <text class="chart-title">{{$t('nutrition_today')}}</text>
          <view class="nutrition-charts">
            <qiun-data-charts
              :canvas2d="true"
              canvas-id="xvmMWWFdeOdEnvVDPjotdobEUaWVmvav"
              type="column"
              :ontouch="true"
              :opts="nutritionOpts"
              :chartData="chartNutritionData"
            />
          </view>
        </view>

      </view>
    </view>

    <!-- 实用工具 -->
    <view class="useful-tools">
      <text class="tools-title">{{$t('tools_title')}}</text>
      <view class="tools-grid">
        <view class="tool" @click="navigateTo('calculator')" animation="fadeInUp">
          <image
            src="https://cdn.pixabay.com/photo/2017/07/06/17/13/calculator-2478633_1280.png"
            :alt="$t('tool_carbon_calculator')"
            class="tool-icon"
            mode="aspectFill"
          />
          <view class="tool-description">
            <text class="tool-name">{{$t('tool_carbon_calculator')}}</text>
            <text class="tool-info">{{$t('tool_carbon_calculator_info')}}</text>
          </view>
        </view>
        <view class="tool" @click="navigateTo('recommend')" animation="fadeInUp" animation-delay="0.2s">
          <image
            src="https://cdn.pixabay.com/photo/2020/03/12/18/37/dish-4925892_1280.png"
            :alt="$t('tool_diet_recommendation')"
            class="tool-icon"
            mode="aspectFill"
          />
          <view class="tool-description">
            <text class="tool-name">{{$t('tool_diet_recommendation')}}</text>
            <text class="tool-info">{{$t('tool_diet_recommendation_info')}}</text>
          </view>
        </view>
        <view class="tool" @click="navigateTo('nutrition')" animation="fadeInUp" animation-delay="0.4s">
          <image
            src="https://cdn.pixabay.com/photo/2016/11/14/15/42/calendar-1823848_1280.png"
            :alt="$t('tool_nutrition_calculator')"
            class="tool-icon"
          />
          <view class="tool-description">
            <text class="tool-name">{{$t('tool_nutrition_calculator')}}</text>
            <text class="tool-info">{{$t('tool_nutrition_calculator_info')}}</text>
          </view>
        </view>
        <view class="tool" @click="navigateTo('family')" animation="fadeInUp" animation-delay="0.6s">
          <image
            src="https://cdn.pixabay.com/photo/2016/01/04/14/24/terminal-board-1120961_1280.png"
            :alt="$t('tool_family_recipe')"
            class="tool-icon"
          />
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
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCarbonAndNutritionStore } from '@/stores/carbon_and_nutrition_data.js'
import { onShow } from '@dcloudio/uni-app';
import { useUserStore } from '@/stores/user.js'; // 引入用户 store

const { t, locale } = useI18n()
const showMore = ref(false)

// 获取 Pinia 存储
const carbonNutritionStore = useCarbonAndNutritionStore()
const userStore = useUserStore() // 获取用户 store

// 使用计算属性从用户 store 获取 registered_days
const registeredDays = computed(() => userStore.user.registered_days)

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

// 样式设定作为 computed 属性，使其依赖于 t()
const nutritionOpts = computed(() => ({
  color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE"],
  padding: [15, 15, 0, 15],
  xAxis: {
    disableGrid: false,
    axisLine: true,
    itemCount: 4,
    rotateLabel: true,
    rotateAngle: 60 // 保持倾斜角度不变
  },
  yAxis: { disabled: true },
  extra: {
    column: {
      width: 20,
      type: 'group',
      seriesGap: 5,
    }
  },
  legend: {
    show: true,
    position: "bottom",
    lineHeight: 25
  },
  title: {
    name: t('nutrition_title'), // 假设有对应的翻译键
    fontSize: 15,
    color: "#666666"
  }
}))

const historyOpts = computed(() => ({
  yAxis: {
    disabled: true
  },
  legend: {
    show: true,
    position: "bottom",
    lineHeight: 25
  },
  title: {
    name: t('carbon_history'),
    fontSize: 15,
    color: "#666666"
  }
}))

const ringOpts = computed(() => ({
  rotate: false,
  rotateLock: false,
  color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
  animation: { duration: 0 },
  padding: [5, 5, 5, 5],
  dataLabel: true,
  enableScroll: false,
  legend: { show: true, position: "bottom", lineHeight: 25 },
  title: {
    name: t('total'), // 多语言支持
    fontSize: 15,
    color: "#666666"
  },
  subtitle: {
    name: ringSubtitle.value, // 动态更新
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
}))

// 动态 subtitle，响应式
const ringSubtitle = ref("")

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

// 用于渲染今日和历史图表
function renderCharts() {
  const today = new Date()
  console.log('today', today)
  const dateString = today.getFullYear() + '-'
    + String(today.getMonth() + 1).padStart(2, '0') + '-'
    + String(today.getDate()).padStart(2, '0');
  console.log('dateString', dateString)
  const todayData = getDataByDate(dateString)
  console.log('todayData', todayData)

  // ----------------
  // 1. 今日环形图
  // ----------------
  if (todayData) {
    const mealTypes = ['breakfast', 'lunch', 'dinner', 'other']
    let mealData = []
    let totalCarbonEmission = 0

    mealTypes.forEach(mealType => {
      const emission = todayData.meals[mealType].carbonEmission || 0
      mealData.push({
        name: t(mealType), // 切换语言时需重新翻译
        value: emission
      })
      totalCarbonEmission += emission
    })

    // 更新 chartTodayData
    chartTodayData.value = {
      series: [
        { data: mealData }
      ]
    }

    // 更新 ringOpts 的 subtitle
    ringSubtitle.value = `${totalCarbonEmission.toFixed(1)}Kg`
  }

  // ----------------
  // 2. 今日营养柱状图
  // ----------------
  if (todayData) {
    const nutrients = ['calories', 'protein', 'fat', 'carbohydrates', 'sodium']
    const intakeData = nutrients.map(n => todayData.nutrients.actual[n] || 0)
    const targetData = nutrients.map(n => todayData.nutrients.target[n] || 0)

    chartNutritionData.value = {
      categories: [
        t('energy_unit'),
        t('protein_unit'),
        t('fat_unit'),
        t('carbohydrates_unit'),
        t('sodium_unit')
      ],
      series: [
        { name: t('intake'), data: intakeData },
        { name: t('target_value'), data: targetData }
      ]
    }
  }

  // ----------------
  // 3. 历史碳排放曲线图
  // ----------------
  const categories = []
  const targetData = []
  const actualData = []

  for (let i = 6; i >= 0; i--) {
    const d = new Date()
    d.setDate(d.getDate() - i)

    const year = d.getFullYear()
    const month = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const ds = `${year}-${month}-${day}`
    console.log('date string (ds)', ds)

    categories.push(`${month}/${day}`)

    const dailyData = getDataByDate(ds)
    targetData.push(dailyData ? dailyData.carbonEmission.target : 0)
    actualData.push(dailyData ? dailyData.carbonEmission.actual : 0)
  }

  chartHistoryData.value = {
    categories: categories,
    series: [
      { name: t('target_value'), data: targetData },
      { name: t('actual_value'), data: actualData }
    ]
  }
}

// 监听语言切换，当 locale 改变时重新渲染图表文本
watch(locale, () => {
  renderCharts()
})

// 监听页面显示，获取数据后进行首次渲染
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
  console.log('onShow')

  // 获取数据
  await carbonNutritionStore.getNutritionGoals()
  await carbonNutritionStore.getCarbonGoals()
  await carbonNutritionStore.getNutritionIntakes()
  await carbonNutritionStore.getCarbonIntakes()

  // 调用渲染函数
  renderCharts()
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