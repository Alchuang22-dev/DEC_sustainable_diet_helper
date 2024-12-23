<template>
  <view class="container">
    <view class="goals-wrapper">
      <!-- 碳目标输入框 -->
      <view class="goal-section">
        <text class="label">{{ t('carbon_goal_label') }}</text>
        <input
            class="input"
            type="number"
            v-model.trim="carbonGoalStr"
            :maxlength="10"
        />
      </view>

      <!-- 分隔线 -->
      <view class="divider"></view>

      <!-- 五个营养目标输入框 -->
      <view class="goals-grid">
        <view class="goal-section">
          <text class="label">{{ t('calories_unit') }}</text>
          <input
              class="input"
              type="number"
              v-model.trim="caloriesStr"
              :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('protein_unit') }}</text>
          <input
              class="input"
              type="number"
              v-model.trim="proteinStr"
              :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('fat_unit') }}</text>
          <input
              class="input"
              type="number"
              v-model.trim="fatStr"
              :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('carbohydrates_unit') }}</text>
          <input
              class="input"
              type="number"
              v-model.trim="carbsStr"
              :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('sodium_unit') }}</text>
          <input
              class="input"
              type="number"
              v-model.trim="sodiumStr"
              :maxlength="10"
          />
        </view>
      </view>

      <!-- 底部按钮 -->
      <view class="button-container">
        <button class="btn btn-today" @click="setGoalsForToday">
          {{ t('set_for_today') }}
        </button>
        <button class="btn btn-week" @click="setGoalsForWeek">
          {{ t('set_for_week') }}
        </button>
      </view>

      <!-- 营养建议说明 -->
      <view class="nutrition-guidelines">
        <text class="guidelines-title">营养摄入参考值</text>
        <view class="guidelines-content">
          <text class="guideline-item">成年男性（18-49岁）建议值：</text>
          <text class="guideline-details">· 热量：2250-2400千卡/天</text>
          <text class="guideline-details">· 蛋白质：65-75克/天</text>
          <text class="guideline-details">· 碳水化合物：320-370克/天</text>
          <text class="guideline-details">· 脂肪：60-70克/天</text>
          <text class="guideline-details">· 钠：2000毫克/天</text>

          <text class="guideline-item">成年女性（18-49岁）建议值：</text>
          <text class="guideline-details">· 热量：1800-2000千卡/天</text>
          <text class="guideline-details">· 蛋白质：55-65克/天</text>
          <text class="guideline-details">· 碳水化合物：250-300克/天</text>
          <text class="guideline-details">· 脂肪：50-60克/天</text>
          <text class="guideline-details">· 钠：2000毫克/天</text>
        </view>
      </view>

    </view>
  </view>
</template>

<script setup>
import {ref, onMounted, computed} from 'vue';
import {useI18n} from 'vue-i18n';

import {useCarbonAndNutritionStore} from '@/stores/carbon_and_nutrition_data';
import {useUserStore} from '@/stores/user';

// 获取国际化函数
const {t} = useI18n();

// 获取 Pinia 存储
const carbonAndNutritionStore = useCarbonAndNutritionStore();
const userStore = useUserStore();

// 从 store 中获取 token（如果需要手动使用）
const token = computed(() => userStore.user.token);

// 碳目标和五大营养目标的绑定值（以字符串形式做双向绑定，提交前再转成整数）
const carbonGoalStr = ref('0');
const caloriesStr = ref('0');
const proteinStr = ref('0');
const fatStr = ref('0');
const carbsStr = ref('0');
const sodiumStr = ref('0');

// 日期格式化函数（与首页相同）
const getFormattedDate = (today) => {
  return today.getFullYear() + '-'
      + String(today.getMonth() + 1).padStart(2, '0') + '-'
      + String(today.getDate()).padStart(2, '0');
};

// onMounted 时，获取当日目标并填充
onMounted(async () => {
  try {
    // 拉取已有的碳和营养目标
    await carbonAndNutritionStore.getCarbonGoals();
    await carbonAndNutritionStore.getNutritionGoals();

    const today = new Date();
    const dateString = getFormattedDate(today);

    // 找到当天的碳目标
    const carbonGoalObj = carbonAndNutritionStore.state.carbonGoals.find(item => {
      return item.date.startsWith(dateString);
    });
    if (carbonGoalObj) {
      carbonGoalStr.value = String(Math.round(carbonGoalObj.emission || 0));
    } else {
      carbonGoalStr.value = '0';
    }

    // 找到当天的营养目标
    const nutritionGoalObj = carbonAndNutritionStore.state.nutritionGoals.find(item => {
      return item.date.startsWith(dateString);
    });
    if (nutritionGoalObj) {
      caloriesStr.value = String(Math.round(nutritionGoalObj.calories || 0));
      proteinStr.value = String(Math.round(nutritionGoalObj.protein || 0));
      fatStr.value = String(Math.round(nutritionGoalObj.fat || 0));
      carbsStr.value = String(Math.round(nutritionGoalObj.carbohydrates || 0));
      sodiumStr.value = String(Math.round(nutritionGoalObj.sodium || 0));
    } else {
      caloriesStr.value = '0';
      proteinStr.value = '0';
      fatStr.value = '0';
      carbsStr.value = '0';
      sodiumStr.value = '0';
    }

  } catch (err) {
    console.error(err);
    uni.showToast({
      title: t('fetchGoalsFail'),
      icon: 'none'
    });
  }
});

