<template>
  <view class="container">
    <!-- 输入规范提示 -->
    <view class="input-specification">
      <text class="input-spec-text">{{ t('input_specification') }}</text>
    </view>

    <view class="goals-wrapper">
      <!-- 碳目标输入框 -->
      <view class="goal-section">
        <text class="label">{{ t('carbon_goal_label') }}</text>
        <input
          class="input"
          type="digit"
          step="0.1"
          v-model.trim="carbonGoalStr"
          :maxlength="10"
        />
      </view>

      <!-- 分隔线 -->
      <view class="divider" />

      <!-- 五大营养目标 -->
      <view class="goals-grid">
        <view class="goal-section">
          <text class="label">{{ t('calories_unit') }}</text>
          <input
            class="input"
            type="digit"
            step="0.1"
            v-model.trim="caloriesStr"
            :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('protein_unit') }}</text>
          <input
            class="input"
            type="digit"
            step="0.1"
            v-model.trim="proteinStr"
            :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('fat_unit') }}</text>
          <input
            class="input"
            type="digit"
            step="0.1"
            v-model.trim="fatStr"
            :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('carbohydrates_unit') }}</text>
          <input
            class="input"
            type="digit"
            step="0.1"
            v-model.trim="carbsStr"
            :maxlength="10"
          />
        </view>
        <view class="goal-section">
          <text class="label">{{ t('sodium_unit') }}</text>
          <input
            class="input"
            type="digit"
            step="0.1"
            v-model.trim="sodiumStr"
            :maxlength="10"
          />
        </view>
      </view>

      <!-- 设置按钮 -->
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
        <text class="guidelines-title">{{ t('nutrition_guidelines_title') }}</text>
        <view class="guidelines-content">
          <text class="guideline-item">{{ t('adult_male_guidelines') }}</text>
          <text class="guideline-details">{{ t('calories_range_male') }}</text>
          <text class="guideline-details">{{ t('protein_range_male') }}</text>
          <text class="guideline-details">{{ t('carbs_range_male') }}</text>
          <text class="guideline-details">{{ t('fat_range_male') }}</text>
          <text class="guideline-details">{{ t('sodium_range') }}</text>

          <text class="guideline-item">{{ t('adult_female_guidelines') }}</text>
          <text class="guideline-details">{{ t('calories_range_female') }}</text>
          <text class="guideline-details">{{ t('protein_range_female') }}</text>
          <text class="guideline-details">{{ t('carbs_range_female') }}</text>
          <text class="guideline-details">{{ t('fat_range_female') }}</text>
          <text class="guideline-details">{{ t('sodium_range') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCarbonAndNutritionStore } from '@/stores/carbon_and_nutrition_data'
import { useUserStore } from '@/stores/user'

/* ----------------- Setup ----------------- */
const { t } = useI18n()
const carbonAndNutritionStore = useCarbonAndNutritionStore()
const userStore = useUserStore()

/* ----------------- Reactive & State ----------------- */
const token = computed(() => userStore.user.token)
const carbonGoalStr = ref('0')
const caloriesStr = ref('0')
const proteinStr = ref('0')
const fatStr = ref('0')
const carbsStr = ref('0')
const sodiumStr = ref('0')

/* ----------------- Lifecycle ----------------- */
onMounted(async () => {
  try {
    // 拉取已有的碳和营养目标
    await carbonAndNutritionStore.getCarbonGoals()
    await carbonAndNutritionStore.getNutritionGoals()

    const today = new Date()
    const dateString = getFormattedDate(today)

    // 当天碳目标
    const carbonGoalObj = carbonAndNutritionStore.state.carbonGoals.find(item =>
      item.date.startsWith(dateString)
    )
    carbonGoalStr.value = carbonGoalObj
      ? String(parseFloat(carbonGoalObj.emission).toFixed(1))
      : '0.0'

    // 当天营养目标
    const nutritionGoalObj = carbonAndNutritionStore.state.nutritionGoals.find(item =>
      item.date.startsWith(dateString)
    )
    if (nutritionGoalObj) {
      caloriesStr.value = String(parseFloat(nutritionGoalObj.calories || 0).toFixed(1))
      proteinStr.value = String(parseFloat(nutritionGoalObj.protein || 0).toFixed(1))
      fatStr.value = String(parseFloat(nutritionGoalObj.fat || 0).toFixed(1))
      carbsStr.value = String(parseFloat(nutritionGoalObj.carbohydrates || 0).toFixed(1))
      sodiumStr.value = String(parseFloat(nutritionGoalObj.sodium || 0).toFixed(1))
    } else {
      caloriesStr.value = '0.0'
      proteinStr.value = '0.0'
      fatStr.value = '0.0'
      carbsStr.value = '0.0'
      sodiumStr.value = '0.0'
    }
  } catch (err) {
    console.error(err)
    uni.showToast({
      title: t('fetchGoalsFail'),
      icon: 'none'
    })
  }
})

/* ----------------- Methods ----------------- */
function getFormattedDate(today) {
  return (
    today.getFullYear() +
    '-' +
    String(today.getMonth() + 1).padStart(2, '0') +
    '-' +
    String(today.getDate()).padStart(2, '0')
  )
}

async function setGoalsForToday() {
  if (!validateInputs()) return

  const today = new Date()
  const dateString = getFormattedDate(today)

  const carbonGoals = [
    {
      date: `${dateString}T00:00:00Z`,
      emission: parseFloat(carbonGoalStr.value)
    }
  ]
  const nutritionGoals = [
    {
      date: `${dateString}T00:00:00Z`,
      calories: parseFloat(caloriesStr.value),
      protein: parseFloat(proteinStr.value),
      fat: parseFloat(fatStr.value),
      carbohydrates: parseFloat(carbsStr.value),
      sodium: parseFloat(sodiumStr.value)
    }
  ]

  try {
    await carbonAndNutritionStore.setCarbonGoals(carbonGoals)
    await carbonAndNutritionStore.setNutritionGoals(nutritionGoals)
    uni.showToast({
      title: t('success'),
      icon: 'success'
    })
  } catch (err) {
    console.error(err)
    uni.showToast({ title: t('setFail'), icon: 'none' })
  }
}

async function setGoalsForWeek() {
  if (!validateInputs()) return

  const carbonGoals = []
  const nutritionGoals = []

  const today = new Date()
  for (let i = 0; i < 7; i++) {
    const date = new Date(today)
    date.setDate(today.getDate() + i)
    const dateString = getFormattedDate(date)

    carbonGoals.push({
      date: `${dateString}T00:00:00Z`,
      emission: parseFloat(carbonGoalStr.value)
    })
    nutritionGoals.push({
      date: `${dateString}T00:00:00Z`,
      calories: parseFloat(caloriesStr.value),
      protein: parseFloat(proteinStr.value),
      fat: parseFloat(fatStr.value),
      carbohydrates: parseFloat(carbsStr.value),
      sodium: parseFloat(sodiumStr.value)
    })
  }

  try {
    await carbonAndNutritionStore.setCarbonGoals(carbonGoals)
    await carbonAndNutritionStore.setNutritionGoals(nutritionGoals)
    uni.showToast({
      title: t('success'),
      icon: 'success'
    })
  } catch (err) {
    console.error(err)
    uni.showToast({ title: t('fail'), icon: 'none' })
  }
}

function validateInputs() {
  // 非负数、最多1位小数
  const regex = /^\d+(\.\d)?$/
  const items = [
    carbonGoalStr.value,
    caloriesStr.value,
    proteinStr.value,
    fatStr.value,
    carbsStr.value,
    sodiumStr.value
  ]
  for (const val of items) {
    if (
      val === '' ||
      isNaN(Number(val)) ||
      Number(val) < 0 ||
      !regex.test(val)
    ) {
      uni.showToast({ title: t('inputValidateFail'), icon: 'none' })
      return false
    }
  }
  return true
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
.input-specification {
  margin-bottom: 16px;
}
.input-spec-text {
  font-size: 12px;
  color: #6b7280;
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