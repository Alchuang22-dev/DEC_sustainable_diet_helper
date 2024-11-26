<!-- pagesTool/carbon_calculator/carbon_calculator.vue -->
<template>
  <view class="container" @load="handleLoad">
    <!-- 全屏背景图片 -->
    <image src="../static/background_img.jpg" class="background-image"></image>

    <!-- 头部标题 -->
    <view class="header">
      <text class="header-title">{{ $t('carbon_calculator') }}</text>
    </view>

    <!-- 已添加的食物标题 -->
    <text class="list-title">{{ $t('added_foods') }}</text>

    <!-- 可滑动的食物列表 -->
    <scroll-view scroll-y="true" class="food-list scroll-view">
      <view v-for="(food, index) in foodList" :key="index" class="card-container">
        <uni-card
            :title="food.name || $t('default_food_name')"
            :thumbnail="food.image || 'https://cdn.pixabay.com/photo/2015/05/16/15/03/tomatoes-769999_1280.jpg'"
            :sub-title="`${$t('weight')}: ${food.weight || '1.2kg'} ${$t('price')}: ${food.price || '5元'}`"
            shadow=1
            @click="animateCard(index)"
            :class="{ clicked: food.isAnimating }"
            :extra="`${food.transportMethod} ${food.foodSource}`"
            :style="{ animationDelay: `${index * 0.1}s` }"
        >
          <view class="card-actions">
            <button class="delete-button" @click.stop="handleDelete(index)">{{ $t('delete') }}</button>
            <button class="edit-button" @click.stop="handleEdit(index)">{{ $t('edit') }}</button>
          </view>
        </uni-card>
      </view>
    </scroll-view>

    <!-- 按钮区 -->
    <view class="button-group">
      <button class="primary-button small-button" @click="navigateToAddFood">
        {{ $t('add_food') }}
      </button>
      <button class="secondary-button small-button" @click="saveData">
        {{ $t('save_additions') }}
      </button>
      <button class="calculate-button small-button" @click="calculateData">
        {{ $t('start_calculation') }}
      </button>
    </view>

    <!-- 结果显示（环形图和条形图） -->
    <view class="result" v-if="showResult">
      <text class="result-title">{{ $t('your_carbon_footprint') }}</text>
      <qiun-data-charts
          :canvas2d="true"
          canvas-id="carbonEmissionChart"
          type="ring"
          :opts="ringOpts"
          :chartData="chartEmissionData"
      />
      <text class="result-title">{{ $t('your_nutrition_intake') }}</text>
      <qiun-data-charts
          :canvas2d="true"
          canvas-id="nutritionChart"
          type="bar"
          :opts="barOpts"
          :chartData="chartNutritionData"
      />
      <button class="save-button" @click="saveEmissionData">{{ $t('save') }}</button>
    </view>

  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useI18n } from 'vue-i18n'; // Import useI18n
import { useFoodListStore } from '@/stores/food_list'; // 引入 Pinia Store

// 使用国际化
const { t } = useI18n();

// 使用 Pinia Store
const foodStore = useFoodListStore();

// 解构需要使用的状态和方法
const { foodList, deleteFood, saveFoodList, loadFoodList, fetchAvailableFoods } = foodStore;

// 碳排放数据，仅包含CO2
const emission = ref({
  CO2: 0,
});

const showResult = ref(false);

// 碳排放环形图数据和配置
const chartEmissionData = ref({
  series: [{
    name: t('co2_emission'),
    data: []
  }]
});

// 营养计算条形图数据和配置
const chartNutritionData = ref({
  categories: [t('energy_unit'), t('protein_unit'), t('fat_unit'), t('carbohydrates_unit'), t('sodium_unit')],
  series: [{
    name: t('actual_value'),
    data: [0, 0, 0, 0, 0]
  },
    {
      name: t('target_value'),
      data: [0, 0, 0, 0, 0]
    }
  ]
});


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
    name: "", // 中心显示总排放量
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
});

