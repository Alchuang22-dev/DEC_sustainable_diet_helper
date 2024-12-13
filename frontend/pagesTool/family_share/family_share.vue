<template>
	<view class="container">
		<text class="title">{{ t('shared_calculation_for_family') }}</text>
		<scroll-view scroll-y style="max-height:600rpx;">
			<view v-for="member in allFamilyMembers" :key="member.id" class="member-row">
				<image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" class="avatar"></image>
				<text class="name">{{ member.nickname }}</text>
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
		</scroll-view>
		<view class="total-ratio">
			<text>{{ t('total_ratio') }}: {{ totalRatio.toFixed(2) }} / 1.00</text>
		</view>
		<button class="confirm-button" :disabled="!isValid" @click="submitData">
			{{ t('confirm_submission') }}
		</button>
	</view>
</template>

<script setup>
import {
	ref,
	onMounted,
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

const selectedMealIndex = ref(0);

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
	// 设置默认的选中项（如果需要 Picker，可以保留）
	selectedMealIndex.value = getDefaultMealType();
});

// 计算当前时间对应的餐食类型
const getDefaultMealType = () => {
	const hour = new Date().getHours();
	if (hour >= 5 && hour < 11) return 0; // 早餐
	if (hour >= 11 && hour < 15) return 1; // 午餐
	if (hour >= 15 && hour < 20) return 2; // 晚餐
	return 3; // 其他
};

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

	const requestData = {
		date: new Date().toISOString(),
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
  padding: 20rpx;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.title {
  font-size: 36rpx;
  text-align: center;
  margin-bottom: 20rpx;
  color: #333;
}

.member-row {
  display: flex;
  align-items: center;
  background: #fff;
  padding: 20rpx;
  margin-bottom: 10rpx;
  border-radius: 10rpx;
  width: 100%;
}

.avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 40rpx;
  margin-right: 20rpx;
}

.name {
  font-size: 28rpx;
  flex: 1;
  color: #333;
}

.ratio-input {
  width: 120rpx;
  height: 60rpx;
  background: #fff;
  border: 1rpx solid #ccc;
  border-radius: 5rpx;
  text-align: right;
  padding-right: 10rpx;
  font-size: 24rpx;
  color: #333;
}

.confirm-button {
  width: 100%;
  background: #4CAF50;
  color: #fff;
  text-align: center;
  font-size: 28rpx;
  padding: 20rpx 0;
  border-radius: 20rpx;
  margin-top: 20rpx;
}

.confirm-button:disabled {
  background: #ccc;
}

.total-ratio {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #333;
}
</style>
