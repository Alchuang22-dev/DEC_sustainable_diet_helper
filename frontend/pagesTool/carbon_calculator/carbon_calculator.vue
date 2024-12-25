<template>
  <view class="page-container" @load="handleLoad">
    <image src="../static/background_img.jpg" class="background-image"></image>

    <view class="header-card">
      <text class="header-title">{{ $t('carbon_calculator') }}</text>
    </view>

    <uni-section :title="t('added_foods')" type="line">
      <view class="content-wrapper">
        <!-- 食物列表 -->
        <scroll-view scroll-y="true" class="food-list">
          <uni-collapse>
            <uni-collapse-item
              v-for="(food, index) in displayFoodList"
              :key="food.id"
              :title="food.displayName || $t('default_food_name')"
              :thumb="food.displayImage || 'https://cdn.pixabay.com/photo/2015/05/16/15/03/tomatoes-769999_1280.jpg'"
            >
              <view class="food-details">
                <image
                  :src="food.displayImage || 'https://cdn.pixabay.com/photo/2015/05/16/15/03/tomatoes-769999_1280.jpg'"
                  class="food-image"
                  mode="aspectFill"
                />
                <view class="food-info">
                  <view class="info-grid">
                    <uni-tag
                      :text="t('weight') + ': ' + (food.weight + 'kg')"
                      type="primary"
                      size="small"
                    />
                    <uni-tag
                      :text="t('price') + ': ' + (food.price + t('yuan'))"
                      type="success"
                      size="small"
                    />
                    <uni-tag
                      v-if="food.transportMethod"
                      :text="t(`transport_${food.transportMethod}`)"
                      type="warning"
                      size="small"
                    />
                    <uni-tag
                      v-if="food.foodSource"
                      :text="t(`source_${food.foodSource}`)"
                      type="info"
                      size="small"
                    />
                  </view>
                  <view class="action-row">
                    <uni-icons
                      type="compose"
                      size="20"
                      color="#2979ff"
                      @click.stop="handleEdit(index)"
                    />
                    <uni-icons
                      type="trash"
                      size="20"
                      color="#dd524d"
                      @click.stop="handleDelete(index)"
                    />
                  </view>
                </view>
              </view>
            </uni-collapse-item>

            <!-- 空列表提示 -->
            <view v-if="displayFoodList.length === 0" class="empty-state">
              <text>{{ $t('no_foods_added') }}</text>
            </view>
          </uni-collapse>
        </scroll-view>

        <!-- 操作按钮 -->
        <view class="action-buttons">
          <uni-row :gutter="10">
            <uni-col :span="8">
              <view class="action-button primary" @click="navigateToAddFood">
                <text>{{ $t('add_food') }}</text>
              </view>
            </uni-col>
            <uni-col :span="8">
              <view class="action-button success" @click="saveData">
                <text>{{ $t('save_additions') }}</text>
              </view>
            </uni-col>
            <uni-col :span="8">
              <view class="action-button warning" @click="calculateData">
                <text>{{ $t('start_calculation') }}</text>
              </view>
            </uni-col>
          </uni-row>
        </view>
      </view>
    </uni-section>

    <!-- 计算结果 -->
    <view class="result" v-if="showResult">
      <uni-section :title="t('results')" type="line">
        <view class="charts-container">
          <view class="chart-wrapper">
            <text class="chart-title">{{ $t('your_carbon_footprint') }}</text>
            <qiun-data-charts
              :canvas2d="true"
              type="ring"
              :opts="ringOpts"
              :chartData="chartEmissionData"
            />
          </view>

          <view class="chart-wrapper">
            <text class="chart-title">{{ $t('your_nutrition_intake') }}</text>
            <qiun-data-charts
              :canvas2d="true"
              type="bar"
              :opts="barOpts"
              :chartData="chartNutritionData"
            />
          </view>

          <view class="action-button-container">
            <view class="action-button primary" @click="handleSaveOptions">
              <text>{{ $t('save') }}</text>
            </view>
          </view>
        </view>
      </uni-section>
    </view>
  </view>
</template>

<script setup>
/**
 * 碳排放计算器页面：展示用户添加的食品列表，并可进行碳排放/营养综合计算
 */

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFoodListStore } from '../stores/food_list'
import { useCarbonAndNutritionStore } from '@/stores/carbon_and_nutrition_data'

// 多语言
const { t, locale } = useI18n()

// Pinia
const foodStore = useFoodListStore()
const carbonNutritionStore = useCarbonAndNutritionStore()

// 解构 store 的部分方法
const {
  foodList,
  deleteFood,
  updateFood,
  saveFoodList,
  loadFoodList,
  fetchAvailableFoods,
  availableFoods,
  calculateNutritionAndEmission,
  getFoodName
} = foodStore

// 是否显示结果图表
const showResult = ref(false)

/**
 * 碳排放环形图
 */
const chartEmissionData = ref({
  series: [
    {
      name: t('co2_emission'),
      data: []
    }
  ]
})

