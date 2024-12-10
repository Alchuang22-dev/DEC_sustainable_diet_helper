<template>
  <view class="container">
	<view class="header">
	  <text class="title">{{$t('diet_restriction_label')}}{{$t('and')}}{{$t('preferences_title')}}</text>
	</view>
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view group="preferences" class="preferences">
      <view
        v-for="(preference, index) in preferences"
        :key="index"
        :class="['preference-card', 'color-' + ((index+1) % 4)]"
      >
        <image :src="preference.icon" class="preference-icon" />
        <text class="preference-name">{{ preference.name }}</text>
        <button class="delete-button" @click="removePreference(index)">
          <image src="@/pagesMy/static/delete.svg" class="delete-icon" />
        </button>
      </view>
    </view>
	<view group="preferences" class="preferences">
	  <view
	    v-for="(restriction, index) in dietRestrictions"
	    :key="index"
	    :class="['preference-card', 'color-' + (index % 4)]"
	  >
		<image src='https://cdn.pixabay.com/photo/2015/03/14/14/00/carrots-673184_1280.jpg' class="preference-icon" />
	    <text class="restriction-name">{{ restriction.name }}</text>
	    <button class="delete-button" @click="removeDietRestriction(index)">
	      <image src="@/pagesMy/static/delete.svg" class="delete-icon" />
	    </button>
	  </view>
	</view>
    <view class="add-preference">
      <button @click="showPreferenceOptions">{{$t('add_preference_button')}}</button>
    </view>
	<view class="add-restriction">
	  <uni-combox
	    :placeholder="$t('please_enter_food_name')"
	    v-model="foodNameInput"
	    :candidates="filteredFoods.map(item => displayName(item))"
	    @input="onComboxInput"
	  ></uni-combox>
	  <button @click="addDietRestriction">{{$t('add_restriction_button')}}</button>
	</view>
    <view v-if="showModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{$t('modal_title')}}</text>
        <view
          v-for="(option, index) in preferenceOptions"
          :key="index"
          class="modal-option"
          @click="selectPreference(option)"
        >
          <image :src="option.icon" class="option-icon" />
          <text class="option-name">{{ option.name }}</text>
        </view> 
      </view>
      <view class="button-content">
        <button class="close-button" @click="closeModal">
          {{$t('close_button')}}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref, reactive, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFoodListStore } from '@/stores/food_list'; // 引入 Pinia Store
import { useUserStore } from "@/stores/user.js"; // 引入用户 Store

const { t, locale } = useI18n();

const foodStore = useFoodListStore();
const userStore = useUserStore();

// 获取 availableFoods
const { availableFoods, fetchAvailableFoods, addFood } = foodStore;

// 定义 BASE_URL 为 ref
const BASE_URL = ref('http://122.51.231.155:8095');

// 食品数据
const food = reactive({
  name: '',
  id: null, // 添加 id 字段
  weight: '',
  price: '',
  transportMethod: 'land', // 使用英文标识
  foodSource: 'local', // 使用英文标识
  imagePath: '',
});

// 食物名称输入
const foodNameInput = ref('');
const showFoodList = ref(false);

// 模糊匹配过滤食物列表
const filteredFoods = computed(() => {
  if (foodNameInput.value === '') {
    const currentLang = locale.value;
    if (currentLang === 'zh-Hans') {
      return availableFoods.filter((f) => f.name_zh !== '');
    } else {
      return availableFoods.filter((f) => f.name_en !== '');
    }
  } else {
    const currentLang = locale.value;
    return availableFoods.filter((f) => {
      if (currentLang === 'zh-Hans') {
        return f.name_zh.includes(foodNameInput.value);
      } else {
        return f.name_en.toLowerCase().includes(foodNameInput.value.toLowerCase());
      }
    });
  }
});

// 根据当前语言显示名称
const displayName = (item) => {
  return locale.value === 'zh-Hans' ? item.name_zh : item.name_en;
};

const onComboxInput = (value) => {
  console.log('onComboxInput', value);
  foodNameInput.value = value;
};

