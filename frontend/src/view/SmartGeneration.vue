<template>
  <div class="container mx-auto p-6">
    <h1 class="text-3xl font-bold text-center mb-8">Умная генерация чисел</h1>
    
    <!-- Статус соединения -->
    <div class="mb-4 text-center space-x-2">
      <span 
        :class="connectionStatus.class"
        class="px-3 py-1 rounded-full text-sm font-semibold"
      >
        {{ connectionStatus.text }}
      </span>
      

    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Первый список - все участники -->
      <div class="bg-white rounded-lg shadow-lg p-6">
        <h2 class="text-xl font-bold mb-4">Все участники</h2>
        <div class="space-y-2 max-h-96 overflow-y-auto">
          <div
            v-for="person in allParticipants"
            :key="person.id"
            class="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50"
          >
            <div class="flex-1">
              <div class="font-semibold">{{ person.last_name }} {{ person.first_name }}</div>
              <div class="text-sm text-gray-600">
                Доля: {{ person.fraction }} ({{ person.fraction_from }}-{{ person.fraction_to }})
              </div>
            </div>
            
            <div class="flex items-center space-x-3">
              <!-- Галочка "жив/не жив" -->
              <label class="flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  :checked="person.alive"
                  :disabled="isInGeneration(person.id)"
                  @change="toggleAliveStatus(person.id, $event.target.checked)"
                  class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                <span class="ml-1 text-sm">{{ person.alive ? 'Жив' : 'Мёртв' }}</span>
              </label>
              
              <!-- Кнопка добавить в генерацию -->
              <button
                v-if="!isInGeneration(person.id)"
                @click="addToGeneration(person)"
                :disabled="!person.alive"
                class="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
              >
                Добавить
              </button>
              <button
                v-else
                @click="removeFromGeneration(person.id)"
                class="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
              >
                Убрать
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Второй список - участники для генерации -->
      <div class="bg-white rounded-lg shadow-lg p-6">
        <h2 class="text-xl font-bold mb-4">Список для генерации</h2>
        <div class="mb-4 p-3 bg-gray-100 rounded-lg">
          <div class="text-sm text-gray-700">
            <strong>Общий диапазон долей:</strong> 
            {{ generationRange.from }}-{{ generationRange.to }} 
            (всего: {{ generationRange.total }})
          </div>
        </div>
        
        <div class="space-y-2 max-h-96 overflow-y-auto">
          <div
            v-for="(person, index) in generationList"
            :key="person.id"
            class="flex items-center justify-between p-3 border rounded-lg bg-blue-50"
          >
            <div class="flex-1">
              <div class="font-semibold">{{ person.last_name }} {{ person.first_name }}</div>
              <div class="text-sm text-gray-600">
                Доля: {{ person.fraction }} ({{ person.fraction_from }}-{{ person.fraction_to }})
              </div>
            </div>
            
            <div class="flex items-center space-x-2">
              <!-- Кнопки перемещения -->
              <button
                @click="moveUp(index)"
                :disabled="index === 0"
                class="w-8 h-8 flex items-center justify-center bg-gray-200 rounded hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed"
              >
                ↑
              </button>
              <button
                @click="moveDown(index)"
                :disabled="index === generationList.length - 1"
                class="w-8 h-8 flex items-center justify-center bg-gray-200 rounded hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed"
              >
                ↓
              </button>
              <button
                @click="removeFromGeneration(person.id)"
                class="w-8 h-8 flex items-center justify-center bg-red-500 text-white rounded hover:bg-red-600"
              >
                ×
              </button>
            </div>
          </div>
          
          <div v-if="generationList.length === 0" class="text-center text-gray-500 py-8">
            Список пуст. Добавьте участников из левого списка.
          </div>
        </div>
        
        <!-- Кнопка генерации -->
        <div class="mt-6">
          <button
            @click="generateNumber"
            :disabled="generationList.length === 0 || isGenerating"
            class="w-full py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed font-semibold"
          >
            {{ isGenerating ? 'Генерация...' : 'Сгенерировать число' }}
          </button>
        </div>
        
        <!-- Результат генерации -->
        <div v-if="generationResult" class="mt-4 p-4 bg-green-100 rounded-lg">
          <h3 class="font-semibold text-green-800 mb-2">Результат генерации:</h3>
          <div class="text-lg font-bold text-green-900">Число: {{ generationResult.number }}</div>
          <div v-if="generationResult.winner" class="text-sm text-green-700 mt-1">
            Победитель: {{ generationResult.winner.last_name }} {{ generationResult.winner.first_name }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

// Реактивные данные
const allParticipants = ref([])
const originalParticipants = ref([]) // Исходные данные с сервера
const generationList = ref([])
const serverGenerationList = ref([]) // raw list exactly as server sent (ids or objects)
const socket = ref(null)
const isConnected = ref(false)
const isGenerating = ref(false)
const generationResult = ref(null)
const reservedFractions = ref(new Map()) // Сохраненные доли участников в генерации
const manualDead = ref(new Set()) // ids manually toggled by the user to 'dead'

// Статус соединения
const connectionStatus = computed(() => {
  if (isConnected.value) {
    return {
      text: 'Подключено',
      class: 'bg-green-100 text-green-800'
    }
  } else {
    return {
      text: 'Отключено',
      class: 'bg-red-100 text-red-800'
    }
  }
})

// Пересчет отображаемых долей для основного списка
const recalculateDisplayFractions = () => {
  if (!originalParticipants.value.length) return
  // Устанавливаем `allParticipants` только один раз — при первоначальной загрузке
  if (!allParticipants.value.length) {
    allParticipants.value = originalParticipants.value.map(p => ({ ...p }))
    return
  }

  // При последующих вызовах обновляем только поле alive в allParticipants,
  // чтобы внешний вид (порядок/фракции) оставался как при первоначальной загрузке
  const aliveMap = new Map(originalParticipants.value.map(p => [p.id, !!p.alive]))
  allParticipants.value = allParticipants.value.map(p => ({
    ...p,
    alive: aliveMap.has(p.id) ? aliveMap.get(p.id) : p.alive
  }))
}

// Диапазон долей для генерации
const generationRange = computed(() => {
  if (generationList.value.length === 0) {
    return { from: 0, to: 0, total: 0 }
  }

  const total = generationList.value.reduce((sum, p) => sum + (p.fraction || 0), 0)

  return {
    from: total > 0 ? 1 : 0,
    to: total,
    total: total
  }
})

// Проверка, находится ли участник в списке для генерации (смотрим на raw серверный список)
const isInGeneration = (id) => {
  return serverGenerationList.value.some(item => {
    if (typeof item === 'number') return item === id
    const iid = item.row_id || item.id
    return iid === id
  })
}

// Построить `generationList` (UI) из `serverGenerationList` без изменений самого serverGenerationList
const buildGenerationUIList = () => {
  const ui = []
  for (const item of serverGenerationList.value) {
    if (typeof item === 'number') {
      const p = originalParticipants.value.find(x => x.id === item)
      if (p) {
        const snapshot = {
          id: p.id,
          row_id: p.id,
          first_name: p.first_name,
          last_name: p.last_name,
          fraction: p.fraction || 0,
          fraction_from: p.fraction_from || 0,
          fraction_to: p.fraction_to || 0,
          _fromServer: false // indicate this snapshot was built locally from originalParticipants
        }
        reservedFractions.value.set(snapshot.id, snapshot)
        ui.push(snapshot)
      }
    } else if (typeof item === 'object' && item !== null) {
      const id = item.row_id || item.id
      const snapshot = {
        id: id,
        row_id: id,
        first_name: item.first_name || item.firstName || '',
        last_name: item.last_name || item.lastName || '',
        // Preserve server-sent fraction fields exactly
        fraction: typeof item.fraction !== 'undefined' ? item.fraction : 0,
        fraction_from: typeof item.fraction_from !== 'undefined' ? item.fraction_from : (typeof item.fractionFrom !== 'undefined' ? item.fractionFrom : 0),
        fraction_to: typeof item.fraction_to !== 'undefined' ? item.fraction_to : (typeof item.fractionTo !== 'undefined' ? item.fractionTo : 0),
        _fromServer: true // indicate this snapshot comes directly from server and must not be altered
      }
      reservedFractions.value.set(snapshot.id, snapshot)
      ui.push(snapshot)
    }
  }

  generationList.value = ui
  updateGenerationDisplay()
}

// WebSocket соединение
const connectWebSocket = () => {
  const wsUrl = import.meta.env.VITE_SMART_GEN_URL || 'ws://127.0.0.1:1328'

  try {
    socket.value = new WebSocket(`${wsUrl}/ws`)

    socket.value.onopen = () => {
      console.log('WebSocket подключен')
      isConnected.value = true
    }

    socket.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        // Новый формат: { pre_excluded: [ids], excluded: [ids|objects] }
        if (data && typeof data === 'object' && ('pre_excluded' in data || 'excluded' in data)) {
          // Обновляем статусы "жив/мёртв" на основе pre_excluded (если есть исходные данные)
          const pre = Array.isArray(data.pre_excluded) ? data.pre_excluded : []
          const preSet = new Set(pre)

          if (originalParticipants.value && originalParticipants.value.length) {
            originalParticipants.value = originalParticipants.value.map(p => ({
              ...p,
              alive: !preSet.has(p.id)
            }))
          }

          // Обновляем отображаемый список слева
          recalculateDisplayFractions()

          // Обновляем serverGenerationList ровно как пришло и строим UI-список из него
          const exc = Array.isArray(data.excluded) ? data.excluded : []
          serverGenerationList.value = exc.slice()
          buildGenerationUIList()
          return
        }

        // Старый формат: массив участников или массив "drops"
        if (Array.isArray(data)) {
          const arr = data
          const first = arr.length ? arr[0] : null

          // Если элементы содержат row_id -> это ответ-сниппет (drops) для excluded
          if (first && ("row_id" in first || "RowId" in first)) {
            // Сохраняем serverGenerationList ровно как пришло (объекты)
            serverGenerationList.value = arr.slice()
            buildGenerationUIList()
          } else {
            // Иначе — это начальный список участников
            // Попытка нормализовать ключ `id`
            const normalized = arr.map(item => ({ ...item }))
            normalized.sort((a, b) => (a.id || a.Id || 0) - (b.id || b.Id || 0))
            originalParticipants.value = normalized
            recalculateDisplayFractions()
          }
        }
      } catch (error) {
        console.error('Ошибка парсинга сообщения:', error)
      }
    }

    socket.value.onclose = () => {
  console.log('WebSocket отключен', socket.value && socket.value.readyState)
      isConnected.value = false
      // Попытка переподключения через 3 секунды
      setTimeout(connectWebSocket, 3000)
    }

    socket.value.onerror = (error) => {
  console.error('Ошибка WebSocket:', error)
      isConnected.value = false
    }
  } catch (error) {
    console.error('Не удалось создать WebSocket соединение:', error)
    isConnected.value = false
    // Попытка переподключения через 3 секунды
    setTimeout(connectWebSocket, 3000)
  }
}