/**
 * 点击「为今天设置」
 * 1. 检验输入是否合法（都应是整数且>=0）。
 * 2. 构造仅当天的 goals 数组并调用 setNutritionGoals、setCarbonGoals。
 */
const setGoalsForToday = async () => {
  if (!validateInputs()) return;

  const today = new Date();
  const dateString = getFormattedDate(today);

  const carbonGoals = [
    {
      date: `${dateString}T00:00:00Z`,
      emission: parseInt(carbonGoalStr.value, 10),
    }
  ];
  const nutritionGoals = [
    {
      date: `${dateString}T00:00:00Z`,
      calories: parseInt(caloriesStr.value, 10),
      protein: parseInt(proteinStr.value, 10),
      fat: parseInt(fatStr.value, 10),
      carbohydrates: parseInt(carbsStr.value, 10),
      sodium: parseInt(sodiumStr.value, 10)
    }
  ];

  try {
    await carbonAndNutritionStore.setCarbonGoals(carbonGoals);
    await carbonAndNutritionStore.setNutritionGoals(nutritionGoals);
    uni.showToast({
      title: t('setSuccess'),
      icon: 'success'
    });
  } catch (err) {
    console.error(err);
    uni.showToast({title: t('setFail'), icon: 'none'});
  }
};

/**
 * 点击「为一周设置」
 * 1. 检验输入是否合法（都应是整数且>=0）。
 * 2. 从今天开始，连续 7 天（今天 + 后面 6 天）都设置相同的目标。
 */
const setGoalsForWeek = async () => {
  if (!validateInputs()) return;

  const carbonGoals = [];
  const nutritionGoals = [];

  const today = new Date();

  for (let i = 0; i < 7; i++) {
    const date = new Date(today);
    date.setDate(today.getDate() + i);
    const dateString = getFormattedDate(date);
    console.log('dateString:', dateString);

    carbonGoals.push({
      date: `${dateString}T00:00:00Z`,
      emission: parseInt(carbonGoalStr.value, 10),
    });
    nutritionGoals.push({
      date: `${dateString}T00:00:00Z`,
      calories: parseInt(caloriesStr.value, 10),
      protein: parseInt(proteinStr.value, 10),
      fat: parseInt(fatStr.value, 10),
      carbohydrates: parseInt(carbsStr.value, 10),
      sodium: parseInt(sodiumStr.value, 10)
    });
  }

  try {
    await carbonAndNutritionStore.setCarbonGoals(carbonGoals);
    await carbonAndNutritionStore.setNutritionGoals(nutritionGoals);
    uni.showToast({
      title: t('setSuccess'),
      icon: 'success'
    });
  } catch (err) {
    console.error(err);
    uni.showToast({title: t('setFail'), icon: 'none'});
  }
};

// 校验输入框是否都是有效整数
function validateInputs() {
  // 如果有非数字或负数或非整数
  const items = [
    carbonGoalStr.value,
    caloriesStr.value,
    proteinStr.value,
    fatStr.value,
    carbsStr.value,
    sodiumStr.value
  ];
  for (const val of items) {
    if (
        val === '' ||
        isNaN(Number(val)) ||
        Number(val) < 0 ||
        !Number.isInteger(Number(val))
    ) {
      uni.showToast({title: t('inputValidateFail'), icon: 'none'});
      return false;
    }
  }
  return true;
}
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 12px 8px;
}

.goals-wrapper {
  background-color: #ffffff;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.goal-section {
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.label {
  font-size: 13px;
  color: #374151;
  font-weight: 500;
}

.input {
  padding: 6px 10px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  font-size: 13px;
  background-color: #f9fafb;
  transition: all 0.2s ease;
}

.input:focus {
  border-color: #60a5fa;
  background-color: #ffffff;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.1);
}

.divider {
  height: 1px;
  background-color: #e5e7eb;
  margin: 12px 0;
}

.goals-grid {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.button-container {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  margin-bottom: 16px;
}

.btn {
  flex: 1;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  justify-content: center;
  align-items: center;
}

.btn-today {
  background-color: #3b82f6;
  color: #ffffff;
}

.btn-today:hover {
  background-color: #2563eb;
}

.btn-week {
  background-color: #10b981;
  color: #ffffff;
}

.btn-week:hover {
  background-color: #059669;
}

/* 营养建议说明样式 */
.nutrition-guidelines {
  margin-top: 8px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

.guidelines-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  display: block;
}

.guidelines-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.guideline-item {
  font-size: 13px;
  color: #374151;
  font-weight: 500;
  margin-top: 8px;
}

.guideline-details {
  font-size: 12px;
  color: #6b7280;
  padding-left: 8px;
}

@media (max-width: 640px) {
  .container {
    padding: 8px 6px;
  }

  .goals-wrapper {
    padding: 10px 8px;
  }

  .button-container {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>