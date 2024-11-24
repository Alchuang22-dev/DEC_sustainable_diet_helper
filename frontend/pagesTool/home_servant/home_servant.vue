<template>
  <view class="container">
    <!-- 全屏背景图片 -->
    <image src="/static/images/index/background_index_new.png" class="background-image"></image>

    <!-- 头部标题 -->
    <view v-if="family.id" class="header">
      <text class="header-title">{{ family.name || $t('default_family_name') }}</text>
    </view>

    <!-- 家庭ID -->
    <text v-if="family.id" class="list-title">{{ $t('family_info') + family.id  }}</text>

    <!-- 用户未加入家庭时的视图 -->
    <view v-if="!family.id" class="no-family-view">
      <image src="https://cdn.pixabay.com/photo/2017/01/13/02/31/family-1976162_1280.png" class="centered-image"></image>
      <view class="family-actions">
        <button class="action-button" @click="showCreateFamilyModal = true">{{ t('create_family') }}</button>
        <button class="action-button" @click="showJoinFamilyModal = true">{{ t('join_family') }}</button>
      </view>
    </view>

    <!-- 用户已加入家庭时的家庭管理部分 -->
    <view v-else class="family-management">
      <!-- 提出想吃的菜品 -->
      <view class="dish-proposal">
        <text class="section-title">{{ t('propose_dish') }}</text>
        <input v-model="newDish.name" :placeholder="t('dish_name_placeholder')" class="input"></input>
        <picker mode="selector" :range="dishPreferenceLevels" :value="newDish.preference" @change="onDishPreferenceChange" class="picker">
          <view class="picker-content">{{ $t('dish_preference') }}: {{ dishPreferenceLevels[newDish.preference] }}</view>
        </picker>
        <button class="submit-button" @click="submitDishProposal">{{ t('submit_proposal') }}</button>
      </view>

      <!-- 家庭成员的提议 -->
      <view class="dish-list">
        <text class="section-title">{{ $t('family_dish_proposals') }}</text>
        <view v-for="dish in family.dishProposals" :key="dish.id" class="dish-item">
          <text class="dish-name">{{ dish.name }}</text>
          <text class="dish-preference">{{ $t('preference_level') }}: {{ dishPreferenceLevels[dish.preference] }}</text>
          <text class="dish-proposer">{{ $t('proposed_by') }}: {{ dish.proposer }}</text>
        </view>
      </view>

      <!-- 共享家庭成员的五大营养成分达标情况 -->
      <view class="shared-data">
        <text class="section-title">{{ $t('shared_nutrition_compliance') }}</text>
        <!-- 添加家庭五大营养成分达标情况的图表 -->
        <view class="charts">
          <qiun-data-charts :canvas2d="true" canvas-id="familyNutrientChart" type="column" :opts="nutrientChartOpts" :chartData="nutrientChartData" />
        </view>
      </view>

      <!-- 家庭成员部分（移到最下方） -->
      <view class="family-info">
        <text class="section-title">{{ $t('family_members') }}</text>
        <view class="family-members">
          <view v-for="member in family.members" :key="member.id" class="member">
            <image :src="member.avatar" class="member-avatar"></image>
            <text class="member-name">{{ `${member.name}(${t(member.family_name)})` }}</text>
          </view>
        </view>
        <!-- 添加“管理成员”按钮 -->
        <button class="manage-members-button" @click="manageMembers">{{ t('manage_members') }}</button>
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
  </view>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFamilyStore } from '@/stores/family.js';

// 国际化
const { t } = useI18n();

// Pinia 状态管理
const familyStore = useFamilyStore();
const family = computed(() => familyStore.family);

// 新菜品提议
const newDish = reactive({
  name: '',
  preference: 0,
});

// 菜品偏好级别（响应式以支持语言切换）
const dishPreferenceLevels = computed(() => [
  t('preference_low'),
  t('preference_medium'),
  t('preference_high'),
]);

// 模态框显示状态
const showCreateFamilyModal = ref(false);
const showJoinFamilyModal = ref(false);

// 创建家庭相关
const newFamilyName = ref('');

// 加入家庭相关
const joinFamilyId = ref('');

