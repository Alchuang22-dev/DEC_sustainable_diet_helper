<template>
  <view class="recommendation-section">
    <text class="recommend-title">{{ t('recommendation_title') }}</text>

    <view class="dishes">
      <view
        class="dish"
        v-for="(dish, index) in dishes"
        :key="index"
        :class="'fade-in-up delay-' + (index + 1)"
      >
        <image :src="dish.image" :alt="dish.name" class="dish-image" />
        <view class="dish-title">{{ dish.name }}</view>
        <view class="dish-actions">
          <button
            :class="['like-button', { liked: dish.liked }]"
            @click="likeDish(index)"
          >
            <span v-if="dish.liked">❤️</span>
            <span v-else>🤍</span>
          </button>
          <button class="delete-button" @click="deleteDish(index)">🗑️</button>
        </view>
      </view>
    </view>

    <!-- 生成菜谱按钮 -->
    <view class="button-container">
      <button class="generate-button fade-in-up delay-7" @click="regetRecipe">
        {{ t('change_food') }}
      </button>
      <button class="generate-button fade-in-up delay-7" @click="generateRecipe">
        {{ t('generate_recipe') }}
      </button>
    </view>
  </view>

  <!-- 推荐菜谱列表 -->
  <view class="recipe-boxes" v-if="showRecipeBoxes">
    <view
      class="box fade-in-up delay-6"
      v-for="(recipe, index) in recommendedRecipes"
      :key="recipe.recipe_id"
      @click="goToRecipe(index)"
    >
      <image
        :src="recipe.image_url"
        :alt="recipe.name"
        class="box-image"
      />
      <view class="box-description">
        <text class="box-title">{{ recipe.name }}</text>
        <text class="box-text">{{ parseIngredients(recipe.ingredients) }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useFoodListStore } from '../stores/food_list'
import { useUserStore } from '@/stores/user.js'

/* ----------------- Setup ----------------- */
const { t } = useI18n()
const foodStore = useFoodListStore()
const userStore = useUserStore()

const BASE_URL = ref('http://122.51.231.155:8095')
const token = ref(userStore.user.token)

/* ----------------- Reactive & State ----------------- */
const recommendedRecipes = ref([])
const showRecipeBoxes = ref(false)

const dishes = ref([])
const availableNewDishes = ref([])

const likedIngredients = []
const dislikedIngredients = []

/* ----------------- Lifecycle ----------------- */
onMounted(() => {
  fetchRecommendedDishes()
})

/* ----------------- Methods ----------------- */
async function fetchRecommendedDishes() {
  try {
    const response = await uni.request({
      url: `${BASE_URL.value}/ingredients/recommend`,
      method: 'POST',
      header: {
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      },
      data: {
        use_last_ingredients: true,
        liked_ingredients: likedIngredients,
        disliked_ingredients: dislikedIngredients
      }
    })
    if (response.statusCode === 200) {
      const recommendedIngredients = response.data.recommended_ingredients
      const firstSixIngredients = recommendedIngredients.slice(0, 6)
      if (firstSixIngredients.length === 0) {
        regetRecipe()
      } else {
        dishes.value = firstSixIngredients.map(ingredient => ({
          id: ingredient.id,
          name: t(ingredient.name),
          image: ingredient.image_url,
          liked: false
        }))
        availableNewDishes.value = recommendedIngredients.slice(6).map(ing => ({
          id: ing.id,
          name: t(ing.name),
          image: ing.image_url,
          liked: false
        }))
      }
    } else {
      regetRecipe()
    }
  } catch (error) {
    console.error('请求失败:', error)
  }
}

async function regetRecipe() {
  try {
    dishes.value = []
    availableNewDishes.value = []
    const response = await uni.request({
      url: `${BASE_URL.value}/ingredients/recommend`,
      method: 'POST',
      header: {
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      },
      data: {
        use_last_ingredients: false,
        liked_ingredients: likedIngredients,
        disliked_ingredients: dislikedIngredients
      }
    })
    if (response.statusCode === 200 && response.data.recommended_ingredients) {
      const recommendedIngredients = response.data.recommended_ingredients
      dishes.value = recommendedIngredients.slice(0, 6).map(ingredient => ({
        id: ingredient.id,
        name: t(ingredient.name),
        image: ingredient.image_url,
        liked: false
      }))
      availableNewDishes.value = recommendedIngredients.slice(6).map(ing => ({
        id: ing.id,
        name: t(ing.name),
        image: ing.image_url,
        liked: false
      }))
    }
  } catch (error) {
    console.error('重新请求异常:', error)
  }
}

