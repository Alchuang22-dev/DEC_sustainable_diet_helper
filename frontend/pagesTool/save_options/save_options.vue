<template>
  <view class="container">
    <!-- 头部卡片 -->
    <view class="header-card">
      <text class="title">{{ t('select_save_method') }}</text>
    </view>

    <!-- Picker 部分 -->
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
        <uni-icons type="check" size="24" color="#fff"></uni-icons>
        <text class="button-text">{{ t('save_for_self') }}</text>
      </view>
      <view class="button secondary-button" @click="saveForFamily">
        <uni-icons type="people" size="24" color="#fff"></uni-icons>
        <text class="button-text">{{ t('save_for_family') }}</text>
      </view>
    </view>
  </view>
</template>


<script setup>
import { reactive, ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from "@/stores/user.js";
import { FamilyStatus, useFamilyStore } from "@/stores/family.js";
import { onLoad } from "@dcloudio/uni-app";
import {formatDate} from "../../uni_modules/uni-dateformat/components/uni-dateformat/date-format";

// 国际化
const { t } = useI18n();
const userStore = useUserStore();
const familyStore = useFamilyStore();
const familyStatus = computed(() => familyStore.family.status);

// 用户信息
const uid = computed(() => userStore.user.uid);
const token = computed(() => userStore.user.token);

// 餐食类型
const mealTypesDisplay = [t('breakfast'), t('lunch'), t('dinner'), t('other')];
const mealTypesValue = ['breakfast', 'lunch', 'dinner', 'other'];
const selectedMealIndex = ref(0);

// Picker 选择改变处理
const onPickerChange = (e) => {
  selectedMealIndex.value = e.detail.value;
};

// 数据处理
const carbonEmissionData = reactive({});
const nutritionData = reactive({});
const totalEmission = ref(0);

// 页面加载时处理传递的数据
onLoad((options) => {
  if (options && options.data) {
    try {
      const parsedData = JSON.parse(decodeURIComponent(options.data));
      console.log('接收到的计算结果:', parsedData);

      Object.assign(carbonEmissionData, parsedData.carbonEmission);
      Object.assign(nutritionData, parsedData.nutrition.series[0].data);
      totalEmission.value = carbonEmissionData.series[0].data.reduce((sum, item) => sum + item.value, 0);
      console.log('Total Emission:', totalEmission.value);
    } catch (error) {
      console.error('解析传递的数据失败:', error);
    }
  }
  selectedMealIndex.value = getDefaultMealType();
});

// 根据当前时间获取默认的餐食类型
const getDefaultMealType = () => {
  const hour = new Date().getHours();
  if (hour >= 5 && hour < 11) return 0; // 早餐
  if (hour >= 11 && hour < 15) return 1; // 午餐
  if (hour >= 15 && hour < 20) return 2; // 晚餐
  return 3; // 其他
};

function formatToISO(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  // 组合成指定格式，仍然使用 Z 后缀
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`;
}

// 保存为自己计算
const saveForSelf = () => {
  const today = new Date();
  const haha = formatToISO(today);
  console.log('当前时间:', haha);
  const requestData = {
    date: formatToISO(today),
    meal_type: mealTypesValue[selectedMealIndex.value],
    calories: nutritionData[0] || 0,
    protein: nutritionData[1] || 0,
    fat: nutritionData[2] || 0,
    carbohydrates: nutritionData[3] || 0,
    sodium: nutritionData[4] || 0,
    emission: totalEmission.value || 0,
    user_shares: [{
      user_id: uid.value,
      ratio: 1.0
    }]
  };

  console.log('发送的数据:', requestData);

  // 发送 POST 请求
  uni.request({
    url: 'http://122.51.231.155:8095/nutrition-carbon/shared/nutrition-carbon',
    method: 'POST',
    data: requestData,
    header: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token.value}`
    },
    success: (res) => {
      if (res.statusCode === 200) {
        uni.showToast({
          title: t('save_success'),
          icon: 'success',
          duration: 2000
        });
        uni.navigateBack();
      } else {
        const errorMsg = res.data?.error || t('save_failed');
        uni.showToast({
          title: errorMsg,
          icon: 'none',
          duration: 2000
        });
      }
    },
    fail: () => {
      uni.showToast({
        title: t('save_failed'),
        icon: 'none',
        duration: 2000
      });
    }
  });
};

// 保存为家人计算
const saveForFamily = () => {
  familyStore.getFamilyDetails();
  if (familyStatus.value !== FamilyStatus.JOINED) {
    uni.showToast({
      title: t('join_family_first'),
      icon: 'none',
      duration: 2000
    });
    return;
  }
  uni.navigateTo({
    url: `/pagesTool/family_share/family_share?data=${encodeURIComponent(JSON.stringify({
      carbonEmission: totalEmission.value,
      nutrition: nutritionData,
      mealType: mealTypesValue[selectedMealIndex.value]
    }))}`
  });
};

// 当前选择的餐食类型
const currentMealType = computed(() => {
  return mealTypesValue[selectedMealIndex.value];
});
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

.picker-card {
  width: 100%;
  padding: 20rpx;
  margin-bottom: 30rpx;
  background: #ffffff;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
  border-radius: 15rpx;
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
