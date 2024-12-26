<template>
  <view class="container">
    <!-- 头部卡片 -->
    <view class="header-card">
      <text class="title">{{ t('select_save_method') }}</text>
    </view>

    <!-- Picker 部分：选择餐食类型 -->
    <view class="picker-container">
      <text class="picker-label">{{ t('select_meal_type') }}</text>
      <picker
        mode="selector"
        :range="mealTypesDisplay"
        :value="selectedMealIndex"
        @change="onPickerChange"
        class="picker"
      >
        <view class="picker-content">
          {{ mealTypesDisplay[selectedMealIndex] }}
        </view>
      </picker>
    </view>

    <!-- 按钮组 -->
    <view class="button-group">
      <view class="button primary-button" @click="saveForSelf">
        <uni-icons type="check" size="24" color="#fff" />
        <text class="button-text">{{ t('save_for_self') }}</text>
      </view>
      <view class="button secondary-button" @click="saveForFamily">
        <uni-icons type="people" size="24" color="#fff" />
        <text class="button-text">{{ t('save_for_family') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { reactive, ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user.js'
import { useFamilyStore, FamilyStatus } from '../stores/family.js'

/* ----------------- Setup ----------------- */
const { t } = useI18n()
const userStore = useUserStore()
const familyStore = useFamilyStore()

const familyStatus = computed(() => familyStore.family.status)
const uid = computed(() => userStore.user.uid)
const token = computed(() => userStore.user.token)

const BASE_URL = 'http://xcxcs.uwdjl.cn:8080'

/* ----------------- 数据：餐食类型 ----------------- */
const mealTypesDisplay = [t('breakfast'), t('lunch'), t('dinner'), t('other')]
const mealTypesValue = ['breakfast', 'lunch', 'dinner', 'other']
const selectedMealIndex = ref(0)

/* ----------------- onLoad & 页面数据 ----------------- */
const carbonEmissionData = reactive({})
const nutritionData = reactive({})
const totalEmission = ref(0)

onLoad(options => {
  if (options && options.data) {
    try {
      const parsedData = JSON.parse(decodeURIComponent(options.data))
      Object.assign(carbonEmissionData, parsedData.carbonEmission)
      Object.assign(nutritionData, parsedData.nutrition.series[0].data)
      totalEmission.value = carbonEmissionData.series[0].data.reduce(
        (sum, item) => sum + item.value,
        0
      )
    } catch (error) {
      console.error('解析传递的数据失败:', error)
    }
  }
  selectedMealIndex.value = getDefaultMealType()
})

/* ----------------- Methods ----------------- */
function onPickerChange(e) {
  selectedMealIndex.value = e.detail.value
}
function getDefaultMealType() {
  const hour = new Date().getHours()
  if (hour >= 5 && hour < 11) return 0
  if (hour >= 11 && hour < 15) return 1
  if (hour >= 15 && hour < 20) return 2
  return 3
}

function formatToISO(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`
}

// 保存为自己
function saveForSelf() {
  const isoDate = formatToISO(new Date())
  const requestData = {
    date: isoDate,
    meal_type: mealTypesValue[selectedMealIndex.value],
    calories: nutritionData[0] || 0,
    protein: nutritionData[1] || 0,
    fat: nutritionData[2] || 0,
    carbohydrates: nutritionData[3] || 0,
    sodium: nutritionData[4] || 0,
    emission: totalEmission.value || 0,
    user_shares: [
      {
        user_id: uid.value,
        ratio: 1.0
      }
    ]
  }
  uni.request({
    url: 'https://xcxcs.uwdjl.cn/nutrition-carbon/shared/nutrition-carbon',

    method: 'POST',
    data: requestData,
    header: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token.value}`
    },
    success: res => {
      if (res.statusCode === 200) {
        uni.showToast({ title: t('save_success'), icon: 'success', duration: 2000 })
      } else {
        const errorMsg = res.data?.error || t('save_failed')
        uni.showToast({ title: errorMsg, icon: 'none', duration: 2000 })
      }
      setTimeout(() => {
        uni.navigateBack({ delta: 1 })
      }, 2000)
    },
    fail: () => {
      uni.showToast({ title: t('save_failed'), icon: 'none', duration: 2000 })
      setTimeout(() => {
        uni.navigateBack({ delta: 1 })
      }, 2000)
    }
  })
}

// 保存为家庭
async function saveForFamily() {
  try {
    await familyStore.getFamilyDetails()
    if (familyStore.family.status !== FamilyStatus.JOINED) {
      uni.showToast({
        title: t('join_family_first'),
        icon: 'none',
        duration: 2000
      })
      return
    }
    // 跳转到 family_share
    uni.navigateTo({
      url: `/pagesTool/family_share/family_share?data=${encodeURIComponent(
        JSON.stringify({
          carbonEmission: totalEmission.value,
          nutrition: nutritionData,
          mealType: mealTypesValue[selectedMealIndex.value]
        })
      )}`
    })
  } catch (error) {
    uni.showToast({
      title: t('error_fetch_family_details'),
      icon: 'none',
      duration: 2000
    })
  }
}
</script>

<style scoped>
:root {
  --primary-color: #4caf50;
  --secondary-color: #8bc34a;
  --accent-color: #ff9800;
  --text-color: #333;
  --background-color: #f5f5f5;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30rpx;
  background: #f5f5f5;
  height: 100vh;
  box-sizing: border-box;
}

.header-card {
  width: 100%;
  padding: 20rpx;
  margin-bottom: 30rpx;
  text-align: center;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: var(--primary-color);
}

.picker-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  padding: 30rpx;
  border-radius: 15rpx;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
  margin-bottom: 30rpx;
}

.picker-label {
  font-size: 28rpx;
  color: #666;
  margin-bottom: 20rpx;
}

.picker {
  width: 100%;
}

.picker-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx;
  background: #f9f9f9;
  border: 1rpx solid #ddd;
  border-radius: 10rpx;
}

.button-group {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.button {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20rpx;
  border-radius: 20rpx;
  font-size: 28rpx;
  color: #fff;
  cursor: pointer;
  transition: background-color 0.3s;
}

.primary-button {
  background-color: #4caf9d;
}

.secondary-button {
  background-color: #178d2a;
}

.button-text {
  margin-left: 10rpx;
  font-size: 28rpx;
}
</style>