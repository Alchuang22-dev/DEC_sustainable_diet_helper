<template>
  <view class="container">
    <view class="form-container">
      <!-- 食物名称 -->
      <view class="form-group">
        <text class="label">{{ t('name') }}</text>
        <uni-combox
          :placeholder="t('please_enter_food_name')"
          v-model="foodNameInput"
          :candidates="filteredFoods.map(item => displayName(item))"
          @input="onComboxInput"
          class="combox"
        ></uni-combox>
      </view>

      <!-- 重量 -->
      <view class="form-group">
        <text class="label">{{ t('total_weight') }}</text>
        <input
          class="input"
          type="digit"
          v-model="food.weight"
          :placeholder="t('please_enter_food_weight')"
        />
        <text v-if="weightError" class="error-message">
          {{ t('weight_must_be_positive_number') }}
        </text>
      </view>

      <!-- 价格 -->
      <view class="form-group">
        <text class="label">{{ t('total_price') }}</text>
        <input
          class="input"
          type="digit"
          v-model="food.price"
          :placeholder="t('please_enter_food_price')"
        />
        <text v-if="priceError" class="error-message">
          {{ t('price_must_be_positive_number') }}
        </text>
      </view>

      <!-- 运输方式 -->
      <view class="form-group">
        <text class="label">{{ t('select_transport_method') }}</text>
        <picker
          mode="selector"
          :range="transportMethods"
          :value="transportIndex"
          @change="onTransportChange"
          class="picker"
        >
          <view class="picker-content">
            {{ transportMethods[transportIndex] }}
          </view>
        </picker>
      </view>

      <!-- 食品来源 -->
      <view class="form-group">
        <text class="label">{{ t('select_food_source') }}</text>
        <picker
          mode="selector"
          :range="foodSources"
          :value="sourceIndex"
          @change="onSourceChange"
          class="picker"
        >
          <view class="picker-content">
            {{ foodSources[sourceIndex] }}
          </view>
        </picker>
      </view>

      <!-- 提交按钮 -->
      <button class="submit-button" @click="submitFoodDetails">
        {{ t('submit') }}
      </button>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFoodListStore } from '../stores/food_list'

/* ----------------- Setup ----------------- */
const { t, locale } = useI18n()
const foodStore = useFoodListStore()
const { availableFoods, fetchAvailableFoods, addFood } = foodStore

/* ----------------- transportMethods & foodSources ----------------- */
const transportMethods = [
  t('transport_land'),
  t('transport_sea'),
  t('transport_air')
]
const foodSources = [
  t('source_local'),
  t('source_imported')
]

const transportIndex = ref(0)
const sourceIndex = ref(0)

/* ----------------- Reactive & State ----------------- */
const food = reactive({
  name: '',
  id: null,
  weight: '',
  price: '',
  transportMethod: 'land',
  foodSource: 'local'
})

const foodNameInput = ref('')
const weightError = ref(false)
const priceError = ref(false)

/* ----------------- Computed ----------------- */
const filteredFoods = computed(() => {
  if (foodNameInput.value === '') {
    return locale.value === 'zh-Hans'
      ? availableFoods.filter(f => f.name_zh)
      : availableFoods.filter(f => f.name_en)
  } else {
    if (locale.value === 'zh-Hans') {
      return availableFoods.filter(f => f.name_zh.includes(foodNameInput.value))
    } else {
      return availableFoods.filter(f =>
        f.name_en.toLowerCase().includes(foodNameInput.value.toLowerCase())
      )
    }
  }
})

function displayName(item) {
  return locale.value === 'zh-Hans' ? item.name_zh : item.name_en
}

/* ----------------- Lifecycle ----------------- */
onMounted(() => {
  if (availableFoods.length === 0) {
    fetchAvailableFoods()
  }
})

/* ----------------- Methods ----------------- */
function onComboxInput(value) {
  foodNameInput.value = value
}

function onTransportChange(e) {
  // e.detail.value 是选中的下标
  transportIndex.value = e.detail.value
  // 根据下标赋值
  if (transportIndex.value === 0) {
    food.transportMethod = 'land'
  } else if (transportIndex.value === 1) {
    food.transportMethod = 'sea'
  } else if (transportIndex.value === 2) {
    food.transportMethod = 'air'
  }
}

function onSourceChange(e) {
  sourceIndex.value = e.detail.value
  food.foodSource = sourceIndex.value === 0 ? 'local' : 'imported'
}

function submitFoodDetails() {
  const matchedFood = availableFoods.find(f => displayName(f) === foodNameInput.value)
  if (!matchedFood) {
    uni.showToast({
      title: t('no_matching_food'),
      icon: 'none',
      duration: 2000
    })
    return
  }

  food.name = matchedFood.name_en
  food.id = matchedFood.id

  // 重置错误
  weightError.value = false
  priceError.value = false

  // 验证
  let valid = true
  if (!/^\d+(\.\d+)?$/.test(food.weight) || parseFloat(food.weight) <= 0) {
    weightError.value = true
    valid = false
  }
  if (!/^\d+(\.\d+)?$/.test(food.price) || parseFloat(food.price) <= 0) {
    priceError.value = true
    valid = false
  }
  if (!food.name || !food.weight || !food.price || !food.transportMethod || !food.foodSource) {
    uni.showToast({
      title: t('please_fill_all_fields'),
      icon: 'none'
    })
    valid = false
  }

  if (!valid) return

  const newFood = {
    name: food.name,
    id: food.id,
    weight: parseFloat(food.weight),
    price: parseFloat(food.price),
    transportMethod: food.transportMethod,
    foodSource: food.foodSource,
    isAnimating: false,
    emission: 0,
    calories: 0,
    protein: 0,
    fat: 0,
    carbohydrates: 0,
    sodium: 0
  }

  addFood(newFood)
  uni.showToast({
    title: t('add_success'),
    icon: 'success',
    duration: 2000
  })
  setTimeout(() => {
    uni.navigateBack()
  }, 2000)
}
</script>

<style scoped>
:root {
  --primary-color: #4caf50;
  --secondary-color: #8bc34a;
  --text-color: #333;
  --background-color: #f5f5f5;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--background-color);
  font-family: var(--font-family);
}

.form-container {
  margin: 20rpx;
  padding: 30rpx;
  background-color: #ffffff;
  border-radius: 20rpx;
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.1);
  flex-grow: 1;
}

.form-group {
  margin-bottom: 30rpx;
}

.label {
  display: block;
  margin-bottom: 10rpx;
  font-size: 28rpx;
  font-weight: bold;
  color: var(--text-color);
}

.input {
  padding: 20rpx;
  border: 1rpx solid var(--border-color);
  border-radius: 10rpx;
  font-size: 28rpx;
}

.picker {
  padding: 20rpx;
  border: 1rpx solid var(--border-color);
  border-radius: 10rpx;
  font-size: 28rpx;
  color: #666666;
}

.picker-content {
  padding: 20rpx;
}

.combox {
  width: 100%;
}

.submit-button {
  padding: 20rpx;
  border: none;
  background-color: var(--primary-color);
  color: #ffffff;
  font-size: 32rpx;
  border-radius: 30rpx;
  cursor: pointer;
  width: 100%;
  text-align: center;
  transition: background-color 0.3s ease, transform 0.2s ease;
  margin-top: 20rpx;
}

.submit-button:hover {
  background-color: var(--secondary-color);
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
}

.error-message {
  color: #f44336;
  font-size: 24rpx;
  margin-top: 5rpx;
}
</style>