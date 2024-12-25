// stores/food_list.js
import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const BASE_URL = 'http://xcxcs.uwdjl.cn:8080'

function extractNumber(value) {
  if (typeof value === 'string') {
    const num = parseFloat(value.replace(/[^\d.]/g, ''))
    return isNaN(num) ? 0 : num
  }
  return Number(value) || 0
}

export const useFoodListStore = defineStore('foodList', () => {
  const foodList = reactive([])
  const loaded = ref(false)
  const availableFoods = reactive([])
  const { locale } = useI18n()

  async function fetchAvailableFoods() {
    try {
      const [resEn, resZh] = await Promise.all([
        new Promise((resolve, reject) => {
          uni.request({
            url: `${BASE_URL}/foods/names`,
            method: 'GET',
            data: { lang: 'en' },
            success: resolve,
            fail: reject
          })
        }),
        new Promise((resolve, reject) => {
          uni.request({
            url: `${BASE_URL}/foods/names`,
            method: 'GET',
            data: { lang: 'zh' },
            success: resolve,
            fail: reject
          })
        })
      ])

      if (resEn.statusCode === 200 && resZh.statusCode === 200) {
        const enFoods = resEn.data
        const zhFoods = resZh.data
        const foodMap = {}
        enFoods.forEach(f => {
          foodMap[f.id] = {
            id: f.id,
            name_en: f.name,
            name_zh: '',
            image_url: f.image_url
          }
        })
        zhFoods.forEach(f => {
          if (foodMap[f.id]) {
            foodMap[f.id].name_zh = f.name
          } else {
            foodMap[f.id] = {
              id: f.id,
              name_en: '',
              name_zh: f.name,
              image_url: f.image_url
            }
          }
        })
        const combinedFoods = Object.values(foodMap).filter(f => f.name_en || f.name_zh)
        availableFoods.splice(0, availableFoods.length, ...combinedFoods)
      } else {
        uni.showToast({
          title: '获取食物列表失败',
          icon: 'none',
          duration: 2000
        })
      }
    } catch (err) {
      console.error('网络错误:', err)
      availableFoods.splice(0, availableFoods.length,
        {
          id: 1,
          name_en: 'Apple',
          name_zh: '苹果',
          image_url: 'https://example.com/apple.jpg'
        },
        {
          id: 2,
          name_en: 'Banana',
          name_zh: '香蕉',
          image_url: 'https://example.com/banana.jpg'
        }
      )
      uni.showToast({
        title: '网络错误，无法获取食物列表',
        icon: 'none',
        duration: 2000
      })
    }
  }

  function loadFoodList() {
    const storedFoodList = uni.getStorageSync('foodDetails')
    if (storedFoodList && storedFoodList.length > 0) {
      foodList.splice(0, foodList.length, ...storedFoodList.map(food => ({ ...food, isAnimating: false })))
    }
    loaded.value = true
  }

  function saveFoodList() {
    uni.setStorageSync('foodDetails', foodList)
  }

  function addFood(newFood) {
    foodList.push(newFood)
    saveFoodList()
  }

  function deleteFood(index) {
    foodList.splice(index, 1)
    saveFoodList()
  }

  function updateFood(index, updatedFood) {
    if (foodList[index]) {
      Object.assign(foodList[index], updatedFood)
      saveFoodList()
    }
  }

  function getFoodName(id) {
    const food = availableFoods.find(f => f.id === id)
    if (food) {
      return locale.value === 'zh-Hans' ? food.name_zh : food.name_en
    }
    return ''
  }

  function getFoodImageUrl(id) {
    const food = availableFoods.find(f => f.id === id)
    return food ? food.image_url : ''
  }

  async function calculateNutritionAndEmission() {
    try {
      const requestData = foodList.map(food => ({
        id: Number(food.id),
        price: extractNumber(food.price),
        weight: extractNumber(food.weight)
      }))
      if (!requestData.length) {
        uni.showToast({
          title: '请先添加食物',
          icon: 'none',
          duration: 2000
        })
        return
      }
      const response = await new Promise((resolve, reject) => {
        uni.request({
          url: `${BASE_URL}/foods/calculate`,
          method: 'POST',
          data: requestData,
          header: { 'Content-Type': 'application/json' },
          success: resolve,
          fail: reject
        })
      })
      if (response.statusCode === 200) {
        response.data.forEach(item => {
          const food = foodList.find(f => Number(f.id) === item.id)
          if (food) {
            food.emission = item.emission
            food.calories = item.calories
            food.protein = item.protein
            food.fat = item.fat
            food.carbohydrates = item.carbohydrates
            food.sodium = item.sodium
          }
        })
        saveFoodList()
      } else {
        uni.showToast({
          title: '计算失败',
          icon: 'none',
          duration: 2000
        })
      }
    } catch (err) {
      console.error('请求失败', err)
      uni.showToast({
        title: '请求失败',
        icon: 'none',
        duration: 2000
      })
    }
  }

  if (!loaded.value) {
    loadFoodList()
  }

  return {
    foodList,
    loaded,
    availableFoods,
    fetchAvailableFoods,
    loadFoodList,
    saveFoodList,
    addFood,
    deleteFood,
    updateFood,
    getFoodName,
    getFoodImageUrl,
    calculateNutritionAndEmission
  }
})