// Обновление списка для генерации при изменении данных
const updateGenerationList = () => {
  // Перестроить UI-список из serverGenerationList
  buildGenerationUIList()
}

// Переключение статуса "жив/мёртв" вручную
const toggleAliveStatus = async (userId, alive) => {
  // Если участник в списке генерации, предупреждаем пользователя
  if (isInGeneration(userId)) {
    alert('Нельзя изменить статус участника, пока он в списке генерации. Сначала уберите его из списка.')
    return
  }

  // Обновляем локально
  const idx = originalParticipants.value.findIndex(p => p.id === userId)
  if (idx !== -1) {
    originalParticipants.value[idx].alive = !!alive
    // Обновляем отображение
    recalculateDisplayFractions()
    // Обновляем набор вручную помеченных как мёртвые
    if (!alive) {
      manualDead.value.add(userId)
    } else {
      manualDead.value.delete(userId)
    }
    // Отправляем полное состояние (pre_excluded будет содержать только manualDead)
    await sendStateUpdate()
  }
}

// Отправка полного состояния на сервер в виде { pre_excluded: [...], excluded: [...] }
const sendStateUpdate = async () => {
  if (!socket.value || socket.value.readyState !== WebSocket.OPEN) return

  // pre_excluded - те, кого пометили как мёртвых
  // pre_excluded should contain only ids the user manually toggled to dead
  const preExcluded = Array.from(manualDead.value)

  // excluded - список для генерации (будем отправлять ровно как serverGenerationList)
  const excluded = serverGenerationList.value.map(item => {
    if (typeof item === 'number') return Number(item)
    if (item && typeof item === 'object') return Number(item.row_id || item.id)
    return null
  }).filter(x => x !== null)

  const message = {
    pre_excluded: preExcluded,
    excluded: excluded
  }

  try {
  const payload = JSON.stringify(message)
  console.log('WS SEND:', payload, 'readyState=', socket.value.readyState)
  socket.value.send(payload)
  } catch (e) {
    console.error('Ошибка отправки состояния по WebSocket:', e)
  }
}

