<template>
  <a-space class="w-full" direction="vertical">
    <a-card size="small">
      <div class="text-lg">
        <span> 当前版本号：{{ version }} </span>
      </div>
    </a-card>

    <a-card :body-style="{ padding: '10px' }">
      <a-descriptions :column="3" :data="dependenciesList" bordered title="构建依赖">
      </a-descriptions>
    </a-card>
    <a-card :body-style="{ padding: '10px' }">
      <a-descriptions :column="3" :data="devDependenciesList" bordered title="开发依赖">
      </a-descriptions>
    </a-card>
  </a-space>
</template>

<script lang="ts" setup>
import useAppInfo from '@/hooks/useAppInfo'
import type {DescData} from '@arco-design/web-vue'
import {onMounted, reactive, ref} from 'vue'

const {version, dependencies, devDependencies} = useAppInfo()
const showContact = ref(false)
const dependenciesList = reactive<DescData[]>([])
const devDependenciesList = reactive<DescData[]>([])
onMounted(() => {
  const depValues = Object.values(dependencies)
  Object.keys(dependencies).map((it, index) => {
    dependenciesList.push({
      label: it,
      value: depValues[index],
    })
  })
  const devDepValues = Object.values(devDependencies)
  Object.keys(devDependencies).map((it, index) => {
    devDependenciesList.push({
      label: it,
      value: devDepValues[index],
    })
  })
})
</script>
