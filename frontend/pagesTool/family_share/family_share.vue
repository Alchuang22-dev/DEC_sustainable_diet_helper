<template>
  <view class="container">
    <text class="title">{{ t('shared_calculation_for_family') }}</text>
    <scroll-view scroll-y class="member-list">
      <view v-for="member in allFamilyMembers" :key="member.id" class="member-row">
        <view class="member-info">
          <image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" class="avatar"></image>
          <text class="name">{{ member.nickname }}</text>
        </view>
        <view class="input-container">
          <input
              class="ratio-input"
              type="number"
              v-model.number="memberRatio[member.id]"
              placeholder="0.00"
              step="0.01"
              min="0"
              max="1"
              @input="validateRatio"
          />
        </view>
      </view>
    </scroll-view>
    <view class="footer">
      <view class="total-ratio" :class="{ 'warning': totalRatio.value > 1 }">
        <text>{{ t('total_ratio') }}: {{ totalRatio.toFixed(2) }} / 1.00</text>
      </view>
      <button class="confirm-button" :disabled="!isValid" @click="submitData">
        {{ t('confirm_submission') }}
      </button>
    </view>
  </view>
</template>

<script setup>
import {
	ref,
	reactive,
	computed,
	watch
} from 'vue';
import {
	useFamilyStore
} from "@/stores/family.js";
import {
	useUserStore
} from "@/stores/user.js";
import {
	useI18n
} from 'vue-i18n';
import {
	onLoad, onShow
} from '@dcloudio/uni-app'

// 使用国际化
const {
	t
} = useI18n();

const userStore = useUserStore();
const familyStore = useFamilyStore();
const allFamilyMembers = computed(() => familyStore.family.allMembers);


// 成员比例
const memberRatio = reactive({});

// 初始化比例
onShow(() => {
  familyStore.getFamilyDetails();
  for (let member of allFamilyMembers.value) {
    memberRatio[member.id] = 0;
  }
});

// 获取当前用户的 UID 和 Token
const token = computed(() => userStore.user.token);

// 传递给页面的数据
const carbonEmissionData = ref(0);
const nutritionData = reactive({});
const mealType = ref('');

// 接收传递的数据
onLoad((options) => {
	if (options && options.data) {
		try {
			const parsedData = JSON.parse(decodeURIComponent(options.data));
			console.log('接收到的计算结果:', parsedData);

			// 使用 reactive 保存数据
			carbonEmissionData.value = parsedData.carbonEmission;
      console.log('carbonEmissionData:', carbonEmissionData.value);
			Object.assign(nutritionData, parsedData.nutrition);
      mealType.value = parsedData.mealType;
		} catch (error) {
			console.error('解析传递的数据失败:', error);
		}
	}
});

// 计算总比例
const totalRatio = computed(() => {
	return Object.values(memberRatio).reduce((sum, val) => sum + (parseFloat(val) || 0), 0);
});

// 验证比例
const isValid = computed(() => {
	// 总比例必须小于或等于1，并且每个比例在0到1之间
	return totalRatio.value <= 1 && Object.values(memberRatio).every(val => val >= 0 && val <= 1);
});

// 监听比例变化，实时验证
watch(memberRatio, () => {
	validateRatio();
}, { deep: true });

// 验证比例的方法
const validateRatio = () => {
	// 确保所有比例在0到1之间
	let sum = 0;
	for (let id in memberRatio) {
		let val = parseFloat(memberRatio[id]) || 0;
		if (val < 0) {
			memberRatio[id] = 0;
		} else if (val > 1) {
			memberRatio[id] = 1;
		}
		sum += memberRatio[id];
	}
	if (sum > 1) {
		uni.showToast({
			title: t('total_ratio_cannot_exceed_one'),
			icon: 'none',
			duration: 2000
		});
	}
};

