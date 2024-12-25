<template>
  <view class="recipe-page">
    <image
      :src="imageUrl"
      alt="Recipe Image"
      class="recipe-image"
      mode="widthFix"
    />
    <text class="recipe-name">{{ name }}</text>
    <view class="ingredients-container">
      <view
        v-for="(ingredient, index) in parsedIngredients"
        :key="index"
        :class="['ingredient-tag', `tag-color-${index % tagColors.length}`]"
      >
        {{ ingredient }}
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import {ref} from 'vue'
import {onLoad as uniOnLoad} from '@dcloudio/uni-app'
import {useI18n} from 'vue-i18n'

/* ----------------- i18n ----------------- */
const {t} = useI18n()

/* ----------------- Reactive State ----------------- */
const name = ref('')
const ingredients = ref('')
const imageUrl = ref('')
const parsedIngredients = ref([])

/* ----------------- Tag Colors ----------------- */
const tagColors = ['red', 'green', 'blue', 'orange', 'purple', 'cyan']

/* ----------------- Methods ----------------- */
function parseIngredients(ingredientsStr) {
  try {
    const ingredientsData = JSON.parse(ingredientsStr)
    if (Array.isArray(ingredientsData)) {
      return ingredientsData
    } else if (typeof ingredientsData === 'object' && ingredientsData !== null) {
      return Object.keys(ingredientsData).map(key => t(key))
    } else {
      return [t('ingredients_unavailable')]
    }
  } catch (e) {
    return [t('ingredients_unavailable')]
  }
}

/* ----------------- onLoad ----------------- */
uniOnLoad(options => {
  name.value = decodeURIComponent(options.name || '')
  ingredients.value = decodeURIComponent(options.ingredients || '')
  imageUrl.value = decodeURIComponent(options.image_url || '')
  parsedIngredients.value = parseIngredients(ingredients.value)
})
</script>

<style scoped>
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

.recipe-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.recipe-image {
  width: 100%;
  max-width: 400px;
  border: 2px solid #ccc;
  border-radius: 8px;
  margin-bottom: 20px;
  object-fit: contain;
}

.recipe-name {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 15px;
  text-align: center;
}

.ingredients-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.ingredient-tag {
  padding: 5px 10px;
  border-radius: 15px;
  color: #fff;
  font-size: 14px;
  background-color: #666;
}

.tag-color-0 {
  background-color: #f44336;
}

.tag-color-1 {
  background-color: #4caf50;
}

.tag-color-2 {
  background-color: #2196f3;
}

.tag-color-3 {
  background-color: #ff9800;
}

.tag-color-4 {
  background-color: #9c27b0;
}

.tag-color-5 {
  background-color: #00bcd4;
}
</style>