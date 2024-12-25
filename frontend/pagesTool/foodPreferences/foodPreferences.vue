<template>
  <view class="container">
    <!-- 标题 -->
    <view class="header">
      <text class="title">
        {{$t('diet_restriction_label')}}{{$t('and')}}{{$t('preferences_title')}}
      </text>
    </view>
    <!-- 替换 greeting 卡片为文本 -->
    <view class="greeting-wrap">
      <text class="greeting-text">{{ t('foodpreference_greeting') }}</text>
    </view>

    <!-- 背景图 -->
    <image
      src="/static/images/index/background_img.jpg"
      class="background-image"
    ></image>

    <!-- 使用 uni-card 容器展示 食物偏好 -->
    <uni-card
      title="偏好"
      :is-shadow="true"
      class="preference-card-wrapper"
    >
      <view
        v-for="(preference, index) in preferences"
        :key="index"
        :class="['preference-card', 'color-' + ((index + 1) % 4)]"
      >
        <image :src="preference.icon" class="preference-icon" />
        <text class="preference-name">{{ preference.name }}</text>
        <button class="delete-button" @click="removePreference(index)">
          <image src="@/pagesMy/static/delete.svg" class="delete-icon" />
        </button>
      </view>

      <!-- 添加偏好按钮 -->
      <view class="add-preference">
        <button @click="showPreferenceOptions">{{$t('add_preference_button')}}</button>
      </view>
    </uni-card>

    <!-- 黑名单输入区域 -->
    <view class="add-restriction">
      <uni-combox
        :placeholder="$t('please_enter_food_name')"
        v-model="foodNameInput"
        :candidates="filteredFoods.map(item => displayName(item))"
        @input="onComboxInput"
      ></uni-combox>
      <button @click="addDietRestriction">{{$t('add_restriction_button')}}</button>
    </view>

    <!-- 黑名单列表（与偏好分开） -->
    <uni-card
      title="黑名单"
      :is-shadow="true"
      class="blacklist-card"
    >
      <!-- 若无黑名单，简单提示 -->
      <view v-if="dietRestrictions.length === 0" style="padding: 10px;">
        <text>暂无黑名单食物</text>
      </view>
      <view
        v-for="(restriction, index) in dietRestrictions"
        :key="index"
        :class="['preference-card', 'color-' + (index % 4)]"
      >
        <image
          src="https://cdn.pixabay.com/photo/2015/03/14/14/00/carrots-673184_1280.jpg"
          class="preference-icon"
        />
        <text class="restriction-name">{{ restriction.name }}</text>
        <!-- 删除按钮与偏好卡片一致 -->
        <button class="delete-button" @click="removeDietRestriction(index)">
          <image src="@/pagesMy/static/delete.svg" class="delete-icon" />
        </button>
      </view>
    </uni-card>

    <!-- 选择偏好弹窗 -->
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
        <button class="close-button" @click="closeModal">{{$t('close_button')}}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref, reactive, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFoodListStore } from '../stores/food_list.js';
import { useUserStore } from '@/stores/user.js';

// 引入 i18n
const { t, locale } = useI18n();

// 引入 store
const foodStore = useFoodListStore();
const userStore = useUserStore();

// 获取 store 中的方法
const { availableFoods, fetchAvailableFoods } = foodStore;
const BASE_URL = ref('http://122.51.231.155:8095');

// 定义食物信息
const food = reactive({
  name: '',
  id: null,
  weight: '',
  price: '',
  transportMethod: 'land',
  foodSource: 'local',
  imagePath: '',
});

// 用户输入
const foodNameInput = ref('');
const showFoodList = ref(false);

// 根据语言过滤
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

// 显示当前语言
const displayName = (item) => {
  return locale.value === 'zh-Hans' ? item.name_zh : item.name_en;
};

const onComboxInput = (value) => {
  foodNameInput.value = value;
};

