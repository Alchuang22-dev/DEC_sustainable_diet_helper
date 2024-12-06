<template>
	<view class="container">
		<!-- 全屏背景图片 -->
		<image src="/static/images/index/background_index_new.png" class="background-image"></image>

		<!-- 头部标题 -->
		<view v-if="family.status === FamilyStatus.JOINED" class="header">
			<text class="header-title">{{ family.name || $t('default_family_name') }}</text>
		</view>

		<!-- 家庭ID -->
		<text v-if="family.status === FamilyStatus.JOINED" class="list-title">{{ $t('family_info') + family.id }}</text>

		<!-- 用户未加入家庭时的视图 -->
		<view v-if="family.status === FamilyStatus.NOT_JOINED" class="no-family-view">
			<image src="https://cdn.pixabay.com/photo/2017/01/13/02/31/family-1976162_1280.png" class="centered-image">
			</image>
			<view class="family-actions">
				<button class="action-button" @click="showCreateFamilyModal = true">{{ t('create_family') }}</button>
				<button class="action-button" @click="showJoinFamilyModal = true">{{ t('join_family') }}</button>
			</view>
		</view>

		<!-- 待审核状态的视图 -->
		<view v-if="family.status === FamilyStatus.PENDING_APPROVAL" class="pending-view">
			<image src="https://cdn.pixabay.com/photo/2021/09/20/22/15/hourglass-6641967_1280.png" class="pending-image"></image>
			<text class="pending-text">待管理员审核</text>
			<button class="cancel-button" @click="handleCancelJoin">退出等待</button>
		</view>

		<!-- 用户已加入家庭时的家庭管理部分 -->
		<view v-if="family.status === FamilyStatus.JOINED" class="family-management">
			<!-- 提出想吃的菜品 -->
			<view class="dish-proposal">
				<text class="section-title">{{ t('propose_dish') }}</text>
				<input v-model="newDish.name" :placeholder="t('dish_name_placeholder')" class="input"></input>
				<picker mode="selector" :range="dishPreferenceLevels" :value="newDish.preference"
					@change="onDishPreferenceChange" class="picker">
					<view>{{ $t('dish_preference') }}: {{ dishPreferenceLevels[newDish.preference] }}</view>
				</picker>
				<button class="submit-button" @click="submitDishProposal">{{ t('submit_proposal') }}</button>
			</view>

			<!-- 家庭成员的提议 -->
			<view class="dish-list">
				<text class="section-title">{{ $t('family_dish_proposals') }}</text>
				<scroll-view class="dish-scroll" scroll-y>
					<view v-for="dish in sortedDishProposals" :key="dish.id" class="dish-item">
						<uni-list>
							<uni-list-item :title="dish.name" :note="dishPreferenceLevels[dish.preference]"
								:rightText="dish.proposer" />
						</uni-list>
					</view>
				</scroll-view>
			</view>

			<!-- 共享家庭成员的五大营养成分达标情况 -->
			<view class="shared-data">
				<text class="section-title">{{ $t('shared_family_data') }}</text>
				<!-- 添加家庭碳排放环形图 -->
				<view class="charts">
					<qiun-data-charts :canvas2d="true" canvas-id="familyCarbonChart" type="ring" :opts="carbonRingOpts"
						:chartData="carbonChartData" />
				</view>
				<!-- 添加家庭五大营养成分达标情况的图表 -->
				<view class="charts">
					<qiun-data-charts :canvas2d="true" canvas-id="familyNutrientChart" type="column"
						:opts="nutrientChartOpts" :chartData="nutrientChartData" />
				</view>
			</view>

			<!-- 家庭成员部分 -->
			<view class="family-info">
				<text class="section-title">{{ $t('family_members') }}</text>
				<view class="family-members">
					<view v-for="member in family.members" :key="member.id" class="member">
						<image :src="member.avatar" class="member-avatar"></image>
						<text class="member-name">{{ `${member.nickname}(${t(member.family_name)})` }}</text>
					</view>
				</view>
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
	import {
		ref,
		computed,
		onMounted,
		onUnmounted,
		reactive
	} from 'vue';
	import {
		useI18n
	} from 'vue-i18n';
	import {
		useFamilyStore
	} from '@/stores/family.js';
	import UniList from "@/uni_modules/uni-list/components/uni-list/uni-list.vue";

	// 国际化
	const {
		t
	} = useI18n();

  const token = uni.getStorageSync('token');
  console.log('这个页面', token);


	// Pinia 状态管理
	const familyStore = useFamilyStore();
	const family = computed(() => familyStore.family);
	const FamilyStatus = familyStore.FamilyStatus;

  // familyStore.reset();

  console.log('家庭', family);
  console.log('家庭状态', family.value.status);

	// 定时器引用
	let statusCheckTimer = null;

	// 新菜品提议
	const newDish = reactive({
		name: '',
		preference: 0,
	});

	// 菜品偏好级别
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

	// 定时检查状态
	const startStatusCheck = () => {
		// 清除可能存在的旧定时器
		if (statusCheckTimer) {
			clearInterval(statusCheckTimer);
		}
		// 设置新的定时器
		statusCheckTimer = setInterval(async () => {
			if (family.value.status === FamilyStatus.PENDING_APPROVAL) {
				await familyStore.getFamilyDetails();
			}
		}, 30000); // 30秒
	};

	// 取消加入申请
	const handleCancelJoin = async () => {
		try {
			await familyStore.cancelJoinRequest();
			uni.showToast({
				title: '已取消申请',
				icon: 'success'
			});
		} catch (error) {
			uni.showToast({
				title: '取消失败',
				icon: 'error'
			});
		}
	};

	// 提交菜品提议
	const submitDishProposal = () => {
		if (newDish.name.trim() === '') {
			uni.showToast({
				title: t('dish_name_required'),
				icon: 'none'
			});
			return;
		}
		familyStore.addDishProposal({
			id: Date.now(),
			name: newDish.name,
			preference: newDish.preference,
			proposer: 'You',
		});
		newDish.name = '';
		newDish.preference = 0;
	};

  const createFamily = async () => {
    if (newFamilyName.value.trim() === '') {
      uni.showToast({
        title: t('family_name_required'),
        icon: 'none'
      });
      return;
    }

    console.log('开始创建家庭...');  // 添加日志
    try {
      // 调用前的状态
      console.log('调用前的family状态:', JSON.stringify(familyStore.family));

      const result = await familyStore.createFamily(newFamilyName.value);
      console.log('createFamily返回结果:', result);  // 添加日志

      newFamilyName.value = '';
      showCreateFamilyModal.value = false;

      // 确认更新后的状态
      console.log('创建家庭成功，更新后的family状态:', JSON.stringify(familyStore.family));

      // 显示成功提示
      uni.showToast({
        title: t('create_family_success'),
        icon: 'success'
      });
    } catch (error) {
      console.error('创建家庭失败，错误详情:', error);  // 添加详细错误日志
      uni.showToast({
        title: t('create_family_failed'),
        icon: 'error'
      });
    }
  };

  // 加入家庭
  const joinFamily = async () => {
    if (joinFamilyId.value.trim() === '') {
      uni.showToast({
        title: t('family_id_required'),
        icon: 'none'
      });
      return;
    }

    try {
      // 先搜索家庭
      const searchResult = await familyStore.searchFamily(joinFamilyId.value);

      // 如果搜索不到家庭
      if (!searchResult || !searchResult.id) {
        uni.showToast({
          title: t('family_not_found'),
          icon: 'error'
        });
        return;
      }

      // 搜索到家庭后，调用加入接口
      await familyStore.joinFamily(searchResult.id);
      joinFamilyId.value = '';
      showJoinFamilyModal.value = false;
      startStatusCheck(); // 开始定时检查状态

      uni.showToast({
        title: t('join_request_sent'),
        icon: 'success'
      });

    } catch (error) {
      console.error('Join family error:', error);
      uni.showToast({
        title: t('join_family_failed'),
        icon: 'error'
      });
    }
  };

	// 家庭五大营养成分达标情况数据
	const nutrientChartData = ref({
		categories: [],
		series: []
	});

	// 图表配置
	const nutrientChartOpts = {
		color: ["#1890FF", "#2FC25B"],
		padding: [15, 0, 0, 0],
		xAxis: {
			disableGrid: true,
			axisLine: true,
		},
		yAxis: {},
		extra: {
			column: {
				type: "stack",
				width: 30,
			},
		},
	};

	// 家庭碳排放环形图数据
	const carbonChartData = ref({
		series: [{
			data: []
		}]
	});

	// 环形图配置
	const carbonRingOpts = {
		rotate: false,
		rotateLock: false,
		color: ["#1890FF", "#91CB74", "#FAC858", "#EE6666", "#73C0DE"],
		padding: [0, 0, 0, 0],
		dataLabel: true,
		enableScroll: false,
		legend: {
			show: true,
			position: "right",
			lineHeight: 25
		},
		title: {
			name: t('carbon_total'),
			fontSize: 15,
			color: "#666666"
		},
		subtitle: {
			name: "",
			fontSize: 25,
			color: "#4CAF50"
		},
		extra: {
			ring: {
				ringWidth: 15,
				activeOpacity: 0.5,
				activeRadius: 10,
				offsetAngle: 0,
				labelWidth: 15,
				border: false,
				borderWidth: 3,
				borderColor: "#FFFFFF"
			}
		}
	};

	// 计算属性：排序后的菜品提议
	const sortedDishProposals = computed(() => {
		if (!family.value.dishProposals) return [];
		return [...family.value.dishProposals].sort((a, b) => b.preference - a.preference);
	});

	// 处理偏好选择变化
	const onDishPreferenceChange = (e) => {
		newDish.preference = parseInt(e.detail.value, 10);
	};

  // 修改管理成员方法
  const manageMembers = () => {
      // 获取当前用户ID
      const currentUserId = userStore.userId; // 假设用户store中有userId字段

      if (familyStore.isAdmin(currentUserId)) {
          // 是管理员，允许跳转
          uni.navigateTo({
              url: '/pagesMy/myFamily/myFamily'
          });
      } else {
          // 不是管理员，显示错误提示
          uni.showToast({
              title: t('not_admin'),
              icon: 'error',
              duration: 2000
          });
      }
  };

	// 生命周期钩子
	onMounted(async () => {
		// 初始化营养成分数据
		nutrientChartData.value.categories = [
			t('energy_unit'),
			t('protein_unit'),
			t('fat_unit'),
			t('carbohydrates_unit'),
			t('sodium_unit'),
		];
		nutrientChartData.value.series = [{
				name: t('user_name1'),
				data: [80, 90, 85, 70, 75],
			},
			{
				name: t('user_name2'),
				data: [100, 100, 100, 100, 120],
			}
		];

		// 初始化碳排放数据
		const memberCarbonData = [{
				name: 'Alice',
				value: 2.5
			},
			{
				name: 'Bob',
				value: 3.0
			},
			{
				name: 'Charlie',
				value: 1.5
			}
		];
		carbonChartData.value.series[0].data = memberCarbonData;
		const totalCarbonEmission = memberCarbonData.reduce((sum, item) => sum + item.value, 0);
		carbonRingOpts.subtitle.name = `${totalCarbonEmission.toFixed(1)}Kg`;

		// 如果当前是待审核状态，启动定时检查
		if (family.value.status === FamilyStatus.PENDING_APPROVAL) {
			startStatusCheck();
		}
	});

	onUnmounted(() => {
		// 清除定时器
		if (statusCheckTimer) {
			clearInterval(statusCheckTimer);
		}
	});
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

  /* 添加待审核状态的样式 */
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

	/* Scroll-View 滚动区域 */
	.dish-scroll {
		max-height: 300rpx;
		/* 固定最大高度 */
		overflow-y: auto;
	}

	/* 菜品项 */
	.dish-item {
		padding: 10rpx;
		border-bottom: 1rpx solid #eee;
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