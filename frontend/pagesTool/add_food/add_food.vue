<template>
	<view class="container">

		<!-- 表单容器 -->
		<view class="form-container">
			<view class="form-group">
				<text class="label">名称</text>
				<input class="input" type="text" v-model="food.name" placeholder="请输入食品名称" />
			</view>
			<view class="form-group">
				<text class="label">总重量 (g)</text>
				<input class="input" type="number" v-model="food.weight" placeholder="请输入食品重量" :error="weightError" />
				<text v-if="weightError" class="error-message">重量必须是正整数</text>
			</view>
			<view class="form-group">
				<text class="label">总价格 (元)</text>
				<input class="input" type="number" v-model="food.price" placeholder="请输入食品价格" :error="priceError" />
				<text v-if="priceError" class="error-message">价格必须是正整数</text>
			</view>
			<view class="form-group">
				<text class="label">请选择运输方式</text>
				<picker mode="selector" :range="transportMethods" :value="transportIndex" @change="onTransportChange">
					<view class="picker">
						{{ transportMethods[transportIndex] }}
					</view>
				</picker>
			</view>
			<view class="form-group">
				<text class="label">请选择食品来源</text>
				<picker mode="selector" :range="foodSources" :value="sourceIndex" @change="onSourceChange">
					<view class="picker">
						{{ foodSources[sourceIndex] }}
					</view>
				</picker>
			</view>

			<!-- 图片上传按钮 -->
			<view class="form-group">
				<text class="label">上传食品图片</text>
				<button class="upload-button" @click="uploadImage">拍照上传</button>
				<image v-if="food.imagePath" :src="food.imagePath" class="uploaded-image"></image>
			</view>

			<button class="submit-button" @click="submitFoodDetails">提交</button>
		</view>
	</view>
</template>

<script setup>
	import {
		ref,
		reactive
	} from 'vue';

	// 食品数据
	const food = reactive({
		name: '',
		weight: '',
		price: '',
		transportMethod: '陆运',
		foodSource: '本地',
		imagePath: '', // 图片路径
	});

	// 下拉选项数据
	const transportMethods = ['陆运', '海运', '空运'];
	const foodSources = ['本地', '进口'];

	// 当前选择的索引
	const transportIndex = ref(0);
	const sourceIndex = ref(0);

	// 输入验证错误状态
	const weightError = ref(false);
	const priceError = ref(false);

	// 返回上一页
	const navigateBack = () => {
		uni.navigateBack();
	};

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
					title: '图片上传失败',
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
			foodSource
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
				title: '请填写所有字段',
				icon: 'none',
			});
			valid = false;
		}

		if (!valid) {
			return;
		}

		const newFood = {
			name,
			weight: parseInt(weight),
			price: parseInt(price),
			transportMethod,
			foodSource,
			imagePath: food.imagePath, // 保存图片路径
		};

		// 从本地存储获取已有的食物列表
		let foodList = uni.getStorageSync('foodDetails') || [];
		foodList.push(newFood);
		uni.setStorageSync('foodDetails', foodList);

		uni.showToast({
			title: '添加成功',
			icon: 'success',
			duration: 2000,
		});

		// 返回上一页并刷新
		setTimeout(() => {
			uni.navigateBack();
		}, 2000);
	};
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
</style>