// 当用户选择食物时
const selectFood = (foodItem) => {
  food.name = foodItem.name_en; // 使用英文存储
  food.id = foodItem.id;
  foodNameInput.value = displayName(foodItem); // 显示当前语言
  showFoodList.value = false;
};

// 监听输入变化控制下拉
watch(foodNameInput, (newValue) => {
  showFoodList.value = newValue !== '';
});

// 计算属性：token
const token = computed(() => userStore.user.token);

// 不再包含 greeting
const preferences = ref([]);

// 供选择的偏好列表
const preferenceOptions = ref([
  { name: t('highProtein'), key: 'highProtein', icon: 'https://cdn.pixabay.com/photo/2023/09/22/07/23/ai-generated-8268310_1280.jpg' },
  { name: t('highEnergy'), key: 'highEnergy', icon: 'https://cdn.pixabay.com/photo/2019/06/01/05/45/dumplings-4243484_1280.jpg' },
  { name: t('lowFat'), key: 'lowFat', icon: 'https://cdn.pixabay.com/photo/2023/06/09/18/18/keto-8052361_1280.png' },
  { name: t('lowCH'), key: 'lowCH', icon: 'https://cdn.pixabay.com/photo/2018/09/23/09/31/smoothie-3697014_1280.jpg' },
  { name: t('lowsodium'), key: 'lowsodium', icon: 'https://cdn.pixabay.com/photo/2016/11/19/09/42/berries-1838314_1280.jpg' },
  { name: t('vegan'), key: 'vegan', icon: 'https://cdn.pixabay.com/photo/2019/04/13/19/03/cow-4125323_1280.png' },
  { name: t('vegetarian'), key: 'vegetarian', icon: 'https://cdn.pixabay.com/photo/2016/09/22/18/51/heart-1688029_1280.png' },
  { name: t('glulenFree'), key: 'glulenFree', icon: 'https://cdn.pixabay.com/photo/2011/08/17/12/31/spike-8743_1280.jpg' },
  { name: t('alcoholFree'), key: 'alcoholFree', icon: 'https://cloud.tsinghua.edu.cn/thumbnail/cf9dba3a498247469fd4/1024/alcohol_free.png' },
  { name: t('dairyFree'), key: 'dairyFree', icon: 'https://cdn.pixabay.com/photo/2022/04/04/14/17/milk-7111433_1280.jpg' },
]);

// 弹窗控制
const showModal = ref(false);

// 黑名单数组
const dietRestrictions = ref([]);

// 页面挂载时初始化
onMounted(() => {
  if (!token.value) {
    console.warn('No token found in userStore.');
  }
  if (availableFoods.length === 0) {
    fetchAvailableFoods();
  }
  getDietRestriction();
  // 请求偏好数据
  uni.request({
    url: `${BASE_URL.value}/preferences`,
    method: 'GET',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json',
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log('success to get preference');
        const data = res.data; // 假设返回数组
        // 将后端数据添加到 preferences
        data.forEach((item) => {
          preferences.value.push({
            name: t(item.name),
            key: item.name,
            icon: 'https://via.placeholder.com/50',
          });
        });
      } else {
        console.error('Failed to load preferences:', res.data);
      }
    },
    fail: (err) => {
      console.error('Error fetching preferences:', err);
    },
  });
});

// 删除偏好
const removePreference = (index) => {
  const preferenceToRemove = preferences.value[index];
  uni.request({
    url: `${BASE_URL.value}/preferences`,
    method: 'DELETE',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json',
    },
    data: {
      preference_name: preferenceToRemove.key,
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message);
        preferences.value.splice(index, 1);
      }
    },
    fail: (err) => {
      console.error('Error removing preference:', err);
    },
  });
};

