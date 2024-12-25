// carbon_and_nutrition_data.js
import { defineStore } from 'pinia'
import { reactive, watch } from 'vue'
import { useUserStore } from './user.js'

const BASE_URL = 'http://122.51.231.155:8095'
const STORAGE_KEY = 'carbon_and_nutrition_store_data'

// 递归遍历，对数值字段四舍五入到 1 位小数
function roundNumbers(data) {
  if (Array.isArray(data)) {
    return data.map(item => roundNumbers(item))
  } else if (data !== null && typeof data === 'object') {
    const roundedObj = {}
    for (const key in data) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        roundedObj[key] = roundNumbers(data[key])
      }
    }
    return roundedObj
  } else if (typeof data === 'number') {
    return Math.round(data * 10) / 10
  } else {
    return data
  }
}

// 封装 request + 401 刷新逻辑
function request(config) {
  const userStore = useUserStore()
  return new Promise((resolve, reject) => {
    uni.request({
      ...config,
      success: async res => {
        if (res.statusCode === 401) {
          try {
            await userStore.refreshToken()
            uni.request({
              ...config,
              header: {
                ...config.header,
                Authorization: `Bearer ${userStore.user.token}`
              },
              success: res2 => {
                if (res2.statusCode === 401) reject(new Error('Unauthorized'))
                else resolve(res2)
              },
              fail: err2 => reject(err2)
            })
          } catch (error) {
            uni.navigateTo({ url: '/pagesMy/login/login' })
            reject(new Error('Unauthorized'))
          }
        } else {
          resolve(res)
        }
      },
      fail: err => reject(err)
    })
  })
}