/**
 * 营养条形图
 */
const chartNutritionData = ref({
  categories: [
    t('energy_unit'),
    t('protein_unit'),
    t('fat_unit'),
    t('carbohydrates_unit'),
    t('sodium_unit')
  ],
  series: [
    {
      name: t('intake_value'),
      data: [0, 0, 0, 0, 0]
    },
    {
      name: t('target_value_today'),
      data: [0, 0, 0, 0, 0]
    }
  ]
})

/**
 * 环形图配置
 */
const ringOpts = ref({
  rotate: false,
  rotateLock: false,
  color: ["#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0", "#9966FF"],
  dataLabel: true,
  enableScroll: false,
  legend: {
    show: true,
    position: "right",
    lineHeight: 25
  },
  title: {
    name: t('total_emission'),
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
      activeRadius: 20,
      offsetAngle: 0,
      labelWidth: 15,
      border: false,
      borderWidth: 3,
      borderColor: "#FFFFFF"
    }
  }
})

/**
 * 营养条形图配置
 */
const barOpts = ref({
  color: [
    "#1890FF",
    "#91CB74",
    "#FAC858",
    "#EE6666",
    "#73C0DE",
    "#3CA272",
    "#FC8452",
    "#9A60B4",
    "#ea7ccc"
  ],
  padding: [15, 40, 0, 5],
  enableScroll: false,
  legend: {},
  xAxis: {
    boundaryGap: "justify",
    disableGrid: true,
    min: 0,
    axisLine: false,
    max: 4000,
    disabled: true
  },
  yAxis: {},
  extra: {
    bar: {
      type: "group",
      meterBorde: 1,
      meterFillColor: "#FFFFFF",
      activeBgColor: "#000000",
      activeBgOpacity: 0.08,
      linearType: "custom",
      barBorderCircle: true,
      seriesGap: 2,
      categoryGap: 2
    }
  }
})

/**
 * 计算属性：显示食物列表（根据语言）
 */
const displayFoodList = computed(() => {
  return foodList.map(food => {
    const found = availableFoods.find(f => f.id === food.id)
    const displayName = found
      ? (locale.value === 'zh-Hans' ? found.name_zh : found.name_en)
      : (food.name || t('default_food_name'))

    const displayImage = found?.image_url || ''

    return {
      ...food,
      displayName,
      displayImage
    }
  })
})

/**
 * 动画效果：页面加载时触发
 */
function handleLoad() {
  foodList.forEach((food, index) => {
    setTimeout(() => {
      food.isAnimating = true
      setTimeout(() => {
        food.isAnimating = false
      }, 500)
    }, index * 100)
  })
}

/**
 * 保存列表到本地
 */
function saveData() {
  saveFoodList()
  uni.showToast({
    title: t('save_success'),
    icon: 'success',
    duration: 2000
  })
}

/**
 * 删除某项食物
 */
function handleDelete(index) {
  deleteFood(index)
  uni.showToast({
    title: t('delete_success'),
    icon: 'success',
    duration: 2000
  })
}

/**
 * 跳转到修改页面
 */
function handleEdit(index) {
  uni.navigateTo({
    url: `/pagesTool/modify_food/modify_food?index=${index}`
  })
}

/**
 * 跳转到添加食物页面
 */
function navigateToAddFood() {
  uni.navigateTo({
    url: '/pagesTool/add_food/add_food'
  })
}

/**
 * 小数处理：保留1位
 */
const roundToOneDecimal = (num) => Number(num.toFixed(1))

/**
 * 计算碳排放和营养信息
 */
async function calculateData() {
  try {
    await calculateNutritionAndEmission()

    // 计算总碳排放
    let totalCO2 = 0
    const emissionData = foodList.map(item => {
      totalCO2 += item.emission
      return {
        name: getFoodName(item.id) || t('default_food_name'),
        value: roundToOneDecimal(item.emission)
      }
    })
    chartEmissionData.value.series[0].data = emissionData
    ringOpts.value.subtitle.name = `${roundToOneDecimal(totalCO2)} kg`

    // 计算总营养
    const totalNutrition = {
      calories: 0,
      protein: 0,
      fat: 0,
      carbohydrates: 0,
      sodium: 0
    }
    foodList.forEach(item => {
      totalNutrition.calories += item.calories
      totalNutrition.protein += item.protein
      totalNutrition.fat += item.fat
      totalNutrition.carbohydrates += item.carbohydrates
      totalNutrition.sodium += item.sodium
    })

    chartNutritionData.value.series[0].data = [
      roundToOneDecimal(totalNutrition.calories),
      roundToOneDecimal(totalNutrition.protein),
      roundToOneDecimal(totalNutrition.fat),
      roundToOneDecimal(totalNutrition.carbohydrates),
      roundToOneDecimal(totalNutrition.sodium)
    ]

    // 获取当日目标
    const today = new Date()
    const dateString = [
      today.getFullYear(),
      String(today.getMonth() + 1).padStart(2, '0'),
      String(today.getDate()).padStart(2, '0')
    ].join('-')
    const dateData = carbonNutritionStore.getDataByDate(dateString)

    chartNutritionData.value.series[1].data = [
      roundToOneDecimal(dateData.nutrients.target.calories),
      roundToOneDecimal(dateData.nutrients.target.protein),
      roundToOneDecimal(dateData.nutrients.target.fat),
      roundToOneDecimal(dateData.nutrients.target.carbohydrates),
      roundToOneDecimal(dateData.nutrients.target.sodium)
    ]

    showResult.value = true
    uni.showToast({
      title: t('calculation_success'),
      icon: 'success',
      duration: 2000
    })
  } catch (err) {
    console.error('计算失败:', err)
    uni.showToast({
      title: t('calculation_failed'),
      icon: 'none',
      duration: 2000
    })
  }
}