// 添加偏好
const showPreferenceOptions = () => {
  showModal.value = true;
};
const closeModal = () => {
  showModal.value = false;
};
const selectPreference = (option) => {
  uni.request({
    url: `${BASE_URL.value}/preferences`,
    method: 'POST',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json',
    },
    data: {
      preference_name: option.key,
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message);
        preferences.value.push({ name: option.name, key: option.key, icon: option.icon });
        closeModal();
      }
    },
    fail: (err) => {
      console.error('Error adding preference:', err);
    },
  });
};

// 黑名单相关
const addDietRestriction = () => {
  const matchedFood = availableFoods.find((f) => displayName(f) === foodNameInput.value);
  if (matchedFood) {
    selectFood(matchedFood);
  } else {
    uni.showToast({
      title: t('no_matching_food'),
      icon: 'none',
      duration: 2000,
    });
    return;
  }
  if (foodNameInput.value.trim()) {
    uni.request({
      url: `${BASE_URL.value}/disliked_preferences`,
      method: 'POST',
      header: {
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'application/json',
      },
      data: {
        food_id: food.id,
      },
      success: (res) => {
        if (res.statusCode === 200) {
          console.log(res.data.message);
          dietRestrictions.value.push({
            name: foodNameInput.value.trim(),
            id: food.id,
          });
          foodNameInput.value = '';
        }
      },
      fail: (err) => {
        console.error('Error adding diet restriction:', err);
      },
    });
  } else {
    console.warn('Please enter a valid diet restriction');
  }
};

const removeDietRestriction = (index) => {
  const restrictionToRemove = dietRestrictions.value[index];
  uni.request({
    url: `${BASE_URL.value}/disliked_preferences`,
    method: 'DELETE',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json',
    },
    data: {
      food_id: restrictionToRemove.id,
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log(res.data.message);
        dietRestrictions.value.splice(index, 1);
      }
    },
    fail: (err) => {
      console.error('Error removing diet restriction:', err);
    },
  });
};

const getDietRestriction = () => {
  uni.request({
    url: `${BASE_URL.value}/disliked_preferences`,
    method: 'GET',
    header: {
      Authorization: `Bearer ${token.value}`,
      'Content-Type': 'application/json',
    },
    success: (res) => {
      if (res.statusCode === 200) {
        console.log('Fetched disliked foods:', res.data.disliked_foods);
        dietRestrictions.value = res.data.disliked_foods.map((food) => ({
          name: food.name,
          id: food.id,
        }));
      } else {
        console.error('Failed to load diet restrictions:', res.data);
      }
    },
    fail: (err) => {
      console.error('Error fetching diet restrictions:', err);
    },
  });
};
</script>

<style scoped>
/* 这里仅作示例演示，实际项目中可根据情况适配 uni-ui 的暗色模式、行间距等 */
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
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
  padding: 5px;
}

.header {
  text-align: center;
  margin-bottom: 10px;
}
.title {
  font-size: 24px;
  font-weight: bold;
}

.greeting-wrap {
  margin-bottom: 10px;
  text-align: center;
}
.greeting-text {
  font-size: 16px;
  color: #333;
}

/* uni-card 包裹的偏好/黑名单列表 */
.preference-card-wrapper,
.blacklist-card {
  margin-top: 10px;
}

/* 偏好/黑名单的卡片通用样式 */
.preference-card {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  /* 在使用 uni-card 时，可以去掉显式的 background-color，
     或保留给卡片内容分区用 */
  margin-bottom: 10px;
}
.color-0 {
  background-color: #4ca;         /* 可自行调整颜色 */
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

.preference-name,
.restriction-name {
  flex: 1;
  font-size: 16px;
}

/* 删除按钮统一样式 */
.delete-button {
  background: none;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  padding: 5px;
}
.delete-icon {
  width: 20px;
  height: 20px;
}

/* 添加偏好按钮 */
.add-preference {
  margin-top: 10px;
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

/* 黑名单输入区域 */
.add-restriction {
  margin-top: 20px;
  margin-left: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

/* 弹窗相关样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 99;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 80%;
  max-width: 400px;
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 20px;
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
  justify-content: center;
  width: 100%;
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
</style>