// Отправить только поле excluded (например при удалении) — ровно как serverGenerationList
const sendExcludedUpdate = async () => {
  if (!socket.value || socket.value.readyState !== WebSocket.OPEN) return

  // Build numeric ids for excluded
  const excluded = serverGenerationList.value.map(item => {
    if (typeof item === 'number') return Number(item)
    if (item && typeof item === 'object') return Number(item.row_id || item.id)
    return null
  }).filter(x => x !== null)

  const message = {
    excluded: excluded
  }

  try {
  const payload = JSON.stringify(message)
  console.log('WS SEND excluded:', payload, 'readyState=', socket.value.readyState)
  socket.value.send(payload)
  } catch (e) {
    console.error('Ошибка отправки excluded по WebSocket:', e)
  }
}

// Добавление участника в список генерации
const addToGeneration = async (person) => {
  if (person.alive && !isInGeneration(person.id)) {
    // Находим текущего участника в отображаемом списке для получения актуальной доли
    const currentPerson = allParticipants.value.find(p => p.id === person.id)

    // Сохраняем полную информацию о доле на момент добавления для правого списка
    const snapshotData = {
      id: person.id,
      row_id: person.id,
      first_name: person.first_name,
      last_name: person.last_name,
      fraction: currentPerson ? currentPerson.fraction : person.fraction || 0,
      fraction_from: currentPerson ? currentPerson.fraction_from : person.fraction_from || 0,
      fraction_to: currentPerson ? currentPerson.fraction_to : person.fraction_to || 0
    }

  // Добавляем в serverGenerationList (как сервер ожидает — id)
  serverGenerationList.value.push(person.id)

  // Сохраняем снимок и перестраиваем UI
  reservedFractions.value.set(person.id, snapshotData)
  buildGenerationUIList()

  // Локально помечаем участника как мёртвого
  const idx = originalParticipants.value.findIndex(p => p.id === person.id)
  if (idx !== -1) originalParticipants.value[idx].alive = false

  // Отправляем полное состояние на сервер
  await sendStateUpdate()
  }
}

