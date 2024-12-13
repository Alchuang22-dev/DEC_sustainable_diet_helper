<template>
	<view class="container">
		<view class="header">
			<text class="title">{{ t('select_save_method') }}</text>
		</view>

		<!-- 添加的 Picker 组件 -->
		<view class="picker-container">
			<text class="picker-label">{{ t('select_meal_type') }}</text>
			<picker mode="selector" :range="mealTypesDisplay" :value="selectedMealIndex" @change="onPickerChange">
				<view class="picker">
					{{ mealTypesDisplay[selectedMealIndex] }}
					<uni-icons type="arrow-right" size="24" color="#999"></uni-icons>
				</view>
			</picker>
		</view>

		<view class="button-group">
			<button class="primary-button" @click="saveForSelf">{{ t('save_for_self') }}</button>
			<button class="secondary-button" @click="saveForFamily">{{ t('save_for_family') }}</button>
		</view>
	</view>
</template>

<script setup>
	import {
		reactive,
		ref,
		computed
	} from 'vue';
	import {
		useI18n
	} from 'vue-i18n';
	import {
		useUserStore
	} from "@/stores/user.js";
  import {
    FamilyStatus,
    useFamilyStore
  } from "@/stores/family.js";
	import {
		onLoad
	} from "@dcloudio/uni-app";

	// 使用国际化
	const {
		t
	} = useI18n();
	const userStore = useUserStore();
	const familyStore = useFamilyStore();
  const familyStatus = computed(() => familyStore.family.status);

	// 获取当前用户的 UID
	const uid = computed(() => userStore.user.uid);
	const token = computed(() => userStore.user.token);

	// 定义餐食类型的显示名称和对应的英文值
	const mealTypesDisplay = [t('breakfast'), t('lunch'), t('dinner'), t('other')];
	const mealTypesValue = ['breakfast', 'lunch', 'dinner', 'other'];

	const selectedMealIndex = ref(0);

	// 当用户选择改变时的回调
	const onPickerChange = (e) => {
		selectedMealIndex.value = e.detail.value;
	};

	// 传递给save_options页面的数据
	const carbonEmissionData = reactive({});
	const nutritionData = reactive({});
	const totalEmission = ref(0);

	// 接收传递的数据
	onLoad((options) => {
		if (options && options.data) {
			try {
				const parsedData = JSON.parse(decodeURIComponent(options.data));
				console.log('接收到的计算结果:', parsedData);

				// 使用 reactive 保存数据
				Object.assign(carbonEmissionData, parsedData.carbonEmission);
				Object.assign(nutritionData, parsedData.nutrition.series[0].data);
				totalEmission.value = carbonEmissionData.series[0].data.reduce((sum, item) => sum + item.value, 0);
				console.log('Total Emission:', totalEmission.value);
				console.log('carbonEmissionData:', carbonEmissionData.series[0].data);
				console.log('nutritionData:', nutritionData);
			} catch (error) {
				console.error('解析传递的数据失败:', error);
			}
		}
		// 设置默认的选中项
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

	// 定义 Picker 的选项
	const pickerOptions = ref(mealTypesDisplay);

	// 保存为自己计算
	const saveForSelf = () => {
		// 构建请求数据
		const requestData = {
			date: new Date().toISOString(),
			meal_type: mealTypesValue[selectedMealIndex.value],
			calories: nutritionData[0] || 0,
			protein: nutritionData[1] || 0,
			fat: nutritionData[2] || 0,
			carbohydrates: nutritionData[3] || 0,
			sodium: nutritionData[4] || 0,
			emission: totalEmission.value || 0,
			user_shares: [{
				user_id: uid.value,
				ratio: 1.0
			}]
		};

		console.log('发送的数据:', requestData);

		// 发送 POST 请求到共享营养碳排放接口
		uni.request({
			url: 'http://122.51.231.155:8095/nutrition-carbon/shared/nutrition-carbon',
			method: 'POST',
			data: requestData,
			header: {
				'Content-Type': 'application/json',
				// 假设使用 Bearer Token 进行认证
				'Authorization': `Bearer ${token.value}`
			},
			success: (res) => {
				if (res.statusCode === 200) {
					uni.showToast({
						title: t('save_success'),
						icon: 'success',
						duration: 2000
					});
					uni.navigateBack();
				} else {
					// 如果后端返回了错误信息，显示具体错误
					const errorMsg = res.data && res.data.error ? res.data.error : t('save_failed');
					uni.showToast({
						title: errorMsg,
						icon: 'none',
						duration: 2000
					});
				}
			},
			fail: () => {
				uni.showToast({
					title: t('save_failed'),
					icon: 'none',
					duration: 2000
				});
			}
		});
	};

	// 保存为家人计算
	const saveForFamily = () => {
    // 如果家庭状态不是已加入家庭，则提示用户加入家庭
    if (familyStatus.value !== FamilyStatus.JOINED) {
      uni.showToast({
        title: t('join_family_first'),
        icon: 'none',
        duration: 2000
      });
      return;
    }
		uni.navigateTo({
			url: `/pagesTool/family_share/family_share?data=${encodeURIComponent(JSON.stringify({
			carbonEmission: totalEmission.value,
			nutrition: nutritionData,
			mealType: mealTypesValue[selectedMealIndex.value] // 使用英文餐食类型
		}))}`
		});
	};

	// 获取当前选择的餐食类型
	const currentMealType = computed(() => {
		return mealTypesValue[selectedMealIndex.value];
	});
</script>

<style scoped>
	.container {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		height: 100%;
		background: #f5f5f5;
	}

	.header {
		margin-bottom: 40rpx;
	}

	.title {
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
	}

	.picker-container {
		width: 80%;
		margin-bottom: 40rpx;
	}

	.picker-label {
		font-size: 28rpx;
		color: #333;
		margin-bottom: 10rpx;
	}

	.picker {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx;
		background: #fff;
		border: 1rpx solid #ccc;
		border-radius: 10rpx;
		font-size: 28rpx;
		color: #333;
	}

	.button-group {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}

	.primary-button,
	.secondary-button {
		padding: 20rpx 40rpx;
		border-radius: 20rpx;
		font-size: 28rpx;
		color: #fff;
		border: none;
		text-align: center;
	}

	.primary-button {
		background-color: #4CAF50;
	}

	.secondary-button {
		background-color: #ff9800;
	}
</style>