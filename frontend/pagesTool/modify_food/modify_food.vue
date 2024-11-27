<!-- pagesTool/modify_food/modify_food.vue -->
<template>
  <view class="container">
    <!-- 表单容器 -->
    <view class="form-container">
      <view class="form-group">
        <text class="label">{{ $t('name') }}</text>
        <input
          class="input"
          type="text"
          v-model="foodNameInput"
          @focus="showFoodList = true"
          @blur="onInputBlur"
          :placeholder="$t('please_enter_food_name')"
        />
        <!-- 食物名的模糊匹配下拉列表 -->
        <view v-if="showFoodList && filteredFoods.length > 0" class="food-list">
          <view
            v-for="item in filteredFoods"
            :key="item.id"
            class="food-item"
            @mousedown.prevent
            @click="selectFood(item)"
          >
            {{ item.name }}
          </view>
        </view>
      </view>
      <view class="form-group">
        <text class="label">{{ $t('total_weight') }}</text>
        <input class="input" type="number" v-model="food.weight" :placeholder="$t('please_enter_food_weight')" :error="weightError" />
        <text v-if="weightError" class="error-message">{{ $t('weight_must_be_positive_integer') }}</text>
      </view>
      <view class="form-group">
        <text class="label">{{ $t('total_price') }}</text>
        <input class="input" type="number" v-model="food.price" :placeholder="$t('please_enter_food_price')" :error="priceError" />
        <text v-if="priceError" class="error-message">{{ $t('price_must_be_positive_integer') }}</text>
      </view>
      <view class="form-group">
        <text class="label">{{ $t('select_transport_method') }}</text>
        <picker mode="selector" :range="transportMethods" :value="transportIndex" @change="onTransportChange">
          <view class="picker">
            {{ transportMethods[transportIndex] }}
          </view>
        </picker>
      </view>
      <view class="form-group">
        <text class="label">{{ $t('select_food_source') }}</text>
        <picker mode="selector" :range="foodSources" :value="sourceIndex" @change="onSourceChange">
          <view class="picker">
            {{ foodSources[sourceIndex] }}
          </view>
        </picker>
      </view>

      <!-- 图片上传按钮 -->
      <view class="form-group">
        <text class="label">{{ $t('upload_food_image') }}</text>
        <button class="upload-button" @click="uploadImage">{{ $t('take_photo_upload') }}</button>
        <image v-if="food.imagePath" :src="food.imagePath" class="uploaded-image"></image>
      </view>

      <button class="submit-button" @click="submitFoodDetails">{{ $t('submit') }}</button>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n'; // Import useI18n
import { useFoodListStore } from '@/stores/food_list'; // 引入 Pinia Store
import { onLoad } from '@dcloudio/uni-app'; // 引入 onLoad 钩子

const { t } = useI18n();
const foodStore = useFoodListStore();

// 获取 availableFoods
const { availableFoods, fetchAvailableFoods, updateFood } = foodStore;

// 初始化数据
const options = ref({});
const foodIndex = ref(null);
const existingFood = ref(null);

// 食品数据，初始化为空
const food = reactive({
  name: '',
  id: null, // 添加 id 字段
  weight: '',
  price: '',
  transportMethod: 'transport_land',
  foodSource: 'source_local',
  imagePath: '',
});

// 食物名称输入
const foodNameInput = ref('');
const showFoodList = ref(false);

// 模糊匹配过滤食物列表
const filteredFoods = computed(() => {
  if (foodNameInput.value === '') {
    return availableFoods;
  } else {
    return availableFoods.filter((f) => f.name.includes(foodNameInput.value));
  }
});

// 当用户选择食物时
const selectFood = (foodItem) => {
  food.name = foodItem.name;
  food.id = foodItem.id;
  foodNameInput.value = foodItem.name;
  showFoodList.value = false;
};

// 处理输入框失焦事件
const onInputBlur = () => {
  setTimeout(() => {
    showFoodList.value = false;
  }, 100);
};

// 监听输入变化，控制下拉列表显示
watch(foodNameInput, (newValue) => {
  if (newValue !== '') {
    showFoodList.value = true;
  } else {
    showFoodList.value = false;
  }
});

// 输入验证错误状态
const weightError = ref(false);
const priceError = ref(false);

// 运输方式和食品来源下拉选项数据
const transportMethods = [t('transport_land'), t('transport_sea'), t('transport_air')];
const foodSources = [t('source_local'), t('source_imported')];

// 当前选择的索引
const transportIndex = ref(0);
const sourceIndex = ref(0);

// 运输方式选择改变
const onTransportChange = (e) => {
  transportIndex.value = e.detail.value;
  food.transportMethod = transportMethods[transportIndex.value];
};

// 食品来源选择改变
const onSourceChange = (e) => {
  sourceIndex.value = e.detail.value;
  food.foodSource = foodSources[sourceIndex.value];
};

