<template>
	<view class="recommendation-section">
	    <text class="recommend-title">{{ $t('recommendation_title') }}</text>
	
	    <view class="dishes">
	        <view class="dish" v-for="(dish, index) in dishes" :key="index"
	            :class="'fade-in-up delay-' + (index + 1)">
	            <image :src="dish.image" :alt="dish.name" class="dish-image"></image>
	            <view class="dish-title">{{ dish.name }}</view>
	            <view class="dish-actions">
	                <button :class="['like-button', { liked: dish.liked }]" @click="likeDish(index)">
	                    <span v-if="dish.liked">❤️</span>
	                    <span v-else>🤍</span>
	                </button>
	                <button class="delete-button" @click="deleteDish(index)">
	                    🗑️
	                </button>
	            </view>
			</view>
	    </view>
		
		<!-- 生成菜谱按钮 -->
		<view class="button-container">
		    <button class="generate-button fade-in-up delay-6" @click="generateRecipe">
		        {{$t('change_food')}}
		    </button>
		    <button class="generate-button fade-in-up delay-6" @click="generateRecipe">
		        {{$t('generate_recipe')}}
		    </button>
		</view>
	</view>
	   <view class="recipe-boxes" v-if="showRecipeBoxes">
	            <view class="box fade-in-up delay-6" @click="goToRecipe('dapanji')">
	                <image src="/static/images/dapanji.png" alt="大盘鸡" class="box-image"></image>
	                <view class="box-description">
	                    <text class="box-title">{{$t('recommended_recipe')}}</text>
	                    <text class="box-text">{{$t('recommended_recipe_info')}}</text>
	                </view>
	            </view>
	        </view>
</template>


<script setup>
import { onMounted, ref, reactive, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useFoodListStore } from '@/stores/food_list'; // 引入 Pinia Store
import { useUserStore } from "@/stores/user.js"; // 引入用户 Store

const { t, locale } = useI18n();

const foodStore = useFoodListStore();
const userStore = useUserStore();

// 定义 BASE_URL 为 ref
const BASE_URL = ref('http://122.51.231.155:8095');

//const user_id = ref('');
// 定义 token 为 computed 属性
const token = computed(() => userStore.user.token);

// 响应式数据
const showRecipeBoxes = ref(false);
const dishes = ref([]);  // 推荐的前6个菜品
const availableNewDishes = ref([]);  // 其他的菜品

// 模拟用户的偏好（可以根据实际情况修改）
const likedIngredients = [];  // 用户喜欢的食材ID
const dislikedIngredients = [];  // 用户不喜欢的食材ID

// 方法：从后端获取推荐食材数据
const fetchRecommendedDishes = async () => {
  try {
    const response = await uni.request({
      url: `${BASE_URL.value}/ingredients/recommend`,
      method: 'POST',
      header: {
		"Authorization": `Bearer ${token.value}`, // 替换为实际的 Token 变量
		"Content-Type": "application/json", // 设置请求类型
      },
      data: {
        use_last_ingredients: true,  // 默认为否
        liked_ingredients: likedIngredients,
        disliked_ingredients: dislikedIngredients,
      },
    })

    // 处理成功响应
    if (response.statusCode === 200 && response[1].data.recommended_ingredients) {
      const recommendedIngredients = response[1].data.recommended_ingredients
      // 将前6个推荐食材放入dishes
      dishes.value = recommendedIngredients.slice(0, 6).map((ingredient) => ({
        name: ingredient.name,
        image: `https://via.placeholder.com/300?text=${ingredient.name}`, // 这里可以根据食材生成图片URL
        liked: false,
      }))
      // 将其余食材放入availableNewDishes
      availableNewDishes.value = recommendedIngredients.slice(6).map((ingredient) => ({
        name: ingredient.name,
        image: `https://via.placeholder.com/300?text=${ingredient.name}`,
        liked: false,
      }))
    } else {
      console.error('获取食材推荐失败:', response[1].data)
    }
  } catch (error) {
    console.error('请求失败:', error)
  }
}

// 方法：生成食谱
const generateRecipe = () => {
  console.log("推荐菜谱");
  showRecipeBoxes.value = true;
}

// 方法：跳转到推荐的食谱
const goToRecipe = (recipeName) => {
  uni.navigateTo({
    url: `/pages/recipes/${recipeName}`,
  })
}

// 喜欢菜品
const likeDish = (index) => {
  dishes.value[index].liked = !dishes.value[index].liked;
}

// 删除菜品
const deleteDish = async (index) => {
  const removedDish = dishes.value.splice(index, 1)[0];
  await simulateBackendDelete(removedDish);
  const newDish = await simulateFetchNewDish();
  dishes.value.push(newDish);
}

// 模拟删除请求
const simulateBackendDelete = (dish) => {
  return new Promise((resolve) => {
    console.log(`Simulating deletion of dish: ${dish.name}`)
    setTimeout(() => {
      resolve();
    }, 1000)
  })
}

// 模拟获取新菜品
const simulateFetchNewDish = () => {
  return new Promise((resolve) => {
    if (availableNewDishes.value.length === 0) {
      resolve({ name: t('default_dish'), image: 'https://cdn.pixabay.com/photo/2016/11/18/14/40/pasta-1836457_1280.jpg', liked: false });
      return;
    }
    const randomIndex = Math.floor(Math.random() * availableNewDishes.value.length);
    const newDish = availableNewDishes.value.splice(randomIndex, 1)[0];
    setTimeout(() => {
      resolve(newDish);
    }, 1000)
  })
}

// 获取推荐的食材
onMounted(() => {
  fetchRecommendedDishes();
})
</script>