// 提交数据的方法
const submitData = () => {
	if (totalRatio.value > 1) {
		uni.showToast({
			title: t('total_ratio_cannot_exceed_one'),
			icon: 'none',
			duration: 2000
		});
		return;
	}

	// 可选：提醒用户如果总和 <1，可能会有未分配的比例
	if (totalRatio.value < 1) {
		uni.showToast({
			title: t('total_ratio_less_than_one'),
			icon: 'none',
			duration: 2000
		});
	}

	const userShares = Object.keys(memberRatio).map(id => ({
		user_id: parseInt(id),
		ratio: parseFloat(memberRatio[id]) || 0
	}));

  const today = new Date();
	const requestData = {
    date: today.getFullYear() + '-'
        + String(today.getMonth() + 1).padStart(2, '0') + '-'
        + String(today.getDate()).padStart(2, '0'),
		meal_type: mealType.value,
		calories: nutritionData[0] || 0,
		protein: nutritionData[1] || 0,
		fat: nutritionData[2] || 0,
    carbohydrates: nutritionData[3] || 0,
    sodium: nutritionData[4] || 0,
    emission: carbonEmissionData.value || 0,
    user_shares: userShares
  };

  console.log('发送的数据:', requestData);

  // 发送 POST 请求到共享营养碳排放接口
  uni.request({
    url: 'http://122.51.231.155:8095/nutrition-carbon/shared/nutrition-carbon',
    method: 'POST',
    data: requestData,
    header: {
      'Content-Type': 'application/json',
      // 使用 Bearer Token 进行认证
      'Authorization': `Bearer ${token.value}`
    },
    success: (res) => {
      if (res.statusCode === 200) {
        uni.showToast({
          title: t('submission_success'),
          icon: 'success',
          duration: 2000
        });
        setTimeout(() => {
          uni.navigateBack({
            delta: 2
          });
        }, 2000);
      } else {
        // 如果后端返回了错误信息，显示具体错误
        const errorMsg = res.data && res.data.error ? res.data.error : t('submission_failed');
        uni.showToast({
          title: errorMsg,
          icon: 'none',
          duration: 2000
        });
        setTimeout(() => {
          uni.navigateBack({
            delta: 2
          });
        }, 2000);
      }
    },
    fail: () => {
      uni.showToast({
        title: t('submission_failed'),
        icon: 'none',
        duration: 2000
      });
      setTimeout(() => {
        uni.navigateBack({
          delta: 2
        });
      }, 2000);
    }
  });
};
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #f8f9fa;
  padding: 32rpx;
}

.title {
  font-size: 40rpx;
  font-weight: 600;
  text-align: center;
  margin-bottom: 32rpx;
  color: #4CAF50FF;
  padding: 16rpx;
}

.member-list {
  flex: 1;
  max-height: calc(100vh - 300rpx);
  margin-bottom: 32rpx;
}

.member-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.member-info {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0; /* Prevents flex items from overflowing */
}

.avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 44rpx;
  margin-right: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
}

.name {
  font-size: 32rpx;
  color: #2c3e50;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 24rpx;
}

.input-container {
  min-width: 160rpx; /* Ensures input area doesn't shrink too much */
}

.ratio-input {
  width: 160rpx;
  height: 72rpx;
  background: #f8f9fa;
  border: 2rpx solid #e9ecef;
  border-radius: 12rpx;
  text-align: center;
  font-size: 28rpx;
  color: #2c3e50;
  padding: 0 16rpx;
  transition: border-color 0.3s ease;
}

.ratio-input:focus {
  border-color: #4CAF50;
  background: white;
}

.footer {
  padding: 24rpx 0;
  background: #f8f9fa;
}

.total-ratio {
  text-align: center;
  font-size: 32rpx;
  color: #2c3e50;
  margin-bottom: 24rpx;
  padding: 16rpx;
  background: white;
  border-radius: 12rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.total-ratio.warning {
  color: #dc3545;
  background: #fff5f5;
}

.confirm-button {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  background: #4CAF50;
  color: white;
  font-size: 32rpx;
  font-weight: 600;
  border-radius: 16rpx;
  box-shadow: 0 4rpx 12rpx rgba(76, 175, 80, 0.2);
  transition: all 0.3s ease;
}

.confirm-button:active {
  transform: translateY(2rpx);
  box-shadow: 0 2rpx 6rpx rgba(76, 175, 80, 0.2);
}

.confirm-button:disabled {
  background: #a5d6a7;
  box-shadow: none;
  opacity: 0.7;
}
</style>