// 营养条形图配置
const barOpts = ref({
  color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE", "#3CA272", "#FC8452", "#9A60B4", "#ea7ccc"],
  padding: [15, 30, 0, 5],
  enableScroll: false,
  legend: {},
  xAxis: {
    boundaryGap: "justify",
    disableGrid: false,
    min: 0,
    axisLine: false,
    max: 4000 // Adjusted max value for better visualization
  },
  yAxis: {},
  extra: {
    bar: {
      type: "group",
      width: 30,
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
});

// 页面加载时处理动画
const handleLoad = () => {
  foodList.forEach((food, index) => {
    setTimeout(() => {
      food.isAnimating = true;
      setTimeout(() => {
        food.isAnimating = false;
      }, 500);
    }, index * 100);
  });
};

// 保存数据到本地存储
const saveData = () => {
  saveFoodList();
  uni.showToast({
    title: t('save_success'),
    icon: 'success',
    duration: 2000,
  });
};

// 删除食物项
const handleDelete = (index) => {
  deleteFood(index);
  uni.showToast({
    title: t('delete_success'),
    icon: 'success',
    duration: 2000,
  });
};

// 编辑食物项
const handleEdit = (index) => {
  uni.navigateTo({
    url: `/pagesTool/modify_food/modify_food?index=${index}`,
  });
};

// 添加食物项
const navigateToAddFood = () => {
  uni.navigateTo({
    url: '/pagesTool/add_food/add_food',
  });
};

// 计算碳排放和营养数据
const calculateData = () => {
  // 模拟向后端发送请求
  uni.request({
    url: 'https://mock-api.com/calculateData', // 模拟的后端接口URL
    method: 'POST',
    data: {
      foodList: foodList.map(food => ({
        id: food.id,
        name: food.name,
        weight: food.weight,
        // 其他需要发送的字段
      }))
    },
    success: (res) => {
      if (res.statusCode === 200) {
        const totalData = res.data.totalData;

        // 更新 foodList 中的每个食物项，添加多个字段
        totalData.forEach((item, index) => {
          const foodIndex = foodList.findIndex(food => food.id === item.id);
          if (foodIndex !== -1) {
            foodList[foodIndex].emission = item.emission;
            foodList[foodIndex].calories = item.calories;
            foodList[foodIndex].protein = item.protein;
            foodList[foodIndex].fat = item.fat;
            foodList[foodIndex].carbohydrates = item.carbohydrates;
            foodList[foodIndex].sodium = item.sodium;
          }
        });

        // 更新环形图的数据和总排放量
        let totalCO2 = 0;
        chartEmissionData.value.series[0].data = totalData.map(item => {
          totalCO2 += item.emission;
          return {
            name: item.name,
            value: item.emission
          };
        });

        // 更新环形图中心显示的总排放量
        ringOpts.value.subtitle.name = `${totalCO2} kg`;

        // 更新条形图的营养数据
        const totalNutrition = {
          calories: 0,
          protein: 0,
          fat: 0,
          carbohydrates: 0,
          sodium: 0
        };

        totalData.forEach(item => {
          totalNutrition.calories += item.calories;
          totalNutrition.protein += item.protein;
          totalNutrition.fat += item.fat;
          totalNutrition.carbohydrates += item.carbohydrates;
          totalNutrition.sodium += item.sodium;
        });

        chartNutritionData.value.series[0].data = [
          totalNutrition.calories,
          totalNutrition.protein,
          totalNutrition.fat,
          totalNutrition.carbohydrates,
          totalNutrition.sodium
        ];

        // TODO: 从后端获取用户目标值，现在省去
        chartNutritionData.value.series[1].data = [
          totalNutrition.calories + 100,
          totalNutrition.protein + 100,
          totalNutrition.fat + 100,
          totalNutrition.carbohydrates + 100,
          totalNutrition.sodium + 100
        ];

        // 显示结果
        showResult.value = true;

        // 初始化并绘制环形图
        uni.createSelectorQuery().select('#carbonEmissionChart').fields({
          node: true,
          size: true
        }, (res) => {
          const canvas = res.node;
          const ctx = canvas.getContext('2d');
          const chart = new qCharts({
            canvas: ctx,
            type: 'ring',
            data: chartEmissionData.value,
            options: ringOpts.value
          });
          chart.draw();
        }).exec();

        // 初始化并绘制条形图
        uni.createSelectorQuery().select('#nutritionChart').fields({
          node: true,
          size: true
        }, (res) => {
          const canvas = res.node;
          const ctx = canvas.getContext('2d');
          const chart = new qCharts({
            canvas: ctx,
            type: 'bar',
            data: chartNutritionData.value,
            options: barOpts.value
          });
          chart.draw();
        }).exec();
      } else {
        console.error('计算失败:', res.data.error);
        uni.showToast({
          title: t('calculation_failed'),
          icon: 'none',
          duration: 2000,
        });
      }
    },
    fail: (err) => {
      console.error('请求失败', err);
      uni.showToast({
        title: t('calculation_failed'),
        icon: 'none',
        duration: 2000,
      });
    }
  });
};

// 保存碳排放数据到后端
const saveEmissionData = () => {
  uni.request({
    url: 'https://mock-api.com/saveEmissionData', // 模拟的后端保存接口URL
    method: 'POST',
    data: {
      foodList: foodList.map(food => ({
        id: food.id,
        name: food.name,
        weight: food.weight,
        price: food.price,
        transportMethod: food.transportMethod,
        foodSource: food.foodSource,
        image: food.image,
        emission: food.emission || 0, // 添加emission字段，默认为0
      }))
    },
    success: (res) => {
      if (res.statusCode === 200) {
        uni.showToast({
          title: t('save_success'),
          icon: 'success',
          duration: 2000,
        });
      } else {
        console.error('保存失败:', res.data.error);
        uni.showToast({
          title: t('save_failed'),
          icon: 'none',
          duration: 2000,
        });
      }
    },
    fail: (err) => {
      console.error('保存失败', err);
      uni.showToast({
        title: t('save_failed'),
        icon: 'none',
        duration: 2000,
      });
    }
  });
};

// 页面跳转方法
const navigateTo = (page) => {
  if (page === 'add_food') {
    navigateToAddFood();
  } else if (page === 'recommendMenu') {
    uni.navigateTo({
      url: '/pages/recommendMenu/recommendMenu',
    });
  }
};


// 页面跳转到修改页面
const navigateToModify = (index) => {
  uni.navigateTo({
    url: `/pagesTool/modify_food/modify_food?index=${index}`,
  });
};

// 动画卡片
const animateCard = (index) => {
  const food = foodList[index];
  if (!food) return;

  food.isAnimating = true;

  setTimeout(() => {
    food.isAnimating = false;
  }, 300);
};

// 页面加载时执行
onMounted(() => {
  if (!foodStore.loaded) {
    loadFoodList();
  }
  // 调用获取食物列表的函数
  fetchAvailableFoods();
  handleLoad();
});
</script>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-color: #4CAF50;
  --secondary-color: #8BC34A;
  --accent-color: #FF9800;
  --text-color: #333;
  --background-color: #f5f5f5;
  --border-color: #e0e0e0;
  --font-family: 'Arial', sans-serif;
}

/* 容器 */
.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: var(--background-color);
  font-family: var(--font-family);
  padding-bottom: 50rpx;
  animation: fadeIn 1s ease-in-out;
}