<style scoped>
    /* 通用变量 */
    :root {
        --primary-color: #4CAF50;
        --secondary-color: #2fc25b;
        --background-color: #f5f5f5;
        --card-background: rgba(255, 255, 255, 0.8);
        --text-color: #333;
        --shadow-color: rgba(0, 0, 0, 0.1);
        --font-size-title: 32rpx;
        --font-size-subtitle: 24rpx;
        --transition-duration: 0.5s;
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
        z-index: 0;
        opacity: 0.1;
    }

    /* 推荐区域 */
    .recommendation-section {
        display: flex;
        flex-direction: column;
        align-items: center;
        background-color: rgba(76, 175, 80, 0.1);
        /* 半透明绿色背景 */
        backdrop-filter: blur(2rpx);
        /* 高斯模糊 */
        padding: 30rpx;
        margin: 20rpx;
        border-radius: 15rpx;
        box-shadow: 0 4rpx 10rpx var(--shadow-color);
        z-index: 1;
        position: relative;
        animation: fadeInUp 1s ease;
    }

    .recommend-title {
        text-align: center;
        margin-bottom: 15rpx;
        font-size: 28rpx;
        color: var(--primary-color);
        font-weight: bold;
    }

    /* 菜品宫格布局 */
    .dishes {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 20rpx;
        justify-items: center;
        width: 100%;
    }

    .dish {
        width: 100%;
        max-width: 300rpx;
        background-color: rgba(255, 255, 255, 0.9);
        color: #333;
        border-radius: 10rpx;
        overflow: hidden;
        box-shadow: 0 4rpx 10rpx var(--shadow-color);
        display: flex;
        flex-direction: column;
        align-items: center;
        opacity: 0;
        transform: translateY(20px);
        animation: fadeInUp 0.5s forwards;
    }

    .dish-image {
        width: 100%;
        height: 150rpx;
        object-fit: cover;
    }

    .dish-title {
        padding: 10rpx;
        font-size: 28rpx;
        background-color: #ffe082;
        text-align: center;
        font-weight: bold;
        width: 100%;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* 新增的菜品操作按钮 */
    .dish-actions {
        display: flex;
        justify-content: space-around;
        width: 100%;
        padding: 10rpx 0;
        background-color: #f0f0f0;
    }

    .like-button,
    .delete-button {
        background: none;
        border: none;
        font-size: 32rpx;
        cursor: pointer;
        transition: transform 0.2s, color 0.2s;
    }

    .like-button:hover,
    .delete-button:hover {
        transform: scale(1.2);
    }

    .like-button {
        color: #e91e63; /* 喜欢按钮使用粉色 */
    }

    .like-button.liked {
        color: #ff4081; /* 喜欢状态下更深的粉色 */
    }

    .delete-button {
        color: #f44336; /* 删除按钮使用红色 */
    }

    /* 生成菜谱按钮 */
	.button-container {
	    display: flex;
	    justify-content: space-between; /* 按钮左右排布 */
	    width: 75%;
	    gap: 20rpx; /* 按钮之间的间距 */
	}
	
	.generate-button {
		background-color: var(--primary-color);
		color: #ffffff;
		padding: 20rpx 40rpx;
		border: none;
		border-radius: 30rpx;
		font-size: 32rpx;
		cursor: pointer;
		opacity: 0;
		transform: translateY(20px);
		animation: fadeInUp 0.5s forwards;
		width: auto; /* 修改为自适应宽度 */
		margin: 0; /* 去除按钮的默认外边距 */
	}

	.recommendation-section button {
		width: auto;
		margin: 0 10rpx;
	}

    /* 推荐菜谱 */
    .recipe-boxes {
        background-color: rgba(255, 255, 255, 0.9);
        padding: 30rpx;
        border-radius: 20rpx;
        box-shadow: 0 4rpx 15rpx rgba(0, 0, 0, 0.1);
        margin: 30rpx 20rpx;
        animation: fadeInUp 1s ease;
    }

    .box {
        display: flex;
        align-items: center;
        cursor: pointer;
        opacity: 0;
        transform: translateY(20px);
        animation: fadeInUp 0.5s forwards;
    }

    .box-image {
        width: 160rpx;
        height: 160rpx;
        margin-right: 30rpx;
        border-radius: 10rpx;
        object-fit: cover;
    }

    .box-description {
        flex-grow: 1;
    }

    .box-title {
        font-size: 36rpx;
        font-weight: bold;
        color: var(--primary-color);
        margin-bottom: 10rpx;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .box-text {
        font-size: 28rpx;
        color: var(--text-color);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* 动画效果 */
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
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    /* 动画延迟 */
    .delay-1 {
        animation-delay: 0.3s;
    }

    .delay-2 {
        animation-delay: 0.6s;
    }

    .delay-3 {
        animation-delay: 0.9s;
    }

    .delay-4 {
        animation-delay: 1.2s;
    }

    .delay-5 {
        animation-delay: 1.5s;
    }

    .delay-6 {
        animation-delay: 1.8s;
    }

    /* 响应式设计 */
    @media (max-width: 600px) {
        .dec_header {
            flex-direction: column;
            align-items: center;
        }

        .dec_logo {
            width: 80%;
            margin-bottom: 10rpx;
        }

        .title {
            width: 100%;
            text-align: center;
            margin-left: 0;
        }

        .recommendation-section {
            padding: 20rpx;
        }

        .recommend-image {
            width: 150rpx;
            height: 150rpx;
        }

        .dish-title {
            font-size: 24rpx;
        }

        .generate-button {
            font-size: 28rpx;
            padding: 15rpx 30rpx;
        }

        .box-title {
            font-size: 32rpx;
        }

        .box-text {
            font-size: 24rpx;
        }

        .dishes {
            gap: 15rpx;
        }
    }
</style>