// Удаление участника из списка генерации
const removeFromGeneration = async (userId) => {
  const index = generationList.value.findIndex(p => p.id === userId)
  if (index !== -1) {
    // Удаляем из serverGenerationList
    serverGenerationList.value.splice(index, 1)

    // Удаляем сохраненную долю
    reservedFractions.value.delete(userId)

    // Локально помечаем участника как живого
    const idx = originalParticipants.value.findIndex(p => p.id === userId)
    if (idx !== -1) originalParticipants.value[idx].alive = true

  // Перестроим UI и отправим только excluded (без удалённого id)
  buildGenerationUIList()
  await sendExcludedUpdate()

    // Пересчитываем отображаемые доли
    recalculateDisplayFractions()
  }
}

// Перемещение участника вверх
const moveUp = (index) => {
  if (index > 0) {
    const temp = serverGenerationList.value[index]
    serverGenerationList.value[index] = serverGenerationList.value[index - 1]
    serverGenerationList.value[index - 1] = temp
    buildGenerationUIList()
    // Отправляем обновлённый порядок на сервер
    sendStateUpdate()
  }
}

// Перемещение участника вниз
const moveDown = (index) => {
  if (index < serverGenerationList.value.length - 1) {
    const temp = serverGenerationList.value[index]
    serverGenerationList.value[index] = serverGenerationList.value[index + 1]
    serverGenerationList.value[index + 1] = temp
    buildGenerationUIList()
    // Отправляем обновлённый порядок на сервер
    sendStateUpdate()
  }
}