export const useCarbonAndNutritionStore = defineStore('carbon_and_nutrition', () => {
  function getInitialState() {
    try {
      const storedData = uni.getStorageSync(STORAGE_KEY)
      return storedData
        ? JSON.parse(storedData)
        : {
            nutritionGoals: [],
            carbonGoals: [],
            nutritionIntakes: [],
            carbonIntakes: [],
            sharedNutritionCarbonIntakes: []
          }
    } catch {
      return {
        nutritionGoals: [],
        carbonGoals: [],
        nutritionIntakes: [],
        carbonIntakes: [],
        sharedNutritionCarbonIntakes: []
      }
    }
  }

  const state = reactive(getInitialState())

  function saveToStorage() {
    try {
      uni.setStorageSync(STORAGE_KEY, JSON.stringify(state))
    } catch {}
  }

  function watchState() {
    const keys = [
      'nutritionGoals',
      'carbonGoals',
      'nutritionIntakes',
      'carbonIntakes',
      'sharedNutritionCarbonIntakes'
    ]
    keys.forEach(key => {
      watch(
        () => state[key],
        () => {
          saveToStorage()
        },
        { deep: true }
      )
    })
  }

  function createRequestConfig(config) {
    const userStore = useUserStore()
    return {
      ...config,
      header: {
        Authorization: `Bearer ${userStore.user.token || ''}`,
        'Content-Type': 'application/json',
        ...(config.header || {})
      }
    }
  }

  // 设置营养目标
  async function setNutritionGoals(goals) {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/nutrition/goals`,
          method: 'POST',
          data: goals
        })
      )
      await getNutritionGoals()
      return response.data
    } catch (error) {
      throw error
    }
  }
  // 获取营养目标
  async function getNutritionGoals() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/nutrition/goals`,
          method: 'GET'
        })
      )

      if (response.statusCode === 200 && response.data?.data) {
        state.nutritionGoals = roundNumbers(response.data.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }

  // 设置碳排放目标
  async function setCarbonGoals(goals) {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/carbon/goals`,
          method: 'POST',
          data: goals
        })
      )
      await getCarbonGoals()
      return response.data
    } catch (error) {
      throw error
    }
  }
  // 获取碳排放目标
  async function getCarbonGoals() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/carbon/goals`,
          method: 'GET'
        })
      )

      if (response.statusCode === 200 && Array.isArray(response.data?.data)) {
        state.carbonGoals = roundNumbers(response.data.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }

  // 获取实际营养摄入
  async function getNutritionIntakes() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/nutrition/intakes`,
          method: 'GET'
        })
      )
      if (response.statusCode === 200 && Array.isArray(response.data)) {
        state.nutritionIntakes = roundNumbers(response.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }
  // 获取实际碳排放摄入
  async function getCarbonIntakes() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/carbon/intakes`,
          method: 'GET'
        })
      )
      if (response.statusCode === 200 && Array.isArray(response.data)) {
        state.carbonIntakes = roundNumbers(response.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }

  // 共享营养碳排放
  async function setSharedNutritionCarbonIntake(sharedData) {
    try {
      const rounded = roundNumbers(sharedData)
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/shared/nutrition-carbon`,
          method: 'POST',
          data: rounded
        })
      )
      await getSharedNutritionCarbonIntakes()
      return response.data
    } catch (error) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error)
      } else {
        throw error
      }
    }
  }
  // 获取共享营养碳排放
  async function getSharedNutritionCarbonIntakes() {
    try {
      const response = await request(
        createRequestConfig({
          url: `${BASE_URL}/nutrition-carbon/shared/nutrition-carbon`,
          method: 'GET'
        })
      )
      if (response.statusCode === 200 && Array.isArray(response.data)) {
        state.sharedNutritionCarbonIntakes = roundNumbers(response.data)
      }
      return response.data
    } catch (error) {
      throw error
    }
  }

  // 重置
  function reset() {
    state.nutritionGoals = []
    state.carbonGoals = []
    state.nutritionIntakes = []
    state.carbonIntakes = []
    state.sharedNutritionCarbonIntakes = []
    clearStorage()
  }
  // 清除本地存储
  function clearStorage() {
    try {
      uni.removeStorageSync(STORAGE_KEY)
    } catch {}
  }

  // getDataByDate
  function getDataByDate(dateString) {
    const nutritionGoal = state.nutritionGoals.find(g => g.date.startsWith(dateString))
    const carbonGoal = state.carbonGoals.find(g => g.date.startsWith(dateString))
    const dailyNutritionIntakes = state.nutritionIntakes.filter(i =>
      i.date.startsWith(dateString)
    )
    const dailyCarbonIntakes = state.carbonIntakes.filter(i =>
      i.date.startsWith(dateString)
    )
    const meals = {
      breakfast: {},
      lunch: {},
      dinner: {},
      other: {}
    }
    const totalNutrients = {
      calories: 0,
      protein: 0,
      fat: 0,
      carbohydrates: 0,
      sodium: 0
    }
    let totalCarbonEmission = 0

    for (const intake of dailyNutritionIntakes) {
      const mealType = intake.meal_type || 'other'
      if (!meals[mealType].nutrients) {
        meals[mealType].nutrients = {
          calories: 0,
          protein: 0,
          fat: 0,
          carbohydrates: 0,
          sodium: 0
        }
      }
      meals[mealType].nutrients.calories += intake.calories || 0
      meals[mealType].nutrients.protein += intake.protein || 0
      meals[mealType].nutrients.fat += intake.fat || 0
      meals[mealType].nutrients.carbohydrates += intake.carbohydrates || 0
      meals[mealType].nutrients.sodium += intake.sodium || 0

      totalNutrients.calories += intake.calories || 0
      totalNutrients.protein += intake.protein || 0
      totalNutrients.fat += intake.fat || 0
      totalNutrients.carbohydrates += intake.carbohydrates || 0
      totalNutrients.sodium += intake.sodium || 0
    }

    for (const cIntake of dailyCarbonIntakes) {
      const mealType = cIntake.meal_type || 'other'
      if (!meals[mealType].carbonEmission) {
        meals[mealType].carbonEmission = 0
      }
      meals[mealType].carbonEmission += cIntake.emission || 0
      totalCarbonEmission += cIntake.emission || 0
    }

    // 四舍五入
    for (const key in totalNutrients) {
      totalNutrients[key] = Math.round(totalNutrients[key] * 10) / 10
    }
    totalCarbonEmission = Math.round(totalCarbonEmission * 10) / 10
    for (const meal in meals) {
      if (!meals[meal].nutrients) {
        meals[meal].nutrients = {
          calories: 0,
          protein: 0,
          fat: 0,
          carbohydrates: 0,
          sodium: 0
        }
      } else {
        for (const key in meals[meal].nutrients) {
          meals[meal].nutrients[key] = Math.round(
            meals[meal].nutrients[key] * 10
          ) / 10
        }
      }
      if (!meals[meal].carbonEmission) {
        meals[meal].carbonEmission = 0
      } else {
        meals[meal].carbonEmission = Math.round(meals[meal].carbonEmission * 10) / 10
      }
    }

    if (
      !nutritionGoal &&
      dailyNutritionIntakes.length === 0 &&
      dailyCarbonIntakes.length === 0
    ) {
      // 返回空
      return {
        nutrients: {
          actual: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
          target: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
        },
        carbonEmission: {
          actual: 0,
          target: 0
        },
        meals: {
          breakfast: {
            nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
            carbonEmission: 0
          },
          lunch: {
            nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
            carbonEmission: 0
          },
          dinner: {
            nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
            carbonEmission: 0
          },
          other: {
            nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
            carbonEmission: 0
          }
        }
      }
    }

    return {
      nutrients: {
        actual: totalNutrients,
        target: nutritionGoal
          ? {
              calories: Math.round(nutritionGoal.calories * 10) / 10,
              protein: Math.round(nutritionGoal.protein * 10) / 10,
              fat: Math.round(nutritionGoal.fat * 10) / 10,
              carbohydrates: Math.round(nutritionGoal.carbohydrates * 10) / 10,
              sodium: Math.round(nutritionGoal.sodium * 10) / 10
            }
          : { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
      },
      carbonEmission: {
        actual: totalCarbonEmission,
        target: carbonGoal ? Math.round(carbonGoal.emission * 10) / 10 : 0
      },
      meals
    }
  }

  watchState()

  return {
    state,
    setNutritionGoals,
    getNutritionGoals,
    setCarbonGoals,
    getCarbonGoals,
    getNutritionIntakes,
    getCarbonIntakes,
    setSharedNutritionCarbonIntake,
    getSharedNutritionCarbonIntakes,
    reset,
    clearStorage,
    getDataByDate
  }
})