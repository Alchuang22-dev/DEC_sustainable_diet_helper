<template>
  <view class="container">
    <!-- Daily Nutrition Target Section -->
    <view class="daily-nutrition-target">
      <button class="toggle-button" @click="toggleTargetList">{{ $t('viewNutritionGoals') }}</button>
      <view v-if="showTargetList" class="nutrition-item-list">
        <view v-for="(target, index) in nutritionTargets" :key="index" class="nutrition-item">
          <text>{{ target.name }}    {{ target.currentValue }}{{ target.unit }}/{{ target.targetValue }}{{ target.unit }}</text>
          <button class="delete-target" @click="deleteTarget(index)">&times;</button>
        </view>
      </view>
      <view class="target-input-container">
        <!-- Picker for selecting target name -->
        <picker mode="selector" :range="pickerOptions" :value="selectedIndex" @change="onTargetNameChange">
          <view class="picker">{{ selectedTarget }}</view>
        </picker>
        <input v-model.number="targetValue" type="number" :placeholder="getTargetValuePlaceholder()" />
        <!-- Picker for selecting unit -->
        <picker mode="selector" :range="unitOptions" :value="unitIndex" @change="onUnitChange">
          <view class="picker">{{ unit }}</view>
        </picker>
      </view>
      <view class="add-button-container">
        <button @click="addOrUpdateTarget">{{ $t('addOrUpdateTarget') }}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';

// i18n keys (to be loaded in your i18n configuration)
const { t, locale } = useI18n();

const nutritionData = ref({
  energy: 30.6,
  protein: 20.4,
  fat: 20.4,
  carbohydrates: 20.4,
  sodium: 2
});

const nutritionTargets = ref([
  { name: '蛋白质', currentValue: 113, targetValue: 968, unit: 'g' },
  { name: '碳水化合物', currentValue: 880, targetValue: 1120, unit: 'g' },
  { name: '脂肪', currentValue: 16, targetValue: 32, unit: 'g' },
]);

const showTargetList = ref(false);
const targetValue = ref(null);
const unit = ref('g');
const unitIndex = ref(0);
const selectedIndex = ref(0); // Default index for the nutrition picker
const selectedTarget = ref('蛋白质'); // Default target name

// Define possible units based on the selected target
const unitOptions = ref(['g', 'mg']); // Default units for most nutrients

// Function to get language-specific names for nutrition targets
const getPickerOptions = () => {
  return [
    t('protein'),
    t('carbohydrates'),
    t('fat'),
    t('energy'),
    t('sodium')
  ];
};

const toggleTargetList = () => {
  showTargetList.value = !showTargetList.value;
};

const addOrUpdateTarget = () => {
  if (!selectedTarget.value || !targetValue.value) {
    uni.showToast({ title: t('pleaseFillTargetNameAndValue'), icon: 'none' });
    return;
  }

  const existingTarget = nutritionTargets.value.find(target => target.name === selectedTarget.value);
  
  if (existingTarget) {
    // Update existing target
    existingTarget.targetValue = targetValue.value;
    existingTarget.unit = unit.value;
  } else {
    // Add new target
    nutritionTargets.value.push({
      name: selectedTarget.value,
      currentValue: 0,
      targetValue: targetValue.value,
      unit: unit.value
    });
  }

  targetValue.value = null;
  unit.value = 'g';
  selectedIndex.value = 0; // Reset the index
  selectedTarget.value = t('protein'); // Reset the selected target to default
};

const deleteTarget = (index) => {
  nutritionTargets.value.splice(index, 1);
};

// Handle target name change and adjust the unit accordingly
const onTargetNameChange = (e) => {
  selectedIndex.value = e.detail.value; // Get the index of the selected target
  selectedTarget.value = pickerOptions.value[selectedIndex.value]; // Update the selected target name
  
  // Adjust unit based on the selected target
  if (selectedTarget.value === t('energy')) {
    unit.value = 'kcal';  // Set unit to kcal for Energy
    unitOptions.value = ['kcal']; // Only kcal available for Energy
  } else {
    unit.value = 'g';  // Default unit for other nutrients
    unitOptions.value = ['g', 'mg']; // g and mg for other nutrients
  }
};

// Handle unit change
const onUnitChange = (e) => {
  const selectedUnitIndex = e.detail.value; // Get the selected unit index
  unit.value = unitOptions.value[selectedUnitIndex]; // Update unit based on index
  unitIndex.value = selectedUnitIndex; // Update the unit index
};

// Get dynamic placeholder for target value
const getTargetValuePlaceholder = () => {
  return selectedTarget.value === t('energy') ? 'kcal' : 'g';
};

// Create a reactive value for picker options
const pickerOptions = ref(getPickerOptions());

// Watch for language changes and update picker options
watch(() => locale.value, () => {
  pickerOptions.value = getPickerOptions(); // Update picker options when language changes
});

onMounted(() => {
  // Initialize chart (using Chart.js)
  const ctx = document.getElementById('nutritionChart').getContext('2d');
  new Chart(ctx, {
    type: 'pie',
    data: {
      labels: getPickerOptions(), // Using i18n keys for labels
      datasets: [{
        data: Object.values(nutritionData.value),
        backgroundColor: ['#42a5f5', '#66bb6a', '#ffca28', '#ef5350', '#29b6f6'],
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top',
        }
      }
    }
  });
});
</script>

<style scoped>
.container {
  font-family: 'Arial', sans-serif;
  background-color: #f0f4f7;
  color: #333;
  padding: 0;
}

.daily-nutrition-target {
  margin: 20px;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.nutrition-item {
  position: relative; /* 设置为relative，以便让delete按钮可以绝对定位 */
  padding-right: 40px; /* 为右侧留出空间给删除按钮，避免遮挡文本 */
  padding: 10px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e0e0e0;
}

.nutrition-item button.delete-target {
  position: absolute; /* 设置按钮为绝对定位 */
  top: 50%; /* 垂直居中 */
  right: 10px; /* 将按钮放置在右侧 */
  transform: translateY(-50%); /* 确保按钮垂直居中 */
  font-size: 20px; /* 设置按钮字体大小 */
  color: red; /* 红色字体 */
  background-color: transparent; /* 按钮背景透明 */
  border: none; /* 去除边框 */
  cursor: pointer; /* 手指形状的光标 */
}


.toggle-button {
  cursor: pointer;
  background-color: #ffffff;
  border: none;
  font-size: 20px;
  color: #4CAF50;
}

.target-input-container {
  display: flex;
  margin-top: 20px;
  align-items: center;
}

.target-input-container input {
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
  font-size: 16px;
  margin-right: 10px;
  flex: 1;
}

.target-input-container picker {
  margin-right: 10px;
}

.add-button-container {
  margin-top: 20px;
  text-align: center;
}

.add-button-container button {
  padding: 10px;
  border: none;
  border-radius: 5px;
  background-color: #4CAF50;
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
}

.picker {
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 5px;
}
</style>