async function generateRecipe() {
  try {
    const selectedIngredients = dishes.value.map(dish => dish.id)

    // 步骤1：保存用户选择的食材
    const setResponse = await uni.request({
      url: `${BASE_URL.value}/ingredients/set`,
      method: 'POST',
      header: {
        Authorization: `Bearer ${token.value}`,
        'Content-Type': 'application/json'
      },
      data: {
        selected_ingredients: selectedIngredients
      }
    })

    // 步骤2：根据用户选择的食材推荐菜谱
    if (setResponse.statusCode === 200) {
      const recommendResponse = await uni.request({
        url: `${BASE_URL.value}/recipes/recommend`,
        method: 'POST',
        header: {
          Authorization: `Bearer ${token.value}`,
          'Content-Type': 'application/json'
        },
        data: {
          selected_ingredients: selectedIngredients,
          disliked_ingredients: dislikedIngredients
        }
      })

      if (recommendResponse.statusCode === 200 && recommendResponse.data.recommended_recipes) {
        recommendedRecipes.value = recommendResponse.data.recommended_recipes
        showRecipeBoxes.value = true
      }
    }
  } catch (error) {
    console.error('生成食谱异常:', error)
  }
}

function parseIngredients(ingredientsStr) {
  try {
    const ingredients = JSON.parse(ingredientsStr)
    if (Array.isArray(ingredients)) {
      return ingredients.length > 5
        ? ingredients.slice(0, 5).join(', ') + ', ...'
        : ingredients.join(', ')
    } else if (typeof ingredients === 'object' && ingredients !== null) {
      const keys = Object.keys(ingredients)
      return keys.length > 5
        ? keys.slice(0, 5).map(key => t(key)).join(', ') + ', ...'
        : keys.map(key => t(key)).join(', ')
    } else {
      return t('ingredients_unavailable')
    }
  } catch (e) {
    return t('ingredients_unavailable')
  }
}

function goToRecipe(index) {
  const recipe = recommendedRecipes.value[index]
  uni.navigateTo({
    url: `/pagesTool/recipe/recipe?name=${encodeURIComponent(recipe.name)}&ingredients=${encodeURIComponent(recipe.ingredients)}&image_url=${encodeURIComponent(recipe.image_url)}`
  })
}

function likeDish(index) {
  dishes.value[index].liked = !dishes.value[index].liked
  likedIngredients.push(dishes.value[index].id)
}

async function deleteDish(index) {
  dislikedIngredients.push(dishes.value[index].id)
  const removedDish = dishes.value.splice(index, 1)[0]
  await simulateBackendDelete(removedDish)
  const newDish = await simulateFetchNewDish()
  dishes.value.push(newDish)
}

function simulateBackendDelete(dish) {
  return new Promise(resolve => {
    // 模拟后端删除
    resolve()
  })
}

function simulateFetchNewDish() {
  return new Promise(resolve => {
    if (availableNewDishes.value.length === 0) {
      resolve({
        name: t('default_dish'),
        image: 'https://cdn.pixabay.com/photo/2016/11/18/14/40/pasta-1836457_1280.jpg',
        liked: false
      })
      return
    }
    const randomIndex = Math.floor(Math.random() * availableNewDishes.value.length)
    const newDish = availableNewDishes.value.splice(randomIndex, 1)[0]
    resolve(newDish)
  })
}
</script>

<style scoped>
:root {
  --primary-color: #4caf50;
  --secondary-color: #2fc25b;
  --background-color: #f5f5f5;
  --card-background: rgba(255, 255, 255, 0.8);
  --text-color: #333;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --font-size-title: 32rpx;
  --font-size-subtitle: 24rpx;
  --transition-duration: 0.5s;
}

.recommendation-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: rgba(76, 175, 80, 0.1);
  backdrop-filter: blur(2rpx);
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
  color: #e91e63;
}
.like-button.liked {
  color: #ff4081;
}

.delete-button {
  color: #f44336;
}

.button-container {
  display: flex;
  justify-content: space-between;
  width: 75%;
  margin-top: 20px;
  gap: 20rpx;
}

.generate-button {
  background-color: var(--primary-color);
  color: #ffffff;
  padding: 20rpx 40rpx;
  margin-top: 20px;
  border: none;
  border-radius: 30rpx;
  font-size: 32rpx;
  cursor: pointer;
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.5s forwards;
  width: auto;
  margin: 0;
}

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
  margin-bottom: 20rpx;
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

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

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
.delay-7 {
  animation-delay: 2.1s;
}
</style>