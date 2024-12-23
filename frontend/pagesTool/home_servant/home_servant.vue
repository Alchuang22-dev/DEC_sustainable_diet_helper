<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_index_new.png" class="background-image"></image>

    <!-- 头部标题 -->
    <view v-if="family.status === FamilyStatus.JOINED" class="header">
      <text class="header-title">{{ family.name || $t('default_family_name') }}</text>
    </view>

    <!-- 家庭ID -->
    <text v-if="family.status === FamilyStatus.JOINED" class="list-title">{{ $t('family_info') + family.familyId }}</text>

    <!-- 用户未加入家庭时的视图 -->
    <view v-if="family.status === FamilyStatus.NOT_JOINED" class="no-family-view">
      <image src="https://cdn.pixabay.com/photo/2017/01/13/02/31/family-1976162_1280.png" class="centered-image"></image>
      <view class="family-actions">
        <button class="action-button" @click="showCreateFamilyModal = true">{{ t('create_family') }}</button>
        <button class="action-button" @click="showJoinFamilyModal = true">{{ t('join_family') }}</button>
      </view>
    </view>

    <!-- 待审核状态的视图 -->
    <view v-if="family.status === FamilyStatus.PENDING_APPROVAL" class="pending-view">
      <image src="https://cdn.pixabay.com/photo/2021/09/20/22/15/hourglass-6641967_1280.png" class="pending-image"></image>
      <text class="pending-text">{{ $t('pending_approval') }}</text>
      <button class="cancel-button" @click="handleCancelJoin">{{ $t('cancel_join') }}</button>
    </view>

    <!-- 如果已加入家庭并且是管理员且有waiting_members，则显示提示框 -->
    <uni-card
        v-if="isCurrentUserAdmin && family.waiting_members && family.waiting_members.length > 0"
        class="admin-notice"
        :is-shadow="false"
        :border="false"
    >
      <view class="notice-content">
        <view class="notice-left">
          <uni-icons type="info-filled" size="20" color="#3A86FF"></uni-icons>
          <text class="notice-text">{{ $t('new_join_requests') }}</text>
        </view>
        <text class="notice-btn" @click="goToMyFamily">
          {{ $t('handle_requests') }}
        </text>
      </view>
    </uni-card>

    <!-- 用户已加入家庭时的家庭管理部分 -->
    <view v-if="family.status === FamilyStatus.JOINED" class="family-management">
      <!-- 提出想吃的菜品 -->
      <view class="dish-proposal">
        <text class="section-title">{{ t('propose_dish') }}</text>

        <!-- 使用uni-combox替换原有的input输入框 -->
        <uni-combox
            :placeholder="t('dish_name_placeholder')"
            v-model="foodNameInput"
            :candidates="filteredFoods.map(item => displayName(item))"
            @input="onComboxInput"
        >
        </uni-combox>

        <!-- 偏好选择器保持不变 -->
        <picker
            mode="selector"
            :range="dishPreferenceLevels"
            :value="newDish.preference"
            @change="onDishPreferenceChange"
            class="picker"
        >
          <view>{{ t('dish_preference') }}: {{ dishPreferenceLevels[newDish.preference] }}</view>
        </picker>
        <button class="submit-button" @click="submitDishProposal">{{ t('submit_proposal') }}</button>
      </view>

      <!-- 家庭成员的提议 -->
      <view class="dish-list">
        <text class="section-title">{{ $t('family_dish_proposals') }}</text>
        <scroll-view scroll-y class="proposals-scroll">
          <uni-collapse v-model="activeCollapse" :accordion="true">
            <uni-collapse-item
                v-for="dish in sortedDishProposals"
                :key="dish.id"
                :name="dish.id.toString()"
            >
              <!-- 自定义折叠面板标题 -->
              <template #title>
                <view class="collapse-title">
                  <text class="dish-name">{{ dish.name }}</text>
                  <view class="dish-info">
                    <uni-tag
                        :text="dishPreferenceLevels[dish.preference]"
                        :type="getPreferenceTagType(dish.preference)"
                        size="small"
                    />
                    <text class="proposer-text">{{ dish.proposer }}</text>
                  </view>
                </view>
              </template>
              <!-- 折叠面板内容 -->
              <view class="dish-actions">
                <uni-icons
                    v-if="canDeleteDish(dish)"
                    type="trash"
                    size="24"
                    color="#FF4D4F"
                    @click.stop="confirmDeleteDish(dish.id)"
                    class="delete-icon"
                ></uni-icons>
              </view>
            </uni-collapse-item>
          </uni-collapse>
        </scroll-view>
      </view>

      <!-- 共享家庭成员的五大营养成分达标情况 -->
      <view class="shared-data">
        <text class="section-title">{{ $t('shared_family_data') }}</text>
        <!-- 家庭碳排放环形图 -->
        <view class="charts">
          <qiun-data-charts
              :canvas2d="true"
              canvas-id="familyCarbonChart"
              type="ring"
              :opts="carbonRingOpts"
              :chartData="carbonChartData"
          />
        </view>
      </view>

      <!-- 新增：家庭成员营养超标情况 -->
      <view class="nutrition-over-section">
        <text class="section-title">{{ t('nutrition_over_title') }}</text>
        <view
            v-for="member in membersNutritionOver"
            :key="member.userId"
            class="nutrition-over-row"
        >
          <image
              :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`"
              class="member-avatar"
          />
          <text class="member-name">{{ member.nickname }}</text>

          <!-- 若有超标营养素 -->
          <text class="over-list" v-if="member.overs.length > 0">
            {{ member.overs.join('、') + t('nutrition_over_suffix') }}
          </text>
          <!-- 若无超标 -->
          <text class="over-list" v-else>
            {{ t('no_nutrition_over') }}
          </text>
        </view>
      </view>

      <!-- 家庭成员部分 -->
      <view class="family-info">
        <text class="section-title">{{ $t('family_members') }}</text>
        <view class="family-members">
          <view v-for="member in sortedFamilyMembers" :key="member.id" class="member">
            <image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" class="member-avatar"></image>
            <text class="member-name">{{ `${member.nickname} (${member.role === 'admin' ? $t('admin') : $t('member')})` }}</text>
          </view>
        </view>
        <!-- 仅管理员可见的“管理成员”按钮 -->
        <button v-if="isCurrentUserAdmin" class="manage-members-button" @click="manageMembers">{{ t('manage_members') }}</button>
      </view>

      <!-- 底部操作按钮部分 -->
      <view class="bottom-container">
        <uni-section class="action-section" type="line">
          <view class="bottom-actions">
            <text class="action-btn leave-btn" @click="handleLeaveFamily">
              <uni-icons type="close" size="16" color="#fff"></uni-icons>
              {{ $t('leave_family') }}
            </text>
            <text v-if="isCurrentUserAdmin" class="action-btn break-btn" @click="handleBreakFamily">
              <uni-icons type="trash" size="16" color="#fff"></uni-icons>
              {{ $t('break_family') }}
            </text>
          </view>
        </uni-section>
      </view>
    </view>

    <!-- 创建家庭的模态框 -->
    <view v-if="showCreateFamilyModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{ $t('create_family') }}</text>
        <input v-model="newFamilyName" :placeholder="t('family_name_placeholder')" class="input"></input>
        <button class="modal-button" @click="createFamily">{{ $t('confirm') }}</button>
        <button class="modal-button cancel" @click="showCreateFamilyModal = false">{{ $t('cancel') }}</button>
      </view>
    </view>

    <!-- 加入家庭的模态框 -->
    <view v-if="showJoinFamilyModal" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{ $t('join_family') }}</text>
        <input v-model="joinFamilyId" :placeholder="t('family_id_placeholder')" class="input"></input>
        <button class="modal-button" @click="joinFamily">{{ $t('confirm') }}</button>
        <button class="modal-button cancel" @click="showJoinFamilyModal = false">{{ $t('cancel') }}</button>
      </view>
    </view>

    <!-- 删除确认模态框 -->
    <view v-if="showDeleteConfirm" class="modal">
      <view class="modal-content">
        <text class="modal-title">{{ $t('confirm_delete') }}</text>
        <button class="modal-button" @click="deleteDish(selectedDishId)">{{ $t('confirm') }}</button>
        <button class="modal-button cancel" @click="showDeleteConfirm = false">{{ $t('cancel') }}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFamilyStore, FamilyStatus } from '../../stores/family.js';
import { useUserStore } from "../../stores/user.js";
import { onShow } from '@dcloudio/uni-app';
import { useFoodListStore } from '../../stores/food_list'; // 引入食物列表store

const { t, locale } = useI18n();

const familyStore = useFamilyStore();
const family = computed(() => familyStore.family);
const FamilyStatusEnum = FamilyStatus;
const userStore = useUserStore();

// 确保用户已登录，否则跳转到登录页面
if (!userStore.user.isLoggedIn) {
  uni.navigateTo({ url: '/pagesMy/login/login' });
}

const currentUserId = computed(() => userStore.user.uid);
const isCurrentUserAdmin = computed(() => {
  return familyStore.isAdmin(currentUserId.value);
});

const allFamilyMembers = computed(() => family.value.allMembers || []);
const sortedFamilyMembers = computed(() => {
  const members = [...allFamilyMembers.value];
  const currentUserIndex = members.findIndex(member => member.id === currentUserId.value);
  if (currentUserIndex > -1) {
    const [currentUser] = members.splice(currentUserIndex, 1);
    members.unshift(currentUser);
  }
  return members;
});

const newDish = ref({
  name: '',
  id: null,
  preference: 0,
});

// 菜品偏好级别
const dishPreferenceLevels = computed(() => [
  t('preference_low'),
  t('preference_medium'),
  t('preference_high'),
]);

// 根据偏好级别返回对应的标签类型
const getPreferenceTagType = (preference) => {
  switch (preference) {
    case 0:
      return 'info';    // 低偏好
    case 1:
      return 'warning'; // 中偏好
    case 2:
      return 'success'; // 高偏好
    default:
      return 'default';
  }
};

// 控制uni-collapse的活动项
const activeCollapse = ref([]);

// 模态框显示状态
const showCreateFamilyModal = ref(false);
const showJoinFamilyModal = ref(false);
const showDeleteConfirm = ref(false);

// 创建家庭相关
const newFamilyName = ref('');

// 加入家庭相关
const joinFamilyId = ref('');

// 引入食物列表Store和相关方法
const foodStore = useFoodListStore();
const { availableFoods, fetchAvailableFoods, getFoodName } = foodStore;

// 用于uni-combox的输入框和下拉列表逻辑
const foodNameInput = ref('');
const showFoodList = ref(false);

// 用于删除确认
const selectedDishId = ref(null);

// 根据当前语言显示食物名称
const displayName = (item) => {
  return locale.value === 'zh-Hans' ? item.name_zh : item.name_en;
};

// 根据输入进行过滤
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

// 当用户在combox输入时
const onComboxInput = (value) => {
  foodNameInput.value = value;
};

// 当用户选择下拉项时
const selectFood = (foodItem) => {
  newDish.value.name = foodItem.name_en; // 内部存英文名
  newDish.value.id = foodItem.id;
  foodNameInput.value = displayName(foodItem); // 显示当前语言名称
  showFoodList.value = false;
};

// 提交菜品提议
const submitDishProposal = async () => {
  // 尝试匹配用户已输入的食物名称
  const matchedFood = availableFoods.find((f) => displayName(f) === foodNameInput.value);

  if (matchedFood) {
    // 选择食物
    selectFood(matchedFood);
  } else {
    uni.showToast({
      title: t('no_matching_food'),
      icon: 'none',
      duration: 2000,
    });
    return;
  }

  if (newDish.value.name.trim() === '') {
    uni.showToast({
      title: t('dish_name_required'),
      icon: 'none'
    });
    return;
  }

  try {
    console.log('newDish:', newDish.value);
    await familyStore.addDishProposal({
      dishId: Number(newDish.value.id),     // dish_id传递为数字
      preference: newDish.value.preference  // level_of_desire
    });

    newDish.value.name = '';
    newDish.value.id = null;
    newDish.value.preference = 0;
    foodNameInput.value = '';

    uni.showToast({
      title: t('submit_success'),
      icon: 'success'
    });
  } catch (error) {
    if (error.message === 'DISH_ALREADY_EXISTS') {
      uni.showToast({
        title: t('dish_already_exists'),
        icon: 'none'
      });
    } else {
      uni.showToast({
        title: t('submit_failed'),
        icon: 'error'
      });
    }
  }
};

// 偏好变化
const onDishPreferenceChange = (e) => {
  newDish.value.preference = parseInt(e.detail.value, 10);
};

// 创建家庭
const createFamily = async () => {
  if (newFamilyName.value.trim() === '') {
    uni.showToast({ title: t('family_name_required'), icon: 'none' });
    return;
  }

  try {
    await familyStore.createFamily(newFamilyName.value);
    newFamilyName.value = '';
    showCreateFamilyModal.value = false;
    uni.showToast({ title: t('create_family_success'), icon: 'success' });
    // 获取菜品提议列表
    await familyStore.getDesiredDishes();
  } catch (error) {
    uni.showToast({ title: t('create_family_failed'), icon: 'error' });
  }
};

// 加入家庭
const joinFamily = async () => {
  if (joinFamilyId.value.trim() === '') {
    uni.showToast({ title: t('family_id_required'), icon: 'none' });
    return;
  }

  try {
    const searchResult = await familyStore.searchFamily(joinFamilyId.value);
    if (!searchResult || !searchResult.id) {
      uni.showToast({ title: t('family_not_found'), icon: 'error' });
      return;
    }

    await familyStore.joinFamily(searchResult.id);
    joinFamilyId.value = '';
    showJoinFamilyModal.value = false;
    startStatusCheck();
    uni.showToast({ title: t('join_request_sent'), icon: 'success' });
    // 获取菜品提议列表
    await familyStore.getDesiredDishes();
  } catch (error) {
    uni.showToast({ title: t('join_family_failed'), icon: 'error' });
  }
};

// 取消加入申请
const handleCancelJoin = async () => {
  try {
    await familyStore.cancelJoinRequest();
    uni.showToast({ title: t('cancel_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('cancel_failed'), icon: 'error' });
  }
};

// 跳转到 myFamily 页面
const goToMyFamily = () => {
  uni.navigateTo({ url: '/pagesMy/myFamily/myFamily' });
};

// 退出家庭
const handleLeaveFamily = async () => {
  try {
    await familyStore.leaveFamily();
    uni.showToast({ title: t('leave_family_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('leave_family_failed'), icon: 'error' });
  }
};

// 解散家庭
const handleBreakFamily = async () => {
  try {
    await familyStore.breakFamily();
    uni.showToast({ title: t('break_family_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('break_family_failed'), icon: 'error' });
  }
};

// 处理成员管理
const manageMembers = () => {
  uni.navigateTo({ url: '/pagesMy/myFamily/myFamily' });
};

// 使用foodStore的getFoodName，根据dish_id返回正确的名称
const sortedDishProposals = computed(() => {
  if (!family.value.dishProposals) return [];
  return family.value.dishProposals.map(item => ({
    id: item.dish_id,
    name: getFoodName(item.dish_id) || 'Unknown Dish',
    preference: item.level_of_desire,
    proposer: item.proposer_user ? item.proposer_user.nickname : 'Unknown'
  })).sort((a, b) => b.preference - a.preference);
});

// 确认删除菜品
const confirmDeleteDish = (dishId) => {
  selectedDishId.value = dishId;
  showDeleteConfirm.value = true;
};

// 删除菜品
const deleteDish = async (dishId) => {
  try {
    await familyStore.deleteDesiredDish(dishId);
    uni.showToast({ title: t('delete_success'), icon: 'success' });
  } catch (error) {
    uni.showToast({ title: t('delete_failed'), icon: 'error' });
  } finally {
    showDeleteConfirm.value = false;
    selectedDishId.value = null;
  }
};

// 判断是否可以删除菜品（提议者或管理员）
const canDeleteDish = (dish) => {
  return dish.proposer === userStore.user.nickName;
};

// ========== 家庭碳排放环形图相关 ==========
// 环形图配置
const carbonRingOpts = {
  title: {
    name: '0 / 0',
    fontSize: 14,
    color: '#333'
  },
  subtitle: {
    name: '',
    fontSize: 12,
    color: '#999'
  },
  // 其他配置项可根据你使用的图表组件自行调整
};

// 环形图数据
const carbonChartData = ref({
  series: []
});

// 监听家庭每日数据，更新环形图 series
watch(
    () => family.value.memberDailyData,
    (newVal) => {
      if (!newVal || newVal.length === 0) {
        carbonChartData.value.series = [];
        carbonRingOpts.title.name = '0 / 0';
        return;
      }
      // 构造 series 数据：每个成员的今日碳排放
      carbonChartData.value.series = newVal.map(member => ({
        name: member.nickname,
        data: member.carbon_intake_sum
      }));

      // 更新环形图中心文字：“家庭总排放量 / 家庭目标量”
      const { familySums } = family.value;
      if (familySums && familySums.carbon_intake_sum !== undefined) {
        carbonRingOpts.title.name = `${familySums.carbon_intake_sum || 0} / ${familySums.carbon_goal_sum || 0}`;
      } else {
        carbonRingOpts.title.name = '0 / 0';
      }
    },
    { immediate: true, deep: true }
);

// ========== 营养超标计算 ==========
// 定义五大营养素对应的多语言key
const nutritionKeys = [
  { key: 'calories', i18nKey: 'calories' },
  { key: 'protein', i18nKey: 'protein' },
  { key: 'fat', i18nKey: 'fat' },
  { key: 'carbohydrates', i18nKey: 'carbohydrates' },
  { key: 'sodium', i18nKey: 'sodium' },
];

// 计算每个成员哪些营养超标
const membersNutritionOver = computed(() => {
  if (!family.value.memberDailyData || family.value.memberDailyData.length === 0) {
    return [];
  }
  return family.value.memberDailyData.map(item => {
    const overList = [];
    nutritionKeys.forEach(nut => {
      const goalVal = item.nutrition_goal[nut.key];
      const intakeVal = item.nutrition_intake_sum[nut.key];
      if (intakeVal > goalVal) {
        overList.push(t(nut.i18nKey));
      }
    });
    return {
      userId: item.user_id,
      nickname: item.nickname,
      avatarUrl: item.avatar_url,
      overs: overList
    };
  });
});

// 页面显示时
onShow(async () => {
  try {
    await familyStore.getFamilyDetails();
    if (family.value.status === FamilyStatusEnum.JOINED) {
      // 获取想吃的菜品列表
      await familyStore.getDesiredDishes();
    }
    if (family.value.status === FamilyStatusEnum.PENDING_APPROVAL) {
      startStatusCheck();
    }
  } catch (error) {
    uni.showToast({ title: t('fetch_family_failed'), icon: 'none' });
  }
});

// 页面mounted时获取可用食物列表
onMounted(async () => {
  if (availableFoods.length === 0) {
    await fetchAvailableFoods();
  }
});

// 监听语言变化，若需要可在此做额外处理
watch(locale, () => {
  // 语言切换后的处理逻辑（可选）
});

// 如果有需求，可以加一个轮询或类似功能
const startStatusCheck = () => {
  // 具体实现根据项目需求
};
</script>

<style scoped>
/* 全局样式变量 */
:root {
  --primary-color: #4CAF50;
  --secondary-color: #2fc25b;
  --background-color: #f5f5f5;
  --card-background: rgba(255, 255, 255, 0.8);
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --font-size-title: 32rpx;
  --font-size-subtitle: 24rpx;
}

/* 容器 */
.container {
  display: flex;
  flex-direction: column;
  background-color: var(--background-color);
  min-height: 100vh;
  padding-bottom: 80rpx;
  position: relative;
  overflow: hidden;
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
  opacity: 0.1;
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

/* 家庭ID */
.list-title {
  margin-left: 10rpx;
  font-size: 20rpx;
  font-weight: bold;
  color: var(--text-color);
  margin-bottom: 20rpx;
  text-align: center;
}

/* 用户未加入家庭时的视图 */
.no-family-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20rpx;
}

.centered-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 50%;
  margin-bottom: 30rpx;
  object-fit: cover;
}

/* 待审核状态的样式 */
.pending-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20rpx;
}

.pending-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 50%;
  margin-bottom: 30rpx;
  object-fit: cover;
}

.pending-text {
  font-size: 32rpx;
  color: #666;
  margin-bottom: 40rpx;
}

.cancel-button {
  background-color: #ff4d4f;
  color: #fff;
  padding: 15rpx 40rpx;
  border-radius: 10rpx;
  font-size: 28rpx;
}

/* 管理员提示框样式 */
.admin-notice {
  margin: 20rpx;
  background: linear-gradient(to right, #EEF2FF, #E0E7FF);
  border-radius: 16rpx;
}

.admin-notice :deep(.uni-card__content) {
  padding: 0 !important;
}

.notice-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
}

.notice-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.notice-text {
  color: #3A86FF;
  font-size: 28rpx;
  font-weight: 500;
}

.notice-btn {
  color: #3A86FF;
  font-size: 28rpx;
  font-weight: 600;
  background-color: rgba(58, 134, 255, 0.1);
  padding: 8rpx 24rpx;
  border-radius: 30rpx;
  transition: all 0.3s ease;
}

.notice-btn:active {
  opacity: 0.8;
  transform: scale(0.98);
}

/* 家庭管理部分 */
.family-management {
  flex: 1;
  padding: 20rpx;
}

/* 创建或加入家庭 */
.family-actions {
  display: flex;
  justify-content: center;
  margin-top: 20rpx;
}

.action-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx 30rpx;
  margin: 0 10rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
}

/* 提出想吃的菜品 */
.proposals-scroll {
  max-height: 600rpx;
  margin: 20rpx 0;
}

/* 折叠面板标题样式 */
.collapse-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8rpx 0;
}

.dish-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-right: 16rpx;
}

.dish-info {
  display: flex;
  align-items: center;
  gap: 16rpx;
  flex-shrink: 0;
}

.proposer-text {
  font-size: 24rpx;
  color: #666;
}

/* 折叠面板内容样式 */
.dish-actions {
  display: flex;
  justify-content: flex-end;
  padding: 16rpx 0;
}

/* uni-collapse 样式优化 */
.dish-list :deep(.uni-collapse) {
  background-color: transparent;
}

.dish-list :deep(.uni-collapse-item__title) {
  padding: 16rpx;
}

.dish-list :deep(.uni-collapse-item__title-box) {
  width: 100%;
}

.dish-list :deep(.uni-collapse-item__wrap) {
  background-color: rgba(255, 255, 255, 0.6);
}

.dish-proposal {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.section-title {
  font-size: var(--font-size-title);
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 10rpx;
}

/* 输入框 */
.input {
  width: 100%;
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin-bottom: 20rpx;
  margin-top: 20rpx;
}

.picker {
  width: 100%;
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin-bottom: 20rpx;
  margin-top: 20rpx;
  color: #666666;
}

.submit-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
}

/* 菜品列表 */
.dish-list {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.dish-list .uni-collapse__item {
  border-bottom: 1rpx solid #eee;
}

.dish-list .uni-collapse__title {
  font-size: 28rpx;
  color: #333;
  font-weight: bold;
}

.dish-details {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 10rpx 0;
}

.delete-icon {
  margin-right: 10rpx;
}

/* 共享数据 */
.shared-data {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.charts {
  width: 100%;
  height: 300rpx;
  margin-bottom: 20rpx;
}

/* 新增：营养超标情况 */
.nutrition-over-section {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.nutrition-over-row {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
}

.nutrition-over-row .member-avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}

.nutrition-over-row .member-name {
  font-size: 26rpx;
  font-weight: 500;
  margin-right: 20rpx;
}

.over-list {
  font-size: 24rpx;
  color: #ff4d4f;
}

/* 家庭成员部分 */
.family-info {
  background-color: var(--card-background);
  padding: 20rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.family-members {
  display: flex;
  flex-wrap: wrap;
}

.member {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 25%;
  margin-bottom: 20rpx;
}

.member-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin-bottom: 5rpx;
}

.member-name {
  font-size: 20rpx;
  color: var(--text-color);
}

/* 管理成员按钮 */
.manage-members-button {
  background-color: var(--secondary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
  margin-top: 10rpx;
}

/* 底部操作按钮样式 */
.bottom-container {
  padding: 20rpx 30rpx 40rpx;
  margin-top: auto;
}

.action-section :deep(.uni-section-header) {
  display: none;
}

.bottom-actions {
  display: flex;
  justify-content: center;
  gap: 30rpx;
  padding: 20rpx 0;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 20rpx 40rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-btn:active {
  opacity: 0.8;
  transform: scale(0.98);
}

.leave-btn {
  background: linear-gradient(to right, #FF9800, #F57C00);
  color: #fff;
  box-shadow: 0 4rpx 12rpx rgba(255, 152, 0, 0.3);
}

.break-btn {
  background: linear-gradient(to right, #FF4B4B, #E53935);
  color: #fff;
  box-shadow: 0 4rpx 12rpx rgba(229, 57, 53, 0.3);
}

/* 模态框 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
}

.modal-content {
  background-color: #fff;
  padding: 40rpx;
  border-radius: 16rpx;
  width: 60%;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}

.modal-title {
  font-size: 36rpx;
  color: #333;
  margin-bottom: 40rpx;
  text-align: center;
  font-weight: 500;
}

.modal-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 20rpx;
  border-radius: 12rpx;
  font-size: 24rpx;
  width: 100%;
  text-align: center;
  margin-top: 20rpx;
  font-weight: 500;
}

.modal-button.cancel {
  background-color: #f5f5f5;
  color: #666;
}

/* 动画效果 */
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

</style>