/**
 * 点击 “保存” 按钮，跳转到保存选项页面
 */
function handleSaveOptions() {
  const calculatedData = {
    carbonEmission: chartEmissionData.value,
    nutrition: chartNutritionData.value
  }

  uni.navigateTo({
    url: `/pagesTool/save_options/save_options?data=${encodeURIComponent(JSON.stringify(calculatedData))}`
  })
}

/**
 * 页面加载后初始化
 */
onMounted(() => {
  if (!foodStore.loaded) {
    loadFoodList()
  }
  fetchAvailableFoods()
  handleLoad()
})
</script>

<style scoped>
:root {
  --primary-color: #4CAF50;
  --secondary-color: #8BC34A;
  --accent-color: #FF9800;
  --text-color: #333;
  --background-color: #f5f5f5;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

/* 页面容器 */
.page-container {
  min-height: 100vh;
  padding: 20rpx;
  box-sizing: border-box;
  background-color: #f5f5f5;
  position: relative;
}

/* 背景图 */
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

/* 头部卡片 */
.header-card {
  margin-bottom: 40rpx;
  padding: 0;
  text-align: center;
}

.header-title {
  font-size: 36rpx;
  color: var(--primary-color);
  font-weight: bold;
  animation: slideDown 1s ease-out;
}

/* 内容容器 */
.content-wrapper {
  background-color: #ffffff;
  border-radius: 15rpx;
  padding: 10rpx;
}

/* 食物列表区域 */
.food-list {
  max-height: 600rpx;
  margin-bottom: 20rpx;
}

.food-details {
  display: flex;
  gap: 20rpx;
  padding: 20rpx;
  background-color: #f8f8f8;
  border-radius: 12rpx;
}

.food-image {
  width: 120rpx;
  height: 120rpx;
  border-radius: 10rpx;
  flex-shrink: 0;
}

.food-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15rpx;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10rpx;
}

:deep(.uni-tag) {
  width: 100%;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  margin: 0;
}

/* 操作行 */
.action-row {
  display: flex;
  justify-content: flex-end;
  gap: 30rpx;
  padding-top: 10rpx;
}

/* 操作按钮组 */
.action-buttons {
  padding: 20rpx 0;
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  padding: 20rpx;
  border-radius: 10rpx;
  color: #ffffff;
  font-size: 24rpx;
  transition: all 0.3s ease;
}

.action-button.primary {
  background-color: #2979ff;
}

.action-button.success {
  background-color: #4caf9d;
}

.action-button.warning {
  background-color: #178d2a;
}

.action-button:active {
  transform: scale(0.98);
  opacity: 0.9;
}

:deep(.uni-collapse) {
  background-color: transparent;
}

:deep(.uni-collapse-item__title) {
  background-color: #ffffff !important;
  border-bottom: 1rpx solid #eee;
}

:deep(.uni-collapse-item__wrap) {
  background-color: #ffffff;
}

/* 结果展示 */
.result {
  margin-top: 30rpx;
  animation: fadeIn 0.5s ease-in-out;
}

.charts-container {
  padding: 20rpx;
  background-color: #ffffff;
  border-radius: 15rpx;
}

.chart-wrapper {
  margin-bottom: 40rpx;
}

.chart-title {
  font-size: 28rpx;
  color: #333;
  font-weight: bold;
  text-align: center;
  margin-bottom: 20rpx;
}

.action-button-container {
  display: flex;
  justify-content: flex-end;
  padding: 20rpx 0;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式 */
@media screen and (max-width: 375px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
  .action-button {
    padding: 15rpx;
    font-size: 22rpx;
  }
  .food-image {
    width: 100rpx;
    height: 100rpx;
  }
}

.food-list::-webkit-scrollbar {
  width: 6rpx;
  background-color: transparent;
}
.food-list::-webkit-scrollbar-thumb {
  background-color: #2979ff;
  border-radius: 3rpx;
}

:deep(.uni-row) {
  margin: -5rpx;
}
:deep(.uni-col) {
  padding: 5rpx;
}

/* 空列表提示 */
.empty-state {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20rpx 0;
}

:deep(.uni-icons) {
  padding: 10rpx;
}
</style>