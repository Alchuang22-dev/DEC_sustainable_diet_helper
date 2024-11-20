<template>
	<view class="container">
		<!-- 顶部导航 -->
		<view class="header">
			<u-navbar title="{{$t('nutrition_calendar')}}" fixed="true"></u-navbar>
		</view>

		<!-- 日期选择器 -->
		<view class="date-picker">
			<u-segmented-control :current="currentDateIndex" :values="dateTabs" @change="onDateChange">
			</u-segmented-control>
		</view>

		<!-- 圆形统计信息 -->
		<view class="circle-info">
			<u-circle-progress percent="80" stroke-width="12" :color="['#FF6F61', '#FDE7E4']" class="circle-progress">
				<view class="circle-content">
					<text class="large-text">25065</text>
					<text class="small-text">{{$t('more_than_recommended')}}: 1861</text>
				</view>
			</u-circle-progress>
		</view>

		<!-- 营养数据 -->
		<view class="nutrition-info">
			<view class="nutrition-item" v-for="item in nutritionData" :key="item.label">
				<text class="label">{{ item.label }}</text>
				<u-line-progress :percentage="item.percentage" active :stroke-width="8" :color="item.color">
				</u-line-progress>
				<text class="values">{{ item.value }}</text>
			</view>
		</view>

		<!-- 建议部分 -->
		<view class="suggestion">
			<view class="suggestion-text">
				<u-notice-bar show-icon mode="closeable" color="#FF6F61">
					{{ $t('exercise_to_offset') }}: {{ exerciseSuggestion }}
				</u-notice-bar>
			</view>
		</view>

		<!-- 饮食记录 -->
		<view class="meal-records">
			<view class="meal" v-for="meal in meals" :key="meal.title">
				<view class="meal-header">
					<text class="meal-title">{{ meal.title }}</text>
					<text class="calories">{{ meal.totalCalories }} {{$t('calories')}}</text>
				</view>
				<view class="meal-item" v-for="item in meal.items" :key="item.name">
					<image :src="item.image" class="meal-image"></image>
					<view class="meal-info">
						<text class="meal-name">{{ item.name }}</text>
						<text class="meal-calories">{{ item.calories }} {{$t('calories')}}</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref
	} from 'vue'
	import {
		useI18n
	} from 'vue-i18n'

	const {
		t
	} = useI18n()

	// 当前日期索引和日期标签
	const currentDateIndex = ref(2)
	const dateTabs = [t('monday'), t('tuesday'), t('wednesday'), t('thursday'), t('friday'), t('saturday'), t('sunday')]

	// 圆形图和营养数据
	const nutritionData = ref([{
			label: t('carbohydrates'),
			percentage: 80,
			value: '5934 / 261g',
			color: '#4CAF50'
		},
		{
			label: t('protein'),
			percentage: 90,
			value: '660 / 74g',
			color: '#FFC107'
		},
		{
			label: t('fat'),
			percentage: 70,
			value: '84 / 58g',
			color: '#FF6F61'
		},
	])

	// 建议运动
	const exerciseSuggestion = '4676 ' + t('minutes_fast_walk')

	// 饮食记录
	const meals = ref([{
			title: t('lunch'),
			totalCalories: 6255,
			items: [{
				name: t('steamed_bun'),
				calories: '6255',
				image: 'https://example.com/steamed_bun.jpg'
			}, ],
		},
		{
			title: t('dinner'),
			totalCalories: 20671,
			items: [{
				name: t('rice'),
				calories: '20671',
				image: 'https://example.com/rice.jpg'
			}, ],
		},
	])

	// 日期改变
	const onDateChange = (index) => {
		currentDateIndex.value = index
		// 动态加载数据可以在此处更新 nutritionData 和 meals
	}
</script>

<style scoped>
	.container {
		display: flex;
		flex-direction: column;
		padding-bottom: 50px;
		background-color: #f5f5f5;
	}

	.header {
		position: relative;
		z-index: 10;
	}

	.date-picker {
		margin: 10px 0;
		background-color: #ffffff;
	}

	.circle-info {
		display: flex;
		justify-content: center;
		align-items: center;
		margin: 20px 0;
	}

	.circle-content {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.large-text {
		font-size: 32px;
		font-weight: bold;
		color: #FF6F61;
	}

	.small-text {
		font-size: 14px;
		color: #888;
	}

	.nutrition-info {
		padding: 10px 20px;
		background-color: #fff;
		border-radius: 10px;
		margin: 10px;
	}

	.nutrition-item {
		margin-bottom: 15px;
	}

	.label {
		font-size: 16px;
		font-weight: bold;
		margin-bottom: 5px;
	}

	.values {
		font-size: 14px;
		color: #555;
		margin-top: 5px;
	}

	.suggestion {
		margin: 20px 10px;
	}

	.meal-records {
		padding: 10px 20px;
		background-color: #fff;
		border-radius: 10px;
		margin: 10px;
	}

	.meal-header {
		display: flex;
		justify-content: space-between;
		margin-bottom: 10px;
	}

	.meal-title {
		font-size: 16px;
		font-weight: bold;
	}

	.calories {
		font-size: 14px;
		color: #FF6F61;
	}

	.meal-item {
		display: flex;
		align-items: center;
		margin-bottom: 10px;
	}

	.meal-image {
		width: 50px;
		height: 50px;
		margin-right: 10px;
		border-radius: 5px;
		object-fit: cover;
	}

	.meal-info {
		flex: 1;
	}
</style>