// 上传图片
const uploadImage = () => {
  uni.chooseImage({
    count: 1, // 只选择一张图片
    sizeType: ['original', 'compressed'], // 可以选择原图或压缩图
    sourceType: ['camera'], // 只允许使用相机
    success: (res) => {
      const tempFilePath = res.tempFilePaths[0];
      food.imagePath = tempFilePath;

      // TODO: 集成图像识别功能
      // 您可以在这里调用图像识别 API，将 tempFilePath 发送到服务器进行识别
      // 例如：
      // recognizeFoodImage(tempFilePath).then(recognizedName => {
      //   food.name = recognizedName;
      // });
    },
    fail: (err) => {
      uni.showToast({
        title: t('image_upload_failed'),
        icon: 'none',
        duration: 2000,
      });
      console.error('图片上传失败:', err);
    },
  });
};

// 提交表单
const submitFoodDetails = () => {
  // 重置错误状态
  weightError.value = false;
  priceError.value = false;

  const {
    name,
    weight,
    price,
    transportMethod,
    foodSource,
    imagePath
  } = food;

  // 输入验证
  let valid = true;

  // 验证重量：必须是正整数
  if (!/^\d+$/.test(weight) || parseInt(weight) <= 0) {
    weightError.value = true;
    valid = false;
  }

  // 验证价格：必须是正整数
  if (!/^\d+$/.test(price) || parseInt(price) <= 0) {
    priceError.value = true;
    valid = false;
  }

  if (!name || !weight || !price || !transportMethod || !foodSource) {
    uni.showToast({
      title: t('please_fill_all_fields'),
      icon: 'none',
    });
    valid = false;
  }

  if (!valid) {
    return;
  }

  const updatedFood = {
    name,
    id: food.id, // 添加 id
    weight: parseInt(weight),
    price: parseInt(price),
    transportMethod,
    foodSource,
    image: imagePath,
    // 保留其他字段不变
  };

  // 使用 Store 更新指定的食物项
  updateFood(foodIndex.value, updatedFood);

  uni.showToast({
    title: t('modify_success'),
    icon: 'success',
    duration: 2000,
  });

  // 返回上一页并刷新
  setTimeout(() => {
    uni.navigateBack();
  }, 2000);
};

// 使用 onLoad 获取路由参数并初始化数据
onLoad((loadedOptions) => {
  options.value = loadedOptions;
  console.log('路由参数:', options.value);

  // 获取传入的食物索引
  foodIndex.value = parseInt(options.value.index);
  console.log('食物索引:', foodIndex.value);
  existingFood.value = foodStore.foodList[foodIndex.value];
  console.log('现有食物:', existingFood.value);

  // 错误处理：如果索引无效或食物项不存在，显示提示并返回
  if (isNaN(foodIndex.value) || !existingFood.value) {
    uni.showToast({
      title: t('invalid_food_item'),
      icon: 'none',
      duration: 2000,
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 2000);
    return; // 记得加上 return，避免继续执行后续代码
  }

  // 初始化食品数据
  food.name = existingFood.value.name || '';
  food.id = existingFood.value.id || null;
  food.weight = existingFood.value.weight || '';
  food.price = existingFood.value.price || '';
  food.transportMethod = existingFood.value.transportMethod || 'transport_land';
  food.foodSource = existingFood.value.foodSource || 'source_local';
  food.imagePath = existingFood.value.image || '';

  // 初始化输入框
  foodNameInput.value = food.name;

  // 初始化下拉选项索引
  transportIndex.value = transportMethods.indexOf(existingFood.value.transportMethod);
  sourceIndex.value = foodSources.indexOf(existingFood.value.foodSource);
});

// 页面加载时执行
onMounted(() => {
  // 如果 availableFoods 为空，调用获取函数
  if (availableFoods.length === 0) {
    fetchAvailableFoods();
  }
});
</script>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-color: #4caf50;
  --secondary-color: #8bc34a;
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
}

/* 头部标题 */
.header {
  display: flex;
  align-items: center;
  padding: 20rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid var(--border-color);
  justify-content: flex-start;
}

.back-button {
  font-size: 36rpx;
  margin-right: 20rpx;
  color: var(--primary-color);
  cursor: pointer;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: var(--text-color);
}

/* 表单容器 */
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
  width: 100%;
  padding: 20rpx;
  border: 1rpx solid var(--border-color);
  border-radius: 10rpx;
  font-size: 28rpx;
}

.picker {
  width: 100%;
  padding: 20rpx;
  border: 1rpx solid var(--border-color);
  border-radius: 10rpx;
  font-size: 28rpx;
  color: #666666;
}

.upload-button {
  width: 100%;
  padding: 20rpx;
  border: 1rpx solid var(--border-color);
  border-radius: 10rpx;
  background-color: #f0f0f0;
  font-size: 28rpx;
  color: var(--text-color);
  cursor: pointer;
  text-align: center;
}

.upload-button:hover {
  background-color: #e0e0e0;
}

.uploaded-image {
  width: 100%;
  height: auto;
  margin-top: 20rpx;
  border-radius: 10rpx;
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

/* 添加下拉列表样式 */
.food-list {
  max-height: 300rpx;
  overflow-y: auto;
  border: 1rpx solid var(--border-color);
  border-top: none;
  background-color: #ffffff;
}

.food-item {
  padding: 20rpx;
  font-size: 28rpx;
  color: var(--text-color);
  cursor: pointer;
}

.food-item:hover {
  background-color: #f0f0f0;
}
</style>
