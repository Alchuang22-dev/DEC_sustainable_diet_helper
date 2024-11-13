<template>
  <view class="container">
    <view class="header">
      <text class="back-button" @click="goBack">←</text>
      <text class="title">营养获取概览</text>
      <text class="menu-button">☰</text>
    </view>

    <view class="nutrition-summary">
      <text class="summary-title">您刚刚一餐获取的营养</text>
      <view class="chart-container">
        <qiun-data-charts :canvas2d="true" canvas-id="nutritionChart" type="ring" :opts="ringOpts"
          :chartData="chartNutritionData" />
      </view>
      <view class="nutrition-table">
        <view class="table-row" v-for="(item, index) in nutritionData" :key="index">
          <text>{{ item.name }}</text>
          <text>{{ item.value }}%</text>
        </view>
      </view>
    </view>

    <view class="daily-nutrition-target">
      <button class="toggle-button" @click="toggleNutritionTargets">查看您的一日营养目标</button>
      <view v-if="showNutritionTargets" id="nutrition-target-list">
        <view v-for="(target, index) in nutritionTargets" :key="index" class="nutrition-item">
          <text>{{ target.name }}</text>
          <view class="progress-bar-container">
            <view class="progress-bar" :style="{ width: `${(target.current / target.targetValue) * 100}%` }"></view>
          </view>
          <text class="progress-text">{{ target.current }}/{{ target.targetValue }}{{ target.unit }}</text>
          <button class="delete-target" @click="removeTarget(index)" style="margin-left: auto;">×</button>
        </view>
      </view>
      <view class="target-input-container">
        <input v-model="newTargetName" type="text" placeholder="营养成分" />
        <input v-model.number="newTargetValue" type="number" placeholder="目标值" />
        <picker :value="selectedUnit" :range="unitOptions" @change="onUnitChange">
          <view class="picker">
            {{ unitOptions[selectedUnit] }}
          </view>
        </picker>
        <button @click="addOrUpdateTarget" id="add-target-button">添加/修改目标</button>
      </view>
    </view>
  </view>
</template>

<script>
import { ref, onMounted } from 'vue';

