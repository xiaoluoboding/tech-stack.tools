<template>
  <div class="lg:sticky top-24 z-10">
    <h2 class="text-3xl font-extrabold text-primary" id="#">
      {{ title }}
    </h2>
    <p v-if="subtitle" class="text-sm text-gray-700 dark:text-gray-300 mb-1">
      {{ subtitle }}
    </p>

    <div
      class="grid grid-cols-1 gap-x-4 gap-y-2 my-4 overflow-x-auto sm:flex-wrap sm:justify-center max-h-[1032px]"
    >
      <button
        v-for="item of items"
        :key="item.id"
        type="button"
        :aria-label="item.name"
        class="w-full px-4 py-3 mb-2 text-left flex items-center rounded-lg cursor-pointer card-bd card-bg text-primary transition-colors duration-150 ease-in-out hover:border-slate-500"
        :class="{
          'border-violet-500 bg-violet-100 dark:bg-violet-500 dark:border-violet-300':
            selectedItem === item.name
        }"
        @click="handleSelect(item)"
      >
        <div class="overflow-hidden flex justify-between w-full">
          <div class="truncate">{{ item.name }}</div>
          <div>({{ item.counts }})</div>
        </div>
        <slot name="badge" v-bind="item" />
        <div class="flex-auto" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Item {
  id: string
  name: string
  color: string
  counts: number
}

defineProps<{
  title: string
  subtitle?: string
  items: Item[]
  selectedItem?: string | null
}>()

const emit = defineEmits<{ (e: 'toggle', any): void }>()

const handleSelect = (item) => {
  emit('toggle', item.name)
}
</script>
