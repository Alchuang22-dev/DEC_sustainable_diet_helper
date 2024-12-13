<template>
	<view class="container">
		<text class="title">为家人共同计算</text>
		<scroll-view scroll-y style="max-height:600rpx;">
			<view v-for="member in familyMembers" :key="member.id" class="member-row">
				<image :src="`http://122.51.231.155:8080/static/${member.avatarUrl}`" class="avatar"></image>
				<text class="name">{{ member.nickname }}</text>
				<input class="ratio-input" type="number" v-model="memberRatio[member.id]" placeholder="0.00"
					@input="validateRatio" />
			</view>
		</scroll-view>
		<button class="confirm-button" @click="submitData">确认提交</button>
	</view>
</template>

<script setup>
	import {
		ref,
		onMounted,
		reactive
	} from 'vue'
	import {
		useFamilyStore
	} from "@/stores/family.js";
	import {
		onLoad
	} from '@dcloudio/uni-app'

	// 碳计算器得来的数据
	const carbonEmissionData = reactive({});
	const nutritionData = reactive({});

	onLoad((options) => {
		if (options && options.data) {
			try {
				const parsedData = JSON.parse(decodeURIComponent(options.data));
				console.log('接收到的计算结果:', parsedData);

				// 使用 reactive 保存数据
				Object.assign(carbonEmissionData, parsedData.carbonEmission);
				Object.assign(nutritionData, parsedData.nutrition);
			} catch (error) {
				console.error('解析传递的数据失败:', error);
			}
		}
	});

	const familyStore = useFamilyStore();
	const familyMembers = ref([]);
	const memberRatio = ref({});

	onMounted(() => {
		familyMembers.value = familyStore.family.allMembers;
		familyMembers.value.forEach(m => {
			memberRatio.value[m.id] = 0;
		});
	});

	const validateRatio = () => {
		// 确保所有比例之和不超过1
		let sum = 0;
		for (let id in memberRatio.value) {
			let val = parseFloat(memberRatio.value[id]) || 0;
			if (val < 0) {
				memberRatio.value[id] = 0;
			} else if (val > 1) {
				memberRatio.value[id] = 1;
			}
			sum += val;
		}
		if (sum > 1) {
			uni.showToast({
				title: '比例总和不能超过1',
				icon: 'none',
				duration: 2000
			});
		}
	};

	const submitData = () => {
		let sum = 0;
		for (let id in memberRatio.value) {
			sum += parseFloat(memberRatio.value[id]) || 0;
		}
		if (sum > 1) {
			uni.showToast({
				title: '比例总和不能超过1',
				icon: 'none',
				duration: 2000
			});
			return;
		}

		const data = {
			carbonEmission: carbonEmissionData,
			nutrition: nutritionData,
			user_shares: Object.keys(memberRatio.value).map(id => ({
				user_id: parseInt(id),
				ratio: parseFloat(memberRatio.value[id]) || 0
			}))
		};

		uni.request({
			url: 'http://122.51.231.155:8080/nutrition-carbon/shared/nutrition-carbon',
			method: 'POST',
			data,
			success: (res) => {
				if (res.statusCode === 200) {
					uni.showToast({
						title: '提交成功',
						icon: 'success',
						duration: 2000
					});
					uni.navigateBack({
						delta: 2
					});
				} else {
					uni.showToast({
						title: '提交失败',
						icon: 'none',
						duration: 2000
					});
				}
			},
			fail: () => {
				uni.showToast({
					title: '提交失败',
					icon: 'none',
					duration: 2000
				});
			}
		});
	};
</script>

<style scoped>
	.container {
		padding: 20rpx;
		background: #f5f5f5;
	}

	.title {
		font-size: 36rpx;
		text-align: center;
		margin-bottom: 20rpx;
	}

	.member-row {
		display: flex;
		align-items: center;
		background: #fff;
		padding: 20rpx;
		margin-bottom: 10rpx;
		border-radius: 10rpx;
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
</style>