// 提交菜品提议
const submitDishProposal = () => {
  if (newDish.name.trim() === '') {
    uni.showToast({ title: t('dish_name_required'), icon: 'none' });
    return;
  }
  // 假设 'You' 是当前用户的名称，可以根据实际情况替换
  familyStore.addDishProposal({
    id: Date.now(),
    name: newDish.name,
    preference: newDish.preference,
    proposer: 'You',
  });
  newDish.name = '';
  newDish.preference = 0;
};

// 创建家庭
const createFamily = () => {
  if (newFamilyName.value.trim() === '') {
    uni.showToast({ title: t('family_name_required'), icon: 'none' });
    return;
  }
  familyStore.createFamily(newFamilyName.value);
  newFamilyName.value = '';
  showCreateFamilyModal.value = false;
};

// 加入家庭
const joinFamily = () => {
  if (joinFamilyId.value.trim() === '') {
    uni.showToast({ title: t('family_id_required'), icon: 'none' });
    return;
  }
  familyStore.joinFamily(joinFamilyId.value);
  joinFamilyId.value = '';
  showJoinFamilyModal.value = false;
};

// 家庭五大营养成分达标情况数据（示例）
const nutrientChartData = ref({
  categories: [],
  series: [
    {
      name: t('nutrient_compliance'),
      data: [],
    },
  ],
});

// 图表配置（示例）
const nutrientChartOpts = {
  color: ["#1890FF"],
  padding: [15, 15, 0, 5],
  xAxis: {
    disableGrid: false,
    axisLine: true,
  },
  yAxis: {},
  extra: {
    column: {
      type: "group",
      width: 30,
    },
  },
};

// 模拟获取家庭五大营养成分达标数据
onMounted(() => {
  nutrientChartData.value.categories = [
    t('protein'),
    t('fat'),
    t('carbohydrate'),
    t('vitamin'),
    t('mineral'),
  ];
  nutrientChartData.value.series[0].data = [80, 90, 75, 85, 70]; // 示例数据，表示达标百分比
});

// 处理 <picker> 变化
const onDishPreferenceChange = (e) => {
  newDish.preference = parseInt(e.detail.value, 10);
};

// 管理成员功能（占位）
const manageMembers = () => {
  // 实现管理成员的逻辑
  uni.showToast({ title: t('feature_not_implemented'), icon: 'none' });
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
  /* 添加其他样式如边框或阴影根据需要 */
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
.dish-proposal {
  display: block;
  background-color: var(--card-background);
  padding: 90rpx;
  border-radius: 10rpx;
  box-shadow: 0 2rpx 5rpx var(--shadow-color);
  margin-bottom: 20rpx;
}

.section-title {
  font-size: var(--font-size-title);
  font-weight: bold;
  width: 50%;
  color: var(--primary-color);
  //margin-left: 0rpx;
  margin-bottom: 10rpx;
}

.input {
  width: 100%;
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  margin-bottom: 15rpx;
}

.picker {
  width: 100%;
  margin-bottom: 15rpx;
}

.picker-content {
  padding: 10rpx;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  background-color: #fff;
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

.dish-item {
  padding: 10rpx;
  border-bottom: 1rpx solid #eee;
}

.dish-name {
  font-size: 24rpx;
  color: var(--text-color);
  margin-bottom: 5rpx;
}

.dish-preference,
.dish-proposer {
  font-size: 20rpx;
  color: #666;
  margin-bottom: 5rpx;
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
  height: 300px;
}

/* 家庭成员部分（移到最下方） */
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
  padding: 20rpx;
  border-radius: 10rpx;
  width: 80%;
}

.modal-title {
  font-size: var(--font-size-subtitle);
  color: var(--primary-color);
  margin-bottom: 15rpx;
  text-align: center;
}

.modal-button {
  background-color: var(--primary-color);
  color: #fff;
  padding: 15rpx;
  border-radius: 10rpx;
  font-size: var(--font-size-subtitle);
  width: 100%;
  text-align: center;
  margin-top: 10rpx;
}

.modal-button.cancel {
  background-color: #ccc;
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

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