// Обновление отображения списка генерации с пересчетом диапазонов для генерации
const updateGenerationDisplay = () => {
  let cumulativeFraction = 0

  for (let i = 0; i < generationList.value.length; i++) {
    const p = generationList.value[i]
    // If this item came from server, preserve its ranges and fraction exactly
    if (p._fromServer) {
      cumulativeFraction += p.fraction || 0
      // ensure reserved snapshot kept
      if (reservedFractions.value.has(p.id)) {
        const snap = reservedFractions.value.get(p.id)
        snap.fraction = p.fraction
        snap.fraction_from = p.fraction_from
        snap.fraction_to = p.fraction_to
        reservedFractions.value.set(p.id, snap)
      } else {
        reservedFractions.value.set(p.id, { ...p })
      }
      continue
    }

    const frac = p.fraction || 0
    const from = frac > 0 ? cumulativeFraction + 1 : 0
    const to = frac > 0 ? cumulativeFraction + frac : 0

    p.fraction_from = from
    p.fraction_to = to

    cumulativeFraction += frac

    // Обновляем снимок в reservedFractions
    if (reservedFractions.value.has(p.id)) {
      const snap = reservedFractions.value.get(p.id)
      snap.fraction = p.fraction
      snap.fraction_from = p.fraction_from
      snap.fraction_to = p.fraction_to
      reservedFractions.value.set(p.id, snap)
    } else {
      reservedFractions.value.set(p.id, { ...p })
    }
  }
}

// Генерация числа
const generateNumber = async () => {
  if (generationList.value.length === 0) return
  
  isGenerating.value = true
  generationResult.value = null
  
  try {
    // Получаем диапазон для генерации
    const range = generationRange.value
    if (range.total === 0) {
      throw new Error('Нет активных участников для генерации')
    }
    
    // Генерируем случайное число от 1 до общего количества долей
    const randomNumber = Math.floor(Math.random() * range.total) + 1
    
    // Находим победителя по сохраненным долям
    let currentSum = 0
    let winner = null
    
    for (const person of generationList.value) {
      const personFraction = person.fraction || 0
      currentSum += personFraction
      if (randomNumber <= currentSum && personFraction > 0) {
        winner = person
        break
      }
    }
    
    generationResult.value = {
      number: randomNumber,
      winner: winner
    }
  } catch (error) {
    console.error('Ошибка генерации:', error)
    alert('Ошибка при генерации числа: ' + error.message)
  } finally {
    isGenerating.value = false
  }
}

// Lifecycle hooks
onMounted(() => {
  connectWebSocket()
})

onUnmounted(() => {
  if (socket.value) {
    socket.value.close()
  }
})
</script>

<style scoped>
/* Дополнительные стили при необходимости */
</style>