// 当用户选择食物时
const selectFood = (foodItem) => {
  food.name = foodItem.name_en; // 始终使用英文名称存储
  food.id = foodItem.id;
  foodNameInput.value = displayName(foodItem); // 显示当前语言的名称
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

// 定义 token 为 computed 属性
const token = computed(() => userStore.user.token);

// 初始化 preferences 可以保留初始项
const preferences = ref([
  { name: t('foodpreference_greeting'), key: 'foodpreference_greeting', icon: 'https://via.placeholder.com/50' },
]);

const preferenceOptions = ref([
  { name: t('highProtein'), key: 'highProtein', icon: 'https://via.placeholder.com/50' },
  { name: t('highEnergy'), key: 'highEnergy', icon: 'https://via.placeholder.com/50' },
  { name: t('lowFat'), key: 'lowFat', icon: 'https://via.placeholder.com/50' },
  { name: t('lowCH'), key: 'lowCH', icon: 'https://via.placeholder.com/50' },
  { name: t('lowsodium'), key: 'lowsodium', icon: 'https://via.placeholder.com/50' },
  { name: t('vegan'), key: 'vegan', icon: 'https://via.placeholder.com/50' },
  { name: t('vegetarian'), key: 'vegetarian', icon: 'https://via.placeholder.com/50' },
  { name: t('glulenFree'), key: 'glulenFree', icon: 'https://via.placeholder.com/50' },
  { name: t('alcoholFree'), key: 'alcoholFree', icon: 'https://via.placeholder.com/50' },
  { name: t('dairyFree'), key: 'dairyFree', icon: 'https://via.placeholder.com/50' },
]);

const showModal = ref(false);

const newRestriction = ref(''); // 新增饮食禁忌的输入
const dietRestrictions = ref([]); // 存储用户的饮食禁忌

onMounted(() => {
  // 不再从 localStorage 中获取 token，而是使用 computed token
  if (!token.value) {
    console.warn('No token found in userStore.');
  }
  if (availableFoods.length === 0) {
    fetchAvailableFoods();
  }
   getDietRestriction();
  // 请求偏好数据
  uni.request({
    url: `${BASE_URL.value}/preferences`, // 使用 BASE_URL.value
    method: 'GET',
    header: {
      "Authorization": `Bearer ${token.value}`, // 使用 computed token
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {},
    success: (res) => {
      if (res.statusCode === 200) {
        // 处理返回的数据
        console.log('success to get preference');
        const data = res.data;  // 假设返回的数据格式是 [ ... ]

        // 将后端数据转移到 preferences 数组中
        data.forEach(item => {
          preferences.value.push({
            name: t(item.name), // 通过翻译返回名称
            key: item.name,     // 设置 key 为后端的 name 字段
            icon: 'https://via.placeholder.com/50' // 默认 icon 地址
          });
        });
      } else {
        console.error('Failed to load preferences:', res.data);
      }
    },
    fail: (err) => {
      console.error('Error fetching preferences:', err);
    }
  });
});

const removePreference = (index) => {
  const preferenceToRemove = preferences.value[index];
  console.log(preferenceToRemove.key);
  uni.request({
    url: `${BASE_URL.value}/preferences`, // 使用 BASE_URL.value
    method: 'DELETE',
    header: {
      "Authorization": `Bearer ${token.value}`, // 使用 computed token
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {
      preference_name: preferenceToRemove.key // 使用存储的 key 字段
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message); // 打印成功信息
        // 删除本地数组中的偏好
        preferences.value.splice(index, 1);
      }
    },
    fail: (err) => {
      console.error('Error removing preference:', err);
    }
  });
};

const showPreferenceOptions = () => {
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const selectPreference = (option) => {
  console.log(option.key);
  uni.request({
    url: `${BASE_URL.value}/preferences`, // 使用 BASE_URL.value
    method: 'POST',
    header: {
      "Authorization": `Bearer ${token.value}`, // 使用 computed token
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {
      preference_name: option.key // 使用存储的 key 字段
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message); // 打印成功信息
        // 将新的偏好添加到本地
        preferences.value.push({ name: option.name, key: option.key, icon: option.icon });
        closeModal();
      }
      console.log(res.data);
    },
    fail: (err) => {
      console.log(err); // 修复：不应使用 res.data
      console.error('Error adding preference:', err);
    }
  });
};

// 添加饮食禁忌的方法
const addDietRestriction = () => {
  console.log('输入食品禁忌');
  const matchedFood = availableFoods.find((f) => displayName(f) === foodNameInput.value);
  if (matchedFood) {
    // 如果找到匹配的食物项，使用 selectFood 方法
    selectFood(matchedFood);
  } else {
    // 如果没有找到，提醒用户
    uni.showToast({
      title: t('no_matching_food'),
      icon: 'none',
      duration: 2000,
    });
    return; // 终止提交
  }
  if (foodNameInput.value.trim()) {
    // 先向后端发送请求
	console.log(food.id);
    uni.request({
      url: `${BASE_URL.value}/disliked_preferences`, // 使用 BASE_URL.value
      method: 'POST',
      header: {
        "Authorization": `Bearer ${token.value}`, // 使用 computed token
        "Content-Type": "application/json", // 设置请求类型
      },
      data: {
        food_id: food.id, // 使用食物的 id 来提交
      },
      success: (res) => {
        if (res.statusCode === 200) {
          console.log(res.data.message); // 打印成功信息
          dietRestrictions.value.push({ name: foodNameInput.value.trim(), id: food.id }); // 将禁忌添加到本地
          foodNameInput.value = ''; // 清空输入框
        }
      },
      fail: (err) => {
        console.error('Error adding diet restriction:', err);
      }
    });
  } else {
    console.warn('Please enter a valid diet restriction');
  }
};


// 删除饮食禁忌的方法
const removeDietRestriction = (index) => {
  const restrictionToRemove = dietRestrictions.value[index];
  console.log('Deleting diet restriction for food ID:', restrictionToRemove.id);

  // 向后端发送删除请求
  uni.request({
    url: `${BASE_URL.value}/disliked_preferences`, // 使用 BASE_URL.value
    method: 'DELETE',
    header: {
      "Authorization": `Bearer ${token.value}`, // 使用 computed token
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {
      food_id: restrictionToRemove.id, // 使用食物的 id 来删除
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message); // 打印成功信息
        // 删除本地数组中的禁忌
        dietRestrictions.value.splice(index, 1);
      }
    },
    fail: (err) => {
      console.error('Error removing diet restriction:', err);
    }
  });
};