/* 全屏背景图片 */
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

/* 头部标题 */
.header {
  padding: 40rpx 20rpx 20rpx;
  text-align: center;
}

.header-title {
  font-size: 48rpx;
  color: var(--primary-color);
  font-weight: bold;
  animation: slideDown 1s ease-out;
}

/* 已添加的食物列表 */
.food-list {
  max-height: 600rpx;
  margin: 20rpx 0rpx;
  padding: 20rpx 0rpx;
  background-color: #ffffff;
  border-radius: 20rpx;
  box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.list-title {
  margin-left: 10rpx;
  font-size: 30rpx;
  font-weight: bold;
  color: var(--text-color);
  margin-bottom: 20rpx;
  text-align: center;
}

/* uni-card 相关样式 */
.card-actions {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  position: relative;
  width: 100%;
}

.delete-button,
.edit-button {
  font-size: 18rpx;
  cursor: pointer;
  transition: color 0.3s ease;
}

.edit-button:hover {
  color: #8BC34A;
}

.delete-button:hover {
  color: #f44336;
}

/* 按钮区 */
.button-group {
  display: flex;
  justify-content: flex-start;
  margin: 20rpx 20rpx;
  gap: 20rpx;
}

.small-button {
  padding: 15rpx 30rpx;
  border-radius: 30rpx;
  border: none;
  font-size: 20rpx;
  color: #ffffff;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.primary-button {
  background-color: var(--primary-color);
}

.primary-button:hover {
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
}

.secondary-button {
  background-color: var(--secondary-color);
}

.secondary-button:hover {
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
}

/* 计算部分 */
.calculate-section {
  text-align: center;
  margin: 10rpx 0;
}

.calculate-button {
  border-radius: 30rpx;
  background-color: var(--accent-color);
  border: none;
  color: #ffffff;
  font-size: 20rpx;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.calculate-button:hover {
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
}

/* 结果显示（环形图） */
.result {
  position: relative;
  margin: 20rpx 20rpx;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  border-radius: 20rpx;
  box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
  font-size: 32rpx;
  color: var(--text-color);
  text-align: center;
  animation: fadeIn 1s ease-in-out;
}

.result-title {
  font-size: 36rpx;
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 15rpx;
}

/* 保存按钮样式 */
.save-button {
  position: absolute;
  bottom: 10rpx;
  right: 10rpx;
  padding: 10rpx 20rpx;
  background-color: var(--primary-color);
  color: #ffffff;
  border: none;
  border-radius: 20rpx;
  font-size: 24rpx;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.save-button:hover {
  background-color: var(--secondary-color);
  transform: translateY(-2rpx);
  box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

@keyframes slideDown {
  from {
    transform: translateY(-20rpx);
    opacity: 0;
  }

  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@keyframes popIn {
  0% {
    transform: scale(0.95);
    opacity: 0;
  }

  60% {
    transform: scale(1.05);
    opacity: 1;
  }

  100% {
    transform: scale(1);
  }
}

@keyframes fadeInUp {
  from {
    transform: translateY(20rpx);
    opacity: 0;
  }

  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* 点击动画效果 */
@keyframes clickEffect {
  0% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.05);
  }

  100% {
    transform: scale(1);
  }
}
</style>
