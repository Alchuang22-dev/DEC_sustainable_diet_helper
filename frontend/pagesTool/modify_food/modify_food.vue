<template>
  <view class="container">
    <view class="form-container">
      <!-- 食物名称 -->
      <view class="form-group">
        <text class="label">{{ $t('name') }}</text>
        <uni-combox
          :placeholder="$t('please_enter_food_name')"
          v-model="foodNameInput"
          :candidates="filteredFoods.map(item => displayName(item))"
          @input="onComboxInput"
        ></uni-combox>
      </view>

      <!-- 重量 -->
      <view class="form-group">
        <text class="label">{{ $t('total_weight') }}</text>
        <input
          class="input"
          type="digit"
          v-model="food.weight"
          :placeholder="$t('please_enter_food_weight')"
          :error="weightError"
        />
        <text
          v-if="weightError"
          class="error-message"
        >
          {{ $t('weight_must_be_positive_integer') }}
        </text>
      </view>

      <!-- 价格 -->
      <view class="form-group">
        <text class="label">{{ $t('total_price') }}</text>
        <input
          class="input"
          type="digit"
          v-model="food.price"
          :placeholder="$t('please_enter_food_price')"
          :error="priceError"
        />
        <text
          v-if="priceError"
          class="error-message"
        >
          {{ $t('price_must_be_positive_integer') }}
        </text>
      </view>

      <!-- 运输方式 -->
      <view class="form-group">
        <text class="label">{{ $t('select_transport_method') }}</text>
        <picker
          mode="selector"
          :range="transportMethods"
          :value="transportIndex"
          @change="onTransportChange"
        >
          <view class="picker">
            {{ transportMethods[transportIndex] }}
          </view>
        </picker>
      </view>

      <!-- 食品来源 -->
      <view class="form-group">
        <text class="label">{{ $t('select_food_source') }}</text>
        <picker
          mode="selector"
          :range="foodSources"
          :value="sourceIndex"
          @change="onSourceChange"
        >
          <view class="picker">
            {{ foodSources[sourceIndex] }}
          </view>
        </picker>
      </view>

      <!-- 提交 -->
      <button class="submit-button" @click="submitFoodDetails">
        {{ $t('submit') }}
      </button>
    </view>
  </view>
</template>

<script setup>
/**
 * 修改食物页面：与 add_food 类似，但用于编辑已有的食物项
 */
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFoodListStore } from '../stores/food_list'
import { onLoad } from '@dcloudio/uni-app'

const { t, locale } = useI18n()
const foodStore = useFoodListStore()

// 解构
const {
  availableFoods,
  fetchAvailableFoods,
  updateFood,
  getFoodName
} = foodStore

// 路由参数
const options = ref({})
const foodIndex = ref(null)
const existingFood = ref(null)

const food = reactive({
  name: '',
  id: null,
  weight: '',
  price: '',
  transportMethod: 'land',
  foodSource: 'local'
})

// 食物名称
const foodNameInput = ref('')

// 下拉列表可见性
const showFoodList = ref(false)

// 可选食物过滤
const filteredFoods = computed(() => {
  if (foodNameInput.value === '') {
    return availableFoods
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

// 显示名称
function displayName(item) {
  return locale.value === 'zh-Hans' ? item.name_zh : item.name_en
}

// combox输入
function onComboxInput(value) {
  foodNameInput.value = value
}

// 验证错误
const weightError = ref(false)
const priceError = ref(false)

// 运输方式、食品来源
const transportMethods = [
  t('transport_land'),
  t('transport_sea'),
  t('transport_air')
]
const foodSources = [
  t('source_local'),
  t('source_imported')
]

// 当前索引
const transportIndex = ref(0)
const sourceIndex = ref(0)

// 运输方式
function onTransportChange(e) {
  transportIndex.value = e.detail.value
  food.transportMethod = e.detail.value === 0 ? 'land' : e.detail.value === 1 ? 'sea' : 'air'
}

// 食品来源
function onSourceChange(e) {
  sourceIndex.value = e.detail.value
  food.foodSource = e.detail.value === 0 ? 'local' : 'imported'
}

/**
 * 提交修改
 */
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

  weightError.value = false
  priceError.value = false

  let valid = true
  if (!/^\d+(\.\d+)?$/.test(food.weight) || parseFloat(food.weight) <= 0) {
    weightError.value = true
    valid = false
  }
  if (!/^\d+(\.\d+)?$/.test(food.price) || parseFloat(food.price) <= 0) {
    priceError.value = true
    valid = false
  }
  if (
    !food.name ||
    !food.weight ||
    !food.price ||
    !food.transportMethod ||
    !food.foodSource
  ) {
    uni.showToast({
      title: t('please_fill_all_fields'),
      icon: 'none'
    })
    valid = false
  }
  if (!valid) return

  const updatedFood = {
    name: food.name,
    id: food.id,
    weight: parseFloat(food.weight),
    price: parseFloat(food.price),
    transportMethod: food.transportMethod,
    foodSource: food.foodSource
  }
  updateFood(foodIndex.value, updatedFood)

  uni.showToast({
    title: t('modify_success'),
    icon: 'success',
    duration: 2000
  })

  setTimeout(() => {
    uni.navigateBack()
  }, 2000)
}

// 接收路由参数
onLoad((loadedOptions) => {
  options.value = loadedOptions
  foodIndex.value = parseInt(options.value.index)
  existingFood.value = foodStore.foodList[foodIndex.value]

  // 索引无效
  if (isNaN(foodIndex.value) || !existingFood.value) {
    uni.showToast({
      title: t('invalid_food_item'),
      icon: 'none',
      duration: 2000
    })
    setTimeout(() => {
      uni.navigateBack()
    }, 2000)
    return
  }

  // 初始化
  food.name = existingFood.value.name || ''
  food.id = existingFood.value.id || null
  food.weight = existingFood.value.weight || ''
  food.price = existingFood.value.price || ''
  food.transportMethod = existingFood.value.transportMethod || 'land'
  food.foodSource = existingFood.value.foodSource || 'local'

  foodNameInput.value = getFoodName(food.id)
  transportIndex.value = ['land', 'sea', 'air'].indexOf(food.transportMethod)
  sourceIndex.value = ['local', 'imported'].indexOf(food.foodSource)
})

// 若后端数据为空，则获取
onMounted(() => {
  if (availableFoods.length === 0) {
    fetchAvailableFoods()
  }
})
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