const getDietRestriction = () => {
  uni.request({
    url: `${BASE_URL.value}/disliked_preferences`, // 使用 BASE_URL.value
    method: 'GET',
    header: {
      "Authorization": `Bearer ${token.value}`, // 使用 computed token
      "Content-Type": "application/json", // 设置请求类型
    },
    data: {},
    success: (res) => {
      if (res.statusCode === 200) {
        console.log('Fetched disliked foods:', res.data.disliked_foods);
        // 将返回的禁忌食物名称列表更新到本地数据中，并包含 id
        dietRestrictions.value = res.data.disliked_foods.map(food => ({
          name: food.name,
          id: food.id
        }));
      } else {
        console.error('Failed to load diet restrictions:', res.data);
      }
    },
    fail: (err) => {
      console.error('Error fetching diet restrictions:', err);
    }
  });
};

</script>

<style scoped>
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: -1;
  opacity: 0.1;
}	

.container {
  padding: 20px;
}
.header {
  text-align: center;
  margin-bottom: 20px;
}
.title {
  font-size: 24px;
  font-weight: bold;
}
.preferences {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.preference-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
}
.color-0 {
  background-color: #4ca;
}
.color-1 {
  background-color: #e0f7fa;
}
.color-2 {
  background-color: #ffe0b2;
}
.color-3 {
  background-color: #e1bee7;
}
.preference-icon {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}
.preference-name {
  flex: 1;
  font-size: 16px;
}
.delete-button {
  background: none;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.delete-icon {
  width: 20px;
  height: 20px;
}
.add-preference {
  margin-top: 20px;
  text-align: center;
}
.add-preference button {
  background-color: #4caf50;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column; /* 设置垂直排列 */
  z-index: 2;
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 80%;
  max-width: 400px;
  max-height: 400px;
  overflow-y: auto; /* 允许垂直滚动 */
  margin-bottom: 20px; /* 添加底部间距 */
}

.modal-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 20px;
}

.modal-option {
  display: flex;
  align-items: center;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
}

.option-icon {
  width: 40px;
  height: 40px;
  margin-right: 10px;
}

.option-name {
  font-size: 16px;
}

.button-content {
  display: flex;
  justify-content: center; /* 按钮居中 */
  width: 100%; /* 确保按钮容器宽度充满 */
}

.close-button {
  background-color: #ff4d4f;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 20px;
}

.add-restriction {
  margin-top: 20px;
  text-align: center;
}

.restriction-label {
  font-size: 18px;
  margin-bottom: 10px;
}

.restriction-input {
  width: 80%;
  padding: 8px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: 1px solid #ccc;
}

.restriction-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  background-color: #f8f8f8;
  margin-bottom: 10px;
}
</style>