export default {
  setup() {
    const nutritionData = ref([
      { name: '蛋白质', value: 20.4 },
      { name: '碳水化合物', value: 20.4 },
      { name: '脂肪', value: 20.4 },
      { name: '能量', value: 30.6 },
      { name: '脂固醇', value: 4.1 },
      { name: '钠', value: 2 },
    ]);

    const nutritionTargets = ref([
      { name: '蛋白质', current: 113, targetValue: 968, unit: 'g' },
      { name: '碳水化合物', current: 880, targetValue: 1120, unit: 'g' },
      { name: '脂肪', current: 16, targetValue: 32, unit: 'g' },
    ]);

    const showNutritionTargets = ref(false);
    const newTargetName = ref('');
    const newTargetValue = ref(null);
    const unitOptions = ref(['g', 'mg']);
    const selectedUnit = ref(0);

    const chartNutritionData = ref({
      series: [
        {
          name: '营养成分',
          data: nutritionData.value.map((item) => ({ name: item.name, value: item.value })),
        },
      ],
    });

    const ringOpts = ref({
      rotate: false,
      rotateLock: false,
      color: ['#42a5f5', '#66bb6a', '#ffca28', '#ef5350', '#ab47bc', '#29b6f6'],
      dataLabel: true,
      enableScroll: false,
      legend: {
        show: true,
        position: 'right',
        lineHeight: 25,
      },
      title: {
        name: '营养成分占比',
        fontSize: 15,
        color: '#666666',
      },
      subtitle: {
        name: '',
        fontSize: 25,
        color: '#4CAF50',
      },
      extra: {
        ring: {
          ringWidth: 15,
          activeOpacity: 0.5,
          activeRadius: 20,
          offsetAngle: 0,
          labelWidth: 15,
          border: false,
          borderWidth: 3,
          borderColor: '#FFFFFF',
        },
      },
    });

    const goBack = () => {
      uni.navigateBack();
    };

    const toggleNutritionTargets = () => {
      showNutritionTargets.value = !showNutritionTargets.value;
    };

    const addOrUpdateTarget = () => {
      if (!newTargetName.value || !newTargetValue.value) {
        uni.showToast({
          title: '请填写营养成分和目标值',
          icon: 'none',
        });
        return;
      }

      const existingIndex = nutritionTargets.value.findIndex(
        (target) => target.name === newTargetName.value
      );

      if (existingIndex !== -1) {
        nutritionTargets.value[existingIndex].targetValue = newTargetValue.value;
        nutritionTargets.value[existingIndex].unit = unitOptions.value[selectedUnit.value];
      } else {
        nutritionTargets.value.push({
          name: newTargetName.value,
          current: 0,
          targetValue: newTargetValue.value,
          unit: unitOptions.value[selectedUnit.value],
        });
      }

      // Force the view to update by reassigning the nutritionTargets array
      nutritionTargets.value = [...nutritionTargets.value];

      newTargetName.value = '';
      newTargetValue.value = null;
    };

    const removeTarget = (index) => {
      nutritionTargets.value.splice(index, 1);
      // Force the view to update by reassigning the nutritionTargets array
      nutritionTargets.value = [...nutritionTargets.value];
    };

    const onUnitChange = (e) => {
      selectedUnit.value = e.detail.value;
    };

    return {
      nutritionData,
      nutritionTargets,
      showNutritionTargets,
      newTargetName,
      newTargetValue,
      unitOptions,
      selectedUnit,
      chartNutritionData,
      ringOpts,
      goBack,
      toggleNutritionTargets,
      addOrUpdateTarget,
      removeTarget,
      onUnitChange,
    };
  },
};
</script>

<style scoped>
.container {
  font-family: 'Arial', sans-serif;
  background-color: #f0f4f7;
  color: #333;
  padding: 0;
}

.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  justify-content: space-between;
}

.back-button,
.menu-button {
  font-size: 24px;
  cursor: pointer;
}

.nutrition-summary {
  text-align: center;
  margin: 20px;
}

.summary-title {
  background-color: #4caf50;
  color: #ffffff;
  padding: 15px;
  border-radius: 10px;
}

.chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 20px;
}

.chart {
  width: 100%;
  max-width: 500px;
  height: auto;
  margin: 0 auto;
}

.nutrition-table {
  width: 100%;
  max-width: 400px;
  padding: 20px;
  margin: 0 auto;
  text-align: center;
}

.table-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #e0e0e0;
}

.daily-nutrition-target {
  margin: 20px;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.toggle-button {
  cursor: pointer;
  background-color: #ffffff;
  border: none;
  font-size: 20px;
  color: #4caf50;
}

.nutrition-item {
  border-bottom: 1px solid #e0e0e0;
  padding: 10px 0;
  display: flex;
  align-items: center;
}

.progress-bar-container {
  width: 100%;
  max-width: 200px;
  height: 10px;
  background-color: #e0e0e0;
  border-radius: 5px;
  overflow: hidden;
  margin: 0 10px;
}

.progress-bar {
  height: 100%;
  background-color: #4caf50;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 14px;
  color: #666;
}

.target-input-container {
  display: flex;
  flex-direction: column;
  margin-top: 20px;
}

.target-input-container input[type='text'],
.target-input-container input[type='number'],
.picker {
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  font-size: 16px;
  margin-bottom: 10px;
}

.target-input-container button {
  padding: 10px;
  border: none;
  border-radius: 5px;
  background-color: #4caf50;
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
}

.delete-target {
  cursor: pointer;
  background-color: transparent;
  border: none;
  font-size: 20px;
  color: #f44336;
  margin-left: auto;
